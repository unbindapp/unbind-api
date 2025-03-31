package loki

import (
	"time"
)

type LokiLabelName string

const (
	LokiLabelTeam        LokiLabelName = "unbind_team"
	LokiLabelProject     LokiLabelName = "unbind_project"
	LokiLabelEnvironment LokiLabelName = "unbind_environment"
	LokiLabelService     LokiLabelName = "unbind_service"
)

// LokiLogOptions represents options for filtering and streaming logs from Loki
type LokiLogOptions struct {
	Label         LokiLabelName // Label to filter logs by
	LabelValue    string        // Value of the label to filter logs by
	Since         time.Duration // Get logs from this time ago
	Limit         int           // Number of log lines to get
	SinceTime     *time.Time    // Get logs from a specific time
	SearchPattern string        // Optional text pattern to search for
}

type LogMetadata struct {
	// Metadata to stick on
	ServiceID     string `json:"service_id"`
	TeamID        string `json:"team_id"`
	ProjectID     string `json:"project_id"`
	EnvironmentID string `json:"environment_id"`
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
