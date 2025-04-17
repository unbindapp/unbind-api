package models

import (
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
	Name          string            `json:"name"`
	DNS           string            `json:"dns"`
	Ports         []schema.PortSpec `json:"ports" nullable:"false"`
	TeamID        uuid.UUID         `json:"team_id"`
	ProjectID     uuid.UUID         `json:"project_id"`
	EnvironmentID uuid.UUID         `json:"environment_id"`
	ServiceID     uuid.UUID         `json:"service_id"`
}

// IngressEndpoint represents external DNS information for a Kubernetes ingress
type IngressEndpoint struct {
	Name          string             `json:"name"`
	Hosts         []ExtendedHostSpec `json:"hosts" nullable:"false"`
	TeamID        uuid.UUID          `json:"team_id"`
	ProjectID     uuid.UUID          `json:"project_id"`
	EnvironmentID uuid.UUID          `json:"environment_id"`
	ServiceID     uuid.UUID          `json:"service_id"`
}

type ExtendedHostSpec struct {
	v1.HostSpec
	Issued bool `json:"issued"`
}
