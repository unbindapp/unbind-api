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
	v1 "github.com/unbindapp/unbind-operator/api/v1"
)

// Status enum
type DeploymentStatus string

const (
	DeploymentStatusQueued    DeploymentStatus = "queued"
	DeploymentStatusBuilding  DeploymentStatus = "building"
	DeploymentStatusSucceeded DeploymentStatus = "succeeded"
	DeploymentStatusCancelled DeploymentStatus = "cancelled"
	DeploymentStatusFailed    DeploymentStatus = "failed"
)

var allDeploymentStatuses = []DeploymentStatus{
	DeploymentStatusQueued,
	DeploymentStatusBuilding,
	DeploymentStatusSucceeded,
	DeploymentStatusCancelled,
	DeploymentStatusFailed,
}

// Values provides list valid values for Enum.
func (DeploymentStatus) Values() (kinds []string) {
	for _, s := range allDeploymentStatuses {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u DeploymentStatus) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["DeploymentStatus"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "DeploymentStatus")
		schemaRef.Title = "DeploymentStatus"
		for _, v := range allDeploymentStatuses {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["DeploymentStatus"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/DeploymentStatus"}
}

// Source enum
type DeploymentSource string

const (
	DeploymentSourceManual DeploymentSource = "manual"
	DeploymentSourceGit    DeploymentSource = "git"
)

var allDeploymentSources = []DeploymentSource{
	DeploymentSourceManual,
	DeploymentSourceGit,
}

// Values provides list valid values for Enum.
func (DeploymentSource) Values() (kinds []string) {
	for _, s := range allDeploymentSources {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u DeploymentSource) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["DeploymentSource"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "DeploymentSource")
		schemaRef.Title = "DeploymentSource"
		for _, v := range allDeploymentSources {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["DeploymentSource"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/DeploymentSource"}
}

// Type to keep track of git committer
type GitCommitter struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// Deployment holds the schema definition for the Deployment entity.
type Deployment struct {
	ent.Schema
}

// Mixin of the Deployment.
func (Deployment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the Deployment.
func (Deployment) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("service_id", uuid.UUID{}),
		field.Enum("status").GoType(DeploymentStatus("")),
		// ! TODO - remove default
		field.Enum("source").GoType(DeploymentSource("")).Default(string(DeploymentSourceManual)),
		field.String("error").
			Optional(),
		field.String("commit_sha").
			Optional().
			Nillable(),
		field.String("commit_message").
			Optional().
			Nillable(),
		field.JSON("commit_author", &GitCommitter{}).
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
		field.String("image").
			Optional().
			Nillable().
			Comment("Reference to the image used for the deployment"),
		field.JSON("resource_definition", &v1.ServiceSpec{}).
			Optional().
			Comment("The Kubernetes resource definition for the deployment"),
	}
}

// Edges of the Deployment.
func (Deployment) Edges() []ent.Edge {
	return []ent.Edge{
		// M2O edge to keep track of the service
		edge.From("service", Service.Type).Ref("deployments").Field("service_id").Unique().Required(),
	}
}

// Annotations of the Deployment
func (Deployment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "deployments",
		},
	}
}
