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

func (self *TeamService) GetVariables(ctx context.Context, userID uuid.UUID, bearerToken string, teamID uuid.UUID) ([]*models.VariableResponse, error) {
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

	// Get variables
	variables, err := self.k8s.GetSecretMap(ctx, team.KubernetesSecret, team.Namespace, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(variables))
	i := 0
	for k, v := range variables {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.TeamVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	return variablesResponse, nil
}

// Create variables in bulk
func (self *TeamService) UpsertVariables(ctx context.Context, userID uuid.UUID, bearerToken string, teamID uuid.UUID, newVariables map[string][]byte) ([]*models.VariableResponse, error) {
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

	// make variables
	_, err = self.k8s.UpsertSecretValues(ctx, team.KubernetesSecret, team.Namespace, newVariables, client)
	if err != nil {
		return nil, err
	}

	// Get variables
	variables, err := self.k8s.GetSecretMap(ctx, team.KubernetesSecret, team.Namespace, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(variables))
	i := 0
	for k, v := range variables {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.TeamVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	return variablesResponse, nil
}

// Delete a secret by key
func (self *TeamService) DeleteVariablesByKey(ctx context.Context, userID uuid.UUID, bearerToken string, teamID uuid.UUID, keys []models.VariableDeleteInput) ([]*models.VariableResponse, error) {
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

	// Get variables
	variables, err := self.k8s.GetSecretMap(ctx, team.KubernetesSecret, team.Namespace, client)
	if err != nil {
		return nil, err
	}

	// Remove from map
	for _, secretKey := range keys {
		delete(variables, secretKey.Name)
	}

	// Update variables
	_, err = self.k8s.UpdateSecret(ctx, team.KubernetesSecret, team.Namespace, variables, client)
	if err != nil {
		return nil, err
	}

	variablesResponse := make([]*models.VariableResponse, len(variables))
	i := 0
	for k, v := range variables {
		variablesResponse[i] = &models.VariableResponse{
			Type:  models.TeamVariable,
			Name:  k,
			Value: string(v),
		}
		i++
	}

	return variablesResponse, nil
}
