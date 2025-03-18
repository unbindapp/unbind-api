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

// BuildJob holds the schema definition for the BuildJob entity.
type BuildJob struct {
	ent.Schema
}

// Mixin of the BuildJob.
func (BuildJob) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the BuildJob.
func (BuildJob) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("service_id", uuid.UUID{}),
		field.Enum("status").Values("queued", "running", "completed", "cancelled", "failed"),
		field.String("error").
			Optional(),
		field.Time("started_at").
			Optional().
			Nillable(),
		field.Time("completed_at").
			Optional().
			Nillable(),
		field.String("kubernetes_job_name").
			Optional().
			Comment("The name of the kubernetes job"),
		field.String("kubernetes_job_status").
			Optional().
			Comment("The status of the kubernetes job"),
		field.Int("attempts").
			Default(0),
	}
}

// Edges of the BuildJob.
func (BuildJob) Edges() []ent.Edge {
	return []ent.Edge{
		// M2O edge to keep track of the service
		edge.From("service", Service.Type).Ref("build_jobs").Field("service_id").Unique().Required(),
	}
}

// Annotations of the BuildJob
func (BuildJob) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "build_jobs",
		},
	}
}
