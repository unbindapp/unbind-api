package projects_handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/errdefs"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/server"
	project_service "github.com/unbindapp/unbind-api/internal/services/project"
)

type ListProjectInput struct {
	server.BaseAuthInput
	TeamID uuid.UUID `path:"team_id"`
}

type ListProjectResponse struct {
	Body struct {
		Data []*project_service.ProjectResponse `json:"data"`
	}
}

func (self *HandlerGroup) ListProjects(ctx context.Context, input *ListProjectInput) (*ListProjectResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	projects, err := self.srv.ProjectService.GetProjectsInTeam(ctx, user.ID, input.TeamID)
	if err != nil {
		if errors.Is(err, errdefs.ErrUnauthorized) {
			return nil, huma.Error403Forbidden("Unauthorized")
		}
		if ent.IsNotFound(err) || errors.Is(err, errdefs.ErrNotFound) {
			return nil, huma.Error404NotFound(err.Error())
		}
		log.Error("Error getting projects", "err", err)
		return nil, huma.Error500InternalServerError("Unable to fetch projects")
	}

	resp := &ListProjectResponse{}
	resp.Body.Data = projects
	return resp, nil
}
