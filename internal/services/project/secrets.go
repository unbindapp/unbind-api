package project_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *ProjectService) GetSecrets(ctx context.Context, userID uuid.UUID, bearerToken string, teamID uuid.UUID, projectID uuid.UUID) ([]*models.SecretResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to read system resources
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		// Has permission to read teams
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to read this specific team
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   teamID.String(),
		},
		// Has permission to read projects
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   "*",
		},
		// Has permission to read this specific project
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   projectID.String(),
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

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, project.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	secretResponse := make([]*models.SecretResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		secretResponse[i] = &models.SecretResponse{
			Type:  models.ProjectSecret,
			Key:   k,
			Value: string(v),
		}
		i++
	}

	return secretResponse, nil
}

// Create secrets in bulk
func (self *ProjectService) CreateSecrets(ctx context.Context, userID uuid.UUID, bearerToken string, teamID, projectID uuid.UUID, newSecrets map[string][]byte) ([]*models.SecretResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to read system resources
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		// Has permission to read teams
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to read this specific team
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   teamID.String(),
		},
		// Has permission to read projects
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   "*",
		},
		// Has permission to read this specific project
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   projectID.String(),
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

	// make secrets
	_, err = self.k8s.AddSecretValues(ctx, project.KubernetesSecret, project.Edges.Team.Namespace, newSecrets, client)
	if err != nil {
		return nil, err
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, project.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	secretResponse := make([]*models.SecretResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		secretResponse[i] = &models.SecretResponse{
			Type:  models.ProjectSecret,
			Key:   k,
			Value: string(v),
		}
		i++
	}

	return secretResponse, nil
}

// Delete a secret by key
func (self *ProjectService) DeleteSecretsByKey(ctx context.Context, userID uuid.UUID, bearerToken string, teamID uuid.UUID, projectID uuid.UUID, keys []string) ([]*models.SecretResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to read system resources
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		// Has permission to read teams
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to read this specific team
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   teamID.String(),
		},
		// Has permission to read projects
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   "*",
		},
		// Has permission to read this specific project
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   projectID.String(),
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

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, project.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	// Remove from map
	for _, secretKey := range keys {
		delete(secrets, secretKey)
	}

	// Update secrets
	_, err = self.k8s.UpdateSecret(ctx, project.KubernetesSecret, project.Edges.Team.Namespace, secrets, client)
	if err != nil {
		return nil, err
	}

	secretResponse := make([]*models.SecretResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		secretResponse[i] = &models.SecretResponse{
			Type:  models.ProjectSecret,
			Key:   k,
			Value: string(v),
		}
		i++
	}

	return secretResponse, nil
}
