package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/log"
)

// Read user's admin repositories (that they can configure CI/CD on)
func (self *GithubClient) ReadUserAdminRepositories(ctx context.Context, installation *ent.GithubInstallation) ([]*github.Repository, error) {
	if installation == nil || installation.Edges.GithubApp == nil {
		return nil, fmt.Errorf("Invalid installation")
	}

	// Get authenticated client
	authenticatedClient, err := self.GetAuthenticatedClient(ctx, installation.GithubAppID, installation.ID, installation.Edges.GithubApp.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Error getting authenticated client: %v", err)
	}

	// Get user's organizations
	ghRepositories, _, err := authenticatedClient.Repositories.ListAll(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting user organizations: %v", err)
	}

	adminRepos := make([]*github.Repository, 0)
	for _, repo := range ghRepositories {
		if perms := repo.GetPermissions(); perms != nil {
			log.Infof("Repo %s perms: %v", repo.GetFullName(), perms)
			if isAdmin, ok := perms["admin"]; ok && isAdmin {
				adminRepos = append(adminRepos, repo)
			}
		}
	}
	return adminRepos, nil
}
