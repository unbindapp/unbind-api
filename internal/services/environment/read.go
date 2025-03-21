package environment_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Get a single environment by ID
func (self *EnvironmentService) GetEnvironmentByID(ctx context.Context, requesterUserID uuid.UUID, teamID uuid.UUID, projectID uuid.UUID, environmentID uuid.UUID) (*models.EnvironmentResponse, error) {
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
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   "*",
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   projectID.String(),
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   environmentID.String(),
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Check if the team exists
	team, err := self.repo.Team().GetByID(ctx, teamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, err
	}

	var project *ent.Project
	for _, p := range team.Edges.Projects {
		if p.ID == projectID {
			project = p
			break
		}
	}
	if project == nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Project not found")
	}

	// Get environment
	var environment *ent.Environment
	for _, e := range project.Edges.Environments {
		if e.ID == environmentID {
			environment = e
			break
		}
	}
	if environment == nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Environment not found")
	}

	// Convert to response
	resp := models.TransformEnvironmentEntity(environment)

	// Summarizes services
	counts, providerSummaries, frameworkSummaries, err := self.repo.Service().SummarizeServices(ctx, []uuid.UUID{environmentID})
	if err != nil {
		return nil, err
	}
	resp.ServiceCount, _ = counts[environmentID]
	resp.FrameworkSummary, _ = frameworkSummaries[environmentID]
	resp.ProviderSummary, _ = providerSummaries[environmentID]

	return resp, nil
}
