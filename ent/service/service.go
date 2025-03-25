// Code generated by ent, DO NOT EDIT.

package service

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the service type in the database.
	Label = "service"
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
	// FieldEnvironmentID holds the string denoting the environment_id field in the database.
	FieldEnvironmentID = "environment_id"
	// FieldGithubInstallationID holds the string denoting the github_installation_id field in the database.
	FieldGithubInstallationID = "github_installation_id"
	// FieldGitRepositoryOwner holds the string denoting the git_repository_owner field in the database.
	FieldGitRepositoryOwner = "git_repository_owner"
	// FieldGitRepository holds the string denoting the git_repository field in the database.
	FieldGitRepository = "git_repository"
	// FieldKubernetesSecret holds the string denoting the kubernetes_secret field in the database.
	FieldKubernetesSecret = "kubernetes_secret"
	// EdgeEnvironment holds the string denoting the environment edge name in mutations.
	EdgeEnvironment = "environment"
	// EdgeGithubInstallation holds the string denoting the github_installation edge name in mutations.
	EdgeGithubInstallation = "github_installation"
	// EdgeServiceConfig holds the string denoting the service_config edge name in mutations.
	EdgeServiceConfig = "service_config"
	// EdgeDeployments holds the string denoting the deployments edge name in mutations.
	EdgeDeployments = "deployments"
	// Table holds the table name of the service in the database.
	Table = "services"
	// EnvironmentTable is the table that holds the environment relation/edge.
	EnvironmentTable = "services"
	// EnvironmentInverseTable is the table name for the Environment entity.
	// It exists in this package in order to avoid circular dependency with the "environment" package.
	EnvironmentInverseTable = "environments"
	// EnvironmentColumn is the table column denoting the environment relation/edge.
	EnvironmentColumn = "environment_id"
	// GithubInstallationTable is the table that holds the github_installation relation/edge.
	GithubInstallationTable = "services"
	// GithubInstallationInverseTable is the table name for the GithubInstallation entity.
	// It exists in this package in order to avoid circular dependency with the "githubinstallation" package.
	GithubInstallationInverseTable = "github_installations"
	// GithubInstallationColumn is the table column denoting the github_installation relation/edge.
	GithubInstallationColumn = "github_installation_id"
	// ServiceConfigTable is the table that holds the service_config relation/edge.
	ServiceConfigTable = "service_configs"
	// ServiceConfigInverseTable is the table name for the ServiceConfig entity.
	// It exists in this package in order to avoid circular dependency with the "serviceconfig" package.
	ServiceConfigInverseTable = "service_configs"
	// ServiceConfigColumn is the table column denoting the service_config relation/edge.
	ServiceConfigColumn = "service_id"
	// DeploymentsTable is the table that holds the deployments relation/edge.
	DeploymentsTable = "deployments"
	// DeploymentsInverseTable is the table name for the Deployment entity.
	// It exists in this package in order to avoid circular dependency with the "deployment" package.
	DeploymentsInverseTable = "deployments"
	// DeploymentsColumn is the table column denoting the deployments relation/edge.
	DeploymentsColumn = "service_id"
)

// Columns holds all SQL columns for service fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldName,
	FieldDisplayName,
	FieldDescription,
	FieldEnvironmentID,
	FieldGithubInstallationID,
	FieldGitRepositoryOwner,
	FieldGitRepository,
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
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Service queries.
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

// ByEnvironmentID orders the results by the environment_id field.
func ByEnvironmentID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEnvironmentID, opts...).ToFunc()
}

// ByGithubInstallationID orders the results by the github_installation_id field.
func ByGithubInstallationID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGithubInstallationID, opts...).ToFunc()
}

// ByGitRepositoryOwner orders the results by the git_repository_owner field.
func ByGitRepositoryOwner(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGitRepositoryOwner, opts...).ToFunc()
}

// ByGitRepository orders the results by the git_repository field.
func ByGitRepository(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGitRepository, opts...).ToFunc()
}

// ByKubernetesSecret orders the results by the kubernetes_secret field.
func ByKubernetesSecret(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldKubernetesSecret, opts...).ToFunc()
}

// ByEnvironmentField orders the results by environment field.
func ByEnvironmentField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEnvironmentStep(), sql.OrderByField(field, opts...))
	}
}

// ByGithubInstallationField orders the results by github_installation field.
func ByGithubInstallationField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newGithubInstallationStep(), sql.OrderByField(field, opts...))
	}
}

// ByServiceConfigField orders the results by service_config field.
func ByServiceConfigField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newServiceConfigStep(), sql.OrderByField(field, opts...))
	}
}

// ByDeploymentsCount orders the results by deployments count.
func ByDeploymentsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newDeploymentsStep(), opts...)
	}
}

// ByDeployments orders the results by deployments terms.
func ByDeployments(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newDeploymentsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newEnvironmentStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EnvironmentInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, EnvironmentTable, EnvironmentColumn),
	)
}
func newGithubInstallationStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(GithubInstallationInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, GithubInstallationTable, GithubInstallationColumn),
	)
}
func newServiceConfigStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ServiceConfigInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, false, ServiceConfigTable, ServiceConfigColumn),
	)
}
func newDeploymentsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(DeploymentsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, DeploymentsTable, DeploymentsColumn),
	)
}
