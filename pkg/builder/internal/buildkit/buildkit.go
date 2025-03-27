package buildkit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/session/secrets/secretsprovider"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/moby/buildkit/util/progress/progressui"
	rpBuildkit "github.com/railwayapp/railpack/buildkit"
	"github.com/railwayapp/railpack/core/plan"
	"github.com/tonistiigi/fsutil"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/pkg/builder/config"
)

type BuildWithBuildkitClientOptions struct {
	ImageName   string
	Platform    rpBuildkit.BuildPlatform
	SecretsHash string
	Secrets     map[string]string
	ImportCache string
	ExportCache string
	CacheKey    string
}

func BuildWithBuildkitClient(cfg *config.Config, appDir string, plan *plan.BuildPlan, opts BuildWithBuildkitClientOptions) error {
	ctx := appcontext.Context()

	imageName := opts.ImageName
	if imageName == "" {
		imageName = getImageName(appDir)
	}

	// Prepend registry URL to image name if configured
	var buildCacheRef string
	if cfg.ContainerRegistryHost != "" {
		// Only prepend registry URL if the image name doesn't already have registry information
		// and if it doesn't contain a port number or domain suffix (like '.com')
		if !strings.Contains(imageName, "/") ||
			(!strings.Contains(imageName, ":") && !strings.Contains(strings.Split(imageName, "/")[0], ".")) {

			// Clean registry URL (remove http/https prefix if present)
			registryURL := cfg.ContainerRegistryHost
			registryURL = strings.TrimPrefix(registryURL, "http://")
			registryURL = strings.TrimPrefix(registryURL, "https://")
			registryURL = strings.TrimRight(registryURL, "/")

			// Prepend the registry URL to the image name
			imageName = fmt.Sprintf("%s/%s", registryURL, imageName)

			// Set the cache ref to the registry URL
			buildCacheRef = fmt.Sprintf("%s/%s:buildcache", registryURL, opts.CacheKey)
		}
	}

	c, err := client.New(ctx, cfg.BuildkitHost)
	if err != nil {
		return fmt.Errorf("failed to connect to buildkit: %w", err)
	}
	defer c.Close()

	// Get the buildkit info early so we can ensure we can connect to the buildkit host
	info, err := c.Info(ctx)
	if err != nil {
		return fmt.Errorf("failed to get buildkit info: %w", err)
	}

	buildPlatform := opts.Platform
	if (buildPlatform == rpBuildkit.BuildPlatform{}) {
		buildPlatform = rpBuildkit.DetermineBuildPlatformFromHost()
	}

	llbState, image, err := rpBuildkit.ConvertPlanToLLB(plan, rpBuildkit.ConvertPlanOptions{
		BuildPlatform: buildPlatform,
		SecretsHash:   opts.SecretsHash,
		CacheKey:      opts.CacheKey,
	})
	if err != nil {
		return fmt.Errorf("error converting plan to LLB: %w", err)
	}

	imageBytes, err := json.Marshal(image)
	if err != nil {
		return fmt.Errorf("error marshalling image: %w", err)
	}

	def, err := llbState.Marshal(ctx, llb.LinuxAmd64)
	if err != nil {
		return fmt.Errorf("error marshaling LLB state: %w", err)
	}

	// Setup channel for progress monitoring
	ch := make(chan *client.SolveStatus)
	var pipeW *io.PipeWriter

	progressDone := make(chan bool)
	go func() {
		displayCh := make(chan *client.SolveStatus)
		go func() {
			for s := range ch {
				displayCh <- s
			}
			close(displayCh)
		}()

		display, err := progressui.NewDisplay(os.Stdout, progressui.AutoMode)
		if err != nil {
			log.Error("failed to create progress display", "error", err)
		}

		_, err = display.UpdateFrom(ctx, displayCh)
		if err != nil {
			log.Error("failed to update progress display", "error", err)
		}
		progressDone <- true
	}()

	appFS, err := fsutil.NewFS(appDir)
	if err != nil {
		return fmt.Errorf("error creating FS: %w", err)
	}

	log.Infof("Building image for %s with BuildKit %s", buildPlatform.String(), info.BuildkitVersion.Version)

	secretsMap := make(map[string][]byte)
	for k, v := range opts.Secrets {
		secretsMap[k] = []byte(v)
	}
	secrets := secretsprovider.FromMap(secretsMap)

	// Setting up session attachments for registry auth if needed
	sessionAttachables := []session.Attachable{secrets}

	// Registry authentication setup
	if cfg.ContainerRegistryUser != "" && cfg.ContainerRegistryPassword != "" {
		// Create a new config file
		configFile := configfile.New("")

		// Add the auth entry for the registry
		configFile.AuthConfigs = map[string]types.AuthConfig{
			cfg.ContainerRegistryHost: {
				Username: cfg.ContainerRegistryUser,
				Password: cfg.ContainerRegistryPassword,
				Auth:     "",
			},
		}

		// Create the auth provider configuration
		cfg := authprovider.DockerAuthProviderConfig{
			ConfigFile: configFile,
		}

		sessionAttachables = append(sessionAttachables, authprovider.NewDockerAuthProvider(cfg))
	}

	solveOpts := client.SolveOpt{}

	// Export to registry
	exportAttrs := map[string]string{
		"name":                  imageName,
		"containerimage.config": string(imageBytes),
		"push":                  "true",
		"compression":           "gzip",
		"compression-level":     "3",
		"registry.insecure":     "true",
	}

	solveOpts = client.SolveOpt{
		LocalMounts: map[string]fsutil.FS{
			"context": appFS,
		},
		Session: sessionAttachables,
		Exports: []client.ExportEntry{
			{
				Type:  client.ExporterImage,
				Attrs: exportAttrs,
			},
		},
		CacheImports: []client.CacheOptionsEntry{
			{
				Type: "registry",
				Attrs: map[string]string{
					"ref": buildCacheRef,
				},
			},
		},
		CacheExports: []client.CacheOptionsEntry{
			{
				Type: "registry",
				Attrs: map[string]string{
					"ref":  buildCacheRef,
					"mode": "max",
				},
			},
		},
	}

	// Add cache import if specified
	if opts.ImportCache != "" {
		solveOpts.CacheImports = append(solveOpts.CacheImports, client.CacheOptionsEntry{
			Type:  "gha",
			Attrs: parseKeyValue(opts.ImportCache),
		})
	}

	// Add cache export if specified
	if opts.ExportCache != "" {
		solveOpts.CacheExports = append(solveOpts.CacheExports, client.CacheOptionsEntry{
			Type:  "gha",
			Attrs: parseKeyValue(opts.ExportCache),
		})
	}

	startTime := time.Now()
	_, err = c.Solve(ctx, def, solveOpts, ch)

	// Wait for progress monitoring to complete
	<-progressDone

	if pipeW != nil {
		pipeW.Close()
	}

	if err != nil {
		return fmt.Errorf("failed to solve: %w", err)
	}

	buildDuration := time.Since(startTime)
	log.Infof("Successfully built image in %.2fs", buildDuration.Seconds())

	log.Infof("image name: %s", imageName)
	return nil
}

func getImageName(appDir string) string {
	parts := strings.Split(appDir, string(os.PathSeparator))
	name := parts[len(parts)-1]
	if name == "" {
		name = "railpack-app" // Fallback if path ends in separator
	}
	return name
}

// Helper function to parse key=value strings into a map
func parseKeyValue(s string) map[string]string {
	attrs := make(map[string]string)
	parts := strings.Split(s, ",")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			attrs[kv[0]] = kv[1]
		}
	}
	return attrs
}
