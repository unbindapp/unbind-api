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
func (self *KubeClient) StreamPodLogs(ctx context.Context, podName, namespace string, opts LogOptions, meta LogMetadata, client *kubernetes.Clientset, eventChan chan<- LogEvent) error {
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
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("error reading from log stream: %v", err)
			}

			lineStr := string(line)

			// Apply search pattern filtering if needed
			if opts.SearchPattern != "" && !strings.Contains(lineStr, opts.SearchPattern) {
				continue
			}

			// Parse timestamp if included
			var timestamp time.Time
			var message string
			if opts.Timestamps {
				// K8s timestamp format is typically like: 2023-01-01T12:34:56.123456Z
				parts := strings.SplitN(lineStr, " ", 2)
				if len(parts) == 2 {
					timestamp, err = time.Parse(time.RFC3339Nano, strings.TrimSpace(parts[0]))
					if err == nil {
						message = strings.TrimSpace(parts[1])
					} else {
						message = lineStr
					}
				} else {
					message = lineStr
				}
			} else {
				message = lineStr
			}

			// Send the log event
			select {
			case eventChan <- LogEvent{
				PodName:   podName,
				Timestamp: timestamp,
				Message:   message,
				Metadata:  meta,
			}:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
