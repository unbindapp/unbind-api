package deployments_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/models"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
)

func (self *DeploymentService) GetDeploymentsForService(ctx context.Context, requesterUserId uuid.UUID, input *models.GetDeploymentsInput) ([]*models.DeploymentResponse, *models.DeploymentResponse, *models.PaginationResponseMetadata, error) {
	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserId, []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   input.ServiceID,
		},
	}); err != nil {
		return nil, nil, nil, err
	}

	service, err := self.validateInputs(ctx, input)
	if err != nil {
		return nil, nil, nil, err
	}

	// Huma doesn't support pointers in query so we need to convert zero's to nil
	var cursor *time.Time
	if !input.Cursor.IsZero() {
		cursor = &input.Cursor
	}
	// Get build jobs
	deployments, nextCursor, err := self.repo.Deployment().GetByServiceIDPaginated(ctx, input.ServiceID, input.PerPage, cursor, input.Statuses)
	if err != nil {
		return nil, nil, nil, err
	}

	if err != nil {
		return nil, nil, nil, err
	}

	// Transform response
	resp := models.TransformDeploymentEntities(deployments)

	currentDeployment, err := self.attachInstanceDataToCurrent(ctx, resp, service)
	if err != nil {
		log.Error("Error attaching instance data to current deployment", "err", err, "service_id", service.ID)
		return nil, nil, nil, err
	}

	// Get pagination metadata
	metadata := &models.PaginationResponseMetadata{
		HasNext:        nextCursor != nil,
		NextCursor:     nextCursor,
		PreviousCursor: cursor,
	}

	return resp, currentDeployment, metadata, nil
}

func (self *DeploymentService) GetDeploymentByID(ctx context.Context, requesterUserId uuid.UUID, input *models.GetDeploymentByIDInput) (*models.DeploymentResponse, error) {
	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserId, []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   input.ServiceID,
		},
	}); err != nil {
		return nil, err
	}

	service, err := self.validateInputs(ctx, input)
	if err != nil {
		return nil, err
	}

	// Get deployment
	deployment, err := self.repo.Deployment().GetByID(ctx, input.DeploymentID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, input.DeploymentID.String())
		}
		return nil, err
	}

	if service.CurrentDeploymentID != nil && *service.CurrentDeploymentID == deployment.ID {
		// If this is the current deployment, attach instance data
		return self.attachInstanceDataToCurrent(ctx, []*models.DeploymentResponse{models.TransformDeploymentEntity(deployment)}, service)
	}

	// Get build jobs
	return models.TransformDeploymentEntity(deployment), nil
}

// Attach instance data
func (self *DeploymentService) attachInstanceDataToCurrent(ctx context.Context, deployments []*models.DeploymentResponse, service *ent.Service) (*models.DeploymentResponse, error) {
	if service.Edges.CurrentDeployment == nil {
		return nil, nil
	}
	targetDeployment := models.TransformDeploymentEntity(service.Edges.CurrentDeployment)

	namespace := service.Edges.Environment.Edges.Project.Edges.Team.Namespace
	statuses, err := self.k8s.GetPodContainerStatusByLabels(
		ctx,
		namespace,
		map[string]string{
			"unbind-service": service.ID.String(),
		},
		self.k8s.GetInternalClient(),
	)
	if err != nil {
		log.Error("Error getting pod container status", "err", err, "service_id", service.ID)
		return targetDeployment, err
	}

	// Get target status
	targetStatus := schema.DeploymentStatusWaiting

	// Get expected replicas
	expectedContainerCount := service.Edges.ServiceConfig.Replicas

	pendingCount := 0
	failedCount := 0
	runningCount := int32(0)
	crashingCount := 0
	unknownCount := 0
	restarts := int32(0)
	events := []models.EventRecord{}
	crashingReasons := []string{}
	for _, status := range statuses {
		switch status.Phase {
		case k8s.PodPending:
			pendingCount++
		case k8s.PodFailed:
			failedCount++
		case k8s.PodRunning:
			runningCount++
		default:
			unknownCount++
		}

		for _, instance := range status.Instances {
			restarts += instance.RestartCount
			if instance.IsCrashing {
				crashingCount++
				crashingReasons = append(crashingReasons, instance.CrashLoopReason)
				switch instance.State {
				case k8s.ContainerStateRunning:
					runningCount++
				case k8s.ContainerStateWaiting:
					pendingCount++
				}

				events = append(events, instance.Events...)
			}
		}
	}

	// Determine target status based on counts
	if failedCount > 0 {
		targetStatus = schema.DeploymentStatusCrashing
	} else if crashingCount > 0 {
		targetStatus = schema.DeploymentStatusCrashing
	} else if pendingCount > 0 || unknownCount > 0 {
		targetStatus = schema.DeploymentStatusWaiting
	} else if runningCount >= expectedContainerCount {
		targetStatus = schema.DeploymentStatusActive
	}

	// Attach data to deployment
	for i := range deployments {
		if deployments[i].ID == targetDeployment.ID {
			deployments[i].Status = targetStatus
			targetDeployment.Status = targetStatus
			deployments[i].InstanceEvents = events
			targetDeployment.InstanceEvents = events
			deployments[i].CrashingReasons = crashingReasons
			targetDeployment.CrashingReasons = crashingReasons
			break
		}
	}

	return targetDeployment, nil
}
