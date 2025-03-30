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

	sse.Register(grp, huma.Operation{
		OperationID: "stream-logs",
		Method:      http.MethodGet,
		Path:        "/stream",
		Summary:     "Stream Logs",
		Description: "Stream logs for a team, project, environment, or service",
	}, map[string]any{
		// Mapping of event type name to Go struct for that event.
		"message":      loki.LogEvents{},
		"errorMessage": loki.LogsError{},
	},
		handlers.GetLogsfunc,
	)
}

func (self *HandlerGroup) handleErr(err error, send sse.Sender) {
	if errors.Is(err, errdefs.ErrInvalidInput) {
		send.Data(
			loki.LogsError{
				Code:    400,
				Message: fmt.Sprintf("invalid input %v", err.Error()),
			},
		)
	}
	if errors.Is(err, errdefs.ErrUnauthorized) {
		send.Data(
			loki.LogsError{
				Code:    403,
				Message: fmt.Sprintf("unauthorized %v", err.Error()),
			},
		)
	}
	if ent.IsNotFound(err) || errors.Is(err, errdefs.ErrNotFound) {
		send.Data(
			loki.LogsError{
				Code:    404,
				Message: fmt.Sprintf("entity not found %v", err.Error()),
			},
		)
	}
	log.Error("Unknown error streaming logs", "err", err)
	send.Data(
		loki.LogsError{
			Code:    500,
			Message: fmt.Sprintf("unknown error streaming logs"),
		},
	)
}
