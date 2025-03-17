package server

import (
	"context"
	"net/http"

	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/infrastructure/cache"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	project_service "github.com/unbindapp/unbind-api/internal/services/project"
	service_service "github.com/unbindapp/unbind-api/internal/services/service"
	team_service "github.com/unbindapp/unbind-api/internal/services/team"
	"golang.org/x/oauth2"
)

// EmptyInput can be used when no input is needed.
type EmptyInput struct{}

// BaseAuthInput can be used when no input is needed and the user must be authenticated.
type BaseAuthInput struct {
	Authorization string `header:"Authorization" doc:"Bearer token" required:"true"`
}

// Server implements generated.ServerInterface
type Server struct {
	KubeClient   *k8s.KubeClient
	Cfg          *config.Config
	OauthConfig  *oauth2.Config
	GithubClient *github.GithubClient
	Repository   repositories.RepositoriesInterface
	StringCache  *cache.ValkeyCache[string]
	HttpClient   *http.Client
	// Services
	TeamService    *team_service.TeamService
	ProjectService *project_service.ProjectService
	ServiceService *service_service.ServiceService
}

func (self *Server) GetUserFromContext(ctx context.Context) (user *ent.User, found bool) {
	user, found = ctx.Value("user").(*ent.User)
	return user, found
}
