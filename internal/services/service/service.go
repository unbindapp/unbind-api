package service_service

import (
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate service management with internal permissions and kubernetes RBAC
type ServiceService struct {
	cfg          *config.Config
	repo         repositories.RepositoriesInterface
	githubClient *github.GithubClient
	k8s          *k8s.KubeClient
}

func NewServiceService(cfg *config.Config, repo repositories.RepositoriesInterface, githubClient *github.GithubClient, k8s *k8s.KubeClient) *ServiceService {
	return &ServiceService{
		cfg:          cfg,
		repo:         repo,
		githubClient: githubClient,
		k8s:          k8s,
	}
}
