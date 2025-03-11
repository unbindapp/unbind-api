package teams_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/kubeclient"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/server"
)

type TeamResponse struct {
	Body struct {
		Data []kubeclient.UnbindTeam `json:"data"`
	}
}

// ListTeams handles GET /teams
func (self *HandlerGroup) ListTeams(ctx context.Context, _ *server.EmptyInput) (*TeamResponse, error) {
	teams, err := self.srv.KubeClient.GetUnbindTeams()
	if err != nil {
		log.Error("Error getting teams", "err", err)
		return nil, huma.Error500InternalServerError("Unable to retrieve teams")
	}

	resp := &TeamResponse{}
	resp.Body.Data = teams
	return resp, nil
}
