package project_service

import (
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate project management with internal permissions and kubernetes RBAC
type ProjectService struct {
	repo repositories.RepositoriesInterface
	k8s  *k8s.KubeClient
}

func NewProjectService(repo repositories.RepositoriesInterface, k8sClient *k8s.KubeClient) *ProjectService {
	return &ProjectService{
		repo: repo,
		k8s:  k8sClient,
	}
}
