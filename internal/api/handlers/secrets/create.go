package secrets_handler

import (
	"context"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Create new
type CreateSecretsInput struct {
	server.BaseAuthInput
	Body struct {
		BaseSecretsJSONInput
		Secrets []*struct {
			Name  string `json:"name" validate:"required"`
			Value string `json:"value" validate:"required"`
		} `json:"secrets" validate:"required"`
	}
}

func (self *HandlerGroup) CreateSecrets(ctx context.Context, input *CreateSecretsInput) (*SecretsResponse, error) {
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

	newSecretMap := make(map[string][]byte)
	for _, secret := range input.Body.Secrets {
		newSecretMap[secret.Name] = []byte(secret.Value)
	}

	// Determine which service to use
	var secret []*models.SecretResponse
	var err error
	switch input.Body.Type {
	case models.TeamSecret:
		secret, err = self.srv.TeamService.CreateSecrets(ctx, user.ID, bearerToken, input.Body.TeamID, newSecretMap)
	case models.ProjectSecret:
		secret, err = self.srv.ProjectService.CreateSecrets(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, newSecretMap)
	case models.EnvironmentSecret:
		secret, err = self.srv.EnvironmentService.CreateSecrets(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, newSecretMap)
	case models.ServiceSecret:
		secret, err = self.srv.ServiceService.CreateSecrets(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.ServiceID, newSecretMap)
	}

	if err != nil {
		return nil, handleSecretErr(err)
	}

	resp := &SecretsResponse{}
	resp.Body.Data = secret
	return resp, nil
}
