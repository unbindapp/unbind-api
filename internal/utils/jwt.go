package utils

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Generate JWT to authenticate with github APIs
func GenerateGithubJWT(appID int64, privateKey *rsa.PrivateKey) (string, error) {
	if privateKey == nil {
		return "", fmt.Errorf("privateKey is required")
	}
	now := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		IssuedAt:  now,
		ExpiresAt: now + (10 * 60), // JWT valid for 10 minutes
		Issuer:    fmt.Sprintf("%d", appID),
	})

	return token.SignedString(privateKey)
}
