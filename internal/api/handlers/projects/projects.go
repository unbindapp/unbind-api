package projects_handler

import (
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

type HandlerGroup struct {
	srv *server.Server
}

func RegisterHandlers(server *server.Server, grp *huma.Group) {
	handlers := &HandlerGroup{
		srv: server,
	}

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "create-project",
			Summary:     "Create Project",
			Description: "Create a project",
			Path:        "/create",
			Method:      http.MethodPost,
		},
		handlers.CreateProject,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-projects",
			Summary:     "List Projects",
			Description: "List all projects",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		handlers.ListProjects,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-project",
			Summary:     "Get Project",
			Description: "Get a project by ID",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetProject,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-project",
			Summary:     "Update Project",
			Description: "Update a project",
			Path:        "/update",
			Method:      http.MethodPut,
		},
		handlers.UpdateProject,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "delete-project",
			Summary:     "Delete Project",
			Description: "Delete a project",
			Path:        "/delete",
			Method:      http.MethodDelete,
		},
		handlers.DeleteProject,
	)

}

func (self *HandlerGroup) handleErr(err error) error {
	if errors.Is(err, errdefs.ErrInvalidInput) {
		return huma.Error400BadRequest("invalid input", err)
	}
	if errors.Is(err, errdefs.ErrUnauthorized) {
		return huma.Error403Forbidden("Unauthorized")
	}
	if ent.IsNotFound(err) || errors.Is(err, errdefs.ErrNotFound) {
		return huma.Error404NotFound("entity not found", err)
	}
	log.Error("Error updating project", "err", err)
	return huma.Error500InternalServerError("Unable to delete project")
}
