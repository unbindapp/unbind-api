package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-api/config"
	deployments_handler "github.com/unbindapp/unbind-api/internal/api/handlers/deployments"
	environments_handler "github.com/unbindapp/unbind-api/internal/api/handlers/environments"
	github_handler "github.com/unbindapp/unbind-api/internal/api/handlers/github"
	logintmp_handler "github.com/unbindapp/unbind-api/internal/api/handlers/logintmp"
	logs_handler "github.com/unbindapp/unbind-api/internal/api/handlers/logs"
	projects_handler "github.com/unbindapp/unbind-api/internal/api/handlers/projects"
	service_handler "github.com/unbindapp/unbind-api/internal/api/handlers/service"
	system_handler "github.com/unbindapp/unbind-api/internal/api/handlers/system"
	teams_handler "github.com/unbindapp/unbind-api/internal/api/handlers/teams"
	user_handler "github.com/unbindapp/unbind-api/internal/api/handlers/user"
	variables_handler "github.com/unbindapp/unbind-api/internal/api/handlers/variables"
	webhook_handler "github.com/unbindapp/unbind-api/internal/api/handlers/webhook"
	"github.com/unbindapp/unbind-api/internal/api/middleware"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	"github.com/unbindapp/unbind-api/internal/infrastructure/buildkitd.go"
	"github.com/unbindapp/unbind-api/internal/infrastructure/cache"
	"github.com/unbindapp/unbind-api/internal/infrastructure/database"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	deployments_service "github.com/unbindapp/unbind-api/internal/services/deployments"
	environment_service "github.com/unbindapp/unbind-api/internal/services/environment"
	logs_service "github.com/unbindapp/unbind-api/internal/services/logs"
	project_service "github.com/unbindapp/unbind-api/internal/services/project"
	service_service "github.com/unbindapp/unbind-api/internal/services/service"
	system_service "github.com/unbindapp/unbind-api/internal/services/system"
	team_service "github.com/unbindapp/unbind-api/internal/services/team"
	"github.com/valkey-io/valkey-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/oauth2"
)

func NewHumaConfig(title, version string) huma.Config {
	schemaPrefix := "#/components/schemas/"
	schemasPath := "/schemas"

	registry := huma.NewMapRegistry(schemaPrefix, huma.DefaultSchemaNamer)

	return huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:   title,
				Version: version,
			},
			Components: &huma.Components{
				Schemas: registry,
			},
		},
		OpenAPIPath:   "/openapi",
		DocsPath:      "/docs",
		SchemasPath:   schemasPath,
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
		// * Remove the $schma field
		// CreateHooks: []func(huma.Config) huma.Config{
		// 	func(c huma.Config) huma.Config {
		// 		// Add a link transformer to the API. This adds `Link` headers and
		// 		// puts `$schema` fields in the response body which point to the JSON
		// 		// Schema that describes the response structure.
		// 		// This is a create hook so we get the latest schema path setting.
		// 		linkTransformer := huma.NewSchemaLinkTransformer(schemaPrefix, c.SchemasPath)
		// 		c.OpenAPI.OnAddOperation = append(c.OpenAPI.OnAddOperation, linkTransformer.OnAddOperation)
		// 		c.Transformers = append(c.Transformers, linkTransformer.Transform)
		// 		return c
		// 	},
		// },
	}
}

