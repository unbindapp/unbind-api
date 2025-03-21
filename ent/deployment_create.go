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
	"github.com/unbindapp/unbind-api/ent/deployment"
)

// DeploymentCreate is the builder for creating a Deployment entity.
type DeploymentCreate struct {
	config
	mutation *DeploymentMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (dc *DeploymentCreate) SetCreatedAt(t time.Time) *DeploymentCreate {
	dc.mutation.SetCreatedAt(t)
	return dc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (dc *DeploymentCreate) SetNillableCreatedAt(t *time.Time) *DeploymentCreate {
	if t != nil {
		dc.SetCreatedAt(*t)
	}
	return dc
}

// SetUpdatedAt sets the "updated_at" field.
func (dc *DeploymentCreate) SetUpdatedAt(t time.Time) *DeploymentCreate {
	dc.mutation.SetUpdatedAt(t)
	return dc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (dc *DeploymentCreate) SetNillableUpdatedAt(t *time.Time) *DeploymentCreate {
	if t != nil {
		dc.SetUpdatedAt(*t)
	}
	return dc
}

// SetID sets the "id" field.
func (dc *DeploymentCreate) SetID(u uuid.UUID) *DeploymentCreate {
	dc.mutation.SetID(u)
	return dc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (dc *DeploymentCreate) SetNillableID(u *uuid.UUID) *DeploymentCreate {
	if u != nil {
		dc.SetID(*u)
	}
	return dc
}

// Mutation returns the DeploymentMutation object of the builder.
func (dc *DeploymentCreate) Mutation() *DeploymentMutation {
	return dc.mutation
}

// Save creates the Deployment in the database.
func (dc *DeploymentCreate) Save(ctx context.Context) (*Deployment, error) {
	dc.defaults()
	return withHooks(ctx, dc.sqlSave, dc.mutation, dc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (dc *DeploymentCreate) SaveX(ctx context.Context) *Deployment {
	v, err := dc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (dc *DeploymentCreate) Exec(ctx context.Context) error {
	_, err := dc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (dc *DeploymentCreate) ExecX(ctx context.Context) {
	if err := dc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (dc *DeploymentCreate) defaults() {
	if _, ok := dc.mutation.CreatedAt(); !ok {
		v := deployment.DefaultCreatedAt()
		dc.mutation.SetCreatedAt(v)
	}
	if _, ok := dc.mutation.UpdatedAt(); !ok {
		v := deployment.DefaultUpdatedAt()
		dc.mutation.SetUpdatedAt(v)
	}
	if _, ok := dc.mutation.ID(); !ok {
		v := deployment.DefaultID()
		dc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (dc *DeploymentCreate) check() error {
	if _, ok := dc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Deployment.created_at"`)}
	}
	if _, ok := dc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Deployment.updated_at"`)}
	}
	return nil
}

func (dc *DeploymentCreate) sqlSave(ctx context.Context) (*Deployment, error) {
	if err := dc.check(); err != nil {
		return nil, err
	}
	_node, _spec := dc.createSpec()
	if err := sqlgraph.CreateNode(ctx, dc.driver, _spec); err != nil {
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
	dc.mutation.id = &_node.ID
	dc.mutation.done = true
	return _node, nil
}

func (dc *DeploymentCreate) createSpec() (*Deployment, *sqlgraph.CreateSpec) {
	var (
		_node = &Deployment{config: dc.config}
		_spec = sqlgraph.NewCreateSpec(deployment.Table, sqlgraph.NewFieldSpec(deployment.FieldID, field.TypeUUID))
	)
	_spec.OnConflict = dc.conflict
	if id, ok := dc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := dc.mutation.CreatedAt(); ok {
		_spec.SetField(deployment.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := dc.mutation.UpdatedAt(); ok {
		_spec.SetField(deployment.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Deployment.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.DeploymentUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (dc *DeploymentCreate) OnConflict(opts ...sql.ConflictOption) *DeploymentUpsertOne {
	dc.conflict = opts
	return &DeploymentUpsertOne{
		create: dc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Deployment.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (dc *DeploymentCreate) OnConflictColumns(columns ...string) *DeploymentUpsertOne {
	dc.conflict = append(dc.conflict, sql.ConflictColumns(columns...))
	return &DeploymentUpsertOne{
		create: dc,
	}
}

type (
	// DeploymentUpsertOne is the builder for "upsert"-ing
	//  one Deployment node.
	DeploymentUpsertOne struct {
		create *DeploymentCreate
	}

	// DeploymentUpsert is the "OnConflict" setter.
	DeploymentUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *DeploymentUpsert) SetUpdatedAt(v time.Time) *DeploymentUpsert {
	u.Set(deployment.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *DeploymentUpsert) UpdateUpdatedAt() *DeploymentUpsert {
	u.SetExcluded(deployment.FieldUpdatedAt)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Deployment.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(deployment.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *DeploymentUpsertOne) UpdateNewValues() *DeploymentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(deployment.FieldID)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(deployment.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Deployment.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *DeploymentUpsertOne) Ignore() *DeploymentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *DeploymentUpsertOne) DoNothing() *DeploymentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the DeploymentCreate.OnConflict
// documentation for more info.
func (u *DeploymentUpsertOne) Update(set func(*DeploymentUpsert)) *DeploymentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&DeploymentUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *DeploymentUpsertOne) SetUpdatedAt(v time.Time) *DeploymentUpsertOne {
	return u.Update(func(s *DeploymentUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *DeploymentUpsertOne) UpdateUpdatedAt() *DeploymentUpsertOne {
	return u.Update(func(s *DeploymentUpsert) {
		s.UpdateUpdatedAt()
	})
}

// Exec executes the query.
func (u *DeploymentUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for DeploymentCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *DeploymentUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *DeploymentUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: DeploymentUpsertOne.ID is not supported by MySQL driver. Use DeploymentUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *DeploymentUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// DeploymentCreateBulk is the builder for creating many Deployment entities in bulk.
type DeploymentCreateBulk struct {
	config
	err      error
	builders []*DeploymentCreate
	conflict []sql.ConflictOption
}

// Save creates the Deployment entities in the database.
func (dcb *DeploymentCreateBulk) Save(ctx context.Context) ([]*Deployment, error) {
	if dcb.err != nil {
		return nil, dcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(dcb.builders))
	nodes := make([]*Deployment, len(dcb.builders))
	mutators := make([]Mutator, len(dcb.builders))
	for i := range dcb.builders {
		func(i int, root context.Context) {
			builder := dcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*DeploymentMutation)
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
					_, err = mutators[i+1].Mutate(root, dcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = dcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, dcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, dcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (dcb *DeploymentCreateBulk) SaveX(ctx context.Context) []*Deployment {
	v, err := dcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (dcb *DeploymentCreateBulk) Exec(ctx context.Context) error {
	_, err := dcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (dcb *DeploymentCreateBulk) ExecX(ctx context.Context) {
	if err := dcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Deployment.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.DeploymentUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (dcb *DeploymentCreateBulk) OnConflict(opts ...sql.ConflictOption) *DeploymentUpsertBulk {
	dcb.conflict = opts
	return &DeploymentUpsertBulk{
		create: dcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Deployment.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (dcb *DeploymentCreateBulk) OnConflictColumns(columns ...string) *DeploymentUpsertBulk {
	dcb.conflict = append(dcb.conflict, sql.ConflictColumns(columns...))
	return &DeploymentUpsertBulk{
		create: dcb,
	}
}

// DeploymentUpsertBulk is the builder for "upsert"-ing
// a bulk of Deployment nodes.
type DeploymentUpsertBulk struct {
	create *DeploymentCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Deployment.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(deployment.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *DeploymentUpsertBulk) UpdateNewValues() *DeploymentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(deployment.FieldID)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(deployment.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Deployment.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *DeploymentUpsertBulk) Ignore() *DeploymentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *DeploymentUpsertBulk) DoNothing() *DeploymentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the DeploymentCreateBulk.OnConflict
// documentation for more info.
func (u *DeploymentUpsertBulk) Update(set func(*DeploymentUpsert)) *DeploymentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&DeploymentUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *DeploymentUpsertBulk) SetUpdatedAt(v time.Time) *DeploymentUpsertBulk {
	return u.Update(func(s *DeploymentUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *DeploymentUpsertBulk) UpdateUpdatedAt() *DeploymentUpsertBulk {
	return u.Update(func(s *DeploymentUpsert) {
		s.UpdateUpdatedAt()
	})
}

// Exec executes the query.
func (u *DeploymentUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the DeploymentCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for DeploymentCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *DeploymentUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
