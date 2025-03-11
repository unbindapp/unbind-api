package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/database"
	"github.com/unbindapp/unbind-api/internal/database/repository"
	"github.com/unbindapp/unbind-api/internal/github"
	"github.com/unbindapp/unbind-api/internal/kubeclient"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/middleware"
	"github.com/unbindapp/unbind-api/internal/server"
	github_handler "github.com/unbindapp/unbind-api/internal/server/github"
	teams_handler "github.com/unbindapp/unbind-api/internal/server/teams"
	user_handler "github.com/unbindapp/unbind-api/internal/server/user"
	webhook_handler "github.com/unbindapp/unbind-api/internal/server/webhook"
	"github.com/valkey-io/valkey-go"
	"golang.org/x/oauth2"
)

func startAPI(cfg *config.Config) {
	// Initialize valkey (redis)
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{cfg.ValkeyURL}})
	if err != nil {
		log.Fatal("Failed to create valkey client", "err", err)
	}
	defer client.Close()

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
			RedirectURL: fmt.Sprintf("%s/auth/callback", cfg.ExternalURL),
			Scopes:      []string{"openid", "profile", "email", "offline_access"},
		},
		GithubClient: github.NewGithubClient(cfg),
		StringCache:  database.NewStringCache(client, "unbind"),
		HttpClient:   &http.Client{},
	}

	// New chi router
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://app.unbind.app",
		},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	api := humachi.New(r, huma.DefaultConfig("Unbind API", "1.0.0"))

	// Create middleware
	mw, err := middleware.NewMiddleware(cfg, repo, api)
	if err != nil {
		log.Warnf("Failed to create middleware: %v", err)
	}

	// Add routes
	huma.Get(api, "/healthz", srvImpl.HealthCheck)

	// /auth group
	authGrp := huma.NewGroup(api, "/auth")
	huma.Get(authGrp, "/login", srvImpl.Login)
	huma.Get(authGrp, "/callback", srvImpl.Callback)

	// /user group
	userGrp := huma.NewGroup(api, "/user")
	userHandlers := user_handler.NewHandlerGroup(srvImpl)
	userGrp.UseMiddleware(mw.Authenticate)
	huma.Get(userGrp, "/me", userHandlers.Me)

	ghGroup := huma.NewGroup(api, "/github")
	githubHandlers := github_handler.NewHandlerGroup(srvImpl)
	ghGroup.UseMiddleware(mw.Authenticate)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "app-create",
			Summary:     "Create GitHub App",
			Description: "Begin the workflow to create a GitHub application",
			Path:        "/app/create",
			Method:      http.MethodGet,
		},
		githubHandlers.HandleGithubAppCreate,
	)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "list-apps",
			Summary:     "List Github Apps",
			Description: "List all the GitHub apps connected to our instance",
			Path:        "/apps",
			Method:      http.MethodGet,
		},
		githubHandlers.HandleListGithubApps,
	)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "list-app-installations",
			Summary:     "List Installations",
			Description: "List all github app installations.",
			Path:        "/installations",
			Method:      http.MethodGet,
		},
		githubHandlers.HandleListGithubAppInstallations,
	)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "list-admin-organizations",
			Summary:     "List Admin Organizations for an User Installation",
			Description: "List all admin organizations for a specific user installation, invalid for 'Organization' installations.",
			Path:        "/installation/{installation_id}/organizations",
			Method:      http.MethodGet,
		},
		githubHandlers.HandleListGithubAdminOrganizations,
	)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "list-admin-repos",
			Summary:     "List Repositories",
			Description: "List all repositories the user has admin access of.",
			Path:        "/repositories",
			Method:      http.MethodGet,
		},
		githubHandlers.HandleListGithubAdminRepositories,
	)

	webhookGroup := huma.NewGroup(api, "/webhook")
	webhookHandlers := webhook_handler.NewHandlerGroup(srvImpl)
	huma.Register(
		webhookGroup,
		huma.Operation{
			OperationID: "github-webhook",
			Summary:     "Github Webhook",
			Description: "Handle incoming GitHub webhooks",
			Path:        "/github",
			Method:      http.MethodPost,
		},
		webhookHandlers.HandleGithubWebhook,
	)
	huma.Register(
		webhookGroup,
		huma.Operation{
			OperationID: "app-save",
			Summary:     "Save GitHub app",
			Description: "Save GitHub app via code exchange and redirect to installation",
			Path:        "/github/app/save",
			Method:      http.MethodGet,
		},
		webhookHandlers.HandleGithubAppSave,
	)

	// /teams group
	teamsGrp := huma.NewGroup(api, "/teams")
	teamHandlers := teams_handler.NewHandlerGroup(srvImpl)
	teamsGrp.UseMiddleware(mw.Authenticate)
	huma.Register(
		teamsGrp,
		huma.Operation{
			OperationID: "list-teams",
			Summary:     "List Teams",
			Description: "List all teams the current user is a member of",
			Path:        "/",
			Method:      http.MethodGet,
		},
		teamHandlers.ListTeams,
	)

	// Start the server
	addr := ":8089"
	log.Infof("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
