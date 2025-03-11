package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/log"
)

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
