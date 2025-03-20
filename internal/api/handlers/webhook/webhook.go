package webhook_handler

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
			OperationID: "github-webhook",
			Summary:     "Github Webhook",
			Description: "Handle incoming GitHub webhooks",
			Path:        "/github",
			Method:      http.MethodPost,
		},
		handlers.HandleGithubWebhook,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "app-save",
			Summary:     "Save GitHub App",
			Description: "Save GitHub app via code exchange and redirect to installation",
			Path:        "/github/app/save",
			Method:      http.MethodGet,
		},
		handlers.HandleGithubAppSave,
	)
}
