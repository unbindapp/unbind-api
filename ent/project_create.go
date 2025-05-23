// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/environment"
	"github.com/unbindapp/unbind-api/ent/project"
	"github.com/unbindapp/unbind-api/ent/team"
	"github.com/unbindapp/unbind-api/ent/webhook"
)

// ProjectCreate is the builder for creating a Project entity.
type ProjectCreate struct {
	config
	mutation *ProjectMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (pc *ProjectCreate) SetCreatedAt(t time.Time) *ProjectCreate {
	pc.mutation.SetCreatedAt(t)
	return pc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (pc *ProjectCreate) SetNillableCreatedAt(t *time.Time) *ProjectCreate {
	if t != nil {
		pc.SetCreatedAt(*t)
	}
	return pc
}

// SetUpdatedAt sets the "updated_at" field.
func (pc *ProjectCreate) SetUpdatedAt(t time.Time) *ProjectCreate {
	pc.mutation.SetUpdatedAt(t)
	return pc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (pc *ProjectCreate) SetNillableUpdatedAt(t *time.Time) *ProjectCreate {
	if t != nil {
		pc.SetUpdatedAt(*t)
	}
	return pc
}

// SetKubernetesName sets the "kubernetes_name" field.
func (pc *ProjectCreate) SetKubernetesName(s string) *ProjectCreate {
	pc.mutation.SetKubernetesName(s)
	return pc
}

// SetName sets the "name" field.
func (pc *ProjectCreate) SetName(s string) *ProjectCreate {
	pc.mutation.SetName(s)
	return pc
}

// SetDescription sets the "description" field.
func (pc *ProjectCreate) SetDescription(s string) *ProjectCreate {
	pc.mutation.SetDescription(s)
	return pc
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (pc *ProjectCreate) SetNillableDescription(s *string) *ProjectCreate {
	if s != nil {
		pc.SetDescription(*s)
	}
	return pc
}

// SetStatus sets the "status" field.
func (pc *ProjectCreate) SetStatus(s string) *ProjectCreate {
	pc.mutation.SetStatus(s)
	return pc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (pc *ProjectCreate) SetNillableStatus(s *string) *ProjectCreate {
	if s != nil {
		pc.SetStatus(*s)
	}
	return pc
}

// SetTeamID sets the "team_id" field.
func (pc *ProjectCreate) SetTeamID(u uuid.UUID) *ProjectCreate {
	pc.mutation.SetTeamID(u)
	return pc
}

// SetDefaultEnvironmentID sets the "default_environment_id" field.
func (pc *ProjectCreate) SetDefaultEnvironmentID(u uuid.UUID) *ProjectCreate {
	pc.mutation.SetDefaultEnvironmentID(u)
	return pc
}

// SetNillableDefaultEnvironmentID sets the "default_environment_id" field if the given value is not nil.
func (pc *ProjectCreate) SetNillableDefaultEnvironmentID(u *uuid.UUID) *ProjectCreate {
	if u != nil {
		pc.SetDefaultEnvironmentID(*u)
	}
	return pc
}

// SetKubernetesSecret sets the "kubernetes_secret" field.
func (pc *ProjectCreate) SetKubernetesSecret(s string) *ProjectCreate {
	pc.mutation.SetKubernetesSecret(s)
	return pc
}

// SetID sets the "id" field.
func (pc *ProjectCreate) SetID(u uuid.UUID) *ProjectCreate {
	pc.mutation.SetID(u)
	return pc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (pc *ProjectCreate) SetNillableID(u *uuid.UUID) *ProjectCreate {
	if u != nil {
		pc.SetID(*u)
	}
	return pc
}

// SetTeam sets the "team" edge to the Team entity.
func (pc *ProjectCreate) SetTeam(t *Team) *ProjectCreate {
	return pc.SetTeamID(t.ID)
}

// AddEnvironmentIDs adds the "environments" edge to the Environment entity by IDs.
func (pc *ProjectCreate) AddEnvironmentIDs(ids ...uuid.UUID) *ProjectCreate {
	pc.mutation.AddEnvironmentIDs(ids...)
	return pc
}

// AddEnvironments adds the "environments" edges to the Environment entity.
func (pc *ProjectCreate) AddEnvironments(e ...*Environment) *ProjectCreate {
	ids := make([]uuid.UUID, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return pc.AddEnvironmentIDs(ids...)
}

// SetDefaultEnvironment sets the "default_environment" edge to the Environment entity.
func (pc *ProjectCreate) SetDefaultEnvironment(e *Environment) *ProjectCreate {
	return pc.SetDefaultEnvironmentID(e.ID)
}

// AddProjectWebhookIDs adds the "project_webhooks" edge to the Webhook entity by IDs.
func (pc *ProjectCreate) AddProjectWebhookIDs(ids ...uuid.UUID) *ProjectCreate {
	pc.mutation.AddProjectWebhookIDs(ids...)
	return pc
}

// AddProjectWebhooks adds the "project_webhooks" edges to the Webhook entity.
func (pc *ProjectCreate) AddProjectWebhooks(w ...*Webhook) *ProjectCreate {
	ids := make([]uuid.UUID, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return pc.AddProjectWebhookIDs(ids...)
}

// Mutation returns the ProjectMutation object of the builder.
func (pc *ProjectCreate) Mutation() *ProjectMutation {
	return pc.mutation
}

// Save creates the Project in the database.
func (pc *ProjectCreate) Save(ctx context.Context) (*Project, error) {
	pc.defaults()
	return withHooks(ctx, pc.sqlSave, pc.mutation, pc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (pc *ProjectCreate) SaveX(ctx context.Context) *Project {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pc *ProjectCreate) Exec(ctx context.Context) error {
	_, err := pc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pc *ProjectCreate) ExecX(ctx context.Context) {
	if err := pc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pc *ProjectCreate) defaults() {
	if _, ok := pc.mutation.CreatedAt(); !ok {
		v := project.DefaultCreatedAt()
		pc.mutation.SetCreatedAt(v)
	}
	if _, ok := pc.mutation.UpdatedAt(); !ok {
		v := project.DefaultUpdatedAt()
		pc.mutation.SetUpdatedAt(v)
	}
	if _, ok := pc.mutation.Status(); !ok {
		v := project.DefaultStatus
		pc.mutation.SetStatus(v)
	}
	if _, ok := pc.mutation.ID(); !ok {
		v := project.DefaultID()
		pc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pc *ProjectCreate) check() error {
	if _, ok := pc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Project.created_at"`)}
	}
	if _, ok := pc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Project.updated_at"`)}
	}
	if _, ok := pc.mutation.KubernetesName(); !ok {
		return &ValidationError{Name: "kubernetes_name", err: errors.New(`ent: missing required field "Project.kubernetes_name"`)}
	}
	if v, ok := pc.mutation.KubernetesName(); ok {
		if err := project.KubernetesNameValidator(v); err != nil {
			return &ValidationError{Name: "kubernetes_name", err: fmt.Errorf(`ent: validator failed for field "Project.kubernetes_name": %w`, err)}
		}
	}
	if _, ok := pc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Project.name"`)}
	}
	if _, ok := pc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "Project.status"`)}
	}
	if _, ok := pc.mutation.TeamID(); !ok {
		return &ValidationError{Name: "team_id", err: errors.New(`ent: missing required field "Project.team_id"`)}
	}
	if _, ok := pc.mutation.KubernetesSecret(); !ok {
		return &ValidationError{Name: "kubernetes_secret", err: errors.New(`ent: missing required field "Project.kubernetes_secret"`)}
	}
	if len(pc.mutation.TeamIDs()) == 0 {
		return &ValidationError{Name: "team", err: errors.New(`ent: missing required edge "Project.team"`)}
	}
	return nil
}

func (pc *ProjectCreate) sqlSave(ctx context.Context) (*Project, error) {
	if err := pc.check(); err != nil {
		return nil, err
	}
	_node, _spec := pc.createSpec()
	if err := sqlgraph.CreateNode(ctx, pc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	pc.mutation.id = &_node.ID
	pc.mutation.done = true
	return _node, nil
}

func (pc *ProjectCreate) createSpec() (*Project, *sqlgraph.CreateSpec) {
	var (
		_node = &Project{config: pc.config}
		_spec = sqlgraph.NewCreateSpec(project.Table, sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID))
	)
	_spec.OnConflict = pc.conflict
	if id, ok := pc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := pc.mutation.CreatedAt(); ok {
		_spec.SetField(project.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := pc.mutation.UpdatedAt(); ok {
		_spec.SetField(project.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := pc.mutation.KubernetesName(); ok {
		_spec.SetField(project.FieldKubernetesName, field.TypeString, value)
		_node.KubernetesName = value
	}
	if value, ok := pc.mutation.Name(); ok {
		_spec.SetField(project.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := pc.mutation.Description(); ok {
		_spec.SetField(project.FieldDescription, field.TypeString, value)
		_node.Description = &value
	}
	if value, ok := pc.mutation.Status(); ok {
		_spec.SetField(project.FieldStatus, field.TypeString, value)
		_node.Status = value
	}
	if value, ok := pc.mutation.KubernetesSecret(); ok {
		_spec.SetField(project.FieldKubernetesSecret, field.TypeString, value)
		_node.KubernetesSecret = value
	}
	if nodes := pc.mutation.TeamIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   project.TeamTable,
			Columns: []string{project.TeamColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(team.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.TeamID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.EnvironmentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.EnvironmentsTable,
			Columns: []string{project.EnvironmentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(environment.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.DefaultEnvironmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   project.DefaultEnvironmentTable,
			Columns: []string{project.DefaultEnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(environment.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.DefaultEnvironmentID = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.ProjectWebhooksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.ProjectWebhooksTable,
			Columns: []string{project.ProjectWebhooksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(webhook.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Project.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ProjectUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (pc *ProjectCreate) OnConflict(opts ...sql.ConflictOption) *ProjectUpsertOne {
	pc.conflict = opts
	return &ProjectUpsertOne{
		create: pc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Project.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (pc *ProjectCreate) OnConflictColumns(columns ...string) *ProjectUpsertOne {
	pc.conflict = append(pc.conflict, sql.ConflictColumns(columns...))
	return &ProjectUpsertOne{
		create: pc,
	}
}

type (
	// ProjectUpsertOne is the builder for "upsert"-ing
	//  one Project node.
	ProjectUpsertOne struct {
		create *ProjectCreate
	}

	// ProjectUpsert is the "OnConflict" setter.
	ProjectUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *ProjectUpsert) SetUpdatedAt(v time.Time) *ProjectUpsert {
	u.Set(project.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ProjectUpsert) UpdateUpdatedAt() *ProjectUpsert {
	u.SetExcluded(project.FieldUpdatedAt)
	return u
}

// SetKubernetesName sets the "kubernetes_name" field.
func (u *ProjectUpsert) SetKubernetesName(v string) *ProjectUpsert {
	u.Set(project.FieldKubernetesName, v)
	return u
}

// UpdateKubernetesName sets the "kubernetes_name" field to the value that was provided on create.
func (u *ProjectUpsert) UpdateKubernetesName() *ProjectUpsert {
	u.SetExcluded(project.FieldKubernetesName)
	return u
}

// SetName sets the "name" field.
func (u *ProjectUpsert) SetName(v string) *ProjectUpsert {
	u.Set(project.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ProjectUpsert) UpdateName() *ProjectUpsert {
	u.SetExcluded(project.FieldName)
	return u
}

// SetDescription sets the "description" field.
func (u *ProjectUpsert) SetDescription(v string) *ProjectUpsert {
	u.Set(project.FieldDescription, v)
	return u
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ProjectUpsert) UpdateDescription() *ProjectUpsert {
	u.SetExcluded(project.FieldDescription)
	return u
}

// ClearDescription clears the value of the "description" field.
func (u *ProjectUpsert) ClearDescription() *ProjectUpsert {
	u.SetNull(project.FieldDescription)
	return u
}

// SetStatus sets the "status" field.
func (u *ProjectUpsert) SetStatus(v string) *ProjectUpsert {
	u.Set(project.FieldStatus, v)
	return u
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ProjectUpsert) UpdateStatus() *ProjectUpsert {
	u.SetExcluded(project.FieldStatus)
	return u
}

// SetTeamID sets the "team_id" field.
func (u *ProjectUpsert) SetTeamID(v uuid.UUID) *ProjectUpsert {
	u.Set(project.FieldTeamID, v)
	return u
}

// UpdateTeamID sets the "team_id" field to the value that was provided on create.
func (u *ProjectUpsert) UpdateTeamID() *ProjectUpsert {
	u.SetExcluded(project.FieldTeamID)
	return u
}

// SetDefaultEnvironmentID sets the "default_environment_id" field.
func (u *ProjectUpsert) SetDefaultEnvironmentID(v uuid.UUID) *ProjectUpsert {
	u.Set(project.FieldDefaultEnvironmentID, v)
	return u
}

// UpdateDefaultEnvironmentID sets the "default_environment_id" field to the value that was provided on create.
func (u *ProjectUpsert) UpdateDefaultEnvironmentID() *ProjectUpsert {
	u.SetExcluded(project.FieldDefaultEnvironmentID)
	return u
}

// ClearDefaultEnvironmentID clears the value of the "default_environment_id" field.
func (u *ProjectUpsert) ClearDefaultEnvironmentID() *ProjectUpsert {
	u.SetNull(project.FieldDefaultEnvironmentID)
	return u
}

// SetKubernetesSecret sets the "kubernetes_secret" field.
func (u *ProjectUpsert) SetKubernetesSecret(v string) *ProjectUpsert {
	u.Set(project.FieldKubernetesSecret, v)
	return u
}

// UpdateKubernetesSecret sets the "kubernetes_secret" field to the value that was provided on create.
func (u *ProjectUpsert) UpdateKubernetesSecret() *ProjectUpsert {
	u.SetExcluded(project.FieldKubernetesSecret)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Project.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(project.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ProjectUpsertOne) UpdateNewValues() *ProjectUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(project.FieldID)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(project.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Project.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ProjectUpsertOne) Ignore() *ProjectUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ProjectUpsertOne) DoNothing() *ProjectUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ProjectCreate.OnConflict
// documentation for more info.
func (u *ProjectUpsertOne) Update(set func(*ProjectUpsert)) *ProjectUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ProjectUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *ProjectUpsertOne) SetUpdatedAt(v time.Time) *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ProjectUpsertOne) UpdateUpdatedAt() *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetKubernetesName sets the "kubernetes_name" field.
func (u *ProjectUpsertOne) SetKubernetesName(v string) *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.SetKubernetesName(v)
	})
}

// UpdateKubernetesName sets the "kubernetes_name" field to the value that was provided on create.
func (u *ProjectUpsertOne) UpdateKubernetesName() *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateKubernetesName()
	})
}

// SetName sets the "name" field.
func (u *ProjectUpsertOne) SetName(v string) *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ProjectUpsertOne) UpdateName() *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateName()
	})
}

// SetDescription sets the "description" field.
func (u *ProjectUpsertOne) SetDescription(v string) *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ProjectUpsertOne) UpdateDescription() *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *ProjectUpsertOne) ClearDescription() *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.ClearDescription()
	})
}

// SetStatus sets the "status" field.
func (u *ProjectUpsertOne) SetStatus(v string) *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ProjectUpsertOne) UpdateStatus() *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateStatus()
	})
}

// SetTeamID sets the "team_id" field.
func (u *ProjectUpsertOne) SetTeamID(v uuid.UUID) *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.SetTeamID(v)
	})
}

// UpdateTeamID sets the "team_id" field to the value that was provided on create.
func (u *ProjectUpsertOne) UpdateTeamID() *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateTeamID()
	})
}

// SetDefaultEnvironmentID sets the "default_environment_id" field.
func (u *ProjectUpsertOne) SetDefaultEnvironmentID(v uuid.UUID) *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.SetDefaultEnvironmentID(v)
	})
}

// UpdateDefaultEnvironmentID sets the "default_environment_id" field to the value that was provided on create.
func (u *ProjectUpsertOne) UpdateDefaultEnvironmentID() *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateDefaultEnvironmentID()
	})
}

// ClearDefaultEnvironmentID clears the value of the "default_environment_id" field.
func (u *ProjectUpsertOne) ClearDefaultEnvironmentID() *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.ClearDefaultEnvironmentID()
	})
}

// SetKubernetesSecret sets the "kubernetes_secret" field.
func (u *ProjectUpsertOne) SetKubernetesSecret(v string) *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.SetKubernetesSecret(v)
	})
}

// UpdateKubernetesSecret sets the "kubernetes_secret" field to the value that was provided on create.
func (u *ProjectUpsertOne) UpdateKubernetesSecret() *ProjectUpsertOne {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateKubernetesSecret()
	})
}

// Exec executes the query.
func (u *ProjectUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ProjectCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ProjectUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ProjectUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: ProjectUpsertOne.ID is not supported by MySQL driver. Use ProjectUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ProjectUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ProjectCreateBulk is the builder for creating many Project entities in bulk.
type ProjectCreateBulk struct {
	config
	err      error
	builders []*ProjectCreate
	conflict []sql.ConflictOption
}

// Save creates the Project entities in the database.
func (pcb *ProjectCreateBulk) Save(ctx context.Context) ([]*Project, error) {
	if pcb.err != nil {
		return nil, pcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(pcb.builders))
	nodes := make([]*Project, len(pcb.builders))
	mutators := make([]Mutator, len(pcb.builders))
	for i := range pcb.builders {
		func(i int, root context.Context) {
			builder := pcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ProjectMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, pcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = pcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, pcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, pcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (pcb *ProjectCreateBulk) SaveX(ctx context.Context) []*Project {
	v, err := pcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pcb *ProjectCreateBulk) Exec(ctx context.Context) error {
	_, err := pcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pcb *ProjectCreateBulk) ExecX(ctx context.Context) {
	if err := pcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Project.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ProjectUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (pcb *ProjectCreateBulk) OnConflict(opts ...sql.ConflictOption) *ProjectUpsertBulk {
	pcb.conflict = opts
	return &ProjectUpsertBulk{
		create: pcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Project.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (pcb *ProjectCreateBulk) OnConflictColumns(columns ...string) *ProjectUpsertBulk {
	pcb.conflict = append(pcb.conflict, sql.ConflictColumns(columns...))
	return &ProjectUpsertBulk{
		create: pcb,
	}
}

// ProjectUpsertBulk is the builder for "upsert"-ing
// a bulk of Project nodes.
type ProjectUpsertBulk struct {
	create *ProjectCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Project.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(project.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ProjectUpsertBulk) UpdateNewValues() *ProjectUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(project.FieldID)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(project.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Project.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ProjectUpsertBulk) Ignore() *ProjectUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ProjectUpsertBulk) DoNothing() *ProjectUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ProjectCreateBulk.OnConflict
// documentation for more info.
func (u *ProjectUpsertBulk) Update(set func(*ProjectUpsert)) *ProjectUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ProjectUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *ProjectUpsertBulk) SetUpdatedAt(v time.Time) *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ProjectUpsertBulk) UpdateUpdatedAt() *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetKubernetesName sets the "kubernetes_name" field.
func (u *ProjectUpsertBulk) SetKubernetesName(v string) *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.SetKubernetesName(v)
	})
}

// UpdateKubernetesName sets the "kubernetes_name" field to the value that was provided on create.
func (u *ProjectUpsertBulk) UpdateKubernetesName() *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateKubernetesName()
	})
}

// SetName sets the "name" field.
func (u *ProjectUpsertBulk) SetName(v string) *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ProjectUpsertBulk) UpdateName() *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateName()
	})
}

// SetDescription sets the "description" field.
func (u *ProjectUpsertBulk) SetDescription(v string) *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ProjectUpsertBulk) UpdateDescription() *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *ProjectUpsertBulk) ClearDescription() *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.ClearDescription()
	})
}

// SetStatus sets the "status" field.
func (u *ProjectUpsertBulk) SetStatus(v string) *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ProjectUpsertBulk) UpdateStatus() *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateStatus()
	})
}

// SetTeamID sets the "team_id" field.
func (u *ProjectUpsertBulk) SetTeamID(v uuid.UUID) *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.SetTeamID(v)
	})
}

// UpdateTeamID sets the "team_id" field to the value that was provided on create.
func (u *ProjectUpsertBulk) UpdateTeamID() *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateTeamID()
	})
}

// SetDefaultEnvironmentID sets the "default_environment_id" field.
func (u *ProjectUpsertBulk) SetDefaultEnvironmentID(v uuid.UUID) *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.SetDefaultEnvironmentID(v)
	})
}

// UpdateDefaultEnvironmentID sets the "default_environment_id" field to the value that was provided on create.
func (u *ProjectUpsertBulk) UpdateDefaultEnvironmentID() *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateDefaultEnvironmentID()
	})
}

// ClearDefaultEnvironmentID clears the value of the "default_environment_id" field.
func (u *ProjectUpsertBulk) ClearDefaultEnvironmentID() *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.ClearDefaultEnvironmentID()
	})
}

// SetKubernetesSecret sets the "kubernetes_secret" field.
func (u *ProjectUpsertBulk) SetKubernetesSecret(v string) *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.SetKubernetesSecret(v)
	})
}

// UpdateKubernetesSecret sets the "kubernetes_secret" field to the value that was provided on create.
func (u *ProjectUpsertBulk) UpdateKubernetesSecret() *ProjectUpsertBulk {
	return u.Update(func(s *ProjectUpsert) {
		s.UpdateKubernetesSecret()
	})
}

// Exec executes the query.
func (u *ProjectUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the ProjectCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ProjectCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ProjectUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
