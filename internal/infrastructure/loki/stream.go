package loki

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// We can reuse your existing LogMetadata and LogEvent types

// StreamLokiLogs streams logs from Loki to the provided channel
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

	// Prepare HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Execute the request
	resp, err := self.client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error from Loki: %s, status: %d", string(body), resp.StatusCode)
	}

	// Constants for batching, matching your existing implementation
	const (
		batchSize    = 100
		maxBatchWait = 200 * time.Millisecond
	)

	// Channels for processing
	linesChan := make(chan LokiStreamResponse)
	readErrChan := make(chan error, 1)

	// Start reading from the response
	go func() {
		defer close(linesChan)

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				readErrChan <- err
				return
			}

			if len(line) == 0 {
				continue
			}

			// Decode the Loki stream response (different format than K8s logs)
			var streamResp LokiStreamResponse
			if err := json.Unmarshal(line, &streamResp); err != nil {
				continue // Skip invalid JSON
			}

			linesChan <- streamResp
		}
	}()

	// Same batching logic as your existing code
	batch := make([]LogEvent, 0, batchSize)
	timer := time.NewTimer(maxBatchWait)
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

	// Main loop
	for {
		select {
		case <-ctx.Done():
			sendBatch()
			return nil

		case err := <-readErrChan:
			if err == io.EOF {
				sendBatch()
				return nil
			}
			return fmt.Errorf("error reading from Loki stream: %v", err)

		case <-timer.C:
			sendBatch()
			timer.Reset(maxBatchWait)

		case streamResp, ok := <-linesChan:
			if !ok {
				sendBatch()
				return nil
			}

			// Process the Loki stream message
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
