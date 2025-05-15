package permissions_repo

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/environment"
	"github.com/unbindapp/unbind-api/ent/group"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/ent/predicate"
	"github.com/unbindapp/unbind-api/ent/project"
	entSchema "github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/team"
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

// * Predicate-Based Permission Logic *

// GetPotentialParents returns a list of resource types that can be parents of the given childType.
// This should ideally be derived directly from schema definitions or a central schema utility.
func GetPotentialParents(childType entSchema.ResourceType) []entSchema.ResourceType {
	switch childType {
	case entSchema.ResourceTypeService:
		return []entSchema.ResourceType{entSchema.ResourceTypeEnvironment, entSchema.ResourceTypeProject, entSchema.ResourceTypeTeam}
	case entSchema.ResourceTypeEnvironment:
		return []entSchema.ResourceType{entSchema.ResourceTypeProject, entSchema.ResourceTypeTeam}
	case entSchema.ResourceTypeProject:
		return []entSchema.ResourceType{entSchema.ResourceTypeTeam}
	default:
		return []entSchema.ResourceType{}
	}
}

// GetAccessibleProjectPredicates returns Ent predicates for filtering projects
// that the user has the given action permission for.
// Returns nil predicate and nil error if user is superuser for projects or access is broadly granted (matches all).
// Returns a predicate that matches nothing if no access is found.
// Returns an error if an issue occurs.
func (self *PermissionsRepository) GetAccessibleProjectPredicates(
	ctx context.Context,
	userID uuid.UUID,
	action entSchema.PermittedAction,
) (predicate.Project, error) {
	groupIDs := self.getUserGroupIDs(ctx, userID)
	if len(groupIDs) == 0 {
		// No groups, implies no group-based permissions.
		// Return a predicate that matches nothing.
		return func(s *sql.Selector) { s.Where(sql.False()) }, nil
	}

	impliedActions := self.getImpliedActions(action)

	return self.buildProjectPredicatesInternal(ctx, groupIDs, impliedActions)
}

// buildProjectPredicatesInternal is a helper to construct predicates for Project entities.
func (self *PermissionsRepository) buildProjectPredicatesInternal(
	ctx context.Context,
	groupIDs []uuid.UUID,
	impliedActions []entSchema.PermittedAction,
) (predicate.Project, error) {
	var projectPredicates []predicate.Project

	// 1. Check for "superuser" permission on ResourceTypeProject
	hasSuperuserProjectAccess, err := self.base.DB.Permission.Query().
		Where(
			permission.HasGroupsWith(group.IDIn(groupIDs...)),
			permission.ActionIn(impliedActions...),
			permission.ResourceTypeEQ(entSchema.ResourceTypeProject),
			func(s *sql.Selector) {
				sqljson.ValueEQ(permission.FieldResourceSelector, true, sqljson.Path("superuser"))
			},
		).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking superuser project access: %w", err)
	}

	if hasSuperuserProjectAccess {
		// User has superuser access to all projects, no specific predicate needed (matches all).
		return nil, nil // Returning nil predicate means "all allowed"
	}

	// 2. Direct access to specific Project IDs
	directProjectPermissions, err := self.base.DB.Permission.Query().
		Where(
			permission.HasGroupsWith(group.IDIn(groupIDs...)),
			permission.ActionIn(impliedActions...),
			permission.ResourceTypeEQ(entSchema.ResourceTypeProject),
			func(s *sql.Selector) { // Has an "id" field in resource_selector
				sqljson.ValueIsNotNull(permission.FieldResourceSelector, sqljson.Path("id"))
			},
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching direct project permissions: %w", err)
	}

	var specificProjectIDs []uuid.UUID
	for _, p := range directProjectPermissions {
		if p.ResourceSelector.ID != uuid.Nil {
			specificProjectIDs = append(specificProjectIDs, p.ResourceSelector.ID)
		}
	}
	if len(specificProjectIDs) > 0 {
		projectPredicates = append(projectPredicates, project.IDIn(specificProjectIDs...))
	}

	// 3. Hierarchical access: Projects belonging to accessible Teams
	// Find teams where the user has 'action' permission (direct or superuser for that team type)
	teamPermissions, err := self.base.DB.Permission.Query().
		Where(
			permission.HasGroupsWith(group.IDIn(groupIDs...)),
			permission.ActionIn(impliedActions...),
			permission.ResourceTypeEQ(entSchema.ResourceTypeTeam),
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching team permissions for project hierarchy: %w", err)
	}

	var accessibleTeamIDs []uuid.UUID
	foundGlobalSuperuserForTeamType := false
	for _, tp := range teamPermissions {
		if tp.ResourceSelector.Superuser { // Check if this permission grants superuser for ResourceTypeTeam
			foundGlobalSuperuserForTeamType = true
			break
		}
		if tp.ResourceSelector.ID != uuid.Nil {
			accessibleTeamIDs = append(accessibleTeamIDs, tp.ResourceSelector.ID)
		}
	}

	if foundGlobalSuperuserForTeamType {
		// User has a permission that makes them superuser for ResourceTypeTeam globally.
		// This implies access to ALL teams, and therefore projects in ALL teams.
		return nil, nil // Access to all projects because they are superuser for all teams
	}

	if len(accessibleTeamIDs) > 0 {
		projectPredicates = append(projectPredicates, project.HasTeamWith(team.IDIn(accessibleTeamIDs...)))
	}

	if len(projectPredicates) == 0 {
		// No specific permissions found, return a predicate that matches nothing.
		return func(s *sql.Selector) { s.Where(sql.False()) }, nil
	}

	return project.Or(projectPredicates...), nil
}

