package middleware

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

func (self *Middleware) Authenticate(ctx huma.Context, next func(huma.Context)) {
	authHeader := ctx.Header("Authorization")
	if authHeader == "" {
		_ = huma.WriteErr(self.api, ctx, http.StatusUnauthorized, "Authorization header required")
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		_ = huma.WriteErr(self.api, ctx, http.StatusUnauthorized, "Authorization header must be a Bearer token")
		return
	}

	bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

	// ! TODO - remove tester token check someday
	if bearerToken == self.cfg.AdminTesterToken {
		user, err := self.repository.User().GetByEmail(ctx.Context(), "admin@unbind.app")
		if err != nil {
			log.Errorf("Failed to process user: %v", err)
			_ = huma.WriteErr(self.api, ctx, http.StatusInternalServerError, "Failed to process user")
			return
		}
		ctx = huma.WithValue(ctx, "user", user)
		next(ctx)
		return
	}

	verifier, err := self.getVerifier()
	if err != nil {
		log.Errorf("Failed to get verifier: %v", err)
		_ = huma.WriteErr(self.api, ctx, http.StatusInternalServerError, "Failed to get verifier")
		return
	}
	token, err := verifier.Verify(ctx.Context(), bearerToken)
	if err != nil {
		_ = huma.WriteErr(self.api, ctx, http.StatusUnauthorized, "Invalid token")
		return
	}

	var claims struct {
		Email    string `json:"email"`
		Username string `json:"preferred_username"`
		Subject  string `json:"sub"`
	}
	if err := token.Claims(&claims); err != nil {
		log.Errorf("Failed to parse claims: %v", err)
		_ = huma.WriteErr(self.api, ctx, http.StatusInternalServerError, "Failed to parse claims")
		return
	}

	// Get or create user using Ent
	user, err := self.repository.User().GetByEmail(ctx.Context(), claims.Email)
	if err != nil {
		log.Errorf("Failed to process user: %v", err)
		_ = huma.WriteErr(self.api, ctx, http.StatusInternalServerError, "Failed to process user")
		return
	}

	ctx = huma.WithValue(ctx, "user", user)
	next(ctx)
}
