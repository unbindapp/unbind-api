package template_handler

import (
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
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
			OperationID: "get-template",
			Summary:     "Get Template",
			Description: "Get a template by ID",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetTemplateByID,
	)
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
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "deploy-template",
			Summary:     "Deploy Template",
			Description: "Deploy a template",
			Path:        "/deploy",
			Method:      http.MethodPost,
		},
		handlers.DeployTemplate,
	)
}

func (self *HandlerGroup) handleErr(err error) error {
	if errors.Is(err, errdefs.ErrInvalidInput) {
		return huma.Error400BadRequest("invalid input", err)
	}
	if errors.Is(err, errdefs.ErrUnauthorized) {
		return huma.Error403Forbidden("Unauthorized")
	}
	if ent.IsNotFound(err) || errors.Is(err, errdefs.ErrNotFound) {
		return huma.Error404NotFound("entity not found", err)
	}
	log.Error("Unexpected service error", "err", err)
	return huma.Error500InternalServerError("Unexpected error occured")
}