// GetAccessibleTeamPredicates returns Ent predicates for filtering teams
// that the user has the given action permission for.
func (self *PermissionsRepository) GetAccessibleTeamPredicates(
	ctx context.Context,
	userID uuid.UUID,
	action entSchema.PermittedAction,
) (predicate.Team, error) {
	groupIDs := self.getUserGroupIDs(ctx, userID)
	if len(groupIDs) == 0 {
		return func(s *sql.Selector) { s.Where(sql.False()) }, nil
	}
	impliedActions := self.getImpliedActions(action)
	return self.buildTeamPredicatesInternal(ctx, groupIDs, impliedActions)
}

func (self *PermissionsRepository) buildTeamPredicatesInternal(
	ctx context.Context,
	groupIDs []uuid.UUID,
	impliedActions []entSchema.PermittedAction,
) (predicate.Team, error) {
	var teamPredicates []predicate.Team

	// 1. Check for "superuser" permission on ResourceTypeTeam
	hasSuperuserTeamAccess, err := self.base.DB.Permission.Query().
		Where(
			permission.HasGroupsWith(group.IDIn(groupIDs...)),
			permission.ActionIn(impliedActions...),
			permission.ResourceTypeEQ(entSchema.ResourceTypeTeam),
			func(s *sql.Selector) {
				sqljson.ValueEQ(permission.FieldResourceSelector, true, sqljson.Path("superuser"))
			},
		).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking superuser team access: %w", err)
	}

	if hasSuperuserTeamAccess {
		return nil, nil // User has superuser access to all teams
	}

	// 2. Direct access to specific Team IDs
	directTeamPermissions, err := self.base.DB.Permission.Query().
		Where(
			permission.HasGroupsWith(group.IDIn(groupIDs...)),
			permission.ActionIn(impliedActions...),
			permission.ResourceTypeEQ(entSchema.ResourceTypeTeam),
			func(s *sql.Selector) { // Has an "id" field in resource_selector
				sqljson.ValueIsNotNull(permission.FieldResourceSelector, sqljson.Path("id"))
			},
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching direct team permissions: %w", err)
	}

	var specificTeamIDs []uuid.UUID
	for _, p := range directTeamPermissions {
		if p.ResourceSelector.ID != uuid.Nil {
			specificTeamIDs = append(specificTeamIDs, p.ResourceSelector.ID)
		}
	}
	if len(specificTeamIDs) > 0 {
		teamPredicates = append(teamPredicates, team.IDIn(specificTeamIDs...))
	}

	// Teams are top-level in this hierarchy, so no further hierarchical checks needed here.

	if len(teamPredicates) == 0 {
		// No specific permissions found, and not superuser implies no access.
		return func(s *sql.Selector) { s.Where(sql.False()) }, nil
	}

	return team.Or(teamPredicates...), nil
}

// GetAccessibleEnvironmentPredicates returns Ent predicates for filtering environments.
// It can be scoped by an optional projectID.
func (self *PermissionsRepository) GetAccessibleEnvironmentPredicates(
	ctx context.Context,
	userID uuid.UUID,
	action entSchema.PermittedAction,
	projectID *uuid.UUID, // Optional: For scoping environments to a specific project
) (predicate.Environment, error) {
	groupIDs := self.getUserGroupIDs(ctx, userID)
	if len(groupIDs) == 0 {
		return func(s *sql.Selector) { s.Where(sql.False()) }, nil
	}
	impliedActions := self.getImpliedActions(action)
	return self.buildEnvironmentPredicatesInternal(ctx, userID, groupIDs, action, impliedActions, projectID)
}

