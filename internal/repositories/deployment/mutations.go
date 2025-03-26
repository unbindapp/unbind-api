package deployment_repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/deployment"
	"github.com/unbindapp/unbind-api/ent/schema"
)

func (self *DeploymentRepository) Create(ctx context.Context, serviceID uuid.UUID, CommitSHA, CommitMessage string, committer *schema.GitCommitter, source schema.DeploymentSource) (*ent.Deployment, error) {
	c := self.base.DB.Deployment.Create().
		SetServiceID(serviceID).
		SetStatus(schema.DeploymentStatusQueued).
		SetSource(source).
		SetCommitAuthor(committer)
	if CommitSHA != "" {
		c.SetCommitSha(CommitSHA)
	}
	if CommitMessage != "" {
		c.SetCommitMessage(CommitMessage)
	}
	return c.Save(ctx)
}

func (self *DeploymentRepository) MarkStarted(ctx context.Context, buildJobID uuid.UUID) (*ent.Deployment, error) {
	return self.base.DB.Deployment.UpdateOneID(buildJobID).
		SetStatus(schema.DeploymentStatusBuilding).
		// ! TODO - retry deployments?
		SetAttempts(1).
		SetStartedAt(time.Now()).
		Save(ctx)
}

func (self *DeploymentRepository) MarkFailed(ctx context.Context, buildJobID uuid.UUID, message string) (*ent.Deployment, error) {
	return self.base.DB.Deployment.UpdateOneID(buildJobID).
		SetStatus(schema.DeploymentStatusFailed).
		SetCompletedAt(time.Now()).
		SetError(message).
		Save(ctx)
}

func (self *DeploymentRepository) MarkSucceeded(ctx context.Context, buildJobID uuid.UUID) (*ent.Deployment, error) {
	return self.base.DB.Deployment.UpdateOneID(buildJobID).
		SetStatus(schema.DeploymentStatusSucceeded).
		SetCompletedAt(time.Now()).
		Save(ctx)
}

// Cancels all jobs that are not in a finished state
func (self *DeploymentRepository) MarkCancelled(ctx context.Context, serviceID uuid.UUID) error {
	return self.base.DB.Deployment.Update().
		SetStatus(schema.DeploymentStatusCancelled).
		SetCompletedAt(time.Now()).
		Where(
			deployment.ServiceIDEQ(serviceID),
			deployment.StatusNotIn(schema.DeploymentStatusFailed, schema.DeploymentStatusCancelled, schema.DeploymentStatusSucceeded),
		).
		Exec(ctx)
}

// Mark cancelled by IDs
func (self *DeploymentRepository) MarkAsCancelled(ctx context.Context, jobIDs []uuid.UUID) error {
	return self.base.DB.Deployment.Update().
		SetStatus(schema.DeploymentStatusCancelled).
		SetCompletedAt(time.Now()).
		Where(
			deployment.IDIn(jobIDs...),
			deployment.StatusNotIn(schema.DeploymentStatusBuilding, schema.DeploymentStatusFailed, schema.DeploymentStatusCancelled, schema.DeploymentStatusSucceeded),
		).
		Exec(ctx)
}

// Assigns the kubernetes "Job" name to the build job
func (self *DeploymentRepository) AssignKubernetesJobName(ctx context.Context, buildJobID uuid.UUID, jobName string) (*ent.Deployment, error) {
	return self.base.DB.Deployment.UpdateOneID(buildJobID).
		SetKubernetesJobName(jobName).
		Save(ctx)
}

func (self *DeploymentRepository) SetKubernetesJobStatus(ctx context.Context, buildJobID uuid.UUID, status string) (*ent.Deployment, error) {
	return self.base.DB.Deployment.UpdateOneID(buildJobID).
		SetKubernetesJobStatus(status).
		Save(ctx)
}
