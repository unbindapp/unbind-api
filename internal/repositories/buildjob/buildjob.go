package buildjob_repo

import (
	"github.com/unbindapp/unbind-api/ent"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

// BuildJobRepository handles GitHub-related database operations
//
//go:generate go run -mod=mod github.com/vburenin/ifacemaker -f "*.go" -i BuildJobRepositoryInterface -p buildjob_repo -s BuildJobRepository -o buildjob_repository_iface.go
type BuildJobRepository struct {
	base *repository.BaseRepository
}

// NewBuildJobRepository creates a new GitHub repository
func NewBuildJobRepository(db *ent.Client) *BuildJobRepository {
	return &BuildJobRepository{
		base: &repository.BaseRepository{DB: db},
	}
}
