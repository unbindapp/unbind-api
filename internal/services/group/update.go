package group_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
)

// UpdateGroup updates a group if the user has permission
func (self *GroupService) UpdateGroup(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, input GroupCreateInput) (*ent.Group, error) {
	// Get the existing group
	group, err := self.repo.Group().GetByID(ctx, groupID)
	if err != nil {
		// May be ent.NotFound
		return nil, err
	}

	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to manage system resources
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		// Has permission to manage groups
		{
			Action:       permission.ActionUpdate,
			ResourceType: permission.ResourceTypeGroup,
			ResourceID:   "*",
		},
		// has permission to update this specific group
		{
			Action:       permission.ActionUpdate,
			ResourceType: permission.ResourceTypeGroup,
			ResourceID:   groupID.String(),
		},
	}

	if group.TeamID != nil {
		// For team groups, check team-level permission
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       permission.ActionUpdate,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   group.TeamID.String(),
		})
	}

	// Execute permission checks
	if err := self.repo.Permissions().Check(ctx, userID, permissionChecks); err != nil {
		return nil, err
	}

	// Create update builder
	update := self.repo.Ent().Group.UpdateOneID(groupID)

	// Update fields if provided
	if input.Name != "" {
		update.SetName(input.Name)
	}

	if input.Description != "" {
		update.SetDescription(input.Description)
	}

	if input.IdentityProvider != "" {
		update.SetIdentityProvider(input.IdentityProvider)
	}

	if input.ExternalID != "" {
		update.SetExternalID(input.ExternalID)
	}

	// Cannot change team scope of an existing group
	if input.TeamID != nil && group.TeamID != nil && *input.TeamID != *group.TeamID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "cannot change team scope of an existing group")
	}

	// Cannot change global group to team-scoped or vice versa
	if (input.TeamID == nil && group.TeamID != nil) || (input.TeamID != nil && group.TeamID == nil) {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "cannot change group scope (global/team)")
	}

	// Execute the update
	updatedGroup, err := update.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errdefs.ErrGroupAlreadyExists
		}
		return nil, fmt.Errorf("error updating group: %w", err)
	}

	return updatedGroup, nil
}
