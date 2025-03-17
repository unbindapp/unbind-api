package user_handler

import "github.com/unbindapp/unbind-api/internal/api/server"

type HandlerGroup struct {
	srv *server.Server
}

func NewHandlerGroup(server *server.Server) *HandlerGroup {
	return &HandlerGroup{
		srv: server,
	}
}
