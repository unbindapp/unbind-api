package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
)

type TeamResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// TransformTeamEntity transforms an ent.Team entity into a TeamResponse
func TransformTeamEntity(entity *ent.Team) *TeamResponse {
	response := &TeamResponse{}
	if entity != nil {
		response = &TeamResponse{
			ID:          entity.ID,
			Name:        entity.Name,
			DisplayName: entity.DisplayName,
			Description: entity.Description,
			CreatedAt:   entity.CreatedAt,
		}
	}
	return response
}

// TransformTeamEntities transforms a slice of ent.Team entities into a slice of TeamResponse
func TransformTeamEntities(entities []*ent.Team) []*TeamResponse {
	responses := make([]*TeamResponse, len(entities))
	for i, entity := range entities {
		responses[i] = TransformTeamEntity(entity)
	}
	return responses
}
