package projects_handler

import (
	"context"
	"errors"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/services/models"
	project_service "github.com/unbindapp/unbind-api/internal/services/project"
)

type CreateProjectInput struct {
	server.BaseAuthInput
	Body struct {
		TeamID      uuid.UUID `json:"team_id" required:"true"`
		DisplayName string    `json:"display_name" required:"true"`
		Description string    `json:"description" required:"true"`
	}
}

type CreateProjectResponse struct {
	Body struct {
		Data *models.ProjectResponse `json:"data"`
	}
}

func (self *HandlerGroup) CreateProject(ctx context.Context, input *CreateProjectInput) (*CreateProjectResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	// Generate name
	name, err := utils.GenerateSlug(input.Body.DisplayName)
	if err != nil {
		log.Error("Error generating project name", "err", err)
		return nil, huma.Error500InternalServerError("Unable to generate project name")
	}

	createdProject, err := self.srv.ProjectService.CreateProject(ctx, user.ID, &project_service.CreateProjectInput{
		TeamID:      input.Body.TeamID,
		Name:        name,
		DisplayName: input.Body.DisplayName,
		Description: input.Body.Description,
	}, bearerToken)
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
		log.Error("Error creating project", "err", err)
		return nil, huma.Error500InternalServerError("Unable to create project")
	}

	resp := &CreateProjectResponse{}
	resp.Body.Data = createdProject
	return resp, nil
}
