package templates_handler

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
			Summary:     "List Templates",
			Description: "List all unbindg service templates",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		handlers.ListTemplates,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-template",
			Summary:     "Get Template",
			Description: "Get a specific template with its schema",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetTemplateMetadata,
	)
}
