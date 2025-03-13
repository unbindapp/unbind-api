package teams_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/log"
	team_service "github.com/unbindapp/unbind-api/internal/services/team"
)

type TeamInput struct {
	Authorization string `header:"Authorization"`
}

type TeamResponse struct {
	Body struct {
		Data []*team_service.GetTeamResponse `json:"data"`
	}
}

// ListTeams handles GET /teams
func (self *HandlerGroup) ListTeams(ctx context.Context, input *TeamInput) (*TeamResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	// Get token
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	teams, err := self.srv.TeamService.ListTeams(ctx, user.ID, bearerToken)
	if err != nil {
		log.Error("Error getting teams", "err", err)
		return nil, huma.Error500InternalServerError("Unable to retrieve teams")
	}

	resp := &TeamResponse{}
	resp.Body.Data = teams
	return resp, nil
}
