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

type ListS3SourceInput struct {
	server.BaseAuthInput
	TeamID uuid.UUID `query:"team_id" format:"uuid" required:"true"`
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
	)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &ListS3SourceOutput{}
	resp.Body.Data = s3Sources
	return resp, nil
}
