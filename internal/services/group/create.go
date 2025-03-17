package group_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
)

// Input for creating group
type GroupCreateInput struct {
	Name             string `validate:"required"`
	Description      string `validate:"required"`
	TeamID           *uuid.UUID
	IdentityProvider string
	ExternalID       string
	SuperuserGroup   bool
}

func (self *GroupService) CreateGroup(ctx context.Context, userID uuid.UUID, input *GroupCreateInput) (*ent.Group, error) {
	err := validate.Validator().Struct(input)
	if err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	// Creating a globally scoped group
	if input.TeamID == nil {
		if err := self.repo.Permissions().Check(
			ctx,
			userID,
			[]permissions_repo.PermissionCheck{
				// Has permission to manage system resources
				{
					Action:       permission.ActionManage,
					ResourceType: permission.ResourceTypeSystem,
					ResourceID:   "*",
				},
			},
		); err != nil {
			return nil, err
		}

		// ? Maybe have different scopes that can create global groups
	} else {
		// Verify the user is a member of the team
		isMember, err := self.repo.Team().HasUserWithID(ctx, *input.TeamID, userID)
		if err != nil {
			return nil, err
		}

		if !isMember {
			return nil, errdefs.ErrUnauthorized
		}

		if err := self.repo.Permissions().Check(
			ctx,
			userID,
			[]permissions_repo.PermissionCheck{
				// Has permission to manage team resources
				{
					Action:       permission.ActionManage,
					ResourceType: permission.ResourceTypeTeam,
					ResourceID:   input.TeamID.String(),
				},
			},
		); err != nil {
			return nil, err
		}
	}

	// Start builder
	groupCreate := self.repo.Ent().Group.Create().
		SetName(input.Name).
		SetDescription(input.Description).
		// ! TODO - we should probably have extra restrictions on making superuser groups
		SetSuperuser(input.SuperuserGroup).
		SetNillableTeamID(input.TeamID)

	if input.IdentityProvider != "" {
		groupCreate.SetIdentityProvider(input.IdentityProvider)
	}
	if input.ExternalID != "" {
		groupCreate.SetExternalID(input.ExternalID)
	}

	group, err := groupCreate.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errdefs.ErrGroupAlreadyExists
		}
		return nil, err
	}

	// Add creator as a member
	err = self.repo.Group().AddUser(ctx, group.ID, userID)
	if err != nil {
		return nil, err
	}

	return group, nil
}
