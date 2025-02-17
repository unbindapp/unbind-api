package middleware

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/database/repository"
)

type Middleware struct {
	verifier   *oidc.IDTokenVerifier
	repository repository.RepositoryInterface
}

func NewMiddleware(cfg *config.Config, repository repository.RepositoryInterface) (*Middleware, error) {
	provider, err := oidc.NewProvider(context.Background(), cfg.DexIssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	return &Middleware{
		verifier:   provider.Verifier(&oidc.Config{ClientID: cfg.DexClientID}),
		repository: repository,
	}, nil
}
