package secrets_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// List all
type ListSecretsInput struct {
	server.BaseAuthInput
	BaseSecretsInput
}

type SecretsResponse struct {
	Body struct {
		Data []*models.SecretResponse `json:"data"`
	}
}

func (self *HandlerGroup) ListSecrets(ctx context.Context, input *ListSecretsInput) (*SecretsResponse, error) {
	// Validate input
	if err := ValidateSecretsDependencies(input.Type, input.TeamID, input.ProjectID, input.EnvironmentID, input.ServiceID); err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	// Get team secrets
	var err error
	var teamSecrets []*models.SecretResponse
	var projectSecrets []*models.SecretResponse
	var environmentSecrets []*models.SecretResponse
	var serviceSecrets []*models.SecretResponse

	// Get team secrets always
	teamSecrets, err = self.srv.TeamService.GetSecrets(ctx, user.ID, bearerToken, input.TeamID)
	if err != nil {
		return nil, handleSecretErr(err)
	}

	if input.Type == models.ProjectSecret ||
		input.Type == models.EnvironmentSecret ||
		input.Type == models.ServiceSecret {
		projectSecrets, err = self.srv.ProjectService.GetSecrets(ctx, user.ID, bearerToken, input.TeamID, input.ProjectID)
		if err != nil {
			return nil, handleSecretErr(err)
		}
	}

	if input.Type == models.EnvironmentSecret ||
		input.Type == models.ServiceSecret {
		environmentSecrets, err = self.srv.EnvironmentService.GetSecrets(ctx, user.ID, bearerToken, input.TeamID, input.ProjectID, input.EnvironmentID)
		if err != nil {
			return nil, handleSecretErr(err)
		}
	}

	if input.Type == models.ServiceSecret {
		serviceSecrets, err = self.srv.ServiceService.GetSecrets(ctx, user.ID, bearerToken, input.TeamID, input.ProjectID, input.EnvironmentID, input.ServiceID)
		if err != nil {
			return nil, handleSecretErr(err)
		}
	}

	// Combine
	secrets := append(teamSecrets, projectSecrets...)
	secrets = append(secrets, environmentSecrets...)
	secrets = append(secrets, serviceSecrets...)

	resp := &SecretsResponse{}
	resp.Body.Data = secrets
	return resp, nil
}
