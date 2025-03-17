package projects_handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
	project_service "github.com/unbindapp/unbind-api/internal/services/project"
)

type UpdateProjectInput struct {
	server.BaseAuthInput
	Body struct {
		TeamID      uuid.UUID `json:"team_id" required:"true"`
		ProjectID   uuid.UUID `json:"project_id" required:"true"`
		DisplayName string    `json:"display_name"`
		Description string    `json:"description"`
	}
}

type UpdateProjectResponse struct {
	Body struct {
		Data *models.ProjectResponse `json:"data"`
	}
}

func (self *HandlerGroup) UpdateProject(ctx context.Context, input *UpdateProjectInput) (*UpdateProjectResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	updatedProject, err := self.srv.ProjectService.UpdateProject(ctx, user.ID, &project_service.UpdateProjectInput{
		TeamID:      input.Body.TeamID,
		DisplayName: input.Body.DisplayName,
		Description: input.Body.Description,
	})
	if err != nil {
		if errors.Is(err, errdefs.ErrInvalidInput) {
			return nil, huma.Error400BadRequest(err.Error())
		}
		if errors.Is(err, errdefs.ErrUnauthorized) {
			return nil, huma.Error403Forbidden("Unauthorized")
		}
		if ent.IsNotFound(err) || errors.Is(err, errdefs.ErrNotFound) {
			return nil, huma.Error404NotFound(err.Error())
		}
		log.Error("Error updating project", "err", err)
		return nil, huma.Error500InternalServerError("Unable to update project")
	}

	resp := &UpdateProjectResponse{}
	resp.Body.Data = updatedProject
	return resp, nil
}
