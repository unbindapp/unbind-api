package variables_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

func (self *VariablesService) ResolveSingleReference(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, serviceID, referenceID uuid.UUID) (string, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   serviceID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return "", err
	}

	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return "", err
	}

	// Get the reference
	reference, err := self.repo.Variables().GetReferenceByID(ctx, referenceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Variable reference not found")
		}
		return "", err
	}
	if reference.TargetServiceID != serviceID {
		return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Variable reference not found")
	}

	// Get the namespace
	namespace, err := self.repo.Service().GetDeploymentNamespace(ctx, serviceID)
	if err != nil {
		return "", err
	}
	return self.resolveReference(ctx, client, namespace, reference)
}

// Resolve variable references into map[string]string
func (self *VariablesService) ResolveAllReferences(ctx context.Context, serviceID uuid.UUID) (map[string]string, error) {
	// Get all references
	result := make(map[string]string)
	references, err := self.repo.Variables().GetReferencesForService(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	if len(references) == 0 {
		return result, nil
	}

	// Get team namespace
	namespace, err := self.repo.Service().GetDeploymentNamespace(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	// Get our meta client
	client := self.k8s.GetInternalClient()

	for _, reference := range references {
		// Resolve the reference
		value, err := self.resolveReference(ctx, client, namespace, reference)
		if err != nil {
			return nil, err
		}

		// Add to our result map
		result[reference.TargetName] = value
	}

	return result, nil
}

func (self *VariablesService) resolveReference(ctx context.Context, client *kubernetes.Clientset, namespace string, reference *ent.VariableReference) (string, error) {
	sourceValues := make(map[string]string)
	for _, source := range reference.Sources {
		// The key we want to replace in our template
		sourceKey := fmt.Sprintf("${%s.%s}", source.SourceKubernetesName, source.Key)

		switch source.Type {
		case schema.VariableReferenceTypeVariable:
			// Get from kubernetes secret
			secret, err := self.k8s.GetSecret(ctx, source.SourceKubernetesName, namespace, client)
			if err != nil {
				if !errors.IsNotFound(err) {
					return "", err
				} else {
					if _, err = self.repo.Variables().AttachError(ctx, reference.ID, err); err != nil {
						log.Errorf("Failed to attach error to variable reference %s: %v", reference.ID, err)
						return "", err
					}
					return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, fmt.Sprintf("Unable to resolve variable %s ${%s.%s}", reference.TargetName, source.SourceKubernetesName, source.Key))
				}
			}

			// Get the value from the secret
			value, ok := secret.Data[source.Key]
			if !ok {
				if _, err = self.repo.Variables().AttachError(ctx, reference.ID, err); err != nil {
					log.Errorf("Failed to attach error to variable reference %s: %v", reference.ID, err)
					return "", err
				}
				return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, fmt.Sprintf("Unable to resolve variable %s ${%s.%s}", reference.TargetName, source.SourceKubernetesName, source.Key))
			}
			sourceValues[sourceKey] = string(value)
		case schema.VariableReferenceTypeInternalEndpoint, schema.VariableReferenceTypeExternalEndpoint:
			endpoints, err := self.k8s.DiscoverEndpointsByLabels(ctx, namespace,
				map[string]string{
					source.SourceType.KubernetesLabel(): source.SourceID.String(),
				}, client)
			if err != nil {
				return "", err
			}

			if len(endpoints.Internal) == 0 && len(endpoints.External) == 0 {
				if _, err = self.repo.Variables().AttachError(ctx, reference.ID, err); err != nil {
					log.Errorf("Failed to attach error to variable reference %s: %v", reference.ID, err)
					return "", err
				}
				return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, fmt.Sprintf("Unable to resolve variable %s ${%s.%s}", reference.TargetName, source.SourceKubernetesName, source.Key))
			}

			found := false
			if source.Type == schema.VariableReferenceTypeInternalEndpoint {
				for _, endpoint := range endpoints.Internal {
					if endpoint.KubernetesName == source.Key {
						// Figure out port
						var targetPort *schema.PortSpec
						for _, port := range endpoint.Ports {
							if port.Protocol != nil && *port.Protocol == schema.ProtocolTCP {
								targetPort = &port
								break
							}
						}
						if targetPort == nil {
							log.Warnf("No TCP port found for endpoint %s", endpoint.KubernetesName)
							if _, err = self.repo.Variables().AttachError(ctx, reference.ID, err); err != nil {
								log.Errorf("Failed to attach error to variable reference %s: %v", reference.ID, err)
								return "", err
							}
							return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, fmt.Sprintf("Unable to resolve variable %s ${%s.%s}", reference.TargetName, source.SourceKubernetesName, source.Key))
						}
						sourceValues[sourceKey] = endpoint.DNS
						found = true
					}
				}

				continue
			}

			// External endpoint
			for _, endpoint := range endpoints.External {
				for _, host := range endpoint.Hosts {
					if host.Host == source.Key {
						sourceValues[sourceKey] = host.Host
						found = true
						break
					}
				}
			}

			if !found {
				if _, err = self.repo.Variables().AttachError(ctx, reference.ID, err); err != nil {
					log.Errorf("Failed to attach error to variable reference %s: %v", reference.ID, err)
					return "", err
				}
				return "", errdefs.NewCustomError(errdefs.ErrTypeNotFound, fmt.Sprintf("Unable to resolve ${%s.%s}", source.SourceKubernetesName, source.Key))
			}
		}
	}

	// Replace all references in the template
	template := reference.ValueTemplate
	for k, v := range sourceValues {
		template = strings.ReplaceAll(template, k, v)
	}
	return template, nil
}
