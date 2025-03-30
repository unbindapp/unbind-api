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

// StreamLokiLogs streams logs from Loki tail API using WebSocket.
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

	// First message
	sentFirstMessage := false

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
				streamEvents := make([]LogEvent, len(stream.Values))
				for i, entry := range stream.Values {
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
					streamEvents[i] = LogEvent{
						PodName:   opts.PodName,
						Timestamp: timestamp,
						Message:   message,
						Metadata:  meta,
					}
				}
				allEvents = append(allEvents, streamEvents...)
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
