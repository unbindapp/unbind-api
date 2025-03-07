package repository

import (
	"context"
	"time"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/oauth2code"
)

func (r *Repository) CreateAuthCode(ctx context.Context, code, clientID, scope string, user *ent.User, expiresAt time.Time) (*ent.Oauth2Code, error) {
	return r.DB.Oauth2Code.Create().SetAuthCode(code).SetClientID(clientID).SetScope(scope).SetUser(user).SetExpiresAt(expiresAt).Save(ctx)
}

func (r *Repository) DeleteAuthCode(ctx context.Context, code string) error {
	_, err := r.DB.Oauth2Code.Delete().Where(oauth2code.AuthCode(code)).Exec(ctx)
	return err
}

func (r *Repository) GetAuthCode(ctx context.Context, code string) (*ent.Oauth2Code, error) {
	return r.DB.Oauth2Code.Query().Where(oauth2code.AuthCode(code)).WithUser().Only(ctx)
}
