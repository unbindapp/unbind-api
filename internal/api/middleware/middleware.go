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

func NewMiddleware(cfg *config.Config, repository repositories.RepositoriesInterface, api huma.API) (*Middleware, error) {
	provider, err := oidc.NewProvider(context.Background(), cfg.DexIssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	return &Middleware{
		verifier:   provider.Verifier(&oidc.Config{ClientID: cfg.DexClientID}),
		repository: repository,
		api:        api,
		cfg:        cfg,
	}, nil
}
