package service_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/prometheus"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Get all services in an environment
func (self *ServiceService) GetServicesInEnvironment(ctx context.Context, requesterUserID uuid.UUID, teamID, projectID, environmentID uuid.UUID) ([]*models.ServiceResponse, error) {
	// Step 1: Get accessible service predicates for the user, scoped to the environmentID.
	servicePreds, err := self.repo.Permissions().GetAccessibleServicePredicates(ctx, requesterUserID, schema.ActionViewer, &environmentID)
	if err != nil {
		return nil, fmt.Errorf("error getting accessible service predicates: %w", err)
	}

	// Step 2: Verify parent inputs (team, project, environment) for integrity and clear error messages.
	// The VerifyInputs method already checks existence and relationships.
	_, _, err = self.VerifyInputs(ctx, teamID, projectID, environmentID)
	if err != nil {
		// VerifyInputs already returns specific errors like NotFound or InvalidInput
		return nil, err
	}

	// Step 3: Get services in environment, applying the permission predicate.
	// The GetByEnvironmentID repo method already filters by environmentID.
	// Pass true for withLatestDeployment as per original logic.
	services, err := self.repo.Service().GetByEnvironmentID(ctx, environmentID, servicePreds, true)
	if err != nil {
		return nil, fmt.Errorf("error fetching services for environment %s: %w", environmentID, err)
	}

	// Convert to response
	resp := models.TransformServiceEntities(services)

	// Figure out all of the PVCs in the list
	var pvcIDs []string
	for _, service := range resp {
		for _, vol := range service.Config.Volumes {
			pvcIDs = append(pvcIDs, vol.ID)
		}
	}

	// Query prometheus
	if len(pvcIDs) == 0 {
		return resp, nil
	}

	stats, err := self.promClient.GetPVCsVolumeStats(ctx, pvcIDs)
	if err != nil {
		log.Errorf("Failed to get PVC stats from prometheus: %v", err)
		return resp, nil
	}

	mapStats := make(map[string]prometheus.PVCVolumeStats)
	for _, stat := range stats {
		mapStats[stat.PVCName] = stat
	}

	// Add stats to the response
	for i := range resp {
		for j := range resp[i].Config.Volumes {
			if stat, ok := mapStats[resp[i].Config.Volumes[j].ID]; ok {
				resp[i].Config.Volumes[j].SizeGB = stat.CapacityGB
				resp[i].Config.Volumes[j].UsedGB = stat.UsedGB
			}
		}
	}

	// Return the response
	return resp, nil
}

// Get a service by ID
func (self *ServiceService) GetServiceByID(ctx context.Context, requesterUserID uuid.UUID, teamID, projectID, environmentID, serviceID uuid.UUID) (*models.ServiceResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to admin service
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   serviceID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Verify inputs
	_, _, err := self.VerifyInputs(ctx, teamID, projectID, environmentID)
	if err != nil {
		return nil, err
	}

	// Get services in environment
	service, err := self.repo.Service().GetByID(ctx, serviceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
		}
		return nil, err
	}

	// Convert to response
	resp := models.TransformServiceEntity(service)

	// Figure out all of the PVCs in the list
	var pvcIDs []string
	for _, vol := range resp.Config.Volumes {
		pvcIDs = append(pvcIDs, vol.ID)
	}

	// Query prometheus
	if len(pvcIDs) == 0 {
		return resp, nil
	}

	stats, err := self.promClient.GetPVCsVolumeStats(ctx, pvcIDs)
	if err != nil {
		log.Errorf("Failed to get PVC stats from prometheus: %v", err)
		return resp, nil
	}

	mapStats := make(map[string]prometheus.PVCVolumeStats)
	for _, stat := range stats {
		mapStats[stat.PVCName] = stat
	}

	// Add stats to the response
	for i := range resp.Config.Volumes {
		if stat, ok := mapStats[resp.Config.Volumes[i].ID]; ok {
			resp.Config.Volumes[i].SizeGB = stat.CapacityGB
			resp.Config.Volumes[i].UsedGB = stat.UsedGB
		}
	}

	// Return the response
	return resp, nil
}
