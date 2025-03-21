// Code generated by ent, DO NOT EDIT.

package buildjob

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/predicate"
	"github.com/unbindapp/unbind-api/ent/schema"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldUpdatedAt, v))
}

// ServiceID applies equality check predicate on the "service_id" field. It's identical to ServiceIDEQ.
func ServiceID(v uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldServiceID, v))
}

// Error applies equality check predicate on the "error" field. It's identical to ErrorEQ.
func Error(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldError, v))
}

// StartedAt applies equality check predicate on the "started_at" field. It's identical to StartedAtEQ.
func StartedAt(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldStartedAt, v))
}

// CompletedAt applies equality check predicate on the "completed_at" field. It's identical to CompletedAtEQ.
func CompletedAt(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldCompletedAt, v))
}

// KubernetesJobName applies equality check predicate on the "kubernetes_job_name" field. It's identical to KubernetesJobNameEQ.
func KubernetesJobName(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldKubernetesJobName, v))
}

// KubernetesJobStatus applies equality check predicate on the "kubernetes_job_status" field. It's identical to KubernetesJobStatusEQ.
func KubernetesJobStatus(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldKubernetesJobStatus, v))
}

// Attempts applies equality check predicate on the "attempts" field. It's identical to AttemptsEQ.
func Attempts(v int) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldAttempts, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLTE(FieldUpdatedAt, v))
}

// ServiceIDEQ applies the EQ predicate on the "service_id" field.
func ServiceIDEQ(v uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldServiceID, v))
}

// ServiceIDNEQ applies the NEQ predicate on the "service_id" field.
func ServiceIDNEQ(v uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNEQ(FieldServiceID, v))
}

// ServiceIDIn applies the In predicate on the "service_id" field.
func ServiceIDIn(vs ...uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIn(FieldServiceID, vs...))
}

// ServiceIDNotIn applies the NotIn predicate on the "service_id" field.
func ServiceIDNotIn(vs ...uuid.UUID) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotIn(FieldServiceID, vs...))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v schema.BuildJobStatus) predicate.BuildJob {
	vc := v
	return predicate.BuildJob(sql.FieldEQ(FieldStatus, vc))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v schema.BuildJobStatus) predicate.BuildJob {
	vc := v
	return predicate.BuildJob(sql.FieldNEQ(FieldStatus, vc))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...schema.BuildJobStatus) predicate.BuildJob {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.BuildJob(sql.FieldIn(FieldStatus, v...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...schema.BuildJobStatus) predicate.BuildJob {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.BuildJob(sql.FieldNotIn(FieldStatus, v...))
}

// ErrorEQ applies the EQ predicate on the "error" field.
func ErrorEQ(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldError, v))
}

// ErrorNEQ applies the NEQ predicate on the "error" field.
func ErrorNEQ(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNEQ(FieldError, v))
}

// ErrorIn applies the In predicate on the "error" field.
func ErrorIn(vs ...string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIn(FieldError, vs...))
}

// ErrorNotIn applies the NotIn predicate on the "error" field.
func ErrorNotIn(vs ...string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotIn(FieldError, vs...))
}

// ErrorGT applies the GT predicate on the "error" field.
func ErrorGT(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGT(FieldError, v))
}

// ErrorGTE applies the GTE predicate on the "error" field.
func ErrorGTE(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGTE(FieldError, v))
}

// ErrorLT applies the LT predicate on the "error" field.
func ErrorLT(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLT(FieldError, v))
}

// ErrorLTE applies the LTE predicate on the "error" field.
func ErrorLTE(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLTE(FieldError, v))
}

// ErrorContains applies the Contains predicate on the "error" field.
func ErrorContains(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldContains(FieldError, v))
}

// ErrorHasPrefix applies the HasPrefix predicate on the "error" field.
func ErrorHasPrefix(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldHasPrefix(FieldError, v))
}

// ErrorHasSuffix applies the HasSuffix predicate on the "error" field.
func ErrorHasSuffix(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldHasSuffix(FieldError, v))
}

// ErrorIsNil applies the IsNil predicate on the "error" field.
func ErrorIsNil() predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIsNull(FieldError))
}

// ErrorNotNil applies the NotNil predicate on the "error" field.
func ErrorNotNil() predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotNull(FieldError))
}

// ErrorEqualFold applies the EqualFold predicate on the "error" field.
func ErrorEqualFold(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEqualFold(FieldError, v))
}

// ErrorContainsFold applies the ContainsFold predicate on the "error" field.
func ErrorContainsFold(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldContainsFold(FieldError, v))
}

// StartedAtEQ applies the EQ predicate on the "started_at" field.
func StartedAtEQ(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldStartedAt, v))
}

// StartedAtNEQ applies the NEQ predicate on the "started_at" field.
func StartedAtNEQ(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNEQ(FieldStartedAt, v))
}

