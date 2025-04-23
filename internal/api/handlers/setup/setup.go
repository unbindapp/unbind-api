package setup_handler

// Handles initial bootstrapping requirements

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
			OperationID: "get-setup-status",
			Summary:     "Get Setup Status",
			Description: "Return if Unbind has been bootstrapped or not",
			Path:        "/status",
			Method:      http.MethodGet,
		},
		handlers.GetStatus,
	)

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "create-bootstrap-user",
			Summary:     "Create Bootstrap User",
			Description: "Create the initial user for Unbind",
			Path:        "/create-user",
			Method:      http.MethodPost,
		},
		handlers.CreateUser,
	)
}
