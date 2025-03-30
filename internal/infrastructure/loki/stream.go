package loki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
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
			escapedNamespace := regexp.QuoteMeta(namespace)
			escapedPodName := regexp.QuoteMeta(podName)
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
		// Case insensitive search
		queryStr = fmt.Sprintf("(%s) |~ \"(?i)%s\"", queryStr, regexp.QuoteMeta(opts.SearchPattern))
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

	reqURL.RawQuery = q.Encode()

	log.Infof("Streaming logs with query: %s, URL: %s", queryStr, reqURL.String())

	// Create websocket connection with timeout
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 15 * time.Second // Set a reasonable timeout

	wsConn, resp, err := dialer.DialContext(ctx, reqURL.String(), nil)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("failed to connect to Loki WebSocket: %v, status: %d", err, resp.StatusCode)
		}
		return fmt.Errorf("failed to connect to Loki WebSocket: %v", err)
	}
	defer wsConn.Close()

	// Set the ping handler to respond with pongs to keep the connection alive
	wsConn.SetPingHandler(func(data string) error {
		err := wsConn.WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(5*time.Second))
		if err != nil {
			log.Warnf("Failed to send pong: %v", err)
		}
		return nil
	})

	// Set the pong handler to reset the read deadline when a pong is received
	wsConn.SetPongHandler(func(string) error {
		wsConn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Initial read deadline - this will be extended by pongs and successful reads
	wsConn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Setup context cancellation
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		log.Info("Context done, closing WebSocket connection")
		// Close the connection gracefully when context is done
		wsConn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			time.Now().Add(5*time.Second))
		wsConn.Close()
		close(done)
	}()

	// Initialize a ping ticker to send pings periodically
	pingTicker := time.NewTicker(15 * time.Second)
	defer pingTicker.Stop()

	// Initialize a heartbeat ticker to send empty messages periodically to client
	heartbeatTicker := time.NewTicker(10 * time.Second)
	defer heartbeatTicker.Stop()

	// Confirm connection is successful by sending an initial empty message
	select {
	case eventChan <- LogEvents{Logs: []LogEvent{}}:
	case <-ctx.Done():
		return nil
	}

	// Main loop for handling the WebSocket connection
	go func() {
		for {
			select {
			case <-done:
				return
			case <-pingTicker.C:
				// Send ping to server to keep connection alive
				err := wsConn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second))
				if err != nil {
					log.Warnf("Failed to send ping: %v", err)
				}
			case <-heartbeatTicker.C:
				// Send heartbeat message to keep the client side alive
				select {
				case eventChan <- LogEvents{Logs: []LogEvent{}, IsHeartbeat: true}:
				case <-done:
					return
				default:
					log.Warn("Failed to send heartbeat to client (channel blocked)")
				}
			}
		}
	}()

	// Main read loop
	for {
		// Check if context is done before attempting to read
		select {
		case <-done:
			return nil
		default:
			// Continue with read
		}

		// Read from WebSocket
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			// Check for normal closure
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Info("WebSocket closed normally")
				return nil
			}

			// Check for timeout - we'll try to keep the connection alive
			if strings.Contains(err.Error(), "i/o timeout") || strings.Contains(err.Error(), "deadline exceeded") {
				log.Warn("WebSocket read timeout, attempting to keep connection alive")

				// Send a ping to check if connection is still alive
				err := wsConn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second))
				if err != nil {
					log.Errorf("Failed to send ping after timeout, connection appears dead: %v", err)
					return fmt.Errorf("websocket connection dead: %v", err)
				}

				// Reset the read deadline and continue
				wsConn.SetReadDeadline(time.Now().Add(60 * time.Second))
				continue
			}

			// For any other error, return and let the caller handle reconnection if needed
			log.Errorf("WebSocket read error: %v", err)
			return fmt.Errorf("websocket read error: %v", err)
		}

		// Reset read deadline after successful read
		wsConn.SetReadDeadline(time.Now().Add(60 * time.Second))

		// Parse the message
		var streamResp LokiStreamResponse
		if err := json.Unmarshal(message, &streamResp); err != nil {
			log.Warnf("Failed to unmarshal Loki stream response: %v", err)
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

			namespace, podName, ok := extractNamespaceAndPod(instance)
			if !ok {
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
				} else {
					log.Warnf("Failed to parse timestamp: %v", err)
					// Use current time as fallback
					timestamp = time.Now()
				}

				// Get the message
				message := entry[1]

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

		// Send events from this batch to the channel if there are any
		if len(allEvents) > 0 {
			select {
			case eventChan <- LogEvents{Logs: allEvents}:
			case <-done:
				// Context canceled
				return nil
			default:
				// Channel is blocked, log a warning but continue
				log.Warnf("Event channel blocked, couldn't send %d log events", len(allEvents))
			}
		}
	}
}

func extractNamespaceAndPod(instance string) (namespace, podName string, ok bool) {
	parts := strings.Split(instance, "/")
	if len(parts) < 2 {
		return "", "", false
	}

	namespace = parts[0]

	// Format is "namespace/podname:service"
	serviceParts := strings.Split(parts[1], ":")
	if len(serviceParts) < 1 {
		return "", "", false
	}

	podName = serviceParts[0]
	if namespace == "" || podName == "" {
		return "", "", false
	}

	return namespace, podName, true
}
