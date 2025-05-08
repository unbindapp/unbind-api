package storage_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type CreatePVCInput struct {
	server.BaseAuthInput
	Body *models.CreatePVCInput
}

type CreatePVCResponse struct {
	Body struct {
		Data *k8s.PVCInfo `json:"data"`
	}
}

func (self *HandlerGroup) CreatePVC(ctx context.Context, input *CreatePVCInput) (*CreatePVCResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	pvc, err := self.srv.StorageService.CreatePVC(ctx, user.ID, bearerToken, input.Body)
	if err != nil {
		return nil, self.handleErr(err)
	}

	response := &CreatePVCResponse{}
	response.Body.Data = pvc
	return response, nil
}
