// Code generated by ifacemaker; DO NOT EDIT.

package deployment_repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
)

// DeploymentRepositoryInterface ...
type DeploymentRepositoryInterface interface {
	Create(ctx context.Context, serviceID uuid.UUID, CommitSHA, CommitMessage string, committer *schema.GitCommitter, source schema.DeploymentSource) (*ent.Deployment, error)
	MarkStarted(ctx context.Context, buildJobID uuid.UUID) (*ent.Deployment, error)
	MarkFailed(ctx context.Context, buildJobID uuid.UUID, message string) (*ent.Deployment, error)
	MarkSucceeded(ctx context.Context, buildJobID uuid.UUID) (*ent.Deployment, error)
	// Cancels all jobs that are not in a finished state
	MarkCancelled(ctx context.Context, serviceID uuid.UUID) error
	// Mark cancelled by IDs
	MarkAsCancelled(ctx context.Context, jobIDs []uuid.UUID) error
	// Assigns the kubernetes "Job" name to the build job
	AssignKubernetesJobName(ctx context.Context, buildJobID uuid.UUID, jobName string) (*ent.Deployment, error)
	SetKubernetesJobStatus(ctx context.Context, buildJobID uuid.UUID, status string) (*ent.Deployment, error)
	GetJobsByStatus(ctx context.Context, status schema.DeploymentStatus) ([]*ent.Deployment, error)
	GetByServiceIDPaginated(ctx context.Context, serviceID uuid.UUID, cursor *time.Time, statusFilter []schema.DeploymentStatus) (jobs []*ent.Deployment, nextCursor *time.Time, err error)
}
