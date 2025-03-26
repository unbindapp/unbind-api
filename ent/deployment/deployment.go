// Code generated by ent, DO NOT EDIT.

package deployment

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
)

const (
	// Label holds the string label denoting the deployment type in the database.
	Label = "deployment"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldServiceID holds the string denoting the service_id field in the database.
	FieldServiceID = "service_id"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldSource holds the string denoting the source field in the database.
	FieldSource = "source"
	// FieldError holds the string denoting the error field in the database.
	FieldError = "error"
	// FieldCommitSha holds the string denoting the commit_sha field in the database.
	FieldCommitSha = "commit_sha"
	// FieldCommitMessage holds the string denoting the commit_message field in the database.
	FieldCommitMessage = "commit_message"
	// FieldCommitAuthor holds the string denoting the commit_author field in the database.
	FieldCommitAuthor = "commit_author"
	// FieldStartedAt holds the string denoting the started_at field in the database.
	FieldStartedAt = "started_at"
	// FieldCompletedAt holds the string denoting the completed_at field in the database.
	FieldCompletedAt = "completed_at"
	// FieldKubernetesJobName holds the string denoting the kubernetes_job_name field in the database.
	FieldKubernetesJobName = "kubernetes_job_name"
	// FieldKubernetesJobStatus holds the string denoting the kubernetes_job_status field in the database.
	FieldKubernetesJobStatus = "kubernetes_job_status"
	// FieldAttempts holds the string denoting the attempts field in the database.
	FieldAttempts = "attempts"
	// FieldImage holds the string denoting the image field in the database.
	FieldImage = "image"
	// EdgeService holds the string denoting the service edge name in mutations.
	EdgeService = "service"
	// Table holds the table name of the deployment in the database.
	Table = "deployments"
	// ServiceTable is the table that holds the service relation/edge.
	ServiceTable = "deployments"
	// ServiceInverseTable is the table name for the Service entity.
	// It exists in this package in order to avoid circular dependency with the "service" package.
	ServiceInverseTable = "services"
	// ServiceColumn is the table column denoting the service relation/edge.
	ServiceColumn = "service_id"
)

// Columns holds all SQL columns for deployment fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldServiceID,
	FieldStatus,
	FieldSource,
	FieldError,
	FieldCommitSha,
	FieldCommitMessage,
	FieldCommitAuthor,
	FieldStartedAt,
	FieldCompletedAt,
	FieldKubernetesJobName,
	FieldKubernetesJobStatus,
	FieldAttempts,
	FieldImage,
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
	// DefaultAttempts holds the default value on creation for the "attempts" field.
	DefaultAttempts int
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s schema.DeploymentStatus) error {
	switch s {
	case "queued", "building", "succeeded", "cancelled", "failed":
		return nil
	default:
		return fmt.Errorf("deployment: invalid enum value for status field: %q", s)
	}
}

const DefaultSource schema.DeploymentSource = "manual"

// SourceValidator is a validator for the "source" field enum values. It is called by the builders before save.
func SourceValidator(s schema.DeploymentSource) error {
	switch s {
	case "manual", "git":
		return nil
	default:
		return fmt.Errorf("deployment: invalid enum value for source field: %q", s)
	}
}

// OrderOption defines the ordering options for the Deployment queries.
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

// ByServiceID orders the results by the service_id field.
func ByServiceID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldServiceID, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// BySource orders the results by the source field.
func BySource(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSource, opts...).ToFunc()
}

// ByError orders the results by the error field.
func ByError(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldError, opts...).ToFunc()
}

// ByCommitSha orders the results by the commit_sha field.
func ByCommitSha(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCommitSha, opts...).ToFunc()
}

// ByCommitMessage orders the results by the commit_message field.
func ByCommitMessage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCommitMessage, opts...).ToFunc()
}

// ByStartedAt orders the results by the started_at field.
func ByStartedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStartedAt, opts...).ToFunc()
}

// ByCompletedAt orders the results by the completed_at field.
func ByCompletedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCompletedAt, opts...).ToFunc()
}

// ByKubernetesJobName orders the results by the kubernetes_job_name field.
func ByKubernetesJobName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldKubernetesJobName, opts...).ToFunc()
}

// ByKubernetesJobStatus orders the results by the kubernetes_job_status field.
func ByKubernetesJobStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldKubernetesJobStatus, opts...).ToFunc()
}

// ByAttempts orders the results by the attempts field.
func ByAttempts(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAttempts, opts...).ToFunc()
}

// ByImage orders the results by the image field.
func ByImage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldImage, opts...).ToFunc()
}

// ByServiceField orders the results by service field.
func ByServiceField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newServiceStep(), sql.OrderByField(field, opts...))
	}
}
func newServiceStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ServiceInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ServiceTable, ServiceColumn),
	)
}
