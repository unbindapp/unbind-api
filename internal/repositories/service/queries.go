package service_repo

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/deployment"
	"github.com/unbindapp/unbind-api/ent/githubapp"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
	"github.com/unbindapp/unbind-api/ent/team"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
)

func (self *ServiceRepository) GetByID(ctx context.Context, serviceID uuid.UUID) (svc *ent.Service, err error) {
	svc, err = self.base.DB.Service.Query().
		Where(service.IDEQ(serviceID)).
		WithEnvironment(func(eq *ent.EnvironmentQuery) {
			eq.WithProject(func(pq *ent.ProjectQuery) {
				pq.WithTeam()
			})
		}).
		WithServiceConfig().
		WithDeployments(func(dq *ent.DeploymentQuery) {
			dq.Order(ent.Desc(deployment.FieldCreatedAt))
			dq.Limit(1)
		}).
		WithCurrentDeployment().
		Only(ctx)

	if err != nil {
		return nil, err
	}
	// Get the latest successful deployment for the service, if needed
	if len(svc.Edges.Deployments) == 0 || svc.Edges.Deployments[0].Status != schema.DeploymentStatusSucceeded {
		lastSuccessfulDeployment, err := self.deploymentRepo.GetLastSuccessfulDeployment(ctx, svc.ID)
		if err != nil && !ent.IsNotFound(err) {
			return nil, err
		}
		if !ent.IsNotFound(err) {
			svc.Edges.Deployments = append(svc.Edges.Deployments, lastSuccessfulDeployment)
		}
	}
	return svc, nil
}

func (self *ServiceRepository) GetByName(ctx context.Context, name string) (*ent.Service, error) {
	return self.base.DB.Service.Query().
		Where(service.NameEQ(name)).
		WithEnvironment(func(eq *ent.EnvironmentQuery) {
			eq.WithProject(func(pq *ent.ProjectQuery) {
				pq.WithTeam()
			})
		}).
		Only(ctx)
}

func (self *ServiceRepository) GetByInstallationIDAndRepoName(ctx context.Context, installationID int64, repoName string) ([]*ent.Service, error) {
	return self.base.DB.Service.Query().
		Where(service.GithubInstallationIDEQ(installationID)).
		Where(service.GitRepositoryEQ(repoName)).
		WithServiceConfig().
		WithCurrentDeployment().
		Order(ent.Desc(service.FieldCreatedAt)).
		All(ctx)
}

