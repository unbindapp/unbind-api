package k8s

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/repository/repositories"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// RBACManager integrates unbind groups with Kubernetes RBAC.
type RBACManager struct {
	repo       repositories.RepositoriesInterface
	kubeClient *KubeClient
}

func NewRBACManager(repository repositories.RepositoriesInterface, kubeClient *KubeClient) *RBACManager {
	return &RBACManager{
		repo:       repository,
		kubeClient: kubeClient,
	}
}

// GroupRBACConfig defines the RBAC configuration for a group
type GroupRBACConfig struct {
	// AccessLevel determines the verbs allowed (admin, edit, view)
	AccessLevel string
	// TeamIDs the group has access to
	TeamIDs []string
	// ProjectIDs the group has access to (if empty but TeamIDs has values, access to all team projects)
	ProjectIDs []string
	// Additional namespaces to restrict to (besides team namespaces)
	AdditionalNamespaces []string
	// ResourceTypes to allow access to (pods, services, deployments, etc.)
	ResourceTypes []string
}

// CreateK8sConfigFromPermissions builds a K8s RBAC config from permissions
func (self *RBACManager) CreateK8sConfigFromPermissions(permissions []*ent.Permission) GroupRBACConfig {
	config := GroupRBACConfig{
		AccessLevel: "view", // Default to view access
		TeamIDs:     []string{},
		ProjectIDs:  []string{},
	}

	// Track highest access level based on the permissions
	highestAccessLevel := ""

	// Process permissions
	for _, p := range permissions {
		// Determine access level for this permission
		accessLevel := getAccessLevelFromAction(p.Action)

		// Update highest access level if this is higher
		if compareAccessLevels(accessLevel, highestAccessLevel) > 0 {
			highestAccessLevel = accessLevel
		}

		// Process based on resource type
		switch p.ResourceType {
		case permission.ResourceTypeTeam:
			if p.ResourceID != "*" {
				config.TeamIDs = append(config.TeamIDs, p.ResourceID)
			}
		case permission.ResourceTypeProject:
			if p.ResourceID != "*" {
				config.ProjectIDs = append(config.ProjectIDs, p.ResourceID)
			}
		}
	}

	// Set the access level based on the highest found
	if highestAccessLevel != "" {
		config.AccessLevel = highestAccessLevel
	}

	// Deduplicate team and project IDs
	config.TeamIDs = deduplicateStrings(config.TeamIDs)
	config.ProjectIDs = deduplicateStrings(config.ProjectIDs)

	return config
}

// Helper function to get access level from permission action
func getAccessLevelFromAction(action permission.Action) string {
	switch action {
	case permission.ActionAdmin:
		return "admin"
	case permission.ActionManage:
		return "admin" // Map manage to admin for K8s purposes
	case permission.ActionEdit:
		return "edit"
	case permission.ActionCreate, permission.ActionUpdate:
		return "edit" // Map create/update to edit for K8s
	case permission.ActionRead, permission.ActionView:
		return "view"
	default:
		return "view" // Default to view for unknown actions
	}
}

// Helper function to compare access levels (higher returns positive number)
func compareAccessLevels(a, b string) int {
	levels := map[string]int{
		"admin": 3,
		"edit":  2,
		"view":  1,
		"":      0, // Empty string is lowest
	}

	levelA, okA := levels[a]
	levelB, okB := levels[b]

	if !okA {
		levelA = 0
	}

	if !okB {
		levelB = 0
	}

	return levelA - levelB
}

