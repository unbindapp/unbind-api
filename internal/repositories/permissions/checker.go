package permissions_repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/group"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

// PermissionCheck defines a permission check to be performed
type PermissionCheck struct {
	Action       permission.Action
	ResourceType permission.ResourceType
	ResourceID   string

	// If a custom check  is provided, it will be called in addition to the standard permission checks
	CustomCheck func(ctx context.Context, userID uuid.UUID) error
}

// Check if a user has any of the provided permissions. If any check passes, the permission is granted.
func (self *PermissionsRepository) Check(
	ctx context.Context,
	userID uuid.UUID,
	checks []PermissionCheck,
) error {
	checksPassed := false

	// Nothing to check, so permission is granted by default
	if len(checks) == 0 {
		checksPassed = true
	}

	for _, c := range checks {
		if c.Action != "" && c.ResourceType != "" && c.ResourceID != "" {
			hasPermission, err := self.HasPermission(ctx, userID, c.Action, c.ResourceType, c.ResourceID)
			if err != nil {
				return fmt.Errorf("error checking permission: %w", err)
			}

			if hasPermission {
				checksPassed = true
				break
			}
		}

		if c.CustomCheck != nil {
			if err := c.CustomCheck(ctx, userID); err != nil {
				log.Infof("custom permission check failed: %v", err)
				return errdefs.ErrUnauthorized
			}
		}
	}
	if !checksPassed {
		return errdefs.ErrUnauthorized
	}
	return nil
}

// HasPermission checks if a user has a specific permission on a resource
func (self *PermissionsRepository) HasPermission(
	ctx context.Context,
	userID uuid.UUID,
	action permission.Action,
	resourceType permission.ResourceType,
	resourceID string,
) (bool, error) {
	// Check super user
	isSuperuser, err := self.userRepo.IsSuperUser(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("error checking if user is superuser: %w", err)
	}

	if isSuperuser {
		return true, nil
	}

	// Get all groups the user belongs to
	userGroups, err := self.userRepo.GetGroups(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("error fetching user groups: %w", err)
	}

	// If user doesn't belong to any group, permission is denied
	if len(userGroups) == 0 {
		return false, nil
	}

	// Extract group IDs
	var groupIDs []uuid.UUID
	for _, g := range userGroups {
		groupIDs = append(groupIDs, g.ID)
	}

	// Check for permission match
	match, err := self.checkPermission(ctx, groupIDs, action, resourceType, resourceID)
	if err != nil {
		return false, err
	}

	if match.HasPermission {
		return true, nil
	}

	// Check for implied permissions (e.g., admin implies read, create, update, delete)
	impliedMatch, err := self.checkImpliedPermissions(ctx, groupIDs, action, resourceType, resourceID)
	if err != nil {
		return false, err
	}

	if impliedMatch.HasPermission {
		return true, nil
	}

	// Check for hierarchical permissions (e.g., team -> project)
	hierarchicalMatch, err := self.checkHierarchicalPermissions(ctx, groupIDs, action, resourceType, resourceID)
	if err != nil {
		return false, err
	}

	return hierarchicalMatch.HasPermission, nil
}

// PermissionResult represents the result of a permission check
type PermissionResult struct {
	HasPermission   bool
	PermissionFound bool
	Error           error
}

// checkPermission checks for an match of action, resource type, and resource ID
func (self *PermissionsRepository) checkPermission(
	ctx context.Context,
	groupIDs []uuid.UUID,
	action permission.Action,
	resourceType permission.ResourceType,
	resourceID string,
) (PermissionResult, error) {
	result := PermissionResult{}

	// Convert string action to permission.Action
	var permAction permission.Action
	switch action {
	case permission.ActionRead:
		permAction = permission.ActionRead
	case permission.ActionCreate:
		permAction = permission.ActionCreate
	case permission.ActionUpdate:
		permAction = permission.ActionUpdate
	case permission.ActionDelete:
		permAction = permission.ActionDelete
	case permission.ActionManage:
		permAction = permission.ActionManage
	case permission.ActionAdmin:
		permAction = permission.ActionAdmin
	case permission.ActionEdit:
		permAction = permission.ActionEdit
	case permission.ActionView:
		permAction = permission.ActionView
	default:
		// If action doesn't match any known action, return false
		return result, nil
	}

	// Check if any of the user's groups has the exact permission
	count, err := self.base.DB.Permission.Query().
		Where(
			permission.HasGroupsWith(group.IDIn(groupIDs...)),
			permission.ActionEQ(permAction),
			permission.ResourceTypeEQ(resourceType),
			permission.Or(
				permission.ResourceIDEQ(resourceID),
				permission.ResourceIDEQ("*"),
			),
		).
		Count(ctx)

	if err != nil {
		result.Error = fmt.Errorf("error checking exact permissions: %w", err)
		return result, result.Error
	}

	result.PermissionFound = count > 0
	result.HasPermission = count > 0
	return result, nil
}

