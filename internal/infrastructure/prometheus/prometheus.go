package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/unbindapp/unbind-api/config"
)

type PrometheusClient struct {
	cfg *config.Config
	api v1.API
}

func NewPrometheusClient(cfg *config.Config) (*PrometheusClient, error) {
	client, err := api.NewClient(api.Config{
		Address: cfg.PrometheusEndpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return &PrometheusClient{
		cfg: cfg,
		api: v1.NewAPI(client),
	}, nil
}
