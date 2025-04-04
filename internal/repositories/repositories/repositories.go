package repositories

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	deployment_repo "github.com/unbindapp/unbind-api/internal/repositories/deployment"
	environment_repo "github.com/unbindapp/unbind-api/internal/repositories/environment"
	github_repo "github.com/unbindapp/unbind-api/internal/repositories/github"
	group_repo "github.com/unbindapp/unbind-api/internal/repositories/group"
	oauth_repo "github.com/unbindapp/unbind-api/internal/repositories/oauth"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	project_repo "github.com/unbindapp/unbind-api/internal/repositories/project"
	service_repo "github.com/unbindapp/unbind-api/internal/repositories/service"
	system_repo "github.com/unbindapp/unbind-api/internal/repositories/system"
	team_repo "github.com/unbindapp/unbind-api/internal/repositories/team"
	user_repo "github.com/unbindapp/unbind-api/internal/repositories/user"
)

// Repositories provides access to all repositories
//
//go:generate go run -mod=mod github.com/vburenin/ifacemaker -f "*.go" -i RepositoriesInterface -p repositories -s Repositories -o repositories_iface.go
type Repositories struct {
	db          *ent.Client
	base        *repository.BaseRepository
	github      github_repo.GithubRepositoryInterface
	user        user_repo.UserRepositoryInterface
	oauth       oauth_repo.OauthRepositoryInterface
	group       group_repo.GroupRepositoryInterface
	project     project_repo.ProjectRepositoryInterface
	team        team_repo.TeamRepositoryInterface
	permissions permissions_repo.PermissionsRepositoryInterface
	environment environment_repo.EnvironmentRepositoryInterface
	service     service_repo.ServiceRepositoryInterface
	deployment  deployment_repo.DeploymentRepositoryInterface
	system      system_repo.SystemRepositoryInterface
}

// NewRepositories creates a new Repositories facade
func NewRepositories(db *ent.Client) *Repositories {
	base := repository.NewBaseRepository(db)
	githubRepo := github_repo.NewGithubRepository(db)
	oauthRepo := oauth_repo.NewOauthRepository(db)
	userRepo := user_repo.NewUserRepository(db)
	projectRepo := project_repo.NewProjectRepository(db)
	teamRepo := team_repo.NewTeamRepository(db)
	environmentRepo := environment_repo.NewEnvironmentRepository(db)
	deploymentRepo := deployment_repo.NewDeploymentRepository(db)
	serviceRepo := service_repo.NewServiceRepository(db, deploymentRepo)
	permissionsRepo := permissions_repo.NewPermissionsRepository(db, userRepo, projectRepo, environmentRepo, serviceRepo, teamRepo)
	groupRepo := group_repo.NewGroupRepository(db, permissionsRepo)
	systemRepo := system_repo.NewSystemRepository(db)
	return &Repositories{
		db:          db,
		base:        base,
		github:      githubRepo,
		user:        userRepo,
		oauth:       oauthRepo,
		group:       groupRepo,
		project:     projectRepo,
		team:        teamRepo,
		permissions: permissionsRepo,
		environment: environmentRepo,
		service:     serviceRepo,
		deployment:  deploymentRepo,
		system:      systemRepo,
	}
}

// Ent() returns the ent client
func (r *Repositories) Ent() *ent.Client {
	return r.db
}

// Github returns the GitHub repository
func (r *Repositories) Github() github_repo.GithubRepositoryInterface {
	return r.github
}

// User returns the User repository
func (r *Repositories) User() user_repo.UserRepositoryInterface {
	return r.user
}

// Oauth returns the OAuth repository
func (r *Repositories) Oauth() oauth_repo.OauthRepositoryInterface {
	return r.oauth
}

// Group returns the Group repository
func (r *Repositories) Group() group_repo.GroupRepositoryInterface {
	return r.group
}

// Project returns the Project repository
func (r *Repositories) Project() project_repo.ProjectRepositoryInterface {
	return r.project
}

// Team returns the Team repository
func (r *Repositories) Team() team_repo.TeamRepositoryInterface {
	return r.team
}

// Permissions returns the Permissions repository
func (r *Repositories) Permissions() permissions_repo.PermissionsRepositoryInterface {
	return r.permissions
}

// Environment returns the Environment repository
func (r *Repositories) Environment() environment_repo.EnvironmentRepositoryInterface {
	return r.environment
}

// Service returns the Service repository
func (r *Repositories) Service() service_repo.ServiceRepositoryInterface {
	return r.service
}

// Deployment returns the Deployment repository
func (r *Repositories) Deployment() deployment_repo.DeploymentRepositoryInterface {
	return r.deployment
}

// System returns the System repository
func (r *Repositories) System() system_repo.SystemRepositoryInterface {
	return r.system
}

func (r *Repositories) WithTx(ctx context.Context, fn func(tx repository.TxInterface) error) error {
	return r.base.WithTx(ctx, fn)
}
