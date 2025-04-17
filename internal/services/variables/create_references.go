package variables_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *VariablesService) CreateVariableReference(ctx context.Context, requesterUserID uuid.UUID, input *models.CreateVariableReferenceInput) (*models.VariableReferenceResponse, error) {
	// ! TODO - we're going to need to change all of our permission checks to filter not reject
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   input.TargetServiceID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	if input.ValueTemplate != nil {
		if !utils.ContainsExactlyOneInterpolationMarker(*input.ValueTemplate) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Value template must contain exactly one interpolation marker")
		}
	}

	// Ensure service exists
	_, err := self.repo.Service().GetByID(ctx, input.TargetServiceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
		}
		return nil, err
	}

	resp, err := self.repo.Variables().CreateReference(ctx, input)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Variable reference already exists")
		}
		return nil, err
	}

	return models.TransformVariableReferenceResponseEntity(resp), nil
}
