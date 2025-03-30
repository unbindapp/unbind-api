package logs_service

import (
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/infrastructure/loki"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate logs management with internal permissions and kubernetes RBAC
type LogsService struct {
	repo        repositories.RepositoriesInterface
	k8s         *k8s.KubeClient
	lokiQuerier *loki.LokiLogQuerier
}

func NewLogsService(repo repositories.RepositoriesInterface, k8sClient *k8s.KubeClient, lokiQuerier *loki.LokiLogQuerier) *LogsService {
	return &LogsService{
		repo:        repo,
		k8s:         k8sClient,
		lokiQuerier: lokiQuerier,
	}
}
