package deployctl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/infrastructure/queue"
	"github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	"github.com/valkey-io/valkey-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Valkey key for the queue
const BUILDER_QUEUE_KEY = "unbind:build:queue"

// The request to deploy a service, includes environment for builder image
type DeploymentJobRequest struct {
	ServiceID   uuid.UUID         `json:"service_id"`
	Environment map[string]string `json:"environment"`
}

// Handles triggering builds for services
type DeploymentController struct {
	cfg          *config.Config
	k8s          *k8s.KubeClient
	jobQueue     *queue.Queue[DeploymentJobRequest]
	ctx          context.Context
	cancelFunc   context.CancelFunc
	repo         *repositories.Repositories
	githubClient *github.GithubClient
}

func NewDeploymentController(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, k8s *k8s.KubeClient, valkeyClient valkey.Client, repositories *repositories.Repositories, githubClient *github.GithubClient) *DeploymentController {
	jobQueue := queue.NewQueue[DeploymentJobRequest](valkeyClient, BUILDER_QUEUE_KEY)

	return &DeploymentController{
		cfg:          cfg,
		k8s:          k8s,
		jobQueue:     jobQueue,
		ctx:          ctx,
		cancelFunc:   cancel,
		repo:         repositories,
		githubClient: githubClient,
	}
}

// Start queue processor
func (self *DeploymentController) StartAsync() {
	// Start the job processor
	self.jobQueue.StartProcessor(self.ctx, self.processJob)

	// Start the job status synchronizer
	go self.startStatusSynchronizer()
}

// Stop stops the deployment manager
func (self *DeploymentController) Stop() {
	self.cancelFunc()
}

// startStatusSynchronizer periodically synchronizes job statuses with Kubernetes
func (self *DeploymentController) startStatusSynchronizer() {
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

// Populate build environment
func (self *DeploymentController) PopulateBuildEnvironment(ctx context.Context, serviceID uuid.UUID) (map[string]string, error) {
	// Get the service
	service, err := self.repo.Service().GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	// Get deployment namespace
	namespace, err := self.repo.Service().GetDeploymentNamespace(ctx, service.ID)

	// Get build secrets
	// ! Use our cluster config for this
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting in-cluster config: %v", err)
	}
	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	buildSecrets, err := self.k8s.GetSecretMap(ctx, service.KubernetesSecret, namespace, client)
	if err != nil {
		log.Error("Error getting secrets", "err", err)
		return nil, huma.Error500InternalServerError("Failed to get build secrets")
	}

	// Convert the byte arrays to base64 strings first
	serializableSecrets := make(map[string]string)
	for k, v := range buildSecrets {
		serializableSecrets[k] = base64.StdEncoding.EncodeToString(v)
	}

	// Serialize the map to JSON
	secretsJSON, err := json.Marshal(serializableSecrets)
	if err != nil {
		log.Error("Error marshalling secrets", "err", err)
		return nil, huma.Error500InternalServerError("Failed to marshal secrets")
	}

	// Populate environment
	env := map[string]string{
		"CONTAINER_REGISTRY_HOST":     self.cfg.ContainerRegistryHost,
		"CONTAINER_REGISTRY_USER":     self.cfg.ContainerRegistryUser,
		"CONTAINER_REGISTRY_PASSWORD": self.cfg.ContainerRegistryPassword,
		"DEPLOYMENT_NAMESPACE":        namespace,
		"SERVICE_PUBLIC":              strconv.FormatBool(service.Edges.ServiceConfig.Public),
		"SERVICE_REPLICAS":            strconv.Itoa(int(service.Edges.ServiceConfig.Replicas)),
		"SERVICE_SECRET_NAME":         service.KubernetesSecret,
		"SERVICE_BUILD_SECRETS":       string(secretsJSON),
	}

	// Add Github fields
	if service.GithubInstallationID != nil {
		if service.GitRepository == nil || service.Edges.ServiceConfig.GitBranch == nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Missing required fields for Github service - doesn't have repository or git branch")
		}
		// Get private key for the service's github app.
		// ! TODO - we can probably reduce these queries
		privKey, err := self.repo.Service().GetGithubPrivateKey(ctx, service.ID)
		if err != nil {
			log.Error("Error getting github private key", "err", err)
			return nil, err
		}

		env["GITHUB_APP_PRIVATE_KEY"] = privKey
		env["GITHUB_INSTALLATION_ID"] = strconv.Itoa(int(*service.GithubInstallationID))
		// Get GitHub installation
		installation, err := self.repo.Github().GetInstallationByID(ctx, *service.GithubInstallationID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "GitHub installation not found")
			}
			return nil, err
		}
		env["GITHUB_APP_ID"] = strconv.Itoa(int(installation.GithubAppID))

		// Verify repository access
		canAccess, cloneUrl, err := self.githubClient.VerifyRepositoryAccess(ctx, installation, installation.AccountLogin, *service.GitRepository)
		if err != nil {
			log.Error("Error verifying repository access", "err", err)
			return nil, err
		}

		if !canAccess || err != nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Repository not accessible with the specified GitHub installation")
		}

		env["GITHUB_REPO_URL"] = cloneUrl

		ref := *service.Edges.ServiceConfig.GitBranch
		if !strings.HasPrefix(ref, "refs/head") {
			ref = "refs/heads/" + ref
		}
		env["GIT_REF"] = ref
	}

	if service.Provider != nil {
		env["SERVICE_PROVIDER"] = string(*service.Provider)
	}

	if service.Framework != nil {
		env["SERVICE_FRAMEWORK"] = string(*service.Framework)
	}

	if service.Edges.ServiceConfig.Port != nil {
		env["SERVICE_PORT"] = strconv.Itoa(*service.Edges.ServiceConfig.Port)
	}

	if service.Edges.ServiceConfig.Host != nil {
		env["SERVICE_HOST"] = *service.Edges.ServiceConfig.Host
	}

	// ! TODO - we need to support the custom run commands, the operator supports it

	return env, nil
}

