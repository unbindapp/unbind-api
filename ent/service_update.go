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
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/ent/predicate"
	"github.com/unbindapp/unbind-api/ent/project"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
)

// ServiceUpdate is the builder for updating Service entities.
type ServiceUpdate struct {
	config
	hooks     []Hook
	mutation  *ServiceMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ServiceUpdate builder.
func (su *ServiceUpdate) Where(ps ...predicate.Service) *ServiceUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetUpdatedAt sets the "updated_at" field.
func (su *ServiceUpdate) SetUpdatedAt(t time.Time) *ServiceUpdate {
	su.mutation.SetUpdatedAt(t)
	return su
}

// SetName sets the "name" field.
func (su *ServiceUpdate) SetName(s string) *ServiceUpdate {
	su.mutation.SetName(s)
	return su
}

// SetNillableName sets the "name" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableName(s *string) *ServiceUpdate {
	if s != nil {
		su.SetName(*s)
	}
	return su
}

// SetDisplayName sets the "display_name" field.
func (su *ServiceUpdate) SetDisplayName(s string) *ServiceUpdate {
	su.mutation.SetDisplayName(s)
	return su
}

// SetNillableDisplayName sets the "display_name" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableDisplayName(s *string) *ServiceUpdate {
	if s != nil {
		su.SetDisplayName(*s)
	}
	return su
}

// SetDescription sets the "description" field.
func (su *ServiceUpdate) SetDescription(s string) *ServiceUpdate {
	su.mutation.SetDescription(s)
	return su
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableDescription(s *string) *ServiceUpdate {
	if s != nil {
		su.SetDescription(*s)
	}
	return su
}

// ClearDescription clears the value of the "description" field.
func (su *ServiceUpdate) ClearDescription() *ServiceUpdate {
	su.mutation.ClearDescription()
	return su
}

// SetType sets the "type" field.
func (su *ServiceUpdate) SetType(s service.Type) *ServiceUpdate {
	su.mutation.SetType(s)
	return su
}

// SetNillableType sets the "type" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableType(s *service.Type) *ServiceUpdate {
	if s != nil {
		su.SetType(*s)
	}
	return su
}

// SetSubtype sets the "subtype" field.
func (su *ServiceUpdate) SetSubtype(s service.Subtype) *ServiceUpdate {
	su.mutation.SetSubtype(s)
	return su
}

// SetNillableSubtype sets the "subtype" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableSubtype(s *service.Subtype) *ServiceUpdate {
	if s != nil {
		su.SetSubtype(*s)
	}
	return su
}

// SetProjectID sets the "project_id" field.
func (su *ServiceUpdate) SetProjectID(u uuid.UUID) *ServiceUpdate {
	su.mutation.SetProjectID(u)
	return su
}

// SetNillableProjectID sets the "project_id" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableProjectID(u *uuid.UUID) *ServiceUpdate {
	if u != nil {
		su.SetProjectID(*u)
	}
	return su
}

// SetGithubInstallationID sets the "github_installation_id" field.
func (su *ServiceUpdate) SetGithubInstallationID(i int64) *ServiceUpdate {
	su.mutation.SetGithubInstallationID(i)
	return su
}

// SetNillableGithubInstallationID sets the "github_installation_id" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableGithubInstallationID(i *int64) *ServiceUpdate {
	if i != nil {
		su.SetGithubInstallationID(*i)
	}
	return su
}

// ClearGithubInstallationID clears the value of the "github_installation_id" field.
func (su *ServiceUpdate) ClearGithubInstallationID() *ServiceUpdate {
	su.mutation.ClearGithubInstallationID()
	return su
}

// SetGitRepository sets the "git_repository" field.
func (su *ServiceUpdate) SetGitRepository(s string) *ServiceUpdate {
	su.mutation.SetGitRepository(s)
	return su
}

// SetNillableGitRepository sets the "git_repository" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableGitRepository(s *string) *ServiceUpdate {
	if s != nil {
		su.SetGitRepository(*s)
	}
	return su
}

// ClearGitRepository clears the value of the "git_repository" field.
func (su *ServiceUpdate) ClearGitRepository() *ServiceUpdate {
	su.mutation.ClearGitRepository()
	return su
}

