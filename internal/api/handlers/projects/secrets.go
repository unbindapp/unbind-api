package projects_handler

import (
	"context"
	"errors"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// List all
type ListProjectSecretsInput struct {
	server.BaseAuthInput
	TeamID    uuid.UUID `query:"team_id" required:"true"`
	ProjectID uuid.UUID `query:"project_id" required:"true"`
}

type SecretsResponse struct {
	Body struct {
		Data []*models.SecretResponse `json:"data"`
	}
}

func handleSecretErr(err error) error {
	if errors.Is(err, errdefs.ErrUnauthorized) {
		return huma.Error403Forbidden("Unauthorized")
	}
	if errors.Is(err, errdefs.ErrNotFound) {
		return huma.Error404NotFound(err.Error())
	}
	log.Error("Error getting secrets", "err", err)
	return huma.Error500InternalServerError("Unable to retrieve secrets")
}

func (self *HandlerGroup) ListSecrets(ctx context.Context, input *ListProjectSecretsInput) (*SecretsResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	// Get team secrets
	teamSecrets, err := self.srv.TeamService.GetSecrets(ctx, user.ID, bearerToken, input.TeamID)
	if err != nil {
		return nil, handleSecretErr(err)
	}

	secrets, err := self.srv.ProjectService.GetSecrets(ctx, user.ID, bearerToken, input.TeamID, input.ProjectID)
	if err != nil {
		return nil, handleSecretErr(err)
	}

	// Combine the two sets of secrets
	for _, secret := range teamSecrets {
		secrets = append(secrets, secret)
	}

	resp := &SecretsResponse{}
	resp.Body.Data = secrets
	return resp, nil
}

// Create new
type CreateProjectSecretsInput struct {
	server.BaseAuthInput
	Body struct {
		TeamID    uuid.UUID `json:"team_id" validate:"required"`
		ProjectID uuid.UUID `json:"project_id" validate:"required"`
		Secrets   []*struct {
			Name  string `json:"name" validate:"required"`
			Value string `json:"value" validate:"required"`
		} `json:"secrets" validate:"required"`
	}
}

func (self *HandlerGroup) CreateSecrets(ctx context.Context, input *CreateProjectSecretsInput) (*SecretsResponse, error) {
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
	secret, err := self.srv.ProjectService.CreateSecrets(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, newSecretMap)
	if err != nil {
		return nil, handleSecretErr(err)
	}

	resp := &SecretsResponse{}
	resp.Body.Data = secret
	return resp, nil
}

// Delete secrets
type DeleteProjectSecretsInput struct {
	server.BaseAuthInput
	Body struct {
		TeamID    uuid.UUID                  `json:"team_id" validate:"required"`
		ProjectID uuid.UUID                  `json:"project_id" validate:"required"`
		Secrets   []models.SecretDeleteInput `json:"secrets" validate:"required"`
	}
}

func (self *HandlerGroup) DeleteSecrets(ctx context.Context, input *DeleteProjectSecretsInput) (*SecretsResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	secret, err := self.srv.ProjectService.DeleteSecretsByKey(ctx, user.ID, bearerToken, input.Body.TeamID, input.Body.ProjectID, input.Body.Secrets)
	if err != nil {
		return nil, handleSecretErr(err)
	}

	resp := &SecretsResponse{}
	resp.Body.Data = secret
	return resp, nil
}
