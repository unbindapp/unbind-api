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
)

// StreamLokiLogs streams logs from Loki to the provided channel using WebSocket
func (self *LokiLogQuerier) StreamLokiLogs(
	ctx context.Context,
	opts LokiLogOptions,
	meta LogMetadata,
	eventChan chan<- LogEvents,
) error {
	// Alloy sets the instance like instance="namespace/podname:service"
	queryStr := fmt.Sprintf("{instance=\"%s/%s:service\"}", opts.PodNamespace, opts.PodName)

	if opts.SearchPattern != "" {
		queryStr = fmt.Sprintf("%s |= \"%s\"", queryStr, opts.SearchPattern)
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

	// Constants for batching
	const (
		batchSize    = 100
		maxBatchWait = 200 * time.Millisecond
	)

	// Batching logic
	batch := make([]LogEvent, 0, batchSize)
	timer := time.NewTimer(maxBatchWait)
	defer timer.Stop()
	sentFirstMessage := false

	// Function to send the current batch
	sendBatch := func() {
		if len(batch) == 0 && sentFirstMessage {
			return
		}

		sentFirstMessage = true
		events := make([]LogEvent, len(batch))
		copy(events, batch)

		select {
		case eventChan <- LogEvents{Logs: events}:
			// Successfully sent
		case <-ctx.Done():
			// Context canceled
		}

		batch = batch[:0] // Clear the batch
	}

	// Reset timer
	if !timer.Stop() {
		<-timer.C
	}
	timer.Reset(maxBatchWait)

	// Main loop for receiving WebSocket messages
	for {
		select {
		case <-ctx.Done():
			sendBatch()
			return nil

		case <-timer.C:
			sendBatch()
			timer.Reset(maxBatchWait)

		default:
			// Read from WebSocket
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					sendBatch()
					return nil
				}
				sendBatch()
				return fmt.Errorf("websocket read error: %v", err)
			}

			// Parse the message
			var streamResp LokiStreamResponse
			if err := json.Unmarshal(message, &streamResp); err != nil {
				// Skip invalid JSON
				continue
			}

			// Process the logs
			for _, stream := range streamResp.Streams {
				for _, entry := range stream.Values {
					// Entry format is [timestamp, log message]
					if len(entry) != 2 {
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

					// Create and add log event
					batch = append(batch, LogEvent{
						PodName:   stream.Stream["instance"], // Use the instance label
						Timestamp: timestamp,
						Message:   message,
						Metadata:  meta,
					})

					// Check batch size
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
	}
}
