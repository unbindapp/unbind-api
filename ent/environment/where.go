// Code generated by ent, DO NOT EDIT.

package environment

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldUpdatedAt, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldName, v))
}

// DisplayName applies equality check predicate on the "display_name" field. It's identical to DisplayNameEQ.
func DisplayName(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldDisplayName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldDescription, v))
}

// Active applies equality check predicate on the "active" field. It's identical to ActiveEQ.
func Active(v bool) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldActive, v))
}

// ProjectID applies equality check predicate on the "project_id" field. It's identical to ProjectIDEQ.
func ProjectID(v uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldProjectID, v))
}

// KubernetesSecret applies equality check predicate on the "kubernetes_secret" field. It's identical to KubernetesSecretEQ.
func KubernetesSecret(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldKubernetesSecret, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldLTE(FieldUpdatedAt, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Environment {
	return predicate.Environment(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Environment {
	return predicate.Environment(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Environment {
	return predicate.Environment(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Environment {
	return predicate.Environment(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Environment {
	return predicate.Environment(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Environment {
	return predicate.Environment(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Environment {
	return predicate.Environment(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Environment {
	return predicate.Environment(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Environment {
	return predicate.Environment(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Environment {
	return predicate.Environment(sql.FieldContainsFold(FieldName, v))
}

// DisplayNameEQ applies the EQ predicate on the "display_name" field.
func DisplayNameEQ(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldDisplayName, v))
}

// DisplayNameNEQ applies the NEQ predicate on the "display_name" field.
func DisplayNameNEQ(v string) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldDisplayName, v))
}

// DisplayNameIn applies the In predicate on the "display_name" field.
func DisplayNameIn(vs ...string) predicate.Environment {
	return predicate.Environment(sql.FieldIn(FieldDisplayName, vs...))
}

// DisplayNameNotIn applies the NotIn predicate on the "display_name" field.
func DisplayNameNotIn(vs ...string) predicate.Environment {
	return predicate.Environment(sql.FieldNotIn(FieldDisplayName, vs...))
}

// DisplayNameGT applies the GT predicate on the "display_name" field.
func DisplayNameGT(v string) predicate.Environment {
	return predicate.Environment(sql.FieldGT(FieldDisplayName, v))
}

// DisplayNameGTE applies the GTE predicate on the "display_name" field.
func DisplayNameGTE(v string) predicate.Environment {
	return predicate.Environment(sql.FieldGTE(FieldDisplayName, v))
}

// DisplayNameLT applies the LT predicate on the "display_name" field.
func DisplayNameLT(v string) predicate.Environment {
	return predicate.Environment(sql.FieldLT(FieldDisplayName, v))
}

// DisplayNameLTE applies the LTE predicate on the "display_name" field.
func DisplayNameLTE(v string) predicate.Environment {
	return predicate.Environment(sql.FieldLTE(FieldDisplayName, v))
}

// DisplayNameContains applies the Contains predicate on the "display_name" field.
func DisplayNameContains(v string) predicate.Environment {
	return predicate.Environment(sql.FieldContains(FieldDisplayName, v))
}

// DisplayNameHasPrefix applies the HasPrefix predicate on the "display_name" field.
func DisplayNameHasPrefix(v string) predicate.Environment {
	return predicate.Environment(sql.FieldHasPrefix(FieldDisplayName, v))
}

// DisplayNameHasSuffix applies the HasSuffix predicate on the "display_name" field.
func DisplayNameHasSuffix(v string) predicate.Environment {
	return predicate.Environment(sql.FieldHasSuffix(FieldDisplayName, v))
}

// DisplayNameEqualFold applies the EqualFold predicate on the "display_name" field.
func DisplayNameEqualFold(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEqualFold(FieldDisplayName, v))
}

// DisplayNameContainsFold applies the ContainsFold predicate on the "display_name" field.
func DisplayNameContainsFold(v string) predicate.Environment {
	return predicate.Environment(sql.FieldContainsFold(FieldDisplayName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Environment {
	return predicate.Environment(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Environment {
	return predicate.Environment(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Environment {
	return predicate.Environment(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Environment {
	return predicate.Environment(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Environment {
	return predicate.Environment(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Environment {
	return predicate.Environment(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Environment {
	return predicate.Environment(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Environment {
	return predicate.Environment(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Environment {
	return predicate.Environment(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.Environment {
	return predicate.Environment(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.Environment {
	return predicate.Environment(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Environment {
	return predicate.Environment(sql.FieldContainsFold(FieldDescription, v))
}

// ActiveEQ applies the EQ predicate on the "active" field.
func ActiveEQ(v bool) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldActive, v))
}

// ActiveNEQ applies the NEQ predicate on the "active" field.
func ActiveNEQ(v bool) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldActive, v))
}

// ProjectIDEQ applies the EQ predicate on the "project_id" field.
func ProjectIDEQ(v uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldProjectID, v))
}

// ProjectIDNEQ applies the NEQ predicate on the "project_id" field.
func ProjectIDNEQ(v uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldProjectID, v))
}

// ProjectIDIn applies the In predicate on the "project_id" field.
func ProjectIDIn(vs ...uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldIn(FieldProjectID, vs...))
}

// ProjectIDNotIn applies the NotIn predicate on the "project_id" field.
func ProjectIDNotIn(vs ...uuid.UUID) predicate.Environment {
	return predicate.Environment(sql.FieldNotIn(FieldProjectID, vs...))
}

// KubernetesSecretEQ applies the EQ predicate on the "kubernetes_secret" field.
func KubernetesSecretEQ(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldKubernetesSecret, v))
}

// KubernetesSecretNEQ applies the NEQ predicate on the "kubernetes_secret" field.
func KubernetesSecretNEQ(v string) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldKubernetesSecret, v))
}

// KubernetesSecretIn applies the In predicate on the "kubernetes_secret" field.
func KubernetesSecretIn(vs ...string) predicate.Environment {
	return predicate.Environment(sql.FieldIn(FieldKubernetesSecret, vs...))
}

// KubernetesSecretNotIn applies the NotIn predicate on the "kubernetes_secret" field.
func KubernetesSecretNotIn(vs ...string) predicate.Environment {
	return predicate.Environment(sql.FieldNotIn(FieldKubernetesSecret, vs...))
}

// KubernetesSecretGT applies the GT predicate on the "kubernetes_secret" field.
func KubernetesSecretGT(v string) predicate.Environment {
	return predicate.Environment(sql.FieldGT(FieldKubernetesSecret, v))
}

// KubernetesSecretGTE applies the GTE predicate on the "kubernetes_secret" field.
func KubernetesSecretGTE(v string) predicate.Environment {
	return predicate.Environment(sql.FieldGTE(FieldKubernetesSecret, v))
}

// KubernetesSecretLT applies the LT predicate on the "kubernetes_secret" field.
func KubernetesSecretLT(v string) predicate.Environment {
	return predicate.Environment(sql.FieldLT(FieldKubernetesSecret, v))
}

// KubernetesSecretLTE applies the LTE predicate on the "kubernetes_secret" field.
func KubernetesSecretLTE(v string) predicate.Environment {
	return predicate.Environment(sql.FieldLTE(FieldKubernetesSecret, v))
}

// KubernetesSecretContains applies the Contains predicate on the "kubernetes_secret" field.
func KubernetesSecretContains(v string) predicate.Environment {
	return predicate.Environment(sql.FieldContains(FieldKubernetesSecret, v))
}

// KubernetesSecretHasPrefix applies the HasPrefix predicate on the "kubernetes_secret" field.
func KubernetesSecretHasPrefix(v string) predicate.Environment {
	return predicate.Environment(sql.FieldHasPrefix(FieldKubernetesSecret, v))
}

// KubernetesSecretHasSuffix applies the HasSuffix predicate on the "kubernetes_secret" field.
func KubernetesSecretHasSuffix(v string) predicate.Environment {
	return predicate.Environment(sql.FieldHasSuffix(FieldKubernetesSecret, v))
}

// KubernetesSecretEqualFold applies the EqualFold predicate on the "kubernetes_secret" field.
func KubernetesSecretEqualFold(v string) predicate.Environment {
	return predicate.Environment(sql.FieldEqualFold(FieldKubernetesSecret, v))
}

// KubernetesSecretContainsFold applies the ContainsFold predicate on the "kubernetes_secret" field.
func KubernetesSecretContainsFold(v string) predicate.Environment {
	return predicate.Environment(sql.FieldContainsFold(FieldKubernetesSecret, v))
}

// HasProject applies the HasEdge predicate on the "project" edge.
func HasProject() predicate.Environment {
	return predicate.Environment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectWith applies the HasEdge predicate on the "project" edge with a given conditions (other predicates).
func HasProjectWith(preds ...predicate.Project) predicate.Environment {
	return predicate.Environment(func(s *sql.Selector) {
		step := newProjectStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasServices applies the HasEdge predicate on the "services" edge.
func HasServices() predicate.Environment {
	return predicate.Environment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ServicesTable, ServicesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServicesWith applies the HasEdge predicate on the "services" edge with a given conditions (other predicates).
func HasServicesWith(preds ...predicate.Service) predicate.Environment {
	return predicate.Environment(func(s *sql.Selector) {
		step := newServicesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProjectDefault applies the HasEdge predicate on the "project_default" edge.
func HasProjectDefault() predicate.Environment {
	return predicate.Environment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ProjectDefaultTable, ProjectDefaultColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectDefaultWith applies the HasEdge predicate on the "project_default" edge with a given conditions (other predicates).
func HasProjectDefaultWith(preds ...predicate.Project) predicate.Environment {
	return predicate.Environment(func(s *sql.Selector) {
		step := newProjectDefaultStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Environment) predicate.Environment {
	return predicate.Environment(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Environment) predicate.Environment {
	return predicate.Environment(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Environment) predicate.Environment {
	return predicate.Environment(sql.NotPredicates(p))
}
