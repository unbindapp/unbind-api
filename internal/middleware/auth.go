package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/unbindapp/unbind-api/internal/log"
)

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := m.verifier.Verify(r.Context(), bearerToken)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		var claims struct {
			Email    string `json:"email"`
			Username string `json:"preferred_username"`
			Subject  string `json:"sub"`
		}
		if err := token.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}

		// Get or create user using Ent
		user, err := m.repository.GetOrCreateUser(r.Context(), claims.Email, claims.Username, claims.Subject)
		if err != nil {
			log.Errorf("Failed to process user: %v", err)
			http.Error(w, "Failed to process user", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
