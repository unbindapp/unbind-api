package models

import (
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type DeploymentInputRequirements interface {
	GetTeamID() uuid.UUID
	GetProjectID() uuid.UUID
	GetServiceID() uuid.UUID
	GetEnvironmentID() uuid.UUID
}

type GetDeploymentsInput struct {
	PaginationParams
	Status        schema.DeploymentStatus `query:"status" required:"false" doc:"Filter by status"`
	ID            uuid.UUID               `query:"id" required:"true" doc:"The ID of the build"`
	TeamID        uuid.UUID               `query:"team_id" required:"true" doc:"The ID of the team"`
	ProjectID     uuid.UUID               `query:"project_id" required:"true" doc:"The ID of the project"`
	EnvironmentID uuid.UUID               `query:"environment_id" required:"true" doc:"The ID of the environment"`
	ServiceID     uuid.UUID               `query:"service_id" required:"true" doc:"The ID of the service"`
}

func (self *GetDeploymentsInput) GetTeamID() uuid.UUID {
	return self.TeamID
}

func (self *GetDeploymentsInput) GetProjectID() uuid.UUID {
	return self.ProjectID
}

func (self *GetDeploymentsInput) GetServiceID() uuid.UUID {
	return self.ServiceID
}

func (self *GetDeploymentsInput) GetEnvironmentID() uuid.UUID {
	return self.EnvironmentID
}

// Triggering build

type CreateDeploymentInput struct {
	TeamID        uuid.UUID `validate:"required,uuid4" required:"true" json:"team_id"`
	ProjectID     uuid.UUID `validate:"required,uuid4" required:"true" json:"project_id"`
	ServiceID     uuid.UUID `validate:"required,uuid4" required:"true" json:"service_id"`
	EnvironmentID uuid.UUID `validate:"required,uuid4" required:"true" json:"environment_id"`
}

func (self *CreateDeploymentInput) GetTeamID() uuid.UUID {
	return self.TeamID
}

func (self *CreateDeploymentInput) GetProjectID() uuid.UUID {
	return self.ProjectID
}

func (self *CreateDeploymentInput) GetServiceID() uuid.UUID {
	return self.ServiceID
}

func (self *CreateDeploymentInput) GetEnvironmentID() uuid.UUID {
	return self.EnvironmentID
}
