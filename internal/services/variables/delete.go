package variables_service

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/models"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
)

// Delete a secret by key
func (self *VariablesService) DeleteVariablesByKey(ctx context.Context, userID uuid.UUID, bearerToken string, input models.BaseVariablesJSONInput, keys []models.VariableDeleteInput, referenceIDs []uuid.UUID) (*models.VariableResponse, error) {
	if len(referenceIDs) > 0 && input.Type != schema.VariableReferenceSourceTypeService {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Reference IDs are only valid for service variables")
	}

	var permissionChecks []permissions_repo.PermissionCheck

	switch input.Type {
	case schema.VariableReferenceSourceTypeTeam:
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeTeam,
			ResourceID:   input.TeamID,
		})
	case schema.VariableReferenceSourceTypeProject:
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeProject,
			ResourceID:   input.ProjectID,
		})
	case schema.VariableReferenceSourceTypeEnvironment:
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeEnvironment,
			ResourceID:   input.EnvironmentID,
		})
	case schema.VariableReferenceSourceTypeService:
		permissionChecks = append(permissionChecks, permissions_repo.PermissionCheck{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   input.ServiceID,
		})
	default:
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Invalid variable type")
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Verify input
	team, project, environment, service, secretName, err := self.validateBaseInputs(ctx, input.Type, input.TeamID, input.ProjectID, input.EnvironmentID, input.ServiceID)
	if err != nil {
		return nil, err
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	secrets := make(map[string][]byte)
	keysToDelete := []string{}
	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		// Delete reference IDs
		if input.Type == schema.VariableReferenceSourceTypeService {
			deletedCount, err := self.repo.Variables().DeleteReferences(ctx, tx, input.ServiceID, referenceIDs)
			if err != nil {
				return err
			}
			if deletedCount != len(referenceIDs) {
				return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Invalid reference IDs")
			}
		}

		// Get secrets
		secrets, err = self.k8s.GetSecretMap(ctx, secretName, team.Namespace, client)
		if err != nil {
			return err
		}

		// For updating var mounts
		var variableMounts []*schema.VariableMount
		variableMountsNeedsUpdate := false
		if service != nil && service.Edges.ServiceConfig != nil {
			variableMounts = service.Edges.ServiceConfig.VariableMounts
		}

		// Remove from map
		for _, secretKey := range keys {
			if input.Type == schema.VariableReferenceSourceTypeService {
				if slices.Contains(service.Edges.ServiceConfig.ProtectedVariables, secretKey.Name) {
					continue
				}
			}

			keysToDelete = append(keysToDelete, secretKey.Name)

			// Delete variable mounts if they exist
			indexToDelete := -1
			for i, variableMount := range variableMounts {
				if variableMount.Name == secretKey.Name {
					indexToDelete = i
					break
				}
			}
			if indexToDelete != -1 {
				variableMountsNeedsUpdate = true
				variableMounts = append(variableMounts[:indexToDelete], variableMounts[indexToDelete+1:]...)
			}
			delete(secrets, secretKey.Name)
		}

		// Update variable mounts
		if variableMountsNeedsUpdate {
			if err := self.repo.Service().UpdateVariableMounts(ctx, tx, service.ID, variableMounts); err != nil {
				return err
			}
		}

		// Update secrets
		_, err = self.k8s.UpdateSecret(ctx, secretName, team.Namespace, secrets, client)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	variableResponse := &models.VariableResponse{
		Variables:          make([]*models.VariableResponseItem, len(secrets)),
		VariableReferences: []*models.VariableReferenceResponse{},
	}
	i := 0
	for k, v := range secrets {
		variableResponse.Variables[i] = &models.VariableResponseItem{
			Type:  input.Type,
			Name:  k,
			Value: string(v),
		}
		i++
	}
	models.SortVariableResponse(variableResponse.Variables)

	// Add references if this is for a service
	if input.Type == schema.VariableReferenceSourceTypeService {
		references, err := self.repo.Variables().GetReferencesForService(ctx, input.ServiceID)
		if err != nil {
			return nil, err
		}
		variableResponse.VariableReferences = models.TransformVariableReferenceResponseEntities(references)
	}

	var sourceID uuid.UUID
	switch input.Type {
	case schema.VariableReferenceSourceTypeTeam:
		sourceID = team.ID
	case schema.VariableReferenceSourceTypeProject:
		sourceID = project.ID
	case schema.VariableReferenceSourceTypeEnvironment:
		sourceID = environment.ID
	case schema.VariableReferenceSourceTypeService:
		sourceID = service.ID
	}

	// Get services referencing deleted variables
	services, err := self.repo.Variables().GetServicesReferencingID(ctx, sourceID, keysToDelete)

	// Build labels to restart
	var labelSelectors []string
	// Convert the labels map to a selector string
	// var labelSelectors []string
	// for key, value := range labels {
	// 	labelSelectors = append(labelSelectors, fmt.Sprintf("%s=%s", key, value))
	// }
	// labelSelector = strings.Join(labelSelectors, ",")
	if input.Type == schema.VariableReferenceSourceTypeService {
		labelSelectors = append(labelSelectors, fmt.Sprintf("%s=%s", "unbind-service", service.ID.String()))
	}

	// Add other services
	for _, svc := range services {
		labelSelectors = append(labelSelectors, fmt.Sprintf("%s=%s", "unbind-service", svc.ID.String()))
	}
	labelSelector := strings.Join(labelSelectors, ",")

	err = self.k8s.RollingRestartPodsByLabels(ctx, team.Namespace, labelSelector, client)
	if err != nil {
		log.Error("Failed to restart pods", "err", err)
		return nil, err
	}

	return variableResponse, nil
}
