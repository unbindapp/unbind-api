// Code generated by ent, DO NOT EDIT.

package permission

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Permission {
	return predicate.Permission(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Permission {
	return predicate.Permission(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Permission {
	return predicate.Permission(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Permission {
	return predicate.Permission(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Permission {
	return predicate.Permission(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Permission {
	return predicate.Permission(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Permission {
	return predicate.Permission(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldUpdatedAt, v))
}

// ResourceID applies equality check predicate on the "resource_id" field. It's identical to ResourceIDEQ.
func ResourceID(v string) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldResourceID, v))
}

// Scope applies equality check predicate on the "scope" field. It's identical to ScopeEQ.
func Scope(v string) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldScope, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Permission {
	return predicate.Permission(sql.FieldLTE(FieldUpdatedAt, v))
}

// ActionEQ applies the EQ predicate on the "action" field.
func ActionEQ(v Action) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldAction, v))
}

// ActionNEQ applies the NEQ predicate on the "action" field.
func ActionNEQ(v Action) predicate.Permission {
	return predicate.Permission(sql.FieldNEQ(FieldAction, v))
}

// ActionIn applies the In predicate on the "action" field.
func ActionIn(vs ...Action) predicate.Permission {
	return predicate.Permission(sql.FieldIn(FieldAction, vs...))
}

// ActionNotIn applies the NotIn predicate on the "action" field.
func ActionNotIn(vs ...Action) predicate.Permission {
	return predicate.Permission(sql.FieldNotIn(FieldAction, vs...))
}

// ResourceTypeEQ applies the EQ predicate on the "resource_type" field.
func ResourceTypeEQ(v ResourceType) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldResourceType, v))
}

// ResourceTypeNEQ applies the NEQ predicate on the "resource_type" field.
func ResourceTypeNEQ(v ResourceType) predicate.Permission {
	return predicate.Permission(sql.FieldNEQ(FieldResourceType, v))
}

// ResourceTypeIn applies the In predicate on the "resource_type" field.
func ResourceTypeIn(vs ...ResourceType) predicate.Permission {
	return predicate.Permission(sql.FieldIn(FieldResourceType, vs...))
}

// ResourceTypeNotIn applies the NotIn predicate on the "resource_type" field.
func ResourceTypeNotIn(vs ...ResourceType) predicate.Permission {
	return predicate.Permission(sql.FieldNotIn(FieldResourceType, vs...))
}

// ResourceIDEQ applies the EQ predicate on the "resource_id" field.
func ResourceIDEQ(v string) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldResourceID, v))
}

// ResourceIDNEQ applies the NEQ predicate on the "resource_id" field.
func ResourceIDNEQ(v string) predicate.Permission {
	return predicate.Permission(sql.FieldNEQ(FieldResourceID, v))
}

// ResourceIDIn applies the In predicate on the "resource_id" field.
func ResourceIDIn(vs ...string) predicate.Permission {
	return predicate.Permission(sql.FieldIn(FieldResourceID, vs...))
}

// ResourceIDNotIn applies the NotIn predicate on the "resource_id" field.
func ResourceIDNotIn(vs ...string) predicate.Permission {
	return predicate.Permission(sql.FieldNotIn(FieldResourceID, vs...))
}

// ResourceIDGT applies the GT predicate on the "resource_id" field.
func ResourceIDGT(v string) predicate.Permission {
	return predicate.Permission(sql.FieldGT(FieldResourceID, v))
}

// ResourceIDGTE applies the GTE predicate on the "resource_id" field.
func ResourceIDGTE(v string) predicate.Permission {
	return predicate.Permission(sql.FieldGTE(FieldResourceID, v))
}

// ResourceIDLT applies the LT predicate on the "resource_id" field.
func ResourceIDLT(v string) predicate.Permission {
	return predicate.Permission(sql.FieldLT(FieldResourceID, v))
}

// ResourceIDLTE applies the LTE predicate on the "resource_id" field.
func ResourceIDLTE(v string) predicate.Permission {
	return predicate.Permission(sql.FieldLTE(FieldResourceID, v))
}

// ResourceIDContains applies the Contains predicate on the "resource_id" field.
func ResourceIDContains(v string) predicate.Permission {
	return predicate.Permission(sql.FieldContains(FieldResourceID, v))
}

// ResourceIDHasPrefix applies the HasPrefix predicate on the "resource_id" field.
func ResourceIDHasPrefix(v string) predicate.Permission {
	return predicate.Permission(sql.FieldHasPrefix(FieldResourceID, v))
}

// ResourceIDHasSuffix applies the HasSuffix predicate on the "resource_id" field.
func ResourceIDHasSuffix(v string) predicate.Permission {
	return predicate.Permission(sql.FieldHasSuffix(FieldResourceID, v))
}

