package schema

import (
	"reflect"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema/mixin"
)

// Status enum
type BuildJobStatus string

const (
	BuildJobStatusQueued    BuildJobStatus = "queued"
	BuildJobStatusRunning   BuildJobStatus = "running"
	BuildJobStatusCompleted BuildJobStatus = "completed"
	BuildJobStatusCancelled BuildJobStatus = "cancelled"
	BuildJobStatusFailed    BuildJobStatus = "failed"
)

var allBuildJobStatuses = []BuildJobStatus{
	BuildJobStatusQueued,
	BuildJobStatusRunning,
	BuildJobStatusCompleted,
	BuildJobStatusCancelled,
	BuildJobStatusFailed,
}

// Values provides list valid values for Enum.
func (BuildJobStatus) Values() (kinds []string) {
	for _, s := range allBuildJobStatuses {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u BuildJobStatus) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["BuildJobStatus"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "BuildJobStatus")
		schemaRef.Title = "BuildJobStatus"
		for _, v := range allBuildJobStatuses {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["BuildJobStatus"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/BuildJobStatus"}
}

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
		field.Enum("status").GoType(BuildJobStatus("")),
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
