package github_handler

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
			OperationID: "app-create",
			Summary:     "Create App",
			Description: "Begin the workflow to create a GitHub application",
			Path:        "/app/create",
			Method:      http.MethodGet,
		},
		handlers.HandleGithubAppCreate,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "get-github-app",
			Summary:     "Get App",
			Description: "Get the GitHub app details",
			Path:        "/app/get",
			Method:      http.MethodGet,
		},
		handlers.HandleGetGithubApp,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-apps",
			Summary:     "List Apps",
			Description: "List all the GitHub apps connected to our instance",
			Path:        "/apps",
			Method:      http.MethodGet,
		},
		handlers.HandleListGithubApps,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-app-installations",
			Summary:     "List Installations",
			Description: "List all github app installations.",
			Path:        "/installations",
			Method:      http.MethodGet,
		},
		handlers.HandleListGithubAppInstallations,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-admin-organizations",
			Summary:     "List Admin Organizations",
			Description: "List all admin organizations for a specific user installation, invalid for 'Organization' installations.",
			Path:        "/installation/{installation_id}/organizations",
			Method:      http.MethodGet,
		},
		handlers.HandleListGithubAdminOrganizations,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "list-admin-repos",
			Summary:     "List Repositories",
			Description: "List all repositories the user has admin access of.",
			Path:        "/repositories",
			Method:      http.MethodGet,
		},
		handlers.HandleListGithubAdminRepositories,
	)
	huma.Register(
		grp,
		huma.Operation{
			OperationID: "repo-detail",
			Summary:     "Repository Detail",
			Description: "Get details of a repository (branch, tags, etc.)",
			Path:        "/repositories/info",
			Method:      http.MethodGet,
		},
		handlers.HandleGetGithubRepositoryDetail,
	)
}
