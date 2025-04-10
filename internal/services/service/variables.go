package service_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *ServiceService) GetVariables(ctx context.Context, userID uuid.UUID, bearerToken string, teamID, projectID, environmentID, serviceID uuid.UUID) ([]*models.VariableResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   serviceID,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Verify input
	_, project, err := self.VerifyInputs(ctx, teamID, projectID, environmentID)
	if err != nil {
		return nil, err
	}

	// Get service
	service, err := self.repo.Service().GetByID(ctx, serviceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
		}
		return nil, err
	}

	if service.EnvironmentID != environmentID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Service does not belong to the specified project")
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, service.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.ServiceVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	models.SortVariableResponse(variablesResponse)
	return variablesResponse, nil
}

// Create secrets in bulk
func (self *ServiceService) UpdateVariables(ctx context.Context, userID uuid.UUID, bearerToken string, teamID, projectID, environmentID, serviceID uuid.UUID, behavior models.VariableUpdateBehavior, newVariables map[string][]byte) ([]*models.VariableResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   serviceID,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	_, project, err := self.VerifyInputs(ctx, teamID, projectID, environmentID)
	if err != nil {
		return nil, err
	}

	// Get service
	service, err := self.repo.Service().GetByID(ctx, serviceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
		}
		return nil, err
	}

	if service.EnvironmentID != environmentID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Service does not belong to the specified project")
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	if behavior == models.VariableUpdateBehaviorOverwrite {
		// make secrets
		_, err = self.k8s.OverwriteSecretValues(ctx, service.KubernetesSecret, project.Edges.Team.Namespace, newVariables, client)
		if err != nil {
			return nil, err
		}
	} else {
		// make secrets
		_, err = self.k8s.UpsertSecretValues(ctx, service.KubernetesSecret, project.Edges.Team.Namespace, newVariables, client)
		if err != nil {
			return nil, err
		}
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, service.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.ServiceVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	// Perform a restart of pods...
	err = self.k8s.RollingRestartPodsByLabel(ctx, project.Edges.Team.Namespace, "unbind-service", service.Name, client)
	if err != nil {
		log.Error("Failed to restart pods", "err", err, "service", service.Name)
		return nil, err
	}

	models.SortVariableResponse(variablesResponse)
	return variablesResponse, nil
}

// Delete a secret by key
func (self *ServiceService) DeleteVariablesByKey(ctx context.Context, userID uuid.UUID, bearerToken string, teamID, projectID, environmentID, serviceID uuid.UUID, keys []models.VariableDeleteInput) ([]*models.VariableResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionAdmin,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   serviceID,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Verify inputs
	_, project, err := self.VerifyInputs(ctx, teamID, projectID, environmentID)
	if err != nil {
		return nil, err
	}

	// Get service
	service, err := self.repo.Service().GetByID(ctx, serviceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
		}
		return nil, err
	}

	if service.EnvironmentID != environmentID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Service does not belong to the specified project")
	}
	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, service.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	// Remove from map
	for _, secretKey := range keys {
		delete(secrets, secretKey.Name)
	}

	// Update secrets
	_, err = self.k8s.UpdateSecret(ctx, service.KubernetesSecret, project.Edges.Team.Namespace, secrets, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.ServiceVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	// Perform a restart of pods...
	err = self.k8s.RollingRestartPodsByLabel(ctx, project.Edges.Team.Namespace, "unbind-service", service.Name, client)
	if err != nil {
		log.Error("Failed to restart pods", "err", err, "service", service.Name)
		return nil, err
	}

	models.SortVariableResponse(variablesResponse)
	return variablesResponse, nil
}
