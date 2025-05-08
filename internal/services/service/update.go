package service_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	service_repo "github.com/unbindapp/unbind-api/internal/repositories/service"
	"github.com/unbindapp/unbind-api/internal/services/models"
	webhooks_service "github.com/unbindapp/unbind-api/internal/services/webooks"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// UpdateServiceConfigInput defines the input for updating a service configuration
type UpdateServiceInput struct {
	TeamID        uuid.UUID `format:"uuid" required:"true" json:"team_id"`
	ProjectID     uuid.UUID `format:"uuid" required:"true" json:"project_id"`
	EnvironmentID uuid.UUID `format:"uuid" required:"true" json:"environment_id"`
	ServiceID     uuid.UUID `format:"uuid" required:"true" json:"service_id"`
	Name          *string   `required:"false" json:"name"`
	Description   *string   `required:"false" json:"description"`

	// Configuration
	GitBranch         *string                `json:"git_branch,omitempty" required:"false"`
	GitTag            *string                `json:"git_tag,omitempty" required:"false" doc:"Tag to build from, supports glob patterns"`
	Builder           *schema.ServiceBuilder `json:"builder,omitempty" required:"false"`
	Hosts             []v1.HostSpec          `json:"hosts,omitempty" required:"false"`
	Ports             []schema.PortSpec      `json:"ports,omitempty" required:"false"`
	Replicas          *int32                 `json:"replicas,omitempty" required:"false"`
	AutoDeploy        *bool                  `json:"auto_deploy,omitempty" required:"false"`
	RunCommand        *string                `json:"run_command,omitempty" required:"false"`
	IsPublic          *bool                  `json:"is_public,omitempty" required:"false"`
	Image             *string                `json:"image,omitempty" required:"false"`
	DockerfilePath    *string                `json:"dockerfile_path,omitempty" required:"false" doc:"Optional path to Dockerfile, if using docker builder - set empty string to reset to default"`
	DockerfileContext *string                `json:"dockerfile_context,omitempty" required:"false" doc:"Optional path to Dockerfile context, if using docker builder - set empty string to reset to default"`

	// Databases
	DatabaseConfig       *schema.DatabaseConfig `json:"database_config,omitempty"`
	S3BackupEndpointID   *uuid.UUID             `json:"s3_backup_endpoint_id,omitempty" format:"uuid"`
	S3BackupBucket       *string                `json:"s3_backup_bucket,omitempty"`
	BackupSchedule       *string                `json:"backup_schedule,omitempty" required:"false" doc:"Cron expression for the backup schedule, e.g. '0 0 * * *'"`
	BackupRetentionCount *int                   `json:"backup_retention,omitempty" required:"false" doc:"Number of base backups to retain, e.g. 3"`

	// PVC
	PVCID        *string `json:"pvc_id,omitempty" doc:"ID of the PVC to attach to the service"`
	PVCMountPath *string `json:"pvc_mount_path,omitempty" required:"false" doc:"Mount path for the PVC"`
}

