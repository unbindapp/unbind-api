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
	// Build base URL by joining paths
	baseUrl, err := utils.JoinURLPaths(self.Cfg.ExternalOauth2URL, string(redirectType))
	if err != nil {
		return "", err
	}

	// Create URL object to properly handle query parameter encoding
	u, err := url.Parse(baseUrl)
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
