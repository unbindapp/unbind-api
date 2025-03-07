package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/database"
	"github.com/unbindapp/unbind-api/internal/database/repository"
	"github.com/unbindapp/unbind-api/internal/github"
	"github.com/unbindapp/unbind-api/internal/kubeclient"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/middleware"
	"github.com/unbindapp/unbind-api/internal/server"
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
	r.Use(middleware.Logger)
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
	userGrp.UseMiddleware(mw.Authenticate)
	huma.Get(userGrp, "/me", srvImpl.Me)

	ghGroup := huma.NewGroup(api, "/github")
	ghGroup.UseMiddleware(mw.Authenticate)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "app-create",
			Summary:     "Create Github app",
			Description: "Begin the workflow to create a github application",
			Path:        "/app/create",
			Method:      http.MethodGet,
		},
		srvImpl.HandleGithubAppCreate,
	)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "manifest-create",
			Summary:     "Create Github manifest",
			Description: "Create a manifest that the user can use to create a GitHub app",
			Path:        "/app/manifest",
			Method:      http.MethodPost,
		},
		srvImpl.HandleGithubManifestCreate,
	)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "app-connect",
			Summary:     "Connect github app",
			Description: "Connect the new github app to our instance, via manifest code exchange",
			Path:        "/app/connect",
			Method:      http.MethodPost,
		},
		srvImpl.HandleGithubAppConnect,
	)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "app-install",
			Summary:     "App Install Redirect",
			Description: "Redirects to install the github app",
			Path:        "/app/install/{app_id}",
			Method:      http.MethodPost,
		},
		srvImpl.HandleGithubAppInstall,
	)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "list-apps",
			Summary:     "List Github Apps",
			Description: "List all the github apps connected to our instance",
			Path:        "/apps",
			Method:      http.MethodGet,
		},
		srvImpl.HandleListGithubApps,
	)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "list-app-installations",
			Summary:     "List Installations",
			Description: "List all installations for a specific github app",
			Path:        "/app/{app_id}/installations",
			Method:      http.MethodGet,
		},
		srvImpl.HandleListGithubAppInstallations,
	)

	webhookGroup := huma.NewGroup(api, "/webhook")
	huma.Register(
		webhookGroup,
		huma.Operation{
			OperationID: "github-webhook",
			Summary:     "Github Webhook",
			Description: "Handle incoming github webhooks",
			Path:        "/github",
			Method:      http.MethodPost,
		},
		srvImpl.HandleGithubWebhook,
	)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "app-save",
			Summary:     "Save github app",
			Description: "Save github app via code exchange and redirect to installation",
			Path:        "/github/app/save",
			Method:      http.MethodGet,
		},
		srvImpl.HandleGithubAppSave,
	)

	// !
	// huma.Get(api, "/teams", srvImpl.ListTeams)

	// Start the server
	addr := ":8089"
	log.Infof("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
