package system_service

import (
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/infrastructure/buildkitd.go"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate system management with internal permissions and kubernetes RBAC
type SystemService struct {
	cfg             *config.Config
	repo            repositories.RepositoriesInterface
	buildkitManager *buildkitd.BuildkitSettingsManager
}

func NewSystemService(cfg *config.Config, repo repositories.RepositoriesInterface, buildkitManager *buildkitd.BuildkitSettingsManager) *SystemService {
	return &SystemService{
		cfg:             cfg,
		repo:            repo,
		buildkitManager: buildkitManager,
	}
}