// Helper function to deduplicate string slices
func deduplicateStrings(strSlice []string) []string {
	keys := make(map[string]bool)
	result := []string{}

	for _, item := range strSlice {
		if _, exists := keys[item]; !exists {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

// SyncGroupToK8s creates or updates Kubernetes RBAC resources for a group
func (self *RBACManager) SyncGroupToK8s(ctx context.Context, groupID uuid.UUID, config GroupRBACConfig) error {
	// Get the group from the database
	g, err := self.repo.Group().GetByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	// Create a Kubernetes-safe name for the group
	k8sGroupName := fmt.Sprintf("unbind-group-%s", g.ID.String())

	// Create the ClusterRole with label selectors
	if err := self.createOrUpdateClusterRole(ctx, k8sGroupName, g.Name, config); err != nil {
		return err
	}

	// Create the ClusterRoleBinding
	if err := self.createOrUpdateClusterRoleBinding(ctx, k8sGroupName, g.Name); err != nil {
		return err
	}

	// Update the group in our database to store K8s reference
	err = self.repo.Group().UpdateK8sRoleName(ctx, g, k8sGroupName)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	return nil
}

// createOrUpdateClusterRole creates or updates a ClusterRole for the given group
// ! TODO - we may need to replace this with Roles and RoleBindings, so it can be tied to a specific namespace, this can grant too much power
// ! kubernetes doesn't do label matching with these roles apparently.
func (self *RBACManager) createOrUpdateClusterRole(ctx context.Context, roleName, displayName string, config GroupRBACConfig) error {
	// Determine verbs based on access level
	var verbs []string
	switch config.AccessLevel {
	case "admin":
		verbs = []string{"get", "list", "watch", "create", "update", "patch", "delete"}
	case "edit":
		verbs = []string{"get", "list", "watch", "create", "update", "patch"}
	default: // view
		verbs = []string{"get", "list", "watch"}
	}

	// Default resources if not specified
	resources := config.ResourceTypes
	if len(resources) == 0 {
		// Common resources to access
		resources = []string{
			// Core workload resources
			"pods", "services", "deployments", "statefulsets",
			"replicasets", "daemonsets", "jobs", "cronjobs",

			// Configuration resources
			"configmaps", "secrets",

			// Storage resources
			"persistentvolumeclaims",

			// Networking resources
			"ingresses", "networkpolicies", "services",

			// Observability resources
			"events", "endpoints",

			// Extensions/API resources
			"customresourcedefinitions", "ingressclasses",
		}
	}

	// Build label selectors based on team and project IDs
	labelSelectors := []map[string]string{}

	// Add team label selectors
	for _, teamID := range config.TeamIDs {
		labelSelectors = append(labelSelectors, map[string]string{
			"unbind-team": teamID,
		})
	}

	// Add project label selectors
	for _, projectID := range config.ProjectIDs {
		labelSelectors = append(labelSelectors, map[string]string{
			"unbind-project": projectID,
		})
	}

	// If no selectors were added but we have team/project permissions,
	// default to allow nothing (will be expanded when teams/projects added)
	if len(labelSelectors) == 0 && (len(config.TeamIDs) > 0 || len(config.ProjectIDs) > 0) {
		// Use a dummy selector that won't match anything
		labelSelectors = append(labelSelectors, map[string]string{
			"unbind-no-match": "true",
		})
	}

	// Build rules based on selectors and permissions
	var rules []interface{}

	// If we have no specific selectors and no team/project permissions,
	// create a catch-all rule (global admin)
	if len(labelSelectors) == 0 && len(config.TeamIDs) == 0 && len(config.ProjectIDs) == 0 {
		// Core API group ("") resources
		rules = append(rules, map[string]interface{}{
			"apiGroups": []interface{}{""},
			"resources": interfaceFromStrings(filterResourcesByAPIGroup(resources, "")),
			"verbs":     interfaceFromStrings(verbs),
		})

		// Apps API group resources
		rules = append(rules, map[string]interface{}{
			"apiGroups": []interface{}{"apps"},
			"resources": interfaceFromStrings(filterResourcesByAPIGroup(resources, "apps")),
			"verbs":     interfaceFromStrings(verbs),
		})

		// Batch API group resources
		rules = append(rules, map[string]interface{}{
			"apiGroups": []interface{}{"batch"},
			"resources": interfaceFromStrings(filterResourcesByAPIGroup(resources, "batch")),
			"verbs":     interfaceFromStrings(verbs),
		})

		// Networking API group resources
		rules = append(rules, map[string]interface{}{
			"apiGroups": []interface{}{"networking.k8s.io"},
			"resources": interfaceFromStrings(filterResourcesByAPIGroup(resources, "networking.k8s.io")),
			"verbs":     interfaceFromStrings(verbs),
		})

		// API Extensions API group resources
		rules = append(rules, map[string]interface{}{
			"apiGroups": []interface{}{"apiextensions.k8s.io"},
			"resources": interfaceFromStrings(filterResourcesByAPIGroup(resources, "apiextensions.k8s.io")),
			"verbs":     interfaceFromStrings(verbs),
		})
	} else {
		// Create rules for each label selector
		for _, selector := range labelSelectors {
			// Core API group ("") resources
			if len(filterResourcesByAPIGroup(resources, "")) > 0 {
				rule := map[string]interface{}{
					"apiGroups": []interface{}{""},
					"resources": interfaceFromStrings(filterResourcesByAPIGroup(resources, "")),
					"verbs":     interfaceFromStrings(verbs),
				}
				rule["resourceSelector"] = map[string]interface{}{
					"matchLabels": selector,
				}
				rules = append(rules, rule)
			}

			// Apps API group resources
			if len(filterResourcesByAPIGroup(resources, "apps")) > 0 {
				rule := map[string]interface{}{
					"apiGroups": []interface{}{"apps"},
					"resources": interfaceFromStrings(filterResourcesByAPIGroup(resources, "apps")),
					"verbs":     interfaceFromStrings(verbs),
				}
				rule["resourceSelector"] = map[string]interface{}{
					"matchLabels": selector,
				}
				rules = append(rules, rule)
			}

			// Batch API group resources
			if len(filterResourcesByAPIGroup(resources, "batch")) > 0 {
				rule := map[string]interface{}{
					"apiGroups": []interface{}{"batch"},
					"resources": interfaceFromStrings(filterResourcesByAPIGroup(resources, "batch")),
					"verbs":     interfaceFromStrings(verbs),
				}
				rule["resourceSelector"] = map[string]interface{}{
					"matchLabels": selector,
				}
				rules = append(rules, rule)
			}

			// Networking API group resources
			if len(filterResourcesByAPIGroup(resources, "networking.k8s.io")) > 0 {
				rule := map[string]interface{}{
					"apiGroups": []interface{}{"networking.k8s.io"},
					"resources": interfaceFromStrings(filterResourcesByAPIGroup(resources, "networking.k8s.io")),
					"verbs":     interfaceFromStrings(verbs),
				}
				rule["resourceSelector"] = map[string]interface{}{
					"matchLabels": selector,
				}
				rules = append(rules, rule)
			}

			// API Extensions API group resources
			if len(filterResourcesByAPIGroup(resources, "apiextensions.k8s.io")) > 0 {
				rule := map[string]interface{}{
					"apiGroups": []interface{}{"apiextensions.k8s.io"},
					"resources": interfaceFromStrings(filterResourcesByAPIGroup(resources, "apiextensions.k8s.io")),
					"verbs":     interfaceFromStrings(verbs),
				}
				rule["resourceSelector"] = map[string]interface{}{
					"matchLabels": selector,
				}
				rules = append(rules, rule)
			}
		}
	}

	// Define the ClusterRole
	cr := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRole",
			"metadata": map[string]interface{}{
				"name": roleName,
				"labels": map[string]interface{}{
					"app.kubernetes.io/managed-by": "unbind",
					"unbind.app/group-name":        displayName,
				},
			},
			"rules": rules,
		},
	}

	// Define the resource schema for ClusterRoles
	crResource := schema.GroupVersionResource{
		Group:    "rbac.authorization.k8s.io",
		Version:  "v1",
		Resource: "clusterroles",
	}

	// Try to get existing ClusterRole
	_, err := self.kubeClient.client.Resource(crResource).Get(ctx, roleName, metav1.GetOptions{})
	if err == nil {
		// Update the existing ClusterRole
		_, err = self.kubeClient.client.Resource(crResource).Update(ctx, cr, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update ClusterRole: %w", err)
		}
		log.Printf("ClusterRole '%s' updated successfully", roleName)
	} else {
		// Create a new ClusterRole
		_, err = self.kubeClient.client.Resource(crResource).Create(ctx, cr, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create ClusterRole: %w", err)
		}
		log.Printf("ClusterRole '%s' created successfully", roleName)
	}

	return nil
}

// Helper function to filter resources by API group
func filterResourcesByAPIGroup(resources []string, apiGroup string) []string {
	var filteredResources []string

	switch apiGroup {
	case "": // Core API group
		coreResources := []string{
			"pods", "services", "configmaps", "secrets",
			"persistentvolumeclaims", "events", "endpoints",
		}
		for _, r := range resources {
			for _, cr := range coreResources {
				if r == cr {
					filteredResources = append(filteredResources, r)
					break
				}
			}
		}
	case "apps":
		appsResources := []string{
			"deployments", "statefulsets", "replicasets", "daemonsets",
		}
		for _, r := range resources {
			for _, ar := range appsResources {
				if r == ar {
					filteredResources = append(filteredResources, r)
					break
				}
			}
		}
	case "batch":
		batchResources := []string{
			"jobs", "cronjobs",
		}
		for _, r := range resources {
			for _, br := range batchResources {
				if r == br {
					filteredResources = append(filteredResources, r)
					break
				}
			}
		}
	case "networking.k8s.io":
		networkingResources := []string{
			"ingresses", "networkpolicies", "ingressclasses",
		}
		for _, r := range resources {
			for _, nr := range networkingResources {
				if r == nr {
					filteredResources = append(filteredResources, r)
					break
				}
			}
		}
	case "apiextensions.k8s.io":
		apiExtResources := []string{
			"customresourcedefinitions",
		}
		for _, r := range resources {
			for _, ar := range apiExtResources {
				if r == ar {
					filteredResources = append(filteredResources, r)
					break
				}
			}
		}
	}

	return filteredResources
}

// createOrUpdateClusterRoleBinding creates or updates a ClusterRoleBinding for the given group
func (self *RBACManager) createOrUpdateClusterRoleBinding(ctx context.Context, roleName, groupName string) error {
	bindingName := fmt.Sprintf("binding-%s", roleName)

	// Define the ClusterRoleBinding
	crb := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRoleBinding",
			"metadata": map[string]interface{}{
				"name": bindingName,
				"labels": map[string]interface{}{
					"app.kubernetes.io/managed-by": "unbind",
					"unbind.app/group-name":        groupName,
				},
			},
			"subjects": []interface{}{
				map[string]interface{}{
					"kind":     "Group",
					"name":     "oidc:" + groupName,
					"apiGroup": "rbac.authorization.k8s.io",
				},
			},
			"roleRef": map[string]interface{}{
				"kind":     "ClusterRole",
				"name":     roleName,
				"apiGroup": "rbac.authorization.k8s.io",
			},
		},
	}

	// Define the resource schema for ClusterRoleBindings
	crbResource := schema.GroupVersionResource{
		Group:    "rbac.authorization.k8s.io",
		Version:  "v1",
		Resource: "clusterrolebindings",
	}

	// Try to get existing ClusterRoleBinding
	_, err := self.kubeClient.client.Resource(crbResource).Get(ctx, bindingName, metav1.GetOptions{})
	if err == nil {
		// Update the existing ClusterRoleBinding
		_, err = self.kubeClient.client.Resource(crbResource).Update(ctx, crb, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update ClusterRoleBinding: %w", err)
		}
		log.Printf("ClusterRoleBinding '%s' updated successfully", bindingName)
	} else {
		// Create a new ClusterRoleBinding
		_, err = self.kubeClient.client.Resource(crbResource).Create(ctx, crb, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create ClusterRoleBinding: %w", err)
		}
		log.Printf("ClusterRoleBinding '%s' created successfully", bindingName)
	}

	return nil
}

