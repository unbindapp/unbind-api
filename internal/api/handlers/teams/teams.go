package teams_handler

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
			OperationID: "list-teams",
			Summary:     "List Teams",
			Description: "List all teams the current user is a member of",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		handlers.ListTeams,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-team",
			Summary:     "Get Team",
			Description: "Get a team by ID",
			Path:        "/get",
			Method:      http.MethodGet,
		},
		handlers.GetTeam,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "update-team",
			Summary:     "Update Team",
			Description: "Update a team",
			Path:        "/update",
			Method:      http.MethodPut,
		},
		handlers.UpdateTeam,
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
	log.Error("Error in team handlers", "err", err)
	return huma.Error500InternalServerError("An unexpected error occurred")
}
