package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/unbindapp/unbind-api/internal/generated"
)

// ListProjects handles GET /project
func (s *Server) ListProjects(w http.ResponseWriter, r *http.Request, params generated.ListProjectsParams) {
	projects := []generated.Project{
		{Id: 1, Name: "ABC"},
		{Id: 2, Name: "DEF"},
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, projects)
}
