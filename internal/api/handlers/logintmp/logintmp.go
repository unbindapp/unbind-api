package logintmp_handler

import "github.com/unbindapp/unbind-api/internal/api/server"

type HandlerGroup struct {
	srv *server.Server
}

// ! TODO - get rid of this one, it's just for early temporary oauth testing

func NewHandlerGroup(server *server.Server) *HandlerGroup {
	return &HandlerGroup{
		srv: server,
	}
}
