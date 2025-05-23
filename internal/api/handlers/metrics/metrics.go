package metrics_handler

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
			OperationID: "get-metrics",
			Summary:     "Get Metrics",
			Description: "Get Metrics for a team, project, environment, or service",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetMetrics,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-system-metrics",
			Summary:     "Get System Metrics",
			Description: "Get system level metrics - e.g. Node, Cluster, Region",
			Path:        "/get-system",
			Method:      http.MethodGet,
		},
		handlers.GetNodeMetrics,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-volume-metrics",
			Summary:     "Get Volume Metrics",
			Description: "Get volume level metrics - e.g. PVC",
			Path:        "/get-volume",
			Method:      http.MethodGet,
		},
		handlers.GetVolumeMetrics,
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
	log.Error("Unexpected error in metrics service", "err", err)
	return huma.Error500InternalServerError("An unexpected error occurred")
}
