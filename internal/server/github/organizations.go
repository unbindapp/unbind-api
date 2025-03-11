package github_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/internal/log"
)

// GET Github admin organizations for installation
type GithubAdminOrganizationListInput struct {
	InstallationID int64 `path:"installation_id" required:"true"`
}

type GithubAdminOrganizationListResponse struct {
	Body []*github.Organization
}

func (self *HandlerGroup) HandleListGithubAdminOrganizations(ctx context.Context, input *GithubAdminOrganizationListInput) (*GithubAdminOrganizationListResponse, error) {
	// Get installation
	installation, err := self.srv.Repository.GetGithubInstallationByID(ctx, input.InstallationID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, huma.Error404NotFound("Installation not found")
		}
		log.Error("Error getting github installation", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github installation")
	}
	if installation.AccountType != githubinstallation.AccountTypeUser {
		return nil, huma.Error400BadRequest("Invalid installation type")
	}

	adminOrgs, err := self.srv.GithubClient.ReadUserAdminOrganizations(ctx, installation)
	if err != nil {
		log.Error("Error getting user admin organizations", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get user admin organizations")
	}

	resp := &GithubAdminOrganizationListResponse{}
	resp.Body = adminOrgs
	return resp, nil
}
