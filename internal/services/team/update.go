package team_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type TeamUpdateInput struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid4"`
	DisplayName string    `json:"display_name"`
	Description *string   `json:"description"`
}

// UpdateTeam updates a specific team
func (self *TeamService) UpdateTeam(ctx context.Context, userID uuid.UUID, input *TeamUpdateInput) (*models.TeamResponse, error) {
	// Validate input
	err := validate.Validator().Struct(input)
	if err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	if input.DisplayName == "" && input.Description == nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "No fields to update")
	}

	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to update team
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeTeam,
			ResourceID:   input.ID,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Update the team in the database
	updatedTeam, err := self.repo.Team().Update(ctx, input.ID, input.DisplayName, input.Description)
	if err != nil {
		// May be ent.NotFound
		return nil, err
	}

	return models.TransformTeamEntity(updatedTeam), nil
}
