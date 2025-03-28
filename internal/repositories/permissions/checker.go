package permissions_repo

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/group"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/ent/predicate"
	entSchema "github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
)

// PermissionCheck defines a permission check to be performed
type PermissionCheck struct {
	Action       entSchema.PermittedAction
	ResourceType entSchema.ResourceType
	ResourceID   uuid.UUID
	// If a custom check is provided, it will be called in addition to the standard permission checks
	CustomCheck func(ctx context.Context, userID uuid.UUID) error
}

// Check if a user has any of the provided permissions. If any check passes, the permission is granted.
func (self *PermissionsRepository) Check(
	ctx context.Context,
	userID uuid.UUID,
	checks []PermissionCheck,
) error {
	// Nothing to check, so permission is granted by default
	if len(checks) == 0 {
		return nil
	}

	// Get all groups the user belongs to
	userGroups, err := self.userRepo.GetGroups(ctx, userID)
	if err != nil {
		return fmt.Errorf("error fetching user groups: %w", err)
	}

	// If user doesn't belong to any group, permission is denied
	if len(userGroups) == 0 {
		return errdefs.ErrUnauthorized
	}

	// Extract group IDs
	var groupIDs []uuid.UUID
	for _, g := range userGroups {
		groupIDs = append(groupIDs, g.ID)
	}

	// Check each permission
	for _, check := range checks {
		// Run custom check if provided
		if check.CustomCheck != nil {
			if err := check.CustomCheck(ctx, userID); err == nil {
				return nil // Custom check passed
			}
		}

		// Skip checks with empty values
		if check.Action == "" || check.ResourceType == "" {
			continue
		}

		// Build a comprehensive permission check query
		hasAccess, err := self.checkComprehensivePermission(ctx, groupIDs, check.Action, check.ResourceType, check.ResourceID)
		if err != nil {
			return fmt.Errorf("error checking permission: %w", err)
		}

		if hasAccess {
			return nil // Permission granted
		}
	}

	// No permission check passed
	return errdefs.ErrUnauthorized
}

// checkComprehensivePermission performs a comprehensive permission check including implied permissions and hierarchy
func (self *PermissionsRepository) checkComprehensivePermission(
	ctx context.Context,
	groupIDs []uuid.UUID,
	action entSchema.PermittedAction,
	resourceType entSchema.ResourceType,
	resourceID uuid.UUID,
) (bool, error) {
	// Get the resource hierarchy information
	// This loads parent IDs and types that would grant permissions to this resource
	hierarchyInfo, err := self.getResourceHierarchy(ctx, resourceType, resourceID)
	if err != nil {
		return false, fmt.Errorf("error getting resource hierarchy: %w", err)
	}

	// ! If resource ID is nil just set it to a random value so it won't match ever
	if resourceID == uuid.Nil {
		resourceID = uuid.New()
	}

	// Determine which actions would satisfy this permission check
	impliedActions := self.getImpliedActions(action)

	// Build a query that handles:
	// 1. Direct resource match
	// 2. Resource hierarchy matches
	// 3. Implied permission matches
	// All in a single database query

	// Base query that handles user's groups
	query := self.base.DB.Permission.Query().
		Where(permission.HasGroupsWith(group.IDIn(groupIDs...)))

	// Create a slice of predicates for the different permission checks
	var predicates []predicate.Permission

	// 1. Direct resource match
	directMatch := []predicate.Permission{
		permission.ResourceTypeEQ(resourceType),
		permission.ActionIn(impliedActions...),
	}

	// Build a predicate for the resource selector
	resourcePredicates := []predicate.Permission{
		// Access to specific resource ID
		func(s *sql.Selector) {
			sqljson.ValueEQ(permission.FieldResourceSelector, resourceID.String(), sqljson.Path("id"))
		},
		// Or superuser access
		func(s *sql.Selector) {
			sqljson.ValueEQ(permission.FieldResourceSelector, true, sqljson.Path("superuser"))
		},
	}

	// Add direct match with resource selector predicates
	directResourceMatch := append(directMatch, permission.Or(resourcePredicates...))
	predicates = append(predicates, permission.And(directResourceMatch...))

	// 2. Add hierarchy matches if available
	for _, hierInfo := range hierarchyInfo {
		// For each parent in the hierarchy
		hierMatch := []predicate.Permission{
			permission.ResourceTypeEQ(hierInfo.ResourceType),
			permission.ActionIn(impliedActions...),
		}

		// Build predicate for the parent resource
		hierResourcePredicates := []predicate.Permission{
			// Access to specific parent resource ID
			func(s *sql.Selector) {
				sqljson.ValueEQ(permission.FieldResourceSelector, hierInfo.ResourceID.String(), sqljson.Path("id"))
			},
			// Or superuser access to the parent resource type
			func(s *sql.Selector) {
				sqljson.ValueEQ(permission.FieldResourceSelector, true, sqljson.Path("superuser"))
			},
		}

		// Add hierarchy match with resource selector predicates
		hierResourceMatch := append(hierMatch, permission.Or(hierResourcePredicates...))
		predicates = append(predicates, permission.And(hierResourceMatch...))
	}

	// Add all predicates to the query using OR
	query = query.Where(permission.Or(predicates...))

	// Execute the query
	count, err := query.Count(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking permissions: %w", err)
	}

	return count > 0, nil
}

