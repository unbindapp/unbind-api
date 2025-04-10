package variables_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Create new
type UpsertVariablesInput struct {
	server.BaseAuthInput
	Body struct {
		BaseVariablesJSONInput
		Behavior  models.VariableUpdateBehavior `json:"behavior" default:"upsert" required:"true" doc:"The behavior of the update - upsert or overwrite"`
		Variables []*struct {
			Name  string `json:"name" required:"true"`
			Value string `json:"value" required:"true"`
		} `json:"variables" required:"true"`
	}
}

func (self *HandlerGroup) UpdateVariables(ctx context.Context, input *UpsertVariablesInput) (*VariablesResponse, error) {
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

	variablesUpdateMap := make(map[string][]byte)
	for _, variable := range input.Body.Variables {
		variablesUpdateMap[variable.Name] = []byte(variable.Value)
	}

	// Determine which service to use
	var variables []*models.VariableResponse
	var err error
	switch input.Body.Type {
	case models.TeamVariable:
		variables, err = self.srv.TeamService.UpdateVariables(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.Behavior, variablesUpdateMap)
	case models.ProjectVariable:
		variables, err = self.srv.ProjectService.UpdateVariables(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.Behavior, variablesUpdateMap)
	case models.EnvironmentVariable:
		variables, err = self.srv.EnvironmentService.UpdateVariables(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.Behavior, variablesUpdateMap)
	case models.ServiceVariable:
		variables, err = self.srv.ServiceService.UpdateVariables(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.ServiceID, input.Body.Behavior, variablesUpdateMap)
	}

	if err != nil {
		return nil, handleVariablesErr(err)
	}

	resp := &VariablesResponse{}
	resp.Body.Data = variables
	return resp, nil
}
