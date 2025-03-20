package user_handler

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
			OperationID: "me",
			Summary:     "Get Current User",
			Description: "Get the current user",
			Path:        "/me",
			Method:      http.MethodGet,
		},
		handlers.Me,
	)
}
