package environments_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

type DeleteEnvironmentInput struct {
	server.BaseAuthInput
	Body struct {
		TeamID        uuid.UUID `json:"team_id" required:"true"`
		ProjectID     uuid.UUID `json:"project_id" required:"true"`
		EnvironmentID uuid.UUID `json:"environment_id" required:"true"`
	}
}

type DeleteEnvironmentResponse struct {
	Body struct {
		Data server.DeletedResponse `json:"data"`
	}
}

func (self *HandlerGroup) DeleteEnvironment(ctx context.Context, input *DeleteEnvironmentInput) (*DeleteEnvironmentResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	err := self.srv.EnvironmentService.DeleteEnvironmentByID(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &DeleteEnvironmentResponse{}
	resp.Body.Data = server.DeletedResponse{
		ID:      input.Body.EnvironmentID,
		Deleted: true,
	}
	return resp, nil
}
