package unbindwebhooks_handler

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
			OperationID: "create-webhook",
			Summary:     "Create Webhook",
			Description: "Create a new webhook for a team or project",
			Path:        "/create",
			Method:      http.MethodPost,
		},
		handlers.CreateWebhook,
	)

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-webhook",
			Summary:     "Update Webhook",
			Description: "Update an existing webhook",
			Path:        "/update",
			Method:      http.MethodPut,
		},
		handlers.UpdateWebhook,
	)

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "delete-webhook",
			Summary:     "Delete Webhook",
			Description: "Delete a webhook",
			Path:        "/delete",
			Method:      http.MethodDelete,
		},
		handlers.DeleteWebhook,
	)

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-webhooks",
			Summary:     "List Webhooks",
			Description: "List webhooks for a team or project",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		handlers.ListWebhooks,
	)

	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-webhook",
			Summary:     "Get Webhook",
			Description: "Get a single webhook by ID",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetWebhook,
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
	log.Error("Unexpected unbindwebhooks error", "err", err)
	return huma.Error500InternalServerError("Unexpected error occured")
}
