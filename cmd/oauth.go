package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/golang-jwt/jwt"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/database"
	"github.com/unbindapp/unbind-api/internal/database/repository"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/middleware"
	"github.com/unbindapp/unbind-api/internal/oauth2server"
	"github.com/unbindapp/unbind-api/internal/utils"
	"github.com/valkey-io/valkey-go"
)

const ACCESS_TOKEN_EXP = 1 * time.Minute
const REFRESH_TOKEN_EXP = 24 * time.Hour * 14 // 2 weeks
var ALLOWED_SCOPES = []string{"openid", "profile", "email", "offline_access"}

// Custom token store that allows multiple refresh tokens per user
type customTokenStore struct {
	clientStore *dbClientStore
	repository  *repository.Repository
}

func NewCustomTokenStore(clientStore *dbClientStore, repository *repository.Repository) *customTokenStore {
	return &customTokenStore{
		clientStore: clientStore,
		repository:  repository,
	}
}

// Handle authorization codes and tokens
func (s *customTokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	code := info.GetCode()
	accessToken := info.GetAccess()
	refreshToken := info.GetRefresh()

	// 1) Detect if we are storing an authorization code
	if code != "" && accessToken == "" && refreshToken == "" {
		// We’re in the code stage
		return s.createAuthorizationCode(ctx, info)
	}

	// 2) Otherwise, we’re storing real access/refresh tokens
	return s.createTokens(ctx, info)
}

