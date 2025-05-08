package updater

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	gh "github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/cache"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	github_integration "github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/pkg/release"
	"github.com/valkey-io/valkey-go"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// Updater handles the update process for the application
type Updater struct {
	cfg            *config.Config
	releaseManager *release.Manager
	CurrentVersion string
	k8sClient      *k8s.KubeClient
	httpClient     *http.Client

	// Cache for updates
	valkeyCache *cache.ValkeyCache[*UpdateCacheItem]
}

type UpdateCacheItem struct {
	Updates   []string
	CheckedAt time.Time
}

// New creates a new updater instance
func New(cfg *config.Config, currentVersion string, k8sClient *k8s.KubeClient, valkeyClient valkey.Client) *Updater {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create unauthenticated GitHub client for public repositories
	githubClient := gh.NewClient(httpClient)

	// ! Temporarily hardcoding version for testing
	currentVersion = "v0.0.1"

	// Create string cache
	valkeyCache := cache.NewCache[*UpdateCacheItem](valkeyClient, "unbind-updater")

	return &Updater{
		cfg:            cfg,
		releaseManager: release.NewManager(NewGitHubClientWrapper(githubClient), cfg.ReleaseRepoOverride),
		CurrentVersion: currentVersion,
		k8sClient:      k8sClient,
		httpClient:     httpClient,
		valkeyCache:    valkeyCache,
	}
}

// CheckForUpdates checks if there are any available updates
func (self *Updater) CheckForUpdates(ctx context.Context) ([]string, error) {
	// Check cache first
	cacheItem, err := self.valkeyCache.Get(ctx, "updates")
	if err == nil && cacheItem != nil {
		// Check if  time is older than 10 minutes
		if time.Since(cacheItem.CheckedAt) < 10*time.Minute {
			return cacheItem.Updates, nil
		}
	}

	// Cache expired or empty, fetch new updates
	updates, err := self.releaseManager.AvailableUpdates(ctx, self.CurrentVersion)
	if err != nil {
		log.Errorf("Failed to check for updates, trying to return cache %v", err)

		if cacheItem != nil {
			// Return cached updates if available
			return cacheItem.Updates, nil
		}

		log.Errorf("Failed to check for updates and no cache available: %v", err)
		return []string{}, nil
	}

	// Cache the updates
	cacheItem = &UpdateCacheItem{
		Updates:   updates,
		CheckedAt: time.Now(),
	}
	if err := self.valkeyCache.Set(ctx, "updates", cacheItem); err != nil {
		log.Errorf("Failed to cache updates: %v", err)
	}

	return updates, nil
}

// GetLatestVersion returns the latest available version
func (self *Updater) GetLatestVersion(ctx context.Context) (string, error) {
	version, err := self.releaseManager.GetLatestVersion(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get latest version: %w", err)
	}
	return version, nil
}

// GetUpdatePath returns the ordered list of versions needed to update from current to target
func (self *Updater) GetUpdatePath(ctx context.Context, targetVersion string) ([]string, error) {
	path, err := self.releaseManager.GetUpdatePath(ctx, self.CurrentVersion, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get update path: %w", err)
	}
	return path, nil
}

// UpdateToVersion updates the application to the specified version
func (self *Updater) UpdateToVersion(ctx context.Context, targetVersion string) error {
	// Get the update path
	updatePath, err := self.GetUpdatePath(ctx, targetVersion)
	if err != nil {
		return fmt.Errorf("failed to get update path: %w", err)
	}

	// Apply Kustomize manifests for each version in the path
	for _, version := range updatePath {
		if err := self.applyKustomizeManifests(ctx, version); err != nil {
			// If an error occurs, attempt to rollback to the previous version
			if rollbackErr := self.rollbackToVersion(ctx, self.CurrentVersion); rollbackErr != nil {
				return fmt.Errorf("failed to apply kustomize manifests for version %s and rollback failed: %v (rollback error: %v)", version, err, rollbackErr)
			}
			return fmt.Errorf("failed to apply kustomize manifests for version %s: %w", version, err)
		}
	}

	// Only update deployment images for the final target version
	if err := self.k8sClient.UpdateDeploymentImages(ctx, targetVersion); err != nil {
		// If an error occurs, attempt to rollback to the previous version
		if rollbackErr := self.rollbackToVersion(ctx, self.CurrentVersion); rollbackErr != nil {
			return fmt.Errorf("failed to update deployment images and rollback failed: %v (rollback error: %v)", err, rollbackErr)
		}
		return fmt.Errorf("failed to update deployment images: %w", err)
	}

	return nil
}

