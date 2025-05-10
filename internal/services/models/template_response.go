package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type TemplateResponse struct {
	ID         uuid.UUID                 `json:"id"`
	Name       string                    `json:"name"`
	Version    int                       `json:"version"`
	Immutable  bool                      `json:"immutable"`
	Definition schema.TemplateDefinition `json:"definition"`
	CreatedAt  time.Time                 `json:"created_at"`
}

// TransformTemplateEntity transforms an ent.Template entity into a TemplateResponse
func TransformTemplateEntity(entity *ent.Template) *TemplateResponse {
	response := &TemplateResponse{}
	if entity != nil {
		response = &TemplateResponse{
			ID:         entity.ID,
			Name:       entity.Name,
			Version:    entity.Version,
			Immutable:  entity.Immutable,
			Definition: entity.Definition,
			CreatedAt:  entity.CreatedAt,
		}
	}
	return response
}

// TransformTemplateEntities transforms a slice of ent.Template entities into a slice of TemplateResponse
func TransformTemplateEntities(entities []*ent.Template) []*TemplateResponse {
	responses := make([]*TemplateResponse, len(entities))
	for i, entity := range entities {
		responses[i] = TransformTemplateEntity(entity)
	}
	return responses
}
