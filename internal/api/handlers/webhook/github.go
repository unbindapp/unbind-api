package webhook_handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/go-github/v69/github"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/buildctl"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/valkey-io/valkey-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Connect the new github app to our instance, via manifest code exchange
type HandleGithubAppSaveInput struct {
	Code  string `query:"code" required:"true"`
	State string `query:"state" required:"true"`
}

type HandleGithubAppSaveResponse struct {
	Status int
	Url    string `header:"Location"`
	Cookie string `header:"Set-Cookie"`
}

// Save github app and redirect to installation
func (self *HandlerGroup) HandleGithubAppSave(ctx context.Context, input *HandleGithubAppSaveInput) (*HandleGithubAppSaveResponse, error) {
	// Exchange the code for tokens.
	appConfig, err := self.srv.GithubClient.ManifestCodeConversion(ctx, input.Code)
	if err != nil {
		return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to exchange manifest code: %v", err))
	}

	// Verify state
	state, err := self.srv.StringCache.Getdel(ctx, appConfig.GetName())
	if err != nil {
		if err == valkey.Nil {
			return nil, huma.Error400BadRequest("Invalid state")
		}
		log.Error("Error getting state from cache", "err", err)
		return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to get state: %v", err))
	}

	if state != input.State {
		return nil, huma.Error400BadRequest("Invalid state")
	}

	// Get user id from cache
	userID, err := self.srv.StringCache.Getdel(ctx, input.State)
	if err != nil {
		if err == valkey.Nil {
			return nil, huma.Error400BadRequest("Invalid state")
		}
		log.Error("Error getting user ID from cache", "err", err)
		return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to get user ID: %v", err))
	}
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		log.Error("Error parsing user ID", "err", err)
		return nil, huma.Error500InternalServerError("Failed to determine user ID")
	}

	// Get organization from the cache
	// !  TODO - do we need this? seems like installation URL is the same regardless of org
	_, err = self.srv.StringCache.Getdel(ctx, state+"-org")
	if err != nil && !errors.Is(err, valkey.Nil) {
		log.Error("Error getting organization from the cache", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get organization from cache")
	}

	// Save the app config
	ghApp, err := self.srv.Repository.Github().CreateApp(ctx, appConfig, userIDParsed)
	if err != nil {
		log.Error("Error saving github app", "err", err)
		return nil, huma.Error500InternalServerError("Failed to save github app")
	}

	// create a cookie that stores the state value
	cookie := &http.Cookie{
		Name:     "github_install_state",
		Value:    state,
		Path:     "/",
		MaxAge:   int(3600),
		Secure:   false,
		HttpOnly: true,
	}

	// Redirect URL - this is where GitHub will send users to install your app
	var installationURL string
	installationURL = fmt.Sprintf(
		"https://github.com/apps/%s/installations/new?state=%s",
		url.QueryEscape(ghApp.Name),
		url.QueryEscape(state),
	)

	// Delay the redirect because github will 404 otherwise
	time.Sleep(2 * time.Second)

	return &HandleGithubAppSaveResponse{
		Status: http.StatusTemporaryRedirect,
		Url:    installationURL,
		Cookie: cookie.String(),
	}, nil
}

// HandleGithubWebhook handles incoming GitHub webhooks.
type GithubWebhookInput struct {
	RawBody               []byte
	Sha1SignatureHeader   string `header:"X-Hub-Signature"`
	Sha256SignatureHeader string `header:"X-Hub-Signature-256"`
	EventType             string `header:"X-GitHub-Event"`
}

type GithubWebhookOutput struct {
}

