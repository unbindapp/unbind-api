package project_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type UpdateProjectInput struct {
	TeamID               uuid.UUID  `validate:"required,uuid4"`
	ProjectID            uuid.UUID  `validate:"required,uuid4"`
	DefaultEnvironmentID *uuid.UUID `validate:"omitempty,uuid4"`
	DisplayName          string
	Description          *string
}

func (self *ProjectService) UpdateProject(ctx context.Context, requesterUserID uuid.UUID, input *UpdateProjectInput) (*models.ProjectResponse, error) {
	// Validate input
	if err := validate.Validator().Struct(input); err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}
	if input.DisplayName == "" && input.Description == nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "No fields to update")
	}

	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to update project
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeProject,
			ResourceID:   input.ProjectID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Check if the team exists
	_, err := self.repo.Team().GetByID(ctx, input.TeamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, err
	}

	// Make sure project exists and is in the team
	project, err := self.repo.Project().GetByID(ctx, input.ProjectID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Project not found")
		}
		return nil, err
	}
	if project.TeamID != input.TeamID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Project not in team")
	}

	// Update the project
	project, err = self.repo.Project().Update(ctx, nil, input.ProjectID, input.DefaultEnvironmentID, input.DisplayName, input.Description)
	if err != nil {
		return nil, err
	}

	// Convert to response
	return models.TransformProjectEntity(project), nil
}
