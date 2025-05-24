package teams_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/models"
)

type TeamResponse struct {
	Body struct {
		Data []*models.TeamResponse `json:"data" nullable:"false"`
	}
}

// ListTeams handles GET /teams
func (self *HandlerGroup) ListTeams(ctx context.Context, input *server.BaseAuthInput) (*TeamResponse, error) {
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
		return nil, self.handleErr(err)
	}

	resp := &TeamResponse{}
	resp.Body.Data = teams
	return resp, nil
}

// Get by ID
type GetTeamInput struct {
	server.BaseAuthInput
	TeamID uuid.UUID `query:"team_id" description:"The ID of the team to retrieve" required:"true"`
}

type GetTeamResponse struct {
	Body struct {
		Data *models.TeamResponse `json:"data" nullable:"false"`
	}
}

// GetTeam handles GET /teams/get
func (self *HandlerGroup) GetTeam(ctx context.Context, input *GetTeamInput) (*GetTeamResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	team, err := self.srv.TeamService.GetTeamByID(ctx, user.ID, input.TeamID)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &GetTeamResponse{}
	resp.Body.Data = team
	return resp, nil
}
