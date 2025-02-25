package server

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
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
	// ! TODO - for now we only need 1 github app, so lets check uniqueness, in future we may want more for some reason
	ghApp, err := s.Repository.GetGithubApp(ctx)
	if ghApp != nil {
		log.Info("Github app already exists")
		return nil, huma.Error400BadRequest("Github app already exists")
	}
	if err != nil && !ent.IsNotFound(err) {
		log.Error("Error getting github app", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github app")
	}

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
