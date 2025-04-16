// Code generated by ifacemaker; DO NOT EDIT.

package service_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

// ServiceRepositoryInterface ...
type ServiceRepositoryInterface interface {
	// Create the service
	Create(ctx context.Context, tx repository.TxInterface, name string, displayName string, description string, environmentID uuid.UUID, gitHubInstallationID *int64, gitRepository *string, gitRepositoryOwner *string, kubernetesSecret string) (*ent.Service, error)
	CreateConfig(ctx context.Context, tx repository.TxInterface, input *MutateConfigInput) (*ent.ServiceConfig, error)
	// Update the service
	Update(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID, displayName *string, description *string) error
	// Update service config
	UpdateConfig(ctx context.Context, tx repository.TxInterface, input *MutateConfigInput) error
	Delete(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID) error
	SetCurrentDeployment(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID, deploymentID uuid.UUID) error
	GetByID(ctx context.Context, serviceID uuid.UUID) (svc *ent.Service, err error)
	GetByName(ctx context.Context, name string) (*ent.Service, error)
	GetByInstallationIDAndRepoName(ctx context.Context, installationID int64, repoName string) ([]*ent.Service, error)
	GetByEnvironmentID(ctx context.Context, environmentID uuid.UUID) ([]*ent.Service, error)
	GetGithubPrivateKey(ctx context.Context, serviceID uuid.UUID) (string, error)
	CountDomainCollisons(ctx context.Context, tx repository.TxInterface, domain string) (int, error)
	GetDeploymentNamespace(ctx context.Context, serviceID uuid.UUID) (string, error)
	// Summarize services in environment
	SummarizeServices(ctx context.Context, environmentIDs []uuid.UUID) (counts map[uuid.UUID]int, icons map[uuid.UUID][]string, err error)
	NeedsDeployment(ctx context.Context, service *ent.Service) (NeedsDeploymentResponse, error)
}
