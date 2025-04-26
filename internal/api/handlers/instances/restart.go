package instances_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

// Restart instance
type RestartInstancesInput struct {
	server.BaseAuthInput
	Body struct {
		ServiceID     uuid.UUID `json:"service_id" required:"true"`
		TeamID        uuid.UUID `json:"team_id" required:"true"`
		ProjectID     uuid.UUID `json:"project_id" required:"true"`
		EnvironmentID uuid.UUID `json:"environment_id" required:"true"`
	}
}

type Restarted struct {
	Restarted bool `json:"restarted"`
}

type RestartServicesResponse struct {
	Body struct {
		Data *Restarted `json:"data"`
	}
}

// RestartInstances handles PUT /instances/restart
func (self *HandlerGroup) RestartInstances(ctx context.Context, input *RestartInstancesInput) (*RestartServicesResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	err := self.srv.ServiceService.RestartServiceByID(
		ctx,
		user.ID,
		bearerToken,
		input.Body.TeamID,
		input.Body.ProjectID,
		input.Body.EnvironmentID,
		input.Body.ServiceID,
	)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &RestartServicesResponse{}
	resp.Body.Data = &Restarted{
		Restarted: true,
	}
	return resp, nil
}
