package buildjob_repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/buildjob"
)

func (self *BuildJobRepository) Create(ctx context.Context, serviceID uuid.UUID) (*ent.BuildJob, error) {
	return self.base.DB.BuildJob.Create().
		SetServiceID(serviceID).
		SetStatus(buildjob.StatusQueued).
		Save(ctx)
}

func (self *BuildJobRepository) MarkStarted(ctx context.Context, buildJobID uuid.UUID) (*ent.BuildJob, error) {
	return self.base.DB.BuildJob.UpdateOneID(buildJobID).
		SetStatus(buildjob.StatusRunning).
		SetStartedAt(time.Now()).
		Save(ctx)
}

func (self *BuildJobRepository) MarkFailed(ctx context.Context, buildJobID uuid.UUID, message string) (*ent.BuildJob, error) {
	return self.base.DB.BuildJob.UpdateOneID(buildJobID).
		SetStatus(buildjob.StatusFailed).
		SetCompletedAt(time.Now()).
		SetError(message).
		Save(ctx)
}

func (self *BuildJobRepository) MarkCompleted(ctx context.Context, buildJobID uuid.UUID) (*ent.BuildJob, error) {
	return self.base.DB.BuildJob.UpdateOneID(buildJobID).
		SetStatus(buildjob.StatusCompleted).
		SetCompletedAt(time.Now()).
		Save(ctx)
}

// Cancels all jobs that are not in a finished state
func (self *BuildJobRepository) MarkCancelled(ctx context.Context, serviceID uuid.UUID) error {
	return self.base.DB.BuildJob.Update().
		SetStatus(buildjob.StatusCancelled).
		SetCompletedAt(time.Now()).
		Where(
			buildjob.ServiceIDEQ(serviceID),
			buildjob.StatusNotIn(buildjob.StatusFailed, buildjob.StatusCancelled, buildjob.StatusCompleted),
		).
		Exec(ctx)
}

// Assigns the kubernetes "Job" name to the build job
func (self *BuildJobRepository) AssignKubernetesJobName(ctx context.Context, buildJobID uuid.UUID, jobName string) (*ent.BuildJob, error) {
	return self.base.DB.BuildJob.UpdateOneID(buildJobID).
		SetKubernetesJobName(jobName).
		Save(ctx)
}

func (self *BuildJobRepository) SetKubernetesJobStatus(ctx context.Context, buildJobID uuid.UUID, status string) (*ent.BuildJob, error) {
	return self.base.DB.BuildJob.UpdateOneID(buildJobID).
		SetKubernetesJobStatus(status).
		Save(ctx)
}