func (self *PermissionsRepository) buildEnvironmentPredicatesInternal(
	ctx context.Context,
	userID uuid.UUID,
	groupIDs []uuid.UUID,
	action entSchema.PermittedAction,
	impliedActions []entSchema.PermittedAction,
	projectID *uuid.UUID,
) (predicate.Environment, error) {
	var envPredicates []predicate.Environment

	// 1. Check for "superuser" permission on ResourceTypeEnvironment
	hasSuperuserEnvAccess, err := self.base.DB.Permission.Query().
		Where(
			permission.HasGroupsWith(group.IDIn(groupIDs...)),
			permission.ActionIn(impliedActions...),
			permission.ResourceTypeEQ(entSchema.ResourceTypeEnvironment),
			func(s *sql.Selector) {
				sqljson.ValueEQ(permission.FieldResourceSelector, true, sqljson.Path("superuser"))
			},
		).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking superuser environment access: %w", err)
	}
	if hasSuperuserEnvAccess {
		// If superuser for environments, projectID scope still applies if provided.
		if projectID != nil {
			return environment.HasProjectWith(project.IDEQ(*projectID)), nil
		}
		return nil, nil // Access to all environments (respecting projectID if applicable)
	}

	// 2. Direct access to specific Environment IDs
	directEnvPermissions, err := self.base.DB.Permission.Query().
		Where(
			permission.HasGroupsWith(group.IDIn(groupIDs...)),
			permission.ActionIn(impliedActions...),
			permission.ResourceTypeEQ(entSchema.ResourceTypeEnvironment),
			func(s *sql.Selector) { sqljson.ValueIsNotNull(permission.FieldResourceSelector, sqljson.Path("id")) },
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching direct environment permissions: %w", err)
	}
	var specificEnvIDs []uuid.UUID
	for _, p := range directEnvPermissions {
		if p.ResourceSelector.ID != uuid.Nil {
			specificEnvIDs = append(specificEnvIDs, p.ResourceSelector.ID)
		}
	}
	if len(specificEnvIDs) > 0 {
		envPredicates = append(envPredicates, environment.IDIn(specificEnvIDs...))
	}

	// 3. Hierarchical access: Environments belonging to accessible Projects or Teams
	//    a) Via Project access
	projectPerms, err := self.GetAccessibleProjectPredicates(ctx, userID, action)
	if err != nil {
		return nil, fmt.Errorf("error getting project predicates for env hierarchy: %w", err)
	}
	if projectPerms != nil {
		envPredicates = append(envPredicates, environment.HasProjectWith(projectPerms))
	}

	// Final combination
	var finalPredicate predicate.Environment
	if len(envPredicates) == 0 {
		// If no specific env permissions and not superuser on envs, and no project perms grant access,
		// then no access. However, if projectPerms was nil (superuser on projects), access *is* granted.
		// This case needs to be distinct from "no permissions at all".
		if projectPerms == nil && !hasSuperuserEnvAccess { // Superuser on projects but not on envs directly
			finalPredicate = nil // All environments (potentially scoped by projectID by caller)
		} else {
			finalPredicate = func(s *sql.Selector) { s.Where(sql.False()) } // No access
		}
	} else {
		finalPredicate = environment.Or(envPredicates...)
	}

	// Apply projectID scope if provided and not already handled by superuser_env + projectID case
	if projectID != nil {
		scopePredicate := environment.HasProjectWith(project.IDEQ(*projectID))
		if finalPredicate == nil { // Handles cases like superuser_project access
			return scopePredicate, nil
		}
		return environment.And(finalPredicate, scopePredicate), nil
	}

	return finalPredicate, nil
}

// GetAccessibleServicePredicates returns Ent predicates for filtering services.
// It can be scoped by an optional environmentID.
func (self *PermissionsRepository) GetAccessibleServicePredicates(
	ctx context.Context,
	userID uuid.UUID,
	action entSchema.PermittedAction,
	environmentID *uuid.UUID, // Optional: For scoping services to a specific environment
) (predicate.Service, error) {
	groupIDs := self.getUserGroupIDs(ctx, userID)
	if len(groupIDs) == 0 {
		return func(s *sql.Selector) { s.Where(sql.False()) }, nil
	}
	impliedActions := self.getImpliedActions(action)
	return self.buildServicePredicatesInternal(ctx, userID, groupIDs, action, impliedActions, environmentID)
}

