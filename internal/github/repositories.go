package github

import (
	"context"
	"fmt"
	"sync"
	"time"

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
	Name     string                `json:"name"`
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
			Name:     repository.GetName(),
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

// ! Getting repository details
// GithubRepositoryDetail contains detailed information about a repository
type GithubRepositoryDetail struct {
	ID              int64                  `json:"id"`
	Name            string                 `json:"name"`
	FullName        string                 `json:"fullName"`
	Description     string                 `json:"description"`
	URL             string                 `json:"url"`
	HTMLURL         string                 `json:"htmlUrl"`
	DefaultBranch   string                 `json:"defaultBranch"`
	Language        string                 `json:"language"`
	Private         bool                   `json:"private"`
	Fork            bool                   `json:"fork"`
	Archived        bool                   `json:"archived"`
	Disabled        bool                   `json:"disabled"`
	Size            int                    `json:"size"`
	StargazersCount int                    `json:"stargazersCount"`
	WatchersCount   int                    `json:"watchersCount"`
	ForksCount      int                    `json:"forksCount"`
	OpenIssuesCount int                    `json:"openIssuesCount"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	PushedAt        time.Time              `json:"pushedAt"`
	Branches        []*GithubBranch        `json:"branches"`
	Tags            []*GithubTag           `json:"tags"`
	Owner           *GithubRepositoryOwner `json:"owner"`
}

// GithubBranch contains information about a branch
type GithubBranch struct {
	Name      string `json:"name"`
	Protected bool   `json:"protected"`
	SHA       string `json:"sha"`
}

// GithubTag contains information about a tag
type GithubTag struct {
	Name string `json:"name"`
	SHA  string `json:"sha"`
}

// GetRepositoryDetail retrieves detailed information about a GitHub repository
func (self *GithubClient) GetRepositoryDetail(ctx context.Context, installation *ent.GithubInstallation, owner, repo string) (*GithubRepositoryDetail, error) {
	if installation == nil || installation.Edges.GithubApp == nil {
		return nil, fmt.Errorf("invalid installation: missing app edge or nil")
	}

	// Get authenticated client
	authenticatedClient, err := self.GetAuthenticatedClient(ctx, installation.GithubAppID, installation.ID, installation.Edges.GithubApp.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error getting authenticated client for %s: %v", installation.AccountLogin, err)
	}

	// Get repository information
	ghRepo, _, err := authenticatedClient.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("error getting repository %s/%s: %v", owner, repo, err)
	}

	// Get branches
	branches, err := self.getRepositoryBranches(ctx, authenticatedClient, owner, repo)
	if err != nil {
		log.Warn("Error getting repository branches", "err", err, "owner", owner, "repo", repo)
		// Continue execution rather than failing
	}

	// Get tags
	tags, err := self.getRepositoryTags(ctx, authenticatedClient, owner, repo)
	if err != nil {
		log.Warn("Error getting repository tags", "err", err, "owner", owner, "repo", repo)
		// Continue execution rather than failing
	}

	// Format the response
	repoDetail := &GithubRepositoryDetail{
		ID:              ghRepo.GetID(),
		Name:            ghRepo.GetName(),
		FullName:        ghRepo.GetFullName(),
		Description:     ghRepo.GetDescription(),
		URL:             ghRepo.GetURL(),
		HTMLURL:         ghRepo.GetHTMLURL(),
		DefaultBranch:   ghRepo.GetDefaultBranch(),
		Language:        ghRepo.GetLanguage(),
		Private:         ghRepo.GetPrivate(),
		Fork:            ghRepo.GetFork(),
		Archived:        ghRepo.GetArchived(),
		Disabled:        ghRepo.GetDisabled(),
		Size:            ghRepo.GetSize(),
		StargazersCount: ghRepo.GetStargazersCount(),
		WatchersCount:   ghRepo.GetWatchersCount(),
		ForksCount:      ghRepo.GetForksCount(),
		OpenIssuesCount: ghRepo.GetOpenIssuesCount(),
		CreatedAt:       ghRepo.GetCreatedAt().Time,
		UpdatedAt:       ghRepo.GetUpdatedAt().Time,
		PushedAt:        ghRepo.GetPushedAt().Time,
		Branches:        branches,
		Tags:            tags,
	}

	// Add owner information if available
	if ghRepo.Owner != nil {
		repoDetail.Owner = &GithubRepositoryOwner{
			ID:        ghRepo.Owner.GetID(),
			Name:      ghRepo.Owner.GetName(),
			Login:     ghRepo.Owner.GetLogin(),
			AvatarURL: ghRepo.Owner.GetAvatarURL(),
		}
	}

	return repoDetail, nil
}

// Get branches for a repository
func (self *GithubClient) getRepositoryBranches(ctx context.Context, client *github.Client, owner, repo string) ([]*GithubBranch, error) {
	opts := &github.BranchListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allBranches []*GithubBranch
	for {
		branches, resp, err := client.Repositories.ListBranches(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("error listing branches: %v", err)
		}

		for _, branch := range branches {
			allBranches = append(allBranches, &GithubBranch{
				Name:      branch.GetName(),
				Protected: branch.GetProtected(),
				SHA:       branch.GetCommit().GetSHA(),
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allBranches, nil
}

// getRepositoryTags retrieves a list of tags for a repository
func (self *GithubClient) getRepositoryTags(ctx context.Context, client *github.Client, owner, repo string) ([]*GithubTag, error) {
	opts := &github.ListOptions{
		PerPage: 100,
	}

	var allTags []*GithubTag
	for {
		tags, resp, err := client.Repositories.ListTags(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("error listing tags: %v", err)
		}

		for _, tag := range tags {
			allTags = append(allTags, &GithubTag{
				Name: tag.GetName(),
				SHA:  tag.GetCommit().GetSHA(),
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allTags, nil
}
