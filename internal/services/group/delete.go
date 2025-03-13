package group_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/group"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/log"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repository/permissions"
)

// DeleteGroup deletes a group if the user has permissions
func (self *GroupService) DeleteGroup(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) error {
	// Get the group to determine its type
	groupToDelete, err := self.repo.Group().GetByID(ctx, groupID)
	if err != nil {
		// May be ent.NotFound
		return err
	}

	// Determine permissions based on group type (global vs team)
	if groupToDelete.TeamID == nil {
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
				// Has permission to manage groups
				{
					Action:       permission.ActionManage,
					ResourceType: permission.ResourceTypeGroup,
					ResourceID:   "*",
				},
			},
		); err != nil {
			return err
		}

	} else {
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
				// Has permission to manage team
				{
					Action:       permission.ActionManage,
					ResourceType: permission.ResourceTypeTeam,
					ResourceID:   groupToDelete.TeamID.String(),
				},
				// Has permission to manage groups
				{
					Action:       permission.ActionManage,
					ResourceType: permission.ResourceTypeGroup,
					ResourceID:   "*",
				},
			},
		); err != nil {
			return err
		}
	}

	// If the group has K8s RBAC, clean it up first
	if groupToDelete.K8sRoleName != "" {
		if err := self.rbacManager.DeleteK8sRBAC(ctx, groupID); err != nil {
			log.Warnf("Error cleaning up K8s RBAC: %v", err)
		}
	}

	// Delete the group and all its associations
	_, err = self.repo.Ent().Group.Delete().Where(group.ID(groupID)).Exec(ctx)
	return err
}
