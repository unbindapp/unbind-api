package system_handler

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
			OperationID: "get-system-information",
			Summary:     "Get System Information",
			Description: "Get system information such as external load balancer IP for DNS configurations.",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetSystemInformation,
	)
}
