package models

import (
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type AvailableVariableReferenceResponse struct {
	Variables         []AvailableVariableReference `json:"variables" nullable:"false"`
	InternalEndpoints []AvailableVariableReference `json:"internal_endpoints" nullable:"false"`
	ExternalEndpoints []AvailableVariableReference `json:"external_endpoints" nullable:"false"`
}

type AvailableVariableReference struct {
	Type       schema.VariableReferenceType       `json:"type"`
	Name       string                             `json:"name"`
	SourceType schema.VariableReferenceSourceType `json:"source_type"`
	SourceID   uuid.UUID                          `json:"source_id"`
	Keys       []string                           `json:"keys"`
}

// Define a comparison function for AvailableVariableReference
func compareAvailableVariableReferences(a, b AvailableVariableReference) int {
	// Sort by source type first, team comes first
	aPriority := getSourceTypePriority(a.SourceType)
	bPriority := getSourceTypePriority(b.SourceType)
	if aPriority != bPriority {
		return aPriority - bPriority
	}

	// sort by name
	return strings.Compare(a.Name, b.Name)
}

// Helper function to get priority of SourceType
func getSourceTypePriority(sourceType schema.VariableReferenceSourceType) int {
	switch sourceType {
	case schema.VariableReferenceSourceTypeTeam:
		return 0
	case schema.VariableReferenceSourceTypeProject:
		return 1
	case schema.VariableReferenceSourceTypeEnvironment:
		return 2
	case schema.VariableReferenceSourceTypeService:
		return 3
	default:
		return 4
	}
}

// SecretData represents a Kubernetes secret with its metadata
type SecretData struct {
	ID         uuid.UUID
	Type       schema.VariableReferenceSourceType
	SecretName string
	Keys       []string
}

func TransformAvailableVariableResponse(secretData []SecretData, endpoints *EndpointDiscovery) *AvailableVariableReferenceResponse {
	resp := &AvailableVariableReferenceResponse{}
	resp.Variables = make([]AvailableVariableReference, len(secretData))

	// Process variables
	for i, secret := range secretData {
		resp.Variables[i] = AvailableVariableReference{
			Type:       schema.VariableReferenceTypeVariable,
			Name:       secret.SecretName,
			SourceType: secret.Type,
			SourceID:   secret.ID,
			Keys:       secret.Keys,
		}
	}

	// Make endpoints response
	resp.InternalEndpoints = make([]AvailableVariableReference, len(endpoints.Internal))
	for i, endpoint := range endpoints.Internal {
		resp.InternalEndpoints[i] = AvailableVariableReference{
			Type:       schema.VariableReferenceTypeInternalEndpoint,
			Name:       endpoint.Name,
			SourceType: schema.VariableReferenceSourceTypeService, // Always service
			SourceID:   endpoint.ServiceID,
			Keys:       []string{endpoint.Name},
		}
	}
	resp.ExternalEndpoints = make([]AvailableVariableReference, len(endpoints.External))
	for i, endpoint := range endpoints.External {
		resp.ExternalEndpoints[i] = AvailableVariableReference{
			Type:       schema.VariableReferenceTypeExternalEndpoint,
			Name:       endpoint.Name,
			SourceType: schema.VariableReferenceSourceTypeService, // Always service
			SourceID:   endpoint.ServiceID,
		}

		resp.ExternalEndpoints[i].Keys = make([]string, len(endpoint.Hosts))
		for j, host := range endpoint.Hosts {
			resp.ExternalEndpoints[i].Keys[j] = host.Host
		}
	}

	slices.SortFunc(resp.Variables, compareAvailableVariableReferences)
	slices.SortFunc(resp.InternalEndpoints, compareAvailableVariableReferences)
	slices.SortFunc(resp.ExternalEndpoints, compareAvailableVariableReferences)

	return resp
}

// The actual response object
type VariableReferenceResponse struct {
	ID              uuid.UUID                        `json:"id" doc:"The ID of the variable reference" required:"true"`
	TargetServiceID uuid.UUID                        `json:"target_service_id" required:"true"`
	TargetName      string                           `json:"target_name" required:"true"`
	Sources         []schema.VariableReferenceSource `json:"sources" required:"true" nullable:"false"`
	ValueTemplate   string                           `json:"value_template" required:"true"`
	CreatedAt       time.Time                        `json:"created_at" required:"true"`
}

func TransformVariableReferenceResponseEntity(entity *ent.VariableReference) *VariableReferenceResponse {
	sources := entity.Sources
	if sources == nil {
		sources = []schema.VariableReferenceSource{}
	}
	return &VariableReferenceResponse{
		ID:              entity.ID,
		TargetServiceID: entity.TargetServiceID,
		TargetName:      entity.TargetName,
		Sources:         sources,
		ValueTemplate:   entity.ValueTemplate,
		CreatedAt:       entity.CreatedAt,
	}
}

func TransformVariableReferenceResponseEntities(entities []*ent.VariableReference) []*VariableReferenceResponse {
	responses := make([]*VariableReferenceResponse, len(entities))
	for i, entity := range entities {
		responses[i] = TransformVariableReferenceResponseEntity(entity)
	}
	return responses
}