// UpdateService updates a service and its configuration
func (self *ServiceService) UpdateService(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *UpdateServiceInput) (*models.ServiceResponse, error) {
	// Verify tag if present
	if input.GitTag != nil {
		if !utils.IsValidGlobPattern(*input.GitTag) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Invalid git tag")
		}
	}

	// Check permissions
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to admin service
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   input.ServiceID,
		},
	}

	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Verify inputs
	_, _, err := self.VerifyInputs(ctx, input.TeamID, input.ProjectID, input.EnvironmentID)
	if err != nil {
		return nil, err
	}

	// Perform update
	service, err := self.repo.Service().GetByID(ctx, input.ServiceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
		}
		return nil, err
	}

	if service.Type == schema.ServiceTypeDockerimage || service.Type == schema.ServiceTypeDatabase {
		if input.Builder != nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Cannot update builder for docker image or database service")
		}
	}

	// For database we don't want to set ports
	if service.Type == schema.ServiceTypeDatabase {
		input.Ports = nil

		// Check backup schedule
		if input.BackupSchedule != nil {
			if err := utils.ValidateCronExpression(*input.BackupSchedule); err != nil {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("invalid backup schedule: %s", err))
			}
		}

		// Disallow attaching a PVC to a database service
		if input.PVCID != nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Cannot attach a PVC to a database service")
		}
	}

	// PVC validation, requires a path
	if input.PVCID != nil {
		if input.PVCMountPath == nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "PVC mount path is required")
		}
		// Validate unix-style path
		if !utils.IsValidUnixPath(*input.PVCMountPath) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Invalid PVC mount path")
		}
	}

	// For database we can't set version if deployed
	if service.Type == schema.ServiceTypeDatabase && input.DatabaseConfig != nil && service.DatabaseVersion != nil {
		hasDeployment := len(service.Edges.Deployments) > 0
		if hasDeployment {
			// * special rule that you can't update version if there is a deployment
			if input.DatabaseConfig.Version != "" {
				if input.DatabaseConfig.Version != *service.DatabaseVersion {
					return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Cannot update version for database service with existing deployment")
				}
			}
		}
	}

	// Verify storage size changes if applicable
	if input.DatabaseConfig != nil {
		if input.DatabaseConfig.StorageSize == "" {
			// Set to existing
			if service.Edges.ServiceConfig.DatabaseConfig != nil {
				input.DatabaseConfig.StorageSize = service.Edges.ServiceConfig.DatabaseConfig.StorageSize
				// Sort of a DB migration I guess
				if input.DatabaseConfig.StorageSize == "" {
					input.DatabaseConfig.StorageSize = "1Gi"
				}
			}
		} else {
			// Parse
			newSizeTarget, err := resource.ParseQuantity(input.DatabaseConfig.StorageSize)
			if err != nil {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Invalid storage size")
			}

			// Parse existing (if present)
			if service.Edges.ServiceConfig.DatabaseConfig != nil && service.Edges.ServiceConfig.DatabaseConfig.StorageSize != "" {
				existingSizeTarget, err := resource.ParseQuantity(service.Edges.ServiceConfig.DatabaseConfig.StorageSize)
				if err != nil {
					existingSizeTarget = resource.MustParse("1Gi")
				}
				// Compare
				if newSizeTarget.Cmp(existingSizeTarget) < 0 {
					return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Cannot decrease storage size")
				}
			}
		}
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Check if PVC is in use by a service
	if input.PVCID != nil {
		self.validatePVC(ctx, input.TeamID, input.ProjectID, input.EnvironmentID, *input.PVCID, service.Edges.Environment.Edges.Project.Edges.Team.Namespace, client)
	}

	// Verify backup sources (for databases)
	// Make sure we can read and write to the S3 bucket provided
	if service.Type == schema.ServiceTypeDatabase && input.S3BackupEndpointID != nil && input.S3BackupBucket != nil {
		// Check if the S3 source exists
		s3Source, err := self.repo.S3().GetByID(ctx, *input.S3BackupEndpointID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "S3 endpoint not found")
			}
			return nil, err
		}

		if err := self.verifyS3Access(ctx, s3Source, *input.S3BackupBucket, service.Edges.Environment.Edges.Project.Edges.Team.Namespace, client); err != nil {
			return nil, err
		}
	}

	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		// Update the service
		if err := self.repo.Service().Update(ctx, tx, input.ServiceID, input.Name, input.Description); err != nil {
			return fmt.Errorf("failed to update service: %w", err)
		}

		if len(service.Edges.ServiceConfig.Hosts) < 1 &&
			input.IsPublic != nil && *input.IsPublic && len(input.Hosts) < 1 && service.Type != schema.ServiceTypeDatabase &&
			(len(input.Ports) > 0 || len(service.Edges.ServiceConfig.Ports) > 0) {

			// Figure out ports
			var ports []schema.PortSpec
			if len(input.Ports) > 0 {
				ports = input.Ports
			}

			if len(service.Edges.ServiceConfig.Ports) > 0 {
				ports = append(ports, service.Edges.ServiceConfig.Ports...)
			}

			generatedHost, err := self.generateWildcardHost(ctx, tx, service.KubernetesName, ports)
			if err != nil {
				return fmt.Errorf("failed to generate wildcard host: %w", err)
			}
			if generatedHost == nil {
				input.IsPublic = utils.ToPtr(false)
			} else {
				input.Hosts = append(input.Hosts, *generatedHost)
			}
		}

		// Update the service config
		updateInput := &service_repo.MutateConfigInput{
			ServiceID:            input.ServiceID,
			Builder:              input.Builder,
			GitBranch:            input.GitBranch,
			GitTag:               input.GitTag,
			Ports:                input.Ports,
			Hosts:                input.Hosts,
			Replicas:             input.Replicas,
			AutoDeploy:           input.AutoDeploy,
			RunCommand:           input.RunCommand,
			Public:               input.IsPublic,
			Image:                input.Image,
			DockerfilePath:       input.DockerfilePath,
			DockerfileContext:    input.DockerfileContext,
			DatabaseConfig:       input.DatabaseConfig,
			S3BackupEndpointID:   input.S3BackupEndpointID,
			S3BackupBucket:       input.S3BackupBucket,
			BackupSchedule:       input.BackupSchedule,
			BackupRetentionCount: input.BackupRetentionCount,
			PVCID:                input.PVCID,
			PVCVolumeMountPath:   input.PVCMountPath,
		}
		if err := self.repo.Service().UpdateConfig(ctx, tx, updateInput); err != nil {
			return fmt.Errorf("failed to update service config: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// Re-fetch the service
	service, err = self.repo.Service().GetByID(ctx, input.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to re-fetch service: %w", err)
	}

	deployments, err := self.RedeployServices(ctx, []*ent.Service{service})
	if err != nil {
		return nil, err
	}

	// Update service deployment
	var newDeployment *ent.Deployment
	if len(deployments) > 0 {
		newDeployment = deployments[0]
		if newDeployment.Status == schema.DeploymentStatusSucceeded {
			service.Edges.CurrentDeployment = newDeployment
		}
		service.Edges.Deployments = []*ent.Deployment{
			newDeployment,
		}
	}

	// Trigger webhook
	go func() {
		event := schema.WebhookEventServiceUpdated
		level := webhooks_service.WebhookLevelInfo

		// Get service with edges
		service, err := self.repo.Service().GetByID(context.Background(), service.ID)
		if err != nil {
			log.Errorf("Failed to get service %s: %v", service.ID.String(), err)
			return
		}

		// Construct URL
		basePath, _ := utils.JoinURLPaths(
			self.cfg.ExternalUIUrl,
			input.TeamID.String(),
			"project",
			input.ProjectID.String(),
		)
		url := basePath + "?environment=" + input.EnvironmentID.String() +
			"&service=" + service.ID.String()
		// Get user
		user, err := self.repo.User().GetByID(context.Background(), requesterUserID)
		if err != nil {
			log.Errorf("Failed to get user %s: %v", requesterUserID.String(), err)
			return
		}
		data := webhooks_service.WebhookData{
			Title: "Service Updated",
			Url:   url,
			Fields: []webhooks_service.WebhookDataField{
				{
					Name:  "Service Name",
					Value: service.Name,
				},
				{
					Name:  "Project & Environment",
					Value: fmt.Sprintf("%s > %s", service.Edges.Environment.Edges.Project.Name, service.Edges.Environment.Name),
				},
				{
					Name:  "Updated By",
					Value: user.Email,
				},
			},
		}

		if input.GitBranch != nil {
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Git Branch",
				Value: *input.GitBranch,
			})
		}

		if input.Image != nil {
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Image",
				Value: *input.Image,
			})
		}

		if input.Replicas != nil {
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Replicas",
				Value: fmt.Sprintf("%d", *input.Replicas),
			})
		}

		if input.AutoDeploy != nil {
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Auto Deploy",
				Value: fmt.Sprintf("%t", *input.AutoDeploy),
			})
		}

		if input.RunCommand != nil {
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Run Command",
				Value: *input.RunCommand,
			})
		}

		if input.IsPublic != nil {
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Public",
				Value: fmt.Sprintf("%t", *input.IsPublic),
			})
		}

		if input.DockerfilePath != nil {
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Dockerfile Path",
				Value: *input.DockerfilePath,
			})
		}

		if input.DockerfileContext != nil {
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Dockerfile Context",
				Value: *input.DockerfileContext,
			})
		}

		if len(service.Edges.ServiceConfig.Hosts) > 0 {
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Service URL",
				Value: fmt.Sprintf("https://%s", service.Edges.ServiceConfig.Hosts[0].Host),
			})
		}

		if newDeployment != nil {
			deploymentUrl, _ := utils.JoinURLPaths(self.cfg.ExternalUIUrl, input.TeamID.String(), "project", input.ProjectID.String(), "?environment="+input.EnvironmentID.String(), "&service="+service.ID.String(), "&deployment="+newDeployment.ID.String())
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Deployment",
				Value: deploymentUrl,
			})
		}

		if err := self.webhookService.TriggerWebhooks(context.Background(), level, event, data); err != nil {
			log.Errorf("Failed to trigger webhook %s: %v", event, err)
		}
	}()

	return models.TransformServiceEntity(service), nil
}
