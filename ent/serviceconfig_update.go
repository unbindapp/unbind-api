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
	"github.com/unbindapp/unbind-api/ent/predicate"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
)

// ServiceConfigUpdate is the builder for updating ServiceConfig entities.
type ServiceConfigUpdate struct {
	config
	hooks     []Hook
	mutation  *ServiceConfigMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ServiceConfigUpdate builder.
func (scu *ServiceConfigUpdate) Where(ps ...predicate.ServiceConfig) *ServiceConfigUpdate {
	scu.mutation.Where(ps...)
	return scu
}

// SetUpdatedAt sets the "updated_at" field.
func (scu *ServiceConfigUpdate) SetUpdatedAt(t time.Time) *ServiceConfigUpdate {
	scu.mutation.SetUpdatedAt(t)
	return scu
}

// SetServiceID sets the "service_id" field.
func (scu *ServiceConfigUpdate) SetServiceID(u uuid.UUID) *ServiceConfigUpdate {
	scu.mutation.SetServiceID(u)
	return scu
}

// SetNillableServiceID sets the "service_id" field if the given value is not nil.
func (scu *ServiceConfigUpdate) SetNillableServiceID(u *uuid.UUID) *ServiceConfigUpdate {
	if u != nil {
		scu.SetServiceID(*u)
	}
	return scu
}

// SetGitBranch sets the "git_branch" field.
func (scu *ServiceConfigUpdate) SetGitBranch(s string) *ServiceConfigUpdate {
	scu.mutation.SetGitBranch(s)
	return scu
}

// SetNillableGitBranch sets the "git_branch" field if the given value is not nil.
func (scu *ServiceConfigUpdate) SetNillableGitBranch(s *string) *ServiceConfigUpdate {
	if s != nil {
		scu.SetGitBranch(*s)
	}
	return scu
}

// ClearGitBranch clears the value of the "git_branch" field.
func (scu *ServiceConfigUpdate) ClearGitBranch() *ServiceConfigUpdate {
	scu.mutation.ClearGitBranch()
	return scu
}

// SetHost sets the "host" field.
func (scu *ServiceConfigUpdate) SetHost(s string) *ServiceConfigUpdate {
	scu.mutation.SetHost(s)
	return scu
}

// SetNillableHost sets the "host" field if the given value is not nil.
func (scu *ServiceConfigUpdate) SetNillableHost(s *string) *ServiceConfigUpdate {
	if s != nil {
		scu.SetHost(*s)
	}
	return scu
}

// ClearHost clears the value of the "host" field.
func (scu *ServiceConfigUpdate) ClearHost() *ServiceConfigUpdate {
	scu.mutation.ClearHost()
	return scu
}

// SetPort sets the "port" field.
func (scu *ServiceConfigUpdate) SetPort(i int) *ServiceConfigUpdate {
	scu.mutation.ResetPort()
	scu.mutation.SetPort(i)
	return scu
}

// SetNillablePort sets the "port" field if the given value is not nil.
func (scu *ServiceConfigUpdate) SetNillablePort(i *int) *ServiceConfigUpdate {
	if i != nil {
		scu.SetPort(*i)
	}
	return scu
}

// AddPort adds i to the "port" field.
func (scu *ServiceConfigUpdate) AddPort(i int) *ServiceConfigUpdate {
	scu.mutation.AddPort(i)
	return scu
}

// ClearPort clears the value of the "port" field.
func (scu *ServiceConfigUpdate) ClearPort() *ServiceConfigUpdate {
	scu.mutation.ClearPort()
	return scu
}

// SetReplicas sets the "replicas" field.
func (scu *ServiceConfigUpdate) SetReplicas(i int32) *ServiceConfigUpdate {
	scu.mutation.ResetReplicas()
	scu.mutation.SetReplicas(i)
	return scu
}

// SetNillableReplicas sets the "replicas" field if the given value is not nil.
func (scu *ServiceConfigUpdate) SetNillableReplicas(i *int32) *ServiceConfigUpdate {
	if i != nil {
		scu.SetReplicas(*i)
	}
	return scu
}

// AddReplicas adds i to the "replicas" field.
func (scu *ServiceConfigUpdate) AddReplicas(i int32) *ServiceConfigUpdate {
	scu.mutation.AddReplicas(i)
	return scu
}

// SetAutoDeploy sets the "auto_deploy" field.
func (scu *ServiceConfigUpdate) SetAutoDeploy(b bool) *ServiceConfigUpdate {
	scu.mutation.SetAutoDeploy(b)
	return scu
}

// SetNillableAutoDeploy sets the "auto_deploy" field if the given value is not nil.
func (scu *ServiceConfigUpdate) SetNillableAutoDeploy(b *bool) *ServiceConfigUpdate {
	if b != nil {
		scu.SetAutoDeploy(*b)
	}
	return scu
}

// SetRunCommand sets the "run_command" field.
func (scu *ServiceConfigUpdate) SetRunCommand(s string) *ServiceConfigUpdate {
	scu.mutation.SetRunCommand(s)
	return scu
}

// SetNillableRunCommand sets the "run_command" field if the given value is not nil.
func (scu *ServiceConfigUpdate) SetNillableRunCommand(s *string) *ServiceConfigUpdate {
	if s != nil {
		scu.SetRunCommand(*s)
	}
	return scu
}

// ClearRunCommand clears the value of the "run_command" field.
func (scu *ServiceConfigUpdate) ClearRunCommand() *ServiceConfigUpdate {
	scu.mutation.ClearRunCommand()
	return scu
}

// SetPublic sets the "public" field.
func (scu *ServiceConfigUpdate) SetPublic(b bool) *ServiceConfigUpdate {
	scu.mutation.SetPublic(b)
	return scu
}

// SetNillablePublic sets the "public" field if the given value is not nil.
func (scu *ServiceConfigUpdate) SetNillablePublic(b *bool) *ServiceConfigUpdate {
	if b != nil {
		scu.SetPublic(*b)
	}
	return scu
}

// SetImage sets the "image" field.
func (scu *ServiceConfigUpdate) SetImage(s string) *ServiceConfigUpdate {
	scu.mutation.SetImage(s)
	return scu
}

// SetNillableImage sets the "image" field if the given value is not nil.
func (scu *ServiceConfigUpdate) SetNillableImage(s *string) *ServiceConfigUpdate {
	if s != nil {
		scu.SetImage(*s)
	}
	return scu
}

// ClearImage clears the value of the "image" field.
func (scu *ServiceConfigUpdate) ClearImage() *ServiceConfigUpdate {
	scu.mutation.ClearImage()
	return scu
}

// SetService sets the "service" edge to the Service entity.
func (scu *ServiceConfigUpdate) SetService(s *Service) *ServiceConfigUpdate {
	return scu.SetServiceID(s.ID)
}

// Mutation returns the ServiceConfigMutation object of the builder.
func (scu *ServiceConfigUpdate) Mutation() *ServiceConfigMutation {
	return scu.mutation
}

// ClearService clears the "service" edge to the Service entity.
func (scu *ServiceConfigUpdate) ClearService() *ServiceConfigUpdate {
	scu.mutation.ClearService()
	return scu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (scu *ServiceConfigUpdate) Save(ctx context.Context) (int, error) {
	scu.defaults()
	return withHooks(ctx, scu.sqlSave, scu.mutation, scu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (scu *ServiceConfigUpdate) SaveX(ctx context.Context) int {
	affected, err := scu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (scu *ServiceConfigUpdate) Exec(ctx context.Context) error {
	_, err := scu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scu *ServiceConfigUpdate) ExecX(ctx context.Context) {
	if err := scu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (scu *ServiceConfigUpdate) defaults() {
	if _, ok := scu.mutation.UpdatedAt(); !ok {
		v := serviceconfig.UpdateDefaultUpdatedAt()
		scu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (scu *ServiceConfigUpdate) check() error {
	if scu.mutation.ServiceCleared() && len(scu.mutation.ServiceIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "ServiceConfig.service"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (scu *ServiceConfigUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ServiceConfigUpdate {
	scu.modifiers = append(scu.modifiers, modifiers...)
	return scu
}

func (scu *ServiceConfigUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := scu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(serviceconfig.Table, serviceconfig.Columns, sqlgraph.NewFieldSpec(serviceconfig.FieldID, field.TypeUUID))
	if ps := scu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := scu.mutation.UpdatedAt(); ok {
		_spec.SetField(serviceconfig.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := scu.mutation.GitBranch(); ok {
		_spec.SetField(serviceconfig.FieldGitBranch, field.TypeString, value)
	}
	if scu.mutation.GitBranchCleared() {
		_spec.ClearField(serviceconfig.FieldGitBranch, field.TypeString)
	}
	if value, ok := scu.mutation.Host(); ok {
		_spec.SetField(serviceconfig.FieldHost, field.TypeString, value)
	}
	if scu.mutation.HostCleared() {
		_spec.ClearField(serviceconfig.FieldHost, field.TypeString)
	}
	if value, ok := scu.mutation.Port(); ok {
		_spec.SetField(serviceconfig.FieldPort, field.TypeInt, value)
	}
	if value, ok := scu.mutation.AddedPort(); ok {
		_spec.AddField(serviceconfig.FieldPort, field.TypeInt, value)
	}
	if scu.mutation.PortCleared() {
		_spec.ClearField(serviceconfig.FieldPort, field.TypeInt)
	}
	if value, ok := scu.mutation.Replicas(); ok {
		_spec.SetField(serviceconfig.FieldReplicas, field.TypeInt32, value)
	}
	if value, ok := scu.mutation.AddedReplicas(); ok {
		_spec.AddField(serviceconfig.FieldReplicas, field.TypeInt32, value)
	}
	if value, ok := scu.mutation.AutoDeploy(); ok {
		_spec.SetField(serviceconfig.FieldAutoDeploy, field.TypeBool, value)
	}
	if value, ok := scu.mutation.RunCommand(); ok {
		_spec.SetField(serviceconfig.FieldRunCommand, field.TypeString, value)
	}
	if scu.mutation.RunCommandCleared() {
		_spec.ClearField(serviceconfig.FieldRunCommand, field.TypeString)
	}
	if value, ok := scu.mutation.Public(); ok {
		_spec.SetField(serviceconfig.FieldPublic, field.TypeBool, value)
	}
	if value, ok := scu.mutation.Image(); ok {
		_spec.SetField(serviceconfig.FieldImage, field.TypeString, value)
	}
	if scu.mutation.ImageCleared() {
		_spec.ClearField(serviceconfig.FieldImage, field.TypeString)
	}
	if scu.mutation.ServiceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   serviceconfig.ServiceTable,
			Columns: []string{serviceconfig.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := scu.mutation.ServiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   serviceconfig.ServiceTable,
			Columns: []string{serviceconfig.ServiceColumn},
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
	_spec.AddModifiers(scu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, scu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{serviceconfig.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	scu.mutation.done = true
	return n, nil
}

// ServiceConfigUpdateOne is the builder for updating a single ServiceConfig entity.
type ServiceConfigUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ServiceConfigMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (scuo *ServiceConfigUpdateOne) SetUpdatedAt(t time.Time) *ServiceConfigUpdateOne {
	scuo.mutation.SetUpdatedAt(t)
	return scuo
}

// SetServiceID sets the "service_id" field.
func (scuo *ServiceConfigUpdateOne) SetServiceID(u uuid.UUID) *ServiceConfigUpdateOne {
	scuo.mutation.SetServiceID(u)
	return scuo
}

// SetNillableServiceID sets the "service_id" field if the given value is not nil.
func (scuo *ServiceConfigUpdateOne) SetNillableServiceID(u *uuid.UUID) *ServiceConfigUpdateOne {
	if u != nil {
		scuo.SetServiceID(*u)
	}
	return scuo
}

// SetGitBranch sets the "git_branch" field.
func (scuo *ServiceConfigUpdateOne) SetGitBranch(s string) *ServiceConfigUpdateOne {
	scuo.mutation.SetGitBranch(s)
	return scuo
}

// SetNillableGitBranch sets the "git_branch" field if the given value is not nil.
func (scuo *ServiceConfigUpdateOne) SetNillableGitBranch(s *string) *ServiceConfigUpdateOne {
	if s != nil {
		scuo.SetGitBranch(*s)
	}
	return scuo
}

// ClearGitBranch clears the value of the "git_branch" field.
func (scuo *ServiceConfigUpdateOne) ClearGitBranch() *ServiceConfigUpdateOne {
	scuo.mutation.ClearGitBranch()
	return scuo
}

// SetHost sets the "host" field.
func (scuo *ServiceConfigUpdateOne) SetHost(s string) *ServiceConfigUpdateOne {
	scuo.mutation.SetHost(s)
	return scuo
}

// SetNillableHost sets the "host" field if the given value is not nil.
func (scuo *ServiceConfigUpdateOne) SetNillableHost(s *string) *ServiceConfigUpdateOne {
	if s != nil {
		scuo.SetHost(*s)
	}
	return scuo
}

// ClearHost clears the value of the "host" field.
func (scuo *ServiceConfigUpdateOne) ClearHost() *ServiceConfigUpdateOne {
	scuo.mutation.ClearHost()
	return scuo
}

// SetPort sets the "port" field.
func (scuo *ServiceConfigUpdateOne) SetPort(i int) *ServiceConfigUpdateOne {
	scuo.mutation.ResetPort()
	scuo.mutation.SetPort(i)
	return scuo
}

// SetNillablePort sets the "port" field if the given value is not nil.
func (scuo *ServiceConfigUpdateOne) SetNillablePort(i *int) *ServiceConfigUpdateOne {
	if i != nil {
		scuo.SetPort(*i)
	}
	return scuo
}

// AddPort adds i to the "port" field.
func (scuo *ServiceConfigUpdateOne) AddPort(i int) *ServiceConfigUpdateOne {
	scuo.mutation.AddPort(i)
	return scuo
}

// ClearPort clears the value of the "port" field.
func (scuo *ServiceConfigUpdateOne) ClearPort() *ServiceConfigUpdateOne {
	scuo.mutation.ClearPort()
	return scuo
}

// SetReplicas sets the "replicas" field.
func (scuo *ServiceConfigUpdateOne) SetReplicas(i int32) *ServiceConfigUpdateOne {
	scuo.mutation.ResetReplicas()
	scuo.mutation.SetReplicas(i)
	return scuo
}

// SetNillableReplicas sets the "replicas" field if the given value is not nil.
func (scuo *ServiceConfigUpdateOne) SetNillableReplicas(i *int32) *ServiceConfigUpdateOne {
	if i != nil {
		scuo.SetReplicas(*i)
	}
	return scuo
}

// AddReplicas adds i to the "replicas" field.
func (scuo *ServiceConfigUpdateOne) AddReplicas(i int32) *ServiceConfigUpdateOne {
	scuo.mutation.AddReplicas(i)
	return scuo
}

// SetAutoDeploy sets the "auto_deploy" field.
func (scuo *ServiceConfigUpdateOne) SetAutoDeploy(b bool) *ServiceConfigUpdateOne {
	scuo.mutation.SetAutoDeploy(b)
	return scuo
}

// SetNillableAutoDeploy sets the "auto_deploy" field if the given value is not nil.
func (scuo *ServiceConfigUpdateOne) SetNillableAutoDeploy(b *bool) *ServiceConfigUpdateOne {
	if b != nil {
		scuo.SetAutoDeploy(*b)
	}
	return scuo
}

// SetRunCommand sets the "run_command" field.
func (scuo *ServiceConfigUpdateOne) SetRunCommand(s string) *ServiceConfigUpdateOne {
	scuo.mutation.SetRunCommand(s)
	return scuo
}

// SetNillableRunCommand sets the "run_command" field if the given value is not nil.
func (scuo *ServiceConfigUpdateOne) SetNillableRunCommand(s *string) *ServiceConfigUpdateOne {
	if s != nil {
		scuo.SetRunCommand(*s)
	}
	return scuo
}

// ClearRunCommand clears the value of the "run_command" field.
func (scuo *ServiceConfigUpdateOne) ClearRunCommand() *ServiceConfigUpdateOne {
	scuo.mutation.ClearRunCommand()
	return scuo
}

// SetPublic sets the "public" field.
func (scuo *ServiceConfigUpdateOne) SetPublic(b bool) *ServiceConfigUpdateOne {
	scuo.mutation.SetPublic(b)
	return scuo
}

// SetNillablePublic sets the "public" field if the given value is not nil.
func (scuo *ServiceConfigUpdateOne) SetNillablePublic(b *bool) *ServiceConfigUpdateOne {
	if b != nil {
		scuo.SetPublic(*b)
	}
	return scuo
}

// SetImage sets the "image" field.
func (scuo *ServiceConfigUpdateOne) SetImage(s string) *ServiceConfigUpdateOne {
	scuo.mutation.SetImage(s)
	return scuo
}

// SetNillableImage sets the "image" field if the given value is not nil.
func (scuo *ServiceConfigUpdateOne) SetNillableImage(s *string) *ServiceConfigUpdateOne {
	if s != nil {
		scuo.SetImage(*s)
	}
	return scuo
}

// ClearImage clears the value of the "image" field.
func (scuo *ServiceConfigUpdateOne) ClearImage() *ServiceConfigUpdateOne {
	scuo.mutation.ClearImage()
	return scuo
}

// SetService sets the "service" edge to the Service entity.
func (scuo *ServiceConfigUpdateOne) SetService(s *Service) *ServiceConfigUpdateOne {
	return scuo.SetServiceID(s.ID)
}

// Mutation returns the ServiceConfigMutation object of the builder.
func (scuo *ServiceConfigUpdateOne) Mutation() *ServiceConfigMutation {
	return scuo.mutation
}

// ClearService clears the "service" edge to the Service entity.
func (scuo *ServiceConfigUpdateOne) ClearService() *ServiceConfigUpdateOne {
	scuo.mutation.ClearService()
	return scuo
}

// Where appends a list predicates to the ServiceConfigUpdate builder.
func (scuo *ServiceConfigUpdateOne) Where(ps ...predicate.ServiceConfig) *ServiceConfigUpdateOne {
	scuo.mutation.Where(ps...)
	return scuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (scuo *ServiceConfigUpdateOne) Select(field string, fields ...string) *ServiceConfigUpdateOne {
	scuo.fields = append([]string{field}, fields...)
	return scuo
}

// Save executes the query and returns the updated ServiceConfig entity.
func (scuo *ServiceConfigUpdateOne) Save(ctx context.Context) (*ServiceConfig, error) {
	scuo.defaults()
	return withHooks(ctx, scuo.sqlSave, scuo.mutation, scuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (scuo *ServiceConfigUpdateOne) SaveX(ctx context.Context) *ServiceConfig {
	node, err := scuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (scuo *ServiceConfigUpdateOne) Exec(ctx context.Context) error {
	_, err := scuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scuo *ServiceConfigUpdateOne) ExecX(ctx context.Context) {
	if err := scuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (scuo *ServiceConfigUpdateOne) defaults() {
	if _, ok := scuo.mutation.UpdatedAt(); !ok {
		v := serviceconfig.UpdateDefaultUpdatedAt()
		scuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (scuo *ServiceConfigUpdateOne) check() error {
	if scuo.mutation.ServiceCleared() && len(scuo.mutation.ServiceIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "ServiceConfig.service"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (scuo *ServiceConfigUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ServiceConfigUpdateOne {
	scuo.modifiers = append(scuo.modifiers, modifiers...)
	return scuo
}

func (scuo *ServiceConfigUpdateOne) sqlSave(ctx context.Context) (_node *ServiceConfig, err error) {
	if err := scuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(serviceconfig.Table, serviceconfig.Columns, sqlgraph.NewFieldSpec(serviceconfig.FieldID, field.TypeUUID))
	id, ok := scuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "ServiceConfig.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := scuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, serviceconfig.FieldID)
		for _, f := range fields {
			if !serviceconfig.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != serviceconfig.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := scuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := scuo.mutation.UpdatedAt(); ok {
		_spec.SetField(serviceconfig.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := scuo.mutation.GitBranch(); ok {
		_spec.SetField(serviceconfig.FieldGitBranch, field.TypeString, value)
	}
	if scuo.mutation.GitBranchCleared() {
		_spec.ClearField(serviceconfig.FieldGitBranch, field.TypeString)
	}
	if value, ok := scuo.mutation.Host(); ok {
		_spec.SetField(serviceconfig.FieldHost, field.TypeString, value)
	}
	if scuo.mutation.HostCleared() {
		_spec.ClearField(serviceconfig.FieldHost, field.TypeString)
	}
	if value, ok := scuo.mutation.Port(); ok {
		_spec.SetField(serviceconfig.FieldPort, field.TypeInt, value)
	}
	if value, ok := scuo.mutation.AddedPort(); ok {
		_spec.AddField(serviceconfig.FieldPort, field.TypeInt, value)
	}
	if scuo.mutation.PortCleared() {
		_spec.ClearField(serviceconfig.FieldPort, field.TypeInt)
	}
	if value, ok := scuo.mutation.Replicas(); ok {
		_spec.SetField(serviceconfig.FieldReplicas, field.TypeInt32, value)
	}
	if value, ok := scuo.mutation.AddedReplicas(); ok {
		_spec.AddField(serviceconfig.FieldReplicas, field.TypeInt32, value)
	}
	if value, ok := scuo.mutation.AutoDeploy(); ok {
		_spec.SetField(serviceconfig.FieldAutoDeploy, field.TypeBool, value)
	}
	if value, ok := scuo.mutation.RunCommand(); ok {
		_spec.SetField(serviceconfig.FieldRunCommand, field.TypeString, value)
	}
	if scuo.mutation.RunCommandCleared() {
		_spec.ClearField(serviceconfig.FieldRunCommand, field.TypeString)
	}
	if value, ok := scuo.mutation.Public(); ok {
		_spec.SetField(serviceconfig.FieldPublic, field.TypeBool, value)
	}
	if value, ok := scuo.mutation.Image(); ok {
		_spec.SetField(serviceconfig.FieldImage, field.TypeString, value)
	}
	if scuo.mutation.ImageCleared() {
		_spec.ClearField(serviceconfig.FieldImage, field.TypeString)
	}
	if scuo.mutation.ServiceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   serviceconfig.ServiceTable,
			Columns: []string{serviceconfig.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := scuo.mutation.ServiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   serviceconfig.ServiceTable,
			Columns: []string{serviceconfig.ServiceColumn},
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
	_spec.AddModifiers(scuo.modifiers...)
	_node = &ServiceConfig{config: scuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, scuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{serviceconfig.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	scuo.mutation.done = true
	return _node, nil
}
