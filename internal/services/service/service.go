package service_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/infrastructure/prometheus"
	"github.com/unbindapp/unbind-api/internal/infrastructure/s3"
	"github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/internal/models"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	variables_service "github.com/unbindapp/unbind-api/internal/services/variables"
	webhooks_service "github.com/unbindapp/unbind-api/internal/services/webooks"
	"github.com/unbindapp/unbind-api/pkg/databases"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	"k8s.io/client-go/kubernetes"
)

// Integrate service management with internal permissions and kubernetes RBAC
type ServiceService struct {
	cfg                  *config.Config
	repo                 repositories.RepositoriesInterface
	githubClient         *github.GithubClient
	k8s                  *k8s.KubeClient
	deploymentController *deployctl.DeploymentController
	dbProvider           *databases.DatabaseProvider
	webhookService       *webhooks_service.WebhooksService
	variableService      *variables_service.VariablesService
	promClient           *prometheus.PrometheusClient
}

func NewServiceService(cfg *config.Config,
	repo repositories.RepositoriesInterface,
	githubClient *github.GithubClient,
	k8s *k8s.KubeClient,
	deploymentController *deployctl.DeploymentController,
	dbProvider *databases.DatabaseProvider,
	webhookService *webhooks_service.WebhooksService,
	variableService *variables_service.VariablesService,
	promClient *prometheus.PrometheusClient) *ServiceService {
	return &ServiceService{
		cfg:                  cfg,
		repo:                 repo,
		githubClient:         githubClient,
		k8s:                  k8s,
		deploymentController: deploymentController,
		dbProvider:           dbProvider,
		webhookService:       webhookService,
		variableService:      variableService,
		promClient:           promClient,
	}
}

func (self *ServiceService) VerifyInputs(ctx context.Context, teamID, projectID, environmentID uuid.UUID) (*ent.Environment, *ent.Project, error) {
	// Verify that the environment exists
	environment, err := self.repo.Environment().GetByID(ctx, environmentID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Environment not found")
		}
		return nil, nil, err
	}

	if environment == nil {
		return nil, nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Environment not found")
	}

	if environment.Edges.Project == nil {
		return nil, nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Environment does not belong to a project")
	}

	if environment.Edges.Project.ID != projectID {
		return nil, nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Environment does not belong to the specified project")
	}

	if environment.Edges.Project.Edges.Team == nil {
		return nil, nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Environment does not belong to a team")
	}

	if environment.Edges.Project.Edges.Team.ID != teamID {
		return nil, nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Project does not belong to the specified team")
	}

	return environment, environment.Edges.Project, nil
}

func (self *ServiceService) generateWildcardHost(ctx context.Context, tx repository.TxInterface, kubernetesName string, ports []schema.PortSpec) (*v1.HostSpec, error) {
	settings, err := self.repo.System().GetSystemSettings(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get system settings: %w", err)
	}

	if settings.WildcardBaseURL == nil || *settings.WildcardBaseURL == "" {
		return nil, nil // No wildcard base URL configured
	}

	domain, err := utils.GenerateSubdomain(kubernetesName, *settings.WildcardBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate subdomain: %w", err)
	}

	domainCount, err := self.repo.Service().CountDomainCollisons(ctx, tx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to count domain collisions: %w", err)
	}

	if domainCount > 0 {
		domain, err = utils.GenerateSubdomain(fmt.Sprintf("%s-%d", kubernetesName, domainCount), *settings.WildcardBaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to generate subdomain with suffix: %w", err)
		}
	}

	return &v1.HostSpec{
		Host: domain,
		Path: "/",
		Port: utils.ToPtr(ports[0].Port),
	}, nil
}

func (self *ServiceService) verifyS3Access(ctx context.Context, s3Source *ent.S3, bucket string, namespace string, client *kubernetes.Clientset) error {
	// Retrieve secret from kubernetes
	secret, err := self.k8s.GetSecret(ctx, s3Source.KubernetesSecret, namespace, client)
	if err != nil {
		return err
	}
	accessKeyId := string(secret.Data["access_key_id"])
	secretKey := string(secret.Data["secret_key"])
	if accessKeyId == "" || secretKey == "" {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
			"S3 source secret is missing access key or secret key")
	}

	s3Client, err := s3.NewS3Client(
		ctx,
		s3Source.Endpoint,
		s3Source.Region,
		accessKeyId,
		secretKey,
	)
	if err != nil {
		return err
	}

	// Probe the bucket
	err = s3Client.ProbeBucketRW(ctx, bucket)
	if err != nil {
		// s3 client already transforms into API handler compatible errors
		return err
	}

	return nil
}

