package schema

import (
	"reflect"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema/mixin"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

// Custom types
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

// * Security context
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

// ServiceConfig holds environment-specific configuration for a service
type ServiceConfig struct {
	ent.Schema
}

// Mixin of the ServiceConfig.
func (ServiceConfig) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the ServiceConfig.
func (ServiceConfig) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("service_id", uuid.UUID{}),
		field.Enum("builder").GoType(ServiceBuilder("")),
		field.String("icon").Comment("Icon metadata, unique of framework, provider, database"),
		// For builds from git using Dockerfile
		field.String("dockerfile_path").Optional().Nillable().Comment("Path to Dockerfile if using docker builder"),
		field.String("dockerfile_context").Optional().Nillable().Comment("Path to Dockerfile context if using docker builder"),
		// Provider and framework directly from railpack
		field.Enum("railpack_provider").GoType(enum.Provider("")).Optional().Nillable().Comment("Provider (e.g. Go, Python, Node, Deno)"),
		field.Enum("railpack_framework").GoType(enum.Framework("")).Optional().Nillable().Comment("Framework of service - corresponds mostly to railpack results - e.g. Django, Next, Express, Gin"),
		// Branch to build from (git)
		field.String("git_branch").Optional().Nillable().Comment("Branch to build from"),
		field.String("git_tag").Optional().Nillable().Comment("Tag to build from, supports glob patterns"),
		// Generic CRD configuration
		field.JSON("hosts", []v1.HostSpec{}).Optional().Comment("External domains and paths for the service"),
		field.JSON("ports", []PortSpec{}).Optional().Comment("Container ports to expose"),
		field.Int32("replicas").Default(1).Comment("Number of replicas for the service"),
		field.Bool("auto_deploy").Default(false).Comment("Whether to automatically deploy on git push"),
		field.String("install_command").Optional().Nillable().Comment("Custom install command (railpack only)"),
		field.String("build_command").Optional().Nillable().Comment("Custom build command (railpack only)"),
		field.String("run_command").Optional().Nillable().Comment("Custom run command"),
		field.Bool("is_public").Default(false).Comment("Whether the service is publicly accessible, creates an ingress resource"),
		field.String("image").Optional().Comment("Custom Docker image if not building from git"), // Only applies to type=docker-image
		// Database
		field.String("definition_version").Optional().Nillable().Comment("Version of the database custom resource definition"),
		field.JSON("database_config", &DatabaseConfig{}).Optional().Comment("Database configuration for the service"),
		field.UUID("s3_backup_endpoint_id", uuid.UUID{}).Optional().Nillable().Comment("S3 endpoint backup to"),
		field.String("s3_backup_bucket").Optional().Nillable().Comment("S3 bucket to backup to"),
		field.String("backup_schedule").Default("5 5 * * *").Comment("Cron expression for the backup schedule"),
		field.Int("backup_retention_count").Default(3).Comment("Number of base backups to retain"),
		// Volume
		field.String("volume_name").Optional().Nillable().Comment("Volume name to use for the service"),
		field.String("volume_mount_path").Optional().Nillable().Comment("Volume mount path for the service"),
		// Security context
		field.JSON("security_context", &SecurityContext{}).Optional().Comment("Security context for the service containers."),
	}
}

// Edges of the ServiceConfig.
func (ServiceConfig) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("service", Service.Type).Ref("service_config").Field("service_id").Unique().Required(),
		// O2M to backup sources
		edge.From("s3_backup_endpoint", S3.Type).Ref("service_backup_endpoint").Field("s3_backup_endpoint_id").Unique(),
	}
}

// Annotations of the ServiceConfig
func (ServiceConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "service_configs",
		},
	}
}

// Enums

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
