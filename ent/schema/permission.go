package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/unbindapp/unbind-api/ent/schema/mixin"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

// Mixin of the Permission.
func (Permission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("action").Values("read", "create", "update", "delete", "manage", "admin", "edit", "view"),
		field.Enum("resource_type").Values("team", "project", "group", "environment", "permission", "user", "system", "service").Comment("Type of resource: 'teams', 'projects', 'k8s', etc."),
		field.String("resource_id").NotEmpty().Comment("Specific resource ID or '*' for all resources of this type"),
		field.String("scope").Optional().Comment("For additional filtering (e.g., k8s namespaces, specific fields)"),
		field.JSON("labels", map[string]string{}).Optional().Comment("Resource labels for K8s label selectors"),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("groups", Group.Type).Ref("permissions"),
	}
}

// Annotations of the Permission
func (Permission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "permissions",
		},
	}
}
