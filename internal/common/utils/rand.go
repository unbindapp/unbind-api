package utils

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
)

func GenerateRandomSimpleID(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	charsetLength := big.NewInt(int64(len(charset)))

	result := make([]byte, length)

	// Fill the slice with random characters from the charset
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}

		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}

func GenerateSlug(displayName string) (string, error) {
	// Convert to lowercase
	slug := strings.ToLower(displayName)

	// Replace non-alphanumeric characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// If slug is empty after cleaning, use a default
	if slug == "" {
		slug = "untitled"
	}

	// Generate a random short ID
	shortID, err := GenerateRandomSimpleID(12)
	if err != nil {
		return "", err
	}

	// Combine slug and short ID
	return slug + "-" + shortID, nil
}
