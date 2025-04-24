package deployments_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *DeploymentService) CreateManualDeployment(ctx context.Context, requesterUserId uuid.UUID, input *models.CreateDeploymentInput) (*models.DeploymentResponse, error) {
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

	// Get git information if applicable
	var commitSHA string
	var commitMessage string
	var committer *schema.GitCommitter

	if service.Edges.GithubInstallation != nil && service.GitRepository != nil && service.Edges.ServiceConfig.GitBranch != nil {
		// Branch or sha
		summaryTarget := *service.Edges.ServiceConfig.GitBranch
		isCommitHash := false
		if input.GitSha != nil {
			summaryTarget = *input.GitSha
			isCommitHash = true
		}

		commitSHA, commitMessage, committer, err = self.githubClient.GetCommitSummary(ctx,
			service.Edges.GithubInstallation,
			service.Edges.GithubInstallation.AccountLogin,
			*service.GitRepository,
			summaryTarget,
			isCommitHash)

		// ! TODO - Should we hard fail here?
		if err != nil {
			return nil, err
		}
	}

	// Enqueue build job
	env, err := self.deploymentController.PopulateBuildEnvironment(ctx, input.ServiceID)
	if err != nil {
		return nil, err
	}
	if input.GitSha != nil {
		env["CHECKOUT_COMMIT_SHA"] = *input.GitSha
	}

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
