package models

import (
	"github.com/google/uuid"
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
	Values     map[string]string                  `json:"values"`
}

// SecretData represents a Kubernetes secret with its metadata
type SecretData struct {
	ID         uuid.UUID
	Type       schema.VariableReferenceSourceType
	SecretName string
	Data       map[string][]byte
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
		}

		resp.Variables[i].Values = make(map[string]string)
		for k, v := range secret.Data {
			resp.Variables[i].Values[k] = string(v)
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
			Values: map[string]string{
				endpoint.Name: endpoint.DNS,
			},
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

		resp.ExternalEndpoints[i].Values = make(map[string]string)
		for _, host := range endpoint.Hosts {
			resp.ExternalEndpoints[i].Values[endpoint.Name] = host.Host
		}
	}

	return resp
}
