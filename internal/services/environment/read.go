package environment_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Get a single environment by ID
func (self *EnvironmentService) GetEnvironmentByID(ctx context.Context, requesterUserID uuid.UUID, teamID uuid.UUID, projectID uuid.UUID, environmentID uuid.UUID) (*models.EnvironmentResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeEnvironment,
			ResourceID:   environmentID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Verify inputs
	_, environment, err := self.VerifyInputs(ctx, teamID, projectID, environmentID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	resp := models.TransformEnvironmentEntity(environment)

	// Summarizes services
	counts, providerSummaries, frameworkSummaries, err := self.repo.Service().SummarizeServices(ctx, []uuid.UUID{environmentID})
	if err != nil {
		return nil, err
	}
	resp.ServiceCount, _ = counts[environmentID]
	resp.FrameworkSummary, _ = frameworkSummaries[environmentID]
	resp.ProviderSummary, _ = providerSummaries[environmentID]

	return resp, nil
}
