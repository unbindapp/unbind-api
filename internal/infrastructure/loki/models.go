package loki

import (
	"time"

	"github.com/google/uuid"
)

// LokiLogOptions represents options for filtering and streaming logs from Loki
type LokiLogOptions struct {
	PodNamespace  string            // Kubernetes namespace to query
	PodName       string            // Pod name to query
	Labels        map[string]string // Additional labels to filter by
	Since         time.Duration     // Get logs from this time ago
	Limit         int               // Number of log lines to get
	Follow        bool              // Whether to stream logs (tail)
	SinceTime     *time.Time        // Get logs from a specific time
	SearchPattern string            // Optional text pattern to search for
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
	Logs        []LogEvent `json:"logs" nullable:"false"`
	IsHeartbeat bool       `json:"is_heartbeat,omitempty"`
}

// Stream Error
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

// LokiStreamResponse represents the format of a Loki log stream response
type LokiStreamResponse struct {
	Streams []struct {
		Stream map[string]string `json:"stream"`
		Values [][2]string       `json:"values"` // [timestamp, message]
	} `json:"streams"`
}
