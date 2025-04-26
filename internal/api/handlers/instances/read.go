package instances_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// List instances (pods) for a service
type ListInstancesInput struct {
	server.BaseAuthInput
	models.InstanceStatusInput
}

type ListInstancesResponse struct {
	Body struct {
		Data []k8s.PodContainerStatus `json:"data" nullable:"false"`
	}
}

// ListInstances gets pods/statuses for a service
func (self *HandlerGroup) ListInstances(ctx context.Context, input *ListInstancesInput) (*ListInstancesResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	containers, err := self.srv.InstanceService.GetInstanceStatuses(
		ctx,
		user.ID,
		bearerToken,
		&input.InstanceStatusInput,
	)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &ListInstancesResponse{}
	resp.Body.Data = containers
	return resp, nil
}
