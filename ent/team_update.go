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
	"github.com/unbindapp/unbind-api/ent/project"
	"github.com/unbindapp/unbind-api/ent/s3"
	"github.com/unbindapp/unbind-api/ent/team"
	"github.com/unbindapp/unbind-api/ent/user"
	"github.com/unbindapp/unbind-api/ent/webhook"
)

// TeamUpdate is the builder for updating Team entities.
type TeamUpdate struct {
	config
	hooks     []Hook
	mutation  *TeamMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the TeamUpdate builder.
func (tu *TeamUpdate) Where(ps ...predicate.Team) *TeamUpdate {
	tu.mutation.Where(ps...)
	return tu
}

// SetUpdatedAt sets the "updated_at" field.
func (tu *TeamUpdate) SetUpdatedAt(t time.Time) *TeamUpdate {
	tu.mutation.SetUpdatedAt(t)
	return tu
}

// SetKubernetesName sets the "kubernetes_name" field.
func (tu *TeamUpdate) SetKubernetesName(s string) *TeamUpdate {
	tu.mutation.SetKubernetesName(s)
	return tu
}

// SetNillableKubernetesName sets the "kubernetes_name" field if the given value is not nil.
func (tu *TeamUpdate) SetNillableKubernetesName(s *string) *TeamUpdate {
	if s != nil {
		tu.SetKubernetesName(*s)
	}
	return tu
}

// SetName sets the "name" field.
func (tu *TeamUpdate) SetName(s string) *TeamUpdate {
	tu.mutation.SetName(s)
	return tu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tu *TeamUpdate) SetNillableName(s *string) *TeamUpdate {
	if s != nil {
		tu.SetName(*s)
	}
	return tu
}

// SetNamespace sets the "namespace" field.
func (tu *TeamUpdate) SetNamespace(s string) *TeamUpdate {
	tu.mutation.SetNamespace(s)
	return tu
}

// SetNillableNamespace sets the "namespace" field if the given value is not nil.
func (tu *TeamUpdate) SetNillableNamespace(s *string) *TeamUpdate {
	if s != nil {
		tu.SetNamespace(*s)
	}
	return tu
}

// SetKubernetesSecret sets the "kubernetes_secret" field.
func (tu *TeamUpdate) SetKubernetesSecret(s string) *TeamUpdate {
	tu.mutation.SetKubernetesSecret(s)
	return tu
}

// SetNillableKubernetesSecret sets the "kubernetes_secret" field if the given value is not nil.
func (tu *TeamUpdate) SetNillableKubernetesSecret(s *string) *TeamUpdate {
	if s != nil {
		tu.SetKubernetesSecret(*s)
	}
	return tu
}

// SetDescription sets the "description" field.
func (tu *TeamUpdate) SetDescription(s string) *TeamUpdate {
	tu.mutation.SetDescription(s)
	return tu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (tu *TeamUpdate) SetNillableDescription(s *string) *TeamUpdate {
	if s != nil {
		tu.SetDescription(*s)
	}
	return tu
}

// ClearDescription clears the value of the "description" field.
func (tu *TeamUpdate) ClearDescription() *TeamUpdate {
	tu.mutation.ClearDescription()
	return tu
}

// AddProjectIDs adds the "projects" edge to the Project entity by IDs.
func (tu *TeamUpdate) AddProjectIDs(ids ...uuid.UUID) *TeamUpdate {
	tu.mutation.AddProjectIDs(ids...)
	return tu
}

// AddProjects adds the "projects" edges to the Project entity.
func (tu *TeamUpdate) AddProjects(p ...*Project) *TeamUpdate {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return tu.AddProjectIDs(ids...)
}

// AddS3SourceIDs adds the "s3_sources" edge to the S3 entity by IDs.
func (tu *TeamUpdate) AddS3SourceIDs(ids ...uuid.UUID) *TeamUpdate {
	tu.mutation.AddS3SourceIDs(ids...)
	return tu
}

// AddS3Sources adds the "s3_sources" edges to the S3 entity.
func (tu *TeamUpdate) AddS3Sources(s ...*S3) *TeamUpdate {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return tu.AddS3SourceIDs(ids...)
}

// AddMemberIDs adds the "members" edge to the User entity by IDs.
func (tu *TeamUpdate) AddMemberIDs(ids ...uuid.UUID) *TeamUpdate {
	tu.mutation.AddMemberIDs(ids...)
	return tu
}

// AddMembers adds the "members" edges to the User entity.
func (tu *TeamUpdate) AddMembers(u ...*User) *TeamUpdate {
	ids := make([]uuid.UUID, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return tu.AddMemberIDs(ids...)
}

// AddTeamWebhookIDs adds the "team_webhooks" edge to the Webhook entity by IDs.
func (tu *TeamUpdate) AddTeamWebhookIDs(ids ...uuid.UUID) *TeamUpdate {
	tu.mutation.AddTeamWebhookIDs(ids...)
	return tu
}

// AddTeamWebhooks adds the "team_webhooks" edges to the Webhook entity.
func (tu *TeamUpdate) AddTeamWebhooks(w ...*Webhook) *TeamUpdate {
	ids := make([]uuid.UUID, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tu.AddTeamWebhookIDs(ids...)
}

// Mutation returns the TeamMutation object of the builder.
func (tu *TeamUpdate) Mutation() *TeamMutation {
	return tu.mutation
}

// ClearProjects clears all "projects" edges to the Project entity.
func (tu *TeamUpdate) ClearProjects() *TeamUpdate {
	tu.mutation.ClearProjects()
	return tu
}

// RemoveProjectIDs removes the "projects" edge to Project entities by IDs.
func (tu *TeamUpdate) RemoveProjectIDs(ids ...uuid.UUID) *TeamUpdate {
	tu.mutation.RemoveProjectIDs(ids...)
	return tu
}

// RemoveProjects removes "projects" edges to Project entities.
func (tu *TeamUpdate) RemoveProjects(p ...*Project) *TeamUpdate {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return tu.RemoveProjectIDs(ids...)
}

// ClearS3Sources clears all "s3_sources" edges to the S3 entity.
func (tu *TeamUpdate) ClearS3Sources() *TeamUpdate {
	tu.mutation.ClearS3Sources()
	return tu
}

// RemoveS3SourceIDs removes the "s3_sources" edge to S3 entities by IDs.
func (tu *TeamUpdate) RemoveS3SourceIDs(ids ...uuid.UUID) *TeamUpdate {
	tu.mutation.RemoveS3SourceIDs(ids...)
	return tu
}

// RemoveS3Sources removes "s3_sources" edges to S3 entities.
func (tu *TeamUpdate) RemoveS3Sources(s ...*S3) *TeamUpdate {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return tu.RemoveS3SourceIDs(ids...)
}

// ClearMembers clears all "members" edges to the User entity.
func (tu *TeamUpdate) ClearMembers() *TeamUpdate {
	tu.mutation.ClearMembers()
	return tu
}

// RemoveMemberIDs removes the "members" edge to User entities by IDs.
func (tu *TeamUpdate) RemoveMemberIDs(ids ...uuid.UUID) *TeamUpdate {
	tu.mutation.RemoveMemberIDs(ids...)
	return tu
}

// RemoveMembers removes "members" edges to User entities.
func (tu *TeamUpdate) RemoveMembers(u ...*User) *TeamUpdate {
	ids := make([]uuid.UUID, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return tu.RemoveMemberIDs(ids...)
}

// ClearTeamWebhooks clears all "team_webhooks" edges to the Webhook entity.
func (tu *TeamUpdate) ClearTeamWebhooks() *TeamUpdate {
	tu.mutation.ClearTeamWebhooks()
	return tu
}

// RemoveTeamWebhookIDs removes the "team_webhooks" edge to Webhook entities by IDs.
func (tu *TeamUpdate) RemoveTeamWebhookIDs(ids ...uuid.UUID) *TeamUpdate {
	tu.mutation.RemoveTeamWebhookIDs(ids...)
	return tu
}

// RemoveTeamWebhooks removes "team_webhooks" edges to Webhook entities.
func (tu *TeamUpdate) RemoveTeamWebhooks(w ...*Webhook) *TeamUpdate {
	ids := make([]uuid.UUID, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tu.RemoveTeamWebhookIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tu *TeamUpdate) Save(ctx context.Context) (int, error) {
	tu.defaults()
	return withHooks(ctx, tu.sqlSave, tu.mutation, tu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TeamUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TeamUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TeamUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tu *TeamUpdate) defaults() {
	if _, ok := tu.mutation.UpdatedAt(); !ok {
		v := team.UpdateDefaultUpdatedAt()
		tu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tu *TeamUpdate) check() error {
	if v, ok := tu.mutation.KubernetesName(); ok {
		if err := team.KubernetesNameValidator(v); err != nil {
			return &ValidationError{Name: "kubernetes_name", err: fmt.Errorf(`ent: validator failed for field "Team.kubernetes_name": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (tu *TeamUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *TeamUpdate {
	tu.modifiers = append(tu.modifiers, modifiers...)
	return tu
}

func (tu *TeamUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := tu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(team.Table, team.Columns, sqlgraph.NewFieldSpec(team.FieldID, field.TypeUUID))
	if ps := tu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.UpdatedAt(); ok {
		_spec.SetField(team.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := tu.mutation.KubernetesName(); ok {
		_spec.SetField(team.FieldKubernetesName, field.TypeString, value)
	}
	if value, ok := tu.mutation.Name(); ok {
		_spec.SetField(team.FieldName, field.TypeString, value)
	}
	if value, ok := tu.mutation.Namespace(); ok {
		_spec.SetField(team.FieldNamespace, field.TypeString, value)
	}
	if value, ok := tu.mutation.KubernetesSecret(); ok {
		_spec.SetField(team.FieldKubernetesSecret, field.TypeString, value)
	}
	if value, ok := tu.mutation.Description(); ok {
		_spec.SetField(team.FieldDescription, field.TypeString, value)
	}
	if tu.mutation.DescriptionCleared() {
		_spec.ClearField(team.FieldDescription, field.TypeString)
	}
	if tu.mutation.ProjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.ProjectsTable,
			Columns: []string{team.ProjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.RemovedProjectsIDs(); len(nodes) > 0 && !tu.mutation.ProjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.ProjectsTable,
			Columns: []string{team.ProjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.ProjectsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.ProjectsTable,
			Columns: []string{team.ProjectsColumn},
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
	if tu.mutation.S3SourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.S3SourcesTable,
			Columns: []string{team.S3SourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(s3.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.RemovedS3SourcesIDs(); len(nodes) > 0 && !tu.mutation.S3SourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.S3SourcesTable,
			Columns: []string{team.S3SourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(s3.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.S3SourcesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.S3SourcesTable,
			Columns: []string{team.S3SourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(s3.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tu.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.MembersTable,
			Columns: team.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.RemovedMembersIDs(); len(nodes) > 0 && !tu.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.MembersTable,
			Columns: team.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.MembersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.MembersTable,
			Columns: team.MembersPrimaryKey,
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
	if tu.mutation.TeamWebhooksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamWebhooksTable,
			Columns: []string{team.TeamWebhooksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(webhook.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.RemovedTeamWebhooksIDs(); len(nodes) > 0 && !tu.mutation.TeamWebhooksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamWebhooksTable,
			Columns: []string{team.TeamWebhooksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(webhook.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.TeamWebhooksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamWebhooksTable,
			Columns: []string{team.TeamWebhooksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(webhook.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(tu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{team.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tu.mutation.done = true
	return n, nil
}

// TeamUpdateOne is the builder for updating a single Team entity.
type TeamUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *TeamMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdatedAt sets the "updated_at" field.
func (tuo *TeamUpdateOne) SetUpdatedAt(t time.Time) *TeamUpdateOne {
	tuo.mutation.SetUpdatedAt(t)
	return tuo
}

// SetKubernetesName sets the "kubernetes_name" field.
func (tuo *TeamUpdateOne) SetKubernetesName(s string) *TeamUpdateOne {
	tuo.mutation.SetKubernetesName(s)
	return tuo
}

// SetNillableKubernetesName sets the "kubernetes_name" field if the given value is not nil.
func (tuo *TeamUpdateOne) SetNillableKubernetesName(s *string) *TeamUpdateOne {
	if s != nil {
		tuo.SetKubernetesName(*s)
	}
	return tuo
}

// SetName sets the "name" field.
func (tuo *TeamUpdateOne) SetName(s string) *TeamUpdateOne {
	tuo.mutation.SetName(s)
	return tuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tuo *TeamUpdateOne) SetNillableName(s *string) *TeamUpdateOne {
	if s != nil {
		tuo.SetName(*s)
	}
	return tuo
}

// SetNamespace sets the "namespace" field.
func (tuo *TeamUpdateOne) SetNamespace(s string) *TeamUpdateOne {
	tuo.mutation.SetNamespace(s)
	return tuo
}

// SetNillableNamespace sets the "namespace" field if the given value is not nil.
func (tuo *TeamUpdateOne) SetNillableNamespace(s *string) *TeamUpdateOne {
	if s != nil {
		tuo.SetNamespace(*s)
	}
	return tuo
}

// SetKubernetesSecret sets the "kubernetes_secret" field.
func (tuo *TeamUpdateOne) SetKubernetesSecret(s string) *TeamUpdateOne {
	tuo.mutation.SetKubernetesSecret(s)
	return tuo
}

// SetNillableKubernetesSecret sets the "kubernetes_secret" field if the given value is not nil.
func (tuo *TeamUpdateOne) SetNillableKubernetesSecret(s *string) *TeamUpdateOne {
	if s != nil {
		tuo.SetKubernetesSecret(*s)
	}
	return tuo
}

// SetDescription sets the "description" field.
func (tuo *TeamUpdateOne) SetDescription(s string) *TeamUpdateOne {
	tuo.mutation.SetDescription(s)
	return tuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (tuo *TeamUpdateOne) SetNillableDescription(s *string) *TeamUpdateOne {
	if s != nil {
		tuo.SetDescription(*s)
	}
	return tuo
}

// ClearDescription clears the value of the "description" field.
func (tuo *TeamUpdateOne) ClearDescription() *TeamUpdateOne {
	tuo.mutation.ClearDescription()
	return tuo
}

// AddProjectIDs adds the "projects" edge to the Project entity by IDs.
func (tuo *TeamUpdateOne) AddProjectIDs(ids ...uuid.UUID) *TeamUpdateOne {
	tuo.mutation.AddProjectIDs(ids...)
	return tuo
}

// AddProjects adds the "projects" edges to the Project entity.
func (tuo *TeamUpdateOne) AddProjects(p ...*Project) *TeamUpdateOne {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return tuo.AddProjectIDs(ids...)
}

// AddS3SourceIDs adds the "s3_sources" edge to the S3 entity by IDs.
func (tuo *TeamUpdateOne) AddS3SourceIDs(ids ...uuid.UUID) *TeamUpdateOne {
	tuo.mutation.AddS3SourceIDs(ids...)
	return tuo
}

// AddS3Sources adds the "s3_sources" edges to the S3 entity.
func (tuo *TeamUpdateOne) AddS3Sources(s ...*S3) *TeamUpdateOne {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return tuo.AddS3SourceIDs(ids...)
}

// AddMemberIDs adds the "members" edge to the User entity by IDs.
func (tuo *TeamUpdateOne) AddMemberIDs(ids ...uuid.UUID) *TeamUpdateOne {
	tuo.mutation.AddMemberIDs(ids...)
	return tuo
}

// AddMembers adds the "members" edges to the User entity.
func (tuo *TeamUpdateOne) AddMembers(u ...*User) *TeamUpdateOne {
	ids := make([]uuid.UUID, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return tuo.AddMemberIDs(ids...)
}

// AddTeamWebhookIDs adds the "team_webhooks" edge to the Webhook entity by IDs.
func (tuo *TeamUpdateOne) AddTeamWebhookIDs(ids ...uuid.UUID) *TeamUpdateOne {
	tuo.mutation.AddTeamWebhookIDs(ids...)
	return tuo
}

// AddTeamWebhooks adds the "team_webhooks" edges to the Webhook entity.
func (tuo *TeamUpdateOne) AddTeamWebhooks(w ...*Webhook) *TeamUpdateOne {
	ids := make([]uuid.UUID, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tuo.AddTeamWebhookIDs(ids...)
}

// Mutation returns the TeamMutation object of the builder.
func (tuo *TeamUpdateOne) Mutation() *TeamMutation {
	return tuo.mutation
}

// ClearProjects clears all "projects" edges to the Project entity.
func (tuo *TeamUpdateOne) ClearProjects() *TeamUpdateOne {
	tuo.mutation.ClearProjects()
	return tuo
}

// RemoveProjectIDs removes the "projects" edge to Project entities by IDs.
func (tuo *TeamUpdateOne) RemoveProjectIDs(ids ...uuid.UUID) *TeamUpdateOne {
	tuo.mutation.RemoveProjectIDs(ids...)
	return tuo
}

// RemoveProjects removes "projects" edges to Project entities.
func (tuo *TeamUpdateOne) RemoveProjects(p ...*Project) *TeamUpdateOne {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return tuo.RemoveProjectIDs(ids...)
}

// ClearS3Sources clears all "s3_sources" edges to the S3 entity.
func (tuo *TeamUpdateOne) ClearS3Sources() *TeamUpdateOne {
	tuo.mutation.ClearS3Sources()
	return tuo
}

// RemoveS3SourceIDs removes the "s3_sources" edge to S3 entities by IDs.
func (tuo *TeamUpdateOne) RemoveS3SourceIDs(ids ...uuid.UUID) *TeamUpdateOne {
	tuo.mutation.RemoveS3SourceIDs(ids...)
	return tuo
}

// RemoveS3Sources removes "s3_sources" edges to S3 entities.
func (tuo *TeamUpdateOne) RemoveS3Sources(s ...*S3) *TeamUpdateOne {
	ids := make([]uuid.UUID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return tuo.RemoveS3SourceIDs(ids...)
}

// ClearMembers clears all "members" edges to the User entity.
func (tuo *TeamUpdateOne) ClearMembers() *TeamUpdateOne {
	tuo.mutation.ClearMembers()
	return tuo
}

// RemoveMemberIDs removes the "members" edge to User entities by IDs.
func (tuo *TeamUpdateOne) RemoveMemberIDs(ids ...uuid.UUID) *TeamUpdateOne {
	tuo.mutation.RemoveMemberIDs(ids...)
	return tuo
}

// RemoveMembers removes "members" edges to User entities.
func (tuo *TeamUpdateOne) RemoveMembers(u ...*User) *TeamUpdateOne {
	ids := make([]uuid.UUID, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return tuo.RemoveMemberIDs(ids...)
}

// ClearTeamWebhooks clears all "team_webhooks" edges to the Webhook entity.
func (tuo *TeamUpdateOne) ClearTeamWebhooks() *TeamUpdateOne {
	tuo.mutation.ClearTeamWebhooks()
	return tuo
}

// RemoveTeamWebhookIDs removes the "team_webhooks" edge to Webhook entities by IDs.
func (tuo *TeamUpdateOne) RemoveTeamWebhookIDs(ids ...uuid.UUID) *TeamUpdateOne {
	tuo.mutation.RemoveTeamWebhookIDs(ids...)
	return tuo
}

// RemoveTeamWebhooks removes "team_webhooks" edges to Webhook entities.
func (tuo *TeamUpdateOne) RemoveTeamWebhooks(w ...*Webhook) *TeamUpdateOne {
	ids := make([]uuid.UUID, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tuo.RemoveTeamWebhookIDs(ids...)
}

// Where appends a list predicates to the TeamUpdate builder.
func (tuo *TeamUpdateOne) Where(ps ...predicate.Team) *TeamUpdateOne {
	tuo.mutation.Where(ps...)
	return tuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tuo *TeamUpdateOne) Select(field string, fields ...string) *TeamUpdateOne {
	tuo.fields = append([]string{field}, fields...)
	return tuo
}

// Save executes the query and returns the updated Team entity.
func (tuo *TeamUpdateOne) Save(ctx context.Context) (*Team, error) {
	tuo.defaults()
	return withHooks(ctx, tuo.sqlSave, tuo.mutation, tuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TeamUpdateOne) SaveX(ctx context.Context) *Team {
	node, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tuo *TeamUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TeamUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tuo *TeamUpdateOne) defaults() {
	if _, ok := tuo.mutation.UpdatedAt(); !ok {
		v := team.UpdateDefaultUpdatedAt()
		tuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tuo *TeamUpdateOne) check() error {
	if v, ok := tuo.mutation.KubernetesName(); ok {
		if err := team.KubernetesNameValidator(v); err != nil {
			return &ValidationError{Name: "kubernetes_name", err: fmt.Errorf(`ent: validator failed for field "Team.kubernetes_name": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (tuo *TeamUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *TeamUpdateOne {
	tuo.modifiers = append(tuo.modifiers, modifiers...)
	return tuo
}

func (tuo *TeamUpdateOne) sqlSave(ctx context.Context) (_node *Team, err error) {
	if err := tuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(team.Table, team.Columns, sqlgraph.NewFieldSpec(team.FieldID, field.TypeUUID))
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Team.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, team.FieldID)
		for _, f := range fields {
			if !team.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != team.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tuo.mutation.UpdatedAt(); ok {
		_spec.SetField(team.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := tuo.mutation.KubernetesName(); ok {
		_spec.SetField(team.FieldKubernetesName, field.TypeString, value)
	}
	if value, ok := tuo.mutation.Name(); ok {
		_spec.SetField(team.FieldName, field.TypeString, value)
	}
	if value, ok := tuo.mutation.Namespace(); ok {
		_spec.SetField(team.FieldNamespace, field.TypeString, value)
	}
	if value, ok := tuo.mutation.KubernetesSecret(); ok {
		_spec.SetField(team.FieldKubernetesSecret, field.TypeString, value)
	}
	if value, ok := tuo.mutation.Description(); ok {
		_spec.SetField(team.FieldDescription, field.TypeString, value)
	}
	if tuo.mutation.DescriptionCleared() {
		_spec.ClearField(team.FieldDescription, field.TypeString)
	}
	if tuo.mutation.ProjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.ProjectsTable,
			Columns: []string{team.ProjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.RemovedProjectsIDs(); len(nodes) > 0 && !tuo.mutation.ProjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.ProjectsTable,
			Columns: []string{team.ProjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(project.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.ProjectsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.ProjectsTable,
			Columns: []string{team.ProjectsColumn},
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
	if tuo.mutation.S3SourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.S3SourcesTable,
			Columns: []string{team.S3SourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(s3.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.RemovedS3SourcesIDs(); len(nodes) > 0 && !tuo.mutation.S3SourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.S3SourcesTable,
			Columns: []string{team.S3SourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(s3.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.S3SourcesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.S3SourcesTable,
			Columns: []string{team.S3SourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(s3.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tuo.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.MembersTable,
			Columns: team.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.RemovedMembersIDs(); len(nodes) > 0 && !tuo.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.MembersTable,
			Columns: team.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.MembersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.MembersTable,
			Columns: team.MembersPrimaryKey,
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
	if tuo.mutation.TeamWebhooksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamWebhooksTable,
			Columns: []string{team.TeamWebhooksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(webhook.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.RemovedTeamWebhooksIDs(); len(nodes) > 0 && !tuo.mutation.TeamWebhooksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamWebhooksTable,
			Columns: []string{team.TeamWebhooksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(webhook.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.TeamWebhooksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamWebhooksTable,
			Columns: []string{team.TeamWebhooksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(webhook.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(tuo.modifiers...)
	_node = &Team{config: tuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{team.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tuo.mutation.done = true
	return _node, nil
}
