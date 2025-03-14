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
	"github.com/unbindapp/unbind-api/internal/github"
	"github.com/unbindapp/unbind-api/internal/k8s"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/middleware"
	"github.com/unbindapp/unbind-api/internal/repository/repositories"
	"github.com/unbindapp/unbind-api/internal/server"
	github_handler "github.com/unbindapp/unbind-api/internal/server/github"
	projects_handler "github.com/unbindapp/unbind-api/internal/server/projects"
	teams_handler "github.com/unbindapp/unbind-api/internal/server/teams"
	user_handler "github.com/unbindapp/unbind-api/internal/server/user"
	webhook_handler "github.com/unbindapp/unbind-api/internal/server/webhook"
	project_service "github.com/unbindapp/unbind-api/internal/services/project"
	team_service "github.com/unbindapp/unbind-api/internal/services/team"
	"github.com/valkey-io/valkey-go"
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
	repo := repositories.NewRepositories(db)

	// Create kubernetes client
	kubeClient := k8s.NewKubeClient(cfg)

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
			Scopes:      []string{"openid", "profile", "email", "offline_access", "groups"},
		},
		GithubClient:   github.NewGithubClient(cfg),
		StringCache:    database.NewStringCache(client, "unbind"),
		HttpClient:     &http.Client{},
		TeamService:    team_service.NewTeamService(repo, kubeClient),
		ProjectService: project_service.NewProjectService(repo, kubeClient),
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
			URL: cfg.ExternalURL,
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
		srvImpl.HealthCheck,
	)

	// /auth group
	authGroup := huma.NewGroup(api, "/auth")
	authGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Auth"}
		next(op)
	})
	huma.Register(
		authGroup,
		huma.Operation{
			OperationID: "login",
			Summary:     "Login",
			Description: "Login",
			Path:        "/login",
			Method:      http.MethodGet,
		},
		srvImpl.Login)
	huma.Register(
		authGroup,
		huma.Operation{
			OperationID: "callback",
			Summary:     "Callback",
			Description: "Callback",
			Path:        "/callback",
			Method:      http.MethodGet,
		},
		srvImpl.Callback,
	)

	// /user group
	userGroup := huma.NewGroup(api, "/user")
	userGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"User"}
		next(op)
	})
	userHandlers := user_handler.NewHandlerGroup(srvImpl)
	userGroup.UseMiddleware(mw.Authenticate)
	huma.Register(
		userGroup,
		huma.Operation{
			OperationID: "me",
			Summary:     "Get User (Me)",
			Description: "Get the current user",
			Path:        "/me",
			Method:      http.MethodGet,
		},
		userHandlers.Me,
	)

	ghGroup := huma.NewGroup(api, "/github")
	ghGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"GitHub"}
		next(op)
	})
	githubHandlers := github_handler.NewHandlerGroup(srvImpl)
	ghGroup.UseMiddleware(mw.Authenticate)
	huma.Register(
		ghGroup,
		huma.Operation{
			OperationID: "app-create",
			Summary:     "Create App",
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
			Summary:     "List Apps",
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
			Summary:     "List Admin Organizations",
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
	webhookGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Webhook"}
		next(op)
	})
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
			Summary:     "Save GitHub App",
			Description: "Save GitHub app via code exchange and redirect to installation",
			Path:        "/github/app/save",
			Method:      http.MethodGet,
		},
		webhookHandlers.HandleGithubAppSave,
	)

	// /teams
	teamsGroup := huma.NewGroup(api, "/teams")
	teamsGroup.UseMiddleware(mw.Authenticate)
	teamsGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Teams"}
		next(op)
	})
	teamHandlers := teams_handler.NewHandlerGroup(srvImpl)
	huma.Register(
		teamsGroup,
		huma.Operation{
			OperationID: "list-teams",
			Summary:     "List Teams",
			Description: "List all teams the current user is a member of",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		teamHandlers.ListTeams,
	)
	huma.Register(
		teamsGroup,
		huma.Operation{
			OperationID: "update-team",
			Summary:     "Update Team",
			Description: "Update a team",
			Path:        "/update",
			Method:      http.MethodPut,
		},
		teamHandlers.UpdateTeam,
	)

	// /projects group
	projectsGroup := huma.NewGroup(api, "/projects")
	projectsGroup.UseMiddleware(mw.Authenticate)
	projectsGroup.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {
		op.Tags = []string{"Projects"}
		next(op)
	})
	projectHandlers := projects_handler.NewHandlerGroup(srvImpl)
	huma.Register(
		projectsGroup,
		huma.Operation{
			OperationID: "create-project",
			Summary:     "Create Project",
			Description: "Create a project",
			Path:        "/create",
			Method:      http.MethodPost,
		},
		projectHandlers.CreateProject,
	)
	huma.Register(
		projectsGroup,
		huma.Operation{
			OperationID: "list-projects",
			Summary:     "List Projects",
			Description: "List all projects",
			Path:        "/list",
			Method:      http.MethodGet,
		},
		projectHandlers.ListProjects,
	)
	huma.Register(
		projectsGroup,
		huma.Operation{
			OperationID: "update-project",
			Summary:     "Update Project",
			Description: "Update a project",
			Path:        "/update",
			Method:      http.MethodPut,
		},
		projectHandlers.UpdateProject,
	)

	// Start the server
	addr := ":8089"
	log.Infof("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
