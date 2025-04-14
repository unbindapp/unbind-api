package webhooks_service

import (
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate webhooks management with internal permissions and kubernetes RBAC
type WebhooksService struct {
	repo repositories.RepositoriesInterface
}

func NewWebhooksService(repo repositories.RepositoriesInterface) *WebhooksService {
	return &WebhooksService{
		repo: repo,
	}
}
