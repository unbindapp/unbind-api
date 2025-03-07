package oauth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/middleware"
	"github.com/valkey-io/valkey-go"
)

// StartOauth2Server starts the OAuth2 server API
func StartOauth2Server(cfg *config.Config) {
	// Initialize valkey (redis)
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{cfg.ValkeyURL}})
	if err != nil {
		log.Fatal("Failed to create valkey client", "err", err)
	}
	defer client.Close()

	oauth2Srv := setupOAuthServer(cfg, client)

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

	// Start server
	addr := ":8090"
	log.Infof("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