// ResourceIDEqualFold applies the EqualFold predicate on the "resource_id" field.
func ResourceIDEqualFold(v string) predicate.Permission {
	return predicate.Permission(sql.FieldEqualFold(FieldResourceID, v))
}

// ResourceIDContainsFold applies the ContainsFold predicate on the "resource_id" field.
func ResourceIDContainsFold(v string) predicate.Permission {
	return predicate.Permission(sql.FieldContainsFold(FieldResourceID, v))
}

// ScopeEQ applies the EQ predicate on the "scope" field.
func ScopeEQ(v string) predicate.Permission {
	return predicate.Permission(sql.FieldEQ(FieldScope, v))
}

// ScopeNEQ applies the NEQ predicate on the "scope" field.
func ScopeNEQ(v string) predicate.Permission {
	return predicate.Permission(sql.FieldNEQ(FieldScope, v))
}

// ScopeIn applies the In predicate on the "scope" field.
func ScopeIn(vs ...string) predicate.Permission {
	return predicate.Permission(sql.FieldIn(FieldScope, vs...))
}

// ScopeNotIn applies the NotIn predicate on the "scope" field.
func ScopeNotIn(vs ...string) predicate.Permission {
	return predicate.Permission(sql.FieldNotIn(FieldScope, vs...))
}

// ScopeGT applies the GT predicate on the "scope" field.
func ScopeGT(v string) predicate.Permission {
	return predicate.Permission(sql.FieldGT(FieldScope, v))
}

// ScopeGTE applies the GTE predicate on the "scope" field.
func ScopeGTE(v string) predicate.Permission {
	return predicate.Permission(sql.FieldGTE(FieldScope, v))
}

// ScopeLT applies the LT predicate on the "scope" field.
func ScopeLT(v string) predicate.Permission {
	return predicate.Permission(sql.FieldLT(FieldScope, v))
}

// ScopeLTE applies the LTE predicate on the "scope" field.
func ScopeLTE(v string) predicate.Permission {
	return predicate.Permission(sql.FieldLTE(FieldScope, v))
}

// ScopeContains applies the Contains predicate on the "scope" field.
func ScopeContains(v string) predicate.Permission {
	return predicate.Permission(sql.FieldContains(FieldScope, v))
}

// ScopeHasPrefix applies the HasPrefix predicate on the "scope" field.
func ScopeHasPrefix(v string) predicate.Permission {
	return predicate.Permission(sql.FieldHasPrefix(FieldScope, v))
}

// ScopeHasSuffix applies the HasSuffix predicate on the "scope" field.
func ScopeHasSuffix(v string) predicate.Permission {
	return predicate.Permission(sql.FieldHasSuffix(FieldScope, v))
}

// ScopeIsNil applies the IsNil predicate on the "scope" field.
func ScopeIsNil() predicate.Permission {
	return predicate.Permission(sql.FieldIsNull(FieldScope))
}

// ScopeNotNil applies the NotNil predicate on the "scope" field.
func ScopeNotNil() predicate.Permission {
	return predicate.Permission(sql.FieldNotNull(FieldScope))
}

// ScopeEqualFold applies the EqualFold predicate on the "scope" field.
func ScopeEqualFold(v string) predicate.Permission {
	return predicate.Permission(sql.FieldEqualFold(FieldScope, v))
}

// ScopeContainsFold applies the ContainsFold predicate on the "scope" field.
func ScopeContainsFold(v string) predicate.Permission {
	return predicate.Permission(sql.FieldContainsFold(FieldScope, v))
}

// LabelsIsNil applies the IsNil predicate on the "labels" field.
func LabelsIsNil() predicate.Permission {
	return predicate.Permission(sql.FieldIsNull(FieldLabels))
}

// LabelsNotNil applies the NotNil predicate on the "labels" field.
func LabelsNotNil() predicate.Permission {
	return predicate.Permission(sql.FieldNotNull(FieldLabels))
}

// HasGroups applies the HasEdge predicate on the "groups" edge.
func HasGroups() predicate.Permission {
	return predicate.Permission(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, GroupsTable, GroupsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasGroupsWith applies the HasEdge predicate on the "groups" edge with a given conditions (other predicates).
func HasGroupsWith(preds ...predicate.Group) predicate.Permission {
	return predicate.Permission(func(s *sql.Selector) {
		step := newGroupsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Permission) predicate.Permission {
	return predicate.Permission(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Permission) predicate.Permission {
	return predicate.Permission(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Permission) predicate.Permission {
	return predicate.Permission(sql.NotPredicates(p))
}
