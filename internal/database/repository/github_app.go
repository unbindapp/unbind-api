package repository

import (
	"context"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubapp"
)

// GetGithubApp returns the GithubApp entity., ent.NotFoundError if not found.s
func (r *Repository) GetGithubApp(ctx context.Context) (*ent.GithubApp, error) {
	return r.DB.GithubApp.Query().First(ctx)
}

// Get all github apps returns a slice of GithubApp entities.
func (r *Repository) GetGithubApps(ctx context.Context, withInstallations bool) ([]*ent.GithubApp, error) {
	q := r.DB.GithubApp.Query()
	if withInstallations {
		q.WithInstallations()
	}
	return q.All(ctx)
}

func (r *Repository) CreateGithubApp(ctx context.Context, app *github.AppConfig) (*ent.GithubApp, error) {
	return r.DB.GithubApp.Create().
		SetID(app.GetID()).
		SetClientID(app.GetClientID()).
		SetClientSecret(app.GetClientSecret()).
		SetWebhookSecret(app.GetWebhookSecret()).
		SetPrivateKey(app.GetPEM()).
		SetName(app.GetName()).
		Save(ctx)
}

func (r *Repository) GetGithubAppByID(ctx context.Context, ID int64) (*ent.GithubApp, error) {
	return r.DB.GithubApp.Query().Where(githubapp.ID(ID)).Only(ctx)
}
