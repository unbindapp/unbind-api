package system_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	system_service "github.com/unbindapp/unbind-api/internal/services/system"
)

type BuildkitSettingsUpdateInput struct {
	server.BaseAuthInput
	Body struct {
		Replicas    int `json:"replicas" nullable:"false"`
		Parallelism int `json:"max_parallelism" nullable:"false"`
	}
}

type BuildkitSettingsUpdateResponse struct {
	Body struct {
		Data *system_service.BuildkitSettingsResponse `json:"settings" nullable:"false"`
	}
}

func (self *HandlerGroup) UpdateBuildkitSettings(ctx context.Context, input *BuildkitSettingsUpdateInput) (*BuildkitSettingsUpdateResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	settings, err := self.srv.SystemService.UpdateBuildkitSettings(ctx, user.ID, input.Body.Replicas, input.Body.Parallelism)
	if err != nil {
		return nil, self.handleErr(err)
	}
	resp := &BuildkitSettingsUpdateResponse{}
	resp.Body.Data = settings
	return resp, nil
}
