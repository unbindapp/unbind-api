package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"golang.org/x/sync/errgroup"
)

// CursorPaginationParams contains parameters for cursor-based pagination
type CursorPaginationParams struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

// CursorPaginationResponse contains pagination metadata with cursor
type CursorPaginationResponse struct {
	NextCursor     string `json:"next_cursor"`
	PreviousCursor string `json:"previous_cursor,omitempty"`
	Limit          int    `json:"limit"`
	HasMore        bool   `json:"has_more"`
}

// RepositoriesCursorResponse contains paginated repositories with cursor info
type RepositoriesCursorResponse struct {
	Repositories []*GithubRepository      `json:"repositories"`
	Pagination   CursorPaginationResponse `json:"pagination"`
}

type cursorData struct {
	InstallationPages map[int64]int `json:"i_p"`
	PrevCursor        string        `json:"prev"`
}

// Read user's admin repositories (that they can configure CI/CD on)
func (self *GithubClient) ReadUserAdminRepositoriesCursor(
	ctx context.Context,
	installations []*ent.GithubInstallation,
	params CursorPaginationParams,
) (*RepositoriesCursorResponse, error) {
	// Validate and set defaults
	if params.Limit < 1 || params.Limit > 30 {
		params.Limit = 30 // Default
	}

	// Process cursor
	installationPages := make(map[int64]int)
	var prevCursor string

	if params.Cursor != "" {
		cursor, err := decodeCursor(params.Cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %v", err)
		}
		installationPages = cursor.InstallationPages
		prevCursor = params.Cursor
	} else {
		for _, inst := range installations {
			if inst != nil {
				installationPages[inst.ID] = 1
			}
		}
	}

	var mu sync.Mutex
	allAdminRepos := make([]*GithubRepository, 0)
	nextInstallationPages := make(map[int64]int)
	var hasMore bool

	// limit concurrency
	const maxConcurrency = 3
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrency)

	// Process installations concurrently with a limit
	for _, installation := range installations {
		if installation == nil || installation.Edges.GithubApp == nil {
			log.Warnf("Invalid installation to read repositories from, missing appedge or nil: %v", installation)
			continue
		}

		// Copy the installation variable to avoid closure issues
		inst := installation

		// Check if we need to process this installation
		page, exists := installationPages[inst.ID]
		if !exists {
			continue
		}

		g.Go(func() error {
			authenticatedClient, err := self.GetAuthenticatedClient(gctx, inst.GithubAppID, inst.ID, inst.Edges.GithubApp.PrivateKey)
			if err != nil {
				return fmt.Errorf("Error getting authenticated client for %s: %v", inst.AccountLogin, err)
			}
			defer authenticatedClient.Client().CloseIdleConnections()

			// Process repositories with proper pagination
			adminRepos, nextPage, err := self.fetchUserAdminRepos(gctx, authenticatedClient, inst, page, params.Limit)
			if err != nil {
				return err
			}

			if nextPage > 0 {
				mu.Lock()
				nextInstallationPages[inst.ID] = nextPage
				hasMore = true
				mu.Unlock()
			}

			// Format and add to the result slice in a thread-safe way
			mu.Lock()
			allAdminRepos = append(allAdminRepos, adminRepos...)
			mu.Unlock()

			return nil
		})
	}

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		return nil, err
	}

	uniqueRepos := removeDuplicateRepositories(allAdminRepos)

	if len(uniqueRepos) > params.Limit {
		uniqueRepos = uniqueRepos[:params.Limit]
		hasMore = true
	}

	var nextCursor string
	if hasMore {
		cursorData := cursorData{
			InstallationPages: nextInstallationPages,
			PrevCursor:        prevCursor,
		}
		nextCursor, _ = encodeCursor(cursorData)
	}

	return &RepositoriesCursorResponse{
		Repositories: uniqueRepos,
		Pagination: CursorPaginationResponse{
			NextCursor:     nextCursor,
			PreviousCursor: prevCursor,
			Limit:          params.Limit,
			HasMore:        hasMore,
		},
	}, nil
}

