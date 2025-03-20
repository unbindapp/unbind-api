package projects_handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
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
