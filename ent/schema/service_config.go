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
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
)

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
		field.Enum("type").GoType(ServiceType("")).Comment("Type of service"),
		field.Enum("builder").GoType(ServiceBuilder("")),
		field.Enum("provider").GoType(enum.Provider("")).Optional().Nillable().Comment("Provider (e.g. Go, Python, Node, Deno)"),
		field.Enum("framework").GoType(enum.Framework("")).Optional().Nillable().Comment("Framework of service - corresponds mostly to railpack results - e.g. Django, Next, Express, Gin"),
		field.String("git_branch").Optional().Nillable().Comment("Branch to build from"),
		field.String("host").Optional().Nillable().Comment("External domain for the service, e.g., unbind.app"),
		field.Int("port").Optional().Nillable().Comment("Main container port"),
		field.Int32("replicas").Default(2).Comment("Number of replicas for the service"),
		field.Bool("auto_deploy").Default(false).Comment("Whether to automatically deploy on git push"),
		field.String("run_command").Optional().Nillable().Comment("Custom run command"),
		field.Bool("public").Default(false).Comment("Whether the service is publicly accessible, creates an ingress resource"),
		field.String("image").Optional().Comment("Custom Docker image if not building from git"),
	}
}

// Edges of the ServiceConfig.
func (ServiceConfig) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("service", Service.Type).Ref("service_config").Field("service_id").Unique().Required(),
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
// Type enum
type ServiceType string

const (
	ServiceTypeGit        ServiceType = "git"
	ServiceTypeDockerfile ServiceType = "dockerfile"
)

var allServiceTypes = []ServiceType{
	ServiceTypeGit,
	ServiceTypeDockerfile,
}

// Values provides list valid values for Enum.
func (s ServiceType) Values() (kinds []string) {
	for _, s := range allServiceTypes {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u ServiceType) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["ServiceType"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "ServiceType")
		schemaRef.Title = "ServiceType"
		for _, v := range allServiceTypes {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["ServiceType"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/ServiceType"}
}

// Builder enum
type ServiceBuilder string

const (
	ServiceBuilderNixpacks ServiceBuilder = "nixpacks"
	ServiceBuilderRailpack ServiceBuilder = "railpack"
	ServiceBuilderDocker   ServiceBuilder = "docker"
)

var allServiceBuilders = []ServiceBuilder{
	ServiceBuilderNixpacks,
	ServiceBuilderRailpack,
	ServiceBuilderDocker,
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