// ResourceHierarchyInfo represents a parent resource in the hierarchy
type ResourceHierarchyInfo struct {
	ResourceType entSchema.ResourceType
	ResourceID   uuid.UUID
}

// getResourceHierarchy returns the hierarchy of resources that would grant permissions to this resource
func (self *PermissionsRepository) getResourceHierarchy(
	ctx context.Context,
	resourceType entSchema.ResourceType,
	resourceID uuid.UUID,
) ([]ResourceHierarchyInfo, error) {
	var hierarchy []ResourceHierarchyInfo

	switch resourceType {
	case entSchema.ResourceTypeProject:
		// Projects belong to teams
		project, err := self.projectRepo.GetByID(ctx, resourceID)
		if err != nil {
			if ent.IsNotFound(err) {
				return hierarchy, nil
			}
			return nil, fmt.Errorf("error fetching project: %w", err)
		}

		// Add team to hierarchy if the project has a team
		if project.Edges.Team != nil {
			hierarchy = append(hierarchy, ResourceHierarchyInfo{
				ResourceType: entSchema.ResourceTypeTeam,
				ResourceID:   project.Edges.Team.ID,
			})
		}

	case entSchema.ResourceTypeEnvironment:
		// Environments belong to projects, which belong to teams
		env, err := self.environmentRepo.GetByID(ctx, resourceID)
		if err != nil {
			if ent.IsNotFound(err) {
				return hierarchy, nil
			}
			return nil, fmt.Errorf("error fetching environment: %w", err)
		}

		// Add project to hierarchy if the environment has a project
		if env.Edges.Project != nil {
			hierarchy = append(hierarchy, ResourceHierarchyInfo{
				ResourceType: entSchema.ResourceTypeProject,
				ResourceID:   env.Edges.Project.ID,
			})

			// Add team to hierarchy if the project has a team
			if env.Edges.Project.Edges.Team != nil {
				hierarchy = append(hierarchy, ResourceHierarchyInfo{
					ResourceType: entSchema.ResourceTypeTeam,
					ResourceID:   env.Edges.Project.Edges.Team.ID,
				})
			}
		}

	case entSchema.ResourceTypeService:
		// Services belong to environments, which belong to projects, which belong to teams
		service, err := self.serviceRepo.GetByID(ctx, resourceID)
		if err != nil {
			if ent.IsNotFound(err) {
				return hierarchy, nil
			}
			return nil, fmt.Errorf("error fetching service: %w", err)
		}

		// Add environment to hierarchy if the service has an environment
		if service.Edges.Environment != nil {
			hierarchy = append(hierarchy, ResourceHierarchyInfo{
				ResourceType: entSchema.ResourceTypeEnvironment,
				ResourceID:   service.Edges.Environment.ID,
			})

			// Add project to hierarchy if the environment has a project
			if service.Edges.Environment.Edges.Project != nil {
				hierarchy = append(hierarchy, ResourceHierarchyInfo{
					ResourceType: entSchema.ResourceTypeProject,
					ResourceID:   service.Edges.Environment.Edges.Project.ID,
				})

				// Add team to hierarchy if the project has a team
				if service.Edges.Environment.Edges.Project.Edges.Team != nil {
					hierarchy = append(hierarchy, ResourceHierarchyInfo{
						ResourceType: entSchema.ResourceTypeTeam,
						ResourceID:   service.Edges.Environment.Edges.Project.Edges.Team.ID,
					})
				}
			}
		}
	}

	return hierarchy, nil
}

