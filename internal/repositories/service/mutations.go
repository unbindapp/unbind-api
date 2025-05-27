package service_repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
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
	TemplateID           *uuid.UUID
	TemplateInstanceID   *uuid.UUID
	ServiceGroupID       *uuid.UUID
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
		SetNillableTemplateID(input.TemplateID).
		SetNillableTemplateInstanceID(input.TemplateInstanceID).
		SetNillableServiceGroupID(input.ServiceGroupID).
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
	Icon                    *string
	Ports                   []schema.PortSpec
	Hosts                   []schema.HostSpec
	Replicas                *int32
	AutoDeploy              *bool
	InstallCommand          *string
	BuildCommand            *string
	RunCommand              *string
	Public                  *bool
	Image                   *string
	DockerfilePath          *string
	DockerfileContext       *string
	CustomDefinitionVersion *string
	DatabaseConfig          *schema.DatabaseConfig
	S3BackupSourceID        *uuid.UUID
	S3BackupBucket          *string
	BackupSchedule          *string
	BackupRetentionCount    *int
	SecurityContext         *schema.SecurityContext
	HealthCheck             *schema.HealthCheck
	VariableMounts          []*schema.VariableMount
	ProtectedVariables      *[]string
	Volumes                 *[]schema.ServiceVolume
	InitContainers          []*schema.InitContainer
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
	if input.Icon != nil {
		icon = *input.Icon
	} else if service.Database != nil {
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
		SetNillableInstallCommand(input.InstallCommand).
		SetNillableBuildCommand(input.BuildCommand).
		SetNillableRunCommand(input.RunCommand).
		SetNillableIsPublic(input.Public).
		SetNillableImage(input.Image).
		SetNillableDockerfilePath(input.DockerfilePath).
		SetNillableDockerfileContext(input.DockerfileContext).
		SetNillableDefinitionVersion(input.CustomDefinitionVersion).
		SetNillableS3BackupSourceID(input.S3BackupSourceID).
		SetNillableS3BackupBucket(input.S3BackupBucket).
		SetNillableBackupSchedule(input.BackupSchedule).
		SetNillableBackupRetentionCount(input.BackupRetentionCount)

	if input.InitContainers != nil {
		c.SetInitContainers(input.InitContainers)
	}

	if input.Volumes != nil {
		c.SetVolumes(*input.Volumes)
	}

	if input.ProtectedVariables != nil {
		c.SetProtectedVariables(*input.ProtectedVariables)
	}

	if len(input.VariableMounts) > 0 {
		c.SetVariableMounts(input.VariableMounts)
	}

	if input.HealthCheck != nil {
		c.SetHealthCheck(input.HealthCheck)
	}

	if input.SecurityContext != nil {
		c.SetSecurityContext(input.SecurityContext)
	}

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
		SetNillableIsPublic(input.Public).
		SetNillableImage(input.Image).
		SetNillableDefinitionVersion(input.CustomDefinitionVersion).
		SetNillableS3BackupSourceID(input.S3BackupSourceID).
		SetNillableBackupSchedule(input.BackupSchedule).
		SetNillableBackupRetentionCount(input.BackupRetentionCount)

	if input.InitContainers != nil {
		if len(input.InitContainers) > 0 {
			upd.SetInitContainers(input.InitContainers)
		} else {
			upd.ClearInitContainers()
		}
	}

	if input.ProtectedVariables != nil {
		upd.SetProtectedVariables(*input.ProtectedVariables)
	}

	if len(input.VariableMounts) > 0 {
		upd.SetVariableMounts(input.VariableMounts)
	} else if input.VariableMounts != nil {
		upd.ClearVariableMounts()
	}

	if input.HealthCheck != nil {
		if input.HealthCheck.Type == schema.HealthCheckTypeNone {
			upd.ClearHealthCheck()
		} else {
			upd.SetHealthCheck(input.HealthCheck)
		}
	}

	if input.SecurityContext != nil {
		upd.SetSecurityContext(input.SecurityContext)
	}

	if input.InstallCommand != nil {
		if *input.InstallCommand == "" {
			upd.ClearInstallCommand()
		} else {
			upd.SetInstallCommand(*input.InstallCommand)
		}
	}

	if input.BuildCommand != nil {
		if *input.BuildCommand == "" {
			upd.ClearBuildCommand()
		} else {
			upd.SetBuildCommand(*input.BuildCommand)
		}
	}

	if input.RunCommand != nil {
		if *input.RunCommand == "" {
			upd.ClearRunCommand()
		} else {
			upd.SetRunCommand(*input.RunCommand)
		}
	}

	if input.Volumes != nil {
		upd.SetVolumes(*input.Volumes)
	}

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

func (self *ServiceRepository) UpdateVariableMounts(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID, variableMounts []*schema.VariableMount) error {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.ServiceConfig.Update().
		Where(serviceconfig.ServiceID(serviceID)).
		SetVariableMounts(variableMounts).
		Exec(ctx)
}

func (self *ServiceRepository) UpdateDatabaseStorageSize(
	ctx context.Context,
	tx repository.TxInterface,
	serviceID uuid.UUID,
	newSize string,
) (*schema.DatabaseConfig, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	newSize = utils.EnsureSuffix(newSize, "Gi")

	svcConfig, err := db.Service.Query().
		Where(service.IDEQ(serviceID)).
		QueryServiceConfig().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if svcConfig.DatabaseConfig != nil {
		svcConfig.DatabaseConfig.StorageSize = newSize
	} else {
		svcConfig.DatabaseConfig = &schema.DatabaseConfig{
			StorageSize: newSize,
		}
	}

	updatedCfg, err := db.ServiceConfig.UpdateOneID(svcConfig.ID).
		SetDatabaseConfig(svcConfig.DatabaseConfig).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return updatedCfg.DatabaseConfig, nil
}
