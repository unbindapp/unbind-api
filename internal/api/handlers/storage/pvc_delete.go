package storage_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type DeletePVCInput struct {
	server.BaseAuthInput
	Body *models.DeletePVCInput
}

type DeletePVCResponse struct {
	Body struct {
		Data server.DeletedResponse `json:"data"`
	}
}

func (self *HandlerGroup) DeletePVC(ctx context.Context, input *DeletePVCInput) (*DeletePVCResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	err := self.srv.StorageService.DeletePVC(ctx, user.ID, bearerToken, input.Body)
	if err != nil {
		return nil, self.handleErr(err)
	}

	response := &DeletePVCResponse{}
	response.Body.Data = server.DeletedResponse{
		ID:      input.Body.ID,
		Deleted: true,
	}
	return response, nil
}