// EnqueueDeploymentJob adds a deployment to the queue
func (self *DeploymentController) EnqueueDeploymentJob(ctx context.Context, req DeploymentJobRequest) (job *ent.Deployment, err error) {
	// Cancel any existing queued jobs
	if err := self.cancelExistingJobs(ctx, req.ServiceID); err != nil {
		return nil, fmt.Errorf("failed to cancel existing jobs: %w", err)
	}

	// Create a record in the database
	job, err = self.repo.Deployment().Create(
		ctx,
		req.ServiceID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create deployment record: %w", err)
	}

	// Add to the queue
	err = self.jobQueue.Enqueue(ctx, job.ID.String(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue job: %w", err)
	}

	return job, nil
}

// cancelExistingJobs marks all pending jobs for a service as cancelled in the DB
// and removes them from the queue
func (self *DeploymentController) cancelExistingJobs(ctx context.Context, serviceID uuid.UUID) error {
	// 1. Get all queued jobs for this service from the queue
	queuedJobs, err := self.jobQueue.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get jobs from queue: %w", err)
	}

	// Keep track of job IDs to mark as cancelled
	var jobIDsToCancel []uuid.UUID

	// 2. Remove matching jobs from the queue
	for _, item := range queuedJobs {
		if item.Data.ServiceID == serviceID {
			// Remove from queue
			if err := self.jobQueue.Remove(ctx, item.ID); err != nil {
				log.Errorf("Failed to remove job %s from queue: %v", item.ID, err)
			}
			idParsed, _ := uuid.Parse(item.ID)
			jobIDsToCancel = append(jobIDsToCancel, idParsed)
		}
	}

	// 3. Mark the jobs as cancelled in the database
	if len(jobIDsToCancel) > 0 {
		if err := self.repo.Deployment().MarkAsCancelled(ctx, jobIDsToCancel); err != nil {
			return fmt.Errorf("failed to mark jobs as cancelled: %w", err)
		}
	}

	return nil
}

// processJob processes a job from the queue
func (self *DeploymentController) processJob(ctx context.Context, item *queue.QueueItem[DeploymentJobRequest]) error {
	jobID, _ := uuid.Parse(item.ID)
	req := item.Data

	// Update the job status in the database
	err := self.repo.Deployment().MarkCancelled(ctx, req.ServiceID)
	if err != nil {
		log.Warnf("Failed to mark job as cancelled: %v service: %s", err, req.ServiceID)
	}
	_, err = self.repo.Deployment().MarkStarted(ctx, jobID)

	if err != nil {
		return fmt.Errorf("failed to mark job started: %w", err)
	}

	// Start the actual Kubernetes job
	k8sJobName, err := self.k8s.CreateDeployment(ctx, req.ServiceID.String(), jobID.String(), req.Environment)
	if err != nil {
		log.Error("Failed to create Kubernetes job", "err", err)

		// Update status to failed
		_, dbErr := self.repo.Deployment().MarkFailed(ctx, jobID, err.Error())

		if dbErr != nil {
			log.Error("Failed to update job failure status", "err", dbErr)
		}

		return err
	}

	// Update the Kubernetes job name in the database
	_, err = self.repo.Deployment().AssignKubernetesJobName(ctx, jobID, k8sJobName)

	if err != nil {
		log.Error("Failed to update Kubernetes job name in database", "err", err, "jobID", jobID, "k8sJobName", k8sJobName)
	}

	return nil
}

// SyncJobStatuses synchronizes the status of all processing jobs with Kubernetes
func (self *DeploymentController) SyncJobStatuses(ctx context.Context) error {
	// Get all job marked running status
	jobs, err := self.repo.Deployment().GetJobsByStatus(ctx, schema.DeploymentStatusRunning)
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
			_, err = self.repo.Deployment().MarkCompleted(ctx, job.ID)
		case k8s.JobFailed:
			_, err = self.repo.Deployment().MarkFailed(ctx, job.ID, k8sStatus.FailureReason)
		default:
			_, err = self.repo.Deployment().SetKubernetesJobStatus(ctx, job.ID, k8sStatus.ConditionType.String())
		}

		if err != nil {
			log.Error("Failed to update job status", "err", err, "jobID", job.ID)
		}
	}

	return nil
}
