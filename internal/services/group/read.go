package group_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repository/permissions"
)

// GetGroupByID retrieves a group by its ID
func (self *GroupService) GetGroupByID(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) (*ent.Group, error) {
	// Retrieve the group
	group, err := self.repo.Group().GetByID(ctx, groupID)
	if err != nil {
		// May be entNotFound err
		return nil, err
	}

	// Check if the user is a member, in which case they should always have permission
	isMember, err := self.repo.Group().HasUserWithID(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		permissionChecks := []permissions_repo.PermissionCheck{
			// Has permission to read system resources
			{
				Action:       permission.ActionRead,
				ResourceType: permission.ResourceTypeSystem,
				ResourceID:   "*",
			},
			// Has permission to read groups
			{
				Action:       permission.ActionRead,
				ResourceType: permission.ResourceTypeGroup,
				ResourceID:   "*",
			},
			// Has permission to read this specific group
			{
				Action:       permission.ActionRead,
				ResourceType: permission.ResourceTypeGroup,
				ResourceID:   groupID.String(),
			},
		}
		if group.TeamID != nil {
			// For team groups, check team-level permission
			permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
				Action:       permission.ActionRead,
				ResourceType: permission.ResourceTypeTeam,
				ResourceID:   group.TeamID.String(),
			})
		}

		if err := self.repo.Permissions().Check(
			ctx,
			userID,
			permissionChecks,
		); err != nil {
			return nil, err
		}
	}

	return group, nil
}

// ListGroups retrieves all groups the user has permission to view
func (self *GroupService) ListGroups(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID) ([]*ent.Group, error) {
	// Start with a base query
	query := self.repo.Ent().Group.Query()

	// Check if user is a member of the team
	isMember, err := self.repo.Team().HasUserWithID(ctx, *teamID, userID)
	if err != nil {
		return nil, err
	}

	// Always give members access to their own team's groups
	if !isMember {
		permissionChecks := []permissions_repo.PermissionCheck{
			// Has permission to read system resources
			{
				Action:       permission.ActionRead,
				ResourceType: permission.ResourceTypeSystem,
				ResourceID:   "*",
			},
			// Has permission to read groups
			{
				Action:       permission.ActionRead,
				ResourceType: permission.ResourceTypeGroup,
				ResourceID:   "*",
			},
		}
		if teamID != nil {
			// For team groups, check team-level permission
			permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
				Action:       permission.ActionRead,
				ResourceType: permission.ResourceTypeTeam,
				ResourceID:   teamID.String(),
			})
		}

		if err := self.repo.Permissions().Check(
			ctx,
			userID,
			permissionChecks,
		); err != nil {
			return nil, err
		}
	}

	// Execute the query
	return query.All(ctx)
}
