package deployments_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
	corev1 "k8s.io/api/core/v1"
)

func (self *DeploymentService) resolveReferences(ctx context.Context, service *ent.Service) ([]corev1.EnvVar, error) {
	// Resolve references
	additionalEnv, err := self.variableService.ResolveAllReferences(ctx, service.ID)
	if err != nil {
		return nil, err
	}
	envVars := make([]corev1.EnvVar, len(additionalEnv))
	i := 0
	for k, v := range additionalEnv {
		envVars[i] = corev1.EnvVar{
			Name:  k,
			Value: v,
		}
		i++
	}

	return envVars, nil
}

func (self *DeploymentService) CreateRedeployment(ctx context.Context, requesterUserId uuid.UUID, input *models.RedeployExistingDeploymentInput) (*models.DeploymentResponse, error) {
	if err := validate.Validator().Struct(input); err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	// Editor can create deployments
	if err := self.repo.Permissions().Check(ctx, requesterUserId, []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   input.ServiceID,
		},
	}); err != nil {
		return nil, err
	}

	service, err := self.validateInputs(ctx, input)
	if err != nil {
		return nil, err
	}

	// Get existing deployment
	deployment, err := self.repo.Deployment().GetByID(ctx, input.DeploymentID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Deployment not found")
		}
		return nil, err
	}

	if deployment.Image != nil && deployment.ResourceDefinition != nil {
		canPullImage, _ := self.imageChecker.CanPullImage(ctx, *deployment.Image)

		if canPullImage {
			// Update env
			envVars, err := self.resolveReferences(ctx, service)
			if err != nil {
				return nil, err
			}
			deployment.ResourceDefinition.Spec.EnvVars = envVars
			// We can just, re-deploy the existing CRD spec
			// Deploy to kubernetes
			_, _, err = self.k8s.DeployUnbindService(ctx, deployment.ResourceDefinition)
			if err != nil {
				return nil, err
			}
			// Update current deployment on DB
			if err := self.repo.Service().SetCurrentDeployment(ctx, nil, deployment.ServiceID, deployment.ID); err != nil {
				return nil, err
			}
			return models.TransformDeploymentEntity(deployment), nil
		}
	}

	// For docker-image builds, if we don't have a deployment image we can get it from the config
	// re-fetch service
	service, err = self.repo.Service().GetByID(ctx, input.ServiceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
		}
		return nil, err
	}

	if (service.Edges.ServiceConfig.Type == schema.ServiceTypeDockerimage || service.Edges.ServiceConfig.Type == schema.ServiceTypeDatabase) && deployment.ResourceDefinition != nil {
		deployment.ResourceDefinition.Spec.Config.Image = service.Edges.ServiceConfig.Image
		envVars, err := self.resolveReferences(ctx, service)
		if err != nil {
			return nil, err
		}
		deployment.ResourceDefinition.Spec.EnvVars = envVars
		// Copy deployment
		deployment.Image = utils.ToPtr(service.Edges.ServiceConfig.Image)
		newDeployment, err := self.repo.Deployment().CreateCopy(ctx, nil, deployment)
		if err != nil {
			return nil, err
		}
		_, _, err = self.k8s.DeployUnbindService(ctx, deployment.ResourceDefinition)
		if err != nil {
			// Mark failed
			if _, err := self.repo.Deployment().MarkFailed(ctx, nil, newDeployment.ID, err.Error(), time.Now()); err != nil {
				return nil, err
			}
		}
		// Mark as succeeded
		newDeployment, err = self.repo.Deployment().MarkSucceeded(ctx, nil, newDeployment.ID, time.Now())
		if err != nil {
			return nil, err
		}
		// Update the service with the new deployment
		if err := self.repo.Service().SetCurrentDeployment(ctx, nil, service.ID, newDeployment.ID); err != nil {
			return nil, err
		}
		return models.TransformDeploymentEntity(newDeployment), nil
	}

	// If we ended up here not as a git service we failed
	if service.Edges.ServiceConfig.Type != schema.ServiceTypeGithub || deployment.CommitSha == nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Unable to re-deploy service")
	}

	// Build a full deployment
	// Get git information if applicable
	var commitSHA string
	var commitMessage string
	var committer *schema.GitCommitter

	if service.Edges.GithubInstallation != nil && service.GitRepository != nil && service.Edges.ServiceConfig.GitBranch != nil {
		commitSHA, commitMessage, committer, err = self.githubClient.GetCommitSummary(ctx,
			service.Edges.GithubInstallation,
			service.Edges.GithubInstallation.AccountLogin,
			*service.GitRepository,
			*deployment.CommitSha,
			true)

		if err != nil {
			return nil, err
		}
	}

	// Enqueue build job
	env, err := self.deploymentController.PopulateBuildEnvironment(ctx, input.ServiceID)
	if err != nil {
		return nil, err
	}
	env["CHECKOUT_COMMIT_SHA"] = *deployment.CommitSha

	job, err := self.deploymentController.EnqueueDeploymentJob(ctx, deployctl.DeploymentJobRequest{
		ServiceID:     input.ServiceID,
		Environment:   env,
		Source:        schema.DeploymentSourceManual,
		CommitSHA:     commitSHA,
		CommitMessage: commitMessage,
		Committer:     committer,
	})
	if err != nil {
		return nil, err
	}

	return models.TransformDeploymentEntity(job), nil
}
