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
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/environment"
	"github.com/unbindapp/unbind-api/ent/predicate"
	"github.com/unbindapp/unbind-api/ent/project"
	"github.com/unbindapp/unbind-api/ent/service"
)

// EnvironmentUpdate is the builder for updating Environment entities.
type EnvironmentUpdate struct {
	config
	hooks     []Hook
	mutation  *EnvironmentMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the EnvironmentUpdate builder.
func (eu *EnvironmentUpdate) Where(ps ...predicate.Environment) *EnvironmentUpdate {
	eu.mutation.Where(ps...)
	return eu
}

// SetUpdatedAt sets the "updated_at" field.
func (eu *EnvironmentUpdate) SetUpdatedAt(t time.Time) *EnvironmentUpdate {
	eu.mutation.SetUpdatedAt(t)
	return eu
}

// SetName sets the "name" field.
func (eu *EnvironmentUpdate) SetName(s string) *EnvironmentUpdate {
	eu.mutation.SetName(s)
	return eu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (eu *EnvironmentUpdate) SetNillableName(s *string) *EnvironmentUpdate {
	if s != nil {
		eu.SetName(*s)
	}
	return eu
}

// SetDisplayName sets the "display_name" field.
func (eu *EnvironmentUpdate) SetDisplayName(s string) *EnvironmentUpdate {
	eu.mutation.SetDisplayName(s)
	return eu
}

// SetNillableDisplayName sets the "display_name" field if the given value is not nil.
func (eu *EnvironmentUpdate) SetNillableDisplayName(s *string) *EnvironmentUpdate {
	if s != nil {
		eu.SetDisplayName(*s)
	}
	return eu
}

// SetDescription sets the "description" field.
func (eu *EnvironmentUpdate) SetDescription(s string) *EnvironmentUpdate {
	eu.mutation.SetDescription(s)
	return eu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (eu *EnvironmentUpdate) SetNillableDescription(s *string) *EnvironmentUpdate {
	if s != nil {
		eu.SetDescription(*s)
	}
	return eu
}

// SetActive sets the "active" field.
func (eu *EnvironmentUpdate) SetActive(b bool) *EnvironmentUpdate {
	eu.mutation.SetActive(b)
	return eu
}

// SetNillableActive sets the "active" field if the given value is not nil.
func (eu *EnvironmentUpdate) SetNillableActive(b *bool) *EnvironmentUpdate {
	if b != nil {
		eu.SetActive(*b)
	}
	return eu
}

// SetProjectID sets the "project_id" field.
func (eu *EnvironmentUpdate) SetProjectID(u uuid.UUID) *EnvironmentUpdate {
	eu.mutation.SetProjectID(u)
	return eu
}

// SetNillableProjectID sets the "project_id" field if the given value is not nil.
func (eu *EnvironmentUpdate) SetNillableProjectID(u *uuid.UUID) *EnvironmentUpdate {
	if u != nil {
		eu.SetProjectID(*u)
	}
	return eu
}

// SetProject sets the "project" edge to the Project entity.
func (eu *EnvironmentUpdate) SetProject(p *Project) *EnvironmentUpdate {
	return eu.SetProjectID(p.ID)
}

// AddServiceIDs adds the "services" edge to the Service entity by IDs.
func (eu *EnvironmentUpdate) AddServiceIDs(ids ...uuid.UUID) *EnvironmentUpdate {
	eu.mutation.AddServiceIDs(ids...)
	return eu
}

// AddServices adds the "services" edges to the Service entity.
func (eu *EnvironmentUpdate) AddServices(s ...*Service) *EnvironmentUpdate {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return eu.AddServiceIDs(ids...)
}

// Mutation returns the EnvironmentMutation object of the builder.
func (eu *EnvironmentUpdate) Mutation() *EnvironmentMutation {
	return eu.mutation
}

// ClearProject clears the "project" edge to the Project entity.
func (eu *EnvironmentUpdate) ClearProject() *EnvironmentUpdate {
	eu.mutation.ClearProject()
	return eu
}

// ClearServices clears all "services" edges to the Service entity.
func (eu *EnvironmentUpdate) ClearServices() *EnvironmentUpdate {
	eu.mutation.ClearServices()
	return eu
}

// RemoveServiceIDs removes the "services" edge to Service entities by IDs.
func (eu *EnvironmentUpdate) RemoveServiceIDs(ids ...uuid.UUID) *EnvironmentUpdate {
	eu.mutation.RemoveServiceIDs(ids...)
	return eu
}

// RemoveServices removes "services" edges to Service entities.
func (eu *EnvironmentUpdate) RemoveServices(s ...*Service) *EnvironmentUpdate {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return eu.RemoveServiceIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (eu *EnvironmentUpdate) Save(ctx context.Context) (int, error) {
	eu.defaults()
	return withHooks(ctx, eu.sqlSave, eu.mutation, eu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (eu *EnvironmentUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *EnvironmentUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *EnvironmentUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (eu *EnvironmentUpdate) defaults() {
	if _, ok := eu.mutation.UpdatedAt(); !ok {
		v := environment.UpdateDefaultUpdatedAt()
		eu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (eu *EnvironmentUpdate) check() error {
	if v, ok := eu.mutation.Name(); ok {
		if err := environment.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Environment.name": %w`, err)}
		}
	}
	if eu.mutation.ProjectCleared() && len(eu.mutation.ProjectIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Environment.project"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (eu *EnvironmentUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *EnvironmentUpdate {
	eu.modifiers = append(eu.modifiers, modifiers...)
	return eu
}

func (eu *EnvironmentUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := eu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(environment.Table, environment.Columns, sqlgraph.NewFieldSpec(environment.FieldID, field.TypeUUID))
	if ps := eu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eu.mutation.UpdatedAt(); ok {
		_spec.SetField(environment.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := eu.mutation.Name(); ok {
		_spec.SetField(environment.FieldName, field.TypeString, value)
	}
	if value, ok := eu.mutation.DisplayName(); ok {
		_spec.SetField(environment.FieldDisplayName, field.TypeString, value)
	}
	if value, ok := eu.mutation.Description(); ok {
		_spec.SetField(environment.FieldDescription, field.TypeString, value)
	}
	if value, ok := eu.mutation.Active(); ok {
		_spec.SetField(environment.FieldActive, field.TypeBool, value)
	}
	if eu.mutation.ProjectCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   environment.ProjectTable,
			Columns: []string{environment.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   environment.ProjectTable,
			Columns: []string{environment.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.mutation.ServicesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ServicesTable,
			Columns: []string{environment.ServicesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.RemovedServicesIDs(); len(nodes) > 0 && !eu.mutation.ServicesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ServicesTable,
			Columns: []string{environment.ServicesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.ServicesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ServicesTable,
			Columns: []string{environment.ServicesColumn},
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
	_spec.AddModifiers(eu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, eu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{environment.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	eu.mutation.done = true
	return n, nil
}

// EnvironmentUpdateOne is the builder for updating a single Environment entity.
type EnvironmentUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *EnvironmentMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (euo *EnvironmentUpdateOne) SetUpdatedAt(t time.Time) *EnvironmentUpdateOne {
	euo.mutation.SetUpdatedAt(t)
	return euo
}

// SetName sets the "name" field.
func (euo *EnvironmentUpdateOne) SetName(s string) *EnvironmentUpdateOne {
	euo.mutation.SetName(s)
	return euo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (euo *EnvironmentUpdateOne) SetNillableName(s *string) *EnvironmentUpdateOne {
	if s != nil {
		euo.SetName(*s)
	}
	return euo
}

// SetDisplayName sets the "display_name" field.
func (euo *EnvironmentUpdateOne) SetDisplayName(s string) *EnvironmentUpdateOne {
	euo.mutation.SetDisplayName(s)
	return euo
}

// SetNillableDisplayName sets the "display_name" field if the given value is not nil.
func (euo *EnvironmentUpdateOne) SetNillableDisplayName(s *string) *EnvironmentUpdateOne {
	if s != nil {
		euo.SetDisplayName(*s)
	}
	return euo
}

// SetDescription sets the "description" field.
func (euo *EnvironmentUpdateOne) SetDescription(s string) *EnvironmentUpdateOne {
	euo.mutation.SetDescription(s)
	return euo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (euo *EnvironmentUpdateOne) SetNillableDescription(s *string) *EnvironmentUpdateOne {
	if s != nil {
		euo.SetDescription(*s)
	}
	return euo
}

// SetActive sets the "active" field.
func (euo *EnvironmentUpdateOne) SetActive(b bool) *EnvironmentUpdateOne {
	euo.mutation.SetActive(b)
	return euo
}

// SetNillableActive sets the "active" field if the given value is not nil.
func (euo *EnvironmentUpdateOne) SetNillableActive(b *bool) *EnvironmentUpdateOne {
	if b != nil {
		euo.SetActive(*b)
	}
	return euo
}

// SetProjectID sets the "project_id" field.
func (euo *EnvironmentUpdateOne) SetProjectID(u uuid.UUID) *EnvironmentUpdateOne {
	euo.mutation.SetProjectID(u)
	return euo
}

// SetNillableProjectID sets the "project_id" field if the given value is not nil.
func (euo *EnvironmentUpdateOne) SetNillableProjectID(u *uuid.UUID) *EnvironmentUpdateOne {
	if u != nil {
		euo.SetProjectID(*u)
	}
	return euo
}

// SetProject sets the "project" edge to the Project entity.
func (euo *EnvironmentUpdateOne) SetProject(p *Project) *EnvironmentUpdateOne {
	return euo.SetProjectID(p.ID)
}

// AddServiceIDs adds the "services" edge to the Service entity by IDs.
func (euo *EnvironmentUpdateOne) AddServiceIDs(ids ...uuid.UUID) *EnvironmentUpdateOne {
	euo.mutation.AddServiceIDs(ids...)
	return euo
}

// AddServices adds the "services" edges to the Service entity.
func (euo *EnvironmentUpdateOne) AddServices(s ...*Service) *EnvironmentUpdateOne {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return euo.AddServiceIDs(ids...)
}

// Mutation returns the EnvironmentMutation object of the builder.
func (euo *EnvironmentUpdateOne) Mutation() *EnvironmentMutation {
	return euo.mutation
}

// ClearProject clears the "project" edge to the Project entity.
func (euo *EnvironmentUpdateOne) ClearProject() *EnvironmentUpdateOne {
	euo.mutation.ClearProject()
	return euo
}

// ClearServices clears all "services" edges to the Service entity.
func (euo *EnvironmentUpdateOne) ClearServices() *EnvironmentUpdateOne {
	euo.mutation.ClearServices()
	return euo
}

// RemoveServiceIDs removes the "services" edge to Service entities by IDs.
func (euo *EnvironmentUpdateOne) RemoveServiceIDs(ids ...uuid.UUID) *EnvironmentUpdateOne {
	euo.mutation.RemoveServiceIDs(ids...)
	return euo
}

// RemoveServices removes "services" edges to Service entities.
func (euo *EnvironmentUpdateOne) RemoveServices(s ...*Service) *EnvironmentUpdateOne {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return euo.RemoveServiceIDs(ids...)
}

// Where appends a list predicates to the EnvironmentUpdate builder.
func (euo *EnvironmentUpdateOne) Where(ps ...predicate.Environment) *EnvironmentUpdateOne {
	euo.mutation.Where(ps...)
	return euo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (euo *EnvironmentUpdateOne) Select(field string, fields ...string) *EnvironmentUpdateOne {
	euo.fields = append([]string{field}, fields...)
	return euo
}

// Save executes the query and returns the updated Environment entity.
func (euo *EnvironmentUpdateOne) Save(ctx context.Context) (*Environment, error) {
	euo.defaults()
	return withHooks(ctx, euo.sqlSave, euo.mutation, euo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (euo *EnvironmentUpdateOne) SaveX(ctx context.Context) *Environment {
	node, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (euo *EnvironmentUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *EnvironmentUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (euo *EnvironmentUpdateOne) defaults() {
	if _, ok := euo.mutation.UpdatedAt(); !ok {
		v := environment.UpdateDefaultUpdatedAt()
		euo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (euo *EnvironmentUpdateOne) check() error {
	if v, ok := euo.mutation.Name(); ok {
		if err := environment.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Environment.name": %w`, err)}
		}
	}
	if euo.mutation.ProjectCleared() && len(euo.mutation.ProjectIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Environment.project"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (euo *EnvironmentUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *EnvironmentUpdateOne {
	euo.modifiers = append(euo.modifiers, modifiers...)
	return euo
}

func (euo *EnvironmentUpdateOne) sqlSave(ctx context.Context) (_node *Environment, err error) {
	if err := euo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(environment.Table, environment.Columns, sqlgraph.NewFieldSpec(environment.FieldID, field.TypeUUID))
	id, ok := euo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Environment.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := euo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, environment.FieldID)
		for _, f := range fields {
			if !environment.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != environment.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := euo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := euo.mutation.UpdatedAt(); ok {
		_spec.SetField(environment.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := euo.mutation.Name(); ok {
		_spec.SetField(environment.FieldName, field.TypeString, value)
	}
	if value, ok := euo.mutation.DisplayName(); ok {
		_spec.SetField(environment.FieldDisplayName, field.TypeString, value)
	}
	if value, ok := euo.mutation.Description(); ok {
		_spec.SetField(environment.FieldDescription, field.TypeString, value)
	}
	if value, ok := euo.mutation.Active(); ok {
		_spec.SetField(environment.FieldActive, field.TypeBool, value)
	}
	if euo.mutation.ProjectCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   environment.ProjectTable,
			Columns: []string{environment.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   environment.ProjectTable,
			Columns: []string{environment.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.mutation.ServicesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ServicesTable,
			Columns: []string{environment.ServicesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.RemovedServicesIDs(); len(nodes) > 0 && !euo.mutation.ServicesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ServicesTable,
			Columns: []string{environment.ServicesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.ServicesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ServicesTable,
			Columns: []string{environment.ServicesColumn},
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
	_spec.AddModifiers(euo.modifiers...)
	_node = &Environment{config: euo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, euo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{environment.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	euo.mutation.done = true
	return _node, nil
}
