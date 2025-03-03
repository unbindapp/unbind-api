// Code generated by ent, DO NOT EDIT.

package githubinstallation

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/unbindapp/unbind-api/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldUpdatedAt, v))
}

// GithubAppID applies equality check predicate on the "github_app_id" field. It's identical to GithubAppIDEQ.
func GithubAppID(v int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldGithubAppID, v))
}

// AccountID applies equality check predicate on the "account_id" field. It's identical to AccountIDEQ.
func AccountID(v int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldAccountID, v))
}

// AccountLogin applies equality check predicate on the "account_login" field. It's identical to AccountLoginEQ.
func AccountLogin(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldAccountLogin, v))
}

// AccountURL applies equality check predicate on the "account_url" field. It's identical to AccountURLEQ.
func AccountURL(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldAccountURL, v))
}

// Suspended applies equality check predicate on the "suspended" field. It's identical to SuspendedEQ.
func Suspended(v bool) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldSuspended, v))
}

// Active applies equality check predicate on the "active" field. It's identical to ActiveEQ.
func Active(v bool) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldActive, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLTE(FieldUpdatedAt, v))
}

// GithubAppIDEQ applies the EQ predicate on the "github_app_id" field.
func GithubAppIDEQ(v int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldGithubAppID, v))
}

// GithubAppIDNEQ applies the NEQ predicate on the "github_app_id" field.
func GithubAppIDNEQ(v int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNEQ(FieldGithubAppID, v))
}

// GithubAppIDIn applies the In predicate on the "github_app_id" field.
func GithubAppIDIn(vs ...int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldIn(FieldGithubAppID, vs...))
}

// GithubAppIDNotIn applies the NotIn predicate on the "github_app_id" field.
func GithubAppIDNotIn(vs ...int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNotIn(FieldGithubAppID, vs...))
}

// AccountIDEQ applies the EQ predicate on the "account_id" field.
func AccountIDEQ(v int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldAccountID, v))
}

// AccountIDNEQ applies the NEQ predicate on the "account_id" field.
func AccountIDNEQ(v int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNEQ(FieldAccountID, v))
}

// AccountIDIn applies the In predicate on the "account_id" field.
func AccountIDIn(vs ...int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldIn(FieldAccountID, vs...))
}

// AccountIDNotIn applies the NotIn predicate on the "account_id" field.
func AccountIDNotIn(vs ...int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNotIn(FieldAccountID, vs...))
}

// AccountIDGT applies the GT predicate on the "account_id" field.
func AccountIDGT(v int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGT(FieldAccountID, v))
}

// AccountIDGTE applies the GTE predicate on the "account_id" field.
func AccountIDGTE(v int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGTE(FieldAccountID, v))
}

// AccountIDLT applies the LT predicate on the "account_id" field.
func AccountIDLT(v int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLT(FieldAccountID, v))
}

// AccountIDLTE applies the LTE predicate on the "account_id" field.
func AccountIDLTE(v int64) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLTE(FieldAccountID, v))
}

// AccountLoginEQ applies the EQ predicate on the "account_login" field.
func AccountLoginEQ(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldAccountLogin, v))
}

// AccountLoginNEQ applies the NEQ predicate on the "account_login" field.
func AccountLoginNEQ(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNEQ(FieldAccountLogin, v))
}

// AccountLoginIn applies the In predicate on the "account_login" field.
func AccountLoginIn(vs ...string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldIn(FieldAccountLogin, vs...))
}

// AccountLoginNotIn applies the NotIn predicate on the "account_login" field.
func AccountLoginNotIn(vs ...string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNotIn(FieldAccountLogin, vs...))
}

// AccountLoginGT applies the GT predicate on the "account_login" field.
func AccountLoginGT(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGT(FieldAccountLogin, v))
}

// AccountLoginGTE applies the GTE predicate on the "account_login" field.
func AccountLoginGTE(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGTE(FieldAccountLogin, v))
}

// AccountLoginLT applies the LT predicate on the "account_login" field.
func AccountLoginLT(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLT(FieldAccountLogin, v))
}

// AccountLoginLTE applies the LTE predicate on the "account_login" field.
func AccountLoginLTE(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLTE(FieldAccountLogin, v))
}

// AccountLoginContains applies the Contains predicate on the "account_login" field.
func AccountLoginContains(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldContains(FieldAccountLogin, v))
}

// AccountLoginHasPrefix applies the HasPrefix predicate on the "account_login" field.
func AccountLoginHasPrefix(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldHasPrefix(FieldAccountLogin, v))
}

// AccountLoginHasSuffix applies the HasSuffix predicate on the "account_login" field.
func AccountLoginHasSuffix(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldHasSuffix(FieldAccountLogin, v))
}

// AccountLoginEqualFold applies the EqualFold predicate on the "account_login" field.
func AccountLoginEqualFold(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEqualFold(FieldAccountLogin, v))
}

// AccountLoginContainsFold applies the ContainsFold predicate on the "account_login" field.
func AccountLoginContainsFold(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldContainsFold(FieldAccountLogin, v))
}

// AccountTypeEQ applies the EQ predicate on the "account_type" field.
func AccountTypeEQ(v AccountType) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldAccountType, v))
}

