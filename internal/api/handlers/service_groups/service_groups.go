package servicegroups_handler

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
			OperationID: "list-service-groups",
			Summary:     "List Service Groups",
			Description: "List all services groups in an environment",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		handlers.ListServiceGroups,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-service-group",
			Summary:     "Get Service Group",
			Description: "Get a single service group by ID",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetServiceGroup,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "create-service-group",
			Summary:     "Create Service Group",
			Description: "Create a service group",
			Path:        "/create",
			Method:      http.MethodPost,
		},
		handlers.CreateServiceGroup,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-service-group",
			Summary:     "Update Service Group",
			Description: "Update a service group",
			Path:        "/update",
			Method:      http.MethodPut,
		},
		handlers.UpdateServiceGroup,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "delete-service-group",
			Summary:     "Delete Service Group",
			Description: "Delete a service group",
			Path:        "/delete",
			Method:      http.MethodDelete,
		},
		handlers.DeleteServiceGroup,
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
	log.Error("Unexpected service error", "err", err)
	return huma.Error500InternalServerError("Unexpected error occured")
}
