package main

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/database"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/middleware"
	"github.com/unbindapp/unbind-api/internal/oauth2server"
	"github.com/unbindapp/unbind-api/internal/repository/repositories"
	"github.com/unbindapp/unbind-api/internal/utils"
	"github.com/valkey-io/valkey-go"
)

const ACCESS_TOKEN_EXP = 2 * time.Minute
const REFRESH_TOKEN_EXP = 24 * time.Hour * 14 // 2 weeks
var ALLOWED_SCOPES = []string{"openid", "profile", "email", "offline_access", "groups"}

// Setup the go-oauth2 server
func setupOAuthServer(ctx context.Context, cfg *config.Config, valkey valkey.Client) *oauth2server.Oauth2Server {
	manager := manage.NewDefaultManager()

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

	// Load private key
	pkey, _, err := repo.Oauth().GetOrGenerateJWTPrivateKey(ctx)
	if err != nil {
		log.Fatalf("Failed to get private key: %v", err)
	}

	// Use our custom token store
	clientStoreCache := database.NewCache[CacheClientInto](valkey, "auth")
	clientStore := NewDBClientStore(ctx, clientStoreCache)
	tokenStore := NewCustomTokenStore(clientStore, repo)
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)
	keyID := "unbind-oauth-key"

	// Create the access token generator
	accessTokenGen := &accessTokenGenerator{
		keyID:      keyID,
		PrivateKey: pkey,
	}
	manager.MapAccessGenerate(accessTokenGen)

	// Register the client for Dex
	dexCallbackUrl, _ := utils.JoinURLPaths(cfg.DexIssuerUrlExternal, "/callback")
	clientStore.Set("dex-client", &models.Client{
		ID:     "dex-client",
		Secret: "dex-secret",
		Domain: dexCallbackUrl,
	})

	// Configure the OAuth2 server
	manager.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    ACCESS_TOKEN_EXP,
		RefreshTokenExp:   REFRESH_TOKEN_EXP,
		IsGenerateRefresh: true,
	})

	// Configure the password grant
	manager.SetPasswordTokenCfg(&manage.Config{
		AccessTokenExp:    ACCESS_TOKEN_EXP,
		RefreshTokenExp:   REFRESH_TOKEN_EXP,
		IsGenerateRefresh: true,
	})

	// Create Oauth2 server
	oauth2Server := &oauth2server.Oauth2Server{
		Ctx:         ctx,
		Repository:  repo,
		Cfg:         cfg,
		PrivateKey:  pkey,
		Kid:         keyID,
		StringCache: database.NewStringCache(valkey, "unbind-oauth-str"),
	}

	// Create the server
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(func(r *http.Request) (string, string, error) {
		// Try HTTP Basic Auth first
		clientID, clientSecret, ok := r.BasicAuth()
		if ok && clientID != "" {
			return clientID, clientSecret, nil
		}

		// Then try form values
		if err := r.ParseForm(); err != nil {
			return "", "", err
		}

		clientID = r.Form.Get("client_id")
		clientSecret = r.Form.Get("client_secret")

		if clientID == "" {
			return "", "", errors.New("client id not found")
		}

		return clientID, clientSecret, nil
	})

	srv.SetClientScopeHandler(func(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error) {
		// Check if the client is allowed to use the requested scope
		parsedScopes := strings.Split(tgr.Scope, " ")
		for _, scope := range parsedScopes {
			if !slices.Contains(ALLOWED_SCOPES, scope) {
				log.Warnf("Client %s requested invalid scope: %s", tgr.ClientID, scope)
				return false, nil
			}
		}
		return true, nil
	})

	srv.SetExtensionFieldsHandler(func(ti oauth2.TokenInfo) map[string]interface{} {
		// Only produce an ID token if the request included the "openid" scope,
		hasOpenid := false
		for _, scope := range strings.Split(ti.GetScope(), " ") {
			if scope == "openid" {
				hasOpenid = true
				break
			}
		}
		if !hasOpenid {
			return nil // do nothing extra
		}

		// 2) Generate the ID token
		idToken, err := generateIDToken(ctx, ti, repo, cfg.ExternalOauth2URL, pkey, keyID)
		if err != nil {
			log.Errorf("Error generating ID token: %v\n", err)
			return nil
		}

		// Return an extra field "id_token"
		return map[string]interface{}{
			"id_token": idToken,
		}
	})

	// Set the password grant handler
	srv.SetPasswordAuthorizationHandler(oauth2Server.PasswordAuthorizationHandler)

	// Set error handler
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Infof("Internal Error: %v", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Infof("Response Error: %v", re.Error.Error())
	})

	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		// Get user ID from query parameters
		userID = r.URL.Query().Get("user_id")
		if userID == "" {
			log.Warn("No user_id found in authorization request")
			return "", errors.New("user not authenticated")
		}

		return userID, nil
	})

	oauth2Server.Srv = srv

	return oauth2Server
}

// StartOauth2Server starts the OAuth2 server API
func StartOauth2Server(cfg *config.Config) {
	// Create base context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize valkey (redis)
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{cfg.ValkeyURL}})
	if err != nil {
		log.Fatal("Failed to create valkey client", "err", err)
	}
	defer client.Close()

	oauth2Srv := setupOAuthServer(ctx, cfg, client)

	// Setup router
	// New chi router
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	// OAuth2/OIDC endpoints
	r.Post("/token", oauth2Srv.HandleToken)
	r.Get("/authorize", oauth2Srv.HandleAuthorize)
	r.Get("/userinfo", oauth2Srv.HandleUserinfo)
	r.Get("/login", oauth2Srv.HandleLoginPage)
	r.Post("/login", oauth2Srv.HandleLoginSubmit)
	r.Get("/.well-known/openid-configuration", oauth2Srv.HandleOpenIDConfiguration)
	r.Get("/.well-known/jwks.json", oauth2Srv.HandleJWKS)

	// Cron job to clean up expired tokens
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Hour().Do(func() {
		err := oauth2Srv.Repository.Oauth().CleanTokenStore(ctx)
		if err != nil {
			log.Error("Failed to clean token store", "err", err)
		}
	})
	s.StartAsync()

	// Start server
	addr := ":8090"
	log.Infof("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Overload()
	if err != nil {
		log.Warn("Error loading .env file:", err)
	}

	cfg := config.NewConfig()
	StartOauth2Server(cfg)
}
