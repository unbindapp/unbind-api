package deployments_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/permission"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *DeploymentService) GetDeploymentsForService(ctx context.Context, requesterUserId uuid.UUID, input *models.GetDeploymentsInput) ([]*models.DeploymentResponse, *models.PaginationResponseMetadata, error) {
	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserId, []permissions_repo.PermissionCheck{
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   input.TeamID.String(),
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   "*",
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   input.ProjectID.String(),
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeService,
			ResourceID:   input.ServiceID.String(),
		},
	}); err != nil {
		return nil, nil, err
	}

	if err := self.validateInputs(ctx, input); err != nil {
		return nil, nil, err
	}

	// Huma doesn't support pointers in query so we need to convert zero's to nil
	var cursor *time.Time
	if !input.Cursor.IsZero() {
		cursor = &input.Cursor
	}
	// Get build jobs
	deployments, nextCursor, err := self.repo.Deployment().GetByServiceIDPaginated(ctx, input.ServiceID, cursor, input.Statuses)
	if err != nil {
		return nil, nil, err
	}

	// Transform response
	resp := models.TransformDeploymentEntities(deployments)

	// Get pagination metadata
	metadata := &models.PaginationResponseMetadata{
		HasNext:        nextCursor != nil,
		NextCursor:     nextCursor,
		PreviousCursor: cursor,
	}

	return resp, metadata, nil

}
