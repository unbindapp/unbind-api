package server

import (
	"context"
	"net/http"

	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	"github.com/unbindapp/unbind-api/internal/infrastructure/cache"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	deployments_service "github.com/unbindapp/unbind-api/internal/services/deployments"
	environment_service "github.com/unbindapp/unbind-api/internal/services/environment"
	logs_service "github.com/unbindapp/unbind-api/internal/services/logs"
	project_service "github.com/unbindapp/unbind-api/internal/services/project"
	service_service "github.com/unbindapp/unbind-api/internal/services/service"
	system_service "github.com/unbindapp/unbind-api/internal/services/system"
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
	KubeClient           *k8s.KubeClient
	Cfg                  *config.Config
	OauthConfig          *oauth2.Config
	GithubClient         *github.GithubClient
	Repository           repositories.RepositoriesInterface
	StringCache          *cache.ValkeyCache[string]
	HttpClient           *http.Client
	DeploymentController *deployctl.DeploymentController
	// Services
	TeamService        *team_service.TeamService
	ProjectService     *project_service.ProjectService
	ServiceService     *service_service.ServiceService
	EnvironmentService *environment_service.EnvironmentService
	LogService         *logs_service.LogsService
	DeploymentService  *deployments_service.DeploymentService
	SystemService      *system_service.SystemService
}

func (self *Server) GetUserFromContext(ctx context.Context) (user *ent.User, found bool) {
	user, found = ctx.Value("user").(*ent.User)
	return user, found
}
