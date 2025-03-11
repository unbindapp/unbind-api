package repository

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/internal/models"
)

func (r *Repository) GetGithubInstallationByID(ctx context.Context, ID int64) (*ent.GithubInstallation, error) {
	return r.DB.GithubInstallation.Query().Where(githubinstallation.ID(ID)).Only(ctx)
}

func (r *Repository) GetGithubInstallationsByAppID(ctx context.Context, appID int64) ([]*ent.GithubInstallation, error) {
	return r.DB.GithubInstallation.Query().Where(githubinstallation.GithubAppID(appID)).All(ctx)
}

func (r *Repository) UpsertGithubInstallation(
	ctx context.Context,
	id int64,
	appID int64,
	accountID int64,
	accountLogin string,
	accountType githubinstallation.AccountType,
	accountURL string,
	repositorySelection githubinstallation.RepositorySelection,
	suspended bool,
	active bool,
	permissions models.GithubInstallationPermissions,
	events []string,
) (*ent.GithubInstallation, error) {
	err := r.DB.GithubInstallation.Create().
		SetID(id).
		SetGithubAppID(appID).
		SetAccountID(accountID).
		SetAccountLogin(accountLogin).
		SetAccountType(accountType).
		SetAccountURL(accountURL).
		SetRepositorySelection(repositorySelection).
		SetSuspended(suspended).
		SetActive(active).
		SetPermissions(permissions).
		SetEvents(events).
		OnConflictColumns(
			githubinstallation.FieldID,
		).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.DB.GithubInstallation.Query().Where(githubinstallation.ID(id)).Only(ctx)
}

func (r *Repository) SetInstallationActive(ctx context.Context, id int64, active bool) (*ent.GithubInstallation, error) {
	return r.DB.GithubInstallation.UpdateOneID(id).
		SetActive(active).
		Save(ctx)
}

func (r *Repository) SetInstallationSuspended(ctx context.Context, id int64, suspended bool) (*ent.GithubInstallation, error) {
	return r.DB.GithubInstallation.UpdateOneID(id).
		SetSuspended(suspended).
		Save(ctx)
}
