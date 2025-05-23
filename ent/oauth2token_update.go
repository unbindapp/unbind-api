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
	"github.com/unbindapp/unbind-api/ent/oauth2token"
	"github.com/unbindapp/unbind-api/ent/predicate"
	"github.com/unbindapp/unbind-api/ent/user"
)

// Oauth2TokenUpdate is the builder for updating Oauth2Token entities.
type Oauth2TokenUpdate struct {
	config
	hooks     []Hook
	mutation  *Oauth2TokenMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the Oauth2TokenUpdate builder.
func (ou *Oauth2TokenUpdate) Where(ps ...predicate.Oauth2Token) *Oauth2TokenUpdate {
	ou.mutation.Where(ps...)
	return ou
}

// SetUpdatedAt sets the "updated_at" field.
func (ou *Oauth2TokenUpdate) SetUpdatedAt(t time.Time) *Oauth2TokenUpdate {
	ou.mutation.SetUpdatedAt(t)
	return ou
}

// SetAccessToken sets the "access_token" field.
func (ou *Oauth2TokenUpdate) SetAccessToken(s string) *Oauth2TokenUpdate {
	ou.mutation.SetAccessToken(s)
	return ou
}

// SetNillableAccessToken sets the "access_token" field if the given value is not nil.
func (ou *Oauth2TokenUpdate) SetNillableAccessToken(s *string) *Oauth2TokenUpdate {
	if s != nil {
		ou.SetAccessToken(*s)
	}
	return ou
}

// SetRefreshToken sets the "refresh_token" field.
func (ou *Oauth2TokenUpdate) SetRefreshToken(s string) *Oauth2TokenUpdate {
	ou.mutation.SetRefreshToken(s)
	return ou
}

// SetNillableRefreshToken sets the "refresh_token" field if the given value is not nil.
func (ou *Oauth2TokenUpdate) SetNillableRefreshToken(s *string) *Oauth2TokenUpdate {
	if s != nil {
		ou.SetRefreshToken(*s)
	}
	return ou
}

// SetClientID sets the "client_id" field.
func (ou *Oauth2TokenUpdate) SetClientID(s string) *Oauth2TokenUpdate {
	ou.mutation.SetClientID(s)
	return ou
}

// SetNillableClientID sets the "client_id" field if the given value is not nil.
func (ou *Oauth2TokenUpdate) SetNillableClientID(s *string) *Oauth2TokenUpdate {
	if s != nil {
		ou.SetClientID(*s)
	}
	return ou
}

// SetExpiresAt sets the "expires_at" field.
func (ou *Oauth2TokenUpdate) SetExpiresAt(t time.Time) *Oauth2TokenUpdate {
	ou.mutation.SetExpiresAt(t)
	return ou
}

// SetNillableExpiresAt sets the "expires_at" field if the given value is not nil.
func (ou *Oauth2TokenUpdate) SetNillableExpiresAt(t *time.Time) *Oauth2TokenUpdate {
	if t != nil {
		ou.SetExpiresAt(*t)
	}
	return ou
}

// SetRevoked sets the "revoked" field.
func (ou *Oauth2TokenUpdate) SetRevoked(b bool) *Oauth2TokenUpdate {
	ou.mutation.SetRevoked(b)
	return ou
}

// SetNillableRevoked sets the "revoked" field if the given value is not nil.
func (ou *Oauth2TokenUpdate) SetNillableRevoked(b *bool) *Oauth2TokenUpdate {
	if b != nil {
		ou.SetRevoked(*b)
	}
	return ou
}

// SetScope sets the "scope" field.
func (ou *Oauth2TokenUpdate) SetScope(s string) *Oauth2TokenUpdate {
	ou.mutation.SetScope(s)
	return ou
}

// SetNillableScope sets the "scope" field if the given value is not nil.
func (ou *Oauth2TokenUpdate) SetNillableScope(s *string) *Oauth2TokenUpdate {
	if s != nil {
		ou.SetScope(*s)
	}
	return ou
}

// SetDeviceInfo sets the "device_info" field.
func (ou *Oauth2TokenUpdate) SetDeviceInfo(s string) *Oauth2TokenUpdate {
	ou.mutation.SetDeviceInfo(s)
	return ou
}

// SetNillableDeviceInfo sets the "device_info" field if the given value is not nil.
func (ou *Oauth2TokenUpdate) SetNillableDeviceInfo(s *string) *Oauth2TokenUpdate {
	if s != nil {
		ou.SetDeviceInfo(*s)
	}
	return ou
}

// ClearDeviceInfo clears the value of the "device_info" field.
func (ou *Oauth2TokenUpdate) ClearDeviceInfo() *Oauth2TokenUpdate {
	ou.mutation.ClearDeviceInfo()
	return ou
}

// SetUserID sets the "user" edge to the User entity by ID.
func (ou *Oauth2TokenUpdate) SetUserID(id uuid.UUID) *Oauth2TokenUpdate {
	ou.mutation.SetUserID(id)
	return ou
}

// SetUser sets the "user" edge to the User entity.
func (ou *Oauth2TokenUpdate) SetUser(u *User) *Oauth2TokenUpdate {
	return ou.SetUserID(u.ID)
}

// Mutation returns the Oauth2TokenMutation object of the builder.
func (ou *Oauth2TokenUpdate) Mutation() *Oauth2TokenMutation {
	return ou.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (ou *Oauth2TokenUpdate) ClearUser() *Oauth2TokenUpdate {
	ou.mutation.ClearUser()
	return ou
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ou *Oauth2TokenUpdate) Save(ctx context.Context) (int, error) {
	ou.defaults()
	return withHooks(ctx, ou.sqlSave, ou.mutation, ou.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ou *Oauth2TokenUpdate) SaveX(ctx context.Context) int {
	affected, err := ou.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ou *Oauth2TokenUpdate) Exec(ctx context.Context) error {
	_, err := ou.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ou *Oauth2TokenUpdate) ExecX(ctx context.Context) {
	if err := ou.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ou *Oauth2TokenUpdate) defaults() {
	if _, ok := ou.mutation.UpdatedAt(); !ok {
		v := oauth2token.UpdateDefaultUpdatedAt()
		ou.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ou *Oauth2TokenUpdate) check() error {
	if ou.mutation.UserCleared() && len(ou.mutation.UserIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Oauth2Token.user"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ou *Oauth2TokenUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *Oauth2TokenUpdate {
	ou.modifiers = append(ou.modifiers, modifiers...)
	return ou
}

func (ou *Oauth2TokenUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ou.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(oauth2token.Table, oauth2token.Columns, sqlgraph.NewFieldSpec(oauth2token.FieldID, field.TypeUUID))
	if ps := ou.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ou.mutation.UpdatedAt(); ok {
		_spec.SetField(oauth2token.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := ou.mutation.AccessToken(); ok {
		_spec.SetField(oauth2token.FieldAccessToken, field.TypeString, value)
	}
	if value, ok := ou.mutation.RefreshToken(); ok {
		_spec.SetField(oauth2token.FieldRefreshToken, field.TypeString, value)
	}
	if value, ok := ou.mutation.ClientID(); ok {
		_spec.SetField(oauth2token.FieldClientID, field.TypeString, value)
	}
	if value, ok := ou.mutation.ExpiresAt(); ok {
		_spec.SetField(oauth2token.FieldExpiresAt, field.TypeTime, value)
	}
	if value, ok := ou.mutation.Revoked(); ok {
		_spec.SetField(oauth2token.FieldRevoked, field.TypeBool, value)
	}
	if value, ok := ou.mutation.Scope(); ok {
		_spec.SetField(oauth2token.FieldScope, field.TypeString, value)
	}
	if value, ok := ou.mutation.DeviceInfo(); ok {
		_spec.SetField(oauth2token.FieldDeviceInfo, field.TypeString, value)
	}
	if ou.mutation.DeviceInfoCleared() {
		_spec.ClearField(oauth2token.FieldDeviceInfo, field.TypeString)
	}
	if ou.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   oauth2token.UserTable,
			Columns: []string{oauth2token.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ou.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   oauth2token.UserTable,
			Columns: []string{oauth2token.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(ou.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, ou.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{oauth2token.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ou.mutation.done = true
	return n, nil
}

// Oauth2TokenUpdateOne is the builder for updating a single Oauth2Token entity.
type Oauth2TokenUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *Oauth2TokenMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (ouo *Oauth2TokenUpdateOne) SetUpdatedAt(t time.Time) *Oauth2TokenUpdateOne {
	ouo.mutation.SetUpdatedAt(t)
	return ouo
}

// SetAccessToken sets the "access_token" field.
func (ouo *Oauth2TokenUpdateOne) SetAccessToken(s string) *Oauth2TokenUpdateOne {
	ouo.mutation.SetAccessToken(s)
	return ouo
}

// SetNillableAccessToken sets the "access_token" field if the given value is not nil.
func (ouo *Oauth2TokenUpdateOne) SetNillableAccessToken(s *string) *Oauth2TokenUpdateOne {
	if s != nil {
		ouo.SetAccessToken(*s)
	}
	return ouo
}

// SetRefreshToken sets the "refresh_token" field.
func (ouo *Oauth2TokenUpdateOne) SetRefreshToken(s string) *Oauth2TokenUpdateOne {
	ouo.mutation.SetRefreshToken(s)
	return ouo
}

// SetNillableRefreshToken sets the "refresh_token" field if the given value is not nil.
func (ouo *Oauth2TokenUpdateOne) SetNillableRefreshToken(s *string) *Oauth2TokenUpdateOne {
	if s != nil {
		ouo.SetRefreshToken(*s)
	}
	return ouo
}

// SetClientID sets the "client_id" field.
func (ouo *Oauth2TokenUpdateOne) SetClientID(s string) *Oauth2TokenUpdateOne {
	ouo.mutation.SetClientID(s)
	return ouo
}

// SetNillableClientID sets the "client_id" field if the given value is not nil.
func (ouo *Oauth2TokenUpdateOne) SetNillableClientID(s *string) *Oauth2TokenUpdateOne {
	if s != nil {
		ouo.SetClientID(*s)
	}
	return ouo
}

// SetExpiresAt sets the "expires_at" field.
func (ouo *Oauth2TokenUpdateOne) SetExpiresAt(t time.Time) *Oauth2TokenUpdateOne {
	ouo.mutation.SetExpiresAt(t)
	return ouo
}

// SetNillableExpiresAt sets the "expires_at" field if the given value is not nil.
func (ouo *Oauth2TokenUpdateOne) SetNillableExpiresAt(t *time.Time) *Oauth2TokenUpdateOne {
	if t != nil {
		ouo.SetExpiresAt(*t)
	}
	return ouo
}

// SetRevoked sets the "revoked" field.
func (ouo *Oauth2TokenUpdateOne) SetRevoked(b bool) *Oauth2TokenUpdateOne {
	ouo.mutation.SetRevoked(b)
	return ouo
}

// SetNillableRevoked sets the "revoked" field if the given value is not nil.
func (ouo *Oauth2TokenUpdateOne) SetNillableRevoked(b *bool) *Oauth2TokenUpdateOne {
	if b != nil {
		ouo.SetRevoked(*b)
	}
	return ouo
}

// SetScope sets the "scope" field.
func (ouo *Oauth2TokenUpdateOne) SetScope(s string) *Oauth2TokenUpdateOne {
	ouo.mutation.SetScope(s)
	return ouo
}

// SetNillableScope sets the "scope" field if the given value is not nil.
func (ouo *Oauth2TokenUpdateOne) SetNillableScope(s *string) *Oauth2TokenUpdateOne {
	if s != nil {
		ouo.SetScope(*s)
	}
	return ouo
}

// SetDeviceInfo sets the "device_info" field.
func (ouo *Oauth2TokenUpdateOne) SetDeviceInfo(s string) *Oauth2TokenUpdateOne {
	ouo.mutation.SetDeviceInfo(s)
	return ouo
}

// SetNillableDeviceInfo sets the "device_info" field if the given value is not nil.
func (ouo *Oauth2TokenUpdateOne) SetNillableDeviceInfo(s *string) *Oauth2TokenUpdateOne {
	if s != nil {
		ouo.SetDeviceInfo(*s)
	}
	return ouo
}

// ClearDeviceInfo clears the value of the "device_info" field.
func (ouo *Oauth2TokenUpdateOne) ClearDeviceInfo() *Oauth2TokenUpdateOne {
	ouo.mutation.ClearDeviceInfo()
	return ouo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (ouo *Oauth2TokenUpdateOne) SetUserID(id uuid.UUID) *Oauth2TokenUpdateOne {
	ouo.mutation.SetUserID(id)
	return ouo
}

// SetUser sets the "user" edge to the User entity.
func (ouo *Oauth2TokenUpdateOne) SetUser(u *User) *Oauth2TokenUpdateOne {
	return ouo.SetUserID(u.ID)
}

// Mutation returns the Oauth2TokenMutation object of the builder.
func (ouo *Oauth2TokenUpdateOne) Mutation() *Oauth2TokenMutation {
	return ouo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (ouo *Oauth2TokenUpdateOne) ClearUser() *Oauth2TokenUpdateOne {
	ouo.mutation.ClearUser()
	return ouo
}

// Where appends a list predicates to the Oauth2TokenUpdate builder.
func (ouo *Oauth2TokenUpdateOne) Where(ps ...predicate.Oauth2Token) *Oauth2TokenUpdateOne {
	ouo.mutation.Where(ps...)
	return ouo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ouo *Oauth2TokenUpdateOne) Select(field string, fields ...string) *Oauth2TokenUpdateOne {
	ouo.fields = append([]string{field}, fields...)
	return ouo
}

// Save executes the query and returns the updated Oauth2Token entity.
func (ouo *Oauth2TokenUpdateOne) Save(ctx context.Context) (*Oauth2Token, error) {
	ouo.defaults()
	return withHooks(ctx, ouo.sqlSave, ouo.mutation, ouo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ouo *Oauth2TokenUpdateOne) SaveX(ctx context.Context) *Oauth2Token {
	node, err := ouo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ouo *Oauth2TokenUpdateOne) Exec(ctx context.Context) error {
	_, err := ouo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ouo *Oauth2TokenUpdateOne) ExecX(ctx context.Context) {
	if err := ouo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ouo *Oauth2TokenUpdateOne) defaults() {
	if _, ok := ouo.mutation.UpdatedAt(); !ok {
		v := oauth2token.UpdateDefaultUpdatedAt()
		ouo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ouo *Oauth2TokenUpdateOne) check() error {
	if ouo.mutation.UserCleared() && len(ouo.mutation.UserIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Oauth2Token.user"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ouo *Oauth2TokenUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *Oauth2TokenUpdateOne {
	ouo.modifiers = append(ouo.modifiers, modifiers...)
	return ouo
}

func (ouo *Oauth2TokenUpdateOne) sqlSave(ctx context.Context) (_node *Oauth2Token, err error) {
	if err := ouo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(oauth2token.Table, oauth2token.Columns, sqlgraph.NewFieldSpec(oauth2token.FieldID, field.TypeUUID))
	id, ok := ouo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Oauth2Token.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ouo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, oauth2token.FieldID)
		for _, f := range fields {
			if !oauth2token.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != oauth2token.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ouo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ouo.mutation.UpdatedAt(); ok {
		_spec.SetField(oauth2token.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := ouo.mutation.AccessToken(); ok {
		_spec.SetField(oauth2token.FieldAccessToken, field.TypeString, value)
	}
	if value, ok := ouo.mutation.RefreshToken(); ok {
		_spec.SetField(oauth2token.FieldRefreshToken, field.TypeString, value)
	}
	if value, ok := ouo.mutation.ClientID(); ok {
		_spec.SetField(oauth2token.FieldClientID, field.TypeString, value)
	}
	if value, ok := ouo.mutation.ExpiresAt(); ok {
		_spec.SetField(oauth2token.FieldExpiresAt, field.TypeTime, value)
	}
	if value, ok := ouo.mutation.Revoked(); ok {
		_spec.SetField(oauth2token.FieldRevoked, field.TypeBool, value)
	}
	if value, ok := ouo.mutation.Scope(); ok {
		_spec.SetField(oauth2token.FieldScope, field.TypeString, value)
	}
	if value, ok := ouo.mutation.DeviceInfo(); ok {
		_spec.SetField(oauth2token.FieldDeviceInfo, field.TypeString, value)
	}
	if ouo.mutation.DeviceInfoCleared() {
		_spec.ClearField(oauth2token.FieldDeviceInfo, field.TypeString)
	}
	if ouo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   oauth2token.UserTable,
			Columns: []string{oauth2token.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ouo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   oauth2token.UserTable,
			Columns: []string{oauth2token.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(ouo.modifiers...)
	_node = &Oauth2Token{config: ouo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ouo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{oauth2token.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ouo.mutation.done = true
	return _node, nil
}
