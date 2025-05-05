package auth_handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"golang.org/x/oauth2"
)

type OauthLoginResponse struct {
	Status int
	Url    string `header:"Location"`
	Cookie string `header:"Set-Cookie"`
}

// Login handles the OAuth login redirect.
func (self *HandlerGroup) DevLogin(ctx context.Context, _ *server.EmptyInput) (*OauthLoginResponse, error) {
	// Generate a random state value for CSRF protection.
	state := fmt.Sprintf("dev-login-%s", uuid.New().String())

	// Build the OAuth2 authentication URL with the state.
	authURL := self.srv.OauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	// Add dev parameter
	authURL += "&initiating_url=dev-login"

	// Create a cookie that stores the state value.
	cookie := &http.Cookie{
		Name:     "state",
		Value:    state,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   false, // Change to true if using HTTPS in production.
		HttpOnly: true,
	}

	// Instead of calling http.Redirect, we return a response struct with the headers.
	// Huma will set the "Location" and "Set-Cookie" headers accordingly.
	return &OauthLoginResponse{
		Status: http.StatusTemporaryRedirect,
		Url:    authURL,
		Cookie: cookie.String(),
	}, nil
}
