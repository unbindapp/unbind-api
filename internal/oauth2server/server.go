package oauth2server

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/url"

	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/infrastructure/cache"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

type Oauth2Server struct {
	Ctx         context.Context
	Cfg         *config.Config
	Repository  repositories.RepositoriesInterface
	Srv         *server.Server
	PrivateKey  *rsa.PrivateKey
	Kid         string
	StringCache *cache.ValkeyCache[string]
}

type RedirectType string

const (
	RedirectLogin     RedirectType = "login"
	RedirectAuthorize RedirectType = "authorize"
)

func (self *Oauth2Server) BuildOauthRedirect(redirectType RedirectType, queryParams map[string]string) (string, error) {
	var baseURL string
	var err error
	allowedUrls := []string{"http://localhost:3000", self.Cfg.ExternalUIUrl}

	if redirectType == RedirectLogin {
		initiatingURL := queryParams["initiating_url"]
		if initiatingURL == "" {
			initiatingURL, _ = utils.JoinURLPaths(self.Cfg.ExternalUIUrl, "/sign-in")
		}
		// Verify that initatingURL is in the allowed URLs
		initiatingURLBase, err := url.Parse(initiatingURL)
		if err != nil {
			// Invalid URL, default to safe option
			initiatingURL, _ = utils.JoinURLPaths(self.Cfg.ExternalUIUrl, "/sign-in")
		} else {
			// Extract base URL (scheme + host)
			initiatingURLBaseStr := fmt.Sprintf("%s://%s", initiatingURLBase.Scheme, initiatingURLBase.Host)

			// Check if base URL is in allowed list
			isAllowed := false
			for _, allowedURL := range allowedUrls {
				if initiatingURLBaseStr == allowedURL {
					isAllowed = true
					break
				}
			}

			// If not allowed, default to safe option
			if !isAllowed {
				initiatingURL, _ = utils.JoinURLPaths(self.Cfg.ExternalUIUrl, "/sign-in")
			}
		}

		baseURL = initiatingURL
	} else {
		baseURL, err = utils.JoinURLPaths(self.Cfg.ExternalOauth2URL, string(redirectType))
	}

	if err != nil {
		return "", err
	}

	// Create URL object to properly handle query parameter encoding
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	// Get query values from URL or create new if none exist
	q := u.Query()

	// Add all provided query parameters
	for key, value := range queryParams {
		q.Set(key, value)
	}

	// Update URL with query parameters
	u.RawQuery = q.Encode()

	return u.String(), nil
}
