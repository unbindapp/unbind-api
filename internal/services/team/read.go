package team_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// ListTeams retrieves all teams the user has permission to view
func (self *TeamService) ListTeams(ctx context.Context, userID uuid.UUID, bearerToken string) ([]*models.TeamResponse, error) {
	// Start with a base query
	query := self.repo.Ent().Team.Query()

	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to read system resources
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		// Has permission to read teams
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
	}
	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Execute the query
	dbTeams, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	// Get teams from kubernetes
	k8sTeams, err := self.k8s.GetUnbindTeams(ctx, bearerToken)
	if err != nil {
		return nil, err
	}

	// Create a map of k8s namespaces
	k8sNamespaces := make(map[string]bool)
	for _, team := range k8sTeams {
		k8sNamespaces[team.Namespace] = true
	}

	// Filter dbTeams to only include those with namespaces in k8sNamespaces
	var filteredTeams []*ent.Team
	for _, team := range dbTeams {
		_, namespaceExists := k8sNamespaces[team.Namespace]
		if !namespaceExists {
			continue
		}
		filteredTeams = append(filteredTeams, team)
	}

	return models.TransformTeamEntities(filteredTeams), nil
}
