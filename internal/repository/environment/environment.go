package environment_repo

import (
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/repository"
)

// EnvironmentRepository handles Environment-related database operations
//
//go:generate go run -mod=mod github.com/vburenin/ifacemaker -f "*.go" -i EnvironmentRepositoryInterface -p environment_repo -s EnvironmentRepository -o environment_repository_iface.go
type EnvironmentRepository struct {
	base *repository.BaseRepository
}

// NewEnvironmentRepository creates a new GitHub repository
func NewEnvironmentRepository(db *ent.Client) *EnvironmentRepository {
	return &EnvironmentRepository{
		base: &repository.BaseRepository{DB: db},
	}
}
