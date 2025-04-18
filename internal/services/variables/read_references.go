package variables_service

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (self *VariablesService) GetAvailableVariableReferences(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, teamID, projectID, environmentID, serviceID uuid.UUID) (*models.AvailableVariableReferenceResponse, error) {
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
		for _, service := range environment.Edges.Services {
			if service.ID == serviceID {
				continue
			}
			serviceSecrets[service.ID] = service.KubernetesSecret
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
		endpoints = eps
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

	return models.TransformAvailableVariableResponse(k8sSecrets, endpoints), nil
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

	switch input.Type {
	case schema.VariableReferenceTypeVariable:
		// Get variable
		secret, err := self.k8s.GetSecret(ctx, input.Name, namespace, client)
		if err != nil {
			if errors.IsNotFound(err) {
				return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Variable not found")
			}
			return "", err
		}
		for k, v := range secret.Data {
			if k == input.Key {
				return string(v), nil
			}
		}
		return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Variable not found")
	case schema.VariableReferenceTypeInternalEndpoint, schema.VariableReferenceTypeExternalEndpoint:
		// Get endpoint
		endpoints, err := self.k8s.DiscoverEndpointsByLabels(
			ctx,
			namespace,
			map[string]string{
				input.SourceType.KubernetesLabel(): input.SourceID.String(),
			},
			client,
		)
		if err != nil {
			return "", err
		}

		if len(endpoints.Internal) == 0 && len(endpoints.External) == 0 {
			return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Endpoint not found")
		}

		if input.Type == schema.VariableReferenceTypeInternalEndpoint {
			for _, endpoint := range endpoints.Internal {
				if endpoint.Name == input.Key {
					// Figure out port
					var targetPort *schema.PortSpec
					for _, port := range endpoint.Ports {
						if port.Protocol != nil && *port.Protocol == schema.ProtocolTCP {
							targetPort = &port
							break
						}
					}
					if targetPort == nil {
						return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, "No TCP port found for endpoint")
					}
					return fmt.Sprintf("%s:%d", endpoint.DNS, targetPort.Port), nil
				}
			}
		} else {
			for _, endpoint := range endpoints.External {
				for _, host := range endpoint.Hosts {
					if host.Host == input.Key {
						return host.Host, nil
					}
				}
			}
		}

	}

	return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Variable not found")
}
