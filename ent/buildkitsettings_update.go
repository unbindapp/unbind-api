// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/unbindapp/unbind-api/ent/buildkitsettings"
	"github.com/unbindapp/unbind-api/ent/predicate"
)

// BuildkitSettingsUpdate is the builder for updating BuildkitSettings entities.
type BuildkitSettingsUpdate struct {
	config
	hooks     []Hook
	mutation  *BuildkitSettingsMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the BuildkitSettingsUpdate builder.
func (bsu *BuildkitSettingsUpdate) Where(ps ...predicate.BuildkitSettings) *BuildkitSettingsUpdate {
	bsu.mutation.Where(ps...)
	return bsu
}

// SetUpdatedAt sets the "updated_at" field.
func (bsu *BuildkitSettingsUpdate) SetUpdatedAt(t time.Time) *BuildkitSettingsUpdate {
	bsu.mutation.SetUpdatedAt(t)
	return bsu
}

// SetMaxParallelism sets the "max_parallelism" field.
func (bsu *BuildkitSettingsUpdate) SetMaxParallelism(i int) *BuildkitSettingsUpdate {
	bsu.mutation.ResetMaxParallelism()
	bsu.mutation.SetMaxParallelism(i)
	return bsu
}

// SetNillableMaxParallelism sets the "max_parallelism" field if the given value is not nil.
func (bsu *BuildkitSettingsUpdate) SetNillableMaxParallelism(i *int) *BuildkitSettingsUpdate {
	if i != nil {
		bsu.SetMaxParallelism(*i)
	}
	return bsu
}

// AddMaxParallelism adds i to the "max_parallelism" field.
func (bsu *BuildkitSettingsUpdate) AddMaxParallelism(i int) *BuildkitSettingsUpdate {
	bsu.mutation.AddMaxParallelism(i)
	return bsu
}

// SetReplicas sets the "replicas" field.
func (bsu *BuildkitSettingsUpdate) SetReplicas(i int) *BuildkitSettingsUpdate {
	bsu.mutation.ResetReplicas()
	bsu.mutation.SetReplicas(i)
	return bsu
}

// SetNillableReplicas sets the "replicas" field if the given value is not nil.
func (bsu *BuildkitSettingsUpdate) SetNillableReplicas(i *int) *BuildkitSettingsUpdate {
	if i != nil {
		bsu.SetReplicas(*i)
	}
	return bsu
}

// AddReplicas adds i to the "replicas" field.
func (bsu *BuildkitSettingsUpdate) AddReplicas(i int) *BuildkitSettingsUpdate {
	bsu.mutation.AddReplicas(i)
	return bsu
}

// Mutation returns the BuildkitSettingsMutation object of the builder.
func (bsu *BuildkitSettingsUpdate) Mutation() *BuildkitSettingsMutation {
	return bsu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (bsu *BuildkitSettingsUpdate) Save(ctx context.Context) (int, error) {
	bsu.defaults()
	return withHooks(ctx, bsu.sqlSave, bsu.mutation, bsu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (bsu *BuildkitSettingsUpdate) SaveX(ctx context.Context) int {
	affected, err := bsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (bsu *BuildkitSettingsUpdate) Exec(ctx context.Context) error {
	_, err := bsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bsu *BuildkitSettingsUpdate) ExecX(ctx context.Context) {
	if err := bsu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (bsu *BuildkitSettingsUpdate) defaults() {
	if _, ok := bsu.mutation.UpdatedAt(); !ok {
		v := buildkitsettings.UpdateDefaultUpdatedAt()
		bsu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bsu *BuildkitSettingsUpdate) check() error {
	if v, ok := bsu.mutation.MaxParallelism(); ok {
		if err := buildkitsettings.MaxParallelismValidator(v); err != nil {
			return &ValidationError{Name: "max_parallelism", err: fmt.Errorf(`ent: validator failed for field "BuildkitSettings.max_parallelism": %w`, err)}
		}
	}
	if v, ok := bsu.mutation.Replicas(); ok {
		if err := buildkitsettings.ReplicasValidator(v); err != nil {
			return &ValidationError{Name: "replicas", err: fmt.Errorf(`ent: validator failed for field "BuildkitSettings.replicas": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (bsu *BuildkitSettingsUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *BuildkitSettingsUpdate {
	bsu.modifiers = append(bsu.modifiers, modifiers...)
	return bsu
}

func (bsu *BuildkitSettingsUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := bsu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(buildkitsettings.Table, buildkitsettings.Columns, sqlgraph.NewFieldSpec(buildkitsettings.FieldID, field.TypeUUID))
	if ps := bsu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := bsu.mutation.UpdatedAt(); ok {
		_spec.SetField(buildkitsettings.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := bsu.mutation.MaxParallelism(); ok {
		_spec.SetField(buildkitsettings.FieldMaxParallelism, field.TypeInt, value)
	}
	if value, ok := bsu.mutation.AddedMaxParallelism(); ok {
		_spec.AddField(buildkitsettings.FieldMaxParallelism, field.TypeInt, value)
	}
	if value, ok := bsu.mutation.Replicas(); ok {
		_spec.SetField(buildkitsettings.FieldReplicas, field.TypeInt, value)
	}
	if value, ok := bsu.mutation.AddedReplicas(); ok {
		_spec.AddField(buildkitsettings.FieldReplicas, field.TypeInt, value)
	}
	_spec.AddModifiers(bsu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, bsu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{buildkitsettings.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	bsu.mutation.done = true
	return n, nil
}

// BuildkitSettingsUpdateOne is the builder for updating a single BuildkitSettings entity.
type BuildkitSettingsUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *BuildkitSettingsMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (bsuo *BuildkitSettingsUpdateOne) SetUpdatedAt(t time.Time) *BuildkitSettingsUpdateOne {
	bsuo.mutation.SetUpdatedAt(t)
	return bsuo
}

// SetMaxParallelism sets the "max_parallelism" field.
func (bsuo *BuildkitSettingsUpdateOne) SetMaxParallelism(i int) *BuildkitSettingsUpdateOne {
	bsuo.mutation.ResetMaxParallelism()
	bsuo.mutation.SetMaxParallelism(i)
	return bsuo
}

// SetNillableMaxParallelism sets the "max_parallelism" field if the given value is not nil.
func (bsuo *BuildkitSettingsUpdateOne) SetNillableMaxParallelism(i *int) *BuildkitSettingsUpdateOne {
	if i != nil {
		bsuo.SetMaxParallelism(*i)
	}
	return bsuo
}

// AddMaxParallelism adds i to the "max_parallelism" field.
func (bsuo *BuildkitSettingsUpdateOne) AddMaxParallelism(i int) *BuildkitSettingsUpdateOne {
	bsuo.mutation.AddMaxParallelism(i)
	return bsuo
}

// SetReplicas sets the "replicas" field.
func (bsuo *BuildkitSettingsUpdateOne) SetReplicas(i int) *BuildkitSettingsUpdateOne {
	bsuo.mutation.ResetReplicas()
	bsuo.mutation.SetReplicas(i)
	return bsuo
}

// SetNillableReplicas sets the "replicas" field if the given value is not nil.
func (bsuo *BuildkitSettingsUpdateOne) SetNillableReplicas(i *int) *BuildkitSettingsUpdateOne {
	if i != nil {
		bsuo.SetReplicas(*i)
	}
	return bsuo
}

// AddReplicas adds i to the "replicas" field.
func (bsuo *BuildkitSettingsUpdateOne) AddReplicas(i int) *BuildkitSettingsUpdateOne {
	bsuo.mutation.AddReplicas(i)
	return bsuo
}

// Mutation returns the BuildkitSettingsMutation object of the builder.
func (bsuo *BuildkitSettingsUpdateOne) Mutation() *BuildkitSettingsMutation {
	return bsuo.mutation
}

// Where appends a list predicates to the BuildkitSettingsUpdate builder.
func (bsuo *BuildkitSettingsUpdateOne) Where(ps ...predicate.BuildkitSettings) *BuildkitSettingsUpdateOne {
	bsuo.mutation.Where(ps...)
	return bsuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (bsuo *BuildkitSettingsUpdateOne) Select(field string, fields ...string) *BuildkitSettingsUpdateOne {
	bsuo.fields = append([]string{field}, fields...)
	return bsuo
}

// Save executes the query and returns the updated BuildkitSettings entity.
func (bsuo *BuildkitSettingsUpdateOne) Save(ctx context.Context) (*BuildkitSettings, error) {
	bsuo.defaults()
	return withHooks(ctx, bsuo.sqlSave, bsuo.mutation, bsuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (bsuo *BuildkitSettingsUpdateOne) SaveX(ctx context.Context) *BuildkitSettings {
	node, err := bsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (bsuo *BuildkitSettingsUpdateOne) Exec(ctx context.Context) error {
	_, err := bsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bsuo *BuildkitSettingsUpdateOne) ExecX(ctx context.Context) {
	if err := bsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (bsuo *BuildkitSettingsUpdateOne) defaults() {
	if _, ok := bsuo.mutation.UpdatedAt(); !ok {
		v := buildkitsettings.UpdateDefaultUpdatedAt()
		bsuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bsuo *BuildkitSettingsUpdateOne) check() error {
	if v, ok := bsuo.mutation.MaxParallelism(); ok {
		if err := buildkitsettings.MaxParallelismValidator(v); err != nil {
			return &ValidationError{Name: "max_parallelism", err: fmt.Errorf(`ent: validator failed for field "BuildkitSettings.max_parallelism": %w`, err)}
		}
	}
	if v, ok := bsuo.mutation.Replicas(); ok {
		if err := buildkitsettings.ReplicasValidator(v); err != nil {
			return &ValidationError{Name: "replicas", err: fmt.Errorf(`ent: validator failed for field "BuildkitSettings.replicas": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (bsuo *BuildkitSettingsUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *BuildkitSettingsUpdateOne {
	bsuo.modifiers = append(bsuo.modifiers, modifiers...)
	return bsuo
}

func (bsuo *BuildkitSettingsUpdateOne) sqlSave(ctx context.Context) (_node *BuildkitSettings, err error) {
	if err := bsuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(buildkitsettings.Table, buildkitsettings.Columns, sqlgraph.NewFieldSpec(buildkitsettings.FieldID, field.TypeUUID))
	id, ok := bsuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "BuildkitSettings.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := bsuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, buildkitsettings.FieldID)
		for _, f := range fields {
			if !buildkitsettings.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != buildkitsettings.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := bsuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := bsuo.mutation.UpdatedAt(); ok {
		_spec.SetField(buildkitsettings.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := bsuo.mutation.MaxParallelism(); ok {
		_spec.SetField(buildkitsettings.FieldMaxParallelism, field.TypeInt, value)
	}
	if value, ok := bsuo.mutation.AddedMaxParallelism(); ok {
		_spec.AddField(buildkitsettings.FieldMaxParallelism, field.TypeInt, value)
	}
	if value, ok := bsuo.mutation.Replicas(); ok {
		_spec.SetField(buildkitsettings.FieldReplicas, field.TypeInt, value)
	}
	if value, ok := bsuo.mutation.AddedReplicas(); ok {
		_spec.AddField(buildkitsettings.FieldReplicas, field.TypeInt, value)
	}
	_spec.AddModifiers(bsuo.modifiers...)
	_node = &BuildkitSettings{config: bsuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, bsuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{buildkitsettings.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	bsuo.mutation.done = true
	return _node, nil
}
