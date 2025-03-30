package loki

import (
	"fmt"
	"net/http"

	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

type LokiLogQuerier struct {
	cfg      *config.Config
	endpoint string
	client   http.Client
}

func NewLokiLogger(cfg *config.Config) (*LokiLogQuerier, error) {
	// Construct endpoint
	reqURLStr, err := utils.JoinURLPaths(cfg.LokiEndpoint, "loki", "api", "v1", "tail")
	if err != nil {
		return nil, fmt.Errorf("Unable to construct loki query URL: %v", err)
	}

	return &LokiLogQuerier{
		cfg:      cfg,
		endpoint: reqURLStr,
		client: http.Client{
			Timeout: 0, // No timeout for streaming
		},
	}, nil
}
