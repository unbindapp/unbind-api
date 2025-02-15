package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/kubeclient"
	"github.com/unbindapp/unbind-api/internal/server"
)

func main() {
	godotenv.Load()

	cfg := config.NewConfig()
	// Initialize config

	// Create kubernetes client
	kubeClient := kubeclient.NewKubeClient(cfg)

	// Implementation
	srvImpl := &server.Server{
		KubeClient: kubeClient,
	}

	// New chi router
	r := chi.NewRouter()
	api := humachi.New(r, huma.DefaultConfig("Unbind API", "1.0.0"))

	// Add routes
	huma.Get(api, "/healthz", srvImpl.HealthCheck)

	// ! TODO - auth stuff
	huma.Get(api, "/teams", srvImpl.ListTeams)

	// Start the server
	addr := ":8089"
	fmt.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
