// Code generated by ifacemaker; DO NOT EDIT.

package service_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
)

// ServiceRepositoryInterface ...
type ServiceRepositoryInterface interface {
	// Create the service
	Create(ctx context.Context, tx repository.TxInterface, name string, displayName string, description string, environmentID uuid.UUID, gitHubInstallationID *int64, gitRepository *string, gitRepositoryOwner *string, kubernetesSecret string) (*ent.Service, error)
	// Create the service config
	CreateConfig(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID, serviceType schema.ServiceType, builder schema.ServiceBuilder, provider *enum.Provider, framework *enum.Framework, gitBranch *string, ports []schema.PortSpec, hosts []schema.HostSpec, replicas *int32, autoDeploy *bool, runCommand *string, public *bool, image *string) (*ent.ServiceConfig, error)
	// Update the service
	Update(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID, displayName *string, description *string) error
	// Update service config
	UpdateConfig(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID, serviceType *schema.ServiceType, builder *schema.ServiceBuilder, gitBranch *string, ports []schema.PortSpec, hosts []schema.HostSpec, replicas *int32, autoDeploy *bool, runCommand *string, public *bool, image *string) error
	Delete(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID) error
	GetByID(ctx context.Context, serviceID uuid.UUID) (*ent.Service, error)
	GetByInstallationIDAndRepoName(ctx context.Context, installationID int64, repoName string) ([]*ent.Service, error)
	GetByEnvironmentID(ctx context.Context, environmentID uuid.UUID) ([]*ent.Service, error)
	GetGithubPrivateKey(ctx context.Context, serviceID uuid.UUID) (string, error)
	CountDomainCollisons(ctx context.Context, tx repository.TxInterface, domain string) (int, error)
	GetDeploymentNamespace(ctx context.Context, serviceID uuid.UUID) (string, error)
	// Summarize services in environment
	SummarizeServices(ctx context.Context, environmentIDs []uuid.UUID) (counts map[uuid.UUID]int, providers map[uuid.UUID][]enum.Provider, frameworks map[uuid.UUID][]enum.Framework, err error)
}
