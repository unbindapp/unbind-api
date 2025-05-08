package storage_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

type DeleteS3EndpointByIDInput struct {
	server.BaseAuthInput
	Body struct {
		ID     uuid.UUID `json:"id" format:"uuid" required:"true"`
		TeamID uuid.UUID `json:"team_id" format:"uuid" required:"true"`
	}
}

type DeleteS3EndpointByIDOutput struct {
	Body struct {
		Data server.DeletedResponse `json:"data"`
	}
}

func (self *HandlerGroup) DeleteS3Endpoint(ctx context.Context, input *DeleteS3EndpointByIDInput) (*DeleteS3EndpointByIDOutput, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	err := self.srv.StorageService.DeleteS3StorageByID(
		ctx,
		user.ID,
		bearerToken,
		input.Body.TeamID,
		input.Body.ID,
	)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &DeleteS3EndpointByIDOutput{}
	resp.Body.Data = server.DeletedResponse{
		ID:      input.Body.ID.String(),
		Deleted: true,
	}
	return resp, nil
}