// StartedAtIn applies the In predicate on the "started_at" field.
func StartedAtIn(vs ...time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIn(FieldStartedAt, vs...))
}

// StartedAtNotIn applies the NotIn predicate on the "started_at" field.
func StartedAtNotIn(vs ...time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotIn(FieldStartedAt, vs...))
}

// StartedAtGT applies the GT predicate on the "started_at" field.
func StartedAtGT(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGT(FieldStartedAt, v))
}

// StartedAtGTE applies the GTE predicate on the "started_at" field.
func StartedAtGTE(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGTE(FieldStartedAt, v))
}

// StartedAtLT applies the LT predicate on the "started_at" field.
func StartedAtLT(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLT(FieldStartedAt, v))
}

// StartedAtLTE applies the LTE predicate on the "started_at" field.
func StartedAtLTE(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLTE(FieldStartedAt, v))
}

// StartedAtIsNil applies the IsNil predicate on the "started_at" field.
func StartedAtIsNil() predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIsNull(FieldStartedAt))
}

// StartedAtNotNil applies the NotNil predicate on the "started_at" field.
func StartedAtNotNil() predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotNull(FieldStartedAt))
}

// CompletedAtEQ applies the EQ predicate on the "completed_at" field.
func CompletedAtEQ(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldCompletedAt, v))
}

// CompletedAtNEQ applies the NEQ predicate on the "completed_at" field.
func CompletedAtNEQ(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNEQ(FieldCompletedAt, v))
}

// CompletedAtIn applies the In predicate on the "completed_at" field.
func CompletedAtIn(vs ...time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIn(FieldCompletedAt, vs...))
}

// CompletedAtNotIn applies the NotIn predicate on the "completed_at" field.
func CompletedAtNotIn(vs ...time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotIn(FieldCompletedAt, vs...))
}

// CompletedAtGT applies the GT predicate on the "completed_at" field.
func CompletedAtGT(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGT(FieldCompletedAt, v))
}

// CompletedAtGTE applies the GTE predicate on the "completed_at" field.
func CompletedAtGTE(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGTE(FieldCompletedAt, v))
}

// CompletedAtLT applies the LT predicate on the "completed_at" field.
func CompletedAtLT(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLT(FieldCompletedAt, v))
}

// CompletedAtLTE applies the LTE predicate on the "completed_at" field.
func CompletedAtLTE(v time.Time) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLTE(FieldCompletedAt, v))
}

// CompletedAtIsNil applies the IsNil predicate on the "completed_at" field.
func CompletedAtIsNil() predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIsNull(FieldCompletedAt))
}

// CompletedAtNotNil applies the NotNil predicate on the "completed_at" field.
func CompletedAtNotNil() predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotNull(FieldCompletedAt))
}

// KubernetesJobNameEQ applies the EQ predicate on the "kubernetes_job_name" field.
func KubernetesJobNameEQ(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldKubernetesJobName, v))
}

// KubernetesJobNameNEQ applies the NEQ predicate on the "kubernetes_job_name" field.
func KubernetesJobNameNEQ(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNEQ(FieldKubernetesJobName, v))
}

// KubernetesJobNameIn applies the In predicate on the "kubernetes_job_name" field.
func KubernetesJobNameIn(vs ...string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIn(FieldKubernetesJobName, vs...))
}

// KubernetesJobNameNotIn applies the NotIn predicate on the "kubernetes_job_name" field.
func KubernetesJobNameNotIn(vs ...string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotIn(FieldKubernetesJobName, vs...))
}

// KubernetesJobNameGT applies the GT predicate on the "kubernetes_job_name" field.
func KubernetesJobNameGT(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGT(FieldKubernetesJobName, v))
}

// KubernetesJobNameGTE applies the GTE predicate on the "kubernetes_job_name" field.
func KubernetesJobNameGTE(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGTE(FieldKubernetesJobName, v))
}

// KubernetesJobNameLT applies the LT predicate on the "kubernetes_job_name" field.
func KubernetesJobNameLT(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLT(FieldKubernetesJobName, v))
}

// KubernetesJobNameLTE applies the LTE predicate on the "kubernetes_job_name" field.
func KubernetesJobNameLTE(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLTE(FieldKubernetesJobName, v))
}

// KubernetesJobNameContains applies the Contains predicate on the "kubernetes_job_name" field.
func KubernetesJobNameContains(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldContains(FieldKubernetesJobName, v))
}

// KubernetesJobNameHasPrefix applies the HasPrefix predicate on the "kubernetes_job_name" field.
func KubernetesJobNameHasPrefix(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldHasPrefix(FieldKubernetesJobName, v))
}

// KubernetesJobNameHasSuffix applies the HasSuffix predicate on the "kubernetes_job_name" field.
func KubernetesJobNameHasSuffix(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldHasSuffix(FieldKubernetesJobName, v))
}

