package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/unbindapp/unbind-api/ent/schema/mixin"
)

// S3 holds the schema definition for the S3 entity.
type S3 struct {
	ent.Schema
}

// Mixin of the S3.
func (S3) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the S3.
func (S3) Fields() []ent.Field {
	return []ent.Field{
		field.String("endpoint"),
		field.String("region"),
		field.Bool("force_path_style").Default(true),
		field.String("kubernetes_secret"),
	}
}

// Edges of the S3.
func (S3) Edges() []ent.Edge {
	return nil
}

// Annotations of the S3
func (S3) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "s3_sources",
		},
	}
}
