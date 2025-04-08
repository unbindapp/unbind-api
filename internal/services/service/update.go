package service_service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	service_repo "github.com/unbindapp/unbind-api/internal/repositories/service"
	"github.com/unbindapp/unbind-api/internal/services/models"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
)

// UpdateServiceConfigInput defines the input for updating a service configuration
type UpdateServiceInput struct {
	TeamID        uuid.UUID `validate:"required,uuid4" required:"true" json:"team_id"`
	ProjectID     uuid.UUID `validate:"required,uuid4" required:"true" json:"project_id"`
	EnvironmentID uuid.UUID `validate:"required,uuid4" required:"true" json:"environment_id"`
	ServiceID     uuid.UUID `validate:"required,uuid4" required:"true" json:"service_id"`
	DisplayName   *string   ` required:"false" json:"display_name"`
	Description   *string   ` required:"false" json:"description"`

	// Configuration
	GitBranch         *string                `json:"git_branch,omitempty" required:"false"`
	Builder           *schema.ServiceBuilder `json:"builder,omitempty" required:"false"`
	Hosts             []v1.HostSpec          `json:"hosts,omitempty" required:"false"`
	Ports             []v1.PortSpec          `json:"ports,omitempty" required:"false"`
	Replicas          *int32                 `json:"replicas,omitempty" required:"false"`
	AutoDeploy        *bool                  `json:"auto_deploy,omitempty" required:"false"`
	RunCommand        *string                `json:"run_command,omitempty" required:"false"`
	Public            *bool                  `json:"public,omitempty" required:"false"`
	Image             *string                `json:"image,omitempty" required:"false"`
	DockerfilePath    *string                `json:"dockerfile_path,omitempty" required:"false" doc:"Optional path to Dockerfile, if using docker builder - set empty string to reset to default"`
	DockerfileContext *string                `json:"dockerfile_context,omitempty" required:"false" doc:"Optional path to Dockerfile context, if using docker builder - set empty string to reset to default"`

	// Databases
	DatabaseConfig *map[string]interface{} `json:"database_config,omitempty"`
}

