package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type TemplateShortResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Version   int       `json:"version"`
	Immutable bool      `json:"immutable"`
	CreatedAt time.Time `json:"created_at"`
}

// TransformTemplateShortEntity transforms an ent.Template entity into a TemplateResponse
func TransformTemplateShortEntity(entity *ent.Template) *TemplateShortResponse {
	response := &TemplateShortResponse{}
	if entity != nil {
		response = &TemplateShortResponse{
			ID:        entity.ID,
			Name:      entity.Name,
			Version:   entity.Version,
			Immutable: entity.Immutable,
			CreatedAt: entity.CreatedAt,
		}
	}
	return response
}

// TransformTemplateShortEntities transforms a slice of ent.Template entities into a slice of TemplateResponse
func TransformTemplateShortEntities(entities []*ent.Template) []*TemplateShortResponse {
	responses := make([]*TemplateShortResponse, len(entities))
	for i, entity := range entities {
		responses[i] = TransformTemplateShortEntity(entity)
	}
	return responses
}

// TemplateWithDefinitionResponse is the response model for a template with its definition
type TemplateWithDefinitionResponse struct {
	ID         uuid.UUID                 `json:"id"`
	Name       string                    `json:"name"`
	Version    int                       `json:"version"`
	Immutable  bool                      `json:"immutable"`
	Definition schema.TemplateDefinition `json:"definition"`
	CreatedAt  time.Time                 `json:"created_at"`
}

// TransformTemplateEntity transforms an ent.Template entity into a TemplateWithDefinitionResponse
func TransformTemplateEntity(entity *ent.Template) *TemplateWithDefinitionResponse {
	response := &TemplateWithDefinitionResponse{}
	for i := range entity.Definition.Services {
		if entity.Definition.Services[i].VariableReferences == nil {
			entity.Definition.Services[i].VariableReferences = []schema.TemplateVariableReference{}
		}
		if entity.Definition.Services[i].Variables == nil {
			entity.Definition.Services[i].Variables = []schema.TemplateVariable{}
		}
		if entity.Definition.Services[i].VariablesMounts == nil {
			entity.Definition.Services[i].VariablesMounts = []*schema.VariableMount{}
		}
		if entity.Definition.Services[i].Volumes == nil {
			entity.Definition.Services[i].Volumes = []schema.TemplateVolume{}
		}
		if entity.Definition.Services[i].DependsOn == nil {
			entity.Definition.Services[i].DependsOn = []int{}
		}
		if entity.Definition.Services[i].HostInputIDs == nil {
			entity.Definition.Services[i].HostInputIDs = []int{}
		}
		if entity.Definition.Services[i].Ports == nil {
			entity.Definition.Services[i].Ports = []schema.PortSpec{}
		}
		if entity.Definition.Services[i].ProtectedVariables == nil {
			entity.Definition.Services[i].ProtectedVariables = []string{}
		}
	}
	if entity != nil {
		response = &TemplateWithDefinitionResponse{
			ID:         entity.ID,
			Name:       entity.Name,
			Version:    entity.Version,
			Definition: entity.Definition,
			Immutable:  entity.Immutable,
			CreatedAt:  entity.CreatedAt,
		}
	}
	return response
}

// TransformTemplateEntities transforms a slice of ent.Template entities into a slice of TemplateWithDefinitionResponse
func TransformTemplateEntities(entities []*ent.Template) []*TemplateWithDefinitionResponse {
	responses := make([]*TemplateWithDefinitionResponse, len(entities))
	for i, entity := range entities {
		responses[i] = TransformTemplateEntity(entity)
	}
	return responses
}
