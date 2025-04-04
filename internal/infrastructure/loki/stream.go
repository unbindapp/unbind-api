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
	opts LokiLogStreamOptions,
	eventChan chan<- LogEvents,
) error {
	queryStr := fmt.Sprintf("{%s=\"%s\"}", opts.Label, opts.LabelValue)

	// Add extra filters
	if opts.RawFilter != "" {
		queryStr = fmt.Sprintf("%s %s", queryStr, opts.RawFilter)
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
		if opts.Limit > 1000 {
			opts.Limit = 1000
		}
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
				case eventChan <- LogEvents{MessageType: LogEventsMessageTypeHeartbeat}:
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
			// Get metadata for this stream
			instance, ok := stream.Stream["instance"]
			if !ok {
				log.Warnf("Stream missing instance label: %s", stream.Stream)
			}
			environmentID, _ := stream.Stream[string(LokiLabelEnvironment)]
			teamID, _ := stream.Stream[string(LokiLabelTeam)]
			projectID, _ := stream.Stream[string(LokiLabelProject)]
			serviceID, _ := stream.Stream[string(LokiLabelService)]
			deploymentID, _ := stream.Stream[string(LokiLabelDeployment)]

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
					PodName:   instance,
					Timestamp: timestamp,
					Message:   message,
					Metadata: LogMetadata{
						TeamID:        teamID,
						ProjectID:     projectID,
						EnvironmentID: environmentID,
						ServiceID:     serviceID,
						DeploymentID:  deploymentID,
					},
				}

				allEvents = append(allEvents, logEvent)
			}
		}

		// Send events from this batch to the channel if there are any
		if len(allEvents) > 0 {
			select {
			case eventChan <- LogEvents{MessageType: LogEventsMessageTypeLog, Logs: allEvents}:
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
