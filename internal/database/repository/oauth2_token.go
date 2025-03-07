package repository

import (
	"context"
	"time"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/oauth2token"
)

func (r *Repository) CreateToken(ctx context.Context, accessToken, refreshToken, clientID, scope string, expiresAt time.Time, user *ent.User) (*ent.Oauth2Token, error) {
	return r.DB.Oauth2Token.Create().SetAccessToken(accessToken).SetRefreshToken(refreshToken).SetClientID(clientID).SetExpiresAt(expiresAt).SetScope(scope).SetUser(user).Save(ctx)
}

func (r *Repository) RevokeAccessToken(ctx context.Context, accessToken string) error {
	// Mark the token as revoked instead of removing it
	dbToken, err := r.DB.Oauth2Token.Query().Where(
		oauth2token.AccessToken(accessToken),
	).Only(ctx)

	if err != nil {
		return err
	}

	_, err = dbToken.Update().SetRevoked(true).Save(ctx)
	return err
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	// Mark the token as revoked instead of removing it
	dbToken, err := r.DB.Oauth2Token.Query().Where(
		oauth2token.RefreshToken(refreshToken),
	).Only(ctx)

	if err != nil {
		return err
	}

	_, err = dbToken.Update().SetRevoked(true).Save(ctx)
	return err
}

func (r *Repository) GetByAccessToken(ctx context.Context, accessToken string) (*ent.Oauth2Token, error) {
	return r.DB.Oauth2Token.Query().Where(
		oauth2token.AccessToken(accessToken),
		oauth2token.Revoked(false),
	).Only(ctx)
}

func (r *Repository) GetByRefreshToken(ctx context.Context, refreshToken string) (*ent.Oauth2Token, error) {
	return r.DB.Oauth2Token.Query().Where(
		oauth2token.RefreshToken(refreshToken),
		oauth2token.Revoked(false),
	).Only(ctx)
}
