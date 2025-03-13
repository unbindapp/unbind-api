package teams_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/k8s"
	"github.com/unbindapp/unbind-api/internal/log"
)

type TeamInput struct {
	Authorization string `header:"Authorization"`
}

type TeamResponse struct {
	Body struct {
		Data []k8s.UnbindTeam `json:"data"`
	}
}

// ListTeams handles GET /teams
func (self *HandlerGroup) ListTeams(ctx context.Context, input *TeamInput) (*TeamResponse, error) {
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	teams, err := self.srv.KubeClient.GetUnbindTeams(ctx, bearerToken)
	if err != nil {
		log.Error("Error getting teams", "err", err)
		return nil, huma.Error500InternalServerError("Unable to retrieve teams")
	}

	resp := &TeamResponse{}
	resp.Body.Data = teams
	return resp, nil
}
