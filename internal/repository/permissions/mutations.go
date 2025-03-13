package permissions_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
)

func (self *PermissionsRepository) Create(ctx context.Context, action permission.Action, resourceType permission.ResourceType, resourceID string, scope string) (*ent.Permission, error) {
	return self.base.DB.Permission.Create().
		SetAction(action).
		SetResourceType(resourceType).
		SetResourceID(resourceID).
		SetScope(scope).
		Save(ctx)
}

func (self *PermissionsRepository) AddToGroup(ctx context.Context, groupID, permissionID uuid.UUID) error {
	_, err := self.base.DB.Group.UpdateOneID(groupID).
		AddPermissionIDs(permissionID).
		Save(ctx)
	return err
}

func (self *PermissionsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return self.base.DB.Permission.DeleteOneID(id).Exec(ctx)
}

func (self *PermissionsRepository) RemoveFromGroup(ctx context.Context, groupID, permissionID uuid.UUID) error {
	_, err := self.base.DB.Group.
		UpdateOneID(groupID).
		RemovePermissionIDs(permissionID).
		Save(ctx)
	return err
}
