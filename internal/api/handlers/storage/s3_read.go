package storage_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/models"
)

type GetS3SourceByIDInput struct {
	server.BaseAuthInput
	ID          uuid.UUID `query:"id" format:"uuid" required:"true"`
	TeamID      uuid.UUID `query:"team_id" format:"uuid" required:"true"`
	WithBuckets bool      `query:"with_buckets"`
}

type GetS3SourceByIDOutput struct {
	Body struct {
		Data *models.S3Response `json:"data"`
	}
}

func (self *HandlerGroup) GetS3SourceByID(ctx context.Context, input *GetS3SourceByIDInput) (*GetS3SourceByIDOutput, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	s3Source, err := self.srv.StorageService.GetS3StorageByID(
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

	resp := &GetS3SourceByIDOutput{}
	resp.Body.Data = s3Source
	return resp, nil
}

type ListS3SourceInput struct {
	server.BaseAuthInput
	TeamID      uuid.UUID `query:"team_id" format:"uuid" required:"true"`
	WithBuckets bool      `query:"with_buckets"`
}

type ListS3SourceOutput struct {
	Body struct {
		Data []*models.S3Response `json:"data" nullable:"false"`
	}
}

func (self *HandlerGroup) ListS3Source(ctx context.Context, input *ListS3SourceInput) (*ListS3SourceOutput, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	s3Sources, err := self.srv.StorageService.ListS3StorageBackends(
		ctx,
		user.ID,
		bearerToken,
		input.TeamID,
		input.WithBuckets,
	)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &ListS3SourceOutput{}
	resp.Body.Data = s3Sources
	return resp, nil
}
