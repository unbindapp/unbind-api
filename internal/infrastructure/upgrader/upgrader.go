package upgrader

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	gh "github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	github_integration "github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/pkg/release"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// Upgrader handles the upgrade process for the application
type Upgrader struct {
	cfg            *config.Config
	releaseManager *release.Manager
	CurrentVersion string
	k8sClient      *k8s.KubeClient
	httpClient     *http.Client
}

// New creates a new upgrader instance
func New(cfg *config.Config, currentVersion string, k8sClient *k8s.KubeClient) *Upgrader {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create unauthenticated GitHub client for public repositories
	githubClient := gh.NewClient(httpClient)

	// ! Temporarily hardcoding version for testing
	currentVersion = "v0.0.1"

	return &Upgrader{
		cfg:            cfg,
		releaseManager: release.NewManager(NewGitHubClientWrapper(githubClient)),
		CurrentVersion: currentVersion,
		k8sClient:      k8sClient,
		httpClient:     httpClient,
	}
}

// CheckForUpdates checks if there are any available updates
func (self *Upgrader) CheckForUpdates(ctx context.Context) ([]string, error) {
	updates, err := self.releaseManager.AvailableUpdates(ctx, self.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	return updates, nil
}

// GetLatestVersion returns the latest available version
func (self *Upgrader) GetLatestVersion(ctx context.Context) (string, error) {
	version, err := self.releaseManager.GetLatestVersion(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get latest version: %w", err)
	}
	return version, nil
}

// GetUpdatePath returns the ordered list of versions needed to update from current to target
func (self *Upgrader) GetUpdatePath(ctx context.Context, targetVersion string) ([]string, error) {
	path, err := self.releaseManager.GetUpdatePath(ctx, self.CurrentVersion, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get update path: %w", err)
	}
	return path, nil
}

// UpgradeToVersion upgrades the application to the specified version
func (self *Upgrader) UpgradeToVersion(ctx context.Context, targetVersion string) error {
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
func (self *Upgrader) applyKustomizeManifests(ctx context.Context, version string) error {
	// Get repository info
	owner, repo := self.releaseManager.GetRepositoryInfo()

	// Create a temporary directory for cloning
	tempDir, err := os.MkdirTemp("", "unbind-upgrade-*")
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
func (self *Upgrader) rollbackToVersion(ctx context.Context, version string) error {
	// Update deployment images to the rollback version
	if err := self.k8sClient.UpdateDeploymentImages(ctx, version); err != nil {
		return fmt.Errorf("failed to rollback deployment images: %w", err)
	}

	return nil
}

// CheckDeploymentsReady checks if all deployments are running with the specified version
func (self *Upgrader) CheckDeploymentsReady(ctx context.Context, version string) (bool, error) {
	return self.k8sClient.CheckDeploymentsReady(ctx, version)
}
