// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/unbindapp/unbind-api/ent/predicate"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
)

// ServiceConfigDelete is the builder for deleting a ServiceConfig entity.
type ServiceConfigDelete struct {
	config
	hooks    []Hook
	mutation *ServiceConfigMutation
}

// Where appends a list predicates to the ServiceConfigDelete builder.
func (scd *ServiceConfigDelete) Where(ps ...predicate.ServiceConfig) *ServiceConfigDelete {
	scd.mutation.Where(ps...)
	return scd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (scd *ServiceConfigDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, scd.sqlExec, scd.mutation, scd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (scd *ServiceConfigDelete) ExecX(ctx context.Context) int {
	n, err := scd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (scd *ServiceConfigDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(serviceconfig.Table, sqlgraph.NewFieldSpec(serviceconfig.FieldID, field.TypeUUID))
	if ps := scd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, scd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	scd.mutation.done = true
	return affected, err
}

// ServiceConfigDeleteOne is the builder for deleting a single ServiceConfig entity.
type ServiceConfigDeleteOne struct {
	scd *ServiceConfigDelete
}

// Where appends a list predicates to the ServiceConfigDelete builder.
func (scdo *ServiceConfigDeleteOne) Where(ps ...predicate.ServiceConfig) *ServiceConfigDeleteOne {
	scdo.scd.mutation.Where(ps...)
	return scdo
}

// Exec executes the deletion query.
func (scdo *ServiceConfigDeleteOne) Exec(ctx context.Context) error {
	n, err := scdo.scd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{serviceconfig.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (scdo *ServiceConfigDeleteOne) ExecX(ctx context.Context) {
	if err := scdo.Exec(ctx); err != nil {
		panic(err)
	}
}
