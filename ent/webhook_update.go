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
	"github.com/unbindapp/unbind-api/ent/project"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/team"
	"github.com/unbindapp/unbind-api/ent/webhook"
)

// WebhookUpdate is the builder for updating Webhook entities.
type WebhookUpdate struct {
	config
	hooks     []Hook
	mutation  *WebhookMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the WebhookUpdate builder.
func (wu *WebhookUpdate) Where(ps ...predicate.Webhook) *WebhookUpdate {
	wu.mutation.Where(ps...)
	return wu
}

// SetUpdatedAt sets the "updated_at" field.
func (wu *WebhookUpdate) SetUpdatedAt(t time.Time) *WebhookUpdate {
	wu.mutation.SetUpdatedAt(t)
	return wu
}

// SetURL sets the "url" field.
func (wu *WebhookUpdate) SetURL(s string) *WebhookUpdate {
	wu.mutation.SetURL(s)
	return wu
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (wu *WebhookUpdate) SetNillableURL(s *string) *WebhookUpdate {
	if s != nil {
		wu.SetURL(*s)
	}
	return wu
}

// SetType sets the "type" field.
func (wu *WebhookUpdate) SetType(st schema.WebhookType) *WebhookUpdate {
	wu.mutation.SetType(st)
	return wu
}

// SetNillableType sets the "type" field if the given value is not nil.
func (wu *WebhookUpdate) SetNillableType(st *schema.WebhookType) *WebhookUpdate {
	if st != nil {
		wu.SetType(*st)
	}
	return wu
}

// SetEvents sets the "events" field.
func (wu *WebhookUpdate) SetEvents(se []schema.WebhookEvent) *WebhookUpdate {
	wu.mutation.SetEvents(se)
	return wu
}

// AppendEvents appends se to the "events" field.
func (wu *WebhookUpdate) AppendEvents(se []schema.WebhookEvent) *WebhookUpdate {
	wu.mutation.AppendEvents(se)
	return wu
}

// SetTeamID sets the "team_id" field.
func (wu *WebhookUpdate) SetTeamID(u uuid.UUID) *WebhookUpdate {
	wu.mutation.SetTeamID(u)
	return wu
}

// SetNillableTeamID sets the "team_id" field if the given value is not nil.
func (wu *WebhookUpdate) SetNillableTeamID(u *uuid.UUID) *WebhookUpdate {
	if u != nil {
		wu.SetTeamID(*u)
	}
	return wu
}

// SetProjectID sets the "project_id" field.
func (wu *WebhookUpdate) SetProjectID(u uuid.UUID) *WebhookUpdate {
	wu.mutation.SetProjectID(u)
	return wu
}

// SetNillableProjectID sets the "project_id" field if the given value is not nil.
func (wu *WebhookUpdate) SetNillableProjectID(u *uuid.UUID) *WebhookUpdate {
	if u != nil {
		wu.SetProjectID(*u)
	}
	return wu
}

// ClearProjectID clears the value of the "project_id" field.
func (wu *WebhookUpdate) ClearProjectID() *WebhookUpdate {
	wu.mutation.ClearProjectID()
	return wu
}

// SetTeam sets the "team" edge to the Team entity.
func (wu *WebhookUpdate) SetTeam(t *Team) *WebhookUpdate {
	return wu.SetTeamID(t.ID)
}

// SetProject sets the "project" edge to the Project entity.
func (wu *WebhookUpdate) SetProject(p *Project) *WebhookUpdate {
	return wu.SetProjectID(p.ID)
}

// Mutation returns the WebhookMutation object of the builder.
func (wu *WebhookUpdate) Mutation() *WebhookMutation {
	return wu.mutation
}

// ClearTeam clears the "team" edge to the Team entity.
func (wu *WebhookUpdate) ClearTeam() *WebhookUpdate {
	wu.mutation.ClearTeam()
	return wu
}

// ClearProject clears the "project" edge to the Project entity.
func (wu *WebhookUpdate) ClearProject() *WebhookUpdate {
	wu.mutation.ClearProject()
	return wu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (wu *WebhookUpdate) Save(ctx context.Context) (int, error) {
	wu.defaults()
	return withHooks(ctx, wu.sqlSave, wu.mutation, wu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (wu *WebhookUpdate) SaveX(ctx context.Context) int {
	affected, err := wu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (wu *WebhookUpdate) Exec(ctx context.Context) error {
	_, err := wu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wu *WebhookUpdate) ExecX(ctx context.Context) {
	if err := wu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (wu *WebhookUpdate) defaults() {
	if _, ok := wu.mutation.UpdatedAt(); !ok {
		v := webhook.UpdateDefaultUpdatedAt()
		wu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (wu *WebhookUpdate) check() error {
	if v, ok := wu.mutation.URL(); ok {
		if err := webhook.URLValidator(v); err != nil {
			return &ValidationError{Name: "url", err: fmt.Errorf(`ent: validator failed for field "Webhook.url": %w`, err)}
		}
	}
	if v, ok := wu.mutation.GetType(); ok {
		if err := webhook.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Webhook.type": %w`, err)}
		}
	}
	if wu.mutation.TeamCleared() && len(wu.mutation.TeamIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Webhook.team"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (wu *WebhookUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *WebhookUpdate {
	wu.modifiers = append(wu.modifiers, modifiers...)
	return wu
}

func (wu *WebhookUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := wu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(webhook.Table, webhook.Columns, sqlgraph.NewFieldSpec(webhook.FieldID, field.TypeUUID))
	if ps := wu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := wu.mutation.UpdatedAt(); ok {
		_spec.SetField(webhook.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := wu.mutation.URL(); ok {
		_spec.SetField(webhook.FieldURL, field.TypeString, value)
	}
	if value, ok := wu.mutation.GetType(); ok {
		_spec.SetField(webhook.FieldType, field.TypeEnum, value)
	}
	if value, ok := wu.mutation.Events(); ok {
		_spec.SetField(webhook.FieldEvents, field.TypeJSON, value)
	}
	if value, ok := wu.mutation.AppendedEvents(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, webhook.FieldEvents, value)
		})
	}
	if wu.mutation.TeamCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   webhook.TeamTable,
			Columns: []string{webhook.TeamColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(team.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wu.mutation.TeamIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   webhook.TeamTable,
			Columns: []string{webhook.TeamColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(team.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if wu.mutation.ProjectCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   webhook.ProjectTable,
			Columns: []string{webhook.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wu.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   webhook.ProjectTable,
			Columns: []string{webhook.ProjectColumn},
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
	_spec.AddModifiers(wu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, wu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{webhook.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	wu.mutation.done = true
	return n, nil
}

// WebhookUpdateOne is the builder for updating a single Webhook entity.
type WebhookUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *WebhookMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (wuo *WebhookUpdateOne) SetUpdatedAt(t time.Time) *WebhookUpdateOne {
	wuo.mutation.SetUpdatedAt(t)
	return wuo
}

// SetURL sets the "url" field.
func (wuo *WebhookUpdateOne) SetURL(s string) *WebhookUpdateOne {
	wuo.mutation.SetURL(s)
	return wuo
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (wuo *WebhookUpdateOne) SetNillableURL(s *string) *WebhookUpdateOne {
	if s != nil {
		wuo.SetURL(*s)
	}
	return wuo
}

// SetType sets the "type" field.
func (wuo *WebhookUpdateOne) SetType(st schema.WebhookType) *WebhookUpdateOne {
	wuo.mutation.SetType(st)
	return wuo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (wuo *WebhookUpdateOne) SetNillableType(st *schema.WebhookType) *WebhookUpdateOne {
	if st != nil {
		wuo.SetType(*st)
	}
	return wuo
}

// SetEvents sets the "events" field.
func (wuo *WebhookUpdateOne) SetEvents(se []schema.WebhookEvent) *WebhookUpdateOne {
	wuo.mutation.SetEvents(se)
	return wuo
}

// AppendEvents appends se to the "events" field.
func (wuo *WebhookUpdateOne) AppendEvents(se []schema.WebhookEvent) *WebhookUpdateOne {
	wuo.mutation.AppendEvents(se)
	return wuo
}

// SetTeamID sets the "team_id" field.
func (wuo *WebhookUpdateOne) SetTeamID(u uuid.UUID) *WebhookUpdateOne {
	wuo.mutation.SetTeamID(u)
	return wuo
}

// SetNillableTeamID sets the "team_id" field if the given value is not nil.
func (wuo *WebhookUpdateOne) SetNillableTeamID(u *uuid.UUID) *WebhookUpdateOne {
	if u != nil {
		wuo.SetTeamID(*u)
	}
	return wuo
}

// SetProjectID sets the "project_id" field.
func (wuo *WebhookUpdateOne) SetProjectID(u uuid.UUID) *WebhookUpdateOne {
	wuo.mutation.SetProjectID(u)
	return wuo
}

// SetNillableProjectID sets the "project_id" field if the given value is not nil.
func (wuo *WebhookUpdateOne) SetNillableProjectID(u *uuid.UUID) *WebhookUpdateOne {
	if u != nil {
		wuo.SetProjectID(*u)
	}
	return wuo
}

// ClearProjectID clears the value of the "project_id" field.
func (wuo *WebhookUpdateOne) ClearProjectID() *WebhookUpdateOne {
	wuo.mutation.ClearProjectID()
	return wuo
}

// SetTeam sets the "team" edge to the Team entity.
func (wuo *WebhookUpdateOne) SetTeam(t *Team) *WebhookUpdateOne {
	return wuo.SetTeamID(t.ID)
}

// SetProject sets the "project" edge to the Project entity.
func (wuo *WebhookUpdateOne) SetProject(p *Project) *WebhookUpdateOne {
	return wuo.SetProjectID(p.ID)
}

// Mutation returns the WebhookMutation object of the builder.
func (wuo *WebhookUpdateOne) Mutation() *WebhookMutation {
	return wuo.mutation
}

// ClearTeam clears the "team" edge to the Team entity.
func (wuo *WebhookUpdateOne) ClearTeam() *WebhookUpdateOne {
	wuo.mutation.ClearTeam()
	return wuo
}

// ClearProject clears the "project" edge to the Project entity.
func (wuo *WebhookUpdateOne) ClearProject() *WebhookUpdateOne {
	wuo.mutation.ClearProject()
	return wuo
}

// Where appends a list predicates to the WebhookUpdate builder.
func (wuo *WebhookUpdateOne) Where(ps ...predicate.Webhook) *WebhookUpdateOne {
	wuo.mutation.Where(ps...)
	return wuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (wuo *WebhookUpdateOne) Select(field string, fields ...string) *WebhookUpdateOne {
	wuo.fields = append([]string{field}, fields...)
	return wuo
}

// Save executes the query and returns the updated Webhook entity.
func (wuo *WebhookUpdateOne) Save(ctx context.Context) (*Webhook, error) {
	wuo.defaults()
	return withHooks(ctx, wuo.sqlSave, wuo.mutation, wuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (wuo *WebhookUpdateOne) SaveX(ctx context.Context) *Webhook {
	node, err := wuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (wuo *WebhookUpdateOne) Exec(ctx context.Context) error {
	_, err := wuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wuo *WebhookUpdateOne) ExecX(ctx context.Context) {
	if err := wuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (wuo *WebhookUpdateOne) defaults() {
	if _, ok := wuo.mutation.UpdatedAt(); !ok {
		v := webhook.UpdateDefaultUpdatedAt()
		wuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (wuo *WebhookUpdateOne) check() error {
	if v, ok := wuo.mutation.URL(); ok {
		if err := webhook.URLValidator(v); err != nil {
			return &ValidationError{Name: "url", err: fmt.Errorf(`ent: validator failed for field "Webhook.url": %w`, err)}
		}
	}
	if v, ok := wuo.mutation.GetType(); ok {
		if err := webhook.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Webhook.type": %w`, err)}
		}
	}
	if wuo.mutation.TeamCleared() && len(wuo.mutation.TeamIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Webhook.team"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (wuo *WebhookUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *WebhookUpdateOne {
	wuo.modifiers = append(wuo.modifiers, modifiers...)
	return wuo
}

func (wuo *WebhookUpdateOne) sqlSave(ctx context.Context) (_node *Webhook, err error) {
	if err := wuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(webhook.Table, webhook.Columns, sqlgraph.NewFieldSpec(webhook.FieldID, field.TypeUUID))
	id, ok := wuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Webhook.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := wuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, webhook.FieldID)
		for _, f := range fields {
			if !webhook.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != webhook.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := wuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := wuo.mutation.UpdatedAt(); ok {
		_spec.SetField(webhook.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := wuo.mutation.URL(); ok {
		_spec.SetField(webhook.FieldURL, field.TypeString, value)
	}
	if value, ok := wuo.mutation.GetType(); ok {
		_spec.SetField(webhook.FieldType, field.TypeEnum, value)
	}
	if value, ok := wuo.mutation.Events(); ok {
		_spec.SetField(webhook.FieldEvents, field.TypeJSON, value)
	}
	if value, ok := wuo.mutation.AppendedEvents(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, webhook.FieldEvents, value)
		})
	}
	if wuo.mutation.TeamCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   webhook.TeamTable,
			Columns: []string{webhook.TeamColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(team.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wuo.mutation.TeamIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   webhook.TeamTable,
			Columns: []string{webhook.TeamColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(team.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if wuo.mutation.ProjectCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   webhook.ProjectTable,
			Columns: []string{webhook.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wuo.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   webhook.ProjectTable,
			Columns: []string{webhook.ProjectColumn},
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
	_spec.AddModifiers(wuo.modifiers...)
	_node = &Webhook{config: wuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, wuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{webhook.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	wuo.mutation.done = true
	return _node, nil
}