// SyncAllGroups synchronizes all groups with Kubernetes RBAC
func (self *RBACManager) SyncAllGroups(ctx context.Context) error {
	groups, err := self.repo.Group().GetAllWithPermissions(ctx)
	if err != nil {
		return fmt.Errorf("failed to query groups: %w", err)
	}

	for _, g := range groups {
		// Determine if this group should have K8s access
		hasK8sAccess := false
		for _, p := range g.Edges.Permissions {
			if p.ResourceType == permission.ResourceTypeTeam || p.ResourceType == permission.ResourceTypeProject {
				hasK8sAccess = true
				break
			}
		}

		if hasK8sAccess {
			// Create the RBAC config from the group's permissions
			config := self.CreateK8sConfigFromPermissions(g.Edges.Permissions)

			// Sync the group to K8s
			if err := self.SyncGroupToK8s(ctx, g.ID, config); err != nil {
				log.Printf("Warning: failed to sync group %s to K8s: %v", g.Name, err)
			}
		} else if g.K8sRoleName != "" {
			// Group no longer has K8s access but did before, remove it
			if err := self.DeleteK8sRBAC(ctx, g.ID); err != nil {
				log.Printf("Warning: failed to delete K8s RBAC for group %s: %v", g.Name, err)
			}
		}
	}

	return nil
}

