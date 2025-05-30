package auth_handler

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
			OperationID: "login",
			Summary:     "Login",
			Description: "Login",
			Path:        "/login",
			Method:      http.MethodPost,
		},
		handlers.LoginSubmit,
	)

	if server.Cfg.EnableDevLoginPage {
		huma.Register(
			grp,
			huma.Operation{
				OperationID: "dev_login",
				Summary:     "Dev Login",
				Description: "Dev Login",
				Path:        "/dev_login",
				Method:      http.MethodGet,
			},
			handlers.DevLogin,
		)
	}

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
