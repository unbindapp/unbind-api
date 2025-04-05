// Code generated by ent, DO NOT EDIT.

package serviceconfig

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/predicate"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldUpdatedAt, v))
}

// ServiceID applies equality check predicate on the "service_id" field. It's identical to ServiceIDEQ.
func ServiceID(v uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldServiceID, v))
}

// DockerfilePath applies equality check predicate on the "dockerfile_path" field. It's identical to DockerfilePathEQ.
func DockerfilePath(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldDockerfilePath, v))
}

// DockerfileContext applies equality check predicate on the "dockerfile_context" field. It's identical to DockerfileContextEQ.
func DockerfileContext(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldDockerfileContext, v))
}

// GitBranch applies equality check predicate on the "git_branch" field. It's identical to GitBranchEQ.
func GitBranch(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldGitBranch, v))
}

// Replicas applies equality check predicate on the "replicas" field. It's identical to ReplicasEQ.
func Replicas(v int32) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldReplicas, v))
}

// AutoDeploy applies equality check predicate on the "auto_deploy" field. It's identical to AutoDeployEQ.
func AutoDeploy(v bool) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldAutoDeploy, v))
}

// RunCommand applies equality check predicate on the "run_command" field. It's identical to RunCommandEQ.
func RunCommand(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldRunCommand, v))
}

// Public applies equality check predicate on the "public" field. It's identical to PublicEQ.
func Public(v bool) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldPublic, v))
}

// Image applies equality check predicate on the "image" field. It's identical to ImageEQ.
func Image(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldImage, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLTE(FieldUpdatedAt, v))
}

// ServiceIDEQ applies the EQ predicate on the "service_id" field.
func ServiceIDEQ(v uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldServiceID, v))
}

// ServiceIDNEQ applies the NEQ predicate on the "service_id" field.
func ServiceIDNEQ(v uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldServiceID, v))
}

// ServiceIDIn applies the In predicate on the "service_id" field.
func ServiceIDIn(vs ...uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIn(FieldServiceID, vs...))
}

// ServiceIDNotIn applies the NotIn predicate on the "service_id" field.
func ServiceIDNotIn(vs ...uuid.UUID) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotIn(FieldServiceID, vs...))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v schema.ServiceType) predicate.ServiceConfig {
	vc := v
	return predicate.ServiceConfig(sql.FieldEQ(FieldType, vc))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v schema.ServiceType) predicate.ServiceConfig {
	vc := v
	return predicate.ServiceConfig(sql.FieldNEQ(FieldType, vc))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...schema.ServiceType) predicate.ServiceConfig {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceConfig(sql.FieldIn(FieldType, v...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...schema.ServiceType) predicate.ServiceConfig {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceConfig(sql.FieldNotIn(FieldType, v...))
}

// BuilderEQ applies the EQ predicate on the "builder" field.
func BuilderEQ(v schema.ServiceBuilder) predicate.ServiceConfig {
	vc := v
	return predicate.ServiceConfig(sql.FieldEQ(FieldBuilder, vc))
}

// BuilderNEQ applies the NEQ predicate on the "builder" field.
func BuilderNEQ(v schema.ServiceBuilder) predicate.ServiceConfig {
	vc := v
	return predicate.ServiceConfig(sql.FieldNEQ(FieldBuilder, vc))
}

// BuilderIn applies the In predicate on the "builder" field.
func BuilderIn(vs ...schema.ServiceBuilder) predicate.ServiceConfig {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceConfig(sql.FieldIn(FieldBuilder, v...))
}

// BuilderNotIn applies the NotIn predicate on the "builder" field.
func BuilderNotIn(vs ...schema.ServiceBuilder) predicate.ServiceConfig {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceConfig(sql.FieldNotIn(FieldBuilder, v...))
}

// DockerfilePathEQ applies the EQ predicate on the "dockerfile_path" field.
func DockerfilePathEQ(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldDockerfilePath, v))
}

// DockerfilePathNEQ applies the NEQ predicate on the "dockerfile_path" field.
func DockerfilePathNEQ(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldDockerfilePath, v))
}

