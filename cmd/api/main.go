package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-co-op/gocron/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/schema"
	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-api/config"
	auth_handler "github.com/unbindapp/unbind-api/internal/api/handlers/auth"
	deployments_handler "github.com/unbindapp/unbind-api/internal/api/handlers/deployments"
	environments_handler "github.com/unbindapp/unbind-api/internal/api/handlers/environments"
	github_handler "github.com/unbindapp/unbind-api/internal/api/handlers/github"
	instances_handler "github.com/unbindapp/unbind-api/internal/api/handlers/instances"
	logs_handler "github.com/unbindapp/unbind-api/internal/api/handlers/logs"
	metrics_handler "github.com/unbindapp/unbind-api/internal/api/handlers/metrics"
	projects_handler "github.com/unbindapp/unbind-api/internal/api/handlers/projects"
	service_handler "github.com/unbindapp/unbind-api/internal/api/handlers/service"
	setup_handler "github.com/unbindapp/unbind-api/internal/api/handlers/setup"
	storage_handler "github.com/unbindapp/unbind-api/internal/api/handlers/storage"
	system_handler "github.com/unbindapp/unbind-api/internal/api/handlers/system"
	teams_handler "github.com/unbindapp/unbind-api/internal/api/handlers/teams"
	template_handler "github.com/unbindapp/unbind-api/internal/api/handlers/templates"
	unbindwebhooks_handler "github.com/unbindapp/unbind-api/internal/api/handlers/unbindwebhooks"
	user_handler "github.com/unbindapp/unbind-api/internal/api/handlers/user"
	variables_handler "github.com/unbindapp/unbind-api/internal/api/handlers/variables"
	webhook_handler "github.com/unbindapp/unbind-api/internal/api/handlers/webhook"
	"github.com/unbindapp/unbind-api/internal/api/middleware"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	"github.com/unbindapp/unbind-api/internal/infrastructure/buildkitd"
	"github.com/unbindapp/unbind-api/internal/infrastructure/cache"
	"github.com/unbindapp/unbind-api/internal/infrastructure/database"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/infrastructure/loki"
	"github.com/unbindapp/unbind-api/internal/infrastructure/prometheus"
	"github.com/unbindapp/unbind-api/internal/infrastructure/registry"
	"github.com/unbindapp/unbind-api/internal/infrastructure/updater"
	"github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	deployments_service "github.com/unbindapp/unbind-api/internal/services/deployments"
	environment_service "github.com/unbindapp/unbind-api/internal/services/environment"
	instance_service "github.com/unbindapp/unbind-api/internal/services/instances"
	logs_service "github.com/unbindapp/unbind-api/internal/services/logs"
	metric_service "github.com/unbindapp/unbind-api/internal/services/metrics"
	project_service "github.com/unbindapp/unbind-api/internal/services/project"
	service_service "github.com/unbindapp/unbind-api/internal/services/service"
	storage_service "github.com/unbindapp/unbind-api/internal/services/storage"
	system_service "github.com/unbindapp/unbind-api/internal/services/system"
	team_service "github.com/unbindapp/unbind-api/internal/services/team"
	templates_service "github.com/unbindapp/unbind-api/internal/services/templates"
	variables_service "github.com/unbindapp/unbind-api/internal/services/variables"
	webhooks_service "github.com/unbindapp/unbind-api/internal/services/webooks"
	"github.com/unbindapp/unbind-api/pkg/databases"
	"github.com/valkey-io/valkey-go"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/oauth2"
)

var Version = "development"
var BuildImage = "ghcr.io/unbindapp/unbind-builder:latest"

// Adding a format for form data
var decoder = schema.NewDecoder()
var urlEncodedFormat = huma.Format{
	Marshal: nil,
	Unmarshal: func(data []byte, v any) error {
		values, err := url.ParseQuery(string(data))
		if err != nil {
			return err
		}

		// WARNING: Dirty workaround!
		// During validation, Huma first parses the body into []any, map[string]any or equivalent for easy validation,
		// before parsing it into the target struct.
		// However, gorilla/schema requires a struct for decoding, so we need to map `url.Values` to a
		// `map[string]any` if this happens.
		// See: https://github.com/danielgtaylor/huma/blob/main/huma.go#L1264
		if vPtr, ok := v.(*interface{}); ok {
			m := map[string]any{}
			for k, v := range values {
				if len(v) > 1 {
					m[k] = v
				} else if len(v) == 1 {
					m[k] = v[0]
				}
			}
			*vPtr = m
			return nil
		}

		// `v` is a struct, try decode normally
		return decoder.Decode(v, values)
	},
}

