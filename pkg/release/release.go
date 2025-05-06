package release

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/google/go-github/v69/github"
	"golang.org/x/mod/semver"
)

const (
	// DefaultOwner is the GitHub organization name for Unbind
	DefaultOwner = "unbindapp"
	// DefaultRepo is the GitHub repository name for Unbind releases
	DefaultRepo = "unbind-releases"
)

// GitHubClientInterface defines the interface for GitHub client operations we need
type GitHubClientInterface interface {
	Repositories() RepositoriesServiceInterface
}

// RepositoriesServiceInterface defines the interface for GitHub repositories service
type RepositoriesServiceInterface interface {
	ListTags(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error)
	ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
}

// Manager handles release management functionality
type Manager struct {
	client      GitHubClientInterface
	repo        string
	metadataURL string
}

// NewManager creates a new release manager
func NewManager(client GitHubClientInterface, releaseRepoOverride string) *Manager {
	repo := DefaultRepo
	if releaseRepoOverride != "" {
		repo = releaseRepoOverride
	}

	return &Manager{
		client:      client,
		repo:        repo,
		metadataURL: fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/master/metadata.json", DefaultOwner, repo),
	}
}

// getPublishedReleases returns a map of tag names to their release status
func (self *Manager) getPublishedReleases(ctx context.Context) (map[string]bool, error) {
	releases, _, err := self.client.Repositories().ListReleases(ctx, DefaultOwner, self.repo, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	published := make(map[string]bool)
	for _, release := range releases {
		if release.GetTagName() != "" {
			published[release.GetTagName()] = true
		}
	}

	return published, nil
}

// AvailableUpdates returns a list of available updates from the current version
// The list is ordered from the next version to the latest version
func (self *Manager) AvailableUpdates(ctx context.Context, currentVersion string) ([]string, error) {
	// Ensure current version has v prefix
	if !strings.HasPrefix(currentVersion, "v") {
		currentVersion = "v" + currentVersion
	}

	// Validate current version
	if !semver.IsValid(currentVersion) {
		return make([]string, 0), nil
	}

	// Get all tags from the repository
	tags, _, err := self.client.Repositories().ListTags(ctx, DefaultOwner, self.repo, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	// Get published releases
	published, err := self.getPublishedReleases(ctx)
	if err != nil {
		return nil, err
	}

	// Get version metadata
	metadata, err := self.GetVersionMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get version metadata: %w", err)
	}

	// Filter and sort valid semver tags that have published releases and metadata
	validVersions := make([]string, 0, len(tags))
	for _, tag := range tags {
		version := tag.GetName()
		if semver.IsValid(version) && published[version] {
			// Only include versions that have metadata
			if _, hasMetadata := metadata[version]; hasMetadata {
				validVersions = append(validVersions, version)
			}
		}
	}

	// Sort versions
	sort.Slice(validVersions, func(i, j int) bool {
		return semver.Compare(validVersions[i], validVersions[j]) < 0
	})

	// Find versions that are newer than current version
	updates := make([]string, 0, len(validVersions))
	for _, version := range validVersions {
		if semver.Compare(version, currentVersion) > 0 {
			updates = append(updates, version)
		}
	}

	return updates, nil
}

// GetLatestVersion returns the latest available version
func (self *Manager) GetLatestVersion(ctx context.Context) (string, error) {
	updates, err := self.AvailableUpdates(ctx, "v0.0.0")
	if err != nil {
		return "", err
	}
	if len(updates) == 0 {
		return "", fmt.Errorf("no versions found")
	}
	return updates[len(updates)-1], nil
}

// GetRepositoryInfo returns the repository owner and name
func (self *Manager) GetRepositoryInfo() (string, string) {
	return DefaultOwner, self.repo
}
