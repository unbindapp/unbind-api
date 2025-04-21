package service_repo

import (
	"context"
	"fmt"

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
	kubernetesName string,
	name string,
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
		SetKubernetesName(kubernetesName).
		SetName(name).
		SetDescription(description).
		SetEnvironmentID(environmentID).
		SetNillableGithubInstallationID(gitHubInstallationID).
		SetNillableGitRepository(gitRepository).
		SetNillableGitRepositoryOwner(gitRepositoryOwner).
		SetKubernetesSecret(kubernetesSecret).
		Save(ctx)
}

// Create the service config
type MutateConfigInput struct {
	ServiceID               uuid.UUID
	ServiceType             schema.ServiceType
	Builder                 *schema.ServiceBuilder
	Provider                *enum.Provider
	Framework               *enum.Framework
	GitBranch               *string
	Ports                   []schema.PortSpec
	Hosts                   []v1.HostSpec
	Replicas                *int32
	AutoDeploy              *bool
	RunCommand              *string
	Public                  *bool
	Image                   *string
	DockerfilePath          *string
	DockerfileContext       *string
	Database                *string
	CustomDefinitionVersion *string
	DatabaseConfig          *map[string]interface{}
	DatabaseVersion         *string
}

func (self *ServiceRepository) CreateConfig(
	ctx context.Context,
	tx repository.TxInterface,
	input *MutateConfigInput,
) (*ent.ServiceConfig, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	if input == nil || input.Builder == nil {
		return nil, fmt.Errorf("builder is missing, but required")
	}

	// Get high level icon
	var icon string
	if input.Database != nil {
		icon = *input.Database
	} else if input.Framework != nil {
		icon = string(*input.Framework)
	} else if input.Provider != nil {
		icon = string(*input.Provider)
	} else {
		icon = string(input.ServiceType)
	}

	c := db.ServiceConfig.Create().
		SetServiceID(input.ServiceID).
		SetType(input.ServiceType).
		SetBuilder(*input.Builder).
		SetIcon(icon).
		SetNillableRailpackProvider(input.Provider).
		SetNillableRailpackFramework(input.Framework).
		SetNillableGitBranch(input.GitBranch).
		SetNillableReplicas(input.Replicas).
		SetNillableAutoDeploy(input.AutoDeploy).
		SetNillableRunCommand(input.RunCommand).
		SetNillablePublic(input.Public).
		SetNillableImage(input.Image).
		SetNillableDockerfilePath(input.DockerfilePath).
		SetNillableDockerfileContext(input.DockerfileContext).
		SetNillableDatabase(input.Database).
		SetNillableDefinitionVersion(input.CustomDefinitionVersion).
		SetNillableDatabaseVersion(input.DatabaseVersion)

	if input.DatabaseConfig != nil {
		c.SetDatabaseConfig(*input.DatabaseConfig)
	}

	if len(input.Ports) > 0 {
		c.SetPorts(input.Ports)
	}

	if len(input.Hosts) > 0 {
		c.SetHosts(input.Hosts)
	}

	return c.Save(ctx)
}

// Update the service
func (self *ServiceRepository) Update(
	ctx context.Context,
	tx repository.TxInterface,
	serviceID uuid.UUID,
	name *string,
	description *string,
) error {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}
	return db.Service.Update().
		Where(service.IDEQ(serviceID)).
		SetNillableName(name).
		SetNillableDescription(description).
		Exec(ctx)
}

// Update service config
func (self *ServiceRepository) UpdateConfig(
	ctx context.Context,
	tx repository.TxInterface,
	input *MutateConfigInput,
) error {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	upd := db.ServiceConfig.Update().
		Where(serviceconfig.ServiceID(input.ServiceID)).
		SetNillableBuilder(input.Builder).
		SetNillableGitBranch(input.GitBranch).
		SetNillableReplicas(input.Replicas).
		SetNillableAutoDeploy(input.AutoDeploy).
		SetNillableRunCommand(input.RunCommand).
		SetNillablePublic(input.Public).
		SetNillableImage(input.Image).
		SetNillableDefinitionVersion(input.CustomDefinitionVersion)

	if input.DatabaseConfig != nil {
		upd.SetDatabaseConfig(*input.DatabaseConfig)
	}

	if input.DatabaseVersion != nil {
		upd.SetDatabaseVersion(*input.DatabaseVersion)
	}

	if input.DockerfilePath != nil {
		if *input.DockerfilePath == "" {
			upd.ClearDockerfilePath()
		} else {
			upd.SetDockerfilePath(*input.DockerfilePath)
		}
	}

	if input.DockerfileContext != nil {
		if *input.DockerfileContext == "" {
			upd.ClearDockerfileContext()
		} else {
			upd.SetDockerfileContext(*input.DockerfileContext)
		}
	}

	if len(input.Ports) > 0 {
		upd.SetPorts(input.Ports)
	}
	if len(input.Hosts) > 0 {
		upd.SetHosts(input.Hosts)
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
