package models

import (
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type AvailableVariableReference struct {
	Type           schema.VariableReferenceType       `json:"type"`
	KubernetesName string                             `json:"kubernetes_name"`
	SourceName     string                             `json:"source_name"`
	SourceIcon     string                             `json:"source_icon"`
	SourceType     schema.VariableReferenceSourceType `json:"source_type"`
	SourceID       uuid.UUID                          `json:"source_id"`
	Keys           []string                           `json:"keys"`
}

// Define a comparison function for AvailableVariableReference
func compareAvailableVariableReferences(a, b AvailableVariableReference) int {
	// Sort by type first
	if a.Type != b.Type {
		return strings.Compare(string(a.Type), string(b.Type))
	}

	// Sort by source type first, team comes first
	aPriority := getSourceTypePriority(a.SourceType)
	bPriority := getSourceTypePriority(b.SourceType)
	if aPriority != bPriority {
		return aPriority - bPriority
	}

	// sort by name
	return strings.Compare(a.KubernetesName, b.KubernetesName)
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

func TransformAvailableVariableResponse(secretData []SecretData, endpoints *EndpointDiscovery, kubernetesNameMap map[uuid.UUID]string, nameMap map[uuid.UUID]string, iconMap map[uuid.UUID]string) []AvailableVariableReference {
	resp := make([]AvailableVariableReference, len(secretData)+len(endpoints.Internal)+len(endpoints.External))

	// Process variables
	i := 0
	for _, secret := range secretData {
		resp[i] = AvailableVariableReference{
			Type:           schema.VariableReferenceTypeVariable,
			SourceName:     nameMap[secret.ID],
			SourceIcon:     iconMap[secret.ID],
			KubernetesName: secret.SecretName,
			SourceType:     secret.Type,
			SourceID:       secret.ID,
			Keys:           secret.Keys,
		}
		i++
	}

	// Make endpoints response
	for _, endpoint := range endpoints.Internal {
		resp[i] = AvailableVariableReference{
			Type:           schema.VariableReferenceTypeInternalEndpoint,
			SourceName:     nameMap[endpoint.ServiceID],
			SourceIcon:     iconMap[endpoint.ServiceID],
			KubernetesName: kubernetesNameMap[endpoint.ServiceID],
			SourceType:     schema.VariableReferenceSourceTypeService, // Always service
			SourceID:       endpoint.ServiceID,
			Keys:           []string{endpoint.KubernetesName},
		}
		i++
	}

	for _, endpoint := range endpoints.External {
		resp[i] = AvailableVariableReference{
			Type:           schema.VariableReferenceTypeExternalEndpoint,
			SourceName:     nameMap[endpoint.ServiceID],
			SourceIcon:     iconMap[endpoint.ServiceID],
			KubernetesName: kubernetesNameMap[endpoint.ServiceID],
			SourceType:     schema.VariableReferenceSourceTypeService, // Always service
			SourceID:       endpoint.ServiceID,
		}

		resp[i].Keys = make([]string, len(endpoint.Hosts))
		for j, host := range endpoint.Hosts {
			resp[i].Keys[j] = host.Host
		}
		i++
	}

	slices.SortFunc(resp, compareAvailableVariableReferences)

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