func startAPI(cfg *config.Config) {
	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalCh
		slog.Info("Received shutdown signal", "signal", sig)
		cancel() // This will propagate cancellation to all derived contexts
	}()

	// Initialize valkey (redis)
	valkeyClient, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{cfg.ValkeyURL}})
	if err != nil {
		log.Fatal("Failed to create valkey client", "err", err)
	}
	defer valkeyClient.Close()

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
	if err := db.Schema.Create(ctx); err != nil {
		log.Fatal("Failed to run migrations", "err", err)
	}
	repo := repositories.NewRepositories(db)

	// Create kubernetes client
	kubeClient := k8s.NewKubeClient(cfg)

	// Create github client
	githubClient := github.NewGithubClient(cfg.GithubURL, cfg)

	// Create deployment controller
	deploymentController := deployctl.NewDeploymentController(ctx, cancel, cfg, kubeClient, valkeyClient, repo, githubClient)

	// Buildkit settings manager
	buildkitSettings := buildkitd.NewBuildkitSettingsManager(repo, kubeClient)

	// Bootstrap
	bootstrapper := &Bootstrapper{
		repos:                   repo,
		buildkitSettingsManager: buildkitSettings,
	}
	if err := bootstrapper.Sync(ctx); err != nil {
		log.Errorf("Failed to sync system settings: %v", err)
	}

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
			RedirectURL: fmt.Sprintf("%s/auth/callback", cfg.ExternalAPIURL),
			Scopes:      []string{"openid", "profile", "email", "offline_access", "groups"},
		},
		GithubClient:         githubClient,
		StringCache:          cache.NewStringCache(valkeyClient, "unbind"),
		HttpClient:           &http.Client{},
		DeploymentController: deploymentController,
		TeamService:          team_service.NewTeamService(repo, kubeClient),
		ProjectService:       project_service.NewProjectService(repo, kubeClient),
		ServiceService:       service_service.NewServiceService(cfg, repo, githubClient, kubeClient, deploymentController),
		EnvironmentService:   environment_service.NewEnvironmentService(repo, kubeClient),
		LogService:           logs_service.NewLogsService(repo, kubeClient),
		DeploymentService:    deployments_service.NewDeploymentService(repo, deploymentController, githubClient),
		SystemService:        system_service.NewSystemService(cfg, repo, buildkitSettings),
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

	config := NewHumaConfig("Unbind API", "1.0.0")
	config.DocsPath = ""
	config.OpenAPI.Servers = []*huma.Server{
		{
			URL: cfg.ExternalAPIURL,
		},
	}
	api := humachi.New(r, config)

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html>
			<html>
				<head>
					<title>API Reference</title>
					<meta charset="utf-8" />
					<meta
						name="viewport"
						content="width=device-width, initial-scale=1" />
				</head>
				<body>
					<script
						id="api-reference"
						data-url="/openapi.json"></script>
					<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
				</body>
			</html>`))
	})

	// Create middleware
	mw, err := middleware.NewMiddleware(cfg, repo, api)
	if err != nil {
		log.Warnf("Failed to create middleware: %v", err)
	}

	type HealthResponse struct {
		Body struct {
			Status string `json:"status"`
		}
	}

	huma.Register(
		api,
		huma.Operation{
			OperationID: "health",
			Summary:     "Health Check",
			Description: "Check if the API is healthy",
			Path:        "/health",
			Method:      http.MethodGet,
			Tags:        []string{"Meta"},
		},
		func(ctx context.Context, i *server.EmptyInput) (*HealthResponse, error) {

			healthResponse := &HealthResponse{}
			healthResponse.Body.Status = "ok"
			return healthResponse, nil
		},
	)

	// /auth group
	authGroup := huma.NewGroup(api, "/auth")
	authGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Auth"}
		next(op)
	})
	logintmp_handler.RegisterHandlers(srvImpl, authGroup)

	// /system group
	systemGroup := huma.NewGroup(api, "/system")
	systemGroup.UseMiddleware(mw.Authenticate)
	systemGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"System"}
		next(op)
	})
	system_handler.RegisterHandlers(srvImpl, systemGroup)

	// /users group
	userGroup := huma.NewGroup(api, "/users")
	userGroup.UseMiddleware(mw.Authenticate)
	userGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Users"}
		next(op)
	})
	user_handler.RegisterHandlers(srvImpl, userGroup)

	ghGroup := huma.NewGroup(api, "/github")
	ghGroup.UseMiddleware(mw.Authenticate)
	ghGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"GitHub"}
		next(op)
	})
	github_handler.RegisterHandlers(srvImpl, ghGroup)

	webhookGroup := huma.NewGroup(api, "/webhook")
	webhookGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Webhook"}
		next(op)
	})
	webhook_handler.RegisterHandlers(srvImpl, webhookGroup)

	// /teams
	teamsGroup := huma.NewGroup(api, "/teams")
	teamsGroup.UseMiddleware(mw.Authenticate)
	teamsGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Teams"}
		next(op)
	})
	teams_handler.RegisterHandlers(srvImpl, teamsGroup)

	// /projects group
	projectsGroup := huma.NewGroup(api, "/projects")
	projectsGroup.UseMiddleware(mw.Authenticate)
	projectsGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Projects"}
		next(op)
	})
	projects_handler.RegisterHandlers(srvImpl, projectsGroup)

	// /environments group
	environmentsGroup := huma.NewGroup(api, "/environments")
	environmentsGroup.UseMiddleware(mw.Authenticate)
	environmentsGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Environments"}
		next(op)
	})
	environments_handler.RegisterHandlers(srvImpl, environmentsGroup)

	// /services group
	servicesGroup := huma.NewGroup(api, "/services")
	servicesGroup.UseMiddleware(mw.Authenticate)
	servicesGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Services"}
		next(op)
	})
	service_handler.RegisterHandlers(srvImpl, servicesGroup)

	// /variables group
	variablesGroup := huma.NewGroup(api, "/variables")
	variablesGroup.UseMiddleware(mw.Authenticate)
	variablesGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Variables"}
		next(op)
	})
	variables_handler.RegisterHandlers(srvImpl, variablesGroup)

	// /logs group
	logsGroup := huma.NewGroup(api, "/logs")
	logsGroup.UseMiddleware(mw.Authenticate)
	logsGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Logs"}
		next(op)
	})
	logs_handler.RegisterHandlers(srvImpl, logsGroup)

	// /deployments group
	deploymentsGroup := huma.NewGroup(api, "/deployments")
	deploymentsGroup.UseMiddleware(mw.Authenticate)
	deploymentsGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Deployments"}
		next(op)
	})
	deployments_handler.RegisterHandlers(srvImpl, deploymentsGroup)

	// Start the server
	addr := ":8089"
	log.Infof("Starting server on %s\n", addr)

	h2s := &http2.Server{}

	server := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(r, h2s),
	}

	// Start deployment queue processeor
	deploymentController.StartAsync()

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for context cancellation (from signal handler)
	<-ctx.Done()
	log.Info("Shutting down server...")

	// Create a shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown the HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Info("Server gracefully stopped")
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Overload()
	if err != nil {
		log.Warn("Error loading .env file:", err)
	}

	cfg := config.NewConfig()
	startAPI(cfg)
}
