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
}

// Manager handles release management functionality
type Manager struct {
	client GitHubClientInterface
}

// NewManager creates a new release manager
func NewManager(client GitHubClientInterface) *Manager {
	return &Manager{
		client: client,
	}
}

// AvailableUpdates returns a list of available updates from the current version
// The list is ordered from the next version to the latest version
func (m *Manager) AvailableUpdates(ctx context.Context, currentVersion string) ([]string, error) {
	// Ensure current version has v prefix
	if !strings.HasPrefix(currentVersion, "v") {
		currentVersion = "v" + currentVersion
	}

	// Validate current version
	if !semver.IsValid(currentVersion) {
		return make([]string, 0), nil
	}

	// Get all tags from the repository
	tags, _, err := m.client.Repositories().ListTags(ctx, DefaultOwner, DefaultRepo, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	// Filter and sort valid semver tags
	validVersions := make([]string, 0, len(tags))
	for _, tag := range tags {
		version := tag.GetName()
		if semver.IsValid(version) {
			validVersions = append(validVersions, version)
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
func (m *Manager) GetLatestVersion(ctx context.Context) (string, error) {
	updates, err := m.AvailableUpdates(ctx, "v0.0.0")
	if err != nil {
		return "", err
	}
	if len(updates) == 0 {
		return "", fmt.Errorf("no versions found")
	}
	return updates[len(updates)-1], nil
}

// GetUpdatePath returns the ordered list of versions needed to update from current to target
func (m *Manager) GetUpdatePath(ctx context.Context, currentVersion, targetVersion string) ([]string, error) {
	// Ensure versions have v prefix
	if !strings.HasPrefix(currentVersion, "v") {
		currentVersion = "v" + currentVersion
	}
	if !strings.HasPrefix(targetVersion, "v") {
		targetVersion = "v" + targetVersion
	}

	// Validate versions
	if !semver.IsValid(currentVersion) {
		return nil, fmt.Errorf("invalid current version: %s", currentVersion)
	}
	if !semver.IsValid(targetVersion) {
		return nil, fmt.Errorf("invalid target version: %s", targetVersion)
	}

	// Get all available updates
	updates, err := m.AvailableUpdates(ctx, currentVersion)
	if err != nil {
		return nil, err
	}

	// If no updates are available, return empty slice
	if len(updates) == 0 {
		return make([]string, 0), nil
	}

	// Check if target version exists in available versions
	targetExists := false
	for _, version := range updates {
		if version == targetVersion {
			targetExists = true
			break
		}
	}

	// If target version doesn't exist in available versions, return empty slice
	if !targetExists {
		return make([]string, 0), nil
	}

	// Filter versions up to target version
	updatePath := make([]string, 0, len(updates))
	for _, version := range updates {
		if semver.Compare(version, targetVersion) <= 0 {
			updatePath = append(updatePath, version)
		}
	}

	return updatePath, nil
}

// GetRepositoryInfo returns the repository owner and name
func (m *Manager) GetRepositoryInfo() (string, string) {
	return DefaultOwner, DefaultRepo
}
