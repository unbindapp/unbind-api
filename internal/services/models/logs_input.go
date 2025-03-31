package models

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type LogType string

const (
	LogTypeTeam        LogType = "team"
	LogTypeProject     LogType = "project"
	LogTypeEnvironment LogType = "environment"
	LogTypeService     LogType = "service"
)

var LogTypeValues = []LogType{
	LogTypeTeam,
	LogTypeProject,
	LogTypeEnvironment,
	LogTypeService,
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u LogType) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["LogType"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "LogType")
		schemaRef.Title = "LogType"
		for _, v := range LogTypeValues {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["LogType"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/LogType"}
}

// LogQueryInput defines the query parameters for log streaming
type LogQueryInput struct {
	Type          LogType   `query:"type" required:"true"`
	TeamID        uuid.UUID `query:"team_id" required:"true"`
	ProjectID     uuid.UUID `query:"project_id" required:"false"`
	EnvironmentID uuid.UUID `query:"environment_id" required:"false"`
	ServiceID     uuid.UUID `query:"service_id" required:"false"`
	Since         string    `query:"since" default:"10m" doc:"Duration to look back (e.g., '1h', '30m')"`
	Tail          int64     `query:"tail" default:"100" doc:"Number of lines to get from the end"`
	Previous      bool      `query:"previous" doc:"Get logs from previous instance"`
	Timestamps    bool      `query:"timestamps" default:"true" doc:"Include timestamps in logs"`
	Filters       string    `query:"filters" doc:"Optional logql filter string"`
}