// KubernetesJobNameIsNil applies the IsNil predicate on the "kubernetes_job_name" field.
func KubernetesJobNameIsNil() predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIsNull(FieldKubernetesJobName))
}

// KubernetesJobNameNotNil applies the NotNil predicate on the "kubernetes_job_name" field.
func KubernetesJobNameNotNil() predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotNull(FieldKubernetesJobName))
}

// KubernetesJobNameEqualFold applies the EqualFold predicate on the "kubernetes_job_name" field.
func KubernetesJobNameEqualFold(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEqualFold(FieldKubernetesJobName, v))
}

// KubernetesJobNameContainsFold applies the ContainsFold predicate on the "kubernetes_job_name" field.
func KubernetesJobNameContainsFold(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldContainsFold(FieldKubernetesJobName, v))
}

// KubernetesJobStatusEQ applies the EQ predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusEQ(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldKubernetesJobStatus, v))
}

// KubernetesJobStatusNEQ applies the NEQ predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusNEQ(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNEQ(FieldKubernetesJobStatus, v))
}

// KubernetesJobStatusIn applies the In predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusIn(vs ...string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIn(FieldKubernetesJobStatus, vs...))
}

// KubernetesJobStatusNotIn applies the NotIn predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusNotIn(vs ...string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotIn(FieldKubernetesJobStatus, vs...))
}

// KubernetesJobStatusGT applies the GT predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusGT(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGT(FieldKubernetesJobStatus, v))
}

// KubernetesJobStatusGTE applies the GTE predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusGTE(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGTE(FieldKubernetesJobStatus, v))
}

// KubernetesJobStatusLT applies the LT predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusLT(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLT(FieldKubernetesJobStatus, v))
}

// KubernetesJobStatusLTE applies the LTE predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusLTE(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLTE(FieldKubernetesJobStatus, v))
}

// KubernetesJobStatusContains applies the Contains predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusContains(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldContains(FieldKubernetesJobStatus, v))
}

// KubernetesJobStatusHasPrefix applies the HasPrefix predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusHasPrefix(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldHasPrefix(FieldKubernetesJobStatus, v))
}

// KubernetesJobStatusHasSuffix applies the HasSuffix predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusHasSuffix(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldHasSuffix(FieldKubernetesJobStatus, v))
}

// KubernetesJobStatusIsNil applies the IsNil predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusIsNil() predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIsNull(FieldKubernetesJobStatus))
}

// KubernetesJobStatusNotNil applies the NotNil predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusNotNil() predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotNull(FieldKubernetesJobStatus))
}

// KubernetesJobStatusEqualFold applies the EqualFold predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusEqualFold(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEqualFold(FieldKubernetesJobStatus, v))
}

// KubernetesJobStatusContainsFold applies the ContainsFold predicate on the "kubernetes_job_status" field.
func KubernetesJobStatusContainsFold(v string) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldContainsFold(FieldKubernetesJobStatus, v))
}

// AttemptsEQ applies the EQ predicate on the "attempts" field.
func AttemptsEQ(v int) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldEQ(FieldAttempts, v))
}

// AttemptsNEQ applies the NEQ predicate on the "attempts" field.
func AttemptsNEQ(v int) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNEQ(FieldAttempts, v))
}

// AttemptsIn applies the In predicate on the "attempts" field.
func AttemptsIn(vs ...int) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldIn(FieldAttempts, vs...))
}

// AttemptsNotIn applies the NotIn predicate on the "attempts" field.
func AttemptsNotIn(vs ...int) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldNotIn(FieldAttempts, vs...))
}

// AttemptsGT applies the GT predicate on the "attempts" field.
func AttemptsGT(v int) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGT(FieldAttempts, v))
}

// AttemptsGTE applies the GTE predicate on the "attempts" field.
func AttemptsGTE(v int) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldGTE(FieldAttempts, v))
}

// AttemptsLT applies the LT predicate on the "attempts" field.
func AttemptsLT(v int) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLT(FieldAttempts, v))
}

// AttemptsLTE applies the LTE predicate on the "attempts" field.
func AttemptsLTE(v int) predicate.BuildJob {
	return predicate.BuildJob(sql.FieldLTE(FieldAttempts, v))
}

// HasService applies the HasEdge predicate on the "service" edge.
func HasService() predicate.BuildJob {
	return predicate.BuildJob(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ServiceTable, ServiceColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServiceWith applies the HasEdge predicate on the "service" edge with a given conditions (other predicates).
func HasServiceWith(preds ...predicate.Service) predicate.BuildJob {
	return predicate.BuildJob(func(s *sql.Selector) {
		step := newServiceStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.BuildJob) predicate.BuildJob {
	return predicate.BuildJob(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.BuildJob) predicate.BuildJob {
	return predicate.BuildJob(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.BuildJob) predicate.BuildJob {
	return predicate.BuildJob(sql.NotPredicates(p))
}
