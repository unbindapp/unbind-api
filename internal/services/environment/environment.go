package environment_service

import (
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate environment management with internal permissions and kubernetes RBAC
type EnvironmentService struct {
	repo repositories.RepositoriesInterface
	k8s  *k8s.KubeClient
}

func NewEnvironmentService(repo repositories.RepositoriesInterface, k8sClient *k8s.KubeClient) *EnvironmentService {
	return &EnvironmentService{
		repo: repo,
		k8s:  k8sClient,
	}
}
