package group_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
)

func (self *GroupService) AddUserToGroup(ctx context.Context, requesterUserID, targetUserID, groupID uuid.UUID) error {
	// Get the group to determine its type
	group, err := self.repo.Group().GetByID(ctx, groupID)
	if err != nil {
		// May be ent.NotFound
		return err
	}

	// Determine permissions based on group type
	if group.TeamID == nil {
		var isAdmin, hasGroupAdmin, hasGlobalGroupAdmin bool
		// Global group - requester must be system admin or group manager
		isAdmin, err = self.repo.Permissions().HasPermission(ctx, requesterUserID, permission.ActionManage, permission.ResourceTypeSystem, "*")
		if err != nil {
			return err
		}

		if !isAdmin {
			hasGroupAdmin, err = self.repo.Permissions().HasPermission(ctx, requesterUserID, permission.ActionManage, permission.ResourceTypeGroup, groupID.String())
			if err != nil {
				return err
			}
		}

		if !isAdmin && !hasGroupAdmin {
			hasGlobalGroupAdmin, err = self.repo.Permissions().HasPermission(ctx, requesterUserID, permission.ActionManage, permission.ResourceTypeGroup, "*")
			if err != nil {
				return err
			}
		}

		if !isAdmin && !hasGroupAdmin && !hasGlobalGroupAdmin {
			return errdefs.ErrUnauthorized
		}
	} else {
		var hasTeamAdmin, hasGroupAdmin bool
		// Team-scoped group - requester must be team admin or group manager
		hasTeamAdmin, err = self.repo.Permissions().HasPermission(ctx, requesterUserID, permission.ActionManage, permission.ResourceTypeTeam, group.TeamID.String())
		if err != nil {
			return err
		}

		if !hasTeamAdmin {
			hasGroupAdmin, err = self.repo.Permissions().HasPermission(ctx, requesterUserID, permission.ActionManage, permission.ResourceTypeGroup, groupID.String())
			if err != nil {
				return err
			}
		}

		if !hasTeamAdmin && !hasGroupAdmin {
			return errdefs.ErrUnauthorized
		}

		// Also verify the target user is a member of the team
		isTeamMember, err := self.repo.Team().HasUserWithID(ctx, *group.TeamID, targetUserID)
		if err != nil {
			return err
		}
		if !isTeamMember {
			return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "target user is not a member of the team")
		}
	}

	// Check if user is already in the group
	isMember, err := self.repo.Group().HasUserWithID(ctx, groupID, targetUserID)
	if err != nil {
		return err
	}

	if isMember {
		return nil // User is already a member, nothing to do
	}

	// Add the user to the group
	return self.repo.Group().AddUser(ctx, groupID, targetUserID)
}

// RemoveUserFromGroup removes a user from a group if the requester has permission
func (self *GroupService) RemoveUserFromGroup(ctx context.Context, requesterUserID, targetUserID, groupID uuid.UUID) error {
	// Get the group to determine its type
	group, err := self.repo.Group().GetByID(ctx, groupID)
	if err != nil {
		// May be ent.NotFound
		return err
	}

	// Users can remove themselves from groups
	if requesterUserID == targetUserID {
		return self.repo.Group().RemoveUser(ctx, groupID, targetUserID)
	}

	// Determine permissions based on group type
	if group.TeamID == nil {
		var isAdmin, hasGroupAdmin, hasGlobalGroupAdmin bool

		// Global group - requester must be system admin or group manager
		isAdmin, err = self.repo.Permissions().HasPermission(ctx, requesterUserID, permission.ActionManage, permission.ResourceTypeSystem, "*")
		if err != nil {
			return err
		}

		if !isAdmin {
			hasGroupAdmin, err = self.repo.Permissions().HasPermission(ctx, requesterUserID, permission.ActionManage, permission.ResourceTypeGroup, groupID.String())
			if err != nil {
				return err
			}
		}

		if !isAdmin && !hasGroupAdmin {
			hasGlobalGroupAdmin, err = self.repo.Permissions().HasPermission(ctx, requesterUserID, permission.ActionManage, permission.ResourceTypeGroup, "*")
			if err != nil {
				return err
			}
		}

		if !isAdmin && !hasGroupAdmin && !hasGlobalGroupAdmin {
			return errdefs.ErrUnauthorized
		}
	} else {
		var hasTeamAdmin, hasGroupAdmin bool

		// Team-scoped group - requester must be team admin or group manager
		hasTeamAdmin, err = self.repo.Permissions().HasPermission(ctx, requesterUserID, permission.ActionManage, permission.ResourceTypeTeam, group.TeamID.String())
		if err != nil {
			return err
		}

		if !hasTeamAdmin {
			hasGroupAdmin, err = self.repo.Permissions().HasPermission(ctx, requesterUserID, permission.ActionManage, permission.ResourceTypeGroup, groupID.String())
			if err != nil {
				return err
			}
		}

		if !hasTeamAdmin && !hasGroupAdmin {
			return errdefs.ErrUnauthorized
		}
	}

	// Check if user is actually in the group
	isMember, err := self.repo.Group().HasUserWithID(ctx, groupID, targetUserID)
	if err != nil {
		return err
	}

	if !isMember {
		return nil // User is not a member, nothing to do
	}

	// Remove the user from the group
	return self.repo.Group().RemoveUser(ctx, groupID, targetUserID)
}