// fetchUserAdminRepos fetches all admin repositories for a user with proper pagination and memory limits
func (self *GithubClient) fetchUserAdminRepos(ctx context.Context,
	client *github.Client,
	inst *ent.GithubInstallation,
	page int,
	perPage int,
) ([]*GithubRepository, int, error) {
	adminRepos := make([]*github.Repository, 0)
	opts := &github.RepositoryListByUserOptions{
		Sort: "updated",
		ListOptions: github.ListOptions{
			PerPage: perPage,
			Page:    page,
		},
	}

	// Get user's repositories with pagination
	ghRepositories, resp, err := client.Repositories.ListByUser(ctx, inst.AccountLogin, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("Error getting repositories for user %s: %v", inst.AccountLogin, err)
	}

	// Filter admin repositories
	for _, repo := range ghRepositories {
		isAdmin := false

		if inst.AccountType == githubinstallation.AccountTypeUser {
			if repo.GetOwner().GetID() == inst.AccountID {
				isAdmin = true
			}
		}

		// Check permissions
		if !isAdmin {
			if perms := repo.GetPermissions(); perms != nil {
				if admin, ok := perms["admin"]; ok && admin {
					isAdmin = true
				}
			}
		}

		if isAdmin {
			adminRepos = append(adminRepos, repo)
		}
	}

	// Check if there are more pages
	nextPage := 0
	if resp.NextPage != 0 {
		nextPage = resp.NextPage
	}

	return formatRepositoryResponse(adminRepos), nextPage, nil
}

func encodeCursor(data cursorData) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(jsonData), nil
}

