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
	"github.com/unbindapp/unbind-api/pkg/templates"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
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
		// For builds from git using Dockerfile
		field.String("dockerfile_path").Optional().Nillable().Comment("Path to Dockerfile if using docker builder"),
		field.String("dockerfile_context").Optional().Nillable().Comment("Path to Dockerfile context if using docker builder"),
		// Provider and framework directly from railpack
		field.Enum("provider").GoType(enum.Provider("")).Optional().Nillable().Comment("Provider (e.g. Go, Python, Node, Deno)"),
		field.Enum("framework").GoType(enum.Framework("")).Optional().Nillable().Comment("Framework of service - corresponds mostly to railpack results - e.g. Django, Next, Express, Gin"),
		// Branch to build from (git)
		field.String("git_branch").Optional().Nillable().Comment("Branch to build from"),
		// Generic CRD configuration
		field.JSON("hosts", []v1.HostSpec{}).Optional().Comment("External domains and paths for the service"),
		field.JSON("ports", []v1.PortSpec{}).Optional().Comment("Container ports to expose"),
		field.Int32("replicas").Default(2).Comment("Number of replicas for the service"),
		field.Bool("auto_deploy").Default(false).Comment("Whether to automatically deploy on git push"),
		field.String("run_command").Optional().Nillable().Comment("Custom run command"),
		field.Bool("public").Default(false).Comment("Whether the service is publicly accessible, creates an ingress resource"),
		field.String("image").Optional().Comment("Custom Docker image if not building from git"), // Only applies to type=docker-image
		// Template configuration
		field.Enum("template_category").GoType(templates.TemplateCategoryName("")).Optional().Nillable().Comment("Template category to use for the service"),
		field.String("template").Optional().Nillable().Comment("Template to use for the service"),
		field.String("template_release_version").Optional().Nillable().Comment("Version of the templates release"),
		field.String("template_version").Optional().Nillable().Comment("Version of the template to use"),
		field.JSON("template_config", map[string]interface{}{}).Optional().Comment("Template configuration for the service"),
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
	ServiceTypeGithub      ServiceType = "github"
	ServiceTypeDockerimage ServiceType = "docker-image"
	ServiceTypeTemplate    ServiceType = "template"
)

var allServiceTypes = []ServiceType{
	ServiceTypeGithub,
	ServiceTypeDockerimage,
	ServiceTypeTemplate,
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
	ServiceBuilderRailpack ServiceBuilder = "railpack"
	ServiceBuilderDocker   ServiceBuilder = "docker"
	ServiceBuilderTemplate ServiceBuilder = "template"
)

var allServiceBuilders = []ServiceBuilder{
	ServiceBuilderRailpack,
	ServiceBuilderDocker,
	ServiceBuilderTemplate,
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
