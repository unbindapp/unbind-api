package servicegroup_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *ServiceGroupService) DeleteServiceGroup(ctx context.Context, requesterUserID uuid.UUID, input *models.DeleteServiceGroupInput) error {
	// Check permissions
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to manage teams
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeEnvironment,
			ResourceID:   input.EnvironmentID,
		},
	}

	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return err
	}

	// Verify inputs
	env, _, err := self.VerifyInputs(ctx, input.TeamID, input.ProjectID, input.EnvironmentID)
	if err != nil {
		return err
	}
	if env.ID != input.EnvironmentID {
		return errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service group not found")
	}

	return self.repo.ServiceGroup().Delete(ctx, input.ID)
}
