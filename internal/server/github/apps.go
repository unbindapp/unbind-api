package github_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/utils"
)

type GitHubAppCreateInput struct {
	RedirectURL  string `query:"redirect_url" required:"true" doc:"The client URL to redirect to after the installation is finished"`
	Organization string `query:"organization" doc:"The organization to install the app for, if any"`
}

type GithubAppCreateResponse struct {
	ContentType string `header:"Content-Type"`
	Body        []byte
}

// Handler to render GitHub page with form submission
func (self *HandlerGroup) HandleGithubAppCreate(ctx context.Context, input *GitHubAppCreateInput) (*GithubAppCreateResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

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
	redirect, err := utils.JoinURLPaths(self.srv.Cfg.ExternalURL, "/webhook/github/app/save")
	if err != nil {
		log.Error("Error building redirect URL", "err", err)
		return nil, huma.Error500InternalServerError("Failed to build redirect URL")
	}

	// Create GitHub app manifest, if not organization we also want organization read permission
	manifest, appName, err := self.srv.GithubClient.CreateAppManifest(redirect, input.RedirectURL, input.Organization == "")

	if err != nil {
		log.Error("Error creating github app manifest", "err", err)
		return nil, huma.Error500InternalServerError("Failed to create github app manifest")
	}

	// Create a unique state to identify this request
	state := uuid.New().String()
	err = self.srv.StringCache.SetWithExpiration(ctx, appName, state, 30*time.Minute)
	if err != nil {
		log.Error("Error setting state in cache", "err", err)
		return nil, huma.Error500InternalServerError("Failed to set state in cache")
	}
	// Set a user ID in the cache
	err = self.srv.StringCache.SetWithExpiration(ctx, state, user.ID.String(), 30*time.Minute)
	if err != nil {
		log.Error("Error setting user ID in cache", "err", err)
		return nil, huma.Error500InternalServerError("Failed to set user ID in cache")
	}

	q := url.Values{}
	q.Set("state", state)
	githubUrl := fmt.Sprintf("%s/settings/apps/new?%s", self.srv.Cfg.GithubURL, q.Encode())

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

// Redirect user to install the app
type HandleGithubAppInstallInput struct {
	AppID int64 `path:"app_id" required:"true"`
}

type HandleGithubAppInstallResponse struct {
	Status int
	Url    string `header:"Location"`
	Cookie string `header:"Set-Cookie"`
}

func (self *HandlerGroup) HandleGithubAppInstall(ctx context.Context, input *HandleGithubAppInstallInput) (*HandleGithubAppInstallResponse, error) {
	// Get the app
	ghApp, err := self.srv.Repository.GetGithubAppByID(ctx, input.AppID)
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

// GET Github apps
type GithubAppListInput struct {
	WithInstallations bool `query:"with_installations"`
}

type GithubAppListResponse struct {
	Body []*ent.GithubApp
}

func (self *HandlerGroup) HandleListGithubApps(ctx context.Context, input *GithubAppListInput) (*GithubAppListResponse, error) {
	apps, err := self.srv.Repository.GetGithubApps(ctx, input.WithInstallations)
	if err != nil {
		log.Error("Error getting github apps", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get github apps")
	}

	resp := &GithubAppListResponse{}
	resp.Body = apps
	return resp, nil
}
