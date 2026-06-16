package middleware

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/auth"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

func (self *Middleware) Authenticate(ctx huma.Context, next func(huma.Context)) {
	token, ok := extractToken(ctx)
	if !ok {
		_ = huma.WriteErr(self.api, ctx, http.StatusUnauthorized, "Authentication required")
		return
	}

	claims, err := self.tokenManager.Verify(token)
	if err != nil {
		_ = huma.WriteErr(self.api, ctx, http.StatusUnauthorized, "Invalid token")
		return
	}

	user, err := self.repository.User().GetByEmail(ctx.Context(), claims.Email)
	if err != nil {
		log.Errorf("Failed to process user: %v", err)
		_ = huma.WriteErr(self.api, ctx, http.StatusInternalServerError, "Failed to process user")
		return
	}

	ctx = huma.WithValue(ctx, "user", user)
	ctx = huma.WithValue(ctx, "bearer_token", token)
	next(ctx)
}

func extractToken(ctx huma.Context) (string, bool) {
	if authHeader := ctx.Header("Authorization"); authHeader != "" {
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return "", false
		}
		return strings.TrimPrefix(authHeader, "Bearer "), true
	}

	cookie, err := huma.ReadCookie(ctx, auth.AccessTokenCookie)
	if err != nil || cookie.Value == "" {
		return "", false
	}
	return cookie.Value, true
}