// DockerfilePathIn applies the In predicate on the "dockerfile_path" field.
func DockerfilePathIn(vs ...string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIn(FieldDockerfilePath, vs...))
}

// DockerfilePathNotIn applies the NotIn predicate on the "dockerfile_path" field.
func DockerfilePathNotIn(vs ...string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotIn(FieldDockerfilePath, vs...))
}

// DockerfilePathGT applies the GT predicate on the "dockerfile_path" field.
func DockerfilePathGT(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGT(FieldDockerfilePath, v))
}

// DockerfilePathGTE applies the GTE predicate on the "dockerfile_path" field.
func DockerfilePathGTE(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGTE(FieldDockerfilePath, v))
}

// DockerfilePathLT applies the LT predicate on the "dockerfile_path" field.
func DockerfilePathLT(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLT(FieldDockerfilePath, v))
}

// DockerfilePathLTE applies the LTE predicate on the "dockerfile_path" field.
func DockerfilePathLTE(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLTE(FieldDockerfilePath, v))
}

// DockerfilePathContains applies the Contains predicate on the "dockerfile_path" field.
func DockerfilePathContains(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldContains(FieldDockerfilePath, v))
}

// DockerfilePathHasPrefix applies the HasPrefix predicate on the "dockerfile_path" field.
func DockerfilePathHasPrefix(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldHasPrefix(FieldDockerfilePath, v))
}

// DockerfilePathHasSuffix applies the HasSuffix predicate on the "dockerfile_path" field.
func DockerfilePathHasSuffix(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldHasSuffix(FieldDockerfilePath, v))
}

// DockerfilePathIsNil applies the IsNil predicate on the "dockerfile_path" field.
func DockerfilePathIsNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIsNull(FieldDockerfilePath))
}

// DockerfilePathNotNil applies the NotNil predicate on the "dockerfile_path" field.
func DockerfilePathNotNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotNull(FieldDockerfilePath))
}

// DockerfilePathEqualFold applies the EqualFold predicate on the "dockerfile_path" field.
func DockerfilePathEqualFold(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEqualFold(FieldDockerfilePath, v))
}

// DockerfilePathContainsFold applies the ContainsFold predicate on the "dockerfile_path" field.
func DockerfilePathContainsFold(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldContainsFold(FieldDockerfilePath, v))
}

// DockerfileContextEQ applies the EQ predicate on the "dockerfile_context" field.
func DockerfileContextEQ(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldDockerfileContext, v))
}

// DockerfileContextNEQ applies the NEQ predicate on the "dockerfile_context" field.
func DockerfileContextNEQ(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldDockerfileContext, v))
}

// DockerfileContextIn applies the In predicate on the "dockerfile_context" field.
func DockerfileContextIn(vs ...string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIn(FieldDockerfileContext, vs...))
}

// DockerfileContextNotIn applies the NotIn predicate on the "dockerfile_context" field.
func DockerfileContextNotIn(vs ...string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotIn(FieldDockerfileContext, vs...))
}

// DockerfileContextGT applies the GT predicate on the "dockerfile_context" field.
func DockerfileContextGT(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGT(FieldDockerfileContext, v))
}

// DockerfileContextGTE applies the GTE predicate on the "dockerfile_context" field.
func DockerfileContextGTE(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGTE(FieldDockerfileContext, v))
}

// DockerfileContextLT applies the LT predicate on the "dockerfile_context" field.
func DockerfileContextLT(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLT(FieldDockerfileContext, v))
}

// DockerfileContextLTE applies the LTE predicate on the "dockerfile_context" field.
func DockerfileContextLTE(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLTE(FieldDockerfileContext, v))
}

// DockerfileContextContains applies the Contains predicate on the "dockerfile_context" field.
func DockerfileContextContains(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldContains(FieldDockerfileContext, v))
}

// DockerfileContextHasPrefix applies the HasPrefix predicate on the "dockerfile_context" field.
func DockerfileContextHasPrefix(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldHasPrefix(FieldDockerfileContext, v))
}

// DockerfileContextHasSuffix applies the HasSuffix predicate on the "dockerfile_context" field.
func DockerfileContextHasSuffix(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldHasSuffix(FieldDockerfileContext, v))
}

