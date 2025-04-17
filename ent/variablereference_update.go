// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/predicate"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/variablereference"
)

// VariableReferenceUpdate is the builder for updating VariableReference entities.
type VariableReferenceUpdate struct {
	config
	hooks     []Hook
	mutation  *VariableReferenceMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the VariableReferenceUpdate builder.
func (vru *VariableReferenceUpdate) Where(ps ...predicate.VariableReference) *VariableReferenceUpdate {
	vru.mutation.Where(ps...)
	return vru
}

// SetUpdatedAt sets the "updated_at" field.
func (vru *VariableReferenceUpdate) SetUpdatedAt(t time.Time) *VariableReferenceUpdate {
	vru.mutation.SetUpdatedAt(t)
	return vru
}

// SetTargetServiceID sets the "target_service_id" field.
func (vru *VariableReferenceUpdate) SetTargetServiceID(u uuid.UUID) *VariableReferenceUpdate {
	vru.mutation.SetTargetServiceID(u)
	return vru
}

// SetNillableTargetServiceID sets the "target_service_id" field if the given value is not nil.
func (vru *VariableReferenceUpdate) SetNillableTargetServiceID(u *uuid.UUID) *VariableReferenceUpdate {
	if u != nil {
		vru.SetTargetServiceID(*u)
	}
	return vru
}

// SetTargetName sets the "target_name" field.
func (vru *VariableReferenceUpdate) SetTargetName(s string) *VariableReferenceUpdate {
	vru.mutation.SetTargetName(s)
	return vru
}

// SetNillableTargetName sets the "target_name" field if the given value is not nil.
func (vru *VariableReferenceUpdate) SetNillableTargetName(s *string) *VariableReferenceUpdate {
	if s != nil {
		vru.SetTargetName(*s)
	}
	return vru
}

// SetSources sets the "sources" field.
func (vru *VariableReferenceUpdate) SetSources(srs []schema.VariableReferenceSource) *VariableReferenceUpdate {
	vru.mutation.SetSources(srs)
	return vru
}

// AppendSources appends srs to the "sources" field.
func (vru *VariableReferenceUpdate) AppendSources(srs []schema.VariableReferenceSource) *VariableReferenceUpdate {
	vru.mutation.AppendSources(srs)
	return vru
}

// SetValueTemplate sets the "value_template" field.
func (vru *VariableReferenceUpdate) SetValueTemplate(s string) *VariableReferenceUpdate {
	vru.mutation.SetValueTemplate(s)
	return vru
}

// SetNillableValueTemplate sets the "value_template" field if the given value is not nil.
func (vru *VariableReferenceUpdate) SetNillableValueTemplate(s *string) *VariableReferenceUpdate {
	if s != nil {
		vru.SetValueTemplate(*s)
	}
	return vru
}

// SetError sets the "error" field.
func (vru *VariableReferenceUpdate) SetError(s string) *VariableReferenceUpdate {
	vru.mutation.SetError(s)
	return vru
}

// SetNillableError sets the "error" field if the given value is not nil.
func (vru *VariableReferenceUpdate) SetNillableError(s *string) *VariableReferenceUpdate {
	if s != nil {
		vru.SetError(*s)
	}
	return vru
}

// ClearError clears the value of the "error" field.
func (vru *VariableReferenceUpdate) ClearError() *VariableReferenceUpdate {
	vru.mutation.ClearError()
	return vru
}

// SetServiceID sets the "service" edge to the Service entity by ID.
func (vru *VariableReferenceUpdate) SetServiceID(id uuid.UUID) *VariableReferenceUpdate {
	vru.mutation.SetServiceID(id)
	return vru
}

// SetService sets the "service" edge to the Service entity.
func (vru *VariableReferenceUpdate) SetService(s *Service) *VariableReferenceUpdate {
	return vru.SetServiceID(s.ID)
}

// Mutation returns the VariableReferenceMutation object of the builder.
func (vru *VariableReferenceUpdate) Mutation() *VariableReferenceMutation {
	return vru.mutation
}

// ClearService clears the "service" edge to the Service entity.
func (vru *VariableReferenceUpdate) ClearService() *VariableReferenceUpdate {
	vru.mutation.ClearService()
	return vru
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (vru *VariableReferenceUpdate) Save(ctx context.Context) (int, error) {
	vru.defaults()
	return withHooks(ctx, vru.sqlSave, vru.mutation, vru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (vru *VariableReferenceUpdate) SaveX(ctx context.Context) int {
	affected, err := vru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (vru *VariableReferenceUpdate) Exec(ctx context.Context) error {
	_, err := vru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vru *VariableReferenceUpdate) ExecX(ctx context.Context) {
	if err := vru.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (vru *VariableReferenceUpdate) defaults() {
	if _, ok := vru.mutation.UpdatedAt(); !ok {
		v := variablereference.UpdateDefaultUpdatedAt()
		vru.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (vru *VariableReferenceUpdate) check() error {
	if vru.mutation.ServiceCleared() && len(vru.mutation.ServiceIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "VariableReference.service"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (vru *VariableReferenceUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *VariableReferenceUpdate {
	vru.modifiers = append(vru.modifiers, modifiers...)
	return vru
}

func (vru *VariableReferenceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := vru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(variablereference.Table, variablereference.Columns, sqlgraph.NewFieldSpec(variablereference.FieldID, field.TypeUUID))
	if ps := vru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := vru.mutation.UpdatedAt(); ok {
		_spec.SetField(variablereference.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := vru.mutation.TargetName(); ok {
		_spec.SetField(variablereference.FieldTargetName, field.TypeString, value)
	}
	if value, ok := vru.mutation.Sources(); ok {
		_spec.SetField(variablereference.FieldSources, field.TypeJSON, value)
	}
	if value, ok := vru.mutation.AppendedSources(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, variablereference.FieldSources, value)
		})
	}
	if value, ok := vru.mutation.ValueTemplate(); ok {
		_spec.SetField(variablereference.FieldValueTemplate, field.TypeString, value)
	}
	if value, ok := vru.mutation.Error(); ok {
		_spec.SetField(variablereference.FieldError, field.TypeString, value)
	}
	if vru.mutation.ErrorCleared() {
		_spec.ClearField(variablereference.FieldError, field.TypeString)
	}
	if vru.mutation.ServiceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   variablereference.ServiceTable,
			Columns: []string{variablereference.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vru.mutation.ServiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   variablereference.ServiceTable,
			Columns: []string{variablereference.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(vru.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, vru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{variablereference.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	vru.mutation.done = true
	return n, nil
}

// VariableReferenceUpdateOne is the builder for updating a single VariableReference entity.
type VariableReferenceUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *VariableReferenceMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (vruo *VariableReferenceUpdateOne) SetUpdatedAt(t time.Time) *VariableReferenceUpdateOne {
	vruo.mutation.SetUpdatedAt(t)
	return vruo
}

// SetTargetServiceID sets the "target_service_id" field.
func (vruo *VariableReferenceUpdateOne) SetTargetServiceID(u uuid.UUID) *VariableReferenceUpdateOne {
	vruo.mutation.SetTargetServiceID(u)
	return vruo
}

// SetNillableTargetServiceID sets the "target_service_id" field if the given value is not nil.
func (vruo *VariableReferenceUpdateOne) SetNillableTargetServiceID(u *uuid.UUID) *VariableReferenceUpdateOne {
	if u != nil {
		vruo.SetTargetServiceID(*u)
	}
	return vruo
}

// SetTargetName sets the "target_name" field.
func (vruo *VariableReferenceUpdateOne) SetTargetName(s string) *VariableReferenceUpdateOne {
	vruo.mutation.SetTargetName(s)
	return vruo
}

// SetNillableTargetName sets the "target_name" field if the given value is not nil.
func (vruo *VariableReferenceUpdateOne) SetNillableTargetName(s *string) *VariableReferenceUpdateOne {
	if s != nil {
		vruo.SetTargetName(*s)
	}
	return vruo
}

// SetSources sets the "sources" field.
func (vruo *VariableReferenceUpdateOne) SetSources(srs []schema.VariableReferenceSource) *VariableReferenceUpdateOne {
	vruo.mutation.SetSources(srs)
	return vruo
}

// AppendSources appends srs to the "sources" field.
func (vruo *VariableReferenceUpdateOne) AppendSources(srs []schema.VariableReferenceSource) *VariableReferenceUpdateOne {
	vruo.mutation.AppendSources(srs)
	return vruo
}

// SetValueTemplate sets the "value_template" field.
func (vruo *VariableReferenceUpdateOne) SetValueTemplate(s string) *VariableReferenceUpdateOne {
	vruo.mutation.SetValueTemplate(s)
	return vruo
}

// SetNillableValueTemplate sets the "value_template" field if the given value is not nil.
func (vruo *VariableReferenceUpdateOne) SetNillableValueTemplate(s *string) *VariableReferenceUpdateOne {
	if s != nil {
		vruo.SetValueTemplate(*s)
	}
	return vruo
}

// SetError sets the "error" field.
func (vruo *VariableReferenceUpdateOne) SetError(s string) *VariableReferenceUpdateOne {
	vruo.mutation.SetError(s)
	return vruo
}

// SetNillableError sets the "error" field if the given value is not nil.
func (vruo *VariableReferenceUpdateOne) SetNillableError(s *string) *VariableReferenceUpdateOne {
	if s != nil {
		vruo.SetError(*s)
	}
	return vruo
}

// ClearError clears the value of the "error" field.
func (vruo *VariableReferenceUpdateOne) ClearError() *VariableReferenceUpdateOne {
	vruo.mutation.ClearError()
	return vruo
}

// SetServiceID sets the "service" edge to the Service entity by ID.
func (vruo *VariableReferenceUpdateOne) SetServiceID(id uuid.UUID) *VariableReferenceUpdateOne {
	vruo.mutation.SetServiceID(id)
	return vruo
}

// SetService sets the "service" edge to the Service entity.
func (vruo *VariableReferenceUpdateOne) SetService(s *Service) *VariableReferenceUpdateOne {
	return vruo.SetServiceID(s.ID)
}

// Mutation returns the VariableReferenceMutation object of the builder.
func (vruo *VariableReferenceUpdateOne) Mutation() *VariableReferenceMutation {
	return vruo.mutation
}

// ClearService clears the "service" edge to the Service entity.
func (vruo *VariableReferenceUpdateOne) ClearService() *VariableReferenceUpdateOne {
	vruo.mutation.ClearService()
	return vruo
}

// Where appends a list predicates to the VariableReferenceUpdate builder.
func (vruo *VariableReferenceUpdateOne) Where(ps ...predicate.VariableReference) *VariableReferenceUpdateOne {
	vruo.mutation.Where(ps...)
	return vruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (vruo *VariableReferenceUpdateOne) Select(field string, fields ...string) *VariableReferenceUpdateOne {
	vruo.fields = append([]string{field}, fields...)
	return vruo
}

// Save executes the query and returns the updated VariableReference entity.
func (vruo *VariableReferenceUpdateOne) Save(ctx context.Context) (*VariableReference, error) {
	vruo.defaults()
	return withHooks(ctx, vruo.sqlSave, vruo.mutation, vruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (vruo *VariableReferenceUpdateOne) SaveX(ctx context.Context) *VariableReference {
	node, err := vruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (vruo *VariableReferenceUpdateOne) Exec(ctx context.Context) error {
	_, err := vruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vruo *VariableReferenceUpdateOne) ExecX(ctx context.Context) {
	if err := vruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (vruo *VariableReferenceUpdateOne) defaults() {
	if _, ok := vruo.mutation.UpdatedAt(); !ok {
		v := variablereference.UpdateDefaultUpdatedAt()
		vruo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (vruo *VariableReferenceUpdateOne) check() error {
	if vruo.mutation.ServiceCleared() && len(vruo.mutation.ServiceIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "VariableReference.service"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (vruo *VariableReferenceUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *VariableReferenceUpdateOne {
	vruo.modifiers = append(vruo.modifiers, modifiers...)
	return vruo
}

func (vruo *VariableReferenceUpdateOne) sqlSave(ctx context.Context) (_node *VariableReference, err error) {
	if err := vruo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(variablereference.Table, variablereference.Columns, sqlgraph.NewFieldSpec(variablereference.FieldID, field.TypeUUID))
	id, ok := vruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "VariableReference.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := vruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, variablereference.FieldID)
		for _, f := range fields {
			if !variablereference.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != variablereference.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := vruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := vruo.mutation.UpdatedAt(); ok {
		_spec.SetField(variablereference.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := vruo.mutation.TargetName(); ok {
		_spec.SetField(variablereference.FieldTargetName, field.TypeString, value)
	}
	if value, ok := vruo.mutation.Sources(); ok {
		_spec.SetField(variablereference.FieldSources, field.TypeJSON, value)
	}
	if value, ok := vruo.mutation.AppendedSources(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, variablereference.FieldSources, value)
		})
	}
	if value, ok := vruo.mutation.ValueTemplate(); ok {
		_spec.SetField(variablereference.FieldValueTemplate, field.TypeString, value)
	}
	if value, ok := vruo.mutation.Error(); ok {
		_spec.SetField(variablereference.FieldError, field.TypeString, value)
	}
	if vruo.mutation.ErrorCleared() {
		_spec.ClearField(variablereference.FieldError, field.TypeString)
	}
	if vruo.mutation.ServiceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   variablereference.ServiceTable,
			Columns: []string{variablereference.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vruo.mutation.ServiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   variablereference.ServiceTable,
			Columns: []string{variablereference.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(vruo.modifiers...)
	_node = &VariableReference{config: vruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, vruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{variablereference.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	vruo.mutation.done = true
	return _node, nil
}
