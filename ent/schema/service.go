package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema/mixin"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
)

// Service holds the schema definition for the Service entity.
type Service struct {
	ent.Schema
}

// Mixin of the Service.
func (Service) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the Service.
func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("display_name"),
		field.String("description").Optional(),
		field.Enum("type").Values("git", "dockerfile").Comment("Type of service"),
		field.Enum("builder").Values("railpack", "docker"),
		field.Enum("provider").GoType(enum.Provider("")).Optional().Nillable().Comment("Provider (e.g. Go, Python, Node, Deno)"),
		field.Enum("framework").GoType(enum.Framework("")).Optional().Nillable().Comment("Framework of service - corresponds mostly to railpack results - e.g. Django, Next, Express, Gin"),
		field.UUID("environment_id", uuid.UUID{}),
		field.Int64("github_installation_id").Optional().Nillable().Comment("Optional reference to GitHub installation"),
		field.String("git_repository").Optional().Nillable().Comment("GitHub repository name"),
		field.String("kubernetes_secret").Comment("Kubernetes secret for this service"),
	}
}

// Edges of the Service.
func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		// M2O edge to keep track of the environment
		edge.From("environment", Environment.Type).Ref("services").Field("environment_id").Unique().Required(),
		// M2O edge to keep track of the GitHub installation
		edge.From("github_installation", GithubInstallation.Type).Ref("services").Field("github_installation_id").Unique(),
		// O2O edge to keep track of the service configuration
		edge.To("service_config", ServiceConfig.Type).Unique(),
		// O2M edge to keep track of deployments
		edge.To("deployments", Deployment.Type),
	}
}

// Annotations of the Service
func (Service) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "services",
		},
	}
}
