package oauth2server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/unbindapp/unbind-api/internal/log"
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

	// Validate client_id is present
	if clientID == "" {
		http.Error(w, "Missing client_id parameter", http.StatusBadRequest)
		return
	}

	// Check if user is authenticated
	if userID == "" {
		// User not authenticated, redirect to login
		// Make sure to properly URL encode all parameters, especially the state
		encodedState := url.QueryEscape(state)
		encodedRedirectURI := url.QueryEscape(redirectURI)
		encodedScope := url.QueryEscape(scope)

		loginURL := fmt.Sprintf("/login?client_id=%s&redirect_uri=%s&response_type=%s&state=%s&scope=%s",
			clientID, encodedRedirectURI, responseType, encodedState, encodedScope)
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
