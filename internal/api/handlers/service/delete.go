package service_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

type DeleteServiceInput struct {
	server.BaseAuthInput
	Body struct {
		TeamID        uuid.UUID `json:"team_id" required:"true"`
		ProjectID     uuid.UUID `json:"project_id" required:"true"`
		EnvironmentID uuid.UUID `json:"environment_id" required:"true"`
		ServiceID     uuid.UUID `json:"service_id" required:"true"`
	}
}

type DeleteServiceResponse struct {
	Body struct {
		Data struct {
			ID      uuid.UUID `json:"id"`
			Deleted bool      `json:"deleted"`
		} `json:"data"`
	}
}

func (self *HandlerGroup) DeleteService(ctx context.Context, input *DeleteServiceInput) (*DeleteServiceResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	err := self.srv.ServiceService.DeleteServiceByID(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.ServiceID)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &DeleteServiceResponse{}
	resp.Body.Data = struct {
		ID      uuid.UUID `json:"id"`
		Deleted bool      `json:"deleted"`
	}{
		ID:      input.Body.ServiceID,
		Deleted: true,
	}
	return resp, nil
}
