package template_handler

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
			OperationID: "list-templates",
			Summary:     "List Available Templates",
			Description: "List all available templates",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		handlers.ListTemplates,
	)
}
