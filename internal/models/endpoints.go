package models

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
)

type EndpointDiscovery struct {
	Internal []ServiceEndpoint `json:"internal" nullable:"false"`
	External []IngressEndpoint `json:"external" nullable:"false"`
}

// ServiceEndpoint represents internal DNS information for a Kubernetes service
type ServiceEndpoint struct {
	KubernetesName string            `json:"kubernetes_name"`
	DNS            string            `json:"dns"`
	Ports          []schema.PortSpec `json:"ports" nullable:"false"`
	TeamID         uuid.UUID         `json:"team_id"`
	ProjectID      uuid.UUID         `json:"project_id"`
	EnvironmentID  uuid.UUID         `json:"environment_id"`
	ServiceID      uuid.UUID         `json:"service_id"`
}

// IngressEndpoint represents external DNS information for a Kubernetes ingress
type IngressEndpoint struct {
	KubernetesName string    `json:"kubernetes_name"`
	IsIngress      bool      `json:"is_ingress"`
	Host           string    `json:"host"`
	Path           string    `json:"path"`
	Port           *int32    `json:"port"`
	TlsStatus      TlsStatus `json:"tls_issued"`
	TeamID         uuid.UUID `json:"team_id"`
	ProjectID      uuid.UUID `json:"project_id"`
	EnvironmentID  uuid.UUID `json:"environment_id"`
	ServiceID      uuid.UUID `json:"service_id"`
}

type ExtendedHostSpec struct {
	v1.HostSpec
	DnsConfigured bool      `json:"dns_configured"`
	Cloudflare    bool      `json:"cloudflare"`
	TlsStatus     TlsStatus `json:"tls_issued"`
}

type TlsStatus string

const (
	TlsStatusPending      TlsStatus = "pending"
	TlsStatusAttempting   TlsStatus = "attempting"
	TlsStatusIssued       TlsStatus = "issued"
	TlsStatusNotAvailable TlsStatus = "not_available"
)

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u TlsStatus) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["TlsStatus"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "TlsStatus")
		schemaRef.Title = "TlsStatus"
		schemaRef.Enum = append(schemaRef.Enum,
			[]any{
				string(TlsStatusPending),
				string(TlsStatusAttempting),
				string(TlsStatusIssued),
				string(TlsStatusNotAvailable),
			}...,
		)
		r.Map()["TlsStatus"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/TlsStatus"}
}
