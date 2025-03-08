package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/github"
	"github.com/unbindapp/unbind-api/internal/log"
)

type GithubCreateManifestInput struct {
	Body struct {
		RedirectURL string `json:"redirect_url"`
	}
}

type GithubCreateManifestResponse struct {
	Body struct {
		PostURL  string                    `json:"post_url"`
		Manifest *github.GitHubAppManifest `json:"manifest"`
	}
}

// Create a manifest that the user can use to create a GitHub app
func (self *Server) HandleGithubManifestCreate(ctx context.Context, input *GithubCreateManifestInput) (*GithubCreateManifestResponse, error) {
	// Create GitHub app manifest
	manifest, _, err := self.GithubClient.CreateAppManifest(input.Body.RedirectURL, input.Body.RedirectURL)

	if err != nil {
		log.Error("Error creating github app manifest", "err", err)
		return nil, huma.Error500InternalServerError("Failed to create github app manifest")
	}

	// Create resp
	resp := &GithubCreateManifestResponse{}
	resp.Body.Manifest = manifest
	resp.Body.PostURL = fmt.Sprintf("%s/settings/apps/new", self.Cfg.GithubURL)
	return resp, nil
}

// Connect the new github app to our instance, via manifest code exchange
type HandleGithubAppConnectInput struct {
	Body struct {
		Code string `json:"code"`
	}
}

type HandleGithubAppConnectResponse struct {
	Body struct {
		Name string `json:"name"`
	}
}

func (self *Server) HandleGithubAppConnect(ctx context.Context, input *HandleGithubAppConnectInput) (*HandleGithubAppConnectResponse, error) {
	// Exchange the code for tokens.
	appConfig, err := self.GithubClient.ManifestCodeConversion(ctx, input.Body.Code)
	if err != nil {
		return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to exchange manifest code: %v", err))
	}

	// Save the app config
	ghApp, err := self.Repository.CreateGithubApp(ctx, appConfig)
	if err != nil {
		log.Error("Error saving github app", "err", err)
		return nil, huma.Error500InternalServerError("Failed to save github app")
	}

	// Return the app name
	resp := &HandleGithubAppConnectResponse{}
	resp.Body.Name = ghApp.Name
	return resp, nil
}

// GET Github apps
type GithubAppListInput struct {
	WithInstallations bool `query:"with_installations"`
}

type GithubAppListResponse struct {
	Body []*ent.GithubApp
}

func (self *Server) HandleListGithubApps(ctx context.Context, input *GithubAppListInput) (*GithubAppListResponse, error) {
	apps, err := self.Repository.GetGithubApps(ctx, input.WithInstallations)
	if err != nil {
		log.Error("Error getting github apps", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github apps")
	}

	resp := &GithubAppListResponse{}
	resp.Body = apps
	return resp, nil
}

// GET Github app installations
type GithubAppInstallationListInput struct {
	AppID int64 `path:"app_id" required:"true"`
}

type GithubAppInstallationListResponse struct {
	Body []*ent.GithubInstallation
}

func (self *Server) HandleListGithubAppInstallations(ctx context.Context, input *GithubAppInstallationListInput) (*GithubAppInstallationListResponse, error) {
	installations, err := self.Repository.GetGithubInstallationsByAppID(ctx, input.AppID)
	if err != nil {
		log.Error("Error getting github installations", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github installations")
	}

	resp := &GithubAppInstallationListResponse{}
	resp.Body = installations
	return resp, nil
}

// Redirect user to install the app
type HandleGithubAppInstallInput struct {
	AppID int64 `path:"app_id" required:"true"`
}

type HandleGithubAppInstallResponse struct {
	Status int
	Url    string `header:"Location"`
	Cookie string `header:"Set-Cookie"`
}

func (self *Server) HandleGithubAppInstall(ctx context.Context, input *HandleGithubAppInstallInput) (*HandleGithubAppInstallResponse, error) {
	// Get the app
	ghApp, err := self.Repository.GetGithubAppByID(ctx, input.AppID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, huma.Error404NotFound("App not found")
		}
		log.Error("Error getting github app", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github app")
	}

	// Create a state parameter to verify the callback
	state := uuid.New().String()

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
	redirectURL := fmt.Sprintf(
		"https://github.com/settings/apps/%s/installations/new?state=%s",
		url.QueryEscape(ghApp.Name),
		url.QueryEscape(state),
	)

	return &HandleGithubAppInstallResponse{
		Status: http.StatusTemporaryRedirect,
		Url:    redirectURL,
		Cookie: cookie.String(),
	}, nil
}
