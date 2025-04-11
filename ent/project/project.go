// Code generated by ent, DO NOT EDIT.

package project

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the project type in the database.
	Label = "project"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDisplayName holds the string denoting the display_name field in the database.
	FieldDisplayName = "display_name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldTeamID holds the string denoting the team_id field in the database.
	FieldTeamID = "team_id"
	// FieldDefaultEnvironmentID holds the string denoting the default_environment_id field in the database.
	FieldDefaultEnvironmentID = "default_environment_id"
	// FieldKubernetesSecret holds the string denoting the kubernetes_secret field in the database.
	FieldKubernetesSecret = "kubernetes_secret"
	// EdgeTeam holds the string denoting the team edge name in mutations.
	EdgeTeam = "team"
	// EdgeEnvironments holds the string denoting the environments edge name in mutations.
	EdgeEnvironments = "environments"
	// EdgeDefaultEnvironment holds the string denoting the default_environment edge name in mutations.
	EdgeDefaultEnvironment = "default_environment"
	// Table holds the table name of the project in the database.
	Table = "projects"
	// TeamTable is the table that holds the team relation/edge.
	TeamTable = "projects"
	// TeamInverseTable is the table name for the Team entity.
	// It exists in this package in order to avoid circular dependency with the "team" package.
	TeamInverseTable = "teams"
	// TeamColumn is the table column denoting the team relation/edge.
	TeamColumn = "team_id"
	// EnvironmentsTable is the table that holds the environments relation/edge.
	EnvironmentsTable = "environments"
	// EnvironmentsInverseTable is the table name for the Environment entity.
	// It exists in this package in order to avoid circular dependency with the "environment" package.
	EnvironmentsInverseTable = "environments"
	// EnvironmentsColumn is the table column denoting the environments relation/edge.
	EnvironmentsColumn = "project_id"
	// DefaultEnvironmentTable is the table that holds the default_environment relation/edge.
	DefaultEnvironmentTable = "projects"
	// DefaultEnvironmentInverseTable is the table name for the Environment entity.
	// It exists in this package in order to avoid circular dependency with the "environment" package.
	DefaultEnvironmentInverseTable = "environments"
	// DefaultEnvironmentColumn is the table column denoting the default_environment relation/edge.
	DefaultEnvironmentColumn = "default_environment_id"
)

// Columns holds all SQL columns for project fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldName,
	FieldDisplayName,
	FieldDescription,
	FieldStatus,
	FieldTeamID,
	FieldDefaultEnvironmentID,
	FieldKubernetesSecret,
}

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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultStatus holds the default value on creation for the "status" field.
	DefaultStatus string
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Project queries.
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

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDisplayName orders the results by the display_name field.
func ByDisplayName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDisplayName, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByTeamID orders the results by the team_id field.
func ByTeamID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTeamID, opts...).ToFunc()
}

// ByDefaultEnvironmentID orders the results by the default_environment_id field.
func ByDefaultEnvironmentID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDefaultEnvironmentID, opts...).ToFunc()
}

// ByKubernetesSecret orders the results by the kubernetes_secret field.
func ByKubernetesSecret(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldKubernetesSecret, opts...).ToFunc()
}

// ByTeamField orders the results by team field.
func ByTeamField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTeamStep(), sql.OrderByField(field, opts...))
	}
}

// ByEnvironmentsCount orders the results by environments count.
func ByEnvironmentsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newEnvironmentsStep(), opts...)
	}
}

// ByEnvironments orders the results by environments terms.
func ByEnvironments(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEnvironmentsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByDefaultEnvironmentField orders the results by default_environment field.
func ByDefaultEnvironmentField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newDefaultEnvironmentStep(), sql.OrderByField(field, opts...))
	}
}
func newTeamStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TeamInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, TeamTable, TeamColumn),
	)
}
func newEnvironmentsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EnvironmentsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, EnvironmentsTable, EnvironmentsColumn),
	)
}
func newDefaultEnvironmentStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(DefaultEnvironmentInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, DefaultEnvironmentTable, DefaultEnvironmentColumn),
	)
}
