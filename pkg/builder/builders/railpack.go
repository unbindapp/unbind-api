package builders

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/railwayapp/railpack/core"
	a "github.com/railwayapp/railpack/core/app"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/pkg/builder/internal/buildkit"
)

func (self *Builder) BuildWithRailpack(ctx context.Context, buildSecrets map[string]string) (imageName, repoName string, err error) {
	// -- Generate image name
	repoName, err = utils.ExtractRepoName(self.config.GitRepoURL)
	if err != nil {
		log.Warnf("Failed to extract repository name: %v", err)
		repoName = fmt.Sprintf("unbind-build-%d", time.Now().Unix())
	}
	outputImage := fmt.Sprintf("%s/%s:%d", self.config.ContainerRegistryHost, repoName, time.Now().Unix())
	cacheKey := fmt.Sprintf("%s/%s:buildcache", self.config.ContainerRegistryHost, repoName)

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

	// --- Railpack build
	buildResult, app, _, err := GenerateBuildResult(tmpDir, buildSecrets)
	if err != nil {
		return "", repoName, fmt.Errorf("failed to generate build result: %v", err)
	}

	core.PrettyPrintBuildResult(buildResult, core.PrintOptions{Version: "unbind-builder"})

	if !buildResult.Success {
		return "", repoName, fmt.Errorf("build failed")
	}

	err = buildkit.BuildWithBuildkitClient(
		self.config,
		app.Source,
		buildResult.Plan,
		buildkit.BuildWithBuildkitClientOptions{
			ImageName: outputImage,
			CacheKey:  cacheKey,
			Secrets:   buildSecrets,
			// ! TODO - add IMport/Export cache
		},
	)

	if err != nil {
		return "", repoName, fmt.Errorf("build failed: %v", err)
	}

	log.Infof("Built image %s", outputImage)
	return outputImage, repoName, nil
}

func GenerateBuildResult(directory string, buildSecrets map[string]string) (*core.BuildResult, *a.App, *a.Environment, error) {
	app, err := a.NewApp(directory)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating app: %w", err)
	}

	log.Infof("Building %s", app.Source)

	env := a.Environment{
		Variables: buildSecrets,
	}

	generateOptions := &core.GenerateBuildPlanOptions{
		RailpackVersion:          "unbind-builder", // ! Add a version
		ErrorMissingStartCommand: true,
	}

	buildResult := core.GenerateBuildPlan(app, &env, generateOptions)

	return buildResult, app, &env, nil
}
