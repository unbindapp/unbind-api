package config

import (
	"net/url"
	"path"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/utils"
)

type Config struct {
	// Root
	ExternalURL       string `env:"EXTERNAL_URL" envDefault:"http://localhost:8089"`
	ExternalOauth2URL string `env:"EXTERNAL_OAUTH2_URL" envDefault:"http://localhost:8090"`
	// Github Specific
	GithubURL        string `env:"GITHUB_URL" envDefault:"https://github.com"` // Override for github enterprise
	GithubWebhookURL string
	// By default we will just use the external URL for the bind and unbind suffixes
	UnbindSuffix string `env:"UNBIND_SUFFIX"`
	// Postgres
	PostgresHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" envDefault:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	PostgresDB       string `env:"POSTGRES_DB" envDefault:"unbind"`
	// Valkey (redis)
	ValkeyURL string `env:"VALKEY_URL" envDefault:"localhost:6379"`
	// Dex (OIDC provider)
	DexIssuerURL         string `env:"DEX_ISSUER_URL"`
	DexIssuerUrlExternal string `env:"DEX_ISSUER_URL_EXTERNAL"`
	DexClientID          string `env:"DEX_CLIENT_ID"`
	DexClientSecret      string `env:"DEX_CLIENT_SECRET"`
	// Kubernetes config, optional - if in cluster it will use the in-cluster config
	KubeConfig string `env:"KUBECONFIG"`
	// ! TODO - remove me some day, for bypassing oauth
	AdminTesterToken string `env:"ADMIN_TESTER_TOKEN"`
}

// Parse environment variables into a Config struct
func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing environment", "err", err)
	}

	// Get suffix if not present
	if cfg.UnbindSuffix == "" {
		suffix, err := utils.ValidateAndExtractDomain(cfg.ExternalURL)
		if err != nil {
			log.Fatal("Error extracting domain from external URL", "err", err)
		}
		cfg.UnbindSuffix = strings.ToLower(suffix)
	}

	// Parse github callback URL
	baseURL, _ := url.Parse(cfg.ExternalURL)
	baseURL.Path = path.Join(baseURL.Path, "webhook/github")
	cfg.GithubWebhookURL = baseURL.String()

	return &cfg
}
