package project_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *ProjectService) GetVariables(ctx context.Context, userID uuid.UUID, bearerToken string, teamID uuid.UUID, projectID uuid.UUID) ([]*models.VariableResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeProject,
			ResourceID:   projectID,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
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

	// Get variables
	variables, err := self.k8s.GetSecretMap(ctx, project.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(variables))
	i := 0
	for k, v := range variables {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.ProjectVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	models.SortVariableResponse(variablesResponse)
	return variablesResponse, nil
}

// Create variables in bulk
func (self *ProjectService) UpsertVariables(ctx context.Context, userID uuid.UUID, bearerToken string, teamID, projectID uuid.UUID, newVariables map[string][]byte) ([]*models.VariableResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeProject,
			ResourceID:   projectID,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
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

	// make variables
	_, err = self.k8s.UpsertSecretValues(ctx, project.KubernetesSecret, project.Edges.Team.Namespace, newVariables, client)
	if err != nil {
		return nil, err
	}

	// Get variables
	variables, err := self.k8s.GetSecretMap(ctx, project.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(variables))
	i := 0
	for k, v := range variables {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.ProjectVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	models.SortVariableResponse(variablesResponse)
	return variablesResponse, nil
}

// Delete a secret by key
func (self *ProjectService) DeleteVariablesByKey(ctx context.Context, userID uuid.UUID, bearerToken string, teamID uuid.UUID, projectID uuid.UUID, keys []models.VariableDeleteInput) ([]*models.VariableResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionAdmin,
			ResourceType: schema.ResourceTypeProject,
			ResourceID:   projectID,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
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

	// Get variables
	variables, err := self.k8s.GetSecretMap(ctx, project.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	// Remove from map
	for _, secretKey := range keys {
		delete(variables, secretKey.Name)
	}

	// Update variables
	_, err = self.k8s.UpdateSecret(ctx, project.KubernetesSecret, project.Edges.Team.Namespace, variables, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(variables))
	i := 0
	for k, v := range variables {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.ProjectVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	models.SortVariableResponse(variablesResponse)
	return variablesResponse, nil
}
