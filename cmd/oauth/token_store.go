package oauth

import (
	"context"
	"fmt"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/unbindapp/unbind-api/internal/database/repository"
)

// Custom token store that using our database backend and handles authorization code flow
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
