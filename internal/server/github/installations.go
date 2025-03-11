package github_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/log"
)

// GET Github app installations
type GithubAppInstallationListInput struct {
	AppID int64 `path:"app_id" required:"true"`
}

type GithubAppInstallationListResponse struct {
	Body []*ent.GithubInstallation
}

func (self *HandlerGroup) HandleListGithubAppInstallations(ctx context.Context, input *GithubAppInstallationListInput) (*GithubAppInstallationListResponse, error) {
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	// ! TODO - RBAC
	installations, err := self.srv.Repository.GetGithubInstallationsByCreator(ctx, user.ID)
	if err != nil {
		log.Error("Error getting github installations", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github installations")
	}

	resp := &GithubAppInstallationListResponse{}
	resp.Body = installations
	return resp, nil
}
