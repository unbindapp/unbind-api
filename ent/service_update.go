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
	"github.com/unbindapp/unbind-api/ent/buildjob"
	"github.com/unbindapp/unbind-api/ent/environment"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/ent/predicate"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
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

// SetBuilder sets the "builder" field.
func (su *ServiceUpdate) SetBuilder(s service.Builder) *ServiceUpdate {
	su.mutation.SetBuilder(s)
	return su
}

// SetNillableBuilder sets the "builder" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableBuilder(s *service.Builder) *ServiceUpdate {
	if s != nil {
		su.SetBuilder(*s)
	}
	return su
}

// SetProvider sets the "provider" field.
func (su *ServiceUpdate) SetProvider(e enum.Provider) *ServiceUpdate {
	su.mutation.SetProvider(e)
	return su
}

// SetNillableProvider sets the "provider" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableProvider(e *enum.Provider) *ServiceUpdate {
	if e != nil {
		su.SetProvider(*e)
	}
	return su
}

// ClearProvider clears the value of the "provider" field.
func (su *ServiceUpdate) ClearProvider() *ServiceUpdate {
	su.mutation.ClearProvider()
	return su
}

// SetFramework sets the "framework" field.
func (su *ServiceUpdate) SetFramework(e enum.Framework) *ServiceUpdate {
	su.mutation.SetFramework(e)
	return su
}

// SetNillableFramework sets the "framework" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableFramework(e *enum.Framework) *ServiceUpdate {
	if e != nil {
		su.SetFramework(*e)
	}
	return su
}

// ClearFramework clears the value of the "framework" field.
func (su *ServiceUpdate) ClearFramework() *ServiceUpdate {
	su.mutation.ClearFramework()
	return su
}

// SetEnvironmentID sets the "environment_id" field.
func (su *ServiceUpdate) SetEnvironmentID(u uuid.UUID) *ServiceUpdate {
	su.mutation.SetEnvironmentID(u)
	return su
}

// SetNillableEnvironmentID sets the "environment_id" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableEnvironmentID(u *uuid.UUID) *ServiceUpdate {
	if u != nil {
		su.SetEnvironmentID(*u)
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

// SetKubernetesSecret sets the "kubernetes_secret" field.
func (su *ServiceUpdate) SetKubernetesSecret(s string) *ServiceUpdate {
	su.mutation.SetKubernetesSecret(s)
	return su
}

// SetNillableKubernetesSecret sets the "kubernetes_secret" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableKubernetesSecret(s *string) *ServiceUpdate {
	if s != nil {
		su.SetKubernetesSecret(*s)
	}
	return su
}

// SetEnvironment sets the "environment" edge to the Environment entity.
func (su *ServiceUpdate) SetEnvironment(e *Environment) *ServiceUpdate {
	return su.SetEnvironmentID(e.ID)
}

// SetGithubInstallation sets the "github_installation" edge to the GithubInstallation entity.
func (su *ServiceUpdate) SetGithubInstallation(g *GithubInstallation) *ServiceUpdate {
	return su.SetGithubInstallationID(g.ID)
}

// SetServiceConfigID sets the "service_config" edge to the ServiceConfig entity by ID.
func (su *ServiceUpdate) SetServiceConfigID(id uuid.UUID) *ServiceUpdate {
	su.mutation.SetServiceConfigID(id)
	return su
}

// SetNillableServiceConfigID sets the "service_config" edge to the ServiceConfig entity by ID if the given value is not nil.
func (su *ServiceUpdate) SetNillableServiceConfigID(id *uuid.UUID) *ServiceUpdate {
	if id != nil {
		su = su.SetServiceConfigID(*id)
	}
	return su
}

// SetServiceConfig sets the "service_config" edge to the ServiceConfig entity.
func (su *ServiceUpdate) SetServiceConfig(s *ServiceConfig) *ServiceUpdate {
	return su.SetServiceConfigID(s.ID)
}

// AddBuildJobIDs adds the "build_jobs" edge to the BuildJob entity by IDs.
func (su *ServiceUpdate) AddBuildJobIDs(ids ...uuid.UUID) *ServiceUpdate {
	su.mutation.AddBuildJobIDs(ids...)
	return su
}

