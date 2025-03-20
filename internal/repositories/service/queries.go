package service_repo

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubapp"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
	"github.com/unbindapp/unbind-api/ent/team"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

func (self *ServiceRepository) GetByID(ctx context.Context, serviceID uuid.UUID) (*ent.Service, error) {
	return self.base.DB.Service.Query().
		Where(service.IDEQ(serviceID)).
		WithEnvironment().
		WithServiceConfig().
		Only(ctx)
}

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

func (self *ServiceRepository) CountDomainCollisons(ctx context.Context, tx repository.TxInterface, domain string) (int, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}
	return db.ServiceConfig.Query().
		Where(
			serviceconfig.HostEqualFold(strings.ToLower(domain)),
		).Count(ctx)
}

func (self *ServiceRepository) GetDeploymentNamespace(ctx context.Context, serviceID uuid.UUID) (string, error) {
	svc, err := self.base.DB.Service.Query().
		Where(service.IDEQ(serviceID)).
		QueryEnvironment().
		QueryProject().
		QueryTeam().Select(team.FieldNamespace).
		Only(ctx)
	if err != nil {
		return "", err
	}
	return svc.Namespace, nil
}
