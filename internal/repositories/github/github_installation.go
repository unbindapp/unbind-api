package github_repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubapp"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/ent/schema"
)

func (self *GithubRepository) GetInstallationByID(ctx context.Context, ID int64) (*ent.GithubInstallation, error) {
	return self.base.DB.GithubInstallation.Query().Where(githubinstallation.ID(ID)).WithGithubApp().Only(ctx)
}

func (self *GithubRepository) GetInstallationsByCreator(ctx context.Context, createdBy uuid.UUID) ([]*ent.GithubInstallation, error) {
	installations, err := self.base.DB.GithubInstallation.Query().
		Where(
			githubinstallation.Active(true),
			githubinstallation.Suspended(false),
		).WithGithubApp(func(gaq *ent.GithubAppQuery) {
		gaq.Where(
			githubapp.CreatedByEQ(createdBy),
		)
	}).All(ctx)

	if err != nil {
		return nil, err
	}

	// ! TODO - we should probably make sure that AccountType and AccountLogin can't be duplicated
	// Use a map to find duplicates based on AccountType and AccountLogin
	seen := make(map[string]int) // map key -> index in result slice
	result := make([]*ent.GithubInstallation, 0, len(installations))

	for _, installation := range installations {
		// Create a unique key based on AccountType and AccountLogin
		key := fmt.Sprintf("%s:%s", installation.AccountType.String(), installation.AccountLogin)

		if existingIdx, exists := seen[key]; exists {
			// Keep newer one
			if installation.CreatedAt.After(result[existingIdx].CreatedAt) {
				result[existingIdx] = installation
			}
			continue
		}

		seen[key] = len(result)
		result = append(result, installation)
	}

	return result, nil
}

func (self *GithubRepository) GetInstallationsByAppID(ctx context.Context, appID int64) ([]*ent.GithubInstallation, error) {
	return self.base.DB.GithubInstallation.Query().Where(githubinstallation.GithubAppID(appID)).All(ctx)
}

func (self *GithubRepository) UpsertInstallation(
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
	permissions schema.GithubInstallationPermissions,
	events []string,
) (*ent.GithubInstallation, error) {
	err := self.base.DB.GithubInstallation.Create().
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

	return self.base.DB.GithubInstallation.Query().Where(githubinstallation.ID(id)).Only(ctx)
}

func (self *GithubRepository) SetInstallationActive(ctx context.Context, id int64, active bool) (*ent.GithubInstallation, error) {
	return self.base.DB.GithubInstallation.UpdateOneID(id).
		SetActive(active).
		Save(ctx)
}

func (self *GithubRepository) SetInstallationSuspended(ctx context.Context, id int64, suspended bool) (*ent.GithubInstallation, error) {
	return self.base.DB.GithubInstallation.UpdateOneID(id).
		SetSuspended(suspended).
		Save(ctx)
}
