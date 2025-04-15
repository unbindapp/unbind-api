package service_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
)

// Get containers for a service
type GetContainersInput struct {
	server.BaseAuthInput
	ServiceID     uuid.UUID `query:"service_id" required:"true"`
	TeamID        uuid.UUID `query:"team_id" required:"true"`
	ProjectID     uuid.UUID `query:"project_id" required:"true"`
	EnvironmentID uuid.UUID `query:"environment_id" required:"true"`
}

type GetContainersResponse struct {
	Body struct {
		Data []k8s.PodContainerStatus `json:"data"`
	}
}

// GetContainers gets containers/statuses for a service
func (self *HandlerGroup) GetContainers(ctx context.Context, input *GetContainersInput) (*GetContainersResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	containers, err := self.srv.ServiceService.GetServiceContainerStatuses(
		ctx,
		user.ID,
		bearerToken,
		input.TeamID,
		input.ProjectID,
		input.EnvironmentID,
		input.ServiceID,
	)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &GetContainersResponse{}
	resp.Body.Data = containers
	return resp, nil
}
