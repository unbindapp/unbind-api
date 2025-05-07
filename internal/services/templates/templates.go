package templates_service

import (
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate templates management with internal permissions and kubernetes RBAC
type TemplatesService struct {
	repo repositories.RepositoriesInterface
	k8s  *k8s.KubeClient
}

func NewTemplatesService(repo repositories.RepositoriesInterface, k8s *k8s.KubeClient) *TemplatesService {
	return &TemplatesService{
		repo: repo,
		k8s:  k8s,
	}
}
