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
		// ! TODO - remove default after migration
		field.String("kubernetes_secret").Default("").Comment("Kubernetes secret for this team"),
		field.String("description").Optional(),
	}
}

// Edges of the Team.
func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("projects", Project.Type),
		edge.From("members", User.Type).Ref("teams"),
		edge.To("groups", Group.Type),
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
