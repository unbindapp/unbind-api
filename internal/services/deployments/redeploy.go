package deployments_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	"github.com/unbindapp/unbind-api/internal/models"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
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
		canPullImage, _ := self.registryTester.CanPullImage(ctx, *deployment.Image)

		if canPullImage {
			// Update env
			envVars, err := self.resolveReferences(ctx, service)
			if err != nil {
				return nil, err
			}
			deployment.ResourceDefinition.Spec.EnvVars = envVars
			// Set volume data to current
			deployment.ResourceDefinition.Spec.Config.Volumes = schema.AsV1Volumes(service.Edges.ServiceConfig.Volumes)
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
	if (service.Type == schema.ServiceTypeDockerimage || service.Type == schema.ServiceTypeDatabase) && deployment.ResourceDefinition != nil {
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
		deployment.ResourceDefinition.Spec.Config.Volumes = schema.AsV1Volumes(service.Edges.ServiceConfig.Volumes)
		// For database, always use latest config
		if service.Type == schema.ServiceTypeDatabase && service.Edges.ServiceConfig.DatabaseConfig != nil {
			deployment.ResourceDefinition.Spec.Config.Database.Config = service.Edges.ServiceConfig.DatabaseConfig.AsV1DatabaseConfig()
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

	// Build a full deployment
	// Get git information if applicable
	var commitMessage string
	var commitSha string

	if deployment.CommitMessage != nil {
		commitMessage = *deployment.CommitMessage
	}

	// Enqueue build job
	env, err := self.deploymentController.PopulateBuildEnvironment(ctx, input.ServiceID, nil)
	if err != nil {
		return nil, err
	}

	if deployment.CommitSha != nil {
		commitSha = *deployment.CommitSha
		env["CHECKOUT_COMMIT_SHA"] = commitSha
	}

	job, err := self.deploymentController.EnqueueDeploymentJob(ctx, deployctl.DeploymentJobRequest{
		ServiceID:     input.ServiceID,
		Environment:   env,
		Source:        schema.DeploymentSourceManual,
		CommitSHA:     commitSha,
		CommitMessage: commitMessage,
		Committer:     deployment.CommitAuthor,
	})
	if err != nil {
		return nil, err
	}

	return models.TransformDeploymentEntity(job), nil
}
