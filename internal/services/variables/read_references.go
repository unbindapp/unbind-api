package variables_service

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *VariablesService) GetAvailableVariableReferences(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, teamID, projectID, environmentID, serviceID uuid.UUID) ([]models.AvailableVariableReference, error) {
	// ! TODO - we're going to need to change all of our permission checks to filter not reject
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeTeam,
			ResourceID:   teamID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Get available variable references
	team, project, environment, _, err := self.validateInputs(ctx, teamID, projectID, environmentID, serviceID)
	if err != nil {
		return nil, err
	}

	// Build a map of names
	nameMap := make(map[uuid.UUID]string)
	nameMap[team.ID] = team.Name
	nameMap[project.ID] = project.Name

	// Base secret names
	teamSecret := team.KubernetesSecret
	projectSecret := project.KubernetesSecret
	environmentSecret := environment.KubernetesSecret

	// They can access all service secrets in the same project
	serviceSecrets := make(map[uuid.UUID]string)
	// Get all environments in this project
	projectEnvironments, err := self.repo.Environment().GetForProject(ctx, nil, project.ID)
	if err != nil {
		return nil, err
	}
	for _, environment := range projectEnvironments {
		nameMap[environment.ID] = environment.Name
		for _, service := range environment.Edges.Services {
			if service.ID == serviceID {
				continue
			}
			serviceSecrets[service.ID] = service.KubernetesSecret
			nameMap[service.ID] = service.Name
		}
	}

	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Use WaitGroup to handle concurrent K8s operations
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Variables to store results
	var k8sSecrets []models.SecretData
	var endpoints *models.EndpointDiscovery
	var secretsErr, endpointsErr error

	// Add two tasks to WaitGroup
	wg.Add(2)

	// Goroutine for getting all kubernetes secrets
	go func() {
		defer wg.Done()
		secrets, err := self.k8s.GetAllSecrets(
			ctx,
			team.ID,
			teamSecret,
			project.ID,
			projectSecret,
			environment.ID,
			environmentSecret,
			serviceSecrets,
			client,
			team.Namespace,
		)

		mu.Lock()
		k8sSecrets = secrets
		secretsErr = err
		mu.Unlock()
	}()

	// Goroutine for getting all DNS/endpoints
	go func() {
		defer wg.Done()
		eps, err := self.k8s.DiscoverEndpointsByLabels(
			ctx,
			team.Namespace,
			map[string]string{
				"unbind-project": project.ID.String(),
			},
			client,
		)

		mu.Lock()
		filteredEPs := &models.EndpointDiscovery{
			External: []models.IngressEndpoint{},
			Internal: []models.ServiceEndpoint{},
		}
		for _, ep := range eps.Internal {
			if ep.ServiceID != serviceID {
				filteredEPs.Internal = append(filteredEPs.Internal, ep)
			}
		}
		for _, ep := range eps.External {
			if ep.ServiceID != serviceID {
				filteredEPs.External = append(filteredEPs.External, ep)
			}
		}
		endpoints = filteredEPs
		endpointsErr = err
		mu.Unlock()
	}()

	// Wait for both operations to complete
	wg.Wait()

	// Check for errors
	if secretsErr != nil {
		return nil, secretsErr
	}

	if endpointsErr != nil {
		return nil, endpointsErr
	}

	return models.TransformAvailableVariableResponse(k8sSecrets, endpoints, nameMap), nil
}
