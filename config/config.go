package config

import (
	"net/url"
	"path"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

type ConfigInterface interface {
	GetPostgresHost() string
	GetPostgresPort() int
	GetPostgresUser() string
	GetPostgresPassword() string
	GetPostgresDB() string
	GetKubeConfig() string
	GetBuildImage() string
	GetSystemNamespace() string
	GetKubeProxyURL() string
	GetBuildkitHost() string
}

type Config struct {
	SystemNamespace string `env:"SYSTEM_NAMESPACE,required"`
	// Root
	ExternalUIUrl  string `env:"EXTERNAL_UI_URL" envDefault:"http://localhost:3000"`
	ExternalAPIURL string `env:"EXTERNAL_API_URL" envDefault:"http://localhost:8089"`
	// This is for generating subdomains
	ExternalWildcardBaseURL string `env:"EXTERNAL_WILDCARD_BASE_URL" envDefault:"http://localhost:8089"`
	ExternalOauth2URL       string `env:"EXTERNAL_OAUTH2_URL" envDefault:"http://localhost:8090"`
	// For unbind custom service definitions
	UnbindServiceDefVersion string `env:"UNBIND_SERVICE_DEF_VERSION" envDefault:"v0.1.23"`
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
	// kube-oidc-proxy
	KubeProxyURL string `env:"KUBE_PROXY_URL" envDefault:"https://kube-oidc-proxy.unbind-system.svc.cluster.local:443"`
	// Builder
	BuildImage string `env:"BUILD_IMAGE" envDefault:"unbindapp/unbind-builder:master-14117576548"`
	// Registry specific
	BootstrapContainerRegistryHost     string `env:"BOOTSTRAP_CONTAINER_REGISTRY_HOST"`
	BootstrapContainerRegistryUser     string `env:"BOOTSTRAP_CONTAINER_REGISTRY_USER"`
	BootstrapContainerRegistryPassword string `env:"BOOTSTRAP_CONTAINER_REGISTRY_PASSWORD"`
	// Buildkit
	BuildkitHost string `env:"BUILDKIT_HOST" envDefault:"tcp://buildkitd.unbind-system:1234"`
	// Logging
	LokiEndpoint string `env:"LOKI_ENDPOINT" envDefault:"http://loki-unbind-gateway.unbind-system.svc.cluster.local"`
	// Metrics
	PrometheusEndpoint string `env:"PROMETHEUS_ENDPOINT" envDefault:"http://kube-prometheus-stack-prometheus.monitoring:9090"`
	// Oauth
	DexConnectorSecret string `env:"DEX_CONNECTOR_SECRET" envDefault:"dex-secret"`
	// ! TODO - remove me some day, for bypassing oauth
	AdminTesterToken string `env:"ADMIN_TESTER_TOKEN"`
}

func (self *Config) GetPostgresHost() string {
	return self.PostgresHost
}

func (self *Config) GetPostgresPort() int {
	return self.PostgresPort
}

func (self *Config) GetPostgresUser() string {
	return self.PostgresUser
}

func (self *Config) GetPostgresPassword() string {
	return self.PostgresPassword
}

func (self *Config) GetPostgresDB() string {
	return self.PostgresDB
}

func (self *Config) GetKubeConfig() string {
	return self.KubeConfig
}

func (self *Config) GetBuildImage() string {
	return self.BuildImage
}

func (self *Config) GetSystemNamespace() string {
	return self.SystemNamespace
}

func (self *Config) GetKubeProxyURL() string {
	return self.KubeProxyURL
}

func (self *Config) GetBuildkitHost() string {
	return self.BuildkitHost
}

// Parse environment variables into a Config struct
func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing environment", "err", err)
	}

	// Get suffix if not present
	if cfg.UnbindSuffix == "" {
		suffix, err := utils.ValidateAndExtractDomain(cfg.ExternalAPIURL)
		if err != nil {
			log.Fatal("Error extracting domain from external URL", "err", err)
		}
		cfg.UnbindSuffix = strings.ToLower(suffix)
	}

	// Parse github callback URL
	baseURL, _ := url.Parse(cfg.ExternalAPIURL)
	baseURL.Path = path.Join(baseURL.Path, "webhook/github")
	cfg.GithubWebhookURL = baseURL.String()

	return &cfg
}
