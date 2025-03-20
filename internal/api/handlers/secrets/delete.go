package secrets_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Delete secrets
type DeleteSecretSecretsInput struct {
	server.BaseAuthInput
	Body struct {
		BaseSecretsJSONInput
		Secrets []models.SecretDeleteInput `json:"secrets" validate:"required"`
	}
}

func (self *HandlerGroup) DeleteSecrets(ctx context.Context, input *DeleteSecretSecretsInput) (*SecretsResponse, error) {
	// Validate input
	if err := ValidateSecretsDependencies(input.Body.Type, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.ServiceID); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	var secret []*models.SecretResponse
	var err error
	// Determine which service to use
	switch input.Body.Type {
	case models.TeamSecret:
		secret, err = self.srv.TeamService.DeleteSecretsByKey(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.Secrets)
	case models.ProjectSecret:
		secret, err = self.srv.ProjectService.DeleteSecretsByKey(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.Secrets)
	case models.EnvironmentSecret:
		secret, err = self.srv.EnvironmentService.DeleteSecretsByKey(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.Secrets)
	case models.ServiceSecret:
		secret, err = self.srv.ServiceService.DeleteSecretsByKey(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.ServiceID, input.Body.Secrets)
	}

	if err != nil {
		return nil, handleSecretErr(err)
	}

	resp := &SecretsResponse{}
	resp.Body.Data = secret
	return resp, nil
}
