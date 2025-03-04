package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/internal/utils"
)

// GitHubAppManifest represents the structure for GitHub App manifest
type GitHubAppManifest struct {
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	URL                string             `json:"url"`
	HookAttributes     HookAttributes     `json:"hook_attributes"`
	RedirectURL        string             `json:"redirect_url"`
	Public             bool               `json:"public"`
	DefaultPermissions DefaultPermissions `json:"default_permissions"`
	DefaultEvents      []string           `json:"default_events"`
}

// HookAttributes contains webhook configuration
type HookAttributes struct {
	URL string `json:"url"`
}

// DefaultPermissions contains permission settings
type DefaultPermissions struct {
	Contents string `json:"contents"`
	Issues   string `json:"issues"`
	Metadata string `json:"metadata"`
}

// CreateAppManifest generates the GitHub App manifest
func (g *GithubClient) CreateAppManifest(redirectUrl string) (*GitHubAppManifest, error) {
	// Generate a random suffix
	suffixRand, err := utils.GenerateRandomSimpleID(5)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random suffix: %w", err)
	}

	return &GitHubAppManifest{
		Name:        fmt.Sprintf("unbind-%s-%s", g.cfg.UnbindSuffix, suffixRand),
		Description: "Application to connect unbind with Github",
		URL:         g.cfg.ExternalURL,
		HookAttributes: HookAttributes{
			URL: g.cfg.GithubWebhookURL,
		},
		RedirectURL: redirectUrl,
		Public:      false,
		DefaultPermissions: DefaultPermissions{
			Contents: "read",
			Issues:   "write",
			Metadata: "read",
		},
		DefaultEvents: []string{"push"},
	}, nil
}

// ManifestCodeConversion gets app configruation from github using the code
func (g *GithubClient) ManifestCodeConversion(ctx context.Context, code string) (*github.AppConfig, error) {
	appConfig, response, err := g.client.Apps.CompleteAppManifest(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange manifest code: %w", err)
	}

	// Check for successful status code (201 Created)
	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return appConfig, nil
}
