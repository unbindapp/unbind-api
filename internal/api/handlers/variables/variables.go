package variables_handler

import (
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
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
			OperationID: "list-variables",
			Summary:     "List Variables",
			Description: "List variables for a service, environment, project, or team",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		handlers.ListVariables,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-variables",
			Summary:     "Update Variables",
			Description: "Create or update variables for a service, environment, project, or team by key",
			Path:        "/update",
			Method:      http.MethodPost,
		},
		handlers.UpdateVariables,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "delete-variables",
			Summary:     "Delete Variables",
			Description: "Delete variables for a service, environment, project, or team",
			Path:        "/delete",
			Method:      http.MethodDelete,
		},
		handlers.DeleteVariables,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-available-references",
			Summary:     "List Available Variable References",
			Description: "List all items that can be references as variables in service configurations",
			Path:        "/references/available",
			Method:      http.MethodGet,
		},
		handlers.ListReferenceableVariables,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "read-variable-reference",
			Summary:     "Read a referenced value",
			Description: "Read a referenced value for a variable",
			Path:        "/references/get",
			Method:      http.MethodGet,
		},
		handlers.ResolveVariableReference,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "create-variable-reference",
			Summary:     "Create Variable Reference",
			Description: "Create a variable reference for a service, environment, project, or team",
			Path:        "/references/create",
			Method:      http.MethodPost,
		},
		handlers.CreateVariableReference,
	)
}

func handleVariablesErr(err error) error {
	if errors.Is(err, errdefs.ErrUnauthorized) {
		return huma.Error403Forbidden("Unauthorized")
	}
	if errors.Is(err, errdefs.ErrNotFound) {
		return huma.Error404NotFound(err.Error())
	}
	log.Error("Error getting variables", "err", err)
	return huma.Error500InternalServerError("Unable to retrieve variables")
}
