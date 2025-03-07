package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/utils"
	"github.com/valkey-io/valkey-go"
)

type GitHubAppCreateInput struct {
	RedirectURL string `query:"redirect_url" required:"true" doc:"The client URL to redirect to after the installation is finished"`
}

type GithubAppCreateResponse struct {
	ContentType string `header:"Content-Type"`
	Body        []byte
}

// Handler to render GitHub page with form submission
func (self *Server) HandleGithubAppCreate(ctx context.Context, input *GitHubAppCreateInput) (*GithubAppCreateResponse, error) {
	// Template for the GitHub form submission page
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>GitHub Form Submit</title>
</head>
<body>
    <script>
        const data = {
            post_url: "{{.PostURL}}",
            manifest: {{.ManifestJSON}}
        };
        const form = document.createElement("form");
        form.method = "post";
        form.action = data.post_url;
        form.style.display = "none";
        
        const input = document.createElement("input");
        input.name = "manifest";
        input.type = "text";
        input.value = JSON.stringify(data.manifest);
        
        form.appendChild(input);
        document.body.appendChild(form);
        form.submit();
    </script>
</body>
</html>`

	// Build redirect
	redirect, err := utils.JoinURLPaths(self.Cfg.ExternalURL, "/webhook/github/app/save")
	if err != nil {
		log.Error("Error building redirect URL", "err", err)
		return nil, huma.Error500InternalServerError("Failed to build redirect URL")
	}

	// Create GitHub app manifest
	manifest, appName, err := self.GithubClient.CreateAppManifest(redirect)

	if err != nil {
		log.Error("Error creating github app manifest", "err", err)
		return nil, huma.Error500InternalServerError("Failed to create github app manifest")
	}

	// Create a unique state to identify this request
	state := uuid.New().String()
	err = self.StringCache.SetWithExpiration(ctx, appName, state, 30*time.Minute)
	if err != nil {
		log.Error("Error setting state in cache", "err", err)
		return nil, huma.Error500InternalServerError("Failed to set state in cache")
	}

	// Store redirect URL for this state
	err = self.StringCache.SetWithExpiration(ctx, state, input.RedirectURL, 30*time.Minute)
	if err != nil {
		log.Error("Error setting redirect URL in cache", "err", err)
		return nil, huma.Error500InternalServerError("Failed to set redirect URL in cache")
	}

	q := url.Values{}
	q.Set("state", state)
	githubUrl := fmt.Sprintf("%s/settings/apps/new?%s", self.Cfg.GithubURL, q.Encode())

	// Create template data struct
	type templateData struct {
		PostURL      string
		ManifestJSON template.JS
	}

	// Convert manifest to JSON
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		log.Error("Error marshaling manifest to JSON", "err", err)
		return nil, huma.Error500InternalServerError("Failed to prepare manifest data")
	}

	// Create data for template
	data := templateData{
		PostURL:      githubUrl,
		ManifestJSON: template.JS(string(manifestJSON)),
	}

	// Parse and execute the template
	t, err := template.New("github-form").Parse(tmpl)
	if err != nil {
		log.Error("Error parsing template", "err", err)
		return nil, huma.Error500InternalServerError("Failed to parse HTML template")
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		log.Error("Error executing template", "err", err)
		return nil, huma.Error500InternalServerError("Failed to render HTML template")
	}

	// Create response
	resp := &GithubAppCreateResponse{
		ContentType: "text/html; charset=utf-8",
		Body:        buf.Bytes(),
	}

	return resp, nil
}

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
func (self *Server) HandleGithubAppSave(ctx context.Context, input *HandleGithubAppSaveInput) (*HandleGithubAppSaveResponse, error) {
	// Exchange the code for tokens.
	appConfig, err := self.GithubClient.ManifestCodeConversion(ctx, input.Code)
	if err != nil {
		return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to exchange manifest code: %v", err))
	}

	// Verify state
	state, err := self.StringCache.Get(ctx, appConfig.GetName())
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

	// Get redirect URL from cache
	redirectURL, err := self.StringCache.Get(ctx, input.State)
	if err != nil {
		if err == valkey.Nil {
			return nil, huma.Error400BadRequest("Invalid state")
		}
		log.Error("Error getting redirect URL from cache", "err", err)
		return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to get redirect URL: %v", err))
	}
	// Erase redirect URL and re-save it with app ID
	err = self.StringCache.Delete(ctx, input.State)
	if err != nil {
		log.Error("Error deleting redirect URL from cache", "err", err)
		return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to delete redirect URL: %v", err))
	}
	err = self.StringCache.SetWithExpiration(ctx, strconv.Itoa(int(appConfig.GetID())), redirectURL, 30*time.Minute)

	// Save the app config
	ghApp, err := self.Repository.CreateGithubApp(ctx, appConfig)
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

	return &HandleGithubAppSaveResponse{
		Status: http.StatusTemporaryRedirect,
		Url:    installationURL,
		Cookie: cookie.String(),
	}, nil
}
