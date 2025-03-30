package loki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

// StreamLokiPodLogs streams logs from Loki tail API using WebSocket for multiple pods using a single connection
func (self *LokiLogQuerier) StreamLokiPodLogs(
	ctx context.Context,
	opts LokiLogOptions,
	podMetadataMap map[string]map[string]LogMetadata, // Map of namespace -> podName -> metadata
	eventChan chan<- LogEvents,
) error {
	if len(podMetadataMap) == 0 {
		return fmt.Errorf("no pods specified for log streaming")
	}

	// Build the OR query for multiple pods
	var queryStr string

	// Use the LogQL regex pattern format: {instance=~"(pattern1|pattern2)"}
	var instancePatterns []string

	for namespace, podsMap := range podMetadataMap {
		for podName := range podsMap {
			// Escape special regex characters in the pod name and namespace
			escapedNamespace := strings.ReplaceAll(namespace, ".", "\\.")
			escapedPodName := strings.ReplaceAll(podName, ".", "\\.")
			pattern := fmt.Sprintf("%s/%s:service", escapedNamespace, escapedPodName)
			instancePatterns = append(instancePatterns, pattern)
		}
	}

	if len(instancePatterns) > 0 {
		// Join patterns with | inside parentheses for proper regex OR syntax
		regexPattern := fmt.Sprintf("(%s)", strings.Join(instancePatterns, "|"))
		queryStr = fmt.Sprintf("{instance=~\"%s\"}", regexPattern)
	} else {
		return fmt.Errorf("no valid pod queries could be constructed")
	}

	// Add search pattern if specified
	if opts.SearchPattern != "" {
		queryStr = fmt.Sprintf("(%s) |= \"%s\"", queryStr, opts.SearchPattern)
	}

	// Build the request URL with parameters
	reqURL, err := url.Parse(self.endpoint)
	if err != nil {
		return fmt.Errorf("Unable to parse loki query URL: %v", err)
	}

	// Change protocol from http to ws
	if reqURL.Scheme == "http" {
		reqURL.Scheme = "ws"
	} else if reqURL.Scheme == "https" {
		reqURL.Scheme = "wss"
	}

	q := reqURL.Query()
	q.Set("query", queryStr)

	// Set time range
	if opts.SinceTime != nil {
		q.Set("start", strconv.FormatInt(opts.SinceTime.UnixNano(), 10))
	} else if opts.Since > 0 {
		startTime := time.Now().Add(-opts.Since)
		q.Set("start", strconv.FormatInt(startTime.UnixNano(), 10))
	}

	// Set limit
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}

	// Set follow mode
	if opts.Follow {
		q.Set("tail", "true")
	}

	reqURL.RawQuery = q.Encode()

	log.Infof("Streaming logs with query: %s", queryStr)

	// Create websocket connection
	dialer := websocket.DefaultDialer
	wsConn, resp, err := dialer.DialContext(ctx, reqURL.String(), nil)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("failed to connect to Loki WebSocket: %v, status: %d", err, resp.StatusCode)
		}
		return fmt.Errorf("failed to connect to Loki WebSocket: %v", err)
	}
	defer wsConn.Close()

	// Setup context cancellation
	go func() {
		<-ctx.Done()
		// Close the connection when context is done
		wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		wsConn.Close()
	}()

	// First message
	sentFirstMessage := false

	// Helper function to extract namespace and pod name from instance label
	extractNamespaceAndPod := func(instance string) (namespace, podName string) {
		parts := strings.Split(instance, "/")
		if len(parts) >= 2 {
			namespace = parts[0]

			// Format is "namespace/podname:service"
			serviceParts := strings.Split(parts[1], ":")
			if len(serviceParts) >= 1 {
				podName = serviceParts[0]
				return namespace, podName
			}
		}
		return "", instance // Return empty namespace and original if parsing fails
	}

	// Main loop for receiving WebSocket messages
	for {
		select {
		case <-ctx.Done():
			return nil

		default:
			// Read from WebSocket
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return nil
				}
				return fmt.Errorf("websocket read error: %v", err)
			}

			// Parse the message
			var streamResp LokiStreamResponse
			if err := json.Unmarshal(message, &streamResp); err != nil {
				log.Error("Failed to unmarshal Loki stream response", "error", err)
				continue
			}

			// Process all logs from this response at once
			var allEvents []LogEvent

			for _, stream := range streamResp.Streams {
				// Extract the instance label
				instance, ok := stream.Stream["instance"]
				if !ok {
					log.Warn("Stream missing instance label", "labels", stream.Stream)
					continue
				}

				namespace, podName := extractNamespaceAndPod(instance)
				if namespace == "" || podName == "" {
					log.Warnf("Failed to parse namespace and pod name from instance: %s", instance)
					continue
				}

				// Get metadata for this pod
				podsInNamespace, ok := podMetadataMap[namespace]
				if !ok {
					log.Warnf("No metadata found for namespace %s", namespace)
					continue
				}

				metadata, ok := podsInNamespace[podName]
				if !ok {
					log.Warnf("No metadata found for pod %s in namespace %s", podName, namespace)
					continue
				}

				for _, entry := range stream.Values {
					// Entry format is [timestamp, log message]
					if len(entry) != 2 {
						log.Warnf("Unprocessable log entry format from loki %v", entry)
						continue
					}

					// Parse timestamp
					var timestamp time.Time
					if ts, err := strconv.ParseInt(entry[0], 10, 64); err == nil {
						// Loki timestamps are in nanoseconds
						timestamp = time.Unix(0, ts)
					}

					// Get the message
					message := entry[1]

					// Apply pattern filter if specified
					if opts.SearchPattern != "" && !strings.Contains(message, opts.SearchPattern) {
						continue
					}

					// Create log event and add it to the collection
					logEvent := LogEvent{
						PodName:   podName,
						Timestamp: timestamp,
						Message:   message,
						Metadata:  metadata,
					}

					allEvents = append(allEvents, logEvent)
				}
			}

			// Make a dummy message if no events
			if len(allEvents) == 0 && !sentFirstMessage {
				allEvents = []LogEvent{}
			}

			// Send events from this batch to the channel
			if len(allEvents) > 0 || !sentFirstMessage {
				sentFirstMessage = true
				select {
				case eventChan <- LogEvents{Logs: allEvents}:
					// Successfully sent
				case <-ctx.Done():
					// Context canceled
					return nil
				}
			}
		}
	}
}
