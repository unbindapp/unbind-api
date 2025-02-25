package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple domain",
			input:    "example",
			expected: "example",
		},
		{
			name:     "Domain with dots",
			input:    "example.com",
			expected: "example-com",
		},
		{
			name:     "Domain with special characters",
			input:    "hello.world!@#$%^&*",
			expected: "hello-world",
		},
		{
			name:     "Domain with multiple dots",
			input:    "sub.domain.example.com",
			expected: "sub-domain-example-com",
		},
		{
			name:     "Domain with leading and trailing special chars",
			input:    ".example.com.",
			expected: "example-com",
		},
		{
			name:     "Domain with consecutive special chars",
			input:    "example..com",
			expected: "example-com",
		},
		{
			name:     "Domain with hyphen",
			input:    "my-domain.com",
			expected: "my-domain-com",
		},
		{
			name:     "Domain with underscore",
			input:    "my_domain.com",
			expected: "my-domain-com",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only special characters",
			input:    "....",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformDomain(tt.input)
			assert.Equal(t, tt.expected, result, "They should be equal")
		})
	}
}

func TestValidateAndExtractDomain(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedResult string
		expectError    bool
		errorContains  string
	}{
		{
			name:           "Valid HTTP URL with port",
			input:          "http://localhost:8099",
			expectedResult: "localhost",
			expectError:    false,
		},
		{
			name:           "Valid HTTPS URL",
			input:          "https://unbind.app",
			expectedResult: "unbind-app",
			expectError:    false,
		},
		{
			name:           "Invalid scheme",
			input:          "ht://localhost:8099",
			expectedResult: "",
			expectError:    true,
			errorContains:  "invalid URL scheme",
		},
		{
			name:           "Valid URL with subdomain",
			input:          "https://sub.example.com",
			expectedResult: "sub-example-com",
			expectError:    false,
		},
		{
			name:           "Valid URL with path",
			input:          "https://example.com/path/to/resource",
			expectedResult: "example-com",
			expectError:    false,
		},
		{
			name:           "Valid URL with query params",
			input:          "https://example.com?param=value",
			expectedResult: "example-com",
			expectError:    false,
		},
		{
			name:           "Valid URL with fragment",
			input:          "https://example.com#section",
			expectedResult: "example-com",
			expectError:    false,
		},
		{
			name:           "Valid URL with username and password",
			input:          "https://user:pass@example.com",
			expectedResult: "example-com",
			expectError:    false,
		},
		{
			name:           "Empty string",
			input:          "",
			expectedResult: "",
			expectError:    true,
		},
		{
			name:           "Invalid URL format",
			input:          "not a url",
			expectedResult: "",
			expectError:    true,
		},
		{
			name:           "Missing scheme",
			input:          "example.com",
			expectedResult: "",
			expectError:    true,
			errorContains:  "invalid URL scheme",
		},
		{
			name:           "FTP URL (invalid scheme)",
			input:          "ftp://example.com",
			expectedResult: "",
			expectError:    true,
			errorContains:  "invalid URL scheme",
		},
		{
			name:           "URL with IP address",
			input:          "http://192.168.1.1",
			expectedResult: "192-168-1-1",
			expectError:    false,
		},
		{
			name:           "URL with IPv6 address",
			input:          "http://[2001:db8::1]",
			expectedResult: "2001-db8-1",
			expectError:    false,
		},
		{
			name:           "URL with unusual characters in domain",
			input:          "https://example-site!.com",
			expectedResult: "example-site-com",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateAndExtractDomain(tt.input)

			// Check error expectation
			if tt.expectError {
				assert.Error(t, err, "Expected an error but got none")
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, "Error message should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Did not expect an error")
				assert.Equal(t, tt.expectedResult, result, "Results should match")
			}
		})
	}
}
