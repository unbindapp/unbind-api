// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/unbindapp/unbind-api/ent/jwtkey"
)

// JWTKeyCreate is the builder for creating a JWTKey entity.
type JWTKeyCreate struct {
	config
	mutation *JWTKeyMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetLabel sets the "label" field.
func (jkc *JWTKeyCreate) SetLabel(s string) *JWTKeyCreate {
	jkc.mutation.SetLabel(s)
	return jkc
}

// SetPrivateKey sets the "private_key" field.
func (jkc *JWTKeyCreate) SetPrivateKey(b []byte) *JWTKeyCreate {
	jkc.mutation.SetPrivateKey(b)
	return jkc
}

// Mutation returns the JWTKeyMutation object of the builder.
func (jkc *JWTKeyCreate) Mutation() *JWTKeyMutation {
	return jkc.mutation
}

// Save creates the JWTKey in the database.
func (jkc *JWTKeyCreate) Save(ctx context.Context) (*JWTKey, error) {
	return withHooks(ctx, jkc.sqlSave, jkc.mutation, jkc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (jkc *JWTKeyCreate) SaveX(ctx context.Context) *JWTKey {
	v, err := jkc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (jkc *JWTKeyCreate) Exec(ctx context.Context) error {
	_, err := jkc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (jkc *JWTKeyCreate) ExecX(ctx context.Context) {
	if err := jkc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (jkc *JWTKeyCreate) check() error {
	if _, ok := jkc.mutation.Label(); !ok {
		return &ValidationError{Name: "label", err: errors.New(`ent: missing required field "JWTKey.label"`)}
	}
	if _, ok := jkc.mutation.PrivateKey(); !ok {
		return &ValidationError{Name: "private_key", err: errors.New(`ent: missing required field "JWTKey.private_key"`)}
	}
	return nil
}

func (jkc *JWTKeyCreate) sqlSave(ctx context.Context) (*JWTKey, error) {
	if err := jkc.check(); err != nil {
		return nil, err
	}
	_node, _spec := jkc.createSpec()
	if err := sqlgraph.CreateNode(ctx, jkc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	jkc.mutation.id = &_node.ID
	jkc.mutation.done = true
	return _node, nil
}

func (jkc *JWTKeyCreate) createSpec() (*JWTKey, *sqlgraph.CreateSpec) {
	var (
		_node = &JWTKey{config: jkc.config}
		_spec = sqlgraph.NewCreateSpec(jwtkey.Table, sqlgraph.NewFieldSpec(jwtkey.FieldID, field.TypeInt))
	)
	_spec.OnConflict = jkc.conflict
	if value, ok := jkc.mutation.Label(); ok {
		_spec.SetField(jwtkey.FieldLabel, field.TypeString, value)
		_node.Label = value
	}
	if value, ok := jkc.mutation.PrivateKey(); ok {
		_spec.SetField(jwtkey.FieldPrivateKey, field.TypeBytes, value)
		_node.PrivateKey = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.JWTKey.Create().
//		SetLabel(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.JWTKeyUpsert) {
//			SetLabel(v+v).
//		}).
//		Exec(ctx)
func (jkc *JWTKeyCreate) OnConflict(opts ...sql.ConflictOption) *JWTKeyUpsertOne {
	jkc.conflict = opts
	return &JWTKeyUpsertOne{
		create: jkc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.JWTKey.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (jkc *JWTKeyCreate) OnConflictColumns(columns ...string) *JWTKeyUpsertOne {
	jkc.conflict = append(jkc.conflict, sql.ConflictColumns(columns...))
	return &JWTKeyUpsertOne{
		create: jkc,
	}
}

type (
	// JWTKeyUpsertOne is the builder for "upsert"-ing
	//  one JWTKey node.
	JWTKeyUpsertOne struct {
		create *JWTKeyCreate
	}

	// JWTKeyUpsert is the "OnConflict" setter.
	JWTKeyUpsert struct {
		*sql.UpdateSet
	}
)

// SetLabel sets the "label" field.
func (u *JWTKeyUpsert) SetLabel(v string) *JWTKeyUpsert {
	u.Set(jwtkey.FieldLabel, v)
	return u
}

// UpdateLabel sets the "label" field to the value that was provided on create.
func (u *JWTKeyUpsert) UpdateLabel() *JWTKeyUpsert {
	u.SetExcluded(jwtkey.FieldLabel)
	return u
}

// SetPrivateKey sets the "private_key" field.
func (u *JWTKeyUpsert) SetPrivateKey(v []byte) *JWTKeyUpsert {
	u.Set(jwtkey.FieldPrivateKey, v)
	return u
}

// UpdatePrivateKey sets the "private_key" field to the value that was provided on create.
func (u *JWTKeyUpsert) UpdatePrivateKey() *JWTKeyUpsert {
	u.SetExcluded(jwtkey.FieldPrivateKey)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.JWTKey.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *JWTKeyUpsertOne) UpdateNewValues() *JWTKeyUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.JWTKey.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *JWTKeyUpsertOne) Ignore() *JWTKeyUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *JWTKeyUpsertOne) DoNothing() *JWTKeyUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the JWTKeyCreate.OnConflict
// documentation for more info.
func (u *JWTKeyUpsertOne) Update(set func(*JWTKeyUpsert)) *JWTKeyUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&JWTKeyUpsert{UpdateSet: update})
	}))
	return u
}

// SetLabel sets the "label" field.
func (u *JWTKeyUpsertOne) SetLabel(v string) *JWTKeyUpsertOne {
	return u.Update(func(s *JWTKeyUpsert) {
		s.SetLabel(v)
	})
}

// UpdateLabel sets the "label" field to the value that was provided on create.
func (u *JWTKeyUpsertOne) UpdateLabel() *JWTKeyUpsertOne {
	return u.Update(func(s *JWTKeyUpsert) {
		s.UpdateLabel()
	})
}

// SetPrivateKey sets the "private_key" field.
func (u *JWTKeyUpsertOne) SetPrivateKey(v []byte) *JWTKeyUpsertOne {
	return u.Update(func(s *JWTKeyUpsert) {
		s.SetPrivateKey(v)
	})
}

// UpdatePrivateKey sets the "private_key" field to the value that was provided on create.
func (u *JWTKeyUpsertOne) UpdatePrivateKey() *JWTKeyUpsertOne {
	return u.Update(func(s *JWTKeyUpsert) {
		s.UpdatePrivateKey()
	})
}

// Exec executes the query.
func (u *JWTKeyUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for JWTKeyCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *JWTKeyUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *JWTKeyUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *JWTKeyUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// JWTKeyCreateBulk is the builder for creating many JWTKey entities in bulk.
type JWTKeyCreateBulk struct {
	config
	err      error
	builders []*JWTKeyCreate
	conflict []sql.ConflictOption
}

// Save creates the JWTKey entities in the database.
func (jkcb *JWTKeyCreateBulk) Save(ctx context.Context) ([]*JWTKey, error) {
	if jkcb.err != nil {
		return nil, jkcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(jkcb.builders))
	nodes := make([]*JWTKey, len(jkcb.builders))
	mutators := make([]Mutator, len(jkcb.builders))
	for i := range jkcb.builders {
		func(i int, root context.Context) {
			builder := jkcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*JWTKeyMutation)
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
					_, err = mutators[i+1].Mutate(root, jkcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = jkcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, jkcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
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
		if _, err := mutators[0].Mutate(ctx, jkcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (jkcb *JWTKeyCreateBulk) SaveX(ctx context.Context) []*JWTKey {
	v, err := jkcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (jkcb *JWTKeyCreateBulk) Exec(ctx context.Context) error {
	_, err := jkcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (jkcb *JWTKeyCreateBulk) ExecX(ctx context.Context) {
	if err := jkcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.JWTKey.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.JWTKeyUpsert) {
//			SetLabel(v+v).
//		}).
//		Exec(ctx)
func (jkcb *JWTKeyCreateBulk) OnConflict(opts ...sql.ConflictOption) *JWTKeyUpsertBulk {
	jkcb.conflict = opts
	return &JWTKeyUpsertBulk{
		create: jkcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.JWTKey.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (jkcb *JWTKeyCreateBulk) OnConflictColumns(columns ...string) *JWTKeyUpsertBulk {
	jkcb.conflict = append(jkcb.conflict, sql.ConflictColumns(columns...))
	return &JWTKeyUpsertBulk{
		create: jkcb,
	}
}

// JWTKeyUpsertBulk is the builder for "upsert"-ing
// a bulk of JWTKey nodes.
type JWTKeyUpsertBulk struct {
	create *JWTKeyCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.JWTKey.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *JWTKeyUpsertBulk) UpdateNewValues() *JWTKeyUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.JWTKey.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *JWTKeyUpsertBulk) Ignore() *JWTKeyUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *JWTKeyUpsertBulk) DoNothing() *JWTKeyUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the JWTKeyCreateBulk.OnConflict
// documentation for more info.
func (u *JWTKeyUpsertBulk) Update(set func(*JWTKeyUpsert)) *JWTKeyUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&JWTKeyUpsert{UpdateSet: update})
	}))
	return u
}

// SetLabel sets the "label" field.
func (u *JWTKeyUpsertBulk) SetLabel(v string) *JWTKeyUpsertBulk {
	return u.Update(func(s *JWTKeyUpsert) {
		s.SetLabel(v)
	})
}

// UpdateLabel sets the "label" field to the value that was provided on create.
func (u *JWTKeyUpsertBulk) UpdateLabel() *JWTKeyUpsertBulk {
	return u.Update(func(s *JWTKeyUpsert) {
		s.UpdateLabel()
	})
}

// SetPrivateKey sets the "private_key" field.
func (u *JWTKeyUpsertBulk) SetPrivateKey(v []byte) *JWTKeyUpsertBulk {
	return u.Update(func(s *JWTKeyUpsert) {
		s.SetPrivateKey(v)
	})
}

// UpdatePrivateKey sets the "private_key" field to the value that was provided on create.
func (u *JWTKeyUpsertBulk) UpdatePrivateKey() *JWTKeyUpsertBulk {
	return u.Update(func(s *JWTKeyUpsert) {
		s.UpdatePrivateKey()
	})
}

// Exec executes the query.
func (u *JWTKeyUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the JWTKeyCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for JWTKeyCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *JWTKeyUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
