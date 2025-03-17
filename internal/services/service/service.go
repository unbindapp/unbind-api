package service_service

import (
	"github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate service management with internal permissions and kubernetes RBAC
type ServiceService struct {
	repo         repositories.RepositoriesInterface
	githubClient *github.GithubClient
}

func NewServiceService(repo repositories.RepositoriesInterface, githubClient *github.GithubClient) *ServiceService {
	return &ServiceService{
		repo:         repo,
		githubClient: githubClient,
	}
}