func NewHumaConfig(title, version string) huma.Config {
	schemaPrefix := "#/components/schemas/"
	schemasPath := "/schemas"

	registry := huma.NewMapRegistry(schemaPrefix, huma.DefaultSchemaNamer)

	cfg := huma.Config{
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
	cfg.Formats["application/x-www-form-urlencoded"] = urlEncodedFormat
	cfg.Formats["x-www-form-urlencoded"] = urlEncodedFormat

	return cfg
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
	log.Infof("Using PostgreSQL database %s@%s:%d", cfg.PostgresUser, cfg.PostgresHost, cfg.PostgresPort)
	// Initialize ent client
	db, err := database.NewEntClient(dbConnInfo)
	if err != nil {
		log.Fatalf("Failed to create ent client: %v", err)
	}
	log.Info("🦋 Running migrations...")

	// Auto-apply migrations from migrations directory
	migrationDir := "/app/migrations"
	if _, err := os.Stat(migrationDir); err != nil {
		_, thisFile, _, _ := runtime.Caller(0)
		migrationDir = filepath.Join(filepath.Dir(thisFile), "../../ent/migrate/migrations")
		if _, err := os.Stat(migrationDir); err != nil {
			log.Fatalf("Migrations directory not found: %v", err)
		}
	}
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationDir),
		dbConnInfo.DSN(),
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Info("🦋 Migrations applied successfully")

	repo := repositories.NewRepositories(db)

	// Do a template sync of our pre-defined stuff
	if err := repo.Template().UpsertPredefinedTemplates(ctx); err != nil {
		log.Errorf("Failed to upsert predefined templates: %v", err)
	}

	// Create kubernetes client
	kubeClient := k8s.NewKubeClient(cfg, repo)

	// Create github client
	githubClient := github.NewGithubClient(cfg.GithubURL, cfg)

	// Buildkit settings manager
	buildkitSettings := buildkitd.NewBuildkitSettingsManager(cfg, repo, kubeClient)

	// Loki log querier
	lokiQuerier, err := loki.NewLokiLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to create Loki log querier, invalid config: %v", err)
	}

	// Prometheus client
	promClient, err := prometheus.NewPrometheusClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Prometheus client: %v", err)
	}

	// Database provider
	dbProvider := databases.NewDatabaseProvider()

	// Bootstrap
	if !cfg.SkipBootstrap {
		bootstrapper := &Bootstrapper{
			cfg:                     cfg,
			kubeClient:              kubeClient,
			repos:                   repo,
			buildkitSettingsManager: buildkitSettings,
		}
		if err := bootstrapper.Sync(ctx); err != nil {
			log.Errorf("Failed to sync system settings: %v", err)
		}
	}

	// Create webhook service
	variableService := variables_service.NewVariablesService(repo, kubeClient)
	webhooksService := webhooks_service.NewWebhooksService(repo)

	// Create deployment controller
	deploymentController := deployctl.NewDeploymentController(ctx, cancel, cfg, kubeClient, valkeyClient, repo, githubClient, webhooksService, variableService)

	// Create services
	teamService := team_service.NewTeamService(repo, kubeClient)
	projectService := project_service.NewProjectService(cfg, repo, kubeClient, webhooksService, deploymentController)
	serviceService := service_service.NewServiceService(cfg, repo, githubClient, kubeClient, deploymentController, dbProvider, webhooksService, variableService)
	environmentService := environment_service.NewEnvironmentService(repo, kubeClient, deploymentController)
	logService := logs_service.NewLogsService(repo, kubeClient, lokiQuerier)
	deploymentService := deployments_service.NewDeploymentService(repo, kubeClient, deploymentController, githubClient, lokiQuerier, registry.NewImageChecker(cfg, repo, kubeClient), variableService)
	systemService := system_service.NewSystemService(cfg, repo, buildkitSettings)
	metricsService := metric_service.NewMetricService(promClient, repo)
	instanceService := instance_service.NewInstanceService(cfg, repo, kubeClient)
	storageService := storage_service.NewStorageService(cfg, repo, kubeClient)
	templateService := templates_service.NewTemplatesService(repo, kubeClient)

	stringCache := cache.NewStringCache(valkeyClient, "unbind")

	// Create OAuth2 config

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.DexClientID,
		ClientSecret: cfg.DexClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.DexIssuerUrlExternal + "/auth",
			TokenURL: cfg.DexIssuerURL + "/token",
		},
		RedirectURL: fmt.Sprintf("%s/auth/callback", cfg.ExternalAPIURL),
		Scopes:      []string{"openid", "profile", "email", "offline_access", "groups"},
	}

	// Implementation
	srvImpl := &server.Server{
		KubeClient:           kubeClient,
		Cfg:                  cfg,
		Repository:           repo,
		GithubClient:         githubClient,
		StringCache:          stringCache,
		HttpClient:           &http.Client{},
		DeploymentController: deploymentController,
		DatabaseProvider:     dbProvider,
		DNSChecker:           utils.NewDNSChecker(),
		UpdateManager:        updater.New(cfg, Version, kubeClient),
		TeamService:          teamService,
		ProjectService:       projectService,
		ServiceService:       serviceService,
		EnvironmentService:   environmentService,
		LogService:           logService,
		DeploymentService:    deploymentService,
		SystemService:        systemService,
		MetricsService:       metricsService,
		WebhooksService:      webhooksService,
		InstanceService:      instanceService,
		VariablesService:     variableService,
		StorageService:       storageService,
		TemplateService:      templateService,
		OauthConfig:          oauthConfig,
	}

	// New chi router
	r := chi.NewRouter()

	allowedOrigins := []string{
		cfg.ExternalUIUrl,
	}

	if cfg.InjectDevOrigins {
		allowedOrigins = append(allowedOrigins, []string{
			"http://localhost:3000",
			"*.unbind.app",
		}...)
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"version": "` + Version + `"}`))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)

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
		mw := middleware.NewMiddleware(cfg, repo, api)

		// /setup group
		setupGroup := huma.NewGroup(api, "/setup")
		setupGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
			op.Tags = []string{"Setup"}
			next(op)
		})
		setup_handler.RegisterHandlers(srvImpl, setupGroup)

		// /auth group
		authGroup := huma.NewGroup(api, "/auth")
		authGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
			op.Tags = []string{"Auth"}
			next(op)
		})
		auth_handler.RegisterHandlers(srvImpl, authGroup)

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

		// /metrics group
		metricsGroup := huma.NewGroup(api, "/metrics")
		metricsGroup.UseMiddleware(mw.Authenticate)
		metricsGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
			op.Tags = []string{"Metrics"}
			next(op)
		})
		metrics_handler.RegisterHandlers(srvImpl, metricsGroup)

		// /unbindwebhooks group
		unbindwebhooksGroup := huma.NewGroup(api, "/unbindwebhooks")
		unbindwebhooksGroup.UseMiddleware(mw.Authenticate)
		unbindwebhooksGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
			op.Tags = []string{"Unbind Webhooks"}
			next(op)
		})
		unbindwebhooks_handler.RegisterHandlers(srvImpl, unbindwebhooksGroup)

		// /instances group
		instancesGroup := huma.NewGroup(api, "/instances")
		instancesGroup.UseMiddleware(mw.Authenticate)
		instancesGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
			op.Tags = []string{"Instances"}
			next(op)
		})
		instances_handler.RegisterHandlers(srvImpl, instancesGroup)

		// /storage group
		storageGroup := huma.NewGroup(api, "/storage")
		storageGroup.UseMiddleware(mw.Authenticate)
		storageGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
			op.Tags = []string{"Storage"}
			next(op)
		})
		storage_handler.RegisterHandlers(srvImpl, storageGroup)

		// /templates group
		templatesGroup := huma.NewGroup(api, "/templates")
		templatesGroup.UseMiddleware(mw.Authenticate)
		templatesGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
			op.Tags = []string{"Templates"}
			next(op)
		})
		template_handler.RegisterHandlers(srvImpl, templatesGroup)
	})

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

	// Start cron jobs
	// Initialize scheduler
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		log.Fatal("Failed to create scheduler", "err", err)
	}

	// ! TODO - we should leverage valkey or something to prevent concurrent runs
	// Clean up test DNS ingresses
	_, err = scheduler.NewJob(
		gocron.DurationJob(10*time.Minute),
		gocron.NewTask(
			func(ctx context.Context) {
				log.Infof("Cleaning up ingress tests.")
				if err := kubeClient.DeleteOldVerificationIngresses(ctx, kubeClient.GetInternalClient()); err != nil {
					log.Error("Failed to delete old verification ingresses", "err", err)
				}
			},
			ctx,
		),
	)
	if err != nil {
		log.Fatal("Failed to create ingress cleanup job", "err", err)
	}

	// Keep database variables in sync
	_, err = scheduler.NewJob(
		gocron.DurationJob(10*time.Minute),
		gocron.NewTask(
			func(ctx context.Context) {
				log.Infof("Syncing database secrets.")
				if err := kubeClient.SyncDatabaseSecrets(ctx); err != nil {
					log.Error("Failed to sync database secrets", "err", err)
				}
			},
			ctx,
		),
	)

	// Start the scheduler
	scheduler.Start()
	defer func() {
		if err := scheduler.Shutdown(); err != nil {
			log.Error("Scheduler shutdown error", "err", err)
		} else {
			log.Info("Scheduler gracefully stopped")
		}
	}()

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
	log.Infof("Starting Unbind API version %s", Version)
	// Load environment variables from .env file
	err := godotenv.Overload()
	if err != nil {
		log.Warn("Error loading .env file:", err)
	}

	cfg := config.NewConfig()
	cfg.BuildImage = BuildImage
	startAPI(cfg)
}
