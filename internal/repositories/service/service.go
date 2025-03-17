package service_repo

import (
	"github.com/unbindapp/unbind-api/ent"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

// ServiceRepository handles service database operations
//
//go:generate go run -mod=mod github.com/vburenin/ifacemaker -f "*.go" -i ServiceRepositoryInterface -p service_repo -s ServiceRepository -o service_repository_iface.go
type ServiceRepository struct {
	base *repository.BaseRepository
}

// NewServiceRepository creates a new repository
func NewServiceRepository(db *ent.Client) *ServiceRepository {
	return &ServiceRepository{
		base: &repository.BaseRepository{DB: db},
	}
}
