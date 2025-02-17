// internal/auth/middleware.go
package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/user"
)

type AuthMiddleware struct {
	verifier *oidc.IDTokenVerifier
	client   *ent.Client
}

func NewAuthMiddleware(issuerURL string, clientID string, client *ent.Client) (*AuthMiddleware, error) {
	provider, err := oidc.NewProvider(context.Background(), issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	return &AuthMiddleware{
		verifier: provider.Verifier(&oidc.Config{ClientID: clientID}),
		client:   client,
	}, nil
}

func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := a.verifier.Verify(r.Context(), bearerToken)
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
		user, err := a.getOrCreateUser(r.Context(), claims.Email, claims.Username, claims.Subject)
		if err != nil {
			http.Error(w, "Failed to process user", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) getOrCreateUser(ctx context.Context, email, username, subject string) (*ent.User, error) {
	// Try to find existing user
	user, err := a.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)

	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}
		// Create new user if not found
		user, err = a.client.User.
			Create().
			SetEmail(email).
			SetUsername(username).
			SetExternalID(subject).
			Save(ctx)

		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
