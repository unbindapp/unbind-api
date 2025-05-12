package schema

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

// * Custom kubernetes-like types
type PortSpec struct {
	InputTemplateID *int `json:"input_template_id,omitempty" required:"false" doc:"For template port inputs"`
	// Will create a node port (public) service
	IsNodePort bool   `json:"is_nodeport" required:"false"`
	NodePort   *int32 `json:"node_port,omitempty" required:"false"`
	// Port is the container port to expose
	Port     int32     `json:"port"`
	Protocol *Protocol `json:"protocol,omitempty" required:"false"`
}

func (self *PortSpec) AsV1PortSpec() v1.PortSpec {
	var protocol *corev1.Protocol
	if self.Protocol != nil {
		protocol = utils.ToPtr(corev1.Protocol(*self.Protocol))
	} else {
		protocol = utils.ToPtr(corev1.ProtocolTCP)
	}
	return v1.PortSpec{
		NodePort: self.NodePort,
		Port:     self.Port,
		Protocol: protocol,
	}
}

type Protocol string

const (
	ProtocolTCP  Protocol = "TCP"
	ProtocolUDP  Protocol = "UDP"
	ProtocolSCTP Protocol = "SCTP"
)

// Values provides list valid values for Enum.
func (s Protocol) Values() (kinds []string) {
	kinds = append(kinds, []string{
		string(ProtocolTCP),
		string(ProtocolUDP),
		string(ProtocolSCTP),
	}...)
	return
}

func AsV1PortSpecs(ports []PortSpec) []v1.PortSpec {
	v1Ports := make([]v1.PortSpec, len(ports))
	for i, port := range ports {
		v1Ports[i] = port.AsV1PortSpec()
	}
	return v1Ports
}

// * Kubernetes Security context
type Capability string

// Adds and removes POSIX capabilities from running containers.
type Capabilities struct {
	Add  []Capability `json:"add,omitempty" protobuf:"bytes,1,rep,name=add,casttype=Capability"`
	Drop []Capability `json:"drop,omitempty" protobuf:"bytes,2,rep,name=drop,casttype=Capability"`
}

type SecurityContext struct {
	Capabilities *Capabilities `json:"capabilities,omitempty" protobuf:"bytes,1,opt,name=capabilities"`
	Privileged   *bool         `json:"privileged,omitempty" protobuf:"varint,2,opt,name=privileged"`
}

func (self *SecurityContext) AsV1SecurityContext() *corev1.SecurityContext {
	if self == nil {
		return nil
	}
	secCtx := &corev1.SecurityContext{}
	if self.Privileged != nil {
		secCtx.Privileged = self.Privileged
	}
	if self.Capabilities != nil {
		secCtx.Capabilities = &corev1.Capabilities{}
		if self.Capabilities.Add != nil {
			secCtx.Capabilities.Add = make([]corev1.Capability, len(self.Capabilities.Add))
			for i, cap := range self.Capabilities.Add {
				secCtx.Capabilities.Add[i] = corev1.Capability(cap)
			}
		}
		if self.Capabilities.Drop != nil {
			secCtx.Capabilities.Drop = make([]corev1.Capability, len(self.Capabilities.Drop))
			for i, cap := range self.Capabilities.Drop {
				secCtx.Capabilities.Drop[i] = corev1.Capability(cap)
			}
		}
	}
	return secCtx
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u Protocol) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["Protocol"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "Protocol")
		schemaRef.Title = "Protocol"
		schemaRef.Enum = append(schemaRef.Enum, []any{
			string(ProtocolTCP),
			string(ProtocolUDP),
			string(ProtocolSCTP),
		}...)
		r.Map()["Protocol"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/Protocol"}
}

type DatabaseConfig struct {
	Version             string `json:"version,omitempty" required:"false" description:"Version of the database"`
	StorageSize         string `json:"storage,omitempty" required:"false" description:"Storage size for the database"`
	DefaultDatabaseName string `json:"defaultDatabaseName,omitempty" required:"false" description:"Default database name"`
	InitDB              string `json:"initdb,omitempty" required:"false" description:"SQL commands to run to initialize the database"`
}

func (self *DatabaseConfig) AsMap() map[string]interface{} {
	ret := make(map[string]interface{})

	if self.Version != "" {
		ret["version"] = self.Version
	}
	if self.StorageSize != "" {
		ret["storage"] = self.StorageSize
	}
	if self.DefaultDatabaseName != "" {
		ret["defaultDatabaseName"] = self.DefaultDatabaseName
	}
	if self.InitDB != "" {
		ret["initdb"] = self.InitDB
	}
	return ret
}

//* Enums

// Builder enum
type ServiceBuilder string

const (
	ServiceBuilderRailpack ServiceBuilder = "railpack"
	ServiceBuilderDocker   ServiceBuilder = "docker"
	ServiceBuilderDatabase ServiceBuilder = "database"
)

var allServiceBuilders = []ServiceBuilder{
	ServiceBuilderRailpack,
	ServiceBuilderDocker,
	ServiceBuilderDatabase,
}

// Values provides list valid values for Enum.
func (s ServiceBuilder) Values() (kinds []string) {
	for _, s := range allServiceBuilders {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u ServiceBuilder) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["ServiceBuilder"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "ServiceBuilder")
		schemaRef.Title = "ServiceBuilder"
		for _, v := range allServiceBuilders {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["ServiceBuilder"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/ServiceBuilder"}
}
