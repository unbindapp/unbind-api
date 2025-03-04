package middleware

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/log"
)

func (m *Middleware) Authenticate(ctx huma.Context, next func(huma.Context)) {
	authHeader := ctx.Header("Authorization")
	if authHeader == "" {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "Authorization header required")
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "Authorization header must be a Bearer token")
		return
	}

	bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := m.verifier.Verify(ctx.Context(), bearerToken)
	if err != nil {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized, "Invalid token")
		return
	}

	var claims struct {
		Email    string `json:"email"`
		Username string `json:"preferred_username"`
		Subject  string `json:"sub"`
	}
	if err := token.Claims(&claims); err != nil {
		log.Errorf("Failed to parse claims: %v", err)
		huma.WriteErr(m.api, ctx, http.StatusInternalServerError, "Failed to parse claims")
		return
	}

	// Get or create user using Ent
	user, err := m.repository.GetOrCreateUser(ctx.Context(), claims.Email, claims.Username, claims.Subject)
	if err != nil {
		log.Errorf("Failed to process user: %v", err)
		huma.WriteErr(m.api, ctx, http.StatusInternalServerError, "Failed to process user")
		return
	}

	ctx = huma.WithValue(ctx, "user", user)
	next(ctx)
}
