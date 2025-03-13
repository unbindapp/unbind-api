// Code generated by ent, DO NOT EDIT.

package group

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldUpdatedAt, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldDescription, v))
}

// Superuser applies equality check predicate on the "superuser" field. It's identical to SuperuserEQ.
func Superuser(v bool) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldSuperuser, v))
}

// K8sRoleName applies equality check predicate on the "k8s_role_name" field. It's identical to K8sRoleNameEQ.
func K8sRoleName(v string) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldK8sRoleName, v))
}

// IdentityProvider applies equality check predicate on the "identity_provider" field. It's identical to IdentityProviderEQ.
func IdentityProvider(v string) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldIdentityProvider, v))
}

// ExternalID applies equality check predicate on the "external_id" field. It's identical to ExternalIDEQ.
func ExternalID(v string) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldExternalID, v))
}

// TeamID applies equality check predicate on the "team_id" field. It's identical to TeamIDEQ.
func TeamID(v uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldTeamID, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Group {
	return predicate.Group(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Group {
	return predicate.Group(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Group {
	return predicate.Group(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Group {
	return predicate.Group(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Group {
	return predicate.Group(sql.FieldLTE(FieldUpdatedAt, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Group {
	return predicate.Group(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Group {
	return predicate.Group(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Group {
	return predicate.Group(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Group {
	return predicate.Group(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Group {
	return predicate.Group(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Group {
	return predicate.Group(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Group {
	return predicate.Group(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Group {
	return predicate.Group(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Group {
	return predicate.Group(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Group {
	return predicate.Group(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Group {
	return predicate.Group(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Group {
	return predicate.Group(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Group {
	return predicate.Group(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Group {
	return predicate.Group(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Group {
	return predicate.Group(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Group {
	return predicate.Group(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Group {
	return predicate.Group(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Group {
	return predicate.Group(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Group {
	return predicate.Group(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Group {
	return predicate.Group(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Group {
	return predicate.Group(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Group {
	return predicate.Group(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.Group {
	return predicate.Group(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.Group {
	return predicate.Group(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Group {
	return predicate.Group(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Group {
	return predicate.Group(sql.FieldContainsFold(FieldDescription, v))
}

// SuperuserEQ applies the EQ predicate on the "superuser" field.
func SuperuserEQ(v bool) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldSuperuser, v))
}

// SuperuserNEQ applies the NEQ predicate on the "superuser" field.
func SuperuserNEQ(v bool) predicate.Group {
	return predicate.Group(sql.FieldNEQ(FieldSuperuser, v))
}

// K8sRoleNameEQ applies the EQ predicate on the "k8s_role_name" field.
func K8sRoleNameEQ(v string) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldK8sRoleName, v))
}

// K8sRoleNameNEQ applies the NEQ predicate on the "k8s_role_name" field.
func K8sRoleNameNEQ(v string) predicate.Group {
	return predicate.Group(sql.FieldNEQ(FieldK8sRoleName, v))
}

// K8sRoleNameIn applies the In predicate on the "k8s_role_name" field.
func K8sRoleNameIn(vs ...string) predicate.Group {
	return predicate.Group(sql.FieldIn(FieldK8sRoleName, vs...))
}

// K8sRoleNameNotIn applies the NotIn predicate on the "k8s_role_name" field.
func K8sRoleNameNotIn(vs ...string) predicate.Group {
	return predicate.Group(sql.FieldNotIn(FieldK8sRoleName, vs...))
}

// K8sRoleNameGT applies the GT predicate on the "k8s_role_name" field.
func K8sRoleNameGT(v string) predicate.Group {
	return predicate.Group(sql.FieldGT(FieldK8sRoleName, v))
}

// K8sRoleNameGTE applies the GTE predicate on the "k8s_role_name" field.
func K8sRoleNameGTE(v string) predicate.Group {
	return predicate.Group(sql.FieldGTE(FieldK8sRoleName, v))
}

// K8sRoleNameLT applies the LT predicate on the "k8s_role_name" field.
func K8sRoleNameLT(v string) predicate.Group {
	return predicate.Group(sql.FieldLT(FieldK8sRoleName, v))
}

// K8sRoleNameLTE applies the LTE predicate on the "k8s_role_name" field.
func K8sRoleNameLTE(v string) predicate.Group {
	return predicate.Group(sql.FieldLTE(FieldK8sRoleName, v))
}

// K8sRoleNameContains applies the Contains predicate on the "k8s_role_name" field.
func K8sRoleNameContains(v string) predicate.Group {
	return predicate.Group(sql.FieldContains(FieldK8sRoleName, v))
}

// K8sRoleNameHasPrefix applies the HasPrefix predicate on the "k8s_role_name" field.
func K8sRoleNameHasPrefix(v string) predicate.Group {
	return predicate.Group(sql.FieldHasPrefix(FieldK8sRoleName, v))
}

// K8sRoleNameHasSuffix applies the HasSuffix predicate on the "k8s_role_name" field.
func K8sRoleNameHasSuffix(v string) predicate.Group {
	return predicate.Group(sql.FieldHasSuffix(FieldK8sRoleName, v))
}

// K8sRoleNameIsNil applies the IsNil predicate on the "k8s_role_name" field.
func K8sRoleNameIsNil() predicate.Group {
	return predicate.Group(sql.FieldIsNull(FieldK8sRoleName))
}

// K8sRoleNameNotNil applies the NotNil predicate on the "k8s_role_name" field.
func K8sRoleNameNotNil() predicate.Group {
	return predicate.Group(sql.FieldNotNull(FieldK8sRoleName))
}

// K8sRoleNameEqualFold applies the EqualFold predicate on the "k8s_role_name" field.
func K8sRoleNameEqualFold(v string) predicate.Group {
	return predicate.Group(sql.FieldEqualFold(FieldK8sRoleName, v))
}

// K8sRoleNameContainsFold applies the ContainsFold predicate on the "k8s_role_name" field.
func K8sRoleNameContainsFold(v string) predicate.Group {
	return predicate.Group(sql.FieldContainsFold(FieldK8sRoleName, v))
}

// IdentityProviderEQ applies the EQ predicate on the "identity_provider" field.
func IdentityProviderEQ(v string) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldIdentityProvider, v))
}

// IdentityProviderNEQ applies the NEQ predicate on the "identity_provider" field.
func IdentityProviderNEQ(v string) predicate.Group {
	return predicate.Group(sql.FieldNEQ(FieldIdentityProvider, v))
}

// IdentityProviderIn applies the In predicate on the "identity_provider" field.
func IdentityProviderIn(vs ...string) predicate.Group {
	return predicate.Group(sql.FieldIn(FieldIdentityProvider, vs...))
}

// IdentityProviderNotIn applies the NotIn predicate on the "identity_provider" field.
func IdentityProviderNotIn(vs ...string) predicate.Group {
	return predicate.Group(sql.FieldNotIn(FieldIdentityProvider, vs...))
}

// IdentityProviderGT applies the GT predicate on the "identity_provider" field.
func IdentityProviderGT(v string) predicate.Group {
	return predicate.Group(sql.FieldGT(FieldIdentityProvider, v))
}

// IdentityProviderGTE applies the GTE predicate on the "identity_provider" field.
func IdentityProviderGTE(v string) predicate.Group {
	return predicate.Group(sql.FieldGTE(FieldIdentityProvider, v))
}

// IdentityProviderLT applies the LT predicate on the "identity_provider" field.
func IdentityProviderLT(v string) predicate.Group {
	return predicate.Group(sql.FieldLT(FieldIdentityProvider, v))
}

// IdentityProviderLTE applies the LTE predicate on the "identity_provider" field.
func IdentityProviderLTE(v string) predicate.Group {
	return predicate.Group(sql.FieldLTE(FieldIdentityProvider, v))
}

// IdentityProviderContains applies the Contains predicate on the "identity_provider" field.
func IdentityProviderContains(v string) predicate.Group {
	return predicate.Group(sql.FieldContains(FieldIdentityProvider, v))
}

// IdentityProviderHasPrefix applies the HasPrefix predicate on the "identity_provider" field.
func IdentityProviderHasPrefix(v string) predicate.Group {
	return predicate.Group(sql.FieldHasPrefix(FieldIdentityProvider, v))
}

// IdentityProviderHasSuffix applies the HasSuffix predicate on the "identity_provider" field.
func IdentityProviderHasSuffix(v string) predicate.Group {
	return predicate.Group(sql.FieldHasSuffix(FieldIdentityProvider, v))
}

// IdentityProviderIsNil applies the IsNil predicate on the "identity_provider" field.
func IdentityProviderIsNil() predicate.Group {
	return predicate.Group(sql.FieldIsNull(FieldIdentityProvider))
}

// IdentityProviderNotNil applies the NotNil predicate on the "identity_provider" field.
func IdentityProviderNotNil() predicate.Group {
	return predicate.Group(sql.FieldNotNull(FieldIdentityProvider))
}

// IdentityProviderEqualFold applies the EqualFold predicate on the "identity_provider" field.
func IdentityProviderEqualFold(v string) predicate.Group {
	return predicate.Group(sql.FieldEqualFold(FieldIdentityProvider, v))
}

// IdentityProviderContainsFold applies the ContainsFold predicate on the "identity_provider" field.
func IdentityProviderContainsFold(v string) predicate.Group {
	return predicate.Group(sql.FieldContainsFold(FieldIdentityProvider, v))
}

// ExternalIDEQ applies the EQ predicate on the "external_id" field.
func ExternalIDEQ(v string) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldExternalID, v))
}

// ExternalIDNEQ applies the NEQ predicate on the "external_id" field.
func ExternalIDNEQ(v string) predicate.Group {
	return predicate.Group(sql.FieldNEQ(FieldExternalID, v))
}

// ExternalIDIn applies the In predicate on the "external_id" field.
func ExternalIDIn(vs ...string) predicate.Group {
	return predicate.Group(sql.FieldIn(FieldExternalID, vs...))
}

// ExternalIDNotIn applies the NotIn predicate on the "external_id" field.
func ExternalIDNotIn(vs ...string) predicate.Group {
	return predicate.Group(sql.FieldNotIn(FieldExternalID, vs...))
}

// ExternalIDGT applies the GT predicate on the "external_id" field.
func ExternalIDGT(v string) predicate.Group {
	return predicate.Group(sql.FieldGT(FieldExternalID, v))
}

// ExternalIDGTE applies the GTE predicate on the "external_id" field.
func ExternalIDGTE(v string) predicate.Group {
	return predicate.Group(sql.FieldGTE(FieldExternalID, v))
}

// ExternalIDLT applies the LT predicate on the "external_id" field.
func ExternalIDLT(v string) predicate.Group {
	return predicate.Group(sql.FieldLT(FieldExternalID, v))
}

// ExternalIDLTE applies the LTE predicate on the "external_id" field.
func ExternalIDLTE(v string) predicate.Group {
	return predicate.Group(sql.FieldLTE(FieldExternalID, v))
}

// ExternalIDContains applies the Contains predicate on the "external_id" field.
func ExternalIDContains(v string) predicate.Group {
	return predicate.Group(sql.FieldContains(FieldExternalID, v))
}

// ExternalIDHasPrefix applies the HasPrefix predicate on the "external_id" field.
func ExternalIDHasPrefix(v string) predicate.Group {
	return predicate.Group(sql.FieldHasPrefix(FieldExternalID, v))
}

// ExternalIDHasSuffix applies the HasSuffix predicate on the "external_id" field.
func ExternalIDHasSuffix(v string) predicate.Group {
	return predicate.Group(sql.FieldHasSuffix(FieldExternalID, v))
}

// ExternalIDIsNil applies the IsNil predicate on the "external_id" field.
func ExternalIDIsNil() predicate.Group {
	return predicate.Group(sql.FieldIsNull(FieldExternalID))
}

// ExternalIDNotNil applies the NotNil predicate on the "external_id" field.
func ExternalIDNotNil() predicate.Group {
	return predicate.Group(sql.FieldNotNull(FieldExternalID))
}

// ExternalIDEqualFold applies the EqualFold predicate on the "external_id" field.
func ExternalIDEqualFold(v string) predicate.Group {
	return predicate.Group(sql.FieldEqualFold(FieldExternalID, v))
}

// ExternalIDContainsFold applies the ContainsFold predicate on the "external_id" field.
func ExternalIDContainsFold(v string) predicate.Group {
	return predicate.Group(sql.FieldContainsFold(FieldExternalID, v))
}

// TeamIDEQ applies the EQ predicate on the "team_id" field.
func TeamIDEQ(v uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldEQ(FieldTeamID, v))
}

// TeamIDNEQ applies the NEQ predicate on the "team_id" field.
func TeamIDNEQ(v uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldNEQ(FieldTeamID, v))
}

// TeamIDIn applies the In predicate on the "team_id" field.
func TeamIDIn(vs ...uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldIn(FieldTeamID, vs...))
}

// TeamIDNotIn applies the NotIn predicate on the "team_id" field.
func TeamIDNotIn(vs ...uuid.UUID) predicate.Group {
	return predicate.Group(sql.FieldNotIn(FieldTeamID, vs...))
}

// TeamIDIsNil applies the IsNil predicate on the "team_id" field.
func TeamIDIsNil() predicate.Group {
	return predicate.Group(sql.FieldIsNull(FieldTeamID))
}

// TeamIDNotNil applies the NotNil predicate on the "team_id" field.
func TeamIDNotNil() predicate.Group {
	return predicate.Group(sql.FieldNotNull(FieldTeamID))
}

// HasUsers applies the HasEdge predicate on the "users" edge.
func HasUsers() predicate.Group {
	return predicate.Group(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, UsersTable, UsersPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUsersWith applies the HasEdge predicate on the "users" edge with a given conditions (other predicates).
func HasUsersWith(preds ...predicate.User) predicate.Group {
	return predicate.Group(func(s *sql.Selector) {
		step := newUsersStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPermissions applies the HasEdge predicate on the "permissions" edge.
func HasPermissions() predicate.Group {
	return predicate.Group(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, PermissionsTable, PermissionsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPermissionsWith applies the HasEdge predicate on the "permissions" edge with a given conditions (other predicates).
func HasPermissionsWith(preds ...predicate.Permission) predicate.Group {
	return predicate.Group(func(s *sql.Selector) {
		step := newPermissionsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTeam applies the HasEdge predicate on the "team" edge.
func HasTeam() predicate.Group {
	return predicate.Group(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, TeamTable, TeamColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTeamWith applies the HasEdge predicate on the "team" edge with a given conditions (other predicates).
func HasTeamWith(preds ...predicate.Team) predicate.Group {
	return predicate.Group(func(s *sql.Selector) {
		step := newTeamStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Group) predicate.Group {
	return predicate.Group(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Group) predicate.Group {
	return predicate.Group(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Group) predicate.Group {
	return predicate.Group(sql.NotPredicates(p))
}
