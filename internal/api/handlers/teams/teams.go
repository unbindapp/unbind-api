package teams_handler

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
			OperationID: "update-team",
			Summary:     "Update Team",
			Description: "Update a team",
			Path:        "/update",
			Method:      http.MethodPut,
		},
		handlers.UpdateTeam,
	)
}
