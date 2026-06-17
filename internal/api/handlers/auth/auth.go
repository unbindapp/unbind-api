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
	handlers := &HandlerGroup{srv: server}

	huma.Register(grp, huma.Operation{
		OperationID: "login",
		Summary:     "Login",
		Description: "Authenticate with email and password, returns session cookies",
		Path:        "/login",
		Method:      http.MethodPost,
	}, handlers.Login)

	huma.Register(grp, huma.Operation{
		OperationID: "logout",
		Summary:     "Logout",
		Description: "Revoke the session and clear cookies",
		Path:        "/logout",
		Method:      http.MethodPost,
	}, handlers.Logout)
}
