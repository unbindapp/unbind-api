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

// Environment holds the schema definition for the Environment entity.
type Environment struct {
	ent.Schema
}

// Mixin of the Environment.
func (Environment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the Environment.
func (Environment) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("display_name"),
		field.String("description").Optional().Nillable(),
		field.Bool("active").Default(true),
		field.UUID("project_id", uuid.UUID{}),
		field.String("kubernetes_secret").Comment("Kubernetes secret for this environment"),
	}
}

// Edges of the Environment.
func (Environment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).Ref("environments").Field("project_id").Unique().Required(),
		edge.To("services", Service.Type).Annotations(
			entsql.Annotation{OnDelete: entsql.Cascade},
		),
		// O2O
		edge.To("project_default", Project.Type),
	}
}

// Annotations of the Environment
func (Environment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "environments",
		},
	}
}