// DeleteK8sRBAC removes Kubernetes RBAC resources for a group
func (self *RBACManager) DeleteK8sRBAC(ctx context.Context, groupID uuid.UUID) error {
	// Get the group from the database
	g, err := self.repo.Group().GetByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	// Check if the group has a K8s role name
	if g.K8sRoleName == "" {
		return nil // Nothing to delete
	}

	roleName := g.K8sRoleName
	bindingName := fmt.Sprintf("binding-%s", roleName)

	// Define the resource schemas
	crResource := schema.GroupVersionResource{
		Group:    "rbac.authorization.k8s.io",
		Version:  "v1",
		Resource: "clusterroles",
	}

	crbResource := schema.GroupVersionResource{
		Group:    "rbac.authorization.k8s.io",
		Version:  "v1",
		Resource: "clusterrolebindings",
	}

	// Delete the ClusterRoleBinding first
	err = self.kubeClient.client.Resource(crbResource).Delete(ctx, bindingName, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Warning: failed to delete ClusterRoleBinding: %v", err)
	} else {
		log.Printf("ClusterRoleBinding '%s' deleted successfully", bindingName)
	}

	// Delete the ClusterRole
	err = self.kubeClient.client.Resource(crResource).Delete(ctx, roleName, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Warning: failed to delete ClusterRole: %v", err)
	} else {
		log.Printf("ClusterRole '%s' deleted successfully", roleName)
	}

	// Update the group in our database to remove K8s reference
	err = self.repo.Group().ClearK8sRoleName(ctx, g)

	return err
}

