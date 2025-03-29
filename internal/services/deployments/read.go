package deployments_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
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

	_, err := self.validateInputs(ctx, input)
	if err != nil {
		return nil, nil, nil, err
	}

	// Huma doesn't support pointers in query so we need to convert zero's to nil
	var cursor *time.Time
	if !input.Cursor.IsZero() {
		cursor = &input.Cursor
	}
	// Get build jobs
	deployments, nextCursor, err := self.repo.Deployment().GetByServiceIDPaginated(ctx, input.ServiceID, cursor, input.Statuses)
	if err != nil {
		return nil, nil, nil, err
	}

	service, err := self.repo.Service().GetByID(ctx, input.ServiceID)

	if err != nil {
		return nil, nil, nil, err
	}

	// Transform response
	resp := models.TransformDeploymentEntities(deployments)
	currentDeployment := models.TransformDeploymentEntity(service.Edges.CurrentDeployment)

	// Get pagination metadata
	metadata := &models.PaginationResponseMetadata{
		HasNext:        nextCursor != nil,
		NextCursor:     nextCursor,
		PreviousCursor: cursor,
	}

	return resp, currentDeployment, metadata, nil
}
