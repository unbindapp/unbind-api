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
