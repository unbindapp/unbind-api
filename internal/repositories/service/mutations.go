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

type CreateServiceInput struct {
	KubernetesName       string
	Name                 string
	ServiceType          schema.ServiceType
	Description          string
	EnvironmentID        uuid.UUID
	GitHubInstallationID *int64
	GitRepository        *string
	GitRepositoryOwner   *string
	KubernetesSecret     string
	Database             *string
	DatabaseVersion      *string
}

// Create the service
func (self *ServiceRepository) Create(
	ctx context.Context,
	tx repository.TxInterface,
	input *CreateServiceInput,
) (*ent.Service, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.Service.Create().
		SetType(input.ServiceType).
		SetKubernetesName(input.KubernetesName).
		SetName(input.Name).
		SetDescription(input.Description).
		SetEnvironmentID(input.EnvironmentID).
		SetNillableGithubInstallationID(input.GitHubInstallationID).
		SetNillableGitRepository(input.GitRepository).
		SetNillableGitRepositoryOwner(input.GitRepositoryOwner).
		SetKubernetesSecret(input.KubernetesSecret).
		SetNillableDatabase(input.Database).
		SetNillableDatabaseVersion(input.DatabaseVersion).Save(ctx)
}

// Create the service config
type MutateConfigInput struct {
	ServiceID               uuid.UUID
	Builder                 *schema.ServiceBuilder
	Provider                *enum.Provider
	Framework               *enum.Framework
	GitBranch               *string
	GitTag                  *string
	Ports                   []schema.PortSpec
	Hosts                   []v1.HostSpec
	Replicas                *int32
	AutoDeploy              *bool
	RunCommand              *string
	Public                  *bool
	Image                   *string
	DockerfilePath          *string
	DockerfileContext       *string
	CustomDefinitionVersion *string
	DatabaseConfig          *schema.DatabaseConfig
	S3BackupEndpointID      *uuid.UUID
	S3BackupBucket          *string
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

	// Get service
	service, err := db.Service.Get(ctx, input.ServiceID)
	if err != nil {
		return nil, err
	}

	// Get high level icon
	var icon string
	if service.Database != nil {
		icon = *service.Database
	} else if input.Framework != nil {
		icon = string(*input.Framework)
	} else if input.Provider != nil {
		icon = string(*input.Provider)
	} else {
		icon = string(service.Type)
	}

	c := db.ServiceConfig.Create().
		SetServiceID(input.ServiceID).
		SetBuilder(*input.Builder).
		SetIcon(icon).
		SetNillableRailpackProvider(input.Provider).
		SetNillableRailpackFramework(input.Framework).
		SetNillableGitBranch(input.GitBranch).
		SetNillableReplicas(input.Replicas).
		SetNillableAutoDeploy(input.AutoDeploy).
		SetNillableRunCommand(input.RunCommand).
		SetNillableIsPublic(input.Public).
		SetNillableImage(input.Image).
		SetNillableDockerfilePath(input.DockerfilePath).
		SetNillableDockerfileContext(input.DockerfileContext).
		SetNillableDefinitionVersion(input.CustomDefinitionVersion).
		SetNillableS3BackupEndpointID(input.S3BackupEndpointID).
		SetNillableS3BackupBucket(input.S3BackupBucket)

	if input.DatabaseConfig != nil {
		c.SetDatabaseConfig(input.DatabaseConfig)
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
		SetNillableReplicas(input.Replicas).
		SetNillableAutoDeploy(input.AutoDeploy).
		SetNillableRunCommand(input.RunCommand).
		SetNillableIsPublic(input.Public).
		SetNillableImage(input.Image).
		SetNillableDefinitionVersion(input.CustomDefinitionVersion).
		SetNillableS3BackupEndpointID(input.S3BackupEndpointID)

	if input.GitBranch != nil {
		if *input.GitBranch == "" {
			upd.ClearGitBranch()
		} else {
			upd.SetGitBranch(*input.GitBranch)
		}
	}
	if input.GitTag != nil {
		if *input.GitTag == "" {
			upd.ClearGitTag()
		} else {
			upd.SetGitTag(*input.GitTag)
		}
	}

	if input.S3BackupBucket != nil {
		if *input.S3BackupBucket == "" {
			upd.ClearS3BackupBucket()
		} else {
			upd.SetS3BackupBucket(*input.S3BackupBucket)
		}
	}

	if input.DatabaseConfig != nil {
		upd.SetDatabaseConfig(input.DatabaseConfig)

		if input.DatabaseConfig.Version != "" {
			// Set on service
			_, err := db.Service.UpdateOneID(input.ServiceID).
				SetDatabaseVersion(input.DatabaseConfig.Version).
				Save(ctx)
			if err != nil {
				return err
			}
		}
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
