package projects_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
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
		return nil, self.handleErr(err)
	}

	resp := &CreateProjectResponse{}
	resp.Body.Data = createdProject
	return resp, nil
}
