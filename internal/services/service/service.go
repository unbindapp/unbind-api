package service_service

import (
	"context"
	"fmt"

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
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	"github.com/unbindapp/unbind-api/internal/services/models"
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

// Add prom metrics to volume
func (self *ServiceService) addPromMetricsToServiceVolumes(ctx context.Context, services []*models.ServiceResponse) error {
	// Figure out all of the PVCs in the list
	var pvcIDs []string
	for _, service := range services {
		for _, vol := range service.Config.Volumes {
			pvcIDs = append(pvcIDs, vol.ID)
		}
	}

	// Query prometheus
	if len(pvcIDs) == 0 {
		return nil
	}

	stats, err := self.promClient.GetPVCsVolumeStats(ctx, pvcIDs)
	if err != nil {
		log.Errorf("Failed to get PVC stats from prometheus: %v", err)
		return nil
	}

	mapStats := make(map[string]prometheus.PVCVolumeStats)
	for _, stat := range stats {
		mapStats[stat.PVCName] = stat
	}

	// Add stats to the response
	for i := range services {
		for j := range services[i].Config.Volumes {
			if stat, ok := mapStats[services[i].Config.Volumes[j].ID]; ok {
				services[i].Config.Volumes[j].SizeGB = stat.CapacityGB
				services[i].Config.Volumes[j].UsedGB = stat.UsedGB
			}
		}
	}

	return nil
}
