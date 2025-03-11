package github_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/github"
	"github.com/unbindapp/unbind-api/internal/log"
)

// GET Github admin organizations for installation
type GithubAdminRepositoryListInput struct {
	InstallationID int64 `path:"installation_id" required:"true"`
}

type GithubAdminRepositoryListResponse struct {
	Body []*github.GithubRepository
}

func (self *HandlerGroup) HandleListGithubAdminRepositories(ctx context.Context, input *GithubAdminRepositoryListInput) (*GithubAdminRepositoryListResponse, error) {
	// Get installation
	installation, err := self.srv.Repository.GetGithubInstallationByID(ctx, input.InstallationID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, huma.Error404NotFound("Installation not found")
		}
		log.Error("Error getting github installation", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github installation")
	}

	adminRepos, err := self.srv.GithubClient.ReadUserAdminRepositories(ctx, installation)
	if err != nil {
		log.Error("Error getting user admin organizations", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get user admin organizations")
	}

	resp := &GithubAdminRepositoryListResponse{}
	resp.Body = adminRepos
	return resp, nil
}
