package permissions_repo

import (
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/repository"
	project_repo "github.com/unbindapp/unbind-api/internal/repository/project"
	team_repo "github.com/unbindapp/unbind-api/internal/repository/team"
	user_repo "github.com/unbindapp/unbind-api/internal/repository/user"
)

// PermissionsRepository handles group and user permissions
//
//go:generate go run -mod=mod github.com/vburenin/ifacemaker -f "*.go" -i PermissionsRepositoryInterface -p permissions_repo -s PermissionsRepository -o permissions_repository_iface.go
type PermissionsRepository struct {
	base        *repository.BaseRepository
	userRepo    user_repo.UserRepositoryInterface
	projectRepo project_repo.ProjectRepositoryInterface
	teamRepo    team_repo.TeamRepositoryInterface
}

// NewPermissionsRepository creates a new GitHub repository
func NewPermissionsRepository(db *ent.Client, userRepo user_repo.UserRepositoryInterface, projectRepo project_repo.ProjectRepositoryInterface, teamRepo team_repo.TeamRepositoryInterface) *PermissionsRepository {
	return &PermissionsRepository{
		base:        &repository.BaseRepository{DB: db},
		userRepo:    userRepo,
		projectRepo: projectRepo,
		teamRepo:    teamRepo,
	}
}
