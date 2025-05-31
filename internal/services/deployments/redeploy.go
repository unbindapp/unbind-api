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

func (self *DeploymentService) redeployExistingImage(ctx context.Context, service *ent.Service, deployment *ent.Deployment) (*models.DeploymentResponse, error) {
	// Update env
	envVars, err := self.resolveReferences(ctx, service)
	if err != nil {
		return nil, err
	}

	// Copy deployment
	newDeployment, err := self.repo.Deployment().CreateCopy(ctx, nil, deployment)
	if err != nil {
		return nil, err
	}

	// Update deployment resource definition
	deployment.ResourceDefinition.Spec.EnvVars = envVars
	deployment.ResourceDefinition.Spec.Config.Volumes = schema.AsV1Volumes(service.Edges.ServiceConfig.Volumes)

	// For docker image services, update the image reference
	if service.Type == schema.ServiceTypeDockerimage {
		deployment.ResourceDefinition.Spec.Config.Image = service.Edges.ServiceConfig.Image
		deployment.Image = utils.ToPtr(service.Edges.ServiceConfig.Image)
	}

	// For database services, always use latest config
	if service.Type == schema.ServiceTypeDatabase && service.Edges.ServiceConfig.DatabaseConfig != nil {
		deployment.ResourceDefinition.Spec.Config.Database.Config = service.Edges.ServiceConfig.DatabaseConfig.AsV1DatabaseConfig()
	}

	// Deploy to kubernetes
	_, _, err = self.k8s.DeployUnbindService(ctx, deployment.ResourceDefinition)
	if err != nil {
		// Mark failed
		if _, err := self.repo.Deployment().MarkFailed(ctx, nil, newDeployment.ID, err.Error(), time.Now()); err != nil {
			return nil, err
		}
		return nil, err
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

	// Check if we can redeploy without rebuilding
	if input.SmartRedeploy && deployment.ResourceDefinition != nil {
		canRedeploy := false

		// For non-database services, check if we can pull the existing image
		if service.Type != schema.ServiceTypeDatabase && deployment.Image != nil {
			canPullImage, _ := self.registryTester.CanPullImage(ctx, *deployment.Image)
			canRedeploy = canPullImage
		}

		// For database and docker image services, we can always redeploy
		if service.Type == schema.ServiceTypeDatabase || service.Type == schema.ServiceTypeDockerimage {
			canRedeploy = true
		}

		if canRedeploy {
			return self.redeployExistingImage(ctx, service, deployment)
		}
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
		ServiceID:         input.ServiceID,
		Environment:       env,
		Source:            schema.DeploymentSourceManual,
		CommitSHA:         commitSha,
		CommitMessage:     commitMessage,
		Committer:         deployment.CommitAuthor,
		DisableBuildCache: input.DisableBuildCache,
	})
	if err != nil {
		return nil, err
	}

	return models.TransformDeploymentEntity(job), nil
}
