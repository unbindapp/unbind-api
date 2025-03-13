package team_service

import (
	"github.com/unbindapp/unbind-api/internal/k8s"
	"github.com/unbindapp/unbind-api/internal/repository/repositories"
)

// Integrate team management with internal permissions and kubernetes RBAC
type TeamService struct {
	repo      repositories.RepositoriesInterface
	k8sClient *k8s.KubeClient
}

func NewTeamService(repo repositories.RepositoriesInterface, k8sClient *k8s.KubeClient) *TeamService {
	return &TeamService{
		repo:      repo,
		k8sClient: k8sClient,
	}
}
