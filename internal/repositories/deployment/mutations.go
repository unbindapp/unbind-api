package deployment_repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/deployment"
	"github.com/unbindapp/unbind-api/ent/schema"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

func (self *DeploymentRepository) Create(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID, CommitSHA, CommitMessage string, committer *schema.GitCommitter, source schema.DeploymentSource) (*ent.Deployment, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	c := db.Deployment.Create().
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

func (self *DeploymentRepository) MarkStarted(ctx context.Context, tx repository.TxInterface, deploymentID uuid.UUID, startedAt time.Time) (*ent.Deployment, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.Deployment.UpdateOneID(deploymentID).
		SetStatus(schema.DeploymentStatusBuilding).
		// ! TODO - retry deployments?
		SetAttempts(1).
		SetStartedAt(startedAt).
		Save(ctx)
}

func (self *DeploymentRepository) MarkFailed(ctx context.Context, tx repository.TxInterface, deploymentID uuid.UUID, message string, failedAt time.Time) (*ent.Deployment, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.Deployment.UpdateOneID(deploymentID).
		Where(
			deployment.StatusNotIn(schema.DeploymentStatusCancelled, schema.DeploymentStatusSucceeded),
		).
		SetStatus(schema.DeploymentStatusFailed).
		SetCompletedAt(failedAt).
		SetError(message).
		Save(ctx)
}

func (self *DeploymentRepository) MarkSucceeded(ctx context.Context, tx repository.TxInterface, deploymentID uuid.UUID, completedAt time.Time) (*ent.Deployment, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.Deployment.UpdateOneID(deploymentID).
		SetStatus(schema.DeploymentStatusSucceeded).
		SetCompletedAt(completedAt).
		Save(ctx)
}

// Cancels all jobs that are not in a finished state
func (self *DeploymentRepository) MarkCancelledExcept(ctx context.Context, serviceID uuid.UUID, deploymentID uuid.UUID) error {
	return self.base.DB.Deployment.Update().
		SetStatus(schema.DeploymentStatusCancelled).
		SetCompletedAt(time.Now()).
		Where(
			deployment.ServiceIDEQ(serviceID),
			deployment.IDNEQ(deploymentID),
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
func (self *DeploymentRepository) AssignKubernetesJobName(ctx context.Context, deploymentID uuid.UUID, jobName string) (*ent.Deployment, error) {
	return self.base.DB.Deployment.UpdateOneID(deploymentID).
		SetKubernetesJobName(jobName).
		Save(ctx)
}

func (self *DeploymentRepository) SetKubernetesJobStatus(ctx context.Context, deploymentID uuid.UUID, status string) (*ent.Deployment, error) {
	return self.base.DB.Deployment.UpdateOneID(deploymentID).
		SetKubernetesJobStatus(status).
		Save(ctx)
}

func (self *DeploymentRepository) AttachDeploymentMetadata(ctx context.Context, tx repository.TxInterface, deploymentID uuid.UUID, imageName string, resourceDefinition *v1.Service) (*ent.Deployment, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	// Prune any sensitive data
	if resourceDefinition != nil {
		resourceDefinition.Spec.EnvVars = []corev1.EnvVar{}
	}

	return db.Deployment.UpdateOneID(deploymentID).
		SetImage(imageName).
		SetResourceDefinition(resourceDefinition).
		Save(ctx)
}

// Create a copy with all metadata, except for failed_at, completed_at, and status
func (self *DeploymentRepository) CreateCopy(ctx context.Context, tx repository.TxInterface, deployment *ent.Deployment) (*ent.Deployment, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.Deployment.Create().
		SetServiceID(deployment.ServiceID).
		SetStatus(schema.DeploymentStatusQueued).
		SetNillableCommitSha(deployment.CommitSha).
		SetNillableCommitMessage(deployment.CommitMessage).
		SetCommitAuthor(deployment.CommitAuthor).
		SetResourceDefinition(deployment.ResourceDefinition).
		SetSource(schema.DeploymentSourceManual).
		SetNillableImage(deployment.Image).
		Save(ctx)
}