func (self *ServiceRepository) GetByEnvironmentID(ctx context.Context, environmentID uuid.UUID) ([]*ent.Service, error) {
	services, err := self.base.DB.Service.Query().
		Where(service.EnvironmentIDEQ(environmentID)).
		WithServiceConfig().
		WithCurrentDeployment().
		Order(ent.Desc(service.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return nil, err
	}

	// Get the latest deployment for each service
	for _, svc := range services {
		latestDeployment, err := self.base.DB.Deployment.Query().
			Where(deployment.ServiceIDEQ(svc.ID)).
			Order(ent.Desc(deployment.FieldCreatedAt)).
			Limit(1).
			First(ctx)

		if err != nil && !ent.IsNotFound(err) {
			return nil, err
		}

		if latestDeployment != nil {
			svc.Edges.Deployments = []*ent.Deployment{latestDeployment}
			if latestDeployment.Status != schema.DeploymentStatusSucceeded {
				// Get last successful deployment if the latest one is not successful
				lastSuccessfulDeployment, err := self.deploymentRepo.GetLastSuccessfulDeployment(ctx, svc.ID)
				if err != nil && !ent.IsNotFound(err) {
					return nil, err
				}
				if !ent.IsNotFound(err) {
					svc.Edges.Deployments = append(svc.Edges.Deployments, lastSuccessfulDeployment)
				}
			}
		}
	}

	return services, nil
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
func (self *ServiceRepository) SummarizeServices(ctx context.Context, environmentIDs []uuid.UUID) (counts map[uuid.UUID]int, providers map[uuid.UUID][]string, err error) {
	counts = make(map[uuid.UUID]int)

	// Maps to not duplicate providers
	providerSets := make(map[uuid.UUID]map[string]struct{})

	// Initialize sets for each environment ID
	for _, envID := range environmentIDs {
		providerSets[envID] = make(map[string]struct{})
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

		if svc.Edges.ServiceConfig.Database != nil {
			providerSets[svc.EnvironmentID][string(*svc.Edges.ServiceConfig.Database)] = struct{}{}
			continue
		}

		if svc.Edges.ServiceConfig.RailpackFramework != nil {
			providerSets[svc.EnvironmentID][string(*svc.Edges.ServiceConfig.RailpackFramework)] = struct{}{}
			continue
		}

		if svc.Edges.ServiceConfig.RailpackProvider != nil {
			providerSets[svc.EnvironmentID][string(*svc.Edges.ServiceConfig.RailpackProvider)] = struct{}{}
			continue
		}

		providerSets[svc.EnvironmentID][string(svc.Edges.ServiceConfig.Type)] = struct{}{}
	}

	// Convert to slices
	providers = make(map[uuid.UUID][]string)

	for envID, providerSet := range providerSets {
		for provider := range providerSet {
			providers[envID] = append(providers[envID], provider)
		}
	}

	return
}

// See if a service needs a new deployment
type NeedsDeploymentResponse string

const (
	NeedsBuildAndDeployment NeedsDeploymentResponse = "needs_build_and_deployment"
	NeedsDeployment         NeedsDeploymentResponse = "needs_deployment"
	NoDeploymentNeeded      NeedsDeploymentResponse = "no_deployment_needed"
)

func (self *ServiceRepository) NeedsDeployment(ctx context.Context, service *ent.Service) (NeedsDeploymentResponse, error) {
	if service.Edges.CurrentDeployment == nil || service.Edges.CurrentDeployment.ResourceDefinition == nil {
		return NoDeploymentNeeded, nil
	}
	// See if a deployment is needed
	// Create a an object with only fields we care to compare
	existingCrd := &v1.Service{
		Spec: v1.ServiceSpec{
			Builder: service.Edges.CurrentDeployment.ResourceDefinition.Spec.Builder,
			Config: v1.ServiceConfigSpec{
				GitBranch:  service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.GitBranch,
				Hosts:      service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.Hosts,
				Replicas:   service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.Replicas,
				Ports:      service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.Ports,
				RunCommand: service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.RunCommand,
				Public:     service.Edges.CurrentDeployment.ResourceDefinition.Spec.Config.Public,
			},
		},
	}
	// Create a new CRD to compare it
	var gitBranch string
	if service.Edges.ServiceConfig.GitBranch != nil {
		gitBranch = *service.Edges.ServiceConfig.GitBranch
		if !strings.HasPrefix(gitBranch, "refs/heads/") {
			gitBranch = fmt.Sprintf("refs/heads/%s", gitBranch)
		}
	}
	newCrd := &v1.Service{
		Spec: v1.ServiceSpec{
			Builder: string(service.Edges.ServiceConfig.Builder),
			Config: v1.ServiceConfigSpec{
				GitBranch:  gitBranch,
				Hosts:      service.Edges.ServiceConfig.Hosts,
				Replicas:   utils.ToPtr(service.Edges.ServiceConfig.Replicas),
				Ports:      service.Edges.ServiceConfig.Ports,
				RunCommand: service.Edges.ServiceConfig.RunCommand,
				Public:     service.Edges.ServiceConfig.Public,
			},
		},
	}

	// Changing builder requires a new build
	if existingCrd.Spec.Builder != newCrd.Spec.Builder {
		return NeedsBuildAndDeployment, nil
	}

	// Branch needs a new build
	if existingCrd.Spec.Config.GitBranch != newCrd.Spec.Config.GitBranch {
		return NeedsBuildAndDeployment, nil
	}

	// Just update the custom resource
	if !reflect.DeepEqual(existingCrd, newCrd) {
		return NeedsDeployment, nil
	}

	return NoDeploymentNeeded, nil
}
