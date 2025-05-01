package oauth2server

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/url"
	"slices"

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

func (self *Oauth2Server) BuildOauthRedirect(redirectType RedirectType, queryParams map[string]string, requestUrl *url.URL) (string, error) {
	var baseURL string
	var err error
	allowedUrls := []string{"http://localhost:3000", self.Cfg.ExternalUIUrl}

	if redirectType == RedirectLogin {
		// Check if the request URL host is in the allowed hosts
		host := requestUrl.Hostname()
		protocol := requestUrl.Scheme
		port := requestUrl.Port()

		signInBaseURL := fmt.Sprintf("%s://%s", protocol, host)
		if port != "" {
			signInBaseURL += ":" + port
		}

		isAllowedUrl := slices.Contains(allowedUrls, signInBaseURL)
		if !isAllowedUrl {
			return "", fmt.Errorf("invalid host: %s", host)
		}

		baseURL, err = utils.JoinURLPaths(signInBaseURL, "sign-in")
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