// SetGitBranch sets the "git_branch" field.
func (su *ServiceUpdate) SetGitBranch(s string) *ServiceUpdate {
	su.mutation.SetGitBranch(s)
	return su
}

// SetNillableGitBranch sets the "git_branch" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableGitBranch(s *string) *ServiceUpdate {
	if s != nil {
		su.SetGitBranch(*s)
	}
	return su
}

// ClearGitBranch clears the value of the "git_branch" field.
func (su *ServiceUpdate) ClearGitBranch() *ServiceUpdate {
	su.mutation.ClearGitBranch()
	return su
}

// SetProject sets the "project" edge to the Project entity.
func (su *ServiceUpdate) SetProject(p *Project) *ServiceUpdate {
	return su.SetProjectID(p.ID)
}

// SetGithubInstallation sets the "github_installation" edge to the GithubInstallation entity.
func (su *ServiceUpdate) SetGithubInstallation(g *GithubInstallation) *ServiceUpdate {
	return su.SetGithubInstallationID(g.ID)
}

// SetServiceConfigsID sets the "service_configs" edge to the ServiceConfig entity by ID.
func (su *ServiceUpdate) SetServiceConfigsID(id uuid.UUID) *ServiceUpdate {
	su.mutation.SetServiceConfigsID(id)
	return su
}

// SetNillableServiceConfigsID sets the "service_configs" edge to the ServiceConfig entity by ID if the given value is not nil.
func (su *ServiceUpdate) SetNillableServiceConfigsID(id *uuid.UUID) *ServiceUpdate {
	if id != nil {
		su = su.SetServiceConfigsID(*id)
	}
	return su
}

// SetServiceConfigs sets the "service_configs" edge to the ServiceConfig entity.
func (su *ServiceUpdate) SetServiceConfigs(s *ServiceConfig) *ServiceUpdate {
	return su.SetServiceConfigsID(s.ID)
}

// Mutation returns the ServiceMutation object of the builder.
func (su *ServiceUpdate) Mutation() *ServiceMutation {
	return su.mutation
}

// ClearProject clears the "project" edge to the Project entity.
func (su *ServiceUpdate) ClearProject() *ServiceUpdate {
	su.mutation.ClearProject()
	return su
}

// ClearGithubInstallation clears the "github_installation" edge to the GithubInstallation entity.
func (su *ServiceUpdate) ClearGithubInstallation() *ServiceUpdate {
	su.mutation.ClearGithubInstallation()
	return su
}

