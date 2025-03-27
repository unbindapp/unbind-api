package service_service

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
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
	GitBranch  *string                `json:"git_branch,omitempty" required:"false"`
	Type       *schema.ServiceType    `json:"type,omitempty" required:"false"`
	Builder    *schema.ServiceBuilder `json:"builder,omitempty" required:"false"`
	Hosts      []v1.HostSpec          `json:"hosts,omitempty" required:"false"`
	Ports      []v1.PortSpec          `json:"ports,omitempty" required:"false"`
	Replicas   *int32                 `json:"replicas,omitempty" required:"false"`
	AutoDeploy *bool                  `json:"auto_deploy,omitempty" required:"false"`
	RunCommand *string                `json:"run_command,omitempty" required:"false"`
	Public     *bool                  `json:"public,omitempty" required:"false"`
	Image      *string                `json:"image,omitempty" required:"false"`
}

// UpdateService updates a service and its configuration
func (self *ServiceService) UpdateService(ctx context.Context, requesterUserID uuid.UUID, input *UpdateServiceInput) (*models.ServiceResponse, error) {
	// Validate input
	if err := validate.Validator().Struct(input); err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	// Check permissions
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to manage teams
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to manage this team
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   input.TeamID.String(),
		},
		// Has permission to manage projects
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   "*",
		},
		// Has permission to manage this project
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   input.ProjectID.String(),
		},
		// Has permission to manage this specific environment
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   input.EnvironmentID.String(),
		},
		// Has permission to update this service
		{
			Action:       permission.ActionUpdate,
			ResourceType: permission.ResourceTypeService,
			ResourceID:   input.ServiceID.String(),
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
	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		// Update the service
		if err := self.repo.Service().Update(ctx, tx, input.ServiceID, input.DisplayName, input.Description); err != nil {
			return fmt.Errorf("failed to update service: %w", err)
		}

		if len(service.Edges.ServiceConfig.Hosts) < 1 && input.Public != nil && *input.Public && len(input.Hosts) < 1 {
			host, err := utils.GenerateSubdomain(service.DisplayName, self.cfg.ExternalWildcardBaseURL)
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
		if err := self.repo.Service().UpdateConfig(ctx,
			tx,
			input.ServiceID,
			input.Type,
			input.Builder,
			input.GitBranch,
			input.Ports,
			input.Hosts,
			input.Replicas,
			input.AutoDeploy,
			input.RunCommand,
			input.Public,
			input.Image); err != nil {
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
	if service.Edges.CurrentDeployment != nil && service.Edges.CurrentDeployment.ResourceDefinition != nil {
		// Create a an object with only fields we care to compare
		existingCrd := &v1.Service{
			Spec: v1.ServiceSpec{
				Config: v1.ServiceConfigSpec{
					GitBranch:  service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.GitBranch,
					Hosts:      service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.Hosts,
					Replicas:   service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.Replicas,
					Ports:      service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.Ports,
					AutoDeploy: service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.AutoDeploy,
					RunCommand: service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.RunCommand,
					Public:     service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.Public,
				},
			},
		}
		// Create a new CRD to compare it
		var gitBranch string
		if service.Edges.ServiceConfig.GitBranch != nil {
			gitBranch = *service.Edges.ServiceConfig.GitBranch
		}
		newCrd := &v1.Service{
			Spec: v1.ServiceSpec{
				Config: v1.ServiceConfigSpec{
					GitBranch:  gitBranch,
					Hosts:      service.Edges.ServiceConfig.Hosts,
					Replicas:   utils.ToPtr(service.Edges.ServiceConfig.Replicas),
					Ports:      service.Edges.ServiceConfig.Ports,
					AutoDeploy: service.Edges.ServiceConfig.AutoDeploy,
					RunCommand: service.Edges.ServiceConfig.RunCommand,
					Public:     service.Edges.ServiceConfig.Public,
				},
			},
		}

		// Changing git branch requires a whole new build
		if existingCrd.Spec.Config.GitBranch != newCrd.Spec.Config.GitBranch {
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
		} else if !reflect.DeepEqual(existingCrd, newCrd) {
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
			crdToDeploy.Spec.Config.GitBranch = gitBranch
			crdToDeploy.Spec.Config.Hosts = service.Edges.ServiceConfig.Hosts
			crdToDeploy.Spec.Config.Replicas = utils.ToPtr(service.Edges.ServiceConfig.Replicas)
			crdToDeploy.Spec.Config.Ports = service.Edges.ServiceConfig.Ports
			crdToDeploy.Spec.Config.AutoDeploy = service.Edges.ServiceConfig.AutoDeploy
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

				deployment, err := self.repo.Deployment().Create(
					ctx,
					tx,
					service.Edges.CurrentDeployment.ServiceID,
					commitSha,
					commitMessage,
					service.Edges.CurrentDeployment.CommitAuthor,
					service.Edges.CurrentDeployment.Source,
				)
				log.Info("Created deployment", "sha", commitSha, "message", commitMessage, "current_deployment", service.Edges.CurrentDeployment.ID)
				if err != nil {
					return err
				}

				// Mark the deployment as started
				if _, err := self.repo.Deployment().MarkStarted(ctx, tx, deployment.ID, time.Now()); err != nil {
					return err
				}

				// Deploy to kubernetes
				_, newService, err := self.k8s.DeployUnbindService(ctx, crdToDeploy)
				if err != nil {
					// Mark failed
					if _, err := self.repo.Deployment().MarkFailed(ctx, tx, deployment.ID, err.Error(), time.Now()); err != nil {
						return err
					}
					// Pass through since we already marked the deployment as failed
					return nil
				}

				// Attach metadata
				if _, err := self.repo.Deployment().AttachDeploymentMetadata(
					ctx,
					tx,
					deployment.ID,
					crdToDeploy.Spec.Config.Image,
					newService,
				); err != nil {
					// Mark failed
					if _, err := self.repo.Deployment().MarkFailed(ctx, tx, deployment.ID, err.Error(), time.Now()); err != nil {
						return err
					}
					// Pass through since we already marked the deployment as failed
					return nil
				}

				// Mark as succeeded
				if _, err := self.repo.Deployment().MarkSucceeded(ctx, tx, deployment.ID, time.Now()); err != nil {
					return err
				}

				// Update the service with the new deployment
				if err := self.repo.Service().SetCurrentDeployment(ctx, tx, service.ID, deployment.ID); err != nil {
					return err
				}

				return nil
			}); err != nil {
				log.Errorf("failed to deploy new CRD: %v", err)
				return nil, err
			}
		}
	}

	return models.TransformServiceEntity(service), nil
}
