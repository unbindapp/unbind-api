package variables_service

import (
	"context"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
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

	// Build a map of names and icons
	kubernetesNameMap := make(map[uuid.UUID]string)
	nameMap := make(map[uuid.UUID]string)
	iconMap := make(map[uuid.UUID]string)
	kubernetesNameMap[team.ID] = team.KubernetesName
	nameMap[team.ID] = team.Name
	iconMap[team.ID] = "team"
	kubernetesNameMap[project.ID] = project.KubernetesName
	nameMap[project.ID] = project.Name
	iconMap[project.ID] = "project"

	// Base secret names
	teamSecret := team.KubernetesSecret
	projectSecret := project.KubernetesSecret
	environmentSecret := environment.KubernetesSecret

	// They can access all service secrets in the same project
	serviceSecrets := make(map[uuid.UUID]string)
	var serviceIDs []uuid.UUID
	kubernetesNameMap[environment.ID] = environment.KubernetesName
	nameMap[environment.ID] = environment.Name
	iconMap[environment.ID] = "environment"

	// Re-fetch environment to populate edges
	environment, err = self.repo.Environment().GetByID(ctx, environment.ID)
	if err != nil {
		return nil, err
	}

	// Re-fetch services to get icons
	environmentServices, err := self.repo.Service().GetByEnvironmentID(ctx, environment.ID, false)
	if err != nil {
		return nil, err
	}

	for _, service := range environmentServices {
		if service.ID == serviceID {
			continue
		}
		serviceSecrets[service.ID] = service.KubernetesSecret
		kubernetesNameMap[service.ID] = service.KubernetesName
		nameMap[service.ID] = service.Name
		iconMap[service.ID] = service.Edges.ServiceConfig.Icon
		serviceIDs = append(serviceIDs, service.ID)
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
				"unbind-environment": environment.ID.String(),
			},
			true,
			client,
		)

		if err != nil {
			log.Errorf("Failed to discover endpoints: %v", err)
			return
		}

		mu.Lock()
		filteredEPs := &models.EndpointDiscovery{
			External: []models.IngressEndpoint{},
			Internal: []models.ServiceEndpoint{},
		}
		for _, ep := range eps.Internal {
			if ep.ServiceID != serviceID && slices.Contains(serviceIDs, ep.ServiceID) {
				filteredEPs.Internal = append(filteredEPs.Internal, ep)
			}
		}
		for _, ep := range eps.External {
			if ep.ServiceID != serviceID && slices.Contains(serviceIDs, ep.ServiceID) {
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

	return models.TransformAvailableVariableResponse(k8sSecrets, endpoints, kubernetesNameMap, nameMap, iconMap), nil
}

// Resolve a variable reference value for a key
func (self *VariablesService) ResolveAvailableReferenceValue(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.ResolveVariableReferenceInput) (string, error) {
	permissionChecks := []permissions_repo.PermissionCheck{}
	switch input.SourceType {
	case schema.VariableReferenceSourceTypeTeam:
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeTeam,
			ResourceID:   input.SourceID,
		})
	case schema.VariableReferenceSourceTypeProject:
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeProject,
			ResourceID:   input.SourceID,
		})
	case schema.VariableReferenceSourceTypeEnvironment:
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeEnvironment,
			ResourceID:   input.SourceID,
		})
	case schema.VariableReferenceSourceTypeService:
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   input.SourceID,
		})
	default:
		return "", errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Invalid source type")
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return "", err
	}

	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return "", err
	}

	// Get team by ID
	namespace, err := self.repo.Team().GetNamespace(ctx, input.TeamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return "", err
	}

	// Create a reference source to reuse the resolution logic
	source := schema.VariableReferenceSource{
		Type:                 input.Type,
		SourceType:           input.SourceType,
		SourceID:             input.SourceID,
		SourceKubernetesName: input.Name,
		Key:                  input.Key,
	}

	value, err := self.resolveSourceValue(ctx, client, namespace, source)
	if err != nil {
		return "", err
	}

	return value, nil
}
