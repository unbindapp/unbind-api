package github

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/config"
)

type GithubClient struct {
	cfg    *config.Config
	client *github.Client
}

func NewGithubClient(cfg *config.Config) *GithubClient {
	httpClient := &http.Client{}
	var githubClient *github.Client

	// Check if we're using GitHub Enterprise
	if cfg.GithubURL != "https://github.com" && cfg.GithubURL != "" {
		// For GitHub Enterprise, we need to set the base URL for the API
		baseURL := cfg.GithubURL
		// Make sure the URL has a trailing slash for the github client
		if !strings.HasSuffix(baseURL, "/") {
			baseURL = baseURL + "/"
		}

		apiURL, _ := url.Parse(baseURL + "api/v3/")
		uploadURL, _ := url.Parse(baseURL + "api/uploads/")

		// Create a GitHub client with enterprise URLs
		githubClient, _ = github.NewClient(httpClient).WithEnterpriseURLs(apiURL.String(), uploadURL.String())
	} else {
		githubClient = github.NewClient(httpClient)
	}

	return &GithubClient{
		cfg:    cfg,
		client: githubClient,
	}
}
