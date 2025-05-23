package oauth2server

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// LoginPageData holds the data to be passed to the login template
type LoginPageData struct {
	ClientID     string
	LoginURI     string
	RedirectURI  string
	ResponseType string
	State        string
	Scope        string
	Error        string
	ErrorMessage string
	PageKey      string
}

// HandleLoginPage renders a styled login form
func (self *Oauth2Server) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	responseType := r.URL.Query().Get("response_type")
	state := r.URL.Query().Get("state")
	scope := r.URL.Query().Get("scope")
	errorParam := r.URL.Query().Get("error")
	// Unique key specific to the rendered form
	pageKey := uuid.NewString()
	loginURL, err := utils.JoinURLPaths(self.Cfg.ExternalAPIURL, "auth", "login")
	if err != nil {
		log.Errorf("Error building login URL: %v\n", err)
		http.Error(w, fmt.Sprintf("Error building login URL: %v", err), http.StatusInternalServerError)
		return
	}

	// Create page data
	data := LoginPageData{
		ClientID:     clientID,
		LoginURI:     loginURL,
		RedirectURI:  redirectURI,
		ResponseType: responseType,
		State:        state,
		Scope:        scope,
		Error:        errorParam,
		PageKey:      pageKey,
	}

	// Set error message based on error code
	if errorParam == "invalid_credentials" {
		data.ErrorMessage = "Invalid email or password. Please try again."
	}

	// Load the template from file
	tmpl, err := loadHtmlTemplate("login.html")
	if err != nil {
		log.Error("Error loading template: ", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Store the state in cache to prevent double submission
	self.StringCache.SetWithExpiration(self.Ctx, pageKey, "1", 30*time.Minute)

	// Set content type and render template
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// load template
func loadHtmlTemplate(name string) (*template.Template, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("Failed to get current file")
	}
	filePath := path.Join(path.Dir(thisFile), "template", name)

	// Load and parse the template from the file
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		// try loading from the root
		filePath = path.Join("app", "template", name)
		tmpl, err = template.ParseFiles(filePath)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse template %s: %w", name, err)
		}
	}

	return tmpl, nil
}
