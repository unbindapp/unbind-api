package deployment_repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/deployment"
	"github.com/unbindapp/unbind-api/ent/schema"
)

func (self *DeploymentRepository) Create(ctx context.Context, serviceID uuid.UUID) (*ent.Deployment, error) {
	return self.base.DB.Deployment.Create().
		SetServiceID(serviceID).
		SetStatus(schema.DeploymentStatusQueued).
		Save(ctx)
}

func (self *DeploymentRepository) MarkStarted(ctx context.Context, buildJobID uuid.UUID) (*ent.Deployment, error) {
	return self.base.DB.Deployment.UpdateOneID(buildJobID).
		SetStatus(schema.DeploymentStatusRunning).
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

func (self *DeploymentRepository) MarkCompleted(ctx context.Context, buildJobID uuid.UUID) (*ent.Deployment, error) {
	return self.base.DB.Deployment.UpdateOneID(buildJobID).
		SetStatus(schema.DeploymentStatusCompleted).
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
			deployment.StatusNotIn(schema.DeploymentStatusFailed, schema.DeploymentStatusCancelled, schema.DeploymentStatusCompleted),
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
			deployment.StatusNotIn(schema.DeploymentStatusRunning, schema.DeploymentStatusFailed, schema.DeploymentStatusCancelled, schema.DeploymentStatusCompleted),
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