// DockerfileContextIsNil applies the IsNil predicate on the "dockerfile_context" field.
func DockerfileContextIsNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIsNull(FieldDockerfileContext))
}

// DockerfileContextNotNil applies the NotNil predicate on the "dockerfile_context" field.
func DockerfileContextNotNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotNull(FieldDockerfileContext))
}

// DockerfileContextEqualFold applies the EqualFold predicate on the "dockerfile_context" field.
func DockerfileContextEqualFold(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEqualFold(FieldDockerfileContext, v))
}

// DockerfileContextContainsFold applies the ContainsFold predicate on the "dockerfile_context" field.
func DockerfileContextContainsFold(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldContainsFold(FieldDockerfileContext, v))
}

// ProviderEQ applies the EQ predicate on the "provider" field.
func ProviderEQ(v enum.Provider) predicate.ServiceConfig {
	vc := v
	return predicate.ServiceConfig(sql.FieldEQ(FieldProvider, vc))
}

// ProviderNEQ applies the NEQ predicate on the "provider" field.
func ProviderNEQ(v enum.Provider) predicate.ServiceConfig {
	vc := v
	return predicate.ServiceConfig(sql.FieldNEQ(FieldProvider, vc))
}

// ProviderIn applies the In predicate on the "provider" field.
func ProviderIn(vs ...enum.Provider) predicate.ServiceConfig {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceConfig(sql.FieldIn(FieldProvider, v...))
}

// ProviderNotIn applies the NotIn predicate on the "provider" field.
func ProviderNotIn(vs ...enum.Provider) predicate.ServiceConfig {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceConfig(sql.FieldNotIn(FieldProvider, v...))
}

// ProviderIsNil applies the IsNil predicate on the "provider" field.
func ProviderIsNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIsNull(FieldProvider))
}

// ProviderNotNil applies the NotNil predicate on the "provider" field.
func ProviderNotNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotNull(FieldProvider))
}

// FrameworkEQ applies the EQ predicate on the "framework" field.
func FrameworkEQ(v enum.Framework) predicate.ServiceConfig {
	vc := v
	return predicate.ServiceConfig(sql.FieldEQ(FieldFramework, vc))
}

// FrameworkNEQ applies the NEQ predicate on the "framework" field.
func FrameworkNEQ(v enum.Framework) predicate.ServiceConfig {
	vc := v
	return predicate.ServiceConfig(sql.FieldNEQ(FieldFramework, vc))
}

// FrameworkIn applies the In predicate on the "framework" field.
func FrameworkIn(vs ...enum.Framework) predicate.ServiceConfig {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceConfig(sql.FieldIn(FieldFramework, v...))
}

// FrameworkNotIn applies the NotIn predicate on the "framework" field.
func FrameworkNotIn(vs ...enum.Framework) predicate.ServiceConfig {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ServiceConfig(sql.FieldNotIn(FieldFramework, v...))
}

// FrameworkIsNil applies the IsNil predicate on the "framework" field.
func FrameworkIsNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIsNull(FieldFramework))
}

// FrameworkNotNil applies the NotNil predicate on the "framework" field.
func FrameworkNotNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotNull(FieldFramework))
}

// GitBranchEQ applies the EQ predicate on the "git_branch" field.
func GitBranchEQ(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldGitBranch, v))
}

// GitBranchNEQ applies the NEQ predicate on the "git_branch" field.
func GitBranchNEQ(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldGitBranch, v))
}

// GitBranchIn applies the In predicate on the "git_branch" field.
func GitBranchIn(vs ...string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIn(FieldGitBranch, vs...))
}

// GitBranchNotIn applies the NotIn predicate on the "git_branch" field.
func GitBranchNotIn(vs ...string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotIn(FieldGitBranch, vs...))
}

// GitBranchGT applies the GT predicate on the "git_branch" field.
func GitBranchGT(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGT(FieldGitBranch, v))
}

// GitBranchGTE applies the GTE predicate on the "git_branch" field.
func GitBranchGTE(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGTE(FieldGitBranch, v))
}

// GitBranchLT applies the LT predicate on the "git_branch" field.
func GitBranchLT(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLT(FieldGitBranch, v))
}

