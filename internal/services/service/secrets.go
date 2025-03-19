package service_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *ServiceService) GetSecrets(ctx context.Context, userID uuid.UUID, bearerToken string, teamID, projectID, environmentID, serviceID uuid.UUID) ([]*models.SecretResponse, error) {
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
		// Has permission to read environments
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   "*",
		},
		// Has permission to read this specific environment
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   environmentID.String(),
		},
		// Has permission to read this specific service
		{
			Action:       permission.ActionRead,
			ResourceType: permission.ResourceTypeService,
			ResourceID:   serviceID.String(),
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// ! TODO - repeat this less
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
	secrets, err := self.k8s.GetSecretMap(ctx, service.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	secretResponse := make([]*models.SecretResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		secretResponse[i] = &models.SecretResponse{
			Type:  models.ServiceSecret,
			Key:   k,
			Value: string(v),
		}
		i++
	}

	return secretResponse, nil
}

// Create secrets in bulk
func (self *ServiceService) CreateSecrets(ctx context.Context, userID uuid.UUID, bearerToken string, teamID, projectID, environmentID, serviceID uuid.UUID, newSecrets map[string][]byte) ([]*models.SecretResponse, error) {
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
		// Has permission to read environments
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   "*",
		},
		// Has permission to read this specific environment
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   environmentID.String(),
		},
		// Has permission to read this specific service
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeService,
			ResourceID:   serviceID.String(),
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
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
	_, err = self.k8s.AddSecretValues(ctx, service.KubernetesSecret, project.Edges.Team.Namespace, newSecrets, client)
	if err != nil {
		return nil, err
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, service.KubernetesSecret, project.Edges.Team.Namespace, client)
	if err != nil {
		return nil, err
	}

	secretResponse := make([]*models.SecretResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		secretResponse[i] = &models.SecretResponse{
			Type:  models.ServiceSecret,
			Key:   k,
			Value: string(v),
		}
		i++
	}

	return secretResponse, nil
}

// Delete a secret by key
func (self *ServiceService) DeleteSecretsByKey(ctx context.Context, userID uuid.UUID, bearerToken string, teamID, projectID, environmentID, serviceID uuid.UUID, keys []models.SecretDeleteInput) ([]*models.SecretResponse, error) {
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
		// Has permission to read environments
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   "*",
		},
		// Has permission to read this specific environment
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   environmentID.String(),
		},
		// Has permission to read this specific service
		{
			Action:       permission.ActionEdit,
			ResourceType: permission.ResourceTypeService,
			ResourceID:   serviceID.String(),
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
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

	secretResponse := make([]*models.SecretResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		secretResponse[i] = &models.SecretResponse{
			Type:  models.ServiceSecret,
			Key:   k,
			Value: string(v),
		}
		i++
	}

	return secretResponse, nil
}
