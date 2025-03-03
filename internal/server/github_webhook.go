package server

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/models"
)

// Redirect user to install the app
type GithubWebhookInput struct {
	RawBody               []byte `contentType:"application/json"`
	Sha1SignatureHeader   string `header:"X-Hub-Signature"`
	Sha256SignatureHeader string `header:"X-Hub-Signature-256"`
	EventType             string `header:"X-GitHub-Event"`
}

type GithubWebhookOutput struct {
}

func (s *Server) HandleGithubWebhook(ctx context.Context, input *GithubWebhookInput) (*GithubWebhookOutput, error) {
	// Since we may have multiple apps, we want to validate against every webhook secret to see if it belongs to any of our apps
	ghApps, err := s.Repository.GetGithubApps(ctx)
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
	var payload []byte
	for _, app := range ghApps {
		err := github.ValidateSignature(signature, input.RawBody, []byte(app.WebhookSecret))
		if err == nil {
			ghApp = app
			break
		}
	}
	if err != nil {
		log.Error("Error validating github webhook", "err", err)
		return nil, huma.Error400BadRequest("Failed to validate github webhook")
	}

	// Parse the webhook event.
	event, err := github.ParseWebHook(input.EventType, payload)
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
			_, err = s.Repository.UpsertGithubInstallation(
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
			_, err := s.Repository.SetInstallationActive(ctx, installationID, false)
			if err != nil {
				log.Error("Error setting installation as inactive", "err", err)
				return nil, huma.Error500InternalServerError("Failed to set installation as inactive")
			}

		case "suspended":
			_, err := s.Repository.SetInstallationSuspended(ctx, installationID, true)
			if err != nil {
				log.Error("Error setting installation as suspended", "err", err)
				return nil, huma.Error500InternalServerError("Failed to set installation as suspended")
			}

		case "unsuspended":
			_, err := s.Repository.SetInstallationSuspended(ctx, installationID, false)
			if err != nil {
				log.Error("Error setting installation as unsuspended", "err", err)
				return nil, huma.Error500InternalServerError("Failed to set installation as unsuspended")
			}
		}
	}

	return &GithubWebhookOutput{}, nil
}
