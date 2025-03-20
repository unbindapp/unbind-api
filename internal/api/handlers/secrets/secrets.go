package secrets_handler

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
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
	Type          models.SecretType `json:"type" required:"true" doc:"The type of secret to fetch"`
	TeamID        uuid.UUID         `query:"team_id" required:"true"`
	ProjectID     uuid.UUID         `query:"project_id" doc:"If present, fetch project secrets"`
	EnvironmentID uuid.UUID         `query:"environment_id" doc:"If present, fetch environment secrets - requires project_id"`
	ServiceID     uuid.UUID         `query:"service_id" doc:"If present, fetch service secrets - requires project_id and environment_id"`
}

type BaseSecretsJSONInput struct {
	Type          models.SecretType `json:"type" required:"true" doc:"The type of secret to fetch"`
	TeamID        uuid.UUID         `json:"team_id" required:"true"`
	ProjectID     uuid.UUID         `json:"project_id" required:"false" doc:"If present without environment_id, mutate team secrets"`
	EnvironmentID uuid.UUID         `json:"environment_id" required:"false" doc:"If present without service_id, mutate environment secrets - requires project_id"`
	ServiceID     uuid.UUID         `json:"service_id" required:"false" doc:"If present, mutate service secrets - requires project_id and environment_id"`
}

func ValidateSecretsDependencies(secretType models.SecretType, teamID, projectID, environmentID, serviceID uuid.UUID) error {
	switch secretType {
	case models.TeamSecret:
		if teamID == uuid.Nil {
			return huma.Error400BadRequest("team_id is required")
		}
	case models.ProjectSecret:
		if teamID == uuid.Nil || projectID == uuid.Nil {
			return huma.Error400BadRequest("team_id and project_id are required")
		}
	case models.EnvironmentSecret:
		if teamID == uuid.Nil || projectID == uuid.Nil || environmentID == uuid.Nil {
			return huma.Error400BadRequest("team_id, project_id, and environment_id are required")
		}
	case models.ServiceSecret:
		if teamID == uuid.Nil || projectID == uuid.Nil || environmentID == uuid.Nil || serviceID == uuid.Nil {
			return huma.Error400BadRequest("team_id, project_id, environment_id, and service_id are required")
		}
	}
	return nil
}
