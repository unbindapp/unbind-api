package service_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubapp"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/ent/service"
)

func (self *ServiceRepository) GetByInstallationIDAndRepoName(ctx context.Context, installationID int64, repoName string) ([]*ent.Service, error) {
	return self.base.DB.Service.Query().
		Where(service.GithubInstallationIDEQ(installationID)).
		Where(service.GitRepositoryEQ(repoName)).
		WithServiceConfig().
		All(ctx)
}

func (self *ServiceRepository) GetGithubPrivateKey(ctx context.Context, serviceID uuid.UUID) (string, error) {
	svc, err := self.base.DB.Service.Query().
		Where(service.IDEQ(serviceID)).
		Only(ctx)
	if err != nil {
		return "", err
	}

	if svc.GithubInstallationID == nil {
		return "", nil
	}

	app, err := self.base.DB.GithubInstallation.Query().
		Where(githubinstallation.IDEQ(*svc.GithubInstallationID)).
		QueryGithubApp().
		Select(githubapp.FieldPrivateKey).
		Only(ctx)
	if err != nil {
		return "", err
	}

	return app.PrivateKey, nil
}
