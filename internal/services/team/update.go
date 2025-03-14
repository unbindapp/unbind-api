package team_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repository/permissions"
	"github.com/unbindapp/unbind-api/internal/validate"
)

type TeamUpdateInput struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	DisplayName string    `json:"display_name" validate:"required"`
}

// UpdateTeam updates a specific team
func (self *TeamService) UpdateTeam(ctx context.Context, userID uuid.UUID, input *TeamUpdateInput) (*GetTeamResponse, error) {
	// Validate input
	err := validate.Validator().Struct(input)
	if err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to update system resources
		{
			Action:       permission.ActionUpdate,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		// Has permission to update teams
		{
			Action:       permission.ActionUpdate,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to update the specific team
		{
			Action:       permission.ActionUpdate,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   input.ID.String(),
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
	updatedTeam, err := self.repo.Team().Update(ctx, input.ID, input.DisplayName)
	if err != nil {
		// May be ent.NotFound
		return nil, err
	}

	return &GetTeamResponse{
		ID:          updatedTeam.ID,
		Name:        updatedTeam.Name,
		DisplayName: updatedTeam.DisplayName,
		CreatedAt:   updatedTeam.CreatedAt,
	}, nil
}
