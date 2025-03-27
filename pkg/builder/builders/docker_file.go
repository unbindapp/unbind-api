package builders

import (
	"context"
	"fmt"
	"os"

	a "github.com/railwayapp/railpack/core/app"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/pkg/builder/internal/buildkit"
)

// Builds a Dockerfile from a git repository
func (self *Builder) BuildDockerfile(ctx context.Context, buildSecrets map[string]string) (imageName, repoName string, err error) {
	// Metadata
	repoName, outputImage, cacheKey := self.GenerateBuildMetadata()

	// -- Create github client
	ghClient := github.NewGithubClient(self.config.GithubURL, nil)

	// -- Clone repository
	log.Infof("Cloning ref '%s' from '%s'", self.config.GitRef, self.config.GitRepoURL)
	// Clone repository to infer information
	tmpDir, err := ghClient.CloneRepository(ctx,
		self.config.GithubAppID,
		self.config.GithubInstallationID,
		self.config.GithubAppPrivateKey,
		self.config.GitRepoURL,
		self.config.GitRef,
	)
	if err != nil {
		log.Error("Error cloning repository", "err", err)
		return "", "", err
	}
	defer os.RemoveAll(tmpDir)

	// Use default Dockerfile if not specified
	if self.config.ServiceDockerfilePath == "" {
		self.config.ServiceDockerfilePath = "Dockerfile"
	}

	// Check if Dockerfile exists
	fullDockerfilePath := fmt.Sprintf("%s/%s", tmpDir, self.config.ServiceDockerfilePath)
	if _, err := os.Stat(fullDockerfilePath); os.IsNotExist(err) {
		return "", repoName, fmt.Errorf("dockerfile not found at path: %s", self.config.ServiceDockerfilePath)
	}

	// Make app from source
	app, err := a.NewApp(tmpDir)
	if err != nil {
		return "", repoName, fmt.Errorf("error creating app: %w", err)
	}

	// Build using BuildKit
	err = buildkit.BuildWithBuildkitClient(
		self.config,
		app.Source,
		buildkit.BuildWithBuildkitClientOptions{
			ImageName:      outputImage,
			CacheKey:       cacheKey,
			Secrets:        buildSecrets,
			DockerfilePath: fullDockerfilePath,
		},
	)
	if err != nil {
		return "", repoName, fmt.Errorf("build failed: %v", err)
	}

	log.Infof("Built image %s from Dockerfile: %s", outputImage, self.config.ServiceDockerfilePath)
	return outputImage, repoName, nil
}
