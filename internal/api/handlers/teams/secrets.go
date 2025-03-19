package teams_handler

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
type ListTeamSecretsInput struct {
	server.BaseAuthInput
	TeamID uuid.UUID `query:"team_id" required:"true"`
}

type SecretsResponse struct {
	Body struct {
		Data []*models.SecretResponse `json:"data"`
	}
}

func (self *HandlerGroup) ListSecrets(ctx context.Context, input *ListTeamSecretsInput) (*SecretsResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	secrets, err := self.srv.TeamService.GetSecrets(ctx, user.ID, bearerToken, input.TeamID)
	if err != nil {
		if errors.Is(err, errdefs.ErrUnauthorized) {
			return nil, huma.Error403Forbidden("Unauthorized")
		}
		if errors.Is(err, errdefs.ErrNotFound) {
			return nil, huma.Error404NotFound(err.Error())
		}
		log.Error("Error getting secrets", "err", err)
		return nil, huma.Error500InternalServerError("Unable to retrieve secrets")
	}

	resp := &SecretsResponse{}
	resp.Body.Data = secrets
	return resp, nil
}

// Add new
type AddTeamSecretInput struct {
	server.BaseAuthInput
	Body struct {
		TeamID  uuid.UUID `json:"team_id" validate:"required"`
		Secrets []*struct {
			Name  string `json:"name" validate:"required"`
			Value string `json:"value" validate:"required"`
		} `json:"secrets" validate:"required"`
	}
}

func (self *HandlerGroup) AddSecret(ctx context.Context, input *AddTeamSecretInput) (*SecretsResponse, error) {
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
	secret, err := self.srv.TeamService.CreateSecrets(ctx, user.ID, bearerToken, input.Body.TeamID, newSecretMap)
	if err != nil {
		if errors.Is(err, errdefs.ErrUnauthorized) {
			return nil, huma.Error403Forbidden("Unauthorized")
		}
		if errors.Is(err, errdefs.ErrNotFound) {
			return nil, huma.Error404NotFound(err.Error())
		}
		log.Error("Error getting secrets", "err", err)
		return nil, huma.Error500InternalServerError("Unable to retrieve secrets")
	}

	resp := &SecretsResponse{}
	resp.Body.Data = secret
	return resp, nil
}

// Delete one
type DeleteTeamSecretInput struct {
	server.BaseAuthInput
	Body struct {
		TeamID  uuid.UUID `json:"team_id" validate:"required"`
		Secrets []*struct {
			Name string `json:"name" validate:"required"`
		} `json:"secrets" validate:"required"`
	}
}

func (self *HandlerGroup) DeleteSecret(ctx context.Context, input *DeleteTeamSecretInput) (*SecretsResponse, error) {
	// Get caller
	user, found := self.srv.GetUserFromContext(ctx)
	if !found {
		log.Error("Error getting user from context")
		return nil, huma.Error401Unauthorized("Unable to retrieve user")
	}
	bearerToken := strings.TrimPrefix(input.Authorization, "Bearer ")

	keysToDelete := make([]string, len(input.Body.Secrets))
	for i, secret := range input.Body.Secrets {
		keysToDelete[i] = secret.Name
	}
	secret, err := self.srv.TeamService.DeleteSecretsByKey(ctx, user.ID, bearerToken, input.Body.TeamID, keysToDelete)
	if err != nil {
		if errors.Is(err, errdefs.ErrUnauthorized) {
			return nil, huma.Error403Forbidden("Unauthorized")
		}
		if errors.Is(err, errdefs.ErrNotFound) {
			return nil, huma.Error404NotFound(err.Error())
		}
		log.Error("Error getting secrets", "err", err)
		return nil, huma.Error500InternalServerError("Unable to retrieve secrets")
	}

	resp := &SecretsResponse{}
	resp.Body.Data = secret
	return resp, nil
}
