package server

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/oauth2"
)

// CallbackInput defines the query parameters for the callback endpoint.
type CallbackInput struct {
	Code string `query:"code" required:"true"`
}

// CallbackResponse defines the JSON structure for the response.
// Huma will automatically encode this as JSON.
type CallbackResponse struct {
	Body struct {
		AccessToken  string    `json:"access_token"`
		TokenType    string    `json:"token_type"`
		RefreshToken string    `json:"refresh_token"`
		Expiry       time.Time `json:"expiry"`
	}
}

// Callback handles the OAuth2 callback.
func (self *Server) Callback(ctx context.Context, in *CallbackInput) (*CallbackResponse, error) {
	// Validate the code parameter.
	if in.Code == "" {
		return nil, huma.Error400BadRequest("No code provided")
	}

	// Exchange the code for tokens.
	oauth2Token, err := self.OauthConfig.Exchange(ctx, in.Code, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to exchange token: %v", err))
	}

	// For development, return the token details as JSON.
	cbResponse := &CallbackResponse{}
	cbResponse.Body.AccessToken = oauth2Token.AccessToken
	cbResponse.Body.TokenType = oauth2Token.TokenType
	cbResponse.Body.RefreshToken = oauth2Token.RefreshToken
	cbResponse.Body.Expiry = oauth2Token.Expiry
	return cbResponse, nil
}
