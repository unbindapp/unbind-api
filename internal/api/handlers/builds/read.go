package builds_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type ListBuildJobsInput struct {
	server.BaseAuthInput
	models.GetBuildJobsInput
}

type ListBuildJobResponseData struct {
	Jobs     []*models.BuildJobResponse         `json:"jobs"`
	Metadata *models.PaginationResponseMetadata `json:"metadata"`
}

type ListBuildJobsResponse struct {
	Body struct {
		Data *ListBuildJobResponseData `json:"data"`
	}
}

func (self *HandlerGroup) ListBuildJobs(ctx context.Context, input *ListBuildJobsInput) (*ListBuildJobsResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	response, metadata, err := self.srv.BuildJobService.GetBuildJobsForService(ctx, user.ID, &input.GetBuildJobsInput)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &ListBuildJobsResponse{}
	resp.Body.Data.Jobs = response
	resp.Body.Data.Metadata = metadata
	return resp, nil
}
