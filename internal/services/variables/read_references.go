package variables_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (self *VariablesService) GetAvailableVariableReferences(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, teamID uuid.UUID) (*models.AvailableVariableReferenceResponse, error) {
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
	team, err := self.repo.Team().GetByID(ctx, teamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, err
	}

	// Get all kubernetes secrets in the team
	teamSecret := team.KubernetesSecret
	projectSecrets := make(map[uuid.UUID]string)
	environmentSecrets := make(map[uuid.UUID]string)
	serviceSecrets := make(map[uuid.UUID]string)
	for _, project := range team.Edges.Projects {
		projectSecrets[project.ID] = project.KubernetesSecret
		for _, environment := range project.Edges.Environments {
			environmentSecrets[environment.ID] = environment.KubernetesSecret
			for _, service := range environment.Edges.Services {
				serviceSecrets[service.ID] = service.KubernetesSecret
			}
		}
	}

	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Get all kubernetes secrets in the team
	k8sSecrets, err := self.k8s.GetAllSecrets(
		ctx,
		team.ID,
		teamSecret,
		projectSecrets,
		environmentSecrets,
		serviceSecrets,
		client,
		team.Namespace,
	)

	if err != nil {
		return nil, err
	}

	// Get all the DNS/endpoints
	endpoints, err := self.k8s.DiscoverEndpointsByLabels(
		ctx,
		team.Namespace,
		map[string]string{
			"unbind-team": team.ID.String(),
		},
		client,
	)

	if err != nil {
		return nil, err
	}

	return models.TransformAvailableVariableResponse(k8sSecrets, endpoints), nil
}

// Resolve a variable reference value for a key
func (self *VariablesService) ResolveAvailableReferenceValue(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.ResolveVariableReferenceInput) (string, error) {
	permissionChecks := []permissions_repo.PermissionCheck{}
	targetLabel := ""
	switch input.SourceType {
	case schema.VariableReferenceSourceTypeTeam:
		targetLabel = "unbind-team"
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeTeam,
			ResourceID:   input.SourceID,
		})
	case schema.VariableReferenceSourceTypeProject:
		targetLabel = "unbind-project"
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeProject,
			ResourceID:   input.SourceID,
		})
	case schema.VariableReferenceSourceTypeEnvironment:
		targetLabel = "unbind-environment"
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeEnvironment,
			ResourceID:   input.SourceID,
		})
	case schema.VariableReferenceSourceTypeService:
		targetLabel = "unbind-service"
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
				targetLabel: input.SourceID.String(),
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
					return fmt.Sprint("%s:%d", endpoint.DNS, targetPort.Port), nil
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
