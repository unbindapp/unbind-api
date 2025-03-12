package oauth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strings"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/database/repository"
)

// Access token generator
// CustomJWTAccessGenerate extends the default JWT generator with more uniqueness
type accessTokenGenerator struct {
	keyID      string
	PrivateKey *rsa.PrivateKey
}

// Token generates a new JWT access token with additional uniqueness claims
func (a *accessTokenGenerator) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (access, refresh string, err error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"aud": data.Client.GetID(),
		"exp": now.Add(data.TokenInfo.GetAccessExpiresIn()).Unix(),
		"sub": data.UserID,
		"iat": now.Unix(),
		"jti": uuid.New().String(),
		"sid": uuid.New().String(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = a.keyID

	// Sign the token
	access, err = token.SignedString(a.PrivateKey)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token if needed
	if isGenRefresh {
		refresh = uuid.New().String()
		refresh = strings.ToUpper(strings.ReplaceAll(refresh, "-", ""))
	}

	return access, refresh, nil
}

// The ID token of the oidc flow
func generateIDToken(ctx context.Context, ti oauth2.TokenInfo, repo *repository.Repository, issuer string, privateKey *rsa.PrivateKey, kid string) (string, error) {
	now := time.Now()

	// Gather the data we need
	userID := ti.GetUserID()
	clientID := ti.GetClientID()

	u, err := repo.GetUserByEmail(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to find user: %w", err)
	}

	claims := jwt.MapClaims{
		"iss":            issuer,
		"jti":            uuid.NewString(),
		"sub":            userID,
		"aud":            clientID,
		"iat":            now.Unix(),
		"exp":            now.Add(ACCESS_TOKEN_EXP).Unix(),
		"email":          u.Email,
		"email_verified": true,
		"name":           "John Doe",
	}

	scopes := strings.Split(ti.GetScope(), " ")
	for _, scope := range scopes {
		if scope == "groups" {
			// ! TODO - dynamic groups
			claims["groups"] = []string{"oidc:users"}
		}
	}

	// Create a token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Optionally set KID in header so clients know which key was used:
	if kid != "" {
		token.Header["kid"] = kid
	}

	// Sign with the RSA private key
	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}