// GitBranchLTE applies the LTE predicate on the "git_branch" field.
func GitBranchLTE(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLTE(FieldGitBranch, v))
}

// GitBranchContains applies the Contains predicate on the "git_branch" field.
func GitBranchContains(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldContains(FieldGitBranch, v))
}

// GitBranchHasPrefix applies the HasPrefix predicate on the "git_branch" field.
func GitBranchHasPrefix(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldHasPrefix(FieldGitBranch, v))
}

// GitBranchHasSuffix applies the HasSuffix predicate on the "git_branch" field.
func GitBranchHasSuffix(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldHasSuffix(FieldGitBranch, v))
}

// GitBranchIsNil applies the IsNil predicate on the "git_branch" field.
func GitBranchIsNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIsNull(FieldGitBranch))
}

// GitBranchNotNil applies the NotNil predicate on the "git_branch" field.
func GitBranchNotNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotNull(FieldGitBranch))
}

// GitBranchEqualFold applies the EqualFold predicate on the "git_branch" field.
func GitBranchEqualFold(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEqualFold(FieldGitBranch, v))
}

// GitBranchContainsFold applies the ContainsFold predicate on the "git_branch" field.
func GitBranchContainsFold(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldContainsFold(FieldGitBranch, v))
}

// HostsIsNil applies the IsNil predicate on the "hosts" field.
func HostsIsNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIsNull(FieldHosts))
}

// HostsNotNil applies the NotNil predicate on the "hosts" field.
func HostsNotNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotNull(FieldHosts))
}

// PortsIsNil applies the IsNil predicate on the "ports" field.
func PortsIsNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIsNull(FieldPorts))
}

// PortsNotNil applies the NotNil predicate on the "ports" field.
func PortsNotNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotNull(FieldPorts))
}

// ReplicasEQ applies the EQ predicate on the "replicas" field.
func ReplicasEQ(v int32) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldReplicas, v))
}

// ReplicasNEQ applies the NEQ predicate on the "replicas" field.
func ReplicasNEQ(v int32) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldReplicas, v))
}

// ReplicasIn applies the In predicate on the "replicas" field.
func ReplicasIn(vs ...int32) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIn(FieldReplicas, vs...))
}

// ReplicasNotIn applies the NotIn predicate on the "replicas" field.
func ReplicasNotIn(vs ...int32) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotIn(FieldReplicas, vs...))
}

// ReplicasGT applies the GT predicate on the "replicas" field.
func ReplicasGT(v int32) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGT(FieldReplicas, v))
}

// ReplicasGTE applies the GTE predicate on the "replicas" field.
func ReplicasGTE(v int32) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGTE(FieldReplicas, v))
}

// ReplicasLT applies the LT predicate on the "replicas" field.
func ReplicasLT(v int32) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLT(FieldReplicas, v))
}

// ReplicasLTE applies the LTE predicate on the "replicas" field.
func ReplicasLTE(v int32) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLTE(FieldReplicas, v))
}

// AutoDeployEQ applies the EQ predicate on the "auto_deploy" field.
func AutoDeployEQ(v bool) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldAutoDeploy, v))
}

// AutoDeployNEQ applies the NEQ predicate on the "auto_deploy" field.
func AutoDeployNEQ(v bool) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldAutoDeploy, v))
}

// RunCommandEQ applies the EQ predicate on the "run_command" field.
func RunCommandEQ(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldRunCommand, v))
}

// RunCommandNEQ applies the NEQ predicate on the "run_command" field.
func RunCommandNEQ(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldRunCommand, v))
}

// RunCommandIn applies the In predicate on the "run_command" field.
func RunCommandIn(vs ...string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIn(FieldRunCommand, vs...))
}

// RunCommandNotIn applies the NotIn predicate on the "run_command" field.
func RunCommandNotIn(vs ...string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotIn(FieldRunCommand, vs...))
}

// RunCommandGT applies the GT predicate on the "run_command" field.
func RunCommandGT(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGT(FieldRunCommand, v))
}

// RunCommandGTE applies the GTE predicate on the "run_command" field.
func RunCommandGTE(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGTE(FieldRunCommand, v))
}

// RunCommandLT applies the LT predicate on the "run_command" field.
func RunCommandLT(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLT(FieldRunCommand, v))
}

