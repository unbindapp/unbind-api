package github

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/internal/log"
)

// Read user's admin repositories (that they can configure CI/CD on)
func (self *GithubClient) ReadUserAdminRepositories(ctx context.Context, installations []*ent.GithubInstallation) ([]*GithubRepository, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	allAdminRepos := make([]*GithubRepository, 0)
	errChan := make(chan error, len(installations))

	for _, installation := range installations {
		if installation == nil || installation.Edges.GithubApp == nil {
			log.Warnf("Invalid installation to read repositories from, missing appedge or nil: %v", installation)
			continue
		}
		wg.Add(1)
		go func(inst *ent.GithubInstallation) {
			defer wg.Done()
			// Get authenticated client
			authenticatedClient, err := self.GetAuthenticatedClient(ctx, inst.GithubAppID, inst.ID, inst.Edges.GithubApp.PrivateKey)
			if err != nil {
				errChan <- fmt.Errorf("Error getting authenticated client for %s: %v", inst.AccountLogin, err)
				return
			}
			// Get user's repositories
			// ! TODO - handle pagination
			ghRepositories, _, err := authenticatedClient.Repositories.ListByUser(ctx, inst.AccountLogin, nil)
			if err != nil {
				errChan <- fmt.Errorf("Error getting repositories for user %s: %v", inst.AccountLogin, err)
				return
			}
			adminRepos := make([]*github.Repository, 0)
			for _, repo := range ghRepositories {
				if inst.AccountType == githubinstallation.AccountTypeUser {
					if repo.GetOwner().GetID() == inst.AccountID {
						adminRepos = append(adminRepos, repo)
						continue
					}
				}
				// ! TODO - figure out organization owners?
				if perms := repo.GetPermissions(); perms != nil {
					log.Infof("Repo %s perms: %v", repo.GetFullName(), perms)
					if isAdmin, ok := perms["admin"]; ok && isAdmin {
						adminRepos = append(adminRepos, repo)
					}
				}
			}
			// Format and add to the result slice in a thread-safe way
			formattedRepos := formatRepositoryResponse(adminRepos)
			mu.Lock()
			allAdminRepos = append(allAdminRepos, formattedRepos...)
			mu.Unlock()
		}(installation)
	}
	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check if any errors occurred
	for err := range errChan {
		if err != nil {
			return nil, err // Return the first error encountered
		}
	}

	// Remove any duplicates
	return removeDuplicateRepositories(allAdminRepos), nil
}

// removeDuplicateRepositories removes duplicate repositories from the slice
// based on the repository ID
func removeDuplicateRepositories(repos []*GithubRepository) []*GithubRepository {
	seen := make(map[int64]bool)
	result := make([]*GithubRepository, 0, len(repos))

	for _, repo := range repos {
		if !seen[repo.ID] {
			seen[repo.ID] = true
			result = append(result, repo)
		}
	}

	return result
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
