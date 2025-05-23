// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/registry"
)

// Registry is the model entity for the Registry schema.
type Registry struct {
	config `json:"-"`
	// ID of the ent.
	// The primary key of the entity.
	ID uuid.UUID `json:"id"`
	// The time at which the entity was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// The time at which the entity was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Host holds the value of the "host" field.
	Host string `json:"host,omitempty"`
	// The name of the kubernetes registry credentials secret, should be located in the unbind system namespace
	KubernetesSecret string `json:"kubernetes_secret,omitempty"`
	// If true, this is the registry that will be used for internal CI/CD
	IsDefault    bool `json:"is_default,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Registry) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case registry.FieldIsDefault:
			values[i] = new(sql.NullBool)
		case registry.FieldHost, registry.FieldKubernetesSecret:
			values[i] = new(sql.NullString)
		case registry.FieldCreatedAt, registry.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case registry.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Registry fields.
func (r *Registry) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case registry.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				r.ID = *value
			}
		case registry.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				r.CreatedAt = value.Time
			}
		case registry.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				r.UpdatedAt = value.Time
			}
		case registry.FieldHost:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field host", values[i])
			} else if value.Valid {
				r.Host = value.String
			}
		case registry.FieldKubernetesSecret:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kubernetes_secret", values[i])
			} else if value.Valid {
				r.KubernetesSecret = value.String
			}
		case registry.FieldIsDefault:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_default", values[i])
			} else if value.Valid {
				r.IsDefault = value.Bool
			}
		default:
			r.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Registry.
// This includes values selected through modifiers, order, etc.
func (r *Registry) Value(name string) (ent.Value, error) {
	return r.selectValues.Get(name)
}

// Update returns a builder for updating this Registry.
// Note that you need to call Registry.Unwrap() before calling this method if this Registry
// was returned from a transaction, and the transaction was committed or rolled back.
func (r *Registry) Update() *RegistryUpdateOne {
	return NewRegistryClient(r.config).UpdateOne(r)
}

// Unwrap unwraps the Registry entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (r *Registry) Unwrap() *Registry {
	_tx, ok := r.config.driver.(*txDriver)
	if !ok {
		panic("ent: Registry is not a transactional entity")
	}
	r.config.driver = _tx.drv
	return r
}

// String implements the fmt.Stringer.
func (r *Registry) String() string {
	var builder strings.Builder
	builder.WriteString("Registry(")
	builder.WriteString(fmt.Sprintf("id=%v, ", r.ID))
	builder.WriteString("created_at=")
	builder.WriteString(r.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(r.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("host=")
	builder.WriteString(r.Host)
	builder.WriteString(", ")
	builder.WriteString("kubernetes_secret=")
	builder.WriteString(r.KubernetesSecret)
	builder.WriteString(", ")
	builder.WriteString("is_default=")
	builder.WriteString(fmt.Sprintf("%v", r.IsDefault))
	builder.WriteByte(')')
	return builder.String()
}

// Registries is a parsable slice of Registry.
type Registries []*Registry