// checkImpliedPermissions checks for permissions that imply the requested permission
// For example, "admin" implies "read", "create", "update", "delete", etc.
func (self *PermissionsRepository) checkImpliedPermissions(
	ctx context.Context,
	groupIDs []uuid.UUID,
	action permission.Action,
	resourceType permission.ResourceType,
	resourceID string,
) (PermissionResult, error) {
	result := PermissionResult{}

	// Define the implications
	var impliedActions []permission.Action

	switch action {
	case permission.ActionRead:
		// "read" is implied by these higher-level permissions
		impliedActions = []permission.Action{
			permission.ActionManage,
			permission.ActionAdmin,
			permission.ActionEdit,
			permission.ActionView,
		}
	case permission.ActionCreate, permission.ActionUpdate, permission.ActionDelete:
		// These are implied by manage and admin
		impliedActions = []permission.Action{
			permission.ActionManage,
			permission.ActionAdmin,
			permission.ActionEdit, // edit implies create and update but not delete
		}
		// Remove edit from implied actions if the action is delete
		if action == permission.ActionDelete {
			impliedActions = []permission.Action{
				permission.ActionManage,
				permission.ActionAdmin,
			}
		}
	case permission.ActionView:
		// view is a subset of read and is only used for K8s permissions
		impliedActions = []permission.Action{
			permission.ActionRead,
			permission.ActionManage,
			permission.ActionAdmin,
			permission.ActionEdit,
		}
	case permission.ActionEdit:
		// edit is implied by manage and admin
		impliedActions = []permission.Action{
			permission.ActionManage,
			permission.ActionAdmin,
		}
	case permission.ActionManage:
		// manage is implied by admin
		impliedActions = []permission.Action{
			permission.ActionAdmin,
		}
	default:
		// No implications for other actions
		return result, nil
	}

	// Check for exact resource ID
	count, err := self.base.DB.Permission.Query().
		Where(
			permission.HasGroupsWith(group.IDIn(groupIDs...)),
			permission.ActionIn(impliedActions...),
			permission.ResourceTypeEQ(resourceType),
			permission.Or(
				permission.ResourceIDEQ(resourceID),
				permission.ResourceIDEQ("*"),
			),
		).
		Count(ctx)

	if err != nil {
		result.Error = fmt.Errorf("error checking implied permissions: %w", err)
		return result, result.Error
	}

	result.PermissionFound = count > 0
	result.HasPermission = count > 0
	return result, nil
}

// checkHierarchicalPermissions checks for permissions in the resource hierarchy
// For example, permission on a team implies permission on its projects
func (self *PermissionsRepository) checkHierarchicalPermissions(
	ctx context.Context,
	groupIDs []uuid.UUID,
	action permission.Action,
	resourceType permission.ResourceType,
	resourceID string,
) (PermissionResult, error) {
	result := PermissionResult{}

	// Define hierarchical relationships
	switch resourceType {
	case permission.ResourceTypeProject:
		// Projects belong to teams
		// If we're checking project permission, also check if the user has the same permission on the parent team
		projectID, err := uuid.Parse(resourceID)
		if err != nil {
			return result, nil // Not a valid UUID, so no hierarchy to check
		}

		// Get the project and team
		project, err := self.projectRepo.GetByID(ctx, projectID)
		if err != nil {
			if ent.IsNotFound(err) {
				return result, nil // Project not found, no hierarchy to check
			}
			result.Error = fmt.Errorf("error fetching project: %w", err)
			return result, result.Error
		}

		// Check for permission on the team
		teamPermResult, err := self.checkPermission(ctx, groupIDs, action, permission.ResourceTypeTeam, project.Edges.Team.ID.String())
		if err != nil {
			result.Error = err
			return result, result.Error
		}

		if teamPermResult.HasPermission {
			result.HasPermission = true
			result.PermissionFound = true
			return result, nil
		}

	case "k8s":
		// Check if resourceID is a valid UUID (might be a team or project ID)
		resourceUUID, err := uuid.Parse(resourceID)
		if err == nil {
			// Try as team ID first
			team, err := self.teamRepo.GetByID(ctx, resourceUUID)
			if err == nil {
				// ResourceID is a team ID, check team permissions
				teamPermResult, err := self.checkPermission(ctx, groupIDs, action, permission.ResourceTypeTeam, team.ID.String())
				if err != nil {
					result.Error = err
					return result, result.Error
				}

				if teamPermResult.HasPermission {
					result.HasPermission = true
					result.PermissionFound = true
					return result, nil
				}
			}

			// Try as project ID
			project, err := self.projectRepo.GetByID(ctx, resourceUUID)
			if err == nil {
				// ResourceID is a project ID, check project permissions
				projectPermResult, err := self.checkPermission(ctx, groupIDs, action, permission.ResourceTypeProject, project.ID.String())
				if err != nil {
					result.Error = err
					return result, result.Error
				}

				if projectPermResult.HasPermission {
					result.HasPermission = true
					result.PermissionFound = true
					return result, nil
				}
			}
		}
	}

	// No hierarchical permission found
	return result, nil
}
