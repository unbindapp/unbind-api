package variables_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
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

		// Remove from map
		for _, secretKey := range keys {
			// Don't allow deletion of special database keys
			if service.Type == schema.ServiceTypeDatabase &&
				(secretKey.Name == "DATABASE_USERNAME" ||
					secretKey.Name == "DATABASE_PASSWORD" ||
					secretKey.Name == "DATABASE_URL") {
				continue
			}
			delete(secrets, secretKey.Name)
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

	// Perform a restart of pods...
	var labelValue string
	switch input.Type {
	case schema.VariableReferenceSourceTypeTeam:
		labelValue = team.ID.String()
	case schema.VariableReferenceSourceTypeProject:
		labelValue = project.ID.String()
	case schema.VariableReferenceSourceTypeEnvironment:
		labelValue = environment.ID.String()
	case schema.VariableReferenceSourceTypeService:
		labelValue = service.ID.String()
	default:
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Invalid variable type")
	}

	err = self.k8s.RollingRestartPodsByLabel(ctx, team.Namespace, input.Type.KubernetesLabel(), labelValue, client)
	if err != nil {
		log.Error("Failed to restart pods", "err", err, "label", input.Type.KubernetesLabel(), "value", labelValue)
		return nil, err
	}

	return variableResponse, nil
}
