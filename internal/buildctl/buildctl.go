package buildctl

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/infrastructure/queue"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	"github.com/valkey-io/valkey-go"
)

// Valkey key for the queue
const BUILDER_QUEUE_KEY = "unbind:build:queue"

// The request to build a service, includes environment for builder image
type BuildJobRequest struct {
	ServiceID   uuid.UUID         `json:"service_id"`
	Environment map[string]string `json:"environment"`
}

// Handles triggering builds for services
type BuildController struct {
	k8s        *k8s.KubeClient
	jobQueue   *queue.Queue[BuildJobRequest]
	ctx        context.Context
	cancelFunc context.CancelFunc
	repo       *repositories.Repositories
}

func NewBuildController(ctx context.Context, k8s *k8s.KubeClient, valkeyClient valkey.Client, repositories *repositories.Repositories) *BuildController {
	jobQueue := queue.NewQueue[BuildJobRequest](valkeyClient, BUILDER_QUEUE_KEY)

	// Create a cancellable context
	ctx, cancelFunc := context.WithCancel(ctx)

	return &BuildController{
		k8s:        k8s,
		jobQueue:   jobQueue,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		repo:       repositories,
	}
}

// Start queue processor
func (self *BuildController) Start() {
	// Start the job processor
	self.jobQueue.StartProcessor(self.ctx, self.processJob)

	// Start the job status synchronizer
	go self.startStatusSynchronizer()
}

// Stop stops the build job manager
func (self *BuildController) Stop() {
	self.cancelFunc()
}

// startStatusSynchronizer periodically synchronizes job statuses with Kubernetes
func (self *BuildController) startStatusSynchronizer() {
	ticker := time.NewTicker(30 * time.Second) // Sync every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-self.ctx.Done():
			return
		case <-ticker.C:
			if err := self.SyncJobStatuses(context.Background()); err != nil {
				log.Error("Failed to sync job statuses", "err", err)
			}
		}
	}
}

// EnqueueBuildJob adds a build job to the queue
func (self *BuildController) EnqueueBuildJob(ctx context.Context, req BuildJobRequest) (jobID string, err error) {
	// Create a record in the database
	job, err := self.repo.BuildJob().Create(
		ctx,
		req.ServiceID,
	)

	if err != nil {
		return "", fmt.Errorf("failed to create build job record: %w", err)
	}

	// Add to the queue
	err = self.jobQueue.Enqueue(ctx, job.ID.String(), req)
	if err != nil {
		return "", fmt.Errorf("failed to enqueue job: %w", err)
	}

	return job.ID.String(), nil
}

// processJob processes a job from the queue
func (self *BuildController) processJob(ctx context.Context, item *queue.QueueItem[BuildJobRequest]) error {
	jobID, _ := uuid.Parse(item.ID)
	req := item.Data

	// Update the job status in the database
	err := self.repo.BuildJob().MarkCancelled(ctx, req.ServiceID)
	if err != nil {
		log.Warnf("Failed to mark job as cancelled: %v service: %s", err, req.ServiceID)
	}
	_, err = self.repo.BuildJob().MarkStarted(ctx, jobID)

	if err != nil {
		return fmt.Errorf("failed to mark job started: %w", err)
	}

	// Start the actual Kubernetes job
	k8sJobName, err := self.k8s.CreateBuildJob(ctx, req.ServiceID.String(), jobID.String(), req.Environment)
	if err != nil {
		log.Error("Failed to create Kubernetes job", "err", err)

		// Update status to failed
		_, dbErr := self.repo.BuildJob().MarkFailed(ctx, jobID, err.Error())

		if dbErr != nil {
			log.Error("Failed to update job failure status", "err", dbErr)
		}

		return err
	}

	// Update the Kubernetes job name in the database
	_, err = self.repo.BuildJob().AssignKubernetesJobName(ctx, jobID, k8sJobName)

	if err != nil {
		log.Error("Failed to update Kubernetes job name in database", "err", err, "jobID", jobID, "k8sJobName", k8sJobName)
	}

	return nil
}

// SyncJobStatuses synchronizes the status of all processing jobs with Kubernetes
func (self *BuildController) SyncJobStatuses(ctx context.Context) error {
	// Get all job marked running status
	jobs, err := self.repo.BuildJob().GetJobsByStatus(ctx, schema.BuildJobStatusRunning)
	if err != nil {
		return fmt.Errorf("failed to query processing jobs: %w", err)
	}

	for _, job := range jobs {
		if job.KubernetesJobName == "" {
			continue
		}

		k8sStatus, err := self.k8s.GetJobStatus(ctx, job.KubernetesJobName)
		if err != nil {
			log.Error("Failed to get Kubernetes job status", "err", err, "jobID", job.ID)
			continue
		}

		// Update based on Kubernetes status
		switch k8sStatus.ConditionType {
		case k8s.JobComplete:
			_, err = self.repo.BuildJob().MarkCompleted(ctx, job.ID)
		case k8s.JobFailed:
			_, err = self.repo.BuildJob().MarkFailed(ctx, job.ID, k8sStatus.FailureReason)
		default:
			_, err = self.repo.BuildJob().SetKubernetesJobStatus(ctx, job.ID, k8sStatus.ConditionType.String())
		}

		if err != nil {
			log.Error("Failed to update job status", "err", err, "jobID", job.ID)
		}
	}

	return nil
}
