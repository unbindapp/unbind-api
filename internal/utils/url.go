package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ValidateAndExtractDomain takes a URL string and:
// 1. Validates if it's a proper URL (must start with http:// or https://)
// 2. Extracts the domain name
// 3. Converts special characters to hyphens
// Returns the transformed domain and an error if invalid
func ValidateAndExtractDomain(urlStr string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// Check if scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("invalid URL scheme: %s", parsedURL.Scheme)
	}

	// Get hostname and remove port if present
	hostname := parsedURL.Hostname()
	if hostname == "" {
		return "", fmt.Errorf("missing hostname in URL")
	}

	// Handle any potential URL encoding in the hostname
	hostname, _ = url.QueryUnescape(hostname)

	// Convert special characters to hyphens and remove trailing/leading hyphens
	transformed := transformDomain(hostname)

	return transformed, nil
}

// transformDomain converts special characters to hyphens
func transformDomain(domain string) string {
	// Replace dots and other special characters with hyphens
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	transformed := reg.ReplaceAllString(domain, "-")

	// Remove leading and trailing hyphens
	transformed = strings.Trim(transformed, "-")

	return transformed
}