func (self *PermissionsRepository) buildServicePredicatesInternal(
	ctx context.Context,
	userID uuid.UUID,
	groupIDs []uuid.UUID,
	action entSchema.PermittedAction,
	impliedActions []entSchema.PermittedAction,
	environmentID *uuid.UUID,
) (predicate.Service, error) {
	var servicePredicates []predicate.Service

	// 1. Check for "superuser" permission on ResourceTypeService
	hasSuperuserServiceAccess, err := self.base.DB.Permission.Query().
		Where(
			permission.HasGroupsWith(group.IDIn(groupIDs...)),
			permission.ActionIn(impliedActions...),
			permission.ResourceTypeEQ(entSchema.ResourceTypeService),
			func(s *sql.Selector) {
				sqljson.ValueEQ(permission.FieldResourceSelector, true, sqljson.Path("superuser"))
			},
		).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking superuser service access: %w", err)
	}
	if hasSuperuserServiceAccess {
		if environmentID != nil {
			return service.HasEnvironmentWith(environment.IDEQ(*environmentID)), nil
		}
		return nil, nil // Access to all services (respecting environmentID if applicable)
	}

	// 2. Direct access to specific Service IDs
	directServicePermissions, err := self.base.DB.Permission.Query().
		Where(
			permission.HasGroupsWith(group.IDIn(groupIDs...)),
			permission.ActionIn(impliedActions...),
			permission.ResourceTypeEQ(entSchema.ResourceTypeService),
			func(s *sql.Selector) { sqljson.ValueIsNotNull(permission.FieldResourceSelector, sqljson.Path("id")) },
		).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching direct service permissions: %w", err)
	}
	var specificServiceIDs []uuid.UUID
	for _, p := range directServicePermissions {
		if p.ResourceSelector.ID != uuid.Nil {
			specificServiceIDs = append(specificServiceIDs, p.ResourceSelector.ID)
		}
	}
	if len(specificServiceIDs) > 0 {
		servicePredicates = append(servicePredicates, service.IDIn(specificServiceIDs...))
	}

	// 3. Hierarchical access: Services belonging to accessible Environments, Projects, or Teams
	//    a) Via Environment access
	// We need projectID to get envPreds if environmentID is nil, but envPreds are scoped with projectID if provided.
	// This needs careful handling of how envPreds are fetched if environmentID itself is the primary scope.
	// For now, assume if environmentID is provided, we primarily care about that specific environment first.
	// If envPreds are needed broadly, the parent function might need to pass the parent projectID explicitly.
	var envProjectID *uuid.UUID // This would be the project ID of the *environmentID if provided, or nil
	if environmentID != nil {
		// If we have a specific environmentID, fetch its projectID for the GetAccessibleEnvironmentPredicates call,
		// as it expects a projectID for broader hierarchical checks if no specific environment perms exist.
		env, err := self.environmentRepo.GetByID(ctx, *environmentID) // Use self.environmentRepo
		if err == nil && env.Edges.Project != nil {                   // Check if project edge is loaded or use env.ProjectID
			// Ensure env.ProjectID is the correct field name if Edges.Project is not always loaded.
			// If env.Edges.Project is not guaranteed, using env.ProjectID (the foreign key field) is safer.
			// Let's assume env.ProjectID is available directly on the env entity.
			envProjectID = &env.ProjectID // Assuming ProjectID is a direct field on ent.Environment
		}
	}

	envPreds, err := self.GetAccessibleEnvironmentPredicates(ctx, userID, action, envProjectID)
	if err != nil {
		return nil, fmt.Errorf("error getting environment predicates for service hierarchy: %w", err)
	}
	if envPreds != nil {
		servicePredicates = append(servicePredicates, service.HasEnvironmentWith(envPreds))
	} else if !hasSuperuserServiceAccess { // nil envPreds means superuser on envs (or further up)
		// and not already superuser on services directly
		// This implies access to all services in all accessible environments.
		// This is a broad grant, effectively handled by returning nil from finalPredicate if no other specific predicates exist.
	}

	// Final combination
	var finalPredicate predicate.Service
	if len(servicePredicates) == 0 {
		if envPreds == nil && !hasSuperuserServiceAccess { // Superuser via hierarchy (e.g. on all envs)
			finalPredicate = nil
		} else {
			finalPredicate = func(s *sql.Selector) { s.Where(sql.False()) } // No access
		}
	} else {
		finalPredicate = service.Or(servicePredicates...)
	}

	// Apply environmentID scope if provided
	if environmentID != nil {
		scopePredicate := service.HasEnvironmentWith(environment.IDEQ(*environmentID))
		if finalPredicate == nil { // Handles cases like superuser_hierarchy access
			return scopePredicate, nil
		}
		return service.And(finalPredicate, scopePredicate), nil
	}

	return finalPredicate, nil
}