// GetGroupsWithK8sAccess returns all groups that have Kubernetes RBAC access
func (self *RBACManager) GetGroupsWithK8sAccess(ctx context.Context) ([]*ent.Group, error) {
	return self.repo.Group().GetAllWithK8sRole(ctx)
}

// CleanupOrphanedRBAC checks for RBAC resources managed by unbind that don't
// have corresponding groups in the database and removes them
func (self *RBACManager) CleanupOrphanedRBAC(ctx context.Context) error {
	// Get all ClusterRoles with our managed-by label
	crResource := schema.GroupVersionResource{
		Group:    "rbac.authorization.k8s.io",
		Version:  "v1",
		Resource: "clusterroles",
	}

	// List all ClusterRoles with our managed-by label
	crList, err := self.kubeClient.client.Resource(crResource).List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=unbind",
	})
	if err != nil {
		return fmt.Errorf("failed to list ClusterRoles: %w", err)
	}

	// Get all groups from the database that have K8s roles
	groups, err := self.repo.Group().GetAllWithK8sRole(ctx)
	if err != nil {
		return fmt.Errorf("failed to get groups with K8s roles: %w", err)
	}

	// Create a map of valid K8s role names
	validRoleNames := make(map[string]bool)
	for _, g := range groups {
		if g.K8sRoleName != "" {
			validRoleNames[g.K8sRoleName] = true
		}
	}

	// Check each ClusterRole to see if it has a corresponding group
	for _, cr := range crList.Items {
		roleName := cr.GetName()

		// Skip if this is a valid role name
		if validRoleNames[roleName] {
			continue
		}

		// This is an orphaned ClusterRole, delete it and its binding
		log.Printf("Found orphaned ClusterRole: %s, cleaning up", roleName)

		// Delete the ClusterRoleBinding
		crbResource := schema.GroupVersionResource{
			Group:    "rbac.authorization.k8s.io",
			Version:  "v1",
			Resource: "clusterrolebindings",
		}

		bindingName := fmt.Sprintf("binding-%s", roleName)
		err = self.kubeClient.client.Resource(crbResource).Delete(ctx, bindingName, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Warning: failed to delete orphaned ClusterRoleBinding %s: %v", bindingName, err)
		} else {
			log.Printf("Deleted orphaned ClusterRoleBinding: %s", bindingName)
		}

		// Delete the ClusterRole
		err = self.kubeClient.client.Resource(crResource).Delete(ctx, roleName, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Warning: failed to delete orphaned ClusterRole %s: %v", roleName, err)
		} else {
			log.Printf("Deleted orphaned ClusterRole: %s", roleName)
		}
	}

	return nil
}

// Helper function to convert a string slice to an interface slice
func interfaceFromStrings(strings []string) []interface{} {
	interfaces := make([]interface{}, len(strings))
	for i, s := range strings {
		interfaces[i] = s
	}
	return interfaces
}