func decodeCursor(cursor string) (cursorData, error) {
	var data cursorData
	jsonData, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// removeDuplicateRepositories removes duplicate repositories from the slice
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

// * Formatting
type GithubRepositoryOwner struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type GithubRepository struct {
	ID        int64                 `json:"id"`
	Name      string                `json:"name"`
	FullName  string                `json:"full_name"`
	HTMLURL   string                `json:"html_url"`
	CloneURL  string                `json:"clone_url"`
	HomePage  string                `json:"homepage"`
	Owner     GithubRepositoryOwner `json:"owner"`
	UpdatedAt time.Time             `json:"updated_at"`
}

func formatRepositoryResponse(repositories []*github.Repository) []*GithubRepository {
	// Pre-allocate the slice to the exact size needed to avoid reallocations
	response := make([]*GithubRepository, 0, len(repositories))

	for _, repository := range repositories {
		// Skip repositories with nil owners to prevent panic
		if repository.Owner == nil {
			log.Warnf("Skipping repository with nil owner: %s", repository.GetName())
			continue
		}

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
			UpdatedAt: repository.GetUpdatedAt().Time,
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
	Ref       string `json:"ref"`
	Protected bool   `json:"protected"`
	SHA       string `json:"sha"`
}

// GithubTag contains information about a tag
type GithubTag struct {
	Name string `json:"name"`
	Ref  string `json:"ref"`
	SHA  string `json:"sha"`
}

// Get details for a repository
func (self *GithubClient) GetRepositoryDetail(ctx context.Context, installation *ent.GithubInstallation, owner, repo string) (*GithubRepositoryDetail, error) {
	if installation == nil || installation.Edges.GithubApp == nil {
		return nil, fmt.Errorf("invalid installation: missing app edge or nil")
	}

	// We'll use a client with a timeout context to ensure we don't hang indefinitely
	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Get authenticated client
	authenticatedClient, err := self.GetAuthenticatedClient(timeoutCtx, installation.GithubAppID, installation.ID, installation.Edges.GithubApp.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error getting authenticated client for %s: %v", installation.AccountLogin, err)
	}
	defer authenticatedClient.Client().CloseIdleConnections()

	// Get repository information
	ghRepo, _, err := authenticatedClient.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("error getting repository %s/%s: %v", owner, repo, err)
	}

	var branches []*GithubBranch
	var tags []*GithubTag
	var wg sync.WaitGroup
	var branchErr, tagErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		branches, branchErr = self.getRepositoryBranches(timeoutCtx, authenticatedClient, owner, repo)
		if branchErr != nil {
			log.Warn("Error getting repository branches", "err", branchErr, "owner", owner, "repo", repo)
			branchErr = nil // Don't fail the entire request for this
		}
	}()

	go func() {
		defer wg.Done()
		tags, tagErr = self.getRepositoryTags(timeoutCtx, authenticatedClient, owner, repo)
		if tagErr != nil {
			log.Warn("Error getting repository tags", "err", tagErr, "owner", owner, "repo", repo)
			tagErr = nil // Don't fail the entire request for this
		}
	}()

	wg.Wait()

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

// Get branches for a repository - improved with context handling
func (self *GithubClient) getRepositoryBranches(ctx context.Context, client *github.Client, owner, repo string) ([]*GithubBranch, error) {
	opts := &github.BranchListOptions{
		ListOptions: github.ListOptions{
			PerPage: 50,
		},
	}

	var allBranches []*GithubBranch
	for {
		// Check if context is canceled
		select {
		case <-ctx.Done():
			return allBranches, ctx.Err()
		default:
		}

		branches, resp, err := client.Repositories.ListBranches(ctx, owner, repo, opts)
		if err != nil {
			return allBranches, fmt.Errorf("error listing branches: %v", err)
		}
		allBranches = make([]*GithubBranch, len(branches))

		for i, branch := range branches {
			allBranches[i] = &GithubBranch{
				Name:      branch.GetName(),
				Protected: branch.GetProtected(),
				SHA:       branch.GetCommit().GetSHA(),
				Ref:       fmt.Sprintf("refs/heads/%s", branch.GetName()),
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allBranches, nil
}

// getRepositoryTags retrieves a list of tags for a repository - improved with context handling
func (self *GithubClient) getRepositoryTags(ctx context.Context, client *github.Client, owner, repo string) ([]*GithubTag, error) {
	opts := &github.ListOptions{
		PerPage: 100,
	}

	var allTags []*GithubTag
	for {
		// Check if context is canceled
		select {
		case <-ctx.Done():
			return allTags, ctx.Err()
		default:
		}

		tags, resp, err := client.Repositories.ListTags(ctx, owner, repo, opts)
		if err != nil {
			return allTags, fmt.Errorf("error listing tags: %v", err)
		}
		allTags = make([]*GithubTag, len(tags))

		for i, tag := range tags {
			allTags[i] = &GithubTag{
				Name: tag.GetName(),
				Ref:  fmt.Sprintf("refs/tags/%s", tag.GetName()),
				SHA:  tag.GetCommit().GetSHA(),
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allTags, nil
}

// VerifyRepositoryAccess with resource cleanup
func (self *GithubClient) VerifyRepositoryAccess(ctx context.Context, installation *ent.GithubInstallation, owner, repo string) (canAccess bool, repoUrl string, err error) {
	if installation == nil || installation.Edges.GithubApp == nil {
		return false, "", fmt.Errorf("invalid installation: missing app edge or nil")
	}

	// Use a short timeout for this simple verification
	timeoutCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Get authenticated client
	authenticatedClient, err := self.GetAuthenticatedClient(timeoutCtx, installation.GithubAppID, installation.ID, installation.Edges.GithubApp.PrivateKey)
	if err != nil {
		return false, "", fmt.Errorf("error getting authenticated client for %s: %v", installation.AccountLogin, err)
	}
	defer authenticatedClient.Client().CloseIdleConnections()

	// See if we can access the repository
	repoResult, resp, err := authenticatedClient.Repositories.Get(ctx, owner, repo)
	if err == nil {
		// Repository found and accessible
		return true, repoResult.GetCloneURL(), nil
	}

	if resp != nil && resp.StatusCode == 404 {
		// Repository either doesn't exist or installation doesn't have access
		return false, "", nil
	}

	log.Errorf("Error verifying repository access: %v", err)
	return false, "", nil
}
