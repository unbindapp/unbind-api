package logs_handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/infrastructure/loki"
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
			OperationID: "query-logs",
			Summary:     "Query Logs",
			Description: "Query logs for a team, project, environment, or service",
			Path:        "/query",
			Method:      http.MethodGet,
		},
		handlers.QueryLogs,
	)

	sse.Register(grp, huma.Operation{
		OperationID: "stream-logs",
		Method:      http.MethodGet,
		Path:        "/stream",
		Summary:     "Stream Logs",
		Description: "Stream logs for a team, project, environment, or service",
	}, map[string]any{
		// Mapping of event type name to Go struct for that event.
		"message": loki.LogEvents{},
	},
		handlers.GetLogsfunc,
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
	log.Error("Unexpected error in logs service", "err", err)
	return huma.Error500InternalServerError("An unexpected error occurred")
}

func (self *HandlerGroup) handleSSEErr(err error, send sse.Sender) {
	if errors.Is(err, errdefs.ErrInvalidInput) {
		send.Data(
			loki.LogEvents{
				MessageType:  loki.LogEventsMessageTypeError,
				ErrorMessage: fmt.Sprintf("invalid input %v", err.Error()),
			},
		)
	}
	if errors.Is(err, errdefs.ErrUnauthorized) {
		send.Data(
			loki.LogEvents{
				MessageType:  loki.LogEventsMessageTypeError,
				ErrorMessage: fmt.Sprintf("unauthorized %v", err.Error()),
			},
		)
	}
	if ent.IsNotFound(err) || errors.Is(err, errdefs.ErrNotFound) {
		send.Data(
			loki.LogEvents{
				MessageType:  loki.LogEventsMessageTypeError,
				ErrorMessage: fmt.Sprintf("entity not found %v", err.Error()),
			},
		)
	}
	log.Error("Unknown error streaming logs", "err", err)
	send.Data(
		loki.LogEvents{
			MessageType:  loki.LogEventsMessageTypeError,
			ErrorMessage: "An unexpected error occurred",
		},
	)
}
