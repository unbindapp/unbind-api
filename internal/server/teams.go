package server

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/kubeclient"
	"github.com/unbindapp/unbind-api/internal/log"
)

type TeamResponse struct {
	Body struct {
		Data []kubeclient.UnbindTeam `json:"data"`
	}
}

// ListTeams handles GET /teams
func (self *Server) ListTeams(ctx context.Context, _ *EmptyInput) (*TeamResponse, error) {
	teams, err := self.KubeClient.GetUnbindTeams()
	if err != nil {
		log.Error("Error getting teams", "err", err)
		return nil, huma.Error500InternalServerError("Unable to retrieve teams")
	}

	resp := &TeamResponse{}
	resp.Body.Data = teams
	return resp, nil
}
