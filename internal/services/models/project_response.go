package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
)

type ProjectResponse struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Description  *string                `json:"description"`
	Status       string                 `json:"status"`
	TeamID       uuid.UUID              `json:"team_id"`
	CreatedAt    time.Time              `json:"created_at"`
	Environments []*EnvironmentResponse `json:"environments" nullable:"false"`
}

func (self *ProjectResponse) AttachServiceSummary(counts map[uuid.UUID]int, providerSummaries map[uuid.UUID][]enum.Provider, frameworkSummaries map[uuid.UUID][]enum.Framework) {
	for _, environment := range self.Environments {
		if count, ok := counts[environment.ID]; ok {
			environment.ServiceCount = count
		}
		if providerSummary, ok := providerSummaries[environment.ID]; ok {
			environment.ProviderSummary = providerSummary
		}
		if frameworkSummary, ok := frameworkSummaries[environment.ID]; ok {
			environment.FrameworkSummary = frameworkSummary
		}
	}
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
