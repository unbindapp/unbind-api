package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/database"
	"github.com/unbindapp/unbind-api/internal/database/repository"
	"github.com/unbindapp/unbind-api/internal/github"
	"github.com/unbindapp/unbind-api/internal/kubeclient"
	"github.com/unbindapp/unbind-api/internal/log"
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
	log.Info("🦋 Running migrations...")
	if err := db.Schema.Create(context.TODO()); err != nil {
		log.Fatal("Failed to run migrations", "err", err)
	}
	repo := repository.NewRepository(db)

	// Create kubernetes client
	kubeClient := kubeclient.NewKubeClient(cfg)

	// Implementation
	srvImpl := &server.Server{
		KubeClient: kubeClient,
		Cfg:        cfg,
		Repository: repo,
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
			Scopes:      []string{"openid", "profile", "email", "offline_access"},
		},
		GithubClient: github.NewGithubClient(cfg),
	}

	// New chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	api := humachi.New(r, huma.DefaultConfig("Unbind API", "1.0.0"))

	// Create middleware
	mw, err := middleware.NewMiddleware(cfg, repo, api)
	if err != nil {
		log.Fatalf("Failed to create middleware: %v", err)
	}

	// Add routes
	huma.Get(api, "/healthz", srvImpl.HealthCheck)

	// /auth group
	authGrp := huma.NewGroup(api, "/auth")
	huma.Get(authGrp, "/login", srvImpl.Login)
	huma.Get(authGrp, "/callback", srvImpl.Callback)

	// /user group
	userGrp := huma.NewGroup(api, "/user")
	userGrp.UseMiddleware(mw.Authenticate)
	huma.Get(userGrp, "/me", srvImpl.Me)

	ghGroup := huma.NewGroup(api, "/github")
	ghGroup.UseMiddleware(mw.Authenticate)
	huma.Post(ghGroup, "/app/manifest", srvImpl.GithubManifestCreate)
	huma.Post(ghGroup, "/app/connect", srvImpl.GithubAppConnect)
	huma.Post(ghGroup, "/app/install/{app_id}", srvImpl.GithubAppInstall)

	webhookGroup := huma.NewGroup(api, "/webhook")
	huma.Post(webhookGroup, "/github", srvImpl.HandleGithubWebhook)

	// !
	// huma.Get(api, "/teams", srvImpl.ListTeams)

	// Start the server
	addr := ":8089"
	fmt.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
