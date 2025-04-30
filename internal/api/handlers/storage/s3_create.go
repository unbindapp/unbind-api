package storage_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type CreateS3Input struct {
	server.BaseAuthInput
	Body *models.S3BackendCreateInput
}

type CreateS3Output struct {
	Body struct {
		Data *models.S3Response `json:"data"`
	}
}

func (self *HandlerGroup) CreateS3(ctx context.Context, input *CreateS3Input) (*CreateS3Output, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	s3source, err := self.srv.StorageService.CreateS3StorageBackend(
		ctx,
		user.ID,
		bearerToken,
		input.Body,
	)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &CreateS3Output{}
	resp.Body.Data = s3source
	return resp, nil
}