// applyKustomizeManifests applies Kustomize manifests for a specific version
func (self *Updater) applyKustomizeManifests(ctx context.Context, version string) error {
	// Get repository info
	owner, repo := self.releaseManager.GetRepositoryInfo()

	// Create a temporary directory for cloning
	tempDir, err := os.MkdirTemp("", "unbind-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Clone the repository using the GitHub integration
	ghClient := github_integration.NewGithubClient("https://github.com", nil)
	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
	tmpDir, err := ghClient.CloneRepository(ctx, 0, 0, "", cloneURL, fmt.Sprintf("refs/tags/%s", version), "")
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Check for version-specific directory
	versionDir := filepath.Join(tmpDir, version)
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		// No version-specific directory found, nothing to do
		return nil
	}

	// Check for kustomization.yaml
	kustomizationPath := filepath.Join(versionDir, "kustomization.yaml")
	if _, err := os.Stat(kustomizationPath); os.IsNotExist(err) {
		// No kustomization.yaml found, nothing to do
		return nil
	}

	// Create a temporary kustomization file with namespace override
	systemNamespace := self.cfg.GetSystemNamespace()
	tempKustomizationPath := filepath.Join(versionDir, "kustomization.yaml.tmp")

	// Read the original kustomization file
	kustomizationContent, err := os.ReadFile(kustomizationPath)
	if err != nil {
		return fmt.Errorf("failed to read kustomization file: %w", err)
	}

	// Create a new kustomization file with namespace override
	newContent := fmt.Sprintf("apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\nnamespace: %s\n\n%s",
		systemNamespace, string(kustomizationContent))

	if err := os.WriteFile(tempKustomizationPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write temporary kustomization file: %w", err)
	}
	defer os.Remove(tempKustomizationPath)

	// Create a filesystem for Kustomize
	fs := filesys.MakeFsOnDisk()

	// Create a Kustomize options
	opts := krusty.MakeDefaultOptions()
	opts.LoadRestrictions = types.LoadRestrictionsNone

	// Create a Kustomize instance
	k := krusty.MakeKustomizer(opts)

	// Build the kustomization
	resMap, err := k.Run(fs, versionDir)
	if err != nil {
		return fmt.Errorf("failed to build kustomization: %w", err)
	}

	// Convert the resources to YAML
	yaml, err := resMap.AsYaml()
	if err != nil {
		return fmt.Errorf("failed to convert resources to YAML: %w", err)
	}

	// Apply the resources using our Kubernetes client
	if err := self.k8sClient.ApplyYAML(ctx, yaml); err != nil {
		return fmt.Errorf("failed to apply resources: %w", err)
	}

	return nil
}

// rollbackToVersion rolls back to a specific version
func (self *Updater) rollbackToVersion(ctx context.Context, version string) error {
	// Update deployment images to the rollback version
	if err := self.k8sClient.UpdateDeploymentImages(ctx, version); err != nil {
		return fmt.Errorf("failed to rollback deployment images: %w", err)
	}

	return nil
}

// CheckDeploymentsReady checks if all deployments are running with the specified version
func (self *Updater) CheckDeploymentsReady(ctx context.Context, version string) (bool, error) {
	return self.k8sClient.CheckDeploymentsReady(ctx, version)
}

// GetNextAvailableVersion returns the next version that can be updated to from the current version
func (self *Updater) GetNextAvailableVersion(ctx context.Context, currentVersion string) (string, error) {
	version, err := self.releaseManager.GetNextAvailableVersion(ctx, currentVersion)
	if err != nil {
		return "", fmt.Errorf("failed to get next available version: %w", err)
	}
	return version, nil
}
