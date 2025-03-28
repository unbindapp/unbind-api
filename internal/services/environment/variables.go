package environment_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *EnvironmentService) GetVariables(ctx context.Context, userID uuid.UUID, bearerToken string, teamID, projectID, environmentID uuid.UUID) ([]*models.VariableResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeEnvironment,
			ResourceID:   environmentID,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Get environment
	environment, err := self.repo.Environment().GetByID(ctx, environmentID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Environment not found")
		}
		return nil, err
	}

	if environment.ProjectID != projectID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Environment does not belong to the specified project")
	}

	// Get project
	project, err := self.repo.Project().GetByID(ctx, projectID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Project not found")
		}
		return nil, err
	}

	if project.TeamID != teamID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Project does not belong to the specified team")
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, environment.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.EnvironmentVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	models.SortVariableResponse(variablesResponse)
	return variablesResponse, nil
}

// Create secrets in bulk
func (self *EnvironmentService) UpsertVariables(ctx context.Context, userID uuid.UUID, bearerToken string, teamID, projectID, environmentID uuid.UUID, newVariables map[string][]byte) ([]*models.VariableResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeEnvironment,
			ResourceID:   environmentID,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Get environment
	environment, err := self.repo.Environment().GetByID(ctx, environmentID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Environment not found")
		}
		return nil, err
	}

	if environment.ProjectID != projectID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Environment does not belong to the specified project")
	}

	// Get project
	project, err := self.repo.Project().GetByID(ctx, projectID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Project not found")
		}
		return nil, err
	}

	if project.TeamID != teamID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Project does not belong to the specified team")
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// make secrets
	_, err = self.k8s.UpsertSecretValues(ctx, environment.KubernetesSecret, project.Edges.Team.Namespace, newVariables, client)
	if err != nil {
		return nil, err
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, environment.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.EnvironmentVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	models.SortVariableResponse(variablesResponse)
	return variablesResponse, nil
}

// Delete a secret by key
func (self *EnvironmentService) DeleteVariablesByKey(ctx context.Context, userID uuid.UUID, bearerToken string, teamID, projectID, environmentID uuid.UUID, keys []models.VariableDeleteInput) ([]*models.VariableResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeEnvironment,
			ResourceID:   environmentID,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Get environment
	environment, err := self.repo.Environment().GetByID(ctx, environmentID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Environment not found")
		}
		return nil, err
	}

	if environment.ProjectID != projectID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Environment does not belong to the specified project")
	}

	// Get project
	project, err := self.repo.Project().GetByID(ctx, projectID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Project not found")
		}
		return nil, err
	}

	if project.TeamID != teamID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Project does not belong to the specified team")
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, environment.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	// Remove from map
	for _, secretKey := range keys {
		delete(secrets, secretKey.Name)
	}

	// Update secrets
	_, err = self.k8s.UpdateSecret(ctx, environment.KubernetesSecret, project.Edges.Team.Namespace, secrets, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.EnvironmentVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	models.SortVariableResponse(variablesResponse)
	return variablesResponse, nil
}
