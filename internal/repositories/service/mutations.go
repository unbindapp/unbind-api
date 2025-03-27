package service_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
)

// Create the service
func (self *ServiceRepository) Create(
	ctx context.Context,
	tx repository.TxInterface,
	name string,
	displayName string,
	description string,
	environmentID uuid.UUID,
	gitHubInstallationID *int64,
	gitRepository *string,
	gitRepositoryOwner *string,
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
		SetEnvironmentID(environmentID).
		SetNillableGithubInstallationID(gitHubInstallationID).
		SetNillableGitRepository(gitRepository).
		SetNillableGitRepositoryOwner(gitRepositoryOwner).
		SetKubernetesSecret(kubernetesSecret).
		Save(ctx)
}

// Create the service config
func (self *ServiceRepository) CreateConfig(
	ctx context.Context,
	tx repository.TxInterface,
	serviceID uuid.UUID,
	serviceType schema.ServiceType,
	builder schema.ServiceBuilder,
	provider *enum.Provider,
	framework *enum.Framework,
	gitBranch *string,
	ports []v1.PortSpec,
	hosts []v1.HostSpec,
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

	c := db.ServiceConfig.Create().
		SetServiceID(serviceID).
		SetType(serviceType).
		SetBuilder(builder).
		SetNillableProvider(provider).
		SetNillableFramework(framework).
		SetNillableGitBranch(gitBranch).
		SetNillableReplicas(replicas).
		SetNillableAutoDeploy(autoDeploy).
		SetNillableRunCommand(runCommand).
		SetNillablePublic(public).
		SetNillableImage(image)

	if len(ports) > 0 {
		c.SetPorts(ports)
	}

	if len(hosts) > 0 {
		c.SetHosts(hosts)
	}

	return c.Save(ctx)
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
	serviceType *schema.ServiceType,
	builder *schema.ServiceBuilder,
	gitBranch *string,
	ports []v1.PortSpec,
	hosts []v1.HostSpec,
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

	upd := db.ServiceConfig.Update().
		Where(serviceconfig.ServiceID(serviceID)).
		SetNillableType(serviceType).
		SetNillableBuilder(builder).
		SetNillableGitBranch(gitBranch).
		SetNillableReplicas(replicas).
		SetNillableAutoDeploy(autoDeploy).
		SetNillableRunCommand(runCommand).
		SetNillablePublic(public).
		SetNillableImage(image)

	if len(ports) > 0 {
		upd.SetPorts(ports)
	}
	if len(hosts) > 0 {
		upd.SetHosts(hosts)
	}

	return upd.Exec(ctx)
}

func (self *ServiceRepository) Delete(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID) error {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.Service.DeleteOneID(serviceID).Exec(ctx)
}

func (self *ServiceRepository) SetCurrentDeployment(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID, deploymentID uuid.UUID) error {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}
	return db.Service.UpdateOneID(serviceID).SetCurrentDeploymentID(deploymentID).Exec(ctx)
}
