package environments_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	environment_service "github.com/unbindapp/unbind-api/internal/services/environment"
	"github.com/unbindapp/unbind-api/internal/models"
)

type CreateEnvironmentInput struct {
	server.BaseAuthInput
	Body *environment_service.CreateEnvironmentInput
}

type CreateEnvironmentResponse struct {
	Body struct {
		Data *models.EnvironmentResponse `json:"data"`
	}
}

func (self *HandlerGroup) CreateEnvironment(ctx context.Context, input *CreateEnvironmentInput) (*CreateEnvironmentResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	createdEnvironment, err := self.srv.EnvironmentService.CreateEnvironment(ctx, user.ID, input.Body, bearerToken)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &CreateEnvironmentResponse{}
	resp.Body.Data = createdEnvironment
	return resp, nil
}
