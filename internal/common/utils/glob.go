package utils

import (
	"path/filepath"
	"strings"
)

// MatchPattern checks if a string matches a glob pattern.
// The pattern can be a simple glob pattern like "v1.*" or "main".
func MatchPattern(pattern, str string) (bool, error) {
	// If the pattern is a simple string without glob characters, do an exact match
	if !strings.ContainsAny(pattern, "*?[]") {
		return pattern == str, nil
	}

	// Convert glob pattern to filepath pattern
	pattern = filepath.Clean(pattern)
	str = filepath.Clean(str)

	// Use filepath.Match for glob pattern matching
	return filepath.Match(pattern, str)
}
