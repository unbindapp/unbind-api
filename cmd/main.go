package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/database"
	"github.com/unbindapp/unbind-api/internal/database/repository"
	"github.com/unbindapp/unbind-api/internal/kubeclient"
	"github.com/unbindapp/unbind-api/internal/middleware"
	"github.com/unbindapp/unbind-api/internal/server"
	"golang.org/x/oauth2"
)

func main() {
	godotenv.Load()

	// Initialize config
	cfg := config.NewConfig()

	// Load database
	dbConnInfo, err := database.GetSqlDbConn(cfg, false)
	if err != nil {
		log.Fatalf("Failed to get database connection info: %v", err)
	}
	// Initialize ent client
	db, err := database.NewEntClient(dbConnInfo)
	if err != nil {
		log.Fatalf("Failed to create ent client: %v", err)
	}
	repo := repository.NewRepository(db)

	// Create kubernetes client
	kubeClient := kubeclient.NewKubeClient(cfg)

	// Implementation
	srvImpl := &server.Server{
		KubeClient: kubeClient,
		Cfg:        cfg,
		// Create an OAuth2 configuration using the Dex
		OauthConfig: &oauth2.Config{
			ClientID:     cfg.DexClientID,
			ClientSecret: cfg.DexClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  cfg.DexIssuerUrlExternal + "/auth",
				TokenURL: cfg.DexIssuerURL + "/token",
			},
			// ! TODO - adjust redirect when necessary
			RedirectURL: "http://localhost:8089/auth/callback",
			Scopes:      []string{"openid", "profile", "email"},
		},
	}

	// Create middleware
	mw, err := middleware.NewMiddleware(cfg, repo)
	if err != nil {
		log.Fatalf("Failed to create middleware: %v", err)
	}

	// New chi router
	r := chi.NewRouter()
	r.Use(mw.Logger)
	api := humachi.New(r, huma.DefaultConfig("Unbind API", "1.0.0"))

	// Add routes
	huma.Get(api, "/healthz", srvImpl.HealthCheck)

	r.Group(func(r chi.Router) {
		huma.Get(api, "/auth/login", srvImpl.Login)
		huma.Get(api, "/auth/callback", srvImpl.Callback)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(mw.Authenticate)

		r.Get("/api/me", func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("user").(*ent.User)
			json.NewEncoder(w).Encode(user)
		})
	})
	huma.Get(api, "/teams", srvImpl.ListTeams)

	// Start the server
	addr := ":8089"
	fmt.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
