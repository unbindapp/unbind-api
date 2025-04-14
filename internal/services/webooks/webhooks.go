package webhooks_service

import (
	"net/http"

	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate webhooks management with internal permissions and kubernetes RBAC
type WebhooksService struct {
	repo       repositories.RepositoriesInterface
	httpClient *http.Client
}

func NewWebhooksService(repo repositories.RepositoriesInterface) *WebhooksService {
	return &WebhooksService{
		repo:       repo,
		httpClient: &http.Client{},
	}
}
