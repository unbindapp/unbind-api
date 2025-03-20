package logintmp_handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
)

type HandlerGroup struct {
	srv *server.Server
}

// ! TODO - get rid of this one, it's just for early temporary oauth testing

func RegisterHandlers(server *server.Server, grp *huma.Group) {
	handlers := &HandlerGroup{
		srv: server,
	}

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "login",
			Summary:     "Login",
			Description: "Login",
			Path:        "/login",
			Method:      http.MethodGet,
		},
		handlers.Login)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "callback",
			Summary:     "Callback",
			Description: "Callback",
			Path:        "/callback",
			Method:      http.MethodGet,
		},
		handlers.Callback,
	)
}
