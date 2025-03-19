package secrets_handler

import (
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

type HandlerGroup struct {
	srv *server.Server
}

func NewHandlerGroup(server *server.Server) *HandlerGroup {
	return &HandlerGroup{
		srv: server,
	}
}

func handleSecretErr(err error) error {
	if errors.Is(err, errdefs.ErrUnauthorized) {
		return huma.Error403Forbidden("Unauthorized")
	}
	if errors.Is(err, errdefs.ErrNotFound) {
		return huma.Error404NotFound(err.Error())
	}
	log.Error("Error getting secrets", "err", err)
	return huma.Error500InternalServerError("Unable to retrieve secrets")
}

// Base inputs
type BaseSecretsInput struct {
	TeamID        uuid.UUID `query:"team_id" required:"true"`
	ProjectID     uuid.UUID `query:"project_id" doc:"If present, fetch project secrets"`
	EnvironmentID uuid.UUID `query:"environment_id" doc:"If present, fetch environment secrets - requires project_id"`
	ServiceID     uuid.UUID `query:"service_id" doc:"If present, fetch service secrets - requires project_id and environment_id"`
}

type BaseSecretsJSONInput struct {
	TeamID        uuid.UUID  `json:"team_id" validate:"required"`
	ProjectID     *uuid.UUID `json:"project_id" doc:"If present without environment_id, mutate team secrets"`
	EnvironmentID *uuid.UUID `json:"environment_id" doc:"If present without service_id, mutate environment secrets - requires project_id"`
	ServiceID     *uuid.UUID `json:"service_id" doc:"If present, mutate service secrets - requires project_id and environment_id"`
}

func ValidateSecretsDependencies[T interface {
	GetTeamID() uuid.UUID
	GetProjectID() *uuid.UUID
	GetEnvironmentID() *uuid.UUID
	GetServiceID() *uuid.UUID
}](input T) error {
	// TeamID is already validated with required tag

	// If EnvironmentID is present, ProjectID must also be present
	if input.GetEnvironmentID() != nil && input.GetEnvironmentID().String() != "" {
		if input.GetProjectID() == nil || input.GetProjectID().String() == "" {
			return fmt.Errorf("if environment_id is provided, project_id must also be provided")
		}
	}

	// If ServiceID is present, EnvironmentID, ProjectID, and TeamID must be present
	if input.GetServiceID() != nil && input.GetServiceID().String() != "" {
		if input.GetEnvironmentID() == nil || input.GetEnvironmentID().String() == "" {
			return fmt.Errorf("if service_id is provided, environment_id must also be provided")
		}
		if input.GetProjectID() == nil || input.GetProjectID().String() == "" {
			return fmt.Errorf("if service_id is provided, project_id must also be provided")
		}
	}

	return nil
}

// Implement getter methods for BaseSecretsInput
// Huma doesn't like pointers as query parameters so we check for uuid.Nil
func (b BaseSecretsInput) GetTeamID() uuid.UUID { return b.TeamID }
func (b BaseSecretsInput) GetProjectID() *uuid.UUID {
	if b.ProjectID == uuid.Nil {
		return nil
	}
	return utils.ToPtr(b.ProjectID)
}
func (b BaseSecretsInput) GetEnvironmentID() *uuid.UUID {
	if b.EnvironmentID == uuid.Nil {
		return nil
	}
	return utils.ToPtr(b.EnvironmentID)
}
func (b BaseSecretsInput) GetServiceID() *uuid.UUID {
	if b.ServiceID == uuid.Nil {
		return nil
	}
	return utils.ToPtr(b.ServiceID)
}

// Implement getter methods for BaseSecretsJSONInput
func (b BaseSecretsJSONInput) GetTeamID() uuid.UUID         { return b.TeamID }
func (b BaseSecretsJSONInput) GetProjectID() *uuid.UUID     { return b.ProjectID }
func (b BaseSecretsJSONInput) GetEnvironmentID() *uuid.UUID { return b.EnvironmentID }
func (b BaseSecretsJSONInput) GetServiceID() *uuid.UUID     { return b.ServiceID }