// getImpliedActions returns a list of actions that would satisfy the requested action
func (self *PermissionsRepository) getImpliedActions(action entSchema.PermittedAction) []entSchema.PermittedAction {
	switch action {
	case entSchema.ActionViewer:
		// Viewer permissions are implied by all other permission levels
		return []entSchema.PermittedAction{
			entSchema.ActionViewer,
			entSchema.ActionEditor,
			entSchema.ActionAdmin,
		}
	case entSchema.ActionEditor:
		// Editor permissions are implied by Editor and Admin
		return []entSchema.PermittedAction{
			entSchema.ActionEditor,
			entSchema.ActionAdmin,
		}
	case entSchema.ActionAdmin:
		// Admin permissions are only implied by Admin
		return []entSchema.PermittedAction{
			entSchema.ActionAdmin,
		}
	default:
		// For any other action, just use that action
		return []entSchema.PermittedAction{action}
	}
}

// GetUserPermissionsForResource returns all permissions a user has for a specific resource
func (self *PermissionsRepository) GetUserPermissionsForResource(
	ctx context.Context,
	userID uuid.UUID,
	resourceType entSchema.ResourceType,
	resourceID uuid.UUID,
) ([]entSchema.PermittedAction, error) {
	result := make([]entSchema.PermittedAction, 0)

	// Check for all possible permission levels
	for _, action := range []entSchema.PermittedAction{
		entSchema.ActionAdmin,
		entSchema.ActionEditor,
		entSchema.ActionViewer,
	} {
		// Use the same comprehensive check we use for authorization
		hasPermission, err := self.checkComprehensivePermission(
			ctx,
			self.getUserGroupIDs(ctx, userID),
			action,
			resourceType,
			resourceID,
		)
		if err != nil {
			return nil, err
		}

		if hasPermission {
			result = append(result, action)
		}
	}

	return result, nil
}

// getUserGroupIDs gets all group IDs for a user
func (self *PermissionsRepository) getUserGroupIDs(ctx context.Context, userID uuid.UUID) []uuid.UUID {
	userGroups, err := self.userRepo.GetGroups(ctx, userID)
	if err != nil {
		return []uuid.UUID{}
	}

	var groupIDs []uuid.UUID
	for _, g := range userGroups {
		groupIDs = append(groupIDs, g.ID)
	}

	return groupIDs
}

// CreatePermission creates a new permission
func (self *PermissionsRepository) CreatePermission(
	ctx context.Context,
	groupID uuid.UUID,
	action entSchema.PermittedAction,
	resourceType entSchema.ResourceType,
	selector entSchema.ResourceSelector,
) (*ent.Permission, error) {
	return self.base.DB.Permission.Create().
		SetAction(action).
		SetResourceType(resourceType).
		SetResourceSelector(selector).
		AddGroupIDs(groupID).
		Save(ctx)
}

// DeletePermission deletes a permission
func (self *PermissionsRepository) DeletePermission(
	ctx context.Context,
	permissionID uuid.UUID,
) error {
	return self.base.DB.Permission.DeleteOneID(permissionID).Exec(ctx)
}

// GetPermissionsByGroup gets all permissions for a group
func (self *PermissionsRepository) GetPermissionsByGroup(
	ctx context.Context,
	groupID uuid.UUID,
) ([]*ent.Permission, error) {
	return self.base.DB.Permission.Query().
		Where(permission.HasGroupsWith(group.ID(groupID))).
		All(ctx)
}
