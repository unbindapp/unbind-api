package repositories

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/repository"
	environment_repo "github.com/unbindapp/unbind-api/internal/repository/environment"
	github_repo "github.com/unbindapp/unbind-api/internal/repository/github"
	group_repo "github.com/unbindapp/unbind-api/internal/repository/group"
	oauth_repo "github.com/unbindapp/unbind-api/internal/repository/oauth"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repository/permissions"
	project_repo "github.com/unbindapp/unbind-api/internal/repository/project"
	team_repo "github.com/unbindapp/unbind-api/internal/repository/team"
	user_repo "github.com/unbindapp/unbind-api/internal/repository/user"
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
}

// NewRepositories creates a new Repositories facade
func NewRepositories(db *ent.Client) *Repositories {
	base := repository.NewBaseRepository(db)
	githubRepo := github_repo.NewGithubRepository(db)
	oauthRepo := oauth_repo.NewOauthRepository(db)
	userRepo := user_repo.NewUserRepository(db)
	projectRepo := project_repo.NewProjectRepository(db)
	teamRepo := team_repo.NewTeamRepository(db)
	permissionsRepo := permissions_repo.NewPermissionsRepository(db, userRepo, projectRepo, teamRepo)
	groupRepo := group_repo.NewGroupRepository(db, permissionsRepo)
	environmentRepo := environment_repo.NewEnvironmentRepository(db)
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

func (r *Repositories) WithTx(ctx context.Context, fn func(tx repository.TxInterface) error) error {
	return r.base.WithTx(ctx, fn)
}