// RunCommandLTE applies the LTE predicate on the "run_command" field.
func RunCommandLTE(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLTE(FieldRunCommand, v))
}

// RunCommandContains applies the Contains predicate on the "run_command" field.
func RunCommandContains(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldContains(FieldRunCommand, v))
}

// RunCommandHasPrefix applies the HasPrefix predicate on the "run_command" field.
func RunCommandHasPrefix(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldHasPrefix(FieldRunCommand, v))
}

// RunCommandHasSuffix applies the HasSuffix predicate on the "run_command" field.
func RunCommandHasSuffix(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldHasSuffix(FieldRunCommand, v))
}

// RunCommandIsNil applies the IsNil predicate on the "run_command" field.
func RunCommandIsNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIsNull(FieldRunCommand))
}

// RunCommandNotNil applies the NotNil predicate on the "run_command" field.
func RunCommandNotNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotNull(FieldRunCommand))
}

// RunCommandEqualFold applies the EqualFold predicate on the "run_command" field.
func RunCommandEqualFold(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEqualFold(FieldRunCommand, v))
}

// RunCommandContainsFold applies the ContainsFold predicate on the "run_command" field.
func RunCommandContainsFold(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldContainsFold(FieldRunCommand, v))
}

// PublicEQ applies the EQ predicate on the "public" field.
func PublicEQ(v bool) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldPublic, v))
}

// PublicNEQ applies the NEQ predicate on the "public" field.
func PublicNEQ(v bool) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldPublic, v))
}

// ImageEQ applies the EQ predicate on the "image" field.
func ImageEQ(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEQ(FieldImage, v))
}

// ImageNEQ applies the NEQ predicate on the "image" field.
func ImageNEQ(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNEQ(FieldImage, v))
}

// ImageIn applies the In predicate on the "image" field.
func ImageIn(vs ...string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIn(FieldImage, vs...))
}

// ImageNotIn applies the NotIn predicate on the "image" field.
func ImageNotIn(vs ...string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotIn(FieldImage, vs...))
}

// ImageGT applies the GT predicate on the "image" field.
func ImageGT(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGT(FieldImage, v))
}

// ImageGTE applies the GTE predicate on the "image" field.
func ImageGTE(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldGTE(FieldImage, v))
}

// ImageLT applies the LT predicate on the "image" field.
func ImageLT(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLT(FieldImage, v))
}

// ImageLTE applies the LTE predicate on the "image" field.
func ImageLTE(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldLTE(FieldImage, v))
}

// ImageContains applies the Contains predicate on the "image" field.
func ImageContains(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldContains(FieldImage, v))
}

// ImageHasPrefix applies the HasPrefix predicate on the "image" field.
func ImageHasPrefix(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldHasPrefix(FieldImage, v))
}

// ImageHasSuffix applies the HasSuffix predicate on the "image" field.
func ImageHasSuffix(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldHasSuffix(FieldImage, v))
}

// ImageIsNil applies the IsNil predicate on the "image" field.
func ImageIsNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldIsNull(FieldImage))
}

// ImageNotNil applies the NotNil predicate on the "image" field.
func ImageNotNil() predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldNotNull(FieldImage))
}

// ImageEqualFold applies the EqualFold predicate on the "image" field.
func ImageEqualFold(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldEqualFold(FieldImage, v))
}

// ImageContainsFold applies the ContainsFold predicate on the "image" field.
func ImageContainsFold(v string) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.FieldContainsFold(FieldImage, v))
}

// HasService applies the HasEdge predicate on the "service" edge.
func HasService() predicate.ServiceConfig {
	return predicate.ServiceConfig(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, ServiceTable, ServiceColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServiceWith applies the HasEdge predicate on the "service" edge with a given conditions (other predicates).
func HasServiceWith(preds ...predicate.Service) predicate.ServiceConfig {
	return predicate.ServiceConfig(func(s *sql.Selector) {
		step := newServiceStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ServiceConfig) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ServiceConfig) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ServiceConfig) predicate.ServiceConfig {
	return predicate.ServiceConfig(sql.NotPredicates(p))
}
