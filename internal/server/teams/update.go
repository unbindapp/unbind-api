package teams_handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/errdefs"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/server"
	team_service "github.com/unbindapp/unbind-api/internal/services/team"
)

type UpdateTeamInput struct {
	server.BaseAuthInput
	TeamID uuid.UUID `path:"team_id"`
	Body   struct {
		DisplayName string `json:"display_name"`
	}
}

type UpdateTeamResponse struct {
	Body struct {
		Data *team_service.GetTeamResponse `json:"data"`
	}
}

// UpdateTeam handles PUT /team/{team_id}
func (self *HandlerGroup) UpdateTeam(ctx context.Context, input *UpdateTeamInput) (*UpdateTeamResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	updatedTeam, err := self.srv.TeamService.UpdateTeam(ctx, user.ID, &team_service.TeamUpdateInput{
		ID:          input.TeamID,
		DisplayName: input.Body.DisplayName,
	})
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, huma.Error404NotFound("Team not found")
		}
		if errors.Is(err, errdefs.ErrUnauthorized) {
			return nil, huma.Error403Forbidden("Unauthorized")
		}
		log.Error("Error getting teams", "err", err)
		return nil, huma.Error500InternalServerError("Unable to update team")
	}

	resp := &UpdateTeamResponse{}
	resp.Body.Data = updatedTeam
	return resp, nil
}
