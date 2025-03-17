package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
)

type ProjectResponse struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Description  string                 `json:"description"`
	Status       string                 `json:"status"`
	TeamID       uuid.UUID              `json:"team_id"`
	CreatedAt    time.Time              `json:"created_at"`
	Environments []*EnvironmentResponse `json:"environments"`
}

// TransformProjectEntity transforms an ent.Project entity into a ProjectResponse
func TransformProjectEntity(entity *ent.Project) *ProjectResponse {
	response := &ProjectResponse{}
	if entity != nil {
		response = &ProjectResponse{
			ID:           entity.ID,
			Name:         entity.Name,
			DisplayName:  entity.DisplayName,
			Description:  entity.Description,
			Status:       entity.Status,
			TeamID:       entity.TeamID,
			CreatedAt:    entity.CreatedAt,
			Environments: TransformEnvironmentEntitities(entity.Edges.Environments),
		}
	}
	return response
}

// TransformProjectEntitities transforms a slice of ent.Project entities into a slice of ProjectResponse
func TransformProjectEntitities(entities []*ent.Project) []*ProjectResponse {
	responses := make([]*ProjectResponse, len(entities))
	for i, entity := range entities {
		responses[i] = TransformProjectEntity(entity)
	}
	return responses
}
