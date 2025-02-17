package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/unbindapp/unbind-api/ent/schema/mixin"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").Unique(),
		field.String("username"),
		field.String("external_id").Unique().Comment(
			"Dex subject ID",
		),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
