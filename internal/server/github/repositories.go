package github_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/github"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/server"
)

// GET Github admin organizations for installation
type GithubAdminRepositoryListResponse struct {
	Body struct {
		Data []*github.GithubRepository `json:"data"`
	}
}

func (self *HandlerGroup) HandleListGithubAdminRepositories(ctx context.Context, input *server.BaseAuthInput) (*GithubAdminRepositoryListResponse, error) {
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	// Get owned installations
	// ! TODO - group RBAC
	installations, err := self.srv.Repository.Github().GetInstallationsByCreator(ctx, user.ID)
	if err != nil {
		log.Error("Error getting github installation", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github installation")
	}
	if len(installations) == 0 {
		return &GithubAdminRepositoryListResponse{
			Body: struct {
				Data []*github.GithubRepository `json:"data"`
			}{
				Data: []*github.GithubRepository{},
			},
		}, nil
	}

	adminRepos, err := self.srv.GithubClient.ReadUserAdminRepositories(ctx, installations)
	if err != nil {
		log.Error("Error getting user admin organizations", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get user admin organizations")
	}

	resp := &GithubAdminRepositoryListResponse{}
	resp.Body.Data = adminRepos
	return resp, nil
}