// ClearServiceConfigs clears the "service_configs" edge to the ServiceConfig entity.
func (su *ServiceUpdate) ClearServiceConfigs() *ServiceUpdate {
	su.mutation.ClearServiceConfigs()
	return su
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *ServiceUpdate) Save(ctx context.Context) (int, error) {
	su.defaults()
	return withHooks(ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *ServiceUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *ServiceUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *ServiceUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (su *ServiceUpdate) defaults() {
	if _, ok := su.mutation.UpdatedAt(); !ok {
		v := service.UpdateDefaultUpdatedAt()
		su.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *ServiceUpdate) check() error {
	if v, ok := su.mutation.Name(); ok {
		if err := service.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Service.name": %w`, err)}
		}
	}
	if v, ok := su.mutation.GetType(); ok {
		if err := service.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Service.type": %w`, err)}
		}
	}
	if v, ok := su.mutation.Subtype(); ok {
		if err := service.SubtypeValidator(v); err != nil {
			return &ValidationError{Name: "subtype", err: fmt.Errorf(`ent: validator failed for field "Service.subtype": %w`, err)}
		}
	}
	if su.mutation.ProjectCleared() && len(su.mutation.ProjectIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Service.project"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (su *ServiceUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ServiceUpdate {
	su.modifiers = append(su.modifiers, modifiers...)
	return su
}

func (su *ServiceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := su.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(service.Table, service.Columns, sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.UpdatedAt(); ok {
		_spec.SetField(service.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := su.mutation.Name(); ok {
		_spec.SetField(service.FieldName, field.TypeString, value)
	}
	if value, ok := su.mutation.DisplayName(); ok {
		_spec.SetField(service.FieldDisplayName, field.TypeString, value)
	}
	if value, ok := su.mutation.Description(); ok {
		_spec.SetField(service.FieldDescription, field.TypeString, value)
	}
	if su.mutation.DescriptionCleared() {
		_spec.ClearField(service.FieldDescription, field.TypeString)
	}
	if value, ok := su.mutation.GetType(); ok {
		_spec.SetField(service.FieldType, field.TypeEnum, value)
	}
	if value, ok := su.mutation.Subtype(); ok {
		_spec.SetField(service.FieldSubtype, field.TypeEnum, value)
	}
	if value, ok := su.mutation.GitRepository(); ok {
		_spec.SetField(service.FieldGitRepository, field.TypeString, value)
	}
	if su.mutation.GitRepositoryCleared() {
		_spec.ClearField(service.FieldGitRepository, field.TypeString)
	}
	if value, ok := su.mutation.GitBranch(); ok {
		_spec.SetField(service.FieldGitBranch, field.TypeString, value)
	}
	if su.mutation.GitBranchCleared() {
		_spec.ClearField(service.FieldGitBranch, field.TypeString)
	}
	if su.mutation.ProjectCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.ProjectTable,
			Columns: []string{service.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.ProjectTable,
			Columns: []string{service.ProjectColumn},
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
	if su.mutation.GithubInstallationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.GithubInstallationTable,
			Columns: []string{service.GithubInstallationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(githubinstallation.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.GithubInstallationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.GithubInstallationTable,
			Columns: []string{service.GithubInstallationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(githubinstallation.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.ServiceConfigsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   service.ServiceConfigsTable,
			Columns: []string{service.ServiceConfigsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceconfig.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.ServiceConfigsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   service.ServiceConfigsTable,
			Columns: []string{service.ServiceConfigsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceconfig.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(su.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{service.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// ServiceUpdateOne is the builder for updating a single Service entity.
type ServiceUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ServiceMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (suo *ServiceUpdateOne) SetUpdatedAt(t time.Time) *ServiceUpdateOne {
	suo.mutation.SetUpdatedAt(t)
	return suo
}

// SetName sets the "name" field.
func (suo *ServiceUpdateOne) SetName(s string) *ServiceUpdateOne {
	suo.mutation.SetName(s)
	return suo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableName(s *string) *ServiceUpdateOne {
	if s != nil {
		suo.SetName(*s)
	}
	return suo
}

// SetDisplayName sets the "display_name" field.
func (suo *ServiceUpdateOne) SetDisplayName(s string) *ServiceUpdateOne {
	suo.mutation.SetDisplayName(s)
	return suo
}

// SetNillableDisplayName sets the "display_name" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableDisplayName(s *string) *ServiceUpdateOne {
	if s != nil {
		suo.SetDisplayName(*s)
	}
	return suo
}

// SetDescription sets the "description" field.
func (suo *ServiceUpdateOne) SetDescription(s string) *ServiceUpdateOne {
	suo.mutation.SetDescription(s)
	return suo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableDescription(s *string) *ServiceUpdateOne {
	if s != nil {
		suo.SetDescription(*s)
	}
	return suo
}

// ClearDescription clears the value of the "description" field.
func (suo *ServiceUpdateOne) ClearDescription() *ServiceUpdateOne {
	suo.mutation.ClearDescription()
	return suo
}

// SetType sets the "type" field.
func (suo *ServiceUpdateOne) SetType(s service.Type) *ServiceUpdateOne {
	suo.mutation.SetType(s)
	return suo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableType(s *service.Type) *ServiceUpdateOne {
	if s != nil {
		suo.SetType(*s)
	}
	return suo
}

// SetSubtype sets the "subtype" field.
func (suo *ServiceUpdateOne) SetSubtype(s service.Subtype) *ServiceUpdateOne {
	suo.mutation.SetSubtype(s)
	return suo
}

// SetNillableSubtype sets the "subtype" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableSubtype(s *service.Subtype) *ServiceUpdateOne {
	if s != nil {
		suo.SetSubtype(*s)
	}
	return suo
}

// SetProjectID sets the "project_id" field.
func (suo *ServiceUpdateOne) SetProjectID(u uuid.UUID) *ServiceUpdateOne {
	suo.mutation.SetProjectID(u)
	return suo
}

// SetNillableProjectID sets the "project_id" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableProjectID(u *uuid.UUID) *ServiceUpdateOne {
	if u != nil {
		suo.SetProjectID(*u)
	}
	return suo
}

// SetGithubInstallationID sets the "github_installation_id" field.
func (suo *ServiceUpdateOne) SetGithubInstallationID(i int64) *ServiceUpdateOne {
	suo.mutation.SetGithubInstallationID(i)
	return suo
}

// SetNillableGithubInstallationID sets the "github_installation_id" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableGithubInstallationID(i *int64) *ServiceUpdateOne {
	if i != nil {
		suo.SetGithubInstallationID(*i)
	}
	return suo
}

// ClearGithubInstallationID clears the value of the "github_installation_id" field.
func (suo *ServiceUpdateOne) ClearGithubInstallationID() *ServiceUpdateOne {
	suo.mutation.ClearGithubInstallationID()
	return suo
}

// SetGitRepository sets the "git_repository" field.
func (suo *ServiceUpdateOne) SetGitRepository(s string) *ServiceUpdateOne {
	suo.mutation.SetGitRepository(s)
	return suo
}

// SetNillableGitRepository sets the "git_repository" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableGitRepository(s *string) *ServiceUpdateOne {
	if s != nil {
		suo.SetGitRepository(*s)
	}
	return suo
}

// ClearGitRepository clears the value of the "git_repository" field.
func (suo *ServiceUpdateOne) ClearGitRepository() *ServiceUpdateOne {
	suo.mutation.ClearGitRepository()
	return suo
}

// SetGitBranch sets the "git_branch" field.
func (suo *ServiceUpdateOne) SetGitBranch(s string) *ServiceUpdateOne {
	suo.mutation.SetGitBranch(s)
	return suo
}

// SetNillableGitBranch sets the "git_branch" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableGitBranch(s *string) *ServiceUpdateOne {
	if s != nil {
		suo.SetGitBranch(*s)
	}
	return suo
}

// ClearGitBranch clears the value of the "git_branch" field.
func (suo *ServiceUpdateOne) ClearGitBranch() *ServiceUpdateOne {
	suo.mutation.ClearGitBranch()
	return suo
}

// SetProject sets the "project" edge to the Project entity.
func (suo *ServiceUpdateOne) SetProject(p *Project) *ServiceUpdateOne {
	return suo.SetProjectID(p.ID)
}

// SetGithubInstallation sets the "github_installation" edge to the GithubInstallation entity.
func (suo *ServiceUpdateOne) SetGithubInstallation(g *GithubInstallation) *ServiceUpdateOne {
	return suo.SetGithubInstallationID(g.ID)
}

// SetServiceConfigsID sets the "service_configs" edge to the ServiceConfig entity by ID.
func (suo *ServiceUpdateOne) SetServiceConfigsID(id uuid.UUID) *ServiceUpdateOne {
	suo.mutation.SetServiceConfigsID(id)
	return suo
}

// SetNillableServiceConfigsID sets the "service_configs" edge to the ServiceConfig entity by ID if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableServiceConfigsID(id *uuid.UUID) *ServiceUpdateOne {
	if id != nil {
		suo = suo.SetServiceConfigsID(*id)
	}
	return suo
}

// SetServiceConfigs sets the "service_configs" edge to the ServiceConfig entity.
func (suo *ServiceUpdateOne) SetServiceConfigs(s *ServiceConfig) *ServiceUpdateOne {
	return suo.SetServiceConfigsID(s.ID)
}

// Mutation returns the ServiceMutation object of the builder.
func (suo *ServiceUpdateOne) Mutation() *ServiceMutation {
	return suo.mutation
}

// ClearProject clears the "project" edge to the Project entity.
func (suo *ServiceUpdateOne) ClearProject() *ServiceUpdateOne {
	suo.mutation.ClearProject()
	return suo
}

// ClearGithubInstallation clears the "github_installation" edge to the GithubInstallation entity.
func (suo *ServiceUpdateOne) ClearGithubInstallation() *ServiceUpdateOne {
	suo.mutation.ClearGithubInstallation()
	return suo
}

// ClearServiceConfigs clears the "service_configs" edge to the ServiceConfig entity.
func (suo *ServiceUpdateOne) ClearServiceConfigs() *ServiceUpdateOne {
	suo.mutation.ClearServiceConfigs()
	return suo
}

// Where appends a list predicates to the ServiceUpdate builder.
func (suo *ServiceUpdateOne) Where(ps ...predicate.Service) *ServiceUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *ServiceUpdateOne) Select(field string, fields ...string) *ServiceUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Service entity.
func (suo *ServiceUpdateOne) Save(ctx context.Context) (*Service, error) {
	suo.defaults()
	return withHooks(ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *ServiceUpdateOne) SaveX(ctx context.Context) *Service {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *ServiceUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *ServiceUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (suo *ServiceUpdateOne) defaults() {
	if _, ok := suo.mutation.UpdatedAt(); !ok {
		v := service.UpdateDefaultUpdatedAt()
		suo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *ServiceUpdateOne) check() error {
	if v, ok := suo.mutation.Name(); ok {
		if err := service.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Service.name": %w`, err)}
		}
	}
	if v, ok := suo.mutation.GetType(); ok {
		if err := service.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Service.type": %w`, err)}
		}
	}
	if v, ok := suo.mutation.Subtype(); ok {
		if err := service.SubtypeValidator(v); err != nil {
			return &ValidationError{Name: "subtype", err: fmt.Errorf(`ent: validator failed for field "Service.subtype": %w`, err)}
		}
	}
	if suo.mutation.ProjectCleared() && len(suo.mutation.ProjectIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Service.project"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (suo *ServiceUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ServiceUpdateOne {
	suo.modifiers = append(suo.modifiers, modifiers...)
	return suo
}

func (suo *ServiceUpdateOne) sqlSave(ctx context.Context) (_node *Service, err error) {
	if err := suo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(service.Table, service.Columns, sqlgraph.NewFieldSpec(service.FieldID, field.TypeUUID))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Service.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, service.FieldID)
		for _, f := range fields {
			if !service.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != service.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.UpdatedAt(); ok {
		_spec.SetField(service.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := suo.mutation.Name(); ok {
		_spec.SetField(service.FieldName, field.TypeString, value)
	}
	if value, ok := suo.mutation.DisplayName(); ok {
		_spec.SetField(service.FieldDisplayName, field.TypeString, value)
	}
	if value, ok := suo.mutation.Description(); ok {
		_spec.SetField(service.FieldDescription, field.TypeString, value)
	}
	if suo.mutation.DescriptionCleared() {
		_spec.ClearField(service.FieldDescription, field.TypeString)
	}
	if value, ok := suo.mutation.GetType(); ok {
		_spec.SetField(service.FieldType, field.TypeEnum, value)
	}
	if value, ok := suo.mutation.Subtype(); ok {
		_spec.SetField(service.FieldSubtype, field.TypeEnum, value)
	}
	if value, ok := suo.mutation.GitRepository(); ok {
		_spec.SetField(service.FieldGitRepository, field.TypeString, value)
	}
	if suo.mutation.GitRepositoryCleared() {
		_spec.ClearField(service.FieldGitRepository, field.TypeString)
	}
	if value, ok := suo.mutation.GitBranch(); ok {
		_spec.SetField(service.FieldGitBranch, field.TypeString, value)
	}
	if suo.mutation.GitBranchCleared() {
		_spec.ClearField(service.FieldGitBranch, field.TypeString)
	}
	if suo.mutation.ProjectCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.ProjectTable,
			Columns: []string{service.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.ProjectTable,
			Columns: []string{service.ProjectColumn},
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
	if suo.mutation.GithubInstallationCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.GithubInstallationTable,
			Columns: []string{service.GithubInstallationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(githubinstallation.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.GithubInstallationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.GithubInstallationTable,
			Columns: []string{service.GithubInstallationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(githubinstallation.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.ServiceConfigsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   service.ServiceConfigsTable,
			Columns: []string{service.ServiceConfigsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceconfig.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.ServiceConfigsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   service.ServiceConfigsTable,
			Columns: []string{service.ServiceConfigsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceconfig.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(suo.modifiers...)
	_node = &Service{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{service.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}
