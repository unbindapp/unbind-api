package models

import (
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type BuildJobInputRequirements interface {
	GetTeamID() uuid.UUID
	GetProjectID() uuid.UUID
	GetServiceID() uuid.UUID
	GetEnvironmentID() uuid.UUID
}

type GetBuildJobsInput struct {
	PaginationParams
	Status        schema.BuildJobStatus `query:"status" required:"false" doc:"Filter by status"`
	ID            uuid.UUID             `query:"id" required:"true" doc:"The ID of the build"`
	TeamID        uuid.UUID             `query:"team_id" required:"true" doc:"The ID of the team"`
	ProjectID     uuid.UUID             `query:"project_id" required:"true" doc:"The ID of the project"`
	EnvironmentID uuid.UUID             `query:"environment_id" required:"true" doc:"The ID of the environment"`
	ServiceID     uuid.UUID             `query:"service_id" required:"true" doc:"The ID of the service"`
}

func (self *GetBuildJobsInput) GetTeamID() uuid.UUID {
	return self.TeamID
}

func (self *GetBuildJobsInput) GetProjectID() uuid.UUID {
	return self.ProjectID
}

func (self *GetBuildJobsInput) GetServiceID() uuid.UUID {
	return self.ServiceID
}

func (self *GetBuildJobsInput) GetEnvironmentID() uuid.UUID {
	return self.EnvironmentID
}

// Triggering build

type CreateBuildJobInput struct {
	TeamID        uuid.UUID `validate:"required,uuid4" required:"true" json:"team_id"`
	ProjectID     uuid.UUID `validate:"required,uuid4" required:"true" json:"project_id"`
	ServiceID     uuid.UUID `validate:"required,uuid4" required:"true" json:"service_id"`
	EnvironmentID uuid.UUID `validate:"required,uuid4" required:"true" json:"environment_id"`
}

func (self *CreateBuildJobInput) GetTeamID() uuid.UUID {
	return self.TeamID
}

func (self *CreateBuildJobInput) GetProjectID() uuid.UUID {
	return self.ProjectID
}

func (self *CreateBuildJobInput) GetServiceID() uuid.UUID {
	return self.ServiceID
}

func (self *CreateBuildJobInput) GetEnvironmentID() uuid.UUID {
	return self.EnvironmentID
}
