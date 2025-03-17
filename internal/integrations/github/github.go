package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/common/utils"
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

// Get an authenticated client for a GitHub App installation
func (self *GithubClient) GetAuthenticatedClient(ctx context.Context, appID int64, installationID int64, appPrivateKey string) (*github.Client, error) {
	privateKey, err := utils.DecodePrivateKey(appPrivateKey)
	if err != nil {
		return nil, err
	}

	bearerToken, err := utils.GenerateGithubJWT(appID, privateKey)
	if err != nil {
		return nil, err
	}

	// Add token to client
	client := self.client.WithAuthToken(bearerToken)

	token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create installation token: %v", err)
	}

	return self.client.WithAuthToken(token.GetToken()), nil
}
