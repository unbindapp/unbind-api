package storage_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/models"
)

type UpdatePVCInput struct {
	server.BaseAuthInput
	Body *models.UpdatePVCInput
}

type UpdatePVCResponse struct {
	Body struct {
		Data *models.PVCInfo `json:"data"`
	}
}

func (self *HandlerGroup) UpdatePVC(ctx context.Context, input *UpdatePVCInput) (*UpdatePVCResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	pvc, err := self.srv.StorageService.UpdatePVC(ctx, user.ID, bearerToken, input.Body)
	if err != nil {
		return nil, self.handleErr(err)
	}

	response := &UpdatePVCResponse{}
	response.Body.Data = pvc
	return response, nil
}
