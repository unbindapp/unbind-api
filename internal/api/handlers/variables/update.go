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
		models.BaseVariablesJSONInput
		Behavior  models.VariableUpdateBehavior `json:"behavior" default:"upsert" required:"true" doc:"The behavior of the update - upsert or overwrite"`
		Variables []*struct {
			Name  string `json:"name" required:"true"`
			Value string `json:"value" required:"true"`
		} `json:"variables" required:"true"`
		VariableReferences *models.MutateVariableReferenceInput `json:"variable_references" required:"false"`
	}
}

func (self *HandlerGroup) UpdateVariables(ctx context.Context, input *UpsertVariablesInput) (*VariablesResponse, error) {
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

	// Update
	variableMap, err := self.srv.VariablesService.UpdateVariables(
		ctx,
		user.ID,
		bearerToken,
		input.Body.VariableReferences,
		input.Body.BaseVariablesJSONInput,
		input.Body.Behavior,
		variablesUpdateMap,
	)
	if err != nil {
		return nil, handleVariablesErr(err)
	}

	resp := &VariablesResponse{}
	resp.Body.Data = variableMap
	return resp, nil
}
