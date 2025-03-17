package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateGithubJWT(t *testing.T) {
	// Generate a test private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "Failed to generate private key for testing")

	appID := int64(12345)

	// Call the function under test
	token, err := GenerateGithubJWT(appID, privateKey)
	assert.NoError(t, err, "GenerateGithubJWT should not return an error")
	assert.NotEmpty(t, token, "Generated token should not be empty")

	// Verify the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return &privateKey.PublicKey, nil
	})

	assert.NoError(t, err, "Token should be parsable")
	assert.True(t, parsedToken.Valid, "Token should be valid")

	// Extract claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok, "Claims should be of type jwt.MapClaims")

	// Verify issuer
	assert.Equal(t, strconv.FormatInt(appID, 10), claims["iss"], "Issuer claim should match app ID")

	// Verify expiration (10 minutes from issuance)
	issuedAt, ok := claims["iat"].(float64)
	assert.True(t, ok, "IssuedAt claim should be present")

	expiresAt, ok := claims["exp"].(float64)
	assert.True(t, ok, "ExpiresAt claim should be present")

	// Verify expiration is 10 minutes (600 seconds) after issuance
	assert.Equal(t, int64(issuedAt)+600, int64(expiresAt), "Token should expire 10 minutes after issuance")
}

func TestGenerateGithubJWT_VerifyAlgorithm(t *testing.T) {
	// Generate a test private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "Failed to generate private key for testing")

	appID := int64(12345)

	// Call the function under test
	token, err := GenerateGithubJWT(appID, privateKey)
	assert.NoError(t, err, "GenerateGithubJWT should not return an error")

	// Verify the token uses RS256 algorithm
	parts := strings.Split(token, ".")
	assert.Equal(t, 3, len(parts), "JWT should have three parts")

	header, err := jwt.DecodeSegment(parts[0])
	assert.NoError(t, err, "Should be able to decode header")
	assert.Contains(t, string(header), "RS256", "Header should specify RS256 algorithm")
}

func TestGenerateGithubJWT_ErrorOnNilKey(t *testing.T) {
	appID := int64(12345)

	// Call with nil key
	token, err := GenerateGithubJWT(appID, nil)
	assert.Error(t, err, "GenerateGithubJWT should return an error with nil private key")
	assert.Empty(t, token, "Token should be empty when error occurs")
}

func TestGenerateGithubJWT_TimeConsistency(t *testing.T) {
	// Generate a test private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "Failed to generate private key for testing")

	appID := int64(12345)

	// Get current time for reference
	beforeTokenTime := time.Now().Unix()

	// Call the function under test
	token, err := GenerateGithubJWT(appID, privateKey)
	assert.NoError(t, err, "GenerateGithubJWT should not return an error")

	// Get time after token generation
	afterTokenTime := time.Now().Unix()

	// Parse token to extract issuedAt claim
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})
	assert.NoError(t, err, "Token should be parsable")

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok, "Claims should be of type jwt.MapClaims")

	issuedAt, ok := claims["iat"].(float64)
	assert.True(t, ok, "IssuedAt claim should be present")

	// Verify that issuedAt is within the time window when the function was called
	assert.GreaterOrEqual(t, int64(issuedAt), beforeTokenTime, "IssuedAt should not be before test started")
	assert.LessOrEqual(t, int64(issuedAt), afterTokenTime, "IssuedAt should not be after test ended")
}
