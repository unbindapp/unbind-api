package service_handler

import (
	"context"
	"errors"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
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

	createdService, err := self.srv.ServiceService.CreateService(ctx, user.ID, input.Body, bearerToken)
	if err != nil {
		if errors.Is(err, errdefs.ErrInvalidInput) {
			return nil, huma.Error400BadRequest(err.Error())
		}
		if errors.Is(err, errdefs.ErrUnauthorized) {
			return nil, huma.Error403Forbidden("Unauthorized")
		}
		if ent.IsNotFound(err) || errors.Is(err, errdefs.ErrNotFound) {
			return nil, huma.Error404NotFound(err.Error())
		}
		log.Error("Error creating service", "err", err)
		return nil, huma.Error500InternalServerError("Unable to create service")
	}

	resp := &CreateServiceResponse{}
	resp.Body.Data = createdService
	return resp, nil
}