func (self *ServiceService) validatePVC(ctx context.Context, teamID, projectID, environmentID uuid.UUID, name, namespace string, client *kubernetes.Clientset) error {
	isInUse, err := self.repo.Service().IsVolumeInUse(ctx, name)
	if err != nil {
		return err
	}
	if isInUse {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "PVC is already in use by another service")
	}

	// Get the actual PVC from k8s
	pvc, err := self.k8s.GetPersistentVolumeClaim(ctx, namespace, name, client)
	if err != nil {
		return err
	}
	if !pvc.IsAvailable {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "PVC is not available")
	}

	// Verify scope
	if pvc.TeamID != teamID {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "PVC not found")
	}

	if pvc.ProjectID != nil && *pvc.ProjectID != projectID {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "PVC not found")
	}

	if pvc.EnvironmentID != nil && *pvc.EnvironmentID != environmentID {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "PVC not found")
	}

	return nil
}

// Add volumes to service config
func (self *ServiceService) getVolumesForServices(ctx context.Context, namespace string, teamID uuid.UUID, services []*ent.Service) (map[uuid.UUID][]models.PVCInfo, error) {
	// Figure out all of the PVCs in the list and map them to services
	var pvcIDs []string
	serviceVolumes := make(map[uuid.UUID][]schema.ServiceVolume) // Map service ID to volumes

	for _, service := range services {
		var volumes []schema.ServiceVolume
		for _, vol := range service.Edges.ServiceConfig.Volumes {
			pvcIDs = append(pvcIDs, vol.ID)
			volumes = append(volumes, vol)
		}
		serviceVolumes[service.ID] = volumes
	}

	if len(pvcIDs) == 0 {
		return nil, nil // No PVCs to process
	}

	// Use channels and goroutines to parallelize the two operations
	type pvcResult struct {
		pvcs []models.PVCInfo
		err  error
	}

	type statsResult struct {
		stats map[string]*prometheus.PVCVolumeStats
		err   error
	}

	pvcChan := make(chan pvcResult, 1)
	statsChan := make(chan statsResult, 1)

	s := time.Now()

	// Get PVCs in parallel
	go func() {
		pvcs, err := self.k8s.ListPersistentVolumeClaims(ctx, namespace, map[string]string{
			"unbind-team": teamID.String(),
		}, self.k8s.GetInternalClient())
		pvcChan <- pvcResult{pvcs: pvcs, err: err}
	}()

	// Get prometheus stats in parallel
	go func() {
		var pvcStats map[string]*prometheus.PVCVolumeStats
		if len(pvcIDs) > 0 {
			stats, err := self.promClient.GetPVCsVolumeStats(ctx, pvcIDs, namespace, self.k8s.GetInternalClient())
			if err != nil {
				log.Errorf("Failed to get PVC stats from prometheus: %v", err)
				pvcStats = make(map[string]*prometheus.PVCVolumeStats) // Empty map so we can still proceed
			} else {
				pvcStats = make(map[string]*prometheus.PVCVolumeStats)
				for _, stat := range stats {
					pvcStats[stat.PVCName] = stat
				}
			}
		} else {
			pvcStats = make(map[string]*prometheus.PVCVolumeStats)
		}
		statsChan <- statsResult{stats: pvcStats, err: nil} // We handle errors above, so always return nil error
	}()

	// Wait for both operations to complete
	pvcRes := <-pvcChan
	statsRes := <-statsChan

	log.Infof("listPersistentVolumeClaims and getPVCsVolumeStats in parallel took %d ms", time.Since(s).Milliseconds())

	// Check for PVC error (this is the critical one)
	if pvcRes.err != nil {
		return nil, pvcRes.err
	}

	// Create a map of PVC ID to PVC for easy lookup
	pvcMap := make(map[string]models.PVCInfo)
	for _, pvc := range pvcRes.pvcs {
		pvcMap[pvc.ID] = pvc
	}

	// Build the result map
	result := make(map[uuid.UUID][]models.PVCInfo)
	for serviceID, volumes := range serviceVolumes {
		var servicePVCs []models.PVCInfo
		for _, volume := range volumes {
			if pvc, exists := pvcMap[volume.ID]; exists {
				// Add mount path from service config
				pvc.MountPath = &volume.MountPath

				// Add prometheus stats if available
				if stat, hasStats := statsRes.stats[pvc.ID]; hasStats {
					pvc.UsedGB = stat.UsedGB
					// CapacityGB should already be set from the PVC itself
				}
				servicePVCs = append(servicePVCs, pvc)
			}
		}
		result[serviceID] = servicePVCs
	}

	return result, nil
}
