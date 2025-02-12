package server

import (
	"net/http"
)

// Server implements generated.ServerInterface
type Server struct {
	// Add fields for DB connections, config, etc. as needed
}

// HealthCheck is your /health endpoint
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
