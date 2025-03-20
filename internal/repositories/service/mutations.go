package service_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

// Create the service
func (self *ServiceRepository) Create(
	ctx context.Context,
	tx repository.TxInterface,
	name string,
	displayName string,
	description string,
	serviceType service.Type,
	builder service.Builder,
	runtime *string,
	framework *string,
	environmentID uuid.UUID,
	gitHubInstallationID *int64,
	gitRepository *string,
	kubernetesSecret string,
) (*ent.Service, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.Service.Create().
		SetName(name).
		SetDisplayName(displayName).
		SetDescription(description).
		SetType(serviceType).
		SetBuilder(builder).
		SetNillableRuntime(runtime).
		SetNillableFramework(framework).
		SetEnvironmentID(environmentID).
		SetNillableGithubInstallationID(gitHubInstallationID).
		SetNillableGitRepository(gitRepository).
		SetKubernetesSecret(kubernetesSecret).
		Save(ctx)
}

// Create the service config
func (self *ServiceRepository) CreateConfig(
	ctx context.Context,
	tx repository.TxInterface,
	serviceID uuid.UUID,
	gitBranch *string,
	port *int,
	host *string,
	replicas *int32,
	autoDeploy *bool,
	runCommand *string,
	public *bool,
	image *string,
) (*ent.ServiceConfig, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.ServiceConfig.Create().
		SetServiceID(serviceID).
		SetNillableGitBranch(gitBranch).
		SetNillablePort(port).
		SetNillableHost(host).
		SetNillableReplicas(replicas).
		SetNillableAutoDeploy(autoDeploy).
		SetNillableRunCommand(runCommand).
		SetNillablePublic(public).
		SetNillableImage(image).
		Save(ctx)
}

// Update the service
func (self *ServiceRepository) Update(
	ctx context.Context,
	tx repository.TxInterface,
	serviceID uuid.UUID,
	displayName *string,
	description *string,
) error {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}
	return db.Service.Update().
		Where(service.IDEQ(serviceID)).
		SetNillableDisplayName(displayName).
		SetNillableDescription(description).
		Exec(ctx)
}

// Update service config
func (self *ServiceRepository) UpdateConfig(
	ctx context.Context,
	tx repository.TxInterface,
	serviceID uuid.UUID,
	gitBranch *string,
	port *int,
	host *string,
	replicas *int32,
	autoDeploy *bool,
	runCommand *string,
	public *bool,
	image *string,
) error {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.ServiceConfig.Update().
		Where(serviceconfig.ServiceID(serviceID)).
		SetNillableGitBranch(gitBranch).
		SetNillablePort(port).
		SetNillableHost(host).
		SetNillableReplicas(replicas).
		SetNillableAutoDeploy(autoDeploy).
		SetNillableRunCommand(runCommand).
		SetNillablePublic(public).
		SetNillableImage(image).
		Exec(ctx)
}
