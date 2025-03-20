package project_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *ProjectService) GetProjectsInTeam(ctx context.Context, requesterUserID uuid.UUID, teamID uuid.UUID, sortBy models.SortByField, sortOrder models.SortOrder) ([]*models.ProjectResponse, error) {
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
	projects, err := self.repo.Project().GetByTeam(ctx, teamID, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	// Transform response
	resp := models.TransformProjectEntitities(projects)

	// Summarizes services
	for _, project := range resp {
		environmentIDs := make([]uuid.UUID, len(project.Environments))
		for i, environment := range project.Environments {
			environmentIDs[i] = environment.ID
		}
		counts, runtimeSummaries, frameworkSummeries, err := self.repo.Service().SummarizeServices(ctx, environmentIDs)
		if err != nil {
			return nil, err
		}
		project.AttachServiceSummary(counts, runtimeSummaries, frameworkSummeries)
	}

	return resp, nil
}

// Get a single project by ID
func (self *ProjectService) GetProjectByID(ctx context.Context, requesterUserID uuid.UUID, teamID uuid.UUID, projectID uuid.UUID) (*models.ProjectResponse, error) {
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
	project, err := self.repo.Project().GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	resp := models.TransformProjectEntity(project)

	// Summarizes services
	environmentIDs := make([]uuid.UUID, len(resp.Environments))
	for i, environment := range resp.Environments {
		environmentIDs[i] = environment.ID
	}
	counts, runtimeSummaries, frameworkSummaries, err := self.repo.Service().SummarizeServices(ctx, environmentIDs)
	if err != nil {
		return nil, err
	}
	resp.AttachServiceSummary(counts, runtimeSummaries, frameworkSummaries)

	return resp, nil
}
