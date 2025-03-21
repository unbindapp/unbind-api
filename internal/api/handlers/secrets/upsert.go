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
type UpsertSecretsInput struct {
	server.BaseAuthInput
	Body struct {
		BaseSecretsJSONInput
		IsBuildSecret bool `json:"is_build_secret" required:"false"`
		Secrets       []*struct {
			Name  string `json:"name" required:"true"`
			Value string `json:"value" required:"true"`
		} `json:"secrets" required:"true"`
	}
}

func (self *HandlerGroup) UpsertSecrets(ctx context.Context, input *UpsertSecretsInput) (*SecretsResponse, error) {
	// Validate input
	if err := ValidateSecretsDependencies(input.Body.Type, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.ServiceID); err != nil {
		return nil, huma.Error400BadRequest("invalid input", err)
	}

	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	secretUpdateMap := make(map[string][]byte)
	for _, secret := range input.Body.Secrets {
		secretUpdateMap[secret.Name] = []byte(secret.Value)
	}

	// Determine which service to use
	var secret []*models.SecretResponse
	var err error
	switch input.Body.Type {
	case models.TeamSecret:
		secret, err = self.srv.TeamService.UpsertSecrets(ctx, user.ID, bearerToken, input.Body.TeamID, secretUpdateMap)
	case models.ProjectSecret:
		secret, err = self.srv.ProjectService.UpsertSecrets(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, secretUpdateMap)
	case models.EnvironmentSecret:
		secret, err = self.srv.EnvironmentService.UpsertSecrets(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, secretUpdateMap)
	case models.ServiceSecret:
		secret, err = self.srv.ServiceService.UpsertSecrets(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.EnvironmentID, input.Body.ServiceID, secretUpdateMap, input.Body.IsBuildSecret)
	}

	if err != nil {
		return nil, handleSecretErr(err)
	}

	resp := &SecretsResponse{}
	resp.Body.Data = secret
	return resp, nil
}
