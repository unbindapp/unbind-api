package repository

import (
	"context"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/ent"
)

// GetGithubApp returns the GithubApp entity., ent.NotFoundError if not found.s
func (r *Repository) GetGithubApp(ctx context.Context) (*ent.GithubApp, error) {
	return r.DB.GithubApp.Query().First(ctx)
}

func (r *Repository) CreateGithubApp(ctx context.Context, app *github.AppConfig) (*ent.GithubApp, error) {
	return r.DB.GithubApp.Create().
		SetGithubAppID(app.GetID()).
		SetClientID(app.GetClientID()).
		SetClientSecret(app.GetClientSecret()).
		SetWebhookSecret(app.GetWebhookSecret()).
		SetPrivateKey(app.GetPEM()).
		Save(ctx)
}
