package servicegroups_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type DeleteServiceGroupInput struct {
	server.BaseAuthInput
	Body *models.DeleteServiceGroupInput
}

type DeleteServiceGroupResponse struct {
	Body struct {
		Data server.DeletedResponse `json:"data"`
	}
}

func (self *HandlerGroup) DeleteServiceGroup(ctx context.Context, input *DeleteServiceGroupInput) (*DeleteServiceGroupResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	err := self.srv.ServiceGroupService.DeleteServiceGroup(ctx, user.ID, input.Body)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &DeleteServiceGroupResponse{}
	resp.Body.Data = server.DeletedResponse{
		ID:      input.Body.ID.String(),
		Deleted: true,
	}
	return resp, nil
}
