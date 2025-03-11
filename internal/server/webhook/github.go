package webhook_handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/go-github/v69/github"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/models"
	"github.com/valkey-io/valkey-go"
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

	// Save the app config
	ghApp, err := self.srv.Repository.CreateGithubApp(ctx, appConfig, userIDParsed)
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
	installationURL := fmt.Sprintf(
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
	ghApps, err := self.srv.Repository.GetGithubApps(ctx, false)
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
			permissions := models.GithubInstallationPermissions{}
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
			_, err = self.srv.Repository.UpsertGithubInstallation(
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
			_, err := self.srv.Repository.SetInstallationActive(ctx, installationID, false)
			if err != nil {
				log.Error("Error setting installation as inactive", "err", err)
				return nil, huma.Error500InternalServerError("Failed to set installation as inactive")
			}

		case "suspended":
			_, err := self.srv.Repository.SetInstallationSuspended(ctx, installationID, true)
			if err != nil {
				log.Error("Error setting installation as suspended", "err", err)
				return nil, huma.Error500InternalServerError("Failed to set installation as suspended")
			}

		case "unsuspended":
			_, err := self.srv.Repository.SetInstallationSuspended(ctx, installationID, false)
			if err != nil {
				log.Error("Error setting installation as unsuspended", "err", err)
				return nil, huma.Error500InternalServerError("Failed to set installation as unsuspended")
			}
		}
	}

	return &GithubWebhookOutput{}, nil
}