// AccountTypeNEQ applies the NEQ predicate on the "account_type" field.
func AccountTypeNEQ(v AccountType) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNEQ(FieldAccountType, v))
}

// AccountTypeIn applies the In predicate on the "account_type" field.
func AccountTypeIn(vs ...AccountType) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldIn(FieldAccountType, vs...))
}

// AccountTypeNotIn applies the NotIn predicate on the "account_type" field.
func AccountTypeNotIn(vs ...AccountType) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNotIn(FieldAccountType, vs...))
}

// AccountURLEQ applies the EQ predicate on the "account_url" field.
func AccountURLEQ(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldAccountURL, v))
}

// AccountURLNEQ applies the NEQ predicate on the "account_url" field.
func AccountURLNEQ(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNEQ(FieldAccountURL, v))
}

// AccountURLIn applies the In predicate on the "account_url" field.
func AccountURLIn(vs ...string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldIn(FieldAccountURL, vs...))
}

// AccountURLNotIn applies the NotIn predicate on the "account_url" field.
func AccountURLNotIn(vs ...string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNotIn(FieldAccountURL, vs...))
}

// AccountURLGT applies the GT predicate on the "account_url" field.
func AccountURLGT(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGT(FieldAccountURL, v))
}

// AccountURLGTE applies the GTE predicate on the "account_url" field.
func AccountURLGTE(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldGTE(FieldAccountURL, v))
}

// AccountURLLT applies the LT predicate on the "account_url" field.
func AccountURLLT(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLT(FieldAccountURL, v))
}

// AccountURLLTE applies the LTE predicate on the "account_url" field.
func AccountURLLTE(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldLTE(FieldAccountURL, v))
}

// AccountURLContains applies the Contains predicate on the "account_url" field.
func AccountURLContains(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldContains(FieldAccountURL, v))
}

// AccountURLHasPrefix applies the HasPrefix predicate on the "account_url" field.
func AccountURLHasPrefix(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldHasPrefix(FieldAccountURL, v))
}

// AccountURLHasSuffix applies the HasSuffix predicate on the "account_url" field.
func AccountURLHasSuffix(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldHasSuffix(FieldAccountURL, v))
}

// AccountURLEqualFold applies the EqualFold predicate on the "account_url" field.
func AccountURLEqualFold(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEqualFold(FieldAccountURL, v))
}

// AccountURLContainsFold applies the ContainsFold predicate on the "account_url" field.
func AccountURLContainsFold(v string) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldContainsFold(FieldAccountURL, v))
}

// RepositorySelectionEQ applies the EQ predicate on the "repository_selection" field.
func RepositorySelectionEQ(v RepositorySelection) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldRepositorySelection, v))
}

// RepositorySelectionNEQ applies the NEQ predicate on the "repository_selection" field.
func RepositorySelectionNEQ(v RepositorySelection) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNEQ(FieldRepositorySelection, v))
}

// RepositorySelectionIn applies the In predicate on the "repository_selection" field.
func RepositorySelectionIn(vs ...RepositorySelection) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldIn(FieldRepositorySelection, vs...))
}

// RepositorySelectionNotIn applies the NotIn predicate on the "repository_selection" field.
func RepositorySelectionNotIn(vs ...RepositorySelection) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNotIn(FieldRepositorySelection, vs...))
}

// SuspendedEQ applies the EQ predicate on the "suspended" field.
func SuspendedEQ(v bool) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldSuspended, v))
}

// SuspendedNEQ applies the NEQ predicate on the "suspended" field.
func SuspendedNEQ(v bool) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNEQ(FieldSuspended, v))
}

// ActiveEQ applies the EQ predicate on the "active" field.
func ActiveEQ(v bool) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldEQ(FieldActive, v))
}

// ActiveNEQ applies the NEQ predicate on the "active" field.
func ActiveNEQ(v bool) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNEQ(FieldActive, v))
}

// PermissionsIsNil applies the IsNil predicate on the "permissions" field.
func PermissionsIsNil() predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldIsNull(FieldPermissions))
}

// PermissionsNotNil applies the NotNil predicate on the "permissions" field.
func PermissionsNotNil() predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNotNull(FieldPermissions))
}

// EventsIsNil applies the IsNil predicate on the "events" field.
func EventsIsNil() predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldIsNull(FieldEvents))
}

// EventsNotNil applies the NotNil predicate on the "events" field.
func EventsNotNil() predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.FieldNotNull(FieldEvents))
}

// HasGithubApps applies the HasEdge predicate on the "github_apps" edge.
func HasGithubApps() predicate.GithubInstallation {
	return predicate.GithubInstallation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, GithubAppsTable, GithubAppsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasGithubAppsWith applies the HasEdge predicate on the "github_apps" edge with a given conditions (other predicates).
func HasGithubAppsWith(preds ...predicate.GithubApp) predicate.GithubInstallation {
	return predicate.GithubInstallation(func(s *sql.Selector) {
		step := newGithubAppsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.GithubInstallation) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.GithubInstallation) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.GithubInstallation) predicate.GithubInstallation {
	return predicate.GithubInstallation(sql.NotPredicates(p))
}
