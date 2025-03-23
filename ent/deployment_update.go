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
	"github.com/unbindapp/unbind-api/ent/deployment"
	"github.com/unbindapp/unbind-api/ent/predicate"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/service"
)

// DeploymentUpdate is the builder for updating Deployment entities.
type DeploymentUpdate struct {
	config
	hooks     []Hook
	mutation  *DeploymentMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the DeploymentUpdate builder.
func (du *DeploymentUpdate) Where(ps ...predicate.Deployment) *DeploymentUpdate {
	du.mutation.Where(ps...)
	return du
}

// SetUpdatedAt sets the "updated_at" field.
func (du *DeploymentUpdate) SetUpdatedAt(t time.Time) *DeploymentUpdate {
	du.mutation.SetUpdatedAt(t)
	return du
}

// SetServiceID sets the "service_id" field.
func (du *DeploymentUpdate) SetServiceID(u uuid.UUID) *DeploymentUpdate {
	du.mutation.SetServiceID(u)
	return du
}

// SetNillableServiceID sets the "service_id" field if the given value is not nil.
func (du *DeploymentUpdate) SetNillableServiceID(u *uuid.UUID) *DeploymentUpdate {
	if u != nil {
		du.SetServiceID(*u)
	}
	return du
}

// SetStatus sets the "status" field.
func (du *DeploymentUpdate) SetStatus(ss schema.DeploymentStatus) *DeploymentUpdate {
	du.mutation.SetStatus(ss)
	return du
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (du *DeploymentUpdate) SetNillableStatus(ss *schema.DeploymentStatus) *DeploymentUpdate {
	if ss != nil {
		du.SetStatus(*ss)
	}
	return du
}

// SetError sets the "error" field.
func (du *DeploymentUpdate) SetError(s string) *DeploymentUpdate {
	du.mutation.SetError(s)
	return du
}

// SetNillableError sets the "error" field if the given value is not nil.
func (du *DeploymentUpdate) SetNillableError(s *string) *DeploymentUpdate {
	if s != nil {
		du.SetError(*s)
	}
	return du
}

// ClearError clears the value of the "error" field.
func (du *DeploymentUpdate) ClearError() *DeploymentUpdate {
	du.mutation.ClearError()
	return du
}

// SetStartedAt sets the "started_at" field.
func (du *DeploymentUpdate) SetStartedAt(t time.Time) *DeploymentUpdate {
	du.mutation.SetStartedAt(t)
	return du
}

// SetNillableStartedAt sets the "started_at" field if the given value is not nil.
func (du *DeploymentUpdate) SetNillableStartedAt(t *time.Time) *DeploymentUpdate {
	if t != nil {
		du.SetStartedAt(*t)
	}
	return du
}

// ClearStartedAt clears the value of the "started_at" field.
func (du *DeploymentUpdate) ClearStartedAt() *DeploymentUpdate {
	du.mutation.ClearStartedAt()
	return du
}

// SetCompletedAt sets the "completed_at" field.
func (du *DeploymentUpdate) SetCompletedAt(t time.Time) *DeploymentUpdate {
	du.mutation.SetCompletedAt(t)
	return du
}

// SetNillableCompletedAt sets the "completed_at" field if the given value is not nil.
func (du *DeploymentUpdate) SetNillableCompletedAt(t *time.Time) *DeploymentUpdate {
	if t != nil {
		du.SetCompletedAt(*t)
	}
	return du
}

// ClearCompletedAt clears the value of the "completed_at" field.
func (du *DeploymentUpdate) ClearCompletedAt() *DeploymentUpdate {
	du.mutation.ClearCompletedAt()
	return du
}

// SetKubernetesJobName sets the "kubernetes_job_name" field.
func (du *DeploymentUpdate) SetKubernetesJobName(s string) *DeploymentUpdate {
	du.mutation.SetKubernetesJobName(s)
	return du
}

// SetNillableKubernetesJobName sets the "kubernetes_job_name" field if the given value is not nil.
func (du *DeploymentUpdate) SetNillableKubernetesJobName(s *string) *DeploymentUpdate {
	if s != nil {
		du.SetKubernetesJobName(*s)
	}
	return du
}

// ClearKubernetesJobName clears the value of the "kubernetes_job_name" field.
func (du *DeploymentUpdate) ClearKubernetesJobName() *DeploymentUpdate {
	du.mutation.ClearKubernetesJobName()
	return du
}

// SetKubernetesJobStatus sets the "kubernetes_job_status" field.
func (du *DeploymentUpdate) SetKubernetesJobStatus(s string) *DeploymentUpdate {
	du.mutation.SetKubernetesJobStatus(s)
	return du
}

// SetNillableKubernetesJobStatus sets the "kubernetes_job_status" field if the given value is not nil.
func (du *DeploymentUpdate) SetNillableKubernetesJobStatus(s *string) *DeploymentUpdate {
	if s != nil {
		du.SetKubernetesJobStatus(*s)
	}
	return du
}

// ClearKubernetesJobStatus clears the value of the "kubernetes_job_status" field.
func (du *DeploymentUpdate) ClearKubernetesJobStatus() *DeploymentUpdate {
	du.mutation.ClearKubernetesJobStatus()
	return du
}

// SetAttempts sets the "attempts" field.
func (du *DeploymentUpdate) SetAttempts(i int) *DeploymentUpdate {
	du.mutation.ResetAttempts()
	du.mutation.SetAttempts(i)
	return du
}

// SetNillableAttempts sets the "attempts" field if the given value is not nil.
func (du *DeploymentUpdate) SetNillableAttempts(i *int) *DeploymentUpdate {
	if i != nil {
		du.SetAttempts(*i)
	}
	return du
}

// AddAttempts adds i to the "attempts" field.
func (du *DeploymentUpdate) AddAttempts(i int) *DeploymentUpdate {
	du.mutation.AddAttempts(i)
	return du
}

// SetService sets the "service" edge to the Service entity.
func (du *DeploymentUpdate) SetService(s *Service) *DeploymentUpdate {
	return du.SetServiceID(s.ID)
}

// Mutation returns the DeploymentMutation object of the builder.
func (du *DeploymentUpdate) Mutation() *DeploymentMutation {
	return du.mutation
}

// ClearService clears the "service" edge to the Service entity.
func (du *DeploymentUpdate) ClearService() *DeploymentUpdate {
	du.mutation.ClearService()
	return du
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (du *DeploymentUpdate) Save(ctx context.Context) (int, error) {
	du.defaults()
	return withHooks(ctx, du.sqlSave, du.mutation, du.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (du *DeploymentUpdate) SaveX(ctx context.Context) int {
	affected, err := du.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (du *DeploymentUpdate) Exec(ctx context.Context) error {
	_, err := du.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (du *DeploymentUpdate) ExecX(ctx context.Context) {
	if err := du.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (du *DeploymentUpdate) defaults() {
	if _, ok := du.mutation.UpdatedAt(); !ok {
		v := deployment.UpdateDefaultUpdatedAt()
		du.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (du *DeploymentUpdate) check() error {
	if v, ok := du.mutation.Status(); ok {
		if err := deployment.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Deployment.status": %w`, err)}
		}
	}
	if du.mutation.ServiceCleared() && len(du.mutation.ServiceIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Deployment.service"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (du *DeploymentUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *DeploymentUpdate {
	du.modifiers = append(du.modifiers, modifiers...)
	return du
}

func (du *DeploymentUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := du.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(deployment.Table, deployment.Columns, sqlgraph.NewFieldSpec(deployment.FieldID, field.TypeUUID))
	if ps := du.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := du.mutation.UpdatedAt(); ok {
		_spec.SetField(deployment.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := du.mutation.Status(); ok {
		_spec.SetField(deployment.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := du.mutation.Error(); ok {
		_spec.SetField(deployment.FieldError, field.TypeString, value)
	}
	if du.mutation.ErrorCleared() {
		_spec.ClearField(deployment.FieldError, field.TypeString)
	}
	if value, ok := du.mutation.StartedAt(); ok {
		_spec.SetField(deployment.FieldStartedAt, field.TypeTime, value)
	}
	if du.mutation.StartedAtCleared() {
		_spec.ClearField(deployment.FieldStartedAt, field.TypeTime)
	}
	if value, ok := du.mutation.CompletedAt(); ok {
		_spec.SetField(deployment.FieldCompletedAt, field.TypeTime, value)
	}
	if du.mutation.CompletedAtCleared() {
		_spec.ClearField(deployment.FieldCompletedAt, field.TypeTime)
	}
	if value, ok := du.mutation.KubernetesJobName(); ok {
		_spec.SetField(deployment.FieldKubernetesJobName, field.TypeString, value)
	}
	if du.mutation.KubernetesJobNameCleared() {
		_spec.ClearField(deployment.FieldKubernetesJobName, field.TypeString)
	}
	if value, ok := du.mutation.KubernetesJobStatus(); ok {
		_spec.SetField(deployment.FieldKubernetesJobStatus, field.TypeString, value)
	}
	if du.mutation.KubernetesJobStatusCleared() {
		_spec.ClearField(deployment.FieldKubernetesJobStatus, field.TypeString)
	}
	if value, ok := du.mutation.Attempts(); ok {
		_spec.SetField(deployment.FieldAttempts, field.TypeInt, value)
	}
	if value, ok := du.mutation.AddedAttempts(); ok {
		_spec.AddField(deployment.FieldAttempts, field.TypeInt, value)
	}
	if du.mutation.ServiceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   deployment.ServiceTable,
			Columns: []string{deployment.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := du.mutation.ServiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   deployment.ServiceTable,
			Columns: []string{deployment.ServiceColumn},
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
	_spec.AddModifiers(du.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, du.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{deployment.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	du.mutation.done = true
	return n, nil
}

// DeploymentUpdateOne is the builder for updating a single Deployment entity.
type DeploymentUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *DeploymentMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (duo *DeploymentUpdateOne) SetUpdatedAt(t time.Time) *DeploymentUpdateOne {
	duo.mutation.SetUpdatedAt(t)
	return duo
}

// SetServiceID sets the "service_id" field.
func (duo *DeploymentUpdateOne) SetServiceID(u uuid.UUID) *DeploymentUpdateOne {
	duo.mutation.SetServiceID(u)
	return duo
}

// SetNillableServiceID sets the "service_id" field if the given value is not nil.
func (duo *DeploymentUpdateOne) SetNillableServiceID(u *uuid.UUID) *DeploymentUpdateOne {
	if u != nil {
		duo.SetServiceID(*u)
	}
	return duo
}

// SetStatus sets the "status" field.
func (duo *DeploymentUpdateOne) SetStatus(ss schema.DeploymentStatus) *DeploymentUpdateOne {
	duo.mutation.SetStatus(ss)
	return duo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (duo *DeploymentUpdateOne) SetNillableStatus(ss *schema.DeploymentStatus) *DeploymentUpdateOne {
	if ss != nil {
		duo.SetStatus(*ss)
	}
	return duo
}

// SetError sets the "error" field.
func (duo *DeploymentUpdateOne) SetError(s string) *DeploymentUpdateOne {
	duo.mutation.SetError(s)
	return duo
}

// SetNillableError sets the "error" field if the given value is not nil.
func (duo *DeploymentUpdateOne) SetNillableError(s *string) *DeploymentUpdateOne {
	if s != nil {
		duo.SetError(*s)
	}
	return duo
}

// ClearError clears the value of the "error" field.
func (duo *DeploymentUpdateOne) ClearError() *DeploymentUpdateOne {
	duo.mutation.ClearError()
	return duo
}

// SetStartedAt sets the "started_at" field.
func (duo *DeploymentUpdateOne) SetStartedAt(t time.Time) *DeploymentUpdateOne {
	duo.mutation.SetStartedAt(t)
	return duo
}

// SetNillableStartedAt sets the "started_at" field if the given value is not nil.
func (duo *DeploymentUpdateOne) SetNillableStartedAt(t *time.Time) *DeploymentUpdateOne {
	if t != nil {
		duo.SetStartedAt(*t)
	}
	return duo
}

// ClearStartedAt clears the value of the "started_at" field.
func (duo *DeploymentUpdateOne) ClearStartedAt() *DeploymentUpdateOne {
	duo.mutation.ClearStartedAt()
	return duo
}

// SetCompletedAt sets the "completed_at" field.
func (duo *DeploymentUpdateOne) SetCompletedAt(t time.Time) *DeploymentUpdateOne {
	duo.mutation.SetCompletedAt(t)
	return duo
}

// SetNillableCompletedAt sets the "completed_at" field if the given value is not nil.
func (duo *DeploymentUpdateOne) SetNillableCompletedAt(t *time.Time) *DeploymentUpdateOne {
	if t != nil {
		duo.SetCompletedAt(*t)
	}
	return duo
}

// ClearCompletedAt clears the value of the "completed_at" field.
func (duo *DeploymentUpdateOne) ClearCompletedAt() *DeploymentUpdateOne {
	duo.mutation.ClearCompletedAt()
	return duo
}

// SetKubernetesJobName sets the "kubernetes_job_name" field.
func (duo *DeploymentUpdateOne) SetKubernetesJobName(s string) *DeploymentUpdateOne {
	duo.mutation.SetKubernetesJobName(s)
	return duo
}

// SetNillableKubernetesJobName sets the "kubernetes_job_name" field if the given value is not nil.
func (duo *DeploymentUpdateOne) SetNillableKubernetesJobName(s *string) *DeploymentUpdateOne {
	if s != nil {
		duo.SetKubernetesJobName(*s)
	}
	return duo
}

// ClearKubernetesJobName clears the value of the "kubernetes_job_name" field.
func (duo *DeploymentUpdateOne) ClearKubernetesJobName() *DeploymentUpdateOne {
	duo.mutation.ClearKubernetesJobName()
	return duo
}

// SetKubernetesJobStatus sets the "kubernetes_job_status" field.
func (duo *DeploymentUpdateOne) SetKubernetesJobStatus(s string) *DeploymentUpdateOne {
	duo.mutation.SetKubernetesJobStatus(s)
	return duo
}

// SetNillableKubernetesJobStatus sets the "kubernetes_job_status" field if the given value is not nil.
func (duo *DeploymentUpdateOne) SetNillableKubernetesJobStatus(s *string) *DeploymentUpdateOne {
	if s != nil {
		duo.SetKubernetesJobStatus(*s)
	}
	return duo
}

// ClearKubernetesJobStatus clears the value of the "kubernetes_job_status" field.
func (duo *DeploymentUpdateOne) ClearKubernetesJobStatus() *DeploymentUpdateOne {
	duo.mutation.ClearKubernetesJobStatus()
	return duo
}

// SetAttempts sets the "attempts" field.
func (duo *DeploymentUpdateOne) SetAttempts(i int) *DeploymentUpdateOne {
	duo.mutation.ResetAttempts()
	duo.mutation.SetAttempts(i)
	return duo
}

// SetNillableAttempts sets the "attempts" field if the given value is not nil.
func (duo *DeploymentUpdateOne) SetNillableAttempts(i *int) *DeploymentUpdateOne {
	if i != nil {
		duo.SetAttempts(*i)
	}
	return duo
}

// AddAttempts adds i to the "attempts" field.
func (duo *DeploymentUpdateOne) AddAttempts(i int) *DeploymentUpdateOne {
	duo.mutation.AddAttempts(i)
	return duo
}

// SetService sets the "service" edge to the Service entity.
func (duo *DeploymentUpdateOne) SetService(s *Service) *DeploymentUpdateOne {
	return duo.SetServiceID(s.ID)
}

// Mutation returns the DeploymentMutation object of the builder.
func (duo *DeploymentUpdateOne) Mutation() *DeploymentMutation {
	return duo.mutation
}

// ClearService clears the "service" edge to the Service entity.
func (duo *DeploymentUpdateOne) ClearService() *DeploymentUpdateOne {
	duo.mutation.ClearService()
	return duo
}

// Where appends a list predicates to the DeploymentUpdate builder.
func (duo *DeploymentUpdateOne) Where(ps ...predicate.Deployment) *DeploymentUpdateOne {
	duo.mutation.Where(ps...)
	return duo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (duo *DeploymentUpdateOne) Select(field string, fields ...string) *DeploymentUpdateOne {
	duo.fields = append([]string{field}, fields...)
	return duo
}

// Save executes the query and returns the updated Deployment entity.
func (duo *DeploymentUpdateOne) Save(ctx context.Context) (*Deployment, error) {
	duo.defaults()
	return withHooks(ctx, duo.sqlSave, duo.mutation, duo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (duo *DeploymentUpdateOne) SaveX(ctx context.Context) *Deployment {
	node, err := duo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (duo *DeploymentUpdateOne) Exec(ctx context.Context) error {
	_, err := duo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (duo *DeploymentUpdateOne) ExecX(ctx context.Context) {
	if err := duo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (duo *DeploymentUpdateOne) defaults() {
	if _, ok := duo.mutation.UpdatedAt(); !ok {
		v := deployment.UpdateDefaultUpdatedAt()
		duo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (duo *DeploymentUpdateOne) check() error {
	if v, ok := duo.mutation.Status(); ok {
		if err := deployment.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Deployment.status": %w`, err)}
		}
	}
	if duo.mutation.ServiceCleared() && len(duo.mutation.ServiceIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Deployment.service"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (duo *DeploymentUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *DeploymentUpdateOne {
	duo.modifiers = append(duo.modifiers, modifiers...)
	return duo
}

func (duo *DeploymentUpdateOne) sqlSave(ctx context.Context) (_node *Deployment, err error) {
	if err := duo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(deployment.Table, deployment.Columns, sqlgraph.NewFieldSpec(deployment.FieldID, field.TypeUUID))
	id, ok := duo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Deployment.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := duo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, deployment.FieldID)
		for _, f := range fields {
			if !deployment.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != deployment.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := duo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := duo.mutation.UpdatedAt(); ok {
		_spec.SetField(deployment.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := duo.mutation.Status(); ok {
		_spec.SetField(deployment.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := duo.mutation.Error(); ok {
		_spec.SetField(deployment.FieldError, field.TypeString, value)
	}
	if duo.mutation.ErrorCleared() {
		_spec.ClearField(deployment.FieldError, field.TypeString)
	}
	if value, ok := duo.mutation.StartedAt(); ok {
		_spec.SetField(deployment.FieldStartedAt, field.TypeTime, value)
	}
	if duo.mutation.StartedAtCleared() {
		_spec.ClearField(deployment.FieldStartedAt, field.TypeTime)
	}
	if value, ok := duo.mutation.CompletedAt(); ok {
		_spec.SetField(deployment.FieldCompletedAt, field.TypeTime, value)
	}
	if duo.mutation.CompletedAtCleared() {
		_spec.ClearField(deployment.FieldCompletedAt, field.TypeTime)
	}
	if value, ok := duo.mutation.KubernetesJobName(); ok {
		_spec.SetField(deployment.FieldKubernetesJobName, field.TypeString, value)
	}
	if duo.mutation.KubernetesJobNameCleared() {
		_spec.ClearField(deployment.FieldKubernetesJobName, field.TypeString)
	}
	if value, ok := duo.mutation.KubernetesJobStatus(); ok {
		_spec.SetField(deployment.FieldKubernetesJobStatus, field.TypeString, value)
	}
	if duo.mutation.KubernetesJobStatusCleared() {
		_spec.ClearField(deployment.FieldKubernetesJobStatus, field.TypeString)
	}
	if value, ok := duo.mutation.Attempts(); ok {
		_spec.SetField(deployment.FieldAttempts, field.TypeInt, value)
	}
	if value, ok := duo.mutation.AddedAttempts(); ok {
		_spec.AddField(deployment.FieldAttempts, field.TypeInt, value)
	}
	if duo.mutation.ServiceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   deployment.ServiceTable,
			Columns: []string{deployment.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := duo.mutation.ServiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   deployment.ServiceTable,
			Columns: []string{deployment.ServiceColumn},
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
	_spec.AddModifiers(duo.modifiers...)
	_node = &Deployment{config: duo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, duo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{deployment.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	duo.mutation.done = true
	return _node, nil
}
