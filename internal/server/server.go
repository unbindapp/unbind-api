package server

import (
	"context"

	"github.com/unbindapp/unbind-api/internal/kubeclient"
)

// Server implements generated.ServerInterface
type Server struct {
	KubeClient *kubeclient.KubeClient
}

// HealthCheck is your /health endpoint
type HealthResponse struct {
	Body struct {
		Status string `json:"status"`
	}
}

func (s *Server) HealthCheck(ctx context.Context, _ *struct{}) (*HealthResponse, error) {
	healthResponse := &HealthResponse{}
	healthResponse.Body.Status = "ok"
	return healthResponse, nil
}
