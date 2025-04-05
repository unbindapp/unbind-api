package deployments_handler

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
			OperationID: "list-deployments",
			Summary:     "List Deployments",
			Description: "List deployments for a specific service",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		handlers.ListDeployments,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-deployment",
			Summary:     "Get Deployment",
			Description: "Get a specific deployment by ID",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetDeploymentByID,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "trigger-deployment",
			Summary:     "Trigger Deployment",
			Description: "Trigger a new deployment for a service manually",
			Path:        "/create",
			Method:      http.MethodPost,
		},
		handlers.CreateDeployment,
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
	log.Error("Unexpected error in deployment handler", "err", err)
	return huma.Error500InternalServerError("Internal server error")
}
