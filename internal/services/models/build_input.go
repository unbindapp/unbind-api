package models

import (
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type GetBuildJobsInput struct {
	PaginationParams
	Status        schema.BuildJobStatus `query:"status" required:"false" doc:"Filter by status"`
	ID            uuid.UUID             `query:"id" required:"true" doc:"The ID of the build"`
	TeamID        uuid.UUID             `query:"team_id" required:"true" doc:"The ID of the team"`
	ProjectID     uuid.UUID             `query:"project_id" required:"true" doc:"The ID of the project"`
	EnvironmentID uuid.UUID             `query:"environment_id" required:"true" doc:"The ID of the environment"`
	ServiceID     uuid.UUID             `query:"service_id" required:"true" doc:"The ID of the service"`
}
