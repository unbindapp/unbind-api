package repository

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/oauth2code"
	"github.com/unbindapp/unbind-api/ent/oauth2token"
)

func (r *Repository) CreateToken(ctx context.Context, accessToken, refreshToken, clientID, scope string, expiresAt time.Time, user *ent.User) (*ent.Oauth2Token, error) {
	return r.DB.Oauth2Token.Create().SetAccessToken(accessToken).SetRefreshToken(refreshToken).SetClientID(clientID).SetExpiresAt(expiresAt).SetScope(scope).SetUser(user).Save(ctx)
}

func (r *Repository) RevokeAccessToken(ctx context.Context, accessToken string) error {
	_, err := r.DB.Oauth2Token.Update().Where(oauth2token.AccessToken(accessToken)).SetRevoked(true).Save(ctx)
	return err
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	_, err := r.DB.Oauth2Token.Update().Where(oauth2token.RefreshToken(refreshToken)).SetRevoked(true).Save(ctx)
	return err
}

func (r *Repository) GetByAccessToken(ctx context.Context, accessToken string) (*ent.Oauth2Token, error) {
	return r.DB.Oauth2Token.Query().Where(
		oauth2token.AccessToken(accessToken),
	).Only(ctx)
}

func (r *Repository) GetByRefreshToken(ctx context.Context, refreshToken string) (*ent.Oauth2Token, error) {
	return r.DB.Oauth2Token.Query().Where(
		oauth2token.RefreshToken(refreshToken),
	).Only(ctx)
}

func (r *Repository) CleanTokenStore(ctx context.Context) (result error) {
	_, err := r.DB.Oauth2Token.Delete().Where(
		oauth2token.Or(
			oauth2token.Revoked(true),
			oauth2token.ExpiresAtLT(time.Now()),
		),
	).Exec(ctx)
	if err != nil {
		multierror.Append(result, err)
	}

	_, err = r.DB.Oauth2Code.Delete().Where(
		oauth2code.ExpiresAtLT(time.Now()),
	).Exec(ctx)

	if err != nil {
		multierror.Append(result, err)
	}

	return result
}
