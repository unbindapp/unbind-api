package service_handler

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
			OperationID: "create-service",
			Summary:     "Create Service",
			Description: "Create a service",
			Path:        "/create",
			Method:      http.MethodPost,
		},
		handlers.CreateService,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-service",
			Summary:     "Update Service",
			Description: "Update a service",
			Path:        "/update",
			Method:      http.MethodPut,
		},
		handlers.UpdateService,
	)
}

func (self *HandlerGroup) handleErr(err error) error {
	if errors.Is(err, errdefs.ErrInvalidInput) {
		return huma.Error400BadRequest(err.Error())
	}
	if errors.Is(err, errdefs.ErrUnauthorized) {
		return huma.Error403Forbidden("Unauthorized")
	}
	if ent.IsNotFound(err) || errors.Is(err, errdefs.ErrNotFound) {
		return huma.Error404NotFound(err.Error())
	}
	log.Error("Unexpected service error", "err", err)
	return huma.Error500InternalServerError("Unexpected error occured")
}
