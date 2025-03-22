package service_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
	service_service "github.com/unbindapp/unbind-api/internal/services/service"
)

type CreateServiceInput struct {
	server.BaseAuthInput
	Body *service_service.CreateServiceInput
}

type CreateServiceResponse struct {
	Body struct {
		Data *models.ServiceResponse `json:"data"`
	}
}

func (self *HandlerGroup) CreateService(ctx context.Context, input *CreateServiceInput) (*CreateServiceResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	if input.Body == nil {
		return nil, huma.Error400BadRequest("Missing body")
	}

	println(input.Body.DisplayName)
	println(input.Body.RepositoryName)
	println(input.Body.RepositoryOwner)

	createdService, err := self.srv.ServiceService.CreateService(ctx, user.ID, input.Body, bearerToken)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &CreateServiceResponse{}
	resp.Body.Data = createdService
	return resp, nil
}
