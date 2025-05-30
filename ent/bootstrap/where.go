// Code generated by ent, DO NOT EDIT.

package bootstrap

import (
	"entgo.io/ent/dialect/sql"
	"github.com/unbindapp/unbind-api/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldLTE(FieldID, id))
}

// IsBootstrapped applies equality check predicate on the "is_bootstrapped" field. It's identical to IsBootstrappedEQ.
func IsBootstrapped(v bool) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldEQ(FieldIsBootstrapped, v))
}

// IsBootstrappedEQ applies the EQ predicate on the "is_bootstrapped" field.
func IsBootstrappedEQ(v bool) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldEQ(FieldIsBootstrapped, v))
}

// IsBootstrappedNEQ applies the NEQ predicate on the "is_bootstrapped" field.
func IsBootstrappedNEQ(v bool) predicate.Bootstrap {
	return predicate.Bootstrap(sql.FieldNEQ(FieldIsBootstrapped, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Bootstrap) predicate.Bootstrap {
	return predicate.Bootstrap(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Bootstrap) predicate.Bootstrap {
	return predicate.Bootstrap(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Bootstrap) predicate.Bootstrap {
	return predicate.Bootstrap(sql.NotPredicates(p))
}
