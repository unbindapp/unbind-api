package server

import (
	"context"
	"net/http"

	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/database"
	"github.com/unbindapp/unbind-api/internal/github"
	"github.com/unbindapp/unbind-api/internal/k8s"
	"github.com/unbindapp/unbind-api/internal/repository/repositories"
	team_service "github.com/unbindapp/unbind-api/internal/services/team"
	"golang.org/x/oauth2"
)

// EmptyInput can be used when no input is needed.
type EmptyInput struct{}

// Server implements generated.ServerInterface
type Server struct {
	KubeClient   *k8s.KubeClient
	Cfg          *config.Config
	OauthConfig  *oauth2.Config
	GithubClient *github.GithubClient
	Repository   repositories.RepositoriesInterface
	StringCache  *database.ValkeyCache[string]
	HttpClient   *http.Client
	// Services
	TeamService *team_service.TeamService
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

func (self *Server) GetUserFromContext(ctx context.Context) (user *ent.User, found bool) {
	user, found = ctx.Value("user").(*ent.User)
	return user, found
}
