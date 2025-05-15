package service_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
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
	return models.TransformServiceEntities(services), nil
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
	return models.TransformServiceEntity(service), nil
}
