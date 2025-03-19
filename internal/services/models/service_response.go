package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/service"
)

// ServiceResponse defines the response structure for service operations
type ServiceResponse struct {
	ID                   uuid.UUID              `json:"id"`
	Name                 string                 `json:"name"`
	DisplayName          string                 `json:"display_name"`
	Description          string                 `json:"description"`
	Type                 service.Type           `json:"type"`
	Builder              service.Builder        `json:"builder"`
	Runtime              *string                `json:"runtime,omitempty"`
	Framework            *string                `json:"framework,omitempty"`
	EnvironmentID        uuid.UUID              `json:"environment_id"`
	GitHubInstallationID *int64                 `json:"github_installation_id,omitempty"`
	GitRepository        *string                `json:"git_repository,omitempty"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	Config               *ServiceConfigResponse `json:"config"`
}

// TransformServiceEntity transforms an ent.Service entity into a ServiceResponse
func TransformServiceEntity(entity *ent.Service) *ServiceResponse {
	response := &ServiceResponse{}
	if entity != nil {
		response = &ServiceResponse{
			ID:                   entity.ID,
			Name:                 entity.Name,
			DisplayName:          entity.DisplayName,
			Description:          entity.Description,
			Type:                 entity.Type,
			Builder:              entity.Builder,
			Runtime:              entity.Runtime,
			Framework:            entity.Framework,
			EnvironmentID:        entity.EnvironmentID,
			GitHubInstallationID: entity.GithubInstallationID,
			GitRepository:        entity.GitRepository,
			CreatedAt:            entity.CreatedAt,
			UpdatedAt:            entity.UpdatedAt,
			Config:               TransformServiceConfigEntity(entity.Edges.ServiceConfig),
		}
	}
	return response
}

// TransformServiceEntities transforms a slice of ent.Service entities into a slice of ServiceResponse
func TransformServiceEntities(entities []*ent.Service) []*ServiceResponse {
	responses := make([]*ServiceResponse, len(entities))
	for i, entity := range entities {
		responses[i] = TransformServiceEntity(entity)
	}
	return responses
}
