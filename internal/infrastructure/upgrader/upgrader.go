package upgrader

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/pkg/release"
)

// Upgrader handles the upgrade process for the application
type Upgrader struct {
	releaseManager *release.Manager
	currentVersion string
	k8sClient      *k8s.KubeClient
	httpClient     *http.Client
}

// New creates a new upgrader instance
func New(currentVersion string, k8sClient *k8s.KubeClient) *Upgrader {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create unauthenticated GitHub client for public repositories
	githubClient := github.NewClient(httpClient)

	return &Upgrader{
		releaseManager: release.NewManager(NewGitHubClientWrapper(githubClient)),
		currentVersion: currentVersion,
		k8sClient:      k8sClient,
		httpClient:     httpClient,
	}
}

// CheckForUpdates checks if there are any available updates
func (u *Upgrader) CheckForUpdates(ctx context.Context) ([]string, error) {
	updates, err := u.releaseManager.AvailableUpdates(ctx, u.currentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	return updates, nil
}

// GetLatestVersion returns the latest available version
func (u *Upgrader) GetLatestVersion(ctx context.Context) (string, error) {
	version, err := u.releaseManager.GetLatestVersion(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get latest version: %w", err)
	}
	return version, nil
}

// GetUpdatePath returns the ordered list of versions needed to update from current to target
func (u *Upgrader) GetUpdatePath(ctx context.Context, targetVersion string) ([]string, error) {
	path, err := u.releaseManager.GetUpdatePath(ctx, u.currentVersion, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get update path: %w", err)
	}
	return path, nil
}

// UpgradeToVersion upgrades the application to the specified version
func (u *Upgrader) UpgradeToVersion(ctx context.Context, targetVersion string) error {
	// Get the update path
	updatePath, err := u.GetUpdatePath(ctx, targetVersion)
	if err != nil {
		return fmt.Errorf("failed to get update path: %w", err)
	}

	// Apply each version update in sequence
	for _, version := range updatePath {
		if err := u.applyVersionUpdate(ctx, version); err != nil {
			// If an error occurs, attempt to rollback to the previous version
			if rollbackErr := u.rollbackToVersion(ctx, u.currentVersion); rollbackErr != nil {
				return fmt.Errorf("failed to apply version %s and rollback failed: %v (rollback error: %v)", version, err, rollbackErr)
			}
			return fmt.Errorf("failed to apply version %s: %w", version, err)
		}
	}

	return nil
}

// applyVersionUpdate applies a single version update
func (u *Upgrader) applyVersionUpdate(ctx context.Context, version string) error {
	// Update deployment images
	if err := u.k8sClient.UpdateDeploymentImages(ctx, version); err != nil {
		return fmt.Errorf("failed to update deployment images: %w", err)
	}

	// Wait for deployments to be ready with a timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	if err := u.k8sClient.WaitForDeploymentsReady(timeoutCtx); err != nil {
		return fmt.Errorf("failed to wait for deployments to be ready: %w", err)
	}

	return nil
}

// rollbackToVersion rolls back to a specific version
func (u *Upgrader) rollbackToVersion(ctx context.Context, version string) error {
	// Update deployment images to the rollback version
	if err := u.k8sClient.UpdateDeploymentImages(ctx, version); err != nil {
		return fmt.Errorf("failed to rollback deployment images: %w", err)
	}

	// Wait for deployments to be ready with a timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	if err := u.k8sClient.WaitForDeploymentsReady(timeoutCtx); err != nil {
		return fmt.Errorf("failed to wait for rollback deployments to be ready: %w", err)
	}

	return nil
}
