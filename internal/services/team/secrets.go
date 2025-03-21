package team_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *TeamService) GetSecrets(ctx context.Context, userID uuid.UUID, bearerToken string, teamID uuid.UUID) ([]*models.SecretResponse, error) {
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
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Get team
	team, err := self.repo.Team().GetByID(ctx, teamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, err
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, team.KubernetesSecret, team.Namespace, client)
	if err != nil {
		return nil, err
	}

	secretResponse := make([]*models.SecretResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		secretResponse[i] = &models.SecretResponse{
			Type:  models.TeamSecret,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	return secretResponse, nil
}

// Create secrets in bulk
func (self *TeamService) UpsertSecrets(ctx context.Context, userID uuid.UUID, bearerToken string, teamID uuid.UUID, newSecrets map[string][]byte) ([]*models.SecretResponse, error) {
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
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Get team
	team, err := self.repo.Team().GetByID(ctx, teamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, err
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// make secrets
	_, err = self.k8s.UpsertSecretValues(ctx, team.KubernetesSecret, team.Namespace, newSecrets, client)
	if err != nil {
		return nil, err
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, team.KubernetesSecret, team.Namespace, client)
	if err != nil {
		return nil, err
	}

	secretResponse := make([]*models.SecretResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		secretResponse[i] = &models.SecretResponse{
			Type:  models.TeamSecret,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	return secretResponse, nil
}

// Delete a secret by key
func (self *TeamService) DeleteSecretsByKey(ctx context.Context, userID uuid.UUID, bearerToken string, teamID uuid.UUID, keys []models.SecretDeleteInput) ([]*models.SecretResponse, error) {
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
	}

	if err := self.repo.Permissions().Check(
		ctx,
		userID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Get team
	team, err := self.repo.Team().GetByID(ctx, teamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, err
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Get secrets
	secrets, err := self.k8s.GetSecretMap(ctx, team.KubernetesSecret, team.Namespace, client)
	if err != nil {
		return nil, err
	}

	// Remove from map
	for _, secretKey := range keys {
		delete(secrets, secretKey.Name)
	}

	// Update secrets
	_, err = self.k8s.UpdateSecret(ctx, team.KubernetesSecret, team.Namespace, secrets, client)
	if err != nil {
		return nil, err
	}

	secretResponse := make([]*models.SecretResponse, len(secrets))
	i := 0
	for k, v := range secrets {
		secretResponse[i] = &models.SecretResponse{
			Type:  models.TeamSecret,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	return secretResponse, nil
}
