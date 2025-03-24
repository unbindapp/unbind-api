package variables_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Delete variables
type DeleteVariablesInput struct {
	server.BaseAuthInput
	Body struct {
		BaseVariablesJSONInput
		Variables       []models.VariableDeleteInput `json:"variables" validate:"required"`
		IsBuildVariable bool                         `json:"is_build_variable" required:"false"`
	}
}

func (self *HandlerGroup) DeleteVariables(ctx context.Context, input *DeleteVariablesInput) (*VariablesResponse, error) {
	// Validate input
	if err := ValidateVariablesDependencies(input.Body.Type, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.ServiceID); err != nil {
		return nil, huma.Error400BadRequest("invalid input", err)
	}

	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	var variables []*models.VariableResponse
	var err error
	// Determine which service to use
	switch input.Body.Type {
	case models.TeamVariable:
		variables, err = self.srv.TeamService.DeleteVariablesByKey(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.Variables)
	case models.ProjectVariable:
		variables, err = self.srv.ProjectService.DeleteVariablesByKey(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.Variables)
	case models.EnvironmentVariable:
		variables, err = self.srv.EnvironmentService.DeleteVariablesByKey(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.Variables)
	case models.ServiceVariable:
		variables, err = self.srv.ServiceService.DeleteVariablesByKey(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.ServiceID, input.Body.Variables, input.Body.IsBuildVariable)
	}

	if err != nil {
		return nil, handleVariablesErr(err)
	}

	resp := &VariablesResponse{}
	resp.Body.Data = variables
	return resp, nil
}
