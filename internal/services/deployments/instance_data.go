package deployments_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/models"
)

// ServiceInstanceData holds instance data for a service
type ServiceInstanceData struct {
	ServiceID       uuid.UUID
	Status          schema.DeploymentStatus
	InstanceEvents  []models.EventRecord
	CrashingReasons []string
}

// AttachInstanceDataToServices efficiently attaches instance data to multiple services in an environment
// This makes a single Kubernetes call per environment instead of per service
func (self *DeploymentService) AttachInstanceDataToServices(ctx context.Context, services []*ent.Service, namespace string) (map[uuid.UUID]*ServiceInstanceData, error) {
	if len(services) == 0 {
		return make(map[uuid.UUID]*ServiceInstanceData), nil
	}

	// Get all pod statuses for the environment in a single call
	statuses, err := self.k8s.GetPodContainerStatusByLabels(
		ctx,
		namespace,
		map[string]string{
			"unbind-environment": services[0].EnvironmentID.String(),
		},
		self.k8s.GetInternalClient(),
	)
	if err != nil {
		log.Error("Error getting pod container status for environment", "err", err, "environment_id", services[0].EnvironmentID)
		return nil, err
	}

	// Group statuses by service ID
	serviceStatuses := make(map[uuid.UUID][]k8s.PodContainerStatus)
	for _, status := range statuses {
		// The ServiceID is already parsed and stored in the PodContainerStatus struct
		if status.ServiceID != uuid.Nil {
			serviceStatuses[status.ServiceID] = append(serviceStatuses[status.ServiceID], status)
		}
	}

	// Calculate instance data for each service
	result := make(map[uuid.UUID]*ServiceInstanceData)
	for _, service := range services {
		if service.Edges.CurrentDeployment == nil || service.Edges.ServiceConfig == nil {
			continue
		}

		statuses := serviceStatuses[service.ID]
		instanceData := self.calculateInstanceData(statuses, service.Edges.ServiceConfig.Replicas)
		result[service.ID] = instanceData
	}

	return result, nil
}

// calculateInstanceData processes pod statuses to determine deployment status and events
func (self *DeploymentService) calculateInstanceData(statuses []k8s.PodContainerStatus, expectedReplicas int32) *ServiceInstanceData {
	// Initialize counters
	pendingCount := 0
	failedCount := 0
	runningCount := int32(0)
	crashingCount := 0
	unknownCount := 0
	events := []models.EventRecord{}
	crashingReasons := []string{}

	// Process each pod status
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

		// Process container instances
		for _, instance := range status.Instances {
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
	targetStatus := schema.DeploymentStatusWaiting
	if failedCount > 0 {
		targetStatus = schema.DeploymentStatusCrashing
	} else if crashingCount > 0 {
		targetStatus = schema.DeploymentStatusCrashing
	} else if pendingCount > 0 || unknownCount > 0 {
		targetStatus = schema.DeploymentStatusWaiting
	} else if runningCount >= expectedReplicas {
		targetStatus = schema.DeploymentStatusActive
	}

	return &ServiceInstanceData{
		Status:          targetStatus,
		InstanceEvents:  events,
		CrashingReasons: crashingReasons,
	}
}

// AttachInstanceDataToDeploymentResponses attaches instance data to deployment responses
func (self *DeploymentService) AttachInstanceDataToDeploymentResponses(deployments []*models.DeploymentResponse, instanceData *ServiceInstanceData, currentDeploymentID uuid.UUID) {
	if instanceData == nil {
		return
	}

	for i := range deployments {
		if deployments[i].ID == currentDeploymentID {
			deployments[i].Status = instanceData.Status
			deployments[i].InstanceEvents = instanceData.InstanceEvents
			deployments[i].CrashingReasons = instanceData.CrashingReasons
		} else {
			if deployments[i].Status == schema.DeploymentStatusBuildSucceeded {
				deployments[i].Status = schema.DeploymentStatusRemoved
			}
		}
	}
}

// AttachInstanceDataToServiceResponse attaches instance data to a single service response
// Returns the instance data that was attached, or nil if no data was available
func (self *DeploymentService) AttachInstanceDataToServiceResponse(service *models.ServiceResponse, instanceDataMap map[uuid.UUID]*ServiceInstanceData) *ServiceInstanceData {
	instanceData := instanceDataMap[service.ID]
	if instanceData == nil {
		return nil
	}

	// Attach to current deployment
	if service.CurrentDeployment != nil {
		service.CurrentDeployment.Status = instanceData.Status
		service.CurrentDeployment.InstanceEvents = instanceData.InstanceEvents
		service.CurrentDeployment.CrashingReasons = instanceData.CrashingReasons
	}

	// Attach to last deployment if it's the current one
	if service.LastDeployment != nil && service.CurrentDeployment != nil &&
		service.LastDeployment.ID == service.CurrentDeployment.ID {
		service.LastDeployment.Status = instanceData.Status
		service.LastDeployment.InstanceEvents = instanceData.InstanceEvents
		service.LastDeployment.CrashingReasons = instanceData.CrashingReasons
	}

	return instanceData
}

// AttachInstanceDataToServiceResponses attaches instance data to multiple service responses
// This is a convenience function that calls AttachInstanceDataToServiceResponse for each service
func (self *DeploymentService) AttachInstanceDataToServiceResponses(services []*models.ServiceResponse, instanceDataMap map[uuid.UUID]*ServiceInstanceData) {
	for _, service := range services {
		self.AttachInstanceDataToServiceResponse(service, instanceDataMap)
	}
}
