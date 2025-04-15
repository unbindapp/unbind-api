package instance_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Get kubernetes container statuses for a service
func (self *InstanceService) GetInstanceStatuses(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.InstanceStatusInput) ([]k8s.PodContainerStatus, error) {
	team, project, environment, service, err := self.validatePermissionsAndParseInputs(ctx, requesterUserID, input.Type, input.TeamID, input.ProjectID, input.EnvironmentID, input.ServiceID)
	if err != nil {
		return nil, err
	}

	// Determine labels
	labels := make(map[string]string)
	switch input.Type {
	case models.InstanceTypeService:
		labels["unbind-service"] = service.ID.String()
	case models.InstanceTypeEnvironment:
		labels["unbind-environment"] = environment.ID.String()
	case models.InstanceTypeProject:
		labels["unbind-project"] = project.ID.String()
	case models.InstanceTypeTeam:
		labels["unbind-team"] = team.ID.String()
	default:
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Invalid instance type")
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	return self.k8s.GetPodContainerStatusByLabels(
		ctx,
		project.Edges.Team.Namespace,
		labels,
		client,
	)
}
