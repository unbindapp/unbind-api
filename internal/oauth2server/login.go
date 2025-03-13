package oauth2server

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/valkey-io/valkey-go"
)

// LoginPageData holds the data to be passed to the login template
type LoginPageData struct {
	ClientID     string
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

	// Create page data
	data := LoginPageData{
		ClientID:     clientID,
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

// HandleLoginSubmit processes the login form submission
func (self *Oauth2Server) HandleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	clientID := r.Form.Get("client_id")
	redirectURI := r.Form.Get("redirect_uri")
	responseType := r.Form.Get("response_type")
	state := r.Form.Get("state")
	scope := r.Form.Get("scope")
	pageKey := r.Form.Get("page_key")

	// Validate credentials against your repository
	user, err := self.Repository.User().Authenticate(r.Context(), username, password)
	if err != nil {
		// Authentication failed - redirect back to login with error
		// Preserve all original query parameters
		q := url.Values{}
		q.Set("client_id", clientID)
		q.Set("redirect_uri", redirectURI)
		q.Set("response_type", responseType)
		q.Set("state", state)
		q.Set("scope", scope)
		q.Set("error", "invalid_credentials")
		q.Set("page_key", pageKey)

		http.Redirect(w, r, "/login?"+q.Encode(), http.StatusFound)
		return
	}

	// Validate state to prevent double submission
	if _, err := self.StringCache.Getdel(self.Ctx, pageKey); err != nil {
		if err == valkey.Nil {
			// Return ok empty status
			w.WriteHeader(http.StatusNoContent)
			return
		}
		log.Error("Error validating request: ", err)
		http.Error(w, "Error validating request", http.StatusInternalServerError)
		return
	}

	// Authentication succeeded - redirect back to authorize endpoint with user info
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", responseType)
	q.Set("state", state)
	q.Set("scope", scope)
	q.Set("user_id", user.Email) // Pass the authenticated user ID

	http.Redirect(w, r, "/authorize?"+q.Encode(), http.StatusFound)
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
