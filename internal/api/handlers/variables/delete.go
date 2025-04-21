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

// Delete variables
type DeleteVariablesInput struct {
	server.BaseAuthInput
	Body struct {
		models.BaseVariablesJSONInput
		Variables    []models.VariableDeleteInput `json:"variables" required:"false" nullable:"false"`
		ReferenceIDs []uuid.UUID                  `json:"reference_ids" required:"false" nullable:"false"`
	}
}

func (self *HandlerGroup) DeleteVariables(ctx context.Context, input *DeleteVariablesInput) (*VariablesResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	variableMap, err := self.srv.VariablesService.DeleteVariablesByKey(
		ctx,
		user.ID,
		bearerToken,
		input.Body.BaseVariablesJSONInput,
		input.Body.Variables,
		input.Body.ReferenceIDs,
	)
	if err != nil {
		return nil, handleVariablesErr(err)
	}

	resp := &VariablesResponse{}
	resp.Body.Data = variableMap
	return resp, nil
}