// AddBuildJobs adds the "build_jobs" edges to the BuildJob entity.
func (su *ServiceUpdate) AddBuildJobs(b ...*BuildJob) *ServiceUpdate {
	ids := make([]uuid.UUID, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return su.AddBuildJobIDs(ids...)
}

// Mutation returns the ServiceMutation object of the builder.
func (su *ServiceUpdate) Mutation() *ServiceMutation {
	return su.mutation
}

// ClearEnvironment clears the "environment" edge to the Environment entity.
func (su *ServiceUpdate) ClearEnvironment() *ServiceUpdate {
	su.mutation.ClearEnvironment()
	return su
}

// ClearGithubInstallation clears the "github_installation" edge to the GithubInstallation entity.
func (su *ServiceUpdate) ClearGithubInstallation() *ServiceUpdate {
	su.mutation.ClearGithubInstallation()
	return su
}

// ClearServiceConfig clears the "service_config" edge to the ServiceConfig entity.
func (su *ServiceUpdate) ClearServiceConfig() *ServiceUpdate {
	su.mutation.ClearServiceConfig()
	return su
}

// ClearBuildJobs clears all "build_jobs" edges to the BuildJob entity.
func (su *ServiceUpdate) ClearBuildJobs() *ServiceUpdate {
	su.mutation.ClearBuildJobs()
	return su
}

// RemoveBuildJobIDs removes the "build_jobs" edge to BuildJob entities by IDs.
func (su *ServiceUpdate) RemoveBuildJobIDs(ids ...uuid.UUID) *ServiceUpdate {
	su.mutation.RemoveBuildJobIDs(ids...)
	return su
}

// RemoveBuildJobs removes "build_jobs" edges to BuildJob entities.
func (su *ServiceUpdate) RemoveBuildJobs(b ...*BuildJob) *ServiceUpdate {
	ids := make([]uuid.UUID, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return su.RemoveBuildJobIDs(ids...)
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
	if v, ok := su.mutation.Builder(); ok {
		if err := service.BuilderValidator(v); err != nil {
			return &ValidationError{Name: "builder", err: fmt.Errorf(`ent: validator failed for field "Service.builder": %w`, err)}
		}
	}
	if v, ok := su.mutation.Provider(); ok {
		if err := service.ProviderValidator(v); err != nil {
			return &ValidationError{Name: "provider", err: fmt.Errorf(`ent: validator failed for field "Service.provider": %w`, err)}
		}
	}
	if v, ok := su.mutation.Framework(); ok {
		if err := service.FrameworkValidator(v); err != nil {
			return &ValidationError{Name: "framework", err: fmt.Errorf(`ent: validator failed for field "Service.framework": %w`, err)}
		}
	}
	if su.mutation.EnvironmentCleared() && len(su.mutation.EnvironmentIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Service.environment"`)
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
	if value, ok := su.mutation.Builder(); ok {
		_spec.SetField(service.FieldBuilder, field.TypeEnum, value)
	}
	if value, ok := su.mutation.Provider(); ok {
		_spec.SetField(service.FieldProvider, field.TypeEnum, value)
	}
	if su.mutation.ProviderCleared() {
		_spec.ClearField(service.FieldProvider, field.TypeEnum)
	}
	if value, ok := su.mutation.Framework(); ok {
		_spec.SetField(service.FieldFramework, field.TypeEnum, value)
	}
	if su.mutation.FrameworkCleared() {
		_spec.ClearField(service.FieldFramework, field.TypeEnum)
	}
	if value, ok := su.mutation.GitRepository(); ok {
		_spec.SetField(service.FieldGitRepository, field.TypeString, value)
	}
	if su.mutation.GitRepositoryCleared() {
		_spec.ClearField(service.FieldGitRepository, field.TypeString)
	}
	if value, ok := su.mutation.KubernetesSecret(); ok {
		_spec.SetField(service.FieldKubernetesSecret, field.TypeString, value)
	}
	if su.mutation.EnvironmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.EnvironmentTable,
			Columns: []string{service.EnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(environment.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.EnvironmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.EnvironmentTable,
			Columns: []string{service.EnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(environment.FieldID, field.TypeUUID),
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
	if su.mutation.ServiceConfigCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   service.ServiceConfigTable,
			Columns: []string{service.ServiceConfigColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceconfig.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.ServiceConfigIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   service.ServiceConfigTable,
			Columns: []string{service.ServiceConfigColumn},
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
	if su.mutation.BuildJobsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.BuildJobsTable,
			Columns: []string{service.BuildJobsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildjob.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedBuildJobsIDs(); len(nodes) > 0 && !su.mutation.BuildJobsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.BuildJobsTable,
			Columns: []string{service.BuildJobsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildjob.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.BuildJobsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.BuildJobsTable,
			Columns: []string{service.BuildJobsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildjob.FieldID, field.TypeUUID),
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

// SetBuilder sets the "builder" field.
func (suo *ServiceUpdateOne) SetBuilder(s service.Builder) *ServiceUpdateOne {
	suo.mutation.SetBuilder(s)
	return suo
}

// SetNillableBuilder sets the "builder" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableBuilder(s *service.Builder) *ServiceUpdateOne {
	if s != nil {
		suo.SetBuilder(*s)
	}
	return suo
}

// SetProvider sets the "provider" field.
func (suo *ServiceUpdateOne) SetProvider(e enum.Provider) *ServiceUpdateOne {
	suo.mutation.SetProvider(e)
	return suo
}

// SetNillableProvider sets the "provider" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableProvider(e *enum.Provider) *ServiceUpdateOne {
	if e != nil {
		suo.SetProvider(*e)
	}
	return suo
}

// ClearProvider clears the value of the "provider" field.
func (suo *ServiceUpdateOne) ClearProvider() *ServiceUpdateOne {
	suo.mutation.ClearProvider()
	return suo
}

// SetFramework sets the "framework" field.
func (suo *ServiceUpdateOne) SetFramework(e enum.Framework) *ServiceUpdateOne {
	suo.mutation.SetFramework(e)
	return suo
}

// SetNillableFramework sets the "framework" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableFramework(e *enum.Framework) *ServiceUpdateOne {
	if e != nil {
		suo.SetFramework(*e)
	}
	return suo
}

// ClearFramework clears the value of the "framework" field.
func (suo *ServiceUpdateOne) ClearFramework() *ServiceUpdateOne {
	suo.mutation.ClearFramework()
	return suo
}

// SetEnvironmentID sets the "environment_id" field.
func (suo *ServiceUpdateOne) SetEnvironmentID(u uuid.UUID) *ServiceUpdateOne {
	suo.mutation.SetEnvironmentID(u)
	return suo
}

// SetNillableEnvironmentID sets the "environment_id" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableEnvironmentID(u *uuid.UUID) *ServiceUpdateOne {
	if u != nil {
		suo.SetEnvironmentID(*u)
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

// SetKubernetesSecret sets the "kubernetes_secret" field.
func (suo *ServiceUpdateOne) SetKubernetesSecret(s string) *ServiceUpdateOne {
	suo.mutation.SetKubernetesSecret(s)
	return suo
}

// SetNillableKubernetesSecret sets the "kubernetes_secret" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableKubernetesSecret(s *string) *ServiceUpdateOne {
	if s != nil {
		suo.SetKubernetesSecret(*s)
	}
	return suo
}

// SetEnvironment sets the "environment" edge to the Environment entity.
func (suo *ServiceUpdateOne) SetEnvironment(e *Environment) *ServiceUpdateOne {
	return suo.SetEnvironmentID(e.ID)
}

// SetGithubInstallation sets the "github_installation" edge to the GithubInstallation entity.
func (suo *ServiceUpdateOne) SetGithubInstallation(g *GithubInstallation) *ServiceUpdateOne {
	return suo.SetGithubInstallationID(g.ID)
}

// SetServiceConfigID sets the "service_config" edge to the ServiceConfig entity by ID.
func (suo *ServiceUpdateOne) SetServiceConfigID(id uuid.UUID) *ServiceUpdateOne {
	suo.mutation.SetServiceConfigID(id)
	return suo
}

// SetNillableServiceConfigID sets the "service_config" edge to the ServiceConfig entity by ID if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableServiceConfigID(id *uuid.UUID) *ServiceUpdateOne {
	if id != nil {
		suo = suo.SetServiceConfigID(*id)
	}
	return suo
}

// SetServiceConfig sets the "service_config" edge to the ServiceConfig entity.
func (suo *ServiceUpdateOne) SetServiceConfig(s *ServiceConfig) *ServiceUpdateOne {
	return suo.SetServiceConfigID(s.ID)
}

// AddBuildJobIDs adds the "build_jobs" edge to the BuildJob entity by IDs.
func (suo *ServiceUpdateOne) AddBuildJobIDs(ids ...uuid.UUID) *ServiceUpdateOne {
	suo.mutation.AddBuildJobIDs(ids...)
	return suo
}

// AddBuildJobs adds the "build_jobs" edges to the BuildJob entity.
func (suo *ServiceUpdateOne) AddBuildJobs(b ...*BuildJob) *ServiceUpdateOne {
	ids := make([]uuid.UUID, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return suo.AddBuildJobIDs(ids...)
}

// Mutation returns the ServiceMutation object of the builder.
func (suo *ServiceUpdateOne) Mutation() *ServiceMutation {
	return suo.mutation
}

// ClearEnvironment clears the "environment" edge to the Environment entity.
func (suo *ServiceUpdateOne) ClearEnvironment() *ServiceUpdateOne {
	suo.mutation.ClearEnvironment()
	return suo
}

// ClearGithubInstallation clears the "github_installation" edge to the GithubInstallation entity.
func (suo *ServiceUpdateOne) ClearGithubInstallation() *ServiceUpdateOne {
	suo.mutation.ClearGithubInstallation()
	return suo
}

// ClearServiceConfig clears the "service_config" edge to the ServiceConfig entity.
func (suo *ServiceUpdateOne) ClearServiceConfig() *ServiceUpdateOne {
	suo.mutation.ClearServiceConfig()
	return suo
}

// ClearBuildJobs clears all "build_jobs" edges to the BuildJob entity.
func (suo *ServiceUpdateOne) ClearBuildJobs() *ServiceUpdateOne {
	suo.mutation.ClearBuildJobs()
	return suo
}

// RemoveBuildJobIDs removes the "build_jobs" edge to BuildJob entities by IDs.
func (suo *ServiceUpdateOne) RemoveBuildJobIDs(ids ...uuid.UUID) *ServiceUpdateOne {
	suo.mutation.RemoveBuildJobIDs(ids...)
	return suo
}

// RemoveBuildJobs removes "build_jobs" edges to BuildJob entities.
func (suo *ServiceUpdateOne) RemoveBuildJobs(b ...*BuildJob) *ServiceUpdateOne {
	ids := make([]uuid.UUID, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return suo.RemoveBuildJobIDs(ids...)
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
	if v, ok := suo.mutation.Builder(); ok {
		if err := service.BuilderValidator(v); err != nil {
			return &ValidationError{Name: "builder", err: fmt.Errorf(`ent: validator failed for field "Service.builder": %w`, err)}
		}
	}
	if v, ok := suo.mutation.Provider(); ok {
		if err := service.ProviderValidator(v); err != nil {
			return &ValidationError{Name: "provider", err: fmt.Errorf(`ent: validator failed for field "Service.provider": %w`, err)}
		}
	}
	if v, ok := suo.mutation.Framework(); ok {
		if err := service.FrameworkValidator(v); err != nil {
			return &ValidationError{Name: "framework", err: fmt.Errorf(`ent: validator failed for field "Service.framework": %w`, err)}
		}
	}
	if suo.mutation.EnvironmentCleared() && len(suo.mutation.EnvironmentIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "Service.environment"`)
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
	if value, ok := suo.mutation.Builder(); ok {
		_spec.SetField(service.FieldBuilder, field.TypeEnum, value)
	}
	if value, ok := suo.mutation.Provider(); ok {
		_spec.SetField(service.FieldProvider, field.TypeEnum, value)
	}
	if suo.mutation.ProviderCleared() {
		_spec.ClearField(service.FieldProvider, field.TypeEnum)
	}
	if value, ok := suo.mutation.Framework(); ok {
		_spec.SetField(service.FieldFramework, field.TypeEnum, value)
	}
	if suo.mutation.FrameworkCleared() {
		_spec.ClearField(service.FieldFramework, field.TypeEnum)
	}
	if value, ok := suo.mutation.GitRepository(); ok {
		_spec.SetField(service.FieldGitRepository, field.TypeString, value)
	}
	if suo.mutation.GitRepositoryCleared() {
		_spec.ClearField(service.FieldGitRepository, field.TypeString)
	}
	if value, ok := suo.mutation.KubernetesSecret(); ok {
		_spec.SetField(service.FieldKubernetesSecret, field.TypeString, value)
	}
	if suo.mutation.EnvironmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.EnvironmentTable,
			Columns: []string{service.EnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(environment.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.EnvironmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   service.EnvironmentTable,
			Columns: []string{service.EnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(environment.FieldID, field.TypeUUID),
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
	if suo.mutation.ServiceConfigCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   service.ServiceConfigTable,
			Columns: []string{service.ServiceConfigColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceconfig.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.ServiceConfigIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   service.ServiceConfigTable,
			Columns: []string{service.ServiceConfigColumn},
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
	if suo.mutation.BuildJobsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.BuildJobsTable,
			Columns: []string{service.BuildJobsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildjob.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedBuildJobsIDs(); len(nodes) > 0 && !suo.mutation.BuildJobsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.BuildJobsTable,
			Columns: []string{service.BuildJobsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildjob.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.BuildJobsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.BuildJobsTable,
			Columns: []string{service.BuildJobsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(buildjob.FieldID, field.TypeUUID),
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
