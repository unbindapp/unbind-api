package variables_service

import (
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate variables management with internal permissions and kubernetes RBAC
type VariablesService struct {
	repo repositories.RepositoriesInterface
	k8s  *k8s.KubeClient
}

func NewVariablesService(repo repositories.RepositoriesInterface, k8s *k8s.KubeClient) *VariablesService {
	return &VariablesService{
		repo: repo,
		k8s:  k8s,
	}
}