func (self *HandlerGroup) HandleGithubWebhook(ctx context.Context, input *GithubWebhookInput) (*GithubWebhookOutput, error) {
	// Since we may have multiple apps, we want to validate against every webhook secret to see if it belongs to any of our apps
	ghApps, err := self.srv.Repository.Github().GetApps(ctx, false)
	if err != nil {
		log.Error("Error getting github apps", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github apps")
	}
	var ghApp *ent.GithubApp

	// Figure out signature for webhook secret validation
	signature := input.Sha256SignatureHeader
	if signature == "" {
		signature = input.Sha1SignatureHeader
	}

	// Validate the payload using the webhook secret.
	for _, app := range ghApps {
		err := github.ValidateSignature(signature, input.RawBody, []byte(app.WebhookSecret))
		if err == nil {
			ghApp = app
			break
		}
	}

	if ghApp == nil {
		log.Error("Received webhook with invalid signature", "input", input.RawBody)
		return nil, huma.Error400BadRequest("Invalid signature")
	}

	// Parse the webhook event.
	event, err := github.ParseWebHook(input.EventType, input.RawBody)
	if err != nil {
		log.Errorf("Could not parse webhook: %v", err)
		return nil, huma.Error400BadRequest("Failed to parse github webhook")
	}

	switch e := event.(type) {
	case *github.InstallationEvent:
		// Common installation data
		installation := e.GetInstallation()
		installationID := installation.GetID()
		account := installation.GetAccount()

		// Check if this event is for our app
		if installation.GetAppID() != ghApp.ID {
			log.Info("Received installation event for different app", "app", e.Installation.GetAppID(), "expected", ghApp.ID)
			return &GithubWebhookOutput{}, nil
		}

		// Determine account type and details
		accountType := githubinstallation.AccountTypeUser
		if strings.EqualFold(account.GetType(), "Organization") {
			accountType = githubinstallation.AccountTypeOrganization
		}

		// Determine repository selection
		repoSelection := githubinstallation.RepositorySelectionAll
		if strings.EqualFold(installation.GetRepositorySelection(), "selected") {
			repoSelection = githubinstallation.RepositorySelectionSelected
		}

		// Process based on action
		switch e.GetAction() {
		case "created", "new_permissions_accepted":
			// Build permissions map
			permissions := schema.GithubInstallationPermissions{}
			if perms := installation.GetPermissions(); perms != nil {
				permissions.Contents = perms.GetContents()
				permissions.Metadata = perms.GetMetadata()
			}

			// Get events
			events := make([]string, 0)
			for _, event := range installation.Events {
				events = append(events, event)
			}

			// Create or update installation in database
			_, err = self.srv.Repository.Github().UpsertInstallation(
				ctx,
				installationID,
				installation.GetAppID(),
				account.GetID(),
				account.GetLogin(),
				accountType,
				account.GetHTMLURL(),
				repoSelection,
				installation.SuspendedAt != nil,
				true,
				permissions,
				events,
			)

			if err != nil {
				log.Error("Error upserting github installation", "err", err)
				return nil, huma.Error500InternalServerError("Failed to upsert github installation")
			}

		case "deleted":
			// Mark as inactive instead of deleting
			_, err := self.srv.Repository.Github().SetInstallationActive(ctx, installationID, false)
			if err != nil {
				log.Error("Error setting installation as inactive", "err", err)
				return nil, huma.Error500InternalServerError("Failed to set installation as inactive")
			}

		case "suspended":
			_, err := self.srv.Repository.Github().SetInstallationSuspended(ctx, installationID, true)
			if err != nil {
				log.Error("Error setting installation as suspended", "err", err)
				return nil, huma.Error500InternalServerError("Failed to set installation as suspended")
			}

		case "unsuspended":
			_, err := self.srv.Repository.Github().SetInstallationSuspended(ctx, installationID, false)
			if err != nil {
				log.Error("Error setting installation as unsuspended", "err", err)
				return nil, huma.Error500InternalServerError("Failed to set installation as unsuspended")
			}
		}
	case *github.PushEvent:
		// Trigger a build if the push event is for a branch we care about
		if e.Repo == nil || e.Installation == nil {
			log.Errorf("Received push event with missing repo or installation %v", e)
			return &GithubWebhookOutput{}, nil
		}
		repoName := e.Repo.GetName()
		repoUrl := e.Repo.GetCloneURL()
		installationID := e.Installation.GetID()
		appID := e.Installation.GetAppID()
		ref := e.GetRef()

		// Get the installation
		installation, err := self.srv.Repository.Github().GetInstallationByID(ctx, installationID)
		if err != nil {
			if ent.IsNotFound(err) {
				log.Info("Received event for installation not found in DB", "id", installationID)
				return &GithubWebhookOutput{}, nil
			}
			log.Error("Error getting installation", "err", err)
			return nil, huma.Error500InternalServerError("Failed to get installation")
		}

		if installation.GithubAppID != appID {
			log.Info("Received push event for different app", "app", appID, "expected", installation.GithubAppID)
			return &GithubWebhookOutput{}, nil
		}

		// Get the services associated with this installation and repo
		services, err := self.srv.Repository.Service().GetByInstallationIDAndRepoName(ctx, installationID, repoName)
		if err != nil {
			log.Error("Error getting services", "err", err)
			return nil, huma.Error500InternalServerError("Failed to get services")
		}

		servicesToBuild := make([]*ent.Service, 0)
		for _, service := range services {
			config := service.Edges.ServiceConfig
			var refToBuild *string
			if config.GitBranch != nil {
				if !strings.Contains(*config.GitBranch, "refs/heads/") {
					refToBuild = utils.ToPtr("refs/heads/" + *config.GitBranch)
				}
			}

			if refToBuild != nil && *refToBuild == ref {
				servicesToBuild = append(servicesToBuild, service)
			}
		}

		if len(servicesToBuild) == 0 {
			// Nothing to do
			return &GithubWebhookOutput{}, nil
		}

		// Trigger builds for each service
		for _, service := range servicesToBuild {
			if !service.Edges.ServiceConfig.AutoDeploy {
				// Skip services that don't have auto-deploy enabled
				continue
			}

			// Get private key for the service's github app.
			privKey, err := self.srv.Repository.Service().GetGithubPrivateKey(ctx, service.ID)
			if err != nil {
				log.Error("Error getting github private key", "err", err)
				return nil, huma.Error500InternalServerError("Failed to get github private key")
			}

			// Get deployment namespace
			namespace, err := self.srv.Repository.Service().GetDeploymentNamespace(ctx, service.ID)

			// Get build secrets
			// ! Use our cluster config for this
			kubeConfig, err := rest.InClusterConfig()
			if err != nil {
				log.Fatalf("Error getting in-cluster config: %v", err)
			}
			client, err := kubernetes.NewForConfig(kubeConfig)
			if err != nil {
				log.Fatalf("Error creating clientset: %v", err)
			}

			buildSecrets, err := self.srv.KubeClient.GetSecretMap(ctx, service.KubernetesBuildSecret, namespace, client)
			if err != nil {
				log.Error("Error getting build secrets", "err", err)
				return nil, huma.Error500InternalServerError("Failed to get build secrets")
			}

			// Convert the byte arrays to base64 strings first
			serializableSecrets := make(map[string]string)
			for k, v := range buildSecrets {
				serializableSecrets[k] = base64.StdEncoding.EncodeToString(v)
			}

			// Serialize the map to JSON
			secretsJSON, err := json.Marshal(serializableSecrets)
			if err != nil {
				log.Error("Error marshalling secrets", "err", err)
				return nil, huma.Error500InternalServerError("Failed to marshal secrets")
			}

			// Create environment for build image
			env := map[string]string{
				"GITHUB_INSTALLATION_ID":      strconv.Itoa(int(installationID)),
				"GITHUB_APP_ID":               strconv.Itoa(int(appID)),
				"GITHUB_APP_PRIVATE_KEY":      privKey,
				"GITHUB_REPO_URL":             repoUrl,
				"GIT_REF":                     ref,
				"CONTAINER_REGISTRY_HOST":     self.srv.Cfg.ContainerRegistryHost,
				"CONTAINER_REGISTRY_USER":     self.srv.Cfg.ContainerRegistryUser,
				"CONTAINER_REGISTRY_PASSWORD": self.srv.Cfg.ContainerRegistryPassword,
				"DEPLOYMENT_NAMESPACE":        namespace,
				"SERVICE_PUBLIC":              strconv.FormatBool(service.Edges.ServiceConfig.Public),
				"SERVICE_REPLICAS":            strconv.Itoa(int(service.Edges.ServiceConfig.Replicas)),
				"SERVICE_SECRET_NAME":         service.KubernetesSecret,
				"SERVICE_BUILD_SECRETS":       string(secretsJSON),
			}

			if service.Provider != nil {
				env["SERVICE_PROVIDER"] = string(*service.Provider)
			}

			if service.Framework != nil {
				env["SERVICE_FRAMEWORK"] = string(*service.Framework)
			}

			if service.Edges.ServiceConfig.Port != nil {
				env["SERVICE_PORT"] = strconv.Itoa(*service.Edges.ServiceConfig.Port)
			}

			if service.Edges.ServiceConfig.Host != nil {
				env["SERVICE_HOST"] = *service.Edges.ServiceConfig.Host
			}

			log.Info("Enqueuing build", "repo", repoName, "branch", ref, "serviceID", service.ID, "installationID", installationID, "appID", appID, "repoUrl", repoUrl)
			jobID, err := self.srv.BuildController.EnqueueBuildJob(
				ctx,
				buildctl.BuildJobRequest{
					ServiceID:   service.ID,
					Environment: env,
				},
			)

			if err != nil {
				log.Error("Error enqueuing build job", "err", err)
				return nil, huma.Error500InternalServerError("Failed to enqueue build job")
			}

			log.Info("Enqueued build job", "jobID", jobID)
		}
	}

	return &GithubWebhookOutput{}, nil
}
