package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/unbindapp/unbind-api/ent/schema/mixin"
)

// Team holds the schema definition for the Team entity.
type Team struct {
	ent.Schema
}

// Mixin of the Team.
func (Team) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the Team.
func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
		field.String("display_name").Comment("Human-readable name"),
		field.String("namespace").Unique().Comment("Kubernetes namespace tied to this team"),
		field.String("kubernetes_secret").Comment("Kubernetes secret for this team"),
		field.String("description").Optional().Nillable(),
	}
}

// Edges of the Team.
func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("projects", Project.Type),
		edge.From("members", User.Type).Ref("teams"),
	}
}

// Annotations of the Team
func (Team) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "teams",
		},
	}
}
