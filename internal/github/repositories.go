package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/log"
)

// Read user's admin repositories (that they can configure CI/CD on)
func (self *GithubClient) ReadUserAdminRepositories(ctx context.Context, installation *ent.GithubInstallation) ([]*GithubRepository, error) {
	if installation == nil || installation.Edges.GithubApp == nil {
		return nil, fmt.Errorf("Invalid installation")
	}

	// Get authenticated client
	authenticatedClient, err := self.GetAuthenticatedClient(ctx, installation.GithubAppID, installation.ID, installation.Edges.GithubApp.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Error getting authenticated client: %v", err)
	}

	// Get user's organizations
	// ! TODO - handle pagination
	ghRepositories, _, err := authenticatedClient.Repositories.ListByUser(ctx, installation.AccountLogin, nil)
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
	return formatRepositoryResponse(adminRepos), nil
}

type GithubRepositoryOwner struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type GithubRepository struct {
	ID       int64                 `json:"id"`
	FullName string                `json:"full_name"`
	HTMLURL  string                `json:"html_url"`
	CloneURL string                `json:"clone_url"`
	HomePage string                `json:"homepage"`
	Owner    GithubRepositoryOwner `json:"owner"`
}

func formatRepositoryResponse(repositories []*github.Repository) []*GithubRepository {
	response := make([]*GithubRepository, 0)
	for _, repository := range repositories {
		response = append(response, &GithubRepository{
			ID:       repository.GetID(),
			FullName: repository.GetFullName(),
			HTMLURL:  repository.GetHTMLURL(),
			CloneURL: repository.GetCloneURL(),
			HomePage: repository.GetHomepage(),
			Owner: GithubRepositoryOwner{
				ID:        repository.Owner.GetID(),
				Name:      repository.Owner.GetName(),
				Login:     repository.Owner.GetLogin(),
				AvatarURL: repository.Owner.GetAvatarURL(),
			},
		})
	}
	return response
}
