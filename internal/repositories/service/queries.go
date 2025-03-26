package service_repo

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/deployment"
	"github.com/unbindapp/unbind-api/ent/githubapp"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
	"github.com/unbindapp/unbind-api/ent/team"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
)

func (self *ServiceRepository) GetByID(ctx context.Context, serviceID uuid.UUID) (*ent.Service, error) {
	return self.base.DB.Service.Query().
		Where(service.IDEQ(serviceID)).
		WithEnvironment().
		WithServiceConfig().
		WithDeployments(func(dq *ent.DeploymentQuery) {
			dq.Order(ent.Desc(deployment.FieldCreatedAt))
			dq.Limit(1)
		}).
		WithCurrentDeployment().
		Only(ctx)
}

func (self *ServiceRepository) GetByInstallationIDAndRepoName(ctx context.Context, installationID int64, repoName string) ([]*ent.Service, error) {
	return self.base.DB.Service.Query().
		Where(service.GithubInstallationIDEQ(installationID)).
		Where(service.GitRepositoryEQ(repoName)).
		WithServiceConfig().
		WithDeployments(func(dq *ent.DeploymentQuery) {
			dq.Order(ent.Desc(deployment.FieldCreatedAt))
			dq.Limit(1)
		}).
		WithCurrentDeployment().
		All(ctx)
}

func (self *ServiceRepository) GetByEnvironmentID(ctx context.Context, environmentID uuid.UUID) ([]*ent.Service, error) {
	return self.base.DB.Service.Query().
		Where(service.EnvironmentIDEQ(environmentID)).
		WithServiceConfig().
		WithDeployments(func(dq *ent.DeploymentQuery) {
			dq.Order(ent.Desc(deployment.FieldCreatedAt))
			dq.Limit(1)
		}).
		WithCurrentDeployment().
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

	// Convert the domain to lowercase for case-insensitive comparison
	lowerDomain := strings.ToLower(domain)

	// Get all services with non-empty hosts JSON field
	configs, err := db.ServiceConfig.Query().
		Where(
			serviceconfig.HostsNotNil(),
		).
		All(ctx)

	if err != nil {
		// Return just the legacy count if we can't query the hosts field
		return 0, nil
	}

	// Count matches in the hosts field manually
	jsonCount := 0
	for _, config := range configs {
		// Skip empty hosts field
		if len(config.Hosts) == 0 {
			continue
		}

		// Check each host in the array
		for _, host := range config.Hosts {
			if strings.EqualFold(host.Host, lowerDomain) {
				jsonCount++
				break
			}
		}
	}

	return jsonCount, nil
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

// Summarize services in environment
func (self *ServiceRepository) SummarizeServices(ctx context.Context, environmentIDs []uuid.UUID) (counts map[uuid.UUID]int, providers map[uuid.UUID][]enum.Provider, frameworks map[uuid.UUID][]enum.Framework, err error) {
	counts = make(map[uuid.UUID]int)

	// Maps to not duplicate providers and frameworks
	providerSets := make(map[uuid.UUID]map[enum.Provider]struct{})
	frameworkSets := make(map[uuid.UUID]map[enum.Framework]struct{})

	// Initialize sets for each environment ID
	for _, envID := range environmentIDs {
		providerSets[envID] = make(map[enum.Provider]struct{})
		frameworkSets[envID] = make(map[enum.Framework]struct{})
	}

	services, err := self.base.DB.Service.Query().
		Select(service.FieldEnvironmentID).
		Where(service.EnvironmentIDIn(environmentIDs...)).
		WithServiceConfig().
		All(ctx)
	if err != nil {
		return
	}

	for _, svc := range services {
		counts[svc.EnvironmentID]++

		if svc.Edges.ServiceConfig == nil {
			continue
		}

		if svc.Edges.ServiceConfig.Provider != nil {
			providerSets[svc.EnvironmentID][*svc.Edges.ServiceConfig.Provider] = struct{}{}
		}

		if svc.Edges.ServiceConfig.Framework != nil {
			frameworkSets[svc.EnvironmentID][*svc.Edges.ServiceConfig.Framework] = struct{}{}
		}
	}

	// Convert to slices
	providers = make(map[uuid.UUID][]enum.Provider)
	frameworks = make(map[uuid.UUID][]enum.Framework)

	for envID, providerSet := range providerSets {
		for provider := range providerSet {
			providers[envID] = append(providers[envID], provider)
		}
	}

	for envID, frameworkSet := range frameworkSets {
		for framework := range frameworkSet {
			frameworks[envID] = append(frameworks[envID], framework)
		}
	}

	return
}
