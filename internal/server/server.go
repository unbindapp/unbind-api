package server

import (
	"context"
	"net/http"

	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/database"
	"github.com/unbindapp/unbind-api/internal/database/repository"
	"github.com/unbindapp/unbind-api/internal/github"
	"github.com/unbindapp/unbind-api/internal/kubeclient"
	"golang.org/x/oauth2"
)

// EmptyInput can be used when no input is needed.
type EmptyInput struct{}

// Server implements generated.ServerInterface
type Server struct {
	KubeClient   *kubeclient.KubeClient
	Cfg          *config.Config
	OauthConfig  *oauth2.Config
	GithubClient *github.GithubClient
	Repository   *repository.Repository
	StringCache  *database.ValkeyCache[string]
	HttpClient   *http.Client
}

// HealthCheck is your /health endpoint
type HealthResponse struct {
	Body struct {
		Status string `json:"status"`
	}
}

func (self *Server) HealthCheck(ctx context.Context, _ *EmptyInput) (*HealthResponse, error) {
	healthResponse := &HealthResponse{}
	healthResponse.Body.Status = "ok"
	return healthResponse, nil
}
