// Code generated by ifacemaker; DO NOT EDIT.

package deployment_repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
)

// DeploymentRepositoryInterface ...
type DeploymentRepositoryInterface interface {
	Create(ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID, CommitSHA, CommitMessage string, GitBranch string, committer *schema.GitCommitter, source schema.DeploymentSource, initialStatus schema.DeploymentStatus) (*ent.Deployment, error)
	MarkQueued(ctx context.Context, tx repository.TxInterface, deploymentID uuid.UUID, queuedAt time.Time) (*ent.Deployment, error)
	MarkStarted(ctx context.Context, tx repository.TxInterface, deploymentID uuid.UUID, startedAt time.Time) (*ent.Deployment, error)
	MarkFailed(ctx context.Context, tx repository.TxInterface, deploymentID uuid.UUID, message string, failedAt time.Time) (*ent.Deployment, error)
	MarkSucceeded(ctx context.Context, tx repository.TxInterface, deploymentID uuid.UUID, completedAt time.Time) (*ent.Deployment, error)
	// Cancels all jobs that are not in a finished state
	MarkCancelledExcept(ctx context.Context, serviceID uuid.UUID, deploymentID uuid.UUID) error
	// Mark cancelled by IDs
	MarkAsCancelled(ctx context.Context, jobIDs []uuid.UUID) error
	// Assigns the kubernetes "Job" name to the build job
	AssignKubernetesJobName(ctx context.Context, deploymentID uuid.UUID, jobName string) (*ent.Deployment, error)
	SetKubernetesJobStatus(ctx context.Context, deploymentID uuid.UUID, status string) (*ent.Deployment, error)
	AttachDeploymentMetadata(ctx context.Context, tx repository.TxInterface, deploymentID uuid.UUID, imageName string, resourceDefinition *v1.Service) (*ent.Deployment, error)
	// Create a copy with all metadata, except for failed_at, completed_at, and status
	CreateCopy(ctx context.Context, tx repository.TxInterface, deployment *ent.Deployment) (*ent.Deployment, error)
	GetByID(ctx context.Context, deploymentID uuid.UUID) (*ent.Deployment, error)
	ExistsInEnvironment(ctx context.Context, deploymentID uuid.UUID, environmentID uuid.UUID) (bool, error)
	ExistsInProject(ctx context.Context, deploymentID uuid.UUID, projectID uuid.UUID) (bool, error)
	ExistsInTeam(ctx context.Context, deploymentID uuid.UUID, teamID uuid.UUID) (bool, error)
	GetLastSuccessfulDeployment(ctx context.Context, serviceID uuid.UUID) (*ent.Deployment, error)
	GetJobsByStatus(ctx context.Context, status schema.DeploymentStatus) ([]*ent.Deployment, error)
	GetByServiceIDPaginated(ctx context.Context, serviceID uuid.UUID, perPage int, cursor *time.Time, statusFilter []schema.DeploymentStatus) (jobs []*ent.Deployment, nextCursor *time.Time, err error)
}
