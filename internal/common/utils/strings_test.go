package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpolationMarker(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"Hello ${}", true},            // Contains exactly one instance
		{"Hello \\${}", false},         // Contains only an escaped instance
		{"Hello ${} and ${}", false},   // Contains two instances
		{"No markers here", false},     // Contains no instances
		{"Mix of \\${} and ${}", true}, // Contains one unescaped and one escaped
	}

	for _, tc := range testCases {
		result := ContainsExactlyOneInterpolationMarker(tc.input)
		assert.Equal(t, tc.expected, result, "Expected %t for input %s, but got %t", tc.expected, tc.input, result)
	}
}
