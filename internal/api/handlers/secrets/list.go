package secrets_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
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
	if err := ValidateSecretsDependencies(input); err != nil {
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
	secrets, err := self.srv.TeamService.GetSecrets(ctx, user.ID, bearerToken, input.TeamID)
	if err != nil {
		return nil, handleSecretErr(err)
	}

	// Get project secrets
	projectSecrets := []*models.SecretResponse{}
	if input.ProjectID != uuid.Nil {
		projectSecrets, err = self.srv.ProjectService.GetSecrets(ctx, user.ID, bearerToken, input.TeamID, input.ProjectID)
		if err != nil {
			return nil, handleSecretErr(err)
		}
	}

	// Get environment secrets
	environmentSecrets := []*models.SecretResponse{}
	if input.EnvironmentID != uuid.Nil {
		environmentSecrets, err = self.srv.EnvironmentService.GetSecrets(ctx, user.ID, bearerToken, input.TeamID, input.ProjectID, input.EnvironmentID)
		if err != nil {
			return nil, handleSecretErr(err)
		}
	}

	// Combine the secrets
	for _, secret := range projectSecrets {
		secrets = append(secrets, secret)
	}
	for _, secret := range environmentSecrets {
		secrets = append(secrets, secret)
	}

	resp := &SecretsResponse{}
	resp.Body.Data = secrets
	return &SecretsResponse{}, nil
}