// createAuthorizationCode is a helper for storing a code
func (s *customTokenStore) createAuthorizationCode(ctx context.Context, info oauth2.TokenInfo) error {
	code := info.GetCode()
	clientID := info.GetClientID()
	userID := info.GetUserID()
	scope := info.GetScope()

	if code == "" {
		return fmt.Errorf("cannot store empty authorization code")
	}
	if clientID == "" {
		return fmt.Errorf("cannot store code for empty client ID")
	}

	// Typically short expiry, e.g. 5-10 minutes
	expiresIn := info.GetCodeExpiresIn()
	if expiresIn <= 0 {
		expiresIn = 5 * time.Minute
	}
	expiresAt := time.Now().Add(expiresIn)

	u, err := s.repository.GetUserByEmail(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	_, err = s.repository.CreateAuthCode(ctx, code, clientID, scope, u, expiresAt)
	return err
}

// createTokens is a helper for storing real access/refresh tokens
func (s *customTokenStore) createTokens(ctx context.Context, info oauth2.TokenInfo) error {
	accessToken := info.GetAccess()
	refreshToken := info.GetRefresh()
	clientID := info.GetClientID()
	userData := info.GetUserID()
	scope := info.GetScope()

	if accessToken == "" {
		return fmt.Errorf("cannot create token with empty access token")
	}
	if refreshToken == "" {
		return fmt.Errorf("cannot create token with empty refresh token")
	}
	if clientID == "" {
		return fmt.Errorf("cannot create token with empty client ID")
	}

	// Find the user by userData (often an email or ID)
	u, err := s.repository.GetUserByEmail(ctx, userData)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Refresh token expiry from info or fallback
	refreshExpiresIn := info.GetRefreshExpiresIn()
	if refreshExpiresIn <= 0 {
		refreshExpiresIn = 24 * time.Hour // or your default
	}
	refreshExpiresAt := time.Now().Add(refreshExpiresIn)

	_, err = s.repository.CreateToken(
		ctx,
		accessToken,
		refreshToken,
		clientID,
		scope,
		refreshExpiresAt,
		u,
	)
	return err
}

// RemoveByCode is called after a successful exchange or if the code is no longer valid.
func (s *customTokenStore) RemoveByCode(ctx context.Context, code string) error {
	// Delete it from DB so it can’t be reused
	return s.repository.DeleteAuthCode(ctx, code)
}

func (s *customTokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return s.repository.RevokeAccessToken(ctx, access)
}

func (s *customTokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return s.repository.RevokeRefreshToken(ctx, refresh)
}

// GetByCode is called when the library needs to exchange an authorization code for tokens.
func (s *customTokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	// Query your DB to see if this code exists
	authCode, err := s.repository.GetAuthCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if authCode == nil {
		return nil, errors.ErrInvalidAuthorizeCode
	}
	if authCode.ExpiresAt.Before(time.Now()) {
		return nil, errors.ErrInvalidAuthorizeCode
	}

	// Build a TokenInfo object with the code fields.
	// The library only needs enough to confirm the code’s existence
	// and associated user/client. The rest of the fields can be minimal.
	token := &models.Token{
		ClientID:      authCode.ClientID,
		UserID:        authCode.Edges.User.Email,
		Code:          authCode.AuthCode,
		CodeCreateAt:  authCode.CreatedAt,
		CodeExpiresIn: authCode.ExpiresAt.Sub(time.Now()),
		Scope:         authCode.Scope,
	}

	return token, nil
}
func (s *customTokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	token, err := s.repository.GetByAccessToken(ctx, access)

	if err != nil {
		return nil, err
	}

	// Get the user associated with the token
	u, err := token.QueryUser().Only(ctx)
	if err != nil {
		return nil, err
	}

	// Fixed clientStore method call - it doesn't take context
	clientInfo, err := s.clientStore.GetByID(ctx, token.ClientID)
	if err != nil {
		return nil, err
	}

	// Create a proper token implementation
	tokenInfo := &models.Token{
		ClientID:        token.ClientID,
		UserID:          u.Email,
		RedirectURI:     clientInfo.GetDomain(), // Using domain instead of missing scope
		Scope:           token.Scope,            // Using default scopes
		Access:          token.AccessToken,
		Refresh:         token.RefreshToken,
		AccessCreateAt:  token.CreatedAt,
		RefreshCreateAt: token.CreatedAt,
	}

	// Calculate remaining access token duration
	accessExpiresDuration := token.ExpiresAt.Sub(time.Now())
	if accessExpiresDuration < 0 {
		accessExpiresDuration = 0
	}

	refreshExpiryTime := token.CreatedAt.Add(REFRESH_TOKEN_EXP)
	refreshExpiresDuration := refreshExpiryTime.Sub(time.Now())
	if refreshExpiresDuration < 0 {
		refreshExpiresDuration = 0
	}

	// Set different expiration durations for access and refresh tokens
	tokenInfo.SetAccessExpiresIn(accessExpiresDuration)
	tokenInfo.SetRefreshExpiresIn(refreshExpiresDuration)

	return tokenInfo, nil
}

func (s *customTokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	token, err := s.repository.GetByRefreshToken(ctx, refresh)

	if err != nil {
		return nil, err
	}

	// Get the user associated with the token
	u, err := token.QueryUser().Only(ctx)
	if err != nil {
		return nil, err
	}

	// Get client information
	clientInfo, err := s.clientStore.GetByID(ctx, token.ClientID)
	if err != nil {
		return nil, err
	}

	// Create a proper token implementation
	tokenInfo := &models.Token{
		ClientID:        token.ClientID,
		UserID:          u.Email,
		RedirectURI:     clientInfo.GetDomain(),
		Scope:           token.Scope,
		Access:          token.AccessToken,
		Refresh:         token.RefreshToken,
		AccessCreateAt:  token.CreatedAt,
		RefreshCreateAt: token.CreatedAt,
	}

	// Calculate remaining access token duration
	accessExpiresDuration := token.ExpiresAt.Sub(time.Now())
	if accessExpiresDuration < 0 {
		accessExpiresDuration = 0
	}

	refreshExpiryTime := token.CreatedAt.Add(REFRESH_TOKEN_EXP)
	refreshExpiresDuration := refreshExpiryTime.Sub(time.Now())
	if refreshExpiresDuration < 0 {
		refreshExpiresDuration = 0
	}

	// Set different expiration durations for access and refresh tokens
	tokenInfo.SetAccessExpiresIn(accessExpiresDuration)
	tokenInfo.SetRefreshExpiresIn(refreshExpiresDuration)

	return tokenInfo, nil
}

func generateIDToken(ti oauth2.TokenInfo, repo *repository.Repository, issuer string, privateKey *rsa.PrivateKey, kid string) (string, error) {
	now := time.Now()

	// Gather the data we need
	userID := ti.GetUserID()
	clientID := ti.GetClientID()

	u, err := repo.GetUserByEmail(context.Background(), userID)
	if err != nil {
		return "", fmt.Errorf("failed to find user: %w", err)
	}

	claims := jwt.MapClaims{
		"iss":            issuer,
		"sub":            userID,
		"aud":            clientID,
		"iat":            now.Unix(),
		"exp":            now.Add(ACCESS_TOKEN_EXP).Unix(),
		"email":          u.Email,
		"email_verified": true,
		"name":           "John Doe",
	}

	// Create a token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Optionally set KID in header so clients know which key was used:
	if kid != "" {
		token.Header["kid"] = kid
	}

	// Sign with the RSA private key
	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

// Persisent client store
type CacheClientInto struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Domain string `json:"domain"`
	UserID string `json:"user_id"`
	Public bool   `json:"public"`
}

func (c CacheClientInto) GetID() string {
	return c.ID
}

func (c CacheClientInto) GetSecret() string {
	return c.Secret
}

func (c CacheClientInto) GetDomain() string {
	return c.Domain
}

func (c CacheClientInto) IsPublic() bool {
	return c.Public
}

func (c CacheClientInto) GetUserID() string {
	return c.UserID
}

type dbClientStore struct {
	cache *database.ValkeyCache[CacheClientInto]
}

func NewDBClientStore(cache *database.ValkeyCache[CacheClientInto]) *dbClientStore {
	return &dbClientStore{
		cache: cache,
	}
}

func (s *dbClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	cacheItem, err := s.cache.Get(ctx, id)
	if err != nil {
		if err == valkey.Nil {
			return nil, errors.ErrInvalidClient
		}
		return nil, err
	}
	return cacheItem, nil
}

func (s *dbClientStore) Set(id string, client oauth2.ClientInfo) error {
	cacheItem := CacheClientInto{
		ID:     client.GetID(),
		Secret: client.GetSecret(),
		Domain: client.GetDomain(),
		UserID: client.GetUserID(),
		Public: client.IsPublic(),
	}
	return s.cache.SetWithExpiration(context.TODO(), id, cacheItem, 30*time.Minute)
}

func setupOAuthServer(cfg *config.Config, valkey valkey.Client) *oauth2server.Oauth2Server {
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
	if err := db.Schema.Create(context.TODO()); err != nil {
		log.Fatal("Failed to run migrations", "err", err)
	}
	repo := repository.NewRepository(db)

	// Load private key
	pkey, pkeyBytes, err := repo.GetOrGenerateJWTPrivateKey(context.Background())
	if err != nil {
		log.Fatalf("Failed to get private key: %v", err)
	}

	// Use our custom token store
	clientStoreCache := database.NewCache[CacheClientInto](valkey, "auth")
	clientStore := NewDBClientStore(clientStoreCache)
	tokenStore := NewCustomTokenStore(clientStore, repo)
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)
	keyID := "unbind-oauth-key"
	jwtGen := generates.NewJWTAccessGenerate(keyID, pkeyBytes, jwt.SigningMethodRS256)
	manager.MapAccessGenerate(jwtGen)

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
		Repository: repo,
		Cfg:        cfg,
		PrivateKey: pkey,
		Kid:        keyID,
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
		idToken, err := generateIDToken(ti, repo, cfg.ExternalOauth2URL, pkey, keyID)
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

func startOauth2Server(cfg *config.Config) {
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
	r.Use(middleware.Logger)
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
