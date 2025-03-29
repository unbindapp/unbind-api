package k8s

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// LogOptions represents options for filtering and streaming logs
type LogOptions struct {
	Namespace     string            // Kubernetes namespace to query
	Labels        map[string]string // Pod labels to filter by
	Since         time.Duration     // Get logs from this time ago (e.g., 1h)
	Tail          int64             // Number of lines to get from the end
	Follow        bool              // Whether to stream logs
	Previous      bool              // Get logs from previous instance
	Timestamps    bool              // Include timestamps
	SinceTime     *metav1.Time      // Get logs from a specific time
	LimitBytes    int64             // Limit bytes of logs returned
	SearchPattern string            // Optional text pattern to grep for
}

type LogMetadata struct {
	// Metadata to stick on
	ServiceID     uuid.UUID `json:"service_id"`
	TeamID        uuid.UUID `json:"team_id"`
	ProjectID     uuid.UUID `json:"project_id"`
	EnvironmentID uuid.UUID `json:"environment_id"`
}

type LogEvents struct {
	// LogEvents is a slice of log events
	Logs []LogEvent `json:"logs" nullable:"false"`
}

// SSE Error
type LogsError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LogEvent represents a log line event sent via SSE
type LogEvent struct {
	PodName   string      `json:"pod_name"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
	Message   string      `json:"message"`
	Metadata  LogMetadata `json:"metadata,omitempty"`
}

// GetPodLogs retrieves logs for a specific pod based on provided options
func (self *KubeClient) GetPodLogs(ctx context.Context, podName string, opts LogOptions, client *kubernetes.Clientset) (string, error) {
	podLogOptions := &corev1.PodLogOptions{
		Follow:     opts.Follow,
		Previous:   opts.Previous,
		Timestamps: opts.Timestamps,
	}

	if opts.Since > 0 {
		podLogOptions.SinceSeconds = utils.ToPtr(int64(opts.Since.Seconds()))
	}

	if opts.SinceTime != nil {
		podLogOptions.SinceTime = opts.SinceTime
	}

	if opts.Tail > 0 {
		podLogOptions.TailLines = &opts.Tail
	}

	if opts.LimitBytes > 0 {
		podLogOptions.LimitBytes = &opts.LimitBytes
	}

	req := client.CoreV1().Pods(opts.Namespace).GetLogs(podName, podLogOptions)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("error opening log stream for pod %s: %v", podName, err)
	}
	defer podLogs.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", fmt.Errorf("error reading log stream for pod %s: %v", podName, err)
	}

	logs := buf.String()

	// Apply search pattern if provided
	if opts.SearchPattern != "" {
		var filteredLines []string
		for _, line := range strings.Split(logs, "\n") {
			if strings.Contains(line, opts.SearchPattern) {
				filteredLines = append(filteredLines, line)
			}
		}
		logs = strings.Join(filteredLines, "\n")
	}

	return logs, nil
}

// StreamPodLogs streams logs from a pod to the provided writer with filtering
func (self *KubeClient) StreamPodLogs(
	ctx context.Context,
	podName, namespace string,
	opts LogOptions,
	meta LogMetadata,
	client *kubernetes.Clientset,
	eventChan chan<- LogEvents,
) error {
	podLogOptions := &corev1.PodLogOptions{
		Follow:     opts.Follow,
		Previous:   opts.Previous,
		Timestamps: opts.Timestamps,
	}
	if opts.Since > 0 {
		podLogOptions.SinceSeconds = utils.ToPtr(int64(opts.Since.Seconds()))
	}
	if opts.SinceTime != nil {
		podLogOptions.SinceTime = opts.SinceTime
	}
	if opts.Tail > 0 {
		podLogOptions.TailLines = &opts.Tail
	}
	if opts.LimitBytes > 0 {
		podLogOptions.LimitBytes = &opts.LimitBytes
	}

	req := client.CoreV1().Pods(namespace).GetLogs(podName, podLogOptions)
	stream, err := req.Stream(ctx)
	if err != nil {
		return fmt.Errorf("error opening log stream for pod %s: %v", podName, err)
	}
	defer stream.Close()

	reader := bufio.NewReader(stream)

	const (
		batchSize    = 100
		maxBatchWait = 200 * time.Millisecond
	)

	// Channels for reading lines
	linesChan := make(chan string)
	readErrChan := make(chan error, 1)

	// Kick off a goroutine to read the lines from the stream
	go func() {
		defer close(linesChan)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				// Send the error and break
				readErrChan <- err
				return
			}
			linesChan <- string(line)
		}
	}()

	batch := make([]LogEvent, 0, batchSize)
	timer := time.NewTimer(maxBatchWait)

	// Function to send the current batch to the channel
	sendBatch := func() {
		var eventsResponse LogEvents
		if len(batch) == 0 {
			return
		} else {
			// Make a copy so we don’t race with subsequent appends
			events := make([]LogEvent, len(batch))
			copy(events, batch)
			eventsResponse = LogEvents{
				Logs: events,
			}
		}

		select {
		case eventChan <- eventsResponse:
			// Successfully sent
		case <-ctx.Done():
			// If context is canceled, just return
		}

		batch = batch[:0]
	}

	// Make sure we start the timer from “now”
	if !timer.Stop() {
		<-timer.C
	}
	timer.Reset(maxBatchWait)

	// Main loop: block on either the timer, the context, or new lines
	for {
		select {
		case <-ctx.Done():
			// Context canceled, flush and quit
			sendBatch()
			return nil

		case err := <-readErrChan:
			// Could be io.EOF or a real error
			if err == io.EOF {
				sendBatch()
				return nil
			}
			return fmt.Errorf("error reading from log stream: %v", err)

		case <-timer.C:
			// Timer fired: flush partial batch
			sendBatch()
			// Reset the timer
			timer.Reset(maxBatchWait)

		case line, ok := <-linesChan:
			// New line from the stream
			if !ok {
				// Channel closed, no more data
				sendBatch()
				return nil
			}

			// Filter if needed
			if opts.SearchPattern != "" && !strings.Contains(line, opts.SearchPattern) {
				continue
			}

			// Extract timestamps if requested
			var timestamp time.Time
			var message string
			if opts.Timestamps {
				parts := strings.SplitN(line, " ", 2)
				if len(parts) == 2 {
					tsCandidate := strings.TrimSpace(parts[0])
					msgCandidate := strings.TrimSpace(parts[1])
					if t, err := time.Parse(time.RFC3339Nano, tsCandidate); err == nil {
						timestamp = t
						message = msgCandidate
					} else {
						// Couldn’t parse, treat entire thing as message
						message = line
					}
				} else {
					message = line
				}
			} else {
				message = line
			}

			batch = append(batch, LogEvent{
				PodName:   podName,
				Timestamp: timestamp,
				Message:   message,
				Metadata:  meta,
			})

			// If the batch is full, send it now and reset timer
			if len(batch) >= batchSize {
				sendBatch()
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(maxBatchWait)
			}
		}
	}
}
