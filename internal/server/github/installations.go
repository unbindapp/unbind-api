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
	installations, err := self.srv.Repository.GetGithubInstallationsByAppID(ctx, input.AppID)
	if err != nil {
		log.Error("Error getting github installations", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github installations")
	}

	resp := &GithubAppInstallationListResponse{}
	resp.Body = installations
	return resp, nil
}
