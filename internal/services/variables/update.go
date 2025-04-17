package variables_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Create secrets in bulk
func (self *VariablesService) UpdateVariables(ctx context.Context, userID uuid.UUID, bearerToken string, input models.BaseVariablesJSONInput, behavior models.VariableUpdateBehavior, newVariables map[string][]byte) ([]*models.VariableResponse, error) {
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

	if behavior == models.VariableUpdateBehaviorOverwrite {
		// make secrets
		_, err = self.k8s.OverwriteSecretValues(ctx, secretName, team.Namespace, newVariables, client)
		if err != nil {
			return nil, err
		}
	} else {
		// make secrets
		_, err = self.k8s.UpsertSecretValues(ctx, secretName, team.Namespace, newVariables, client)
		if err != nil {
			return nil, err
		}
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, secretName, team.Namespace, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		variablesResponse[i] = &models.VariableResponse{
			Type:  input.Type,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	// Perform a restart of pods...
	// Get label target
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

	models.SortVariableResponse(variablesResponse)
	return variablesResponse, nil
}
