package builds_service

import (
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate builds management with internal permissions and kubernetes RBAC
type BuildsService struct {
	repo repositories.RepositoriesInterface
}

func NewBuildsService(repo repositories.RepositoriesInterface) *BuildsService {
	return &BuildsService{
		repo: repo,
	}
}
