package github_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

// GET Github app installations
type GithubAppInstallationListResponse struct {
	Body struct {
		Data []*ent.GithubInstallation `json:"data"`
	}
}

func (self *HandlerGroup) HandleListGithubAppInstallations(ctx context.Context, input *server.BaseAuthInput) (*GithubAppInstallationListResponse, error) {
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	// ! TODO - RBAC
	installations, err := self.srv.Repository.Github().GetInstallationsByCreator(ctx, user.ID)
	if err != nil {
		log.Error("Error getting github installations", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github installations")
	}

	resp := &GithubAppInstallationListResponse{}
	resp.Body.Data = installations
	return resp, nil
}
