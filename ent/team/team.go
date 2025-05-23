// Code generated by ent, DO NOT EDIT.

package team

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the team type in the database.
	Label = "team"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldKubernetesName holds the string denoting the kubernetes_name field in the database.
	FieldKubernetesName = "kubernetes_name"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldNamespace holds the string denoting the namespace field in the database.
	FieldNamespace = "namespace"
	// FieldKubernetesSecret holds the string denoting the kubernetes_secret field in the database.
	FieldKubernetesSecret = "kubernetes_secret"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// EdgeProjects holds the string denoting the projects edge name in mutations.
	EdgeProjects = "projects"
	// EdgeS3Sources holds the string denoting the s3_sources edge name in mutations.
	EdgeS3Sources = "s3_sources"
	// EdgeMembers holds the string denoting the members edge name in mutations.
	EdgeMembers = "members"
	// EdgeTeamWebhooks holds the string denoting the team_webhooks edge name in mutations.
	EdgeTeamWebhooks = "team_webhooks"
	// Table holds the table name of the team in the database.
	Table = "teams"
	// ProjectsTable is the table that holds the projects relation/edge.
	ProjectsTable = "projects"
	// ProjectsInverseTable is the table name for the Project entity.
	// It exists in this package in order to avoid circular dependency with the "project" package.
	ProjectsInverseTable = "projects"
	// ProjectsColumn is the table column denoting the projects relation/edge.
	ProjectsColumn = "team_id"
	// S3SourcesTable is the table that holds the s3_sources relation/edge.
	S3SourcesTable = "s3_sources"
	// S3SourcesInverseTable is the table name for the S3 entity.
	// It exists in this package in order to avoid circular dependency with the "s3" package.
	S3SourcesInverseTable = "s3_sources"
	// S3SourcesColumn is the table column denoting the s3_sources relation/edge.
	S3SourcesColumn = "team_id"
	// MembersTable is the table that holds the members relation/edge. The primary key declared below.
	MembersTable = "user_teams"
	// MembersInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	MembersInverseTable = "users"
	// TeamWebhooksTable is the table that holds the team_webhooks relation/edge.
	TeamWebhooksTable = "webhooks"
	// TeamWebhooksInverseTable is the table name for the Webhook entity.
	// It exists in this package in order to avoid circular dependency with the "webhook" package.
	TeamWebhooksInverseTable = "webhooks"
	// TeamWebhooksColumn is the table column denoting the team_webhooks relation/edge.
	TeamWebhooksColumn = "team_id"
)

// Columns holds all SQL columns for team fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldKubernetesName,
	FieldName,
	FieldNamespace,
	FieldKubernetesSecret,
	FieldDescription,
}

var (
	// MembersPrimaryKey and MembersColumn2 are the table columns denoting the
	// primary key for the members relation (M2M).
	MembersPrimaryKey = []string{"user_id", "team_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// KubernetesNameValidator is a validator for the "kubernetes_name" field. It is called by the builders before save.
	KubernetesNameValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Team queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByKubernetesName orders the results by the kubernetes_name field.
func ByKubernetesName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldKubernetesName, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByNamespace orders the results by the namespace field.
func ByNamespace(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNamespace, opts...).ToFunc()
}

// ByKubernetesSecret orders the results by the kubernetes_secret field.
func ByKubernetesSecret(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldKubernetesSecret, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByProjectsCount orders the results by projects count.
func ByProjectsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newProjectsStep(), opts...)
	}
}

// ByProjects orders the results by projects terms.
func ByProjects(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProjectsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByS3SourcesCount orders the results by s3_sources count.
func ByS3SourcesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newS3SourcesStep(), opts...)
	}
}

// ByS3Sources orders the results by s3_sources terms.
func ByS3Sources(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newS3SourcesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByMembersCount orders the results by members count.
func ByMembersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMembersStep(), opts...)
	}
}

// ByMembers orders the results by members terms.
func ByMembers(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMembersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByTeamWebhooksCount orders the results by team_webhooks count.
func ByTeamWebhooksCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newTeamWebhooksStep(), opts...)
	}
}

// ByTeamWebhooks orders the results by team_webhooks terms.
func ByTeamWebhooks(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTeamWebhooksStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newProjectsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProjectsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ProjectsTable, ProjectsColumn),
	)
}
func newS3SourcesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(S3SourcesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, S3SourcesTable, S3SourcesColumn),
	)
}
func newMembersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MembersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, MembersTable, MembersPrimaryKey...),
	)
}
func newTeamWebhooksStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TeamWebhooksInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, TeamWebhooksTable, TeamWebhooksColumn),
	)
}
