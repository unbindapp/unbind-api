package variables_handler

import (
	"errors"
	"net/http"

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

func RegisterHandlers(server *server.Server, grp *huma.Group) {
	handlers := &HandlerGroup{
		srv: server,
	}

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-variables",
			Summary:     "List Variables",
			Description: "List variables for a service, environment, project, or team",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		handlers.ListVariables,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-variables",
			Summary:     "Update Variables",
			Description: "Create or update variables for a service, environment, project, or team by key",
			Path:        "/update",
			Method:      http.MethodPost,
		},
		handlers.UpdateVariables,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "delete-variables",
			Summary:     "Delete Variables",
			Description: "Delete variables for a service, environment, project, or team",
			Path:        "/delete",
			Method:      http.MethodDelete,
		},
		handlers.DeleteVariables,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-available-references",
			Summary:     "List Available Variable References",
			Description: "List all items that can be references as variables in service configurations",
			Path:        "/references/available",
			Method:      http.MethodGet,
		},
		handlers.ListReferenceableeVariables,
	)
}

func handleVariablesErr(err error) error {
	if errors.Is(err, errdefs.ErrUnauthorized) {
		return huma.Error403Forbidden("Unauthorized")
	}
	if errors.Is(err, errdefs.ErrNotFound) {
		return huma.Error404NotFound(err.Error())
	}
	log.Error("Error getting variables", "err", err)
	return huma.Error500InternalServerError("Unable to retrieve variables")
}

// Base inputs
type BaseVariablesInput struct {
	Type          models.VariableType `query:"type" required:"true" doc:"The type of variable"`
	TeamID        uuid.UUID           `query:"team_id" required:"true"`
	ProjectID     uuid.UUID           `query:"project_id" doc:"If present, fetch project variables"`
	EnvironmentID uuid.UUID           `query:"environment_id" doc:"If present, fetch environment variables - requires project_id"`
	ServiceID     uuid.UUID           `query:"service_id" doc:"If present, fetch service variables - requires project_id and environment_id"`
}

type BaseVariablesJSONInput struct {
	Type          models.VariableType `json:"type" required:"true" doc:"The type of variable"`
	TeamID        uuid.UUID           `json:"team_id" required:"true"`
	ProjectID     uuid.UUID           `json:"project_id" required:"false" doc:"If present without environment_id, mutate team variables"`
	EnvironmentID uuid.UUID           `json:"environment_id" required:"false" doc:"If present without service_id, mutate environment variables - requires project_id"`
	ServiceID     uuid.UUID           `json:"service_id" required:"false" doc:"If present, mutate service variables - requires project_id and environment_id"`
}

func ValidateVariablesDependencies(variableType models.VariableType, teamID, projectID, environmentID, serviceID uuid.UUID) error {
	switch variableType {
	case models.TeamVariable:
		if teamID == uuid.Nil {
			return huma.Error400BadRequest("team_id is required")
		}
	case models.ProjectVariable:
		if teamID == uuid.Nil || projectID == uuid.Nil {
			return huma.Error400BadRequest("team_id and project_id are required")
		}
	case models.EnvironmentVariable:
		if teamID == uuid.Nil || projectID == uuid.Nil || environmentID == uuid.Nil {
			return huma.Error400BadRequest("team_id, project_id, and environment_id are required")
		}
	case models.ServiceVariable:
		if teamID == uuid.Nil || projectID == uuid.Nil || environmentID == uuid.Nil || serviceID == uuid.Nil {
			return huma.Error400BadRequest("team_id, project_id, environment_id, and service_id are required")
		}
	}
	return nil
}
