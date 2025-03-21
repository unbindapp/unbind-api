package builds_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *BuildsService) GetBuildJobsForService(ctx context.Context, requesterUserId uuid.UUID, input *models.GetBuildJobsInput) ([]*models.BuildJobResponse, *models.PaginationResponseMetadata, error) {
	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserId, []permissions_repo.PermissionCheck{
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   input.TeamID.String(),
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   "*",
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   input.ProjectID.String(),
		},
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeService,
			ResourceID:   input.ServiceID.String(),
		},
	}); err != nil {
		return nil, nil, err
	}

	if err := self.validateInputs(ctx, input); err != nil {
		return nil, nil, err
	}

	// Huma doesn't support pointers in query so we need to convert zero's to nil
	var cursor *time.Time
	if !input.Cursor.IsZero() {
		cursor = &input.Cursor
	}
	var status *schema.BuildJobStatus
	if input.Status != "" {
		status = &input.Status
	}
	// Get build jobs
	buildJobs, nextCursor, err := self.repo.BuildJob().GetByServiceIDPaginated(ctx, input.ServiceID, cursor, status)
	if err != nil {
		return nil, nil, err
	}

	// Transform response
	resp := models.TransformBuildJobEntities(buildJobs)

	// Get pagination metadata
	metadata := &models.PaginationResponseMetadata{
		HasNext:        nextCursor != nil,
		NextCursor:     nextCursor,
		PreviousCursor: cursor,
	}

	return resp, metadata, nil

}

func (self *BuildsService) validateInputs(ctx context.Context, input *models.GetBuildJobsInput) error {
	// Get team
	team, err := self.repo.Team().GetByID(ctx, input.TeamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound, input.TeamID.String())
		}
		return err
	}

	// Get project
	project, err := self.repo.Project().GetByID(ctx, input.ProjectID)
	if err != nil {
		if ent.IsNotFound(err) {
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound, input.ProjectID.String())
		}
		return err
	}

	if project.TeamID != team.ID {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "project does not belong to team")
	}

	var environmentID *uuid.UUID
	for _, env := range project.Edges.Environments {
		if env.ID == input.EnvironmentID {
			environmentID = utils.ToPtr(env.ID)
			break
		}
	}

	if environmentID == nil {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "environment does not belong to project")
	}

	// Get service
	service, err := self.repo.Service().GetByID(ctx, input.ServiceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound, "service not found")
		}
		return err
	}

	if service.EnvironmentID != *environmentID {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "service does not belong to environment")
	}
	return nil
}
