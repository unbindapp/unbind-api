package github

import (
	"context"
	"fmt"
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

// Get the token we can use to authenticate with GitHub
func (self *GithubClient) GetInstallationToken(ctx context.Context, appID int64, installationID int64, appPrivateKey string) (string, error) {
	client, err := self.GetAuthenticatedClient(ctx, appID, installationID, appPrivateKey)
	if err != nil {
		return "", err
	}

	token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create installation token: %v", err)
	}

	return token.GetToken(), nil
}

func (self *GithubClient) GetAuthenticatedClient(ctx context.Context, appID int64, installationID int64, appPrivateKey string) (*github.Client, error) {
	// Get the app's installation token
	token, err := self.GetInstallationToken(ctx, appID, installationID, appPrivateKey)
	if err != nil {
		return nil, err
	}

	return self.client.WithAuthToken(token), nil
}
