package environments_handler

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
			OperationID: "get-environment",
			Summary:     "Get Environment",
			Description: "Get an environment by ID",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetEnvironment,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-environments",
			Summary:     "List Environments",
			Description: "List all environments in a project",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		handlers.ListEnvironments,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "create-environment",
			Summary:     "Create Environment",
			Description: "Create a new environment",
			Path:        "/create",
			Method:      http.MethodPost,
		},
		handlers.CreateEnvironment,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-environment",
			Summary:     "Update Environment",
			Description: "Update an existing environment",
			Path:        "/update",
			Method:      http.MethodPut,
		},
		handlers.UpdateEnvironment,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "delete-environment",
			Summary:     "Delete Environment",
			Description: "Delete an environment by ID",
			Path:        "/delete",
			Method:      http.MethodDelete,
		},
		handlers.DeleteEnvironment,
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
	log.Error("Server error in environments handler", "err", err)
	return huma.Error500InternalServerError("Internal server error")
}
