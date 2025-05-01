package oauth2server

import (
	"fmt"
	"net/http"

	"github.com/unbindapp/unbind-api/internal/common/log"
)

// HandleAuthorize handles the OAuth2 authorization request
func (s *Oauth2Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	// Parse the request to get the client_id
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse request", http.StatusBadRequest)
		return
	}

	// Get required parameters
	clientID := r.Form.Get("client_id")
	redirectURI := r.Form.Get("redirect_uri")
	responseType := r.Form.Get("response_type")
	state := r.Form.Get("state")
	scope := r.Form.Get("scope")
	userID := r.Form.Get("user_id")
	initatingURL := r.Form.Get("initiating_url")
	log.Infof("Initiating URL: %s\n", initatingURL)

	// Validate client_id is present
	if clientID == "" {
		http.Error(w, "Missing client_id parameter", http.StatusBadRequest)
		return
	}

	// Check if user is authenticated
	if userID == "" {
		// User not authenticated, redirect to login
		loginURL, err := s.BuildOauthRedirect(RedirectLogin, map[string]string{
			"client_id":     clientID,
			"redirect_uri":  redirectURI,
			"response_type": responseType,
			"state":         state,
			"scope":         scope,
		})
		if err != nil {
			log.Errorf("Error building login URL: %v\n", err)
			http.Error(w, fmt.Sprintf("Error building login URL: %v", err), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, loginURL, http.StatusFound)
		return
	}

	// Set the user authorization handler to use our user ID
	s.Srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (string, error) {
		return userID, nil
	})

	// Handle the authorization request
	err := s.Srv.HandleAuthorizeRequest(w, r)
	if err != nil {
		log.Errorf("Error handling authorize request: %v\n", err)
		http.Error(w, fmt.Sprintf("Authorization error: %v", err), http.StatusBadRequest)
		return
	}
}
