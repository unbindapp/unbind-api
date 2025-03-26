package utils

import (
	"errors"
	"net/url"
	"strings"
)

func ExtractRepoName(gitURL string) (string, error) {
	// Parse the URL
	u, err := url.Parse(gitURL)
	if err != nil || !u.IsAbs() || u.Scheme == "" || u.Host == "" {
		return "", errors.New("invalid URL format")
	}

	// Check if path is empty
	if u.Path == "" {
		return "", errors.New("no repository path found in URL")
	}

	// Clean the path and split
	cleanPath := strings.TrimSuffix(u.Path, ".git")
	cleanPath = strings.TrimPrefix(cleanPath, "/")
	parts := strings.Split(cleanPath, "/")

	// Ensure we have at least org/repo format
	if len(parts) < 2 {
		return "", errors.New("invalid repository path format")
	}

	// Get the repo name (last part)
	repoName := parts[len(parts)-1]
	if repoName == "" {
		return "", errors.New("empty repository name")
	}

	return repoName, nil
}
