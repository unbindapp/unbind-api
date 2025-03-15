package project_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repository/permissions"
)

func (self *ProjectService) GetProjectsInTeam(ctx context.Context, requesterUserID uuid.UUID, teamID uuid.UUID) ([]*ProjectResponse, error) {
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
		// Has permission to read the specific team
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   teamID.String(),
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Check if the team exists
	_, err := self.repo.Team().GetByID(ctx, teamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, err
	}

	// Get projects
	projects, err := self.repo.Project().GetByTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	projectResponse := make([]*ProjectResponse, len(projects))
	for i, project := range projects {
		environmentsResponse := make([]*EnvironmentResponse, len(project.Edges.Environments))
		for j, environment := range project.Edges.Environments {
			environmentsResponse[j] = &EnvironmentResponse{
				ID:          environment.ID,
				Name:        environment.Name,
				DisplayName: environment.DisplayName,
				Description: environment.Description,
				CreatedAt:   environment.CreatedAt,
			}
		}

		projectResponse[i] = &ProjectResponse{
			ID:           project.ID,
			Name:         project.Name,
			DisplayName:  project.DisplayName,
			Description:  project.Description,
			Status:       project.Status,
			TeamID:       project.TeamID,
			CreatedAt:    project.CreatedAt,
			Environments: environmentsResponse,
		}
	}

	return projectResponse, nil
}
