package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type BuildJobResponse struct {
	ID          uuid.UUID             `json:"id"`
	ServiceID   uuid.UUID             `json:"service_id"`
	Status      schema.BuildJobStatus `json:"status"`
	Error       string                `json:"error,omitempty"`
	StartedAt   *time.Time            `json:"started_at,omitempty"`
	CompletedAt *time.Time            `json:"completed_at,omitempty"`
	Attempts    int                   `json:"attempts"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// TransformBuildJobEntity transforms an ent.BuildJob entity into a BuildJobResponse
func TransformBuildJobEntity(entity *ent.BuildJob) *BuildJobResponse {
	response := &BuildJobResponse{}
	if entity != nil {
		response = &BuildJobResponse{
			ID:          entity.ID,
			ServiceID:   entity.ServiceID,
			Status:      entity.Status,
			Error:       entity.Error,
			StartedAt:   entity.StartedAt,
			CompletedAt: entity.CompletedAt,
			Attempts:    entity.Attempts,
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
		}
	}
	return response
}

// Transforms a slice of ent.BuildJob entities into a slice of BuildJobResponse
func TransformBuildJobEntities(entities []*ent.BuildJob) []*BuildJobResponse {
	responses := make([]*BuildJobResponse, len(entities))
	for i, entity := range entities {
		responses[i] = TransformBuildJobEntity(entity)
	}
	return responses
}
