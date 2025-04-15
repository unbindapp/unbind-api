package models

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type InstanceType string

const (
	InstanceTypeTeam        InstanceType = "team"
	InstanceTypeProject     InstanceType = "project"
	InstanceTypeEnvironment InstanceType = "instance"
	InstanceTypeService     InstanceType = "service"
)

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u InstanceType) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["InstanceType"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "InstanceType")
		schemaRef.Title = "InstanceType"
		schemaRef.Enum = append(schemaRef.Enum, string(InstanceTypeTeam))
		schemaRef.Enum = append(schemaRef.Enum, string(InstanceTypeProject))
		schemaRef.Enum = append(schemaRef.Enum, string(InstanceTypeEnvironment))
		schemaRef.Enum = append(schemaRef.Enum, string(InstanceTypeService))
		r.Map()["InstanceType"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/InstanceType"}
}

// InstanceStatusInput defines the query parameters for getting instance statuses
type InstanceStatusInput struct {
	Type          InstanceType `query:"type" required:"true"`
	TeamID        uuid.UUID    `query:"team_id" required:"true"`
	ProjectID     uuid.UUID    `query:"project_id" required:"false"`
	EnvironmentID uuid.UUID    `query:"environment_id" required:"false"`
	ServiceID     uuid.UUID    `query:"service_id" required:"false"`
}
