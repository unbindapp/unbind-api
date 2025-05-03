package storage_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type GetS3EndpointByIDInput struct {
	server.BaseAuthInput
	ID          uuid.UUID `query:"id" format:"uuid" required:"true"`
	TeamID      uuid.UUID `query:"team_id" format:"uuid" required:"true"`
	WithBuckets bool      `query:"with_buckets"`
}

type GetS3EndpointByIDOutput struct {
	Body struct {
		Data *models.S3Response `json:"data"`
	}
}

func (self *HandlerGroup) GetS3EndpointByID(ctx context.Context, input *GetS3EndpointByIDInput) (*GetS3EndpointByIDOutput, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	s3Endpoint, err := self.srv.StorageService.GetS3StorageByID(
		ctx,
		user.ID,
		bearerToken,
		input.TeamID,
		input.ID,
		input.WithBuckets,
	)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &GetS3EndpointByIDOutput{}
	resp.Body.Data = s3Endpoint
	return resp, nil
}

type ListS3EndpointInput struct {
	server.BaseAuthInput
	TeamID      uuid.UUID `query:"team_id" format:"uuid" required:"true"`
	WithBuckets bool      `query:"with_buckets"`
}

type ListS3EndpointOutput struct {
	Body struct {
		Data []*models.S3Response `json:"data" nullable:"false"`
	}
}

func (self *HandlerGroup) ListS3Endpoint(ctx context.Context, input *ListS3EndpointInput) (*ListS3EndpointOutput, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	s3Endpoints, err := self.srv.StorageService.ListS3StorageBackends(
		ctx,
		user.ID,
		bearerToken,
		input.TeamID,
		input.WithBuckets,
	)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &ListS3EndpointOutput{}
	resp.Body.Data = s3Endpoints
	return resp, nil
}
