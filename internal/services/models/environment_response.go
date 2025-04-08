package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
)

type EnvironmentResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"display_name"`
	Description  string    `json:"description"`
	Active       bool      `json:"active"`
	ServiceCount int       `json:"service_count,omitempty"`
	ServiceIcons []string  `json:"service_icons,omitempty" nullable:"false"`
	CreatedAt    time.Time `json:"created_at"`
}

// TransformEnvironmentEntity transforms an ent.Environment entity into an EnvironmentResponse
func TransformEnvironmentEntity(entity *ent.Environment) *EnvironmentResponse {
	response := &EnvironmentResponse{}
	if entity != nil {
		response = &EnvironmentResponse{
			ID:           entity.ID,
			Name:         entity.Name,
			DisplayName:  entity.DisplayName,
			Description:  entity.Description,
			Active:       entity.Active,
			CreatedAt:    entity.CreatedAt,
			ServiceIcons: []string{},
		}
	}
	return response
}

// Transforms a slice of ent.Environment entities into a slice of EnvironmentResponse
func TransformEnvironmentEntitities(entities []*ent.Environment) []*EnvironmentResponse {
	responses := make([]*EnvironmentResponse, len(entities))
	for i, entity := range entities {
		responses[i] = TransformEnvironmentEntity(entity)
	}
	return responses
}
