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

type DeleteProjectInput struct {
	server.BaseAuthInput
	Body struct {
		TeamID    uuid.UUID `json:"team_id" required:"true"`
		ProjectID uuid.UUID `json:"project_id" required:"true"`
	}
}

type DeleteProjectResponse struct {
	Body struct {
		Data struct {
			ID      uuid.UUID `json:"id"`
			Deleted bool      `json:"deleted"`
		} `json:"data"`
	}
}

func (self *HandlerGroup) DeleteProject(ctx context.Context, input *DeleteProjectInput) (*DeleteProjectResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}

	err := self.srv.ProjectService.DeleteProject(ctx, user.ID, &project_service.DeleteProjectInput{
		TeamID:    input.Body.TeamID,
		ProjectID: input.Body.ProjectID,
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
		return nil, huma.Error500InternalServerError("Unable to delete project")
	}

	resp := &DeleteProjectResponse{}
	resp.Body.Data = struct {
		ID      uuid.UUID `json:"id"`
		Deleted bool      `json:"deleted"`
	}{
		ID:      input.Body.ProjectID,
		Deleted: true,
	}
	return resp, nil
}
