package models

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type PvcScope string

const (
	PvcScopeTeam        PvcScope = "team"
	PvcScopeProject     PvcScope = "project"
	PvcScopeEnvironment PvcScope = "environment"
)

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u PvcScope) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["PvcScope"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "PvcScope")
		schemaRef.Title = "PvcScope"
		schemaRef.Enum = append(schemaRef.Enum,
			[]any{
				string(PvcScopeTeam),
				string(PvcScopeProject),
				string(PvcScopeEnvironment),
			}...)
		r.Map()["PvcScope"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/PvcScope"}
}

// * List
type ListPVCInput struct {
	Type          PvcScope  `query:"type" required:"true"`
	TeamID        uuid.UUID `query:"team_id" required:"true"`
	ProjectID     uuid.UUID `query:"project_id" required:"false"`
	EnvironmentID uuid.UUID `query:"environment_id" required:"false"`
}

// * Get
type GetPVCInput struct {
	Type          PvcScope  `query:"type" required:"true"`
	TeamID        uuid.UUID `query:"team_id" required:"true"`
	ProjectID     uuid.UUID `query:"project_id" required:"false"`
	EnvironmentID uuid.UUID `query:"environment_id" required:"false"`
	ID            string    `query:"id" required:"true"`
}

// * Create
type CreatePVCInput struct {
	Type          PvcScope  `query:"type" required:"true"`
	Name          string    `query:"name" required:"true" minLength:"1" maxLength:"63" doc:"Name of the PVC"`
	TeamID        uuid.UUID `query:"team_id" required:"true"`
	ProjectID     uuid.UUID `query:"project_id" required:"false"`
	EnvironmentID uuid.UUID `query:"environment_id" required:"false"`
	Size          string    `query:"size" required:"true" doc:"Size of the PVC (e.g., '10Gi')"`
}
