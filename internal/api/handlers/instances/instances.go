package instances_handler

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
			OperationID: "list-instances",
			Summary:     "List Instances (Pods)",
			Description: "List all instances (pods) for a service, environment, project, or team. Including health status",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		handlers.ListInstances,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-instance-health",
			Summary:     "Get Instance Health",
			Description: "Get the health/status of instances in a service",
			Path:        "/health",
			Method:      http.MethodGet,
		},
		handlers.GetInstanceHealth,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "restart-instances",
			Summary:     "Restart Instances (Pods)",
			Description: "Restart all instances (pods) for a service",
			Path:        "/restart",
			Method:      http.MethodPut,
		},
		handlers.RestartInstances,
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
