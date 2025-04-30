package storage_service

import (
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate storage management with internal permissions and kubernetes RBAC
type StorageService struct {
	cfg  *config.Config
	repo repositories.RepositoriesInterface
	k8s  *k8s.KubeClient
}

func NewStorageService(cfg *config.Config, repo repositories.RepositoriesInterface, k8sClient *k8s.KubeClient) *StorageService {
	return &StorageService{
		cfg:  cfg,
		repo: repo,
		k8s:  k8sClient,
	}
}
