package oauth2server

import (
	"context"
	"crypto/rsa"

	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/unbindapp/unbind-api/config"
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
