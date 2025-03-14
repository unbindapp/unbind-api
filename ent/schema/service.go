package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema/mixin"
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
		field.Enum("type").Values("database", "api", "web", "custom").Comment("Type of service"),
		field.Enum("subtype").Values("react", "go", "node", "next", "other").Comment("Type of service"),
		field.UUID("project_id", uuid.UUID{}),
		field.Int64("github_installation_id").Optional().Nillable().Comment("Optional reference to GitHub installation"),
		field.String("git_repository").Optional().Nillable().Comment("GitHub repository name"),
		field.String("git_branch").Optional().Nillable().Default("main").Comment("Branch to build from"),
	}
}

// Edges of the Service.
func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		// M2O edge to keep track of the project
		edge.From("project", Project.Type).Ref("services").Field("project_id").Unique().Required(),
		// M2O edge to keep track of the GitHub installation
		edge.From("github_installation", GithubInstallation.Type).Ref("services").Field("github_installation_id").Unique(),
		// O2O edge to keep track of the service configuration
		edge.To("service_configs", ServiceConfig.Type).Unique(),
		// ! TODO - Add edge to keep track of build history
		// edge.To("deployments", Deployment.Type),
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
