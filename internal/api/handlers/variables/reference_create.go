package variables_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Create new
type CreateVariableReferenceInput struct {
	server.BaseAuthInput
	Body *models.CreateVariableReferenceInput
}

type CreateVariableReferenceResponse struct {
	Body struct {
		Data *models.VariableReferenceResponse `json:"data"`
	}
}

func (self *HandlerGroup) CreateVariableReference(ctx context.Context, input *CreateVariableReferenceInput) (*CreateVariableReferenceResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	created, err := self.srv.VariablesService.CreateVariableReference(ctx, user.ID, input.Body)

	if err != nil {
		return nil, handleVariablesErr(err)
	}

	resp := &CreateVariableReferenceResponse{}
	resp.Body.Data = created
	return resp, nil
}
