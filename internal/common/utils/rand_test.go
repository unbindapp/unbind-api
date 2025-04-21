package utils

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomSimpleID(t *testing.T) {
	t.Run("generates string of correct length", func(t *testing.T) {
		length := 5
		result, err := GenerateRandomSimpleID(length)

		require.NoError(t, err)
		assert.Equal(t, length, len(result), "Generated string should be of the requested length")
	})

	t.Run("generates only lowercase alphanumeric characters", func(t *testing.T) {
		result, err := GenerateRandomSimpleID(5)

		require.NoError(t, err)
		matched, err := regexp.MatchString("^[a-z0-9]+$", result)
		require.NoError(t, err, "Regex matching error")
		assert.True(t, matched, "String should only contain lowercase letters and numbers")
	})

	t.Run("generates different strings on consecutive calls", func(t *testing.T) {
		// Generate multiple strings and ensure they're different
		results := make(map[string]bool)

		// Generate 100 strings to have a statistically significant sample
		for i := 0; i < 100; i++ {
			str, err := GenerateRandomSimpleID(5)
			require.NoError(t, err)
			results[str] = true
		}

		// If truly random, we should have close to 100 unique strings
		// We use a lower bound to account for possible collisions
		assert.Greater(t, len(results), 95, "Should generate mostly unique strings")
	})

	t.Run("handles zero length", func(t *testing.T) {
		result, err := GenerateRandomSimpleID(0)

		require.NoError(t, err)
		assert.Equal(t, "", result, "Zero length should return empty string")
	})

	t.Run("handles large length", func(t *testing.T) {
		length := 1000
		result, err := GenerateRandomSimpleID(length)

		require.NoError(t, err)
		assert.Equal(t, length, len(result), "Should handle large requested lengths")
	})
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		wantPattern string
	}{
		{
			name:        "Simple name",
			displayName: "New Project",
			wantPattern: "^new-project-[a-z0-9]{6}$",
		},
		{
			name:        "Name with special characters",
			displayName: "Hello World! @#$ 123",
			wantPattern: "^hello-world-123-[a-z0-9]{6}$",
		},
		{
			name:        "Name with multiple spaces",
			displayName: "  Multiple   Spaces  ",
			wantPattern: "^multiple-spaces-[a-z0-9]{6}$",
		},
		{
			name:        "Empty string",
			displayName: "",
			wantPattern: "^untitled-[a-z0-9]{6}$",
		},
		{
			name:        "Non-ASCII characters",
			displayName: "Café Résumé",
			wantPattern: "^caf-r-sum-[a-z0-9]{6}$",
		},
		{
			name:        "Only special characters",
			displayName: "!@#$%^&*()",
			wantPattern: "^untitled-[a-z0-9]{6}$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSlug(tt.displayName)

			// Assert no error occurred
			assert.NoError(t, err)

			// Assert the slug matches the expected pattern
			assert.Regexp(t, regexp.MustCompile(tt.wantPattern), got)
		})
	}
}

func TestGenerateSlugUniqueness(t *testing.T) {
	// Generate multiple slugs from the same name
	displayName := "Test Project"
	numSlugs := 10
	slugs := make([]string, numSlugs)

	for i := 0; i < numSlugs; i++ {
		slug, err := GenerateSlug(displayName)
		assert.NoError(t, err)
		slugs[i] = slug
	}

	// Check that all generated slugs are unique
	for i := 0; i < numSlugs; i++ {
		for j := i + 1; j < numSlugs; j++ {
			assert.NotEqual(t, slugs[i], slugs[j], "Generated slugs should be unique")
		}
	}
}
