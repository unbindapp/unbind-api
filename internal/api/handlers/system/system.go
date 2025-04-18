package system_handler

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
			OperationID: "get-system-information",
			Summary:     "Get System Information",
			Description: "Get system information such as external load balancer IP for DNS configurations.",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetSystemInformation,
	)

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-buildkit-settings",
			Summary:     "Update Buildkit Settings",
			Description: "Update buildkit settings such as replicas and parallelism.",
			Path:        "/buildkit/update",
			Method:      http.MethodPut,
		},
		handlers.UpdateBuildkitSettings,
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
	log.Error("Unexpected system error", "err", err)
	return huma.Error500InternalServerError("Unexpected error occured")
}