// UpdateService updates a service and its configuration
func (self *ServiceService) UpdateService(ctx context.Context, requesterUserID uuid.UUID, input *UpdateServiceInput) (*models.ServiceResponse, error) {
	// Validate input
	if err := validate.Validator().Struct(input); err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
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

	if service.Edges.ServiceConfig.Type == schema.ServiceTypeDockerimage || service.Edges.ServiceConfig.Type == schema.ServiceTypeDatabase {
		if input.Builder != nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Cannot update builder for docker image or database service")
		}
	}

	// For database we don't want to set ports
	if service.Edges.ServiceConfig.Type == schema.ServiceTypeDatabase {
		input.Ports = nil
	}

	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		// Update the service
		if err := self.repo.Service().Update(ctx, tx, input.ServiceID, input.DisplayName, input.Description); err != nil {
			return fmt.Errorf("failed to update service: %w", err)
		}

		if len(service.Edges.ServiceConfig.Hosts) < 1 && input.Public != nil && *input.Public && len(input.Hosts) < 1 && service.Edges.ServiceConfig.Type != schema.ServiceTypeDatabase {
			host, err := utils.GenerateSubdomain(service.Name, self.cfg.ExternalWildcardBaseURL)
			if err != nil {
				log.Warn("failed to generate subdomain", "error", err)
			} else {
				input.Hosts = []v1.HostSpec{
					{
						Host: host,
						Path: "/",
					},
				}
				// Figure out ports
				if len(input.Ports) > 0 {
					input.Hosts[0].Port = &input.Ports[0].Port
				} else {
					if len(service.Edges.ServiceConfig.Ports) > 0 {
						input.Hosts[0].Port = &service.Edges.ServiceConfig.Ports[0].Port
					}
				}
			}
		}

		// Update the service config
		updateInput := &service_repo.MutateConfigInput{
			ServiceID:         input.ServiceID,
			Builder:           input.Builder,
			GitBranch:         input.GitBranch,
			Ports:             input.Ports,
			Hosts:             input.Hosts,
			Replicas:          input.Replicas,
			AutoDeploy:        input.AutoDeploy,
			RunCommand:        input.RunCommand,
			Public:            input.Public,
			Image:             input.Image,
			DockerfilePath:    input.DockerfilePath,
			DockerfileContext: input.DockerfileContext,
			DatabaseConfig:    input.DatabaseConfig,
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

	// See if a deployment is needed
	needsDeploymentType, err := self.repo.Service().NeedsDeployment(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("failed to check if deployment is needed: %w", err)
	}

	// Populated if we deploy something in this process
	var newDeployment *ent.Deployment

	switch needsDeploymentType {
	case service_repo.NeedsBuildAndDeployment:
		// New full build needed
		env, err := self.deploymentController.PopulateBuildEnvironment(ctx, input.ServiceID)
		if err != nil {
			return nil, err
		}

		var commitSHA string
		var commitMessage string
		var committer *schema.GitCommitter

		if service.Edges.GithubInstallation != nil && service.GitRepository != nil && service.Edges.ServiceConfig.GitBranch != nil {
			commitSHA, commitMessage, committer, err = self.githubClient.GetBranchHeadSummary(ctx,
				service.Edges.GithubInstallation,
				service.Edges.GithubInstallation.AccountLogin,
				*service.GitRepository,
				*service.Edges.ServiceConfig.GitBranch)

			// ! TODO - Should we hard fail here?
			if err != nil {
				return nil, err
			}
		}

		_, err = self.deploymentController.EnqueueDeploymentJob(ctx, deployctl.DeploymentJobRequest{
			ServiceID:     input.ServiceID,
			Environment:   env,
			Source:        schema.DeploymentSourceManual,
			CommitSHA:     commitSHA,
			CommitMessage: commitMessage,
			Committer:     committer,
		})
		if err != nil {
			log.Errorf("failed to enqueue deployment job: %v", err)
			return nil, err
		}
	case service_repo.NeedsDeployment:
		// New adhoc deployment needed
		crdToDeploy := &v1.Service{}
		// Metadata
		crdToDeploy.Name = service.Edges.CurrentDeployment.ResourceDefinition.Name
		crdToDeploy.Namespace = service.Edges.CurrentDeployment.ResourceDefinition.Namespace
		crdToDeploy.Kind = service.Edges.CurrentDeployment.ResourceDefinition.Kind
		crdToDeploy.APIVersion = service.Edges.CurrentDeployment.ResourceDefinition.APIVersion
		crdToDeploy.Labels = service.Edges.CurrentDeployment.ResourceDefinition.Labels
		crdToDeploy.Spec = service.Edges.CurrentDeployment.ResourceDefinition.Spec

		// Update the Spec
		var gitBranch string
		if service.Edges.ServiceConfig.GitBranch != nil {
			gitBranch = *service.Edges.ServiceConfig.GitBranch
			if !strings.HasPrefix(gitBranch, "refs/heads/") {
				gitBranch = fmt.Sprintf("refs/heads/%s", gitBranch)
			}
		}
		crdToDeploy.Spec.Config.GitBranch = gitBranch
		crdToDeploy.Spec.Config.Hosts = service.Edges.ServiceConfig.Hosts
		crdToDeploy.Spec.Config.Replicas = utils.ToPtr(service.Edges.ServiceConfig.Replicas)
		crdToDeploy.Spec.Config.Ports = service.Edges.ServiceConfig.Ports
		crdToDeploy.Spec.Config.RunCommand = service.Edges.ServiceConfig.RunCommand
		crdToDeploy.Spec.Config.Public = service.Edges.ServiceConfig.Public

		// Deploy the new CRD
		if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
			// Create a record in the database
			commitSha := ""
			if service.Edges.CurrentDeployment.CommitSha != nil {
				commitSha = *service.Edges.CurrentDeployment.CommitSha
			}
			commitMessage := ""
			if service.Edges.CurrentDeployment.CommitMessage != nil {
				commitMessage = *service.Edges.CurrentDeployment.CommitMessage
			}

			newDeployment, err = self.repo.Deployment().Create(
				ctx,
				tx,
				service.Edges.CurrentDeployment.ServiceID,
				commitSha,
				commitMessage,
				service.Edges.CurrentDeployment.CommitAuthor,
				service.Edges.CurrentDeployment.Source,
			)
			if err != nil {
				return err
			}

			// Mark the deployment as started
			if _, err := self.repo.Deployment().MarkStarted(ctx, tx, newDeployment.ID, time.Now()); err != nil {
				return err
			}

			// Deploy to kubernetes
			_, newService, err := self.k8s.DeployUnbindService(ctx, crdToDeploy)
			if err != nil {
				// Mark failed
				if _, err := self.repo.Deployment().MarkFailed(ctx, tx, newDeployment.ID, err.Error(), time.Now()); err != nil {
					return err
				}
				// Pass through since we already marked the deployment as failed
				return nil
			}

			// Attach metadata
			if _, err := self.repo.Deployment().AttachDeploymentMetadata(
				ctx,
				tx,
				newDeployment.ID,
				crdToDeploy.Spec.Config.Image,
				newService,
			); err != nil {
				// Mark failed
				if _, err := self.repo.Deployment().MarkFailed(ctx, tx, newDeployment.ID, err.Error(), time.Now()); err != nil {
					return err
				}
				// Pass through since we already marked the deployment as failed
				return nil
			}

			// Mark as succeeded
			if _, err := self.repo.Deployment().MarkSucceeded(ctx, tx, newDeployment.ID, time.Now()); err != nil {
				return err
			}

			// Update the service with the new deployment
			if err := self.repo.Service().SetCurrentDeployment(ctx, tx, service.ID, newDeployment.ID); err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Errorf("failed to deploy new CRD: %v", err)
			return nil, err
		}
	}

	// Update service deployment
	if newDeployment != nil {
		if newDeployment.Status == schema.DeploymentStatusSucceeded {
			service.Edges.CurrentDeployment = newDeployment
		}
		service.Edges.Deployments = []*ent.Deployment{
			newDeployment,
		}
	}

	return models.TransformServiceEntity(service), nil
}
