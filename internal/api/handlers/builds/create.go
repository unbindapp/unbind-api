package builds_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type CreateBuildInput struct {
	server.BaseAuthInput
	Body struct {
		models.CreateBuildJobInput
	}
}

type CreateBuildOutput struct {
	Body struct {
		Data *models.BuildJobResponse `json:"data"`
	}
}

func (self *HandlerGroup) CreateBuild(ctx context.Context, input *CreateBuildInput) (*CreateBuildOutput, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	// Create build job
	buildJob, err := self.srv.BuildJobService.CreateManualBuildJob(ctx, user.ID, &input.Body.CreateBuildJobInput)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &CreateBuildOutput{}
	resp.Body.Data = buildJob
	return resp, nil
}
