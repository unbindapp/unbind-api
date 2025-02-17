package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/unbindapp/unbind-api/internal/log"
)

type Config struct {
	// Postgres
	PostgresHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" envDefault:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	PostgresDB       string `env:"POSTGRES_DB" envDefault:"unbind"`
	// Dex (OIDC provider)
	DexIssuerURL         string `env:"DEX_ISSUER_URL"`
	DexIssuerUrlExternal string `env:"DEX_ISSUER_URL_EXTERNAL"`
	DexClientID          string `env:"DEX_CLIENT_ID"`
	DexClientSecret      string `env:"DEX_CLIENT_SECRET"`
	// Kubernetes config, optional - if in cluster it will use the in-cluster config
	KubeConfig string `env:"KUBECONFIG"`
}

// Parse environment variables into a Config struct
func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing environment", "err", err)
	}
	return &cfg
}
