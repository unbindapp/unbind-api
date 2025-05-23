// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/unbindapp/unbind-api/ent/oauth2code"
	"github.com/unbindapp/unbind-api/ent/predicate"
)

// Oauth2CodeDelete is the builder for deleting a Oauth2Code entity.
type Oauth2CodeDelete struct {
	config
	hooks    []Hook
	mutation *Oauth2CodeMutation
}

// Where appends a list predicates to the Oauth2CodeDelete builder.
func (od *Oauth2CodeDelete) Where(ps ...predicate.Oauth2Code) *Oauth2CodeDelete {
	od.mutation.Where(ps...)
	return od
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (od *Oauth2CodeDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, od.sqlExec, od.mutation, od.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (od *Oauth2CodeDelete) ExecX(ctx context.Context) int {
	n, err := od.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (od *Oauth2CodeDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(oauth2code.Table, sqlgraph.NewFieldSpec(oauth2code.FieldID, field.TypeUUID))
	if ps := od.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, od.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	od.mutation.done = true
	return affected, err
}

// Oauth2CodeDeleteOne is the builder for deleting a single Oauth2Code entity.
type Oauth2CodeDeleteOne struct {
	od *Oauth2CodeDelete
}

// Where appends a list predicates to the Oauth2CodeDelete builder.
func (odo *Oauth2CodeDeleteOne) Where(ps ...predicate.Oauth2Code) *Oauth2CodeDeleteOne {
	odo.od.mutation.Where(ps...)
	return odo
}

// Exec executes the deletion query.
func (odo *Oauth2CodeDeleteOne) Exec(ctx context.Context) error {
	n, err := odo.od.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{oauth2code.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (odo *Oauth2CodeDeleteOne) ExecX(ctx context.Context) {
	if err := odo.Exec(ctx); err != nil {
		panic(err)
	}
}
