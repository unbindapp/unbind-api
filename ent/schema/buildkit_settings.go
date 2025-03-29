package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/unbindapp/unbind-api/ent/schema/mixin"
)

// BuildkitSettings holds the schema definition for the BuildkitSettings entity.
type BuildkitSettings struct {
	ent.Schema
}

// Mixin of the BuildkitSettings.
func (BuildkitSettings) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the BuildkitSettings.
func (BuildkitSettings) Fields() []ent.Field {
	return []ent.Field{
		field.Int("max_parallelism").Default(2).Validate(func(i int) error {
			if i < 1 {
				return fmt.Errorf("max_parallelism must be greater than 0")
			}
			return nil
		}).
			Comment("Maximum number of parallel build steps"),
		field.Int("replicas").Default(1).Validate(func(i int) error {
			if i < 1 {
				return fmt.Errorf("replicas must be greater than 0")
			}
			return nil
		}).
			Comment("Number of replicas for the buildkitd deployment"),
	}
}

// Edges of the BuildkitSettings.
func (BuildkitSettings) Edges() []ent.Edge {
	return nil
}

// Annotations of the BuildkitSettings
func (BuildkitSettings) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "buildkit_settings",
		},
	}
}
