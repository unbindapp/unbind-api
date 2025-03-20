package service_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Get all services in an environment
func (self *ServiceService) GetServicesInEnvironment(ctx context.Context, requesterUserID uuid.UUID, teamID, projectID, environmentID uuid.UUID) ([]*models.ServiceResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to read system resources
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		// Has permission to read teams
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to read the specific team
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   teamID.String(),
		},
		// Has permission to read the specific project
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   projectID.String(),
		},
		// Has permission to read the specific environment
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   environmentID.String(),
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
	services, err := self.repo.Service().GetByEnvironmentID(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	return models.TransformServiceEntities(services), nil
}

// Get a service by ID
func (self *ServiceService) GetServiceByID(ctx context.Context, requesterUserID uuid.UUID, teamID, projectID, environmentID, serviceID uuid.UUID) (*models.ServiceResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to read system resources
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		// Has permission to read teams
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to read the specific team
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   teamID.String(),
		},
		// Has permission to read the specific project
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   projectID.String(),
		},
		// Has permission to read the specific environment
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   environmentID.String(),
		},
		// Has permission to read the specific service
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeService,
			ResourceID:   serviceID.String(),
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
	service, err := self.repo.Service().GetByID(ctx, environmentID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
		}
		return nil, err
	}

	// Convert to response
	return models.TransformServiceEntity(service), nil
}
