package auth_handler

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type RefreshInput struct {
	RefreshToken http.Cookie `cookie:"refresh_token"`
}

func (self *HandlerGroup) Refresh(ctx context.Context, input *RefreshInput) (*SessionResponse, error) {
	value := input.RefreshToken.Value
	if value == "" {
		return nil, huma.Error401Unauthorized("Missing refresh token")
	}

	stored, err := self.srv.Repository.Oauth().GetByRefreshToken(ctx, value)
	if err != nil || stored.Revoked || stored.ExpiresAt.Before(time.Now()) {
		return nil, huma.Error401Unauthorized("Invalid refresh token")
	}

	user := stored.Edges.User
	if user == nil {
		return nil, huma.Error401Unauthorized("Invalid refresh token")
	}

	if err := self.srv.Repository.Oauth().RevokeRefreshToken(ctx, value); err != nil {
		return nil, huma.Error500InternalServerError("Failed to rotate session", err)
	}

	return self.issueSession(ctx, user)
}
