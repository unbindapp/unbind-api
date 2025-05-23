package middleware

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

type Middleware struct {
	verifier   *oidc.IDTokenVerifier
	repository repositories.RepositoriesInterface
	api        huma.API
	cfg        *config.Config
}

func NewMiddleware(cfg *config.Config, repository repositories.RepositoriesInterface, api huma.API) *Middleware {
	return &Middleware{
		repository: repository,
		api:        api,
		cfg:        cfg,
	}
}

func (self *Middleware) getVerifier() (*oidc.IDTokenVerifier, error) {
	if self.verifier != nil {
		return self.verifier, nil
	}

	provider, err := oidc.NewProvider(context.Background(), self.cfg.DexIssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}
	self.verifier = provider.Verifier(&oidc.Config{ClientID: self.cfg.DexClientID})

	return self.verifier, nil
}
