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
func (s *Server) GithubManifestCreate(ctx context.Context, input *GithubCreateManifestInput) (*GithubCreateManifestResponse, error) {
	// Create GitHub app manifest
	manifest := s.GithubClient.CreateAppManifest(input.Body.RedirectURL)

	// Create resp
	resp := &GithubCreateManifestResponse{}
	resp.Body.Manifest = manifest
	resp.Body.PostURL = fmt.Sprintf("%s/settings/apps/new", s.Cfg.GithubURL)
	return resp, nil
}

// Connect the new github app to our instance, via manifest code exchange
type GithubAppConnectInput struct {
	Body struct {
		Code string `json:"code"`
	}
}

type GithubAppConnectResponse struct {
	Body struct {
		Name string `json:"name"`
	}
}

func (s *Server) GithubAppConnect(ctx context.Context, input *GithubAppConnectInput) (*GithubAppConnectResponse, error) {
	// Exchange the code for tokens.
	appConfig, err := s.GithubClient.ManifestCodeConversion(ctx, input.Body.Code)
	if err != nil {
		return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to exchange manifest code: %v", err))
	}

	// Save the app config
	ghApp, err := s.Repository.CreateGithubApp(ctx, appConfig)
	if err != nil {
		log.Error("Error saving github app", "err", err)
		return nil, huma.Error500InternalServerError("Failed to save github app")
	}

	// Return the app name
	resp := &GithubAppConnectResponse{}
	resp.Body.Name = ghApp.Name
	return resp, nil
}

// Redirect user to install the app
type GithubAppInstallInput struct {
	AppID int64 `path:"app_id" validate:"required"`
}

type GithubAppInstallResponse struct {
	Status int
	Url    string `header:"Location"`
	Cookie string `header:"Set-Cookie"`
}

func (s *Server) GithubAppInstall(ctx context.Context, input *GithubAppInstallInput) (*GithubAppInstallResponse, error) {
	// Get the app
	ghApp, err := s.Repository.GetGithubAppByID(ctx, input.AppID)
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

	return &GithubAppInstallResponse{
		Status: http.StatusTemporaryRedirect,
		Url:    redirectURL,
		Cookie: cookie.String(),
	}, nil
}
