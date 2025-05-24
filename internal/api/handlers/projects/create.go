package projects_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/models"
	project_service "github.com/unbindapp/unbind-api/internal/services/project"
)

type CreateProjectInput struct {
	server.BaseAuthInput
	Body *project_service.CreateProjectInput
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

	createdProject, err := self.srv.ProjectService.CreateProject(ctx, user.ID, input.Body, bearerToken)
	if err != nil {
		return nil, self.handleErr(err)
	}

	resp := &CreateProjectResponse{}
	resp.Body.Data = createdProject
	return resp, nil
}
