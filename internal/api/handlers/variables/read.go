package variables_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// List all
type ListVariablesInput struct {
	server.BaseAuthInput
	BaseVariablesInput
}

type VariablesResponse struct {
	Body struct {
		Data []*models.VariableResponse `json:"data" nullable:"false"`
	}
}

func (self *HandlerGroup) ListVariables(ctx context.Context, input *ListVariablesInput) (*VariablesResponse, error) {
	// Validate input
	if err := ValidateVariablesDependencies(input.Type, input.TeamID, input.ProjectID, input.EnvironmentID, input.ServiceID); err != nil {
		return nil, huma.Error400BadRequest("invalid input", err)
	}

	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	// Get team variables
	var err error
	var teamVariables []*models.VariableResponse
	var projectVariables []*models.VariableResponse
	var environmentVariables []*models.VariableResponse
	var serviceVariables []*models.VariableResponse

	// Get team variables always
	teamVariables, err = self.srv.TeamService.GetVariables(ctx, user.ID, bearerToken, input.TeamID)
	if err != nil {
		return nil, handleVariablesErr(err)
	}

	if input.Type == models.ProjectVariable ||
		input.Type == models.EnvironmentVariable ||
		input.Type == models.ServiceVariable {
		projectVariables, err = self.srv.ProjectService.GetVariables(ctx, user.ID, bearerToken, input.TeamID, input.ProjectID)
		if err != nil {
			return nil, handleVariablesErr(err)
		}
	}

	if input.Type == models.EnvironmentVariable ||
		input.Type == models.ServiceVariable {
		environmentVariables, err = self.srv.EnvironmentService.GetVariables(ctx, user.ID, bearerToken, input.TeamID, input.ProjectID, input.EnvironmentID)
		if err != nil {
			return nil, handleVariablesErr(err)
		}
	}

	if input.Type == models.ServiceVariable {
		serviceVariables, err = self.srv.ServiceService.GetVariables(ctx, user.ID, bearerToken, input.TeamID, input.ProjectID, input.EnvironmentID, input.ServiceID)
		if err != nil {
			return nil, handleVariablesErr(err)
		}
	}

	// Combine
	variables := append(teamVariables, projectVariables...)
	variables = append(variables, environmentVariables...)
	variables = append(variables, serviceVariables...)

	resp := &VariablesResponse{}
	resp.Body.Data = variables
	return resp, nil
}

// List all
type ListReferenceableVariablesInput struct {
	server.BaseAuthInput
	TeamID uuid.UUID `query:"team_id" required:"true"`
}

type ReferenceableVariablsResponse struct {
	Body struct {
		Data *models.AvailableVariableReferenceResponse `json:"data" nullable:"false"`
	}
}

func (self *HandlerGroup) ListReferenceableeVariables(ctx context.Context, input *ListReferenceableVariablesInput) (*ReferenceableVariablsResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	// Get team variables
	references, err := self.srv.VariablesService.GetAvailableVariableReferences(ctx, user.ID, bearerToken, input.TeamID)
	if err != nil {
		return nil, handleVariablesErr(err)
	}

	resp := &ReferenceableVariablsResponse{}
	resp.Body.Data = references
	return resp, nil
}
