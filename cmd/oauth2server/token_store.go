package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

type customTokenStore struct {
	clientStore *clientStore
	repository  repositories.RepositoriesInterface
}

func NewCustomTokenStore(clientStore *clientStore, repository repositories.RepositoriesInterface) *customTokenStore {
	return &customTokenStore{
		clientStore: clientStore,
		repository:  repository,
	}
}

func (s *customTokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	code := info.GetCode()
	accessToken := info.GetAccess()
	refreshToken := info.GetRefresh()

	if code != "" && accessToken == "" && refreshToken == "" {
		return s.createAuthorizationCode(ctx, info)
	}

	return s.createTokens(ctx, info)
}

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

	expiresIn := info.GetCodeExpiresIn()
	if expiresIn <= 0 {
		expiresIn = 5 * time.Minute
	}
	expiresAt := time.Now().Add(expiresIn)

	u, err := s.repository.User().GetByEmail(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	_, err = s.repository.Oauth().CreateAuthCode(ctx, code, clientID, scope, u, expiresAt)
	return err
}

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

	u, err := s.repository.User().GetByEmail(ctx, userData)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	refreshExpiresIn := info.GetRefreshExpiresIn()
	if refreshExpiresIn <= 0 {
		refreshExpiresIn = 24 * time.Hour
	}
	refreshExpiresAt := time.Now().Add(refreshExpiresIn)

	_, err = s.repository.Oauth().CreateToken(
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

func (s *customTokenStore) RemoveByCode(ctx context.Context, code string) error {
	return s.repository.Oauth().DeleteAuthCode(ctx, code)
}

func (s *customTokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return s.repository.Oauth().RevokeAccessToken(ctx, access)
}

func (s *customTokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return s.repository.Oauth().RevokeRefreshToken(ctx, refresh)
}

func (s *customTokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	authCode, err := s.repository.Oauth().GetAuthCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if authCode == nil {
		return nil, errors.ErrInvalidAuthorizeCode
	}
	if authCode.ExpiresAt.Before(time.Now()) {
		return nil, errors.ErrInvalidAuthorizeCode
	}

	return &models.Token{
		ClientID:      authCode.ClientID,
		UserID:        authCode.Edges.User.Email,
		Code:          authCode.AuthCode,
		CodeCreateAt:  authCode.CreatedAt,
		CodeExpiresIn: time.Until(authCode.ExpiresAt),
		Scope:         authCode.Scope,
	}, nil
}

func (s *customTokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	token, err := s.repository.Oauth().GetByAccessToken(ctx, access)
	if err != nil {
		return nil, err
	}

	u, err := token.QueryUser().Only(ctx)
	if err != nil {
		return nil, err
	}

	clientInfo, err := s.clientStore.GetByID(ctx, token.ClientID)
	if err != nil {
		return nil, err
	}

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

	accessExpiresDuration := time.Until(token.ExpiresAt)
	if accessExpiresDuration < 0 {
		accessExpiresDuration = 0
	}

	refreshExpiryTime := token.CreatedAt.Add(REFRESH_TOKEN_EXP)
	refreshExpiresDuration := time.Until(refreshExpiryTime)
	if refreshExpiresDuration < 0 {
		refreshExpiresDuration = 0
	}

	tokenInfo.SetAccessExpiresIn(accessExpiresDuration)
	tokenInfo.SetRefreshExpiresIn(refreshExpiresDuration)

	return tokenInfo, nil
}

func (s *customTokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	token, err := s.repository.Oauth().GetByRefreshToken(ctx, refresh)
	if err != nil {
		return nil, err
	}

	u, err := token.QueryUser().Only(ctx)
	if err != nil {
		return nil, err
	}

	clientInfo, err := s.clientStore.GetByID(ctx, token.ClientID)
	if err != nil {
		return nil, err
	}

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

	accessExpiresDuration := time.Until(token.ExpiresAt)
	if accessExpiresDuration < 0 {
		accessExpiresDuration = 0
	}

	refreshExpiryTime := token.CreatedAt.Add(REFRESH_TOKEN_EXP)
	refreshExpiresDuration := time.Until(refreshExpiryTime)
	if refreshExpiresDuration < 0 {
		refreshExpiresDuration = 0
	}

	tokenInfo.SetAccessExpiresIn(accessExpiresDuration)
	tokenInfo.SetRefreshExpiresIn(refreshExpiresDuration)

	return tokenInfo, nil
}
