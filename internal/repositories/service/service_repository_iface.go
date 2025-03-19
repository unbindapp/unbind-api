// Code generated by ifacemaker; DO NOT EDIT.

package service_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/service"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

// ServiceRepositoryInterface ...
type ServiceRepositoryInterface interface {
	// Create the service config
	CreateConfig(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID, gitBranch *string, port *int, host *string, replicas *int32, autoDeploy *bool, runCommand *string, public *bool, image *string) (*ent.ServiceConfig, error)
	// Create the service
	Create(ctx context.Context, tx repository.TxInterface, name string, displayName string, description string, serviceType service.Type, builder service.Builder, runtime *string, framework *string, environmentID uuid.UUID, gitHubInstallationID *int64, gitRepository *string, kubernetesSecret string) (*ent.Service, error)
	GetByID(ctx context.Context, serviceID uuid.UUID) (*ent.Service, error)
	GetByInstallationIDAndRepoName(ctx context.Context, installationID int64, repoName string) ([]*ent.Service, error)
	GetGithubPrivateKey(ctx context.Context, serviceID uuid.UUID) (string, error)
	CountDomainCollisons(ctx context.Context, tx repository.TxInterface, domain string) (int, error)
}
