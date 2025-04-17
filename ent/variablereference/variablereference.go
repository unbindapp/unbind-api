// Code generated by ent, DO NOT EDIT.

package variablereference

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
)

const (
	// Label holds the string label denoting the variablereference type in the database.
	Label = "variable_reference"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldTargetServiceID holds the string denoting the target_service_id field in the database.
	FieldTargetServiceID = "target_service_id"
	// FieldTargetName holds the string denoting the target_name field in the database.
	FieldTargetName = "target_name"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldSourceType holds the string denoting the source_type field in the database.
	FieldSourceType = "source_type"
	// FieldSourceID holds the string denoting the source_id field in the database.
	FieldSourceID = "source_id"
	// FieldSourceName holds the string denoting the source_name field in the database.
	FieldSourceName = "source_name"
	// FieldSourceKey holds the string denoting the source_key field in the database.
	FieldSourceKey = "source_key"
	// FieldValueTemplate holds the string denoting the value_template field in the database.
	FieldValueTemplate = "value_template"
	// EdgeService holds the string denoting the service edge name in mutations.
	EdgeService = "service"
	// Table holds the table name of the variablereference in the database.
	Table = "variable_references"
	// ServiceTable is the table that holds the service relation/edge.
	ServiceTable = "variable_references"
	// ServiceInverseTable is the table name for the Service entity.
	// It exists in this package in order to avoid circular dependency with the "service" package.
	ServiceInverseTable = "services"
	// ServiceColumn is the table column denoting the service relation/edge.
	ServiceColumn = "target_service_id"
)

// Columns holds all SQL columns for variablereference fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldTargetServiceID,
	FieldTargetName,
	FieldType,
	FieldSourceType,
	FieldSourceID,
	FieldSourceName,
	FieldSourceKey,
	FieldValueTemplate,
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
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// TypeValidator is a validator for the "type" field enum values. It is called by the builders before save.
func TypeValidator(_type schema.VariableReferenceType) error {
	switch _type {
	case "variable", "external_endpoint", "internal_endpoint":
		return nil
	default:
		return fmt.Errorf("variablereference: invalid enum value for type field: %q", _type)
	}
}

// SourceTypeValidator is a validator for the "source_type" field enum values. It is called by the builders before save.
func SourceTypeValidator(st schema.VariableReferenceSourceType) error {
	switch st {
	case "team", "project", "environment", "service":
		return nil
	default:
		return fmt.Errorf("variablereference: invalid enum value for source_type field: %q", st)
	}
}

// OrderOption defines the ordering options for the VariableReference queries.
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

// ByTargetServiceID orders the results by the target_service_id field.
func ByTargetServiceID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTargetServiceID, opts...).ToFunc()
}

// ByTargetName orders the results by the target_name field.
func ByTargetName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTargetName, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// BySourceType orders the results by the source_type field.
func BySourceType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSourceType, opts...).ToFunc()
}

// BySourceID orders the results by the source_id field.
func BySourceID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSourceID, opts...).ToFunc()
}

// BySourceName orders the results by the source_name field.
func BySourceName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSourceName, opts...).ToFunc()
}

// BySourceKey orders the results by the source_key field.
func BySourceKey(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSourceKey, opts...).ToFunc()
}

// ByValueTemplate orders the results by the value_template field.
func ByValueTemplate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldValueTemplate, opts...).ToFunc()
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
