// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/variablereference"
)

// VariableReference is the model entity for the VariableReference schema.
type VariableReference struct {
	config `json:"-"`
	// ID of the ent.
	// The primary key of the entity.
	ID uuid.UUID `json:"id"`
	// The time at which the entity was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// The time at which the entity was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// TargetServiceID holds the value of the "target_service_id" field.
	TargetServiceID uuid.UUID `json:"target_service_id,omitempty"`
	// TargetName holds the value of the "target_name" field.
	TargetName string `json:"target_name,omitempty"`
	// List of sources for this variable reference, interpolated as ${sourcename.sourcekey}
	Sources []schema.VariableReferenceSource `json:"sources,omitempty"`
	// Optional template for the value, e.g. 'Hello ${a.b} this is my variable ${c.d}'
	ValueTemplate string `json:"value_template,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the VariableReferenceQuery when eager-loading is set.
	Edges        VariableReferenceEdges `json:"edges"`
	selectValues sql.SelectValues
}

// VariableReferenceEdges holds the relations/edges for other nodes in the graph.
type VariableReferenceEdges struct {
	// Service that this variable reference points to
	Service *Service `json:"service,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ServiceOrErr returns the Service value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e VariableReferenceEdges) ServiceOrErr() (*Service, error) {
	if e.Service != nil {
		return e.Service, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: service.Label}
	}
	return nil, &NotLoadedError{edge: "service"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*VariableReference) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case variablereference.FieldSources:
			values[i] = new([]byte)
		case variablereference.FieldTargetName, variablereference.FieldValueTemplate:
			values[i] = new(sql.NullString)
		case variablereference.FieldCreatedAt, variablereference.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case variablereference.FieldID, variablereference.FieldTargetServiceID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the VariableReference fields.
func (vr *VariableReference) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case variablereference.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				vr.ID = *value
			}
		case variablereference.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				vr.CreatedAt = value.Time
			}
		case variablereference.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				vr.UpdatedAt = value.Time
			}
		case variablereference.FieldTargetServiceID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field target_service_id", values[i])
			} else if value != nil {
				vr.TargetServiceID = *value
			}
		case variablereference.FieldTargetName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field target_name", values[i])
			} else if value.Valid {
				vr.TargetName = value.String
			}
		case variablereference.FieldSources:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field sources", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &vr.Sources); err != nil {
					return fmt.Errorf("unmarshal field sources: %w", err)
				}
			}
		case variablereference.FieldValueTemplate:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field value_template", values[i])
			} else if value.Valid {
				vr.ValueTemplate = value.String
			}
		default:
			vr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the VariableReference.
// This includes values selected through modifiers, order, etc.
func (vr *VariableReference) Value(name string) (ent.Value, error) {
	return vr.selectValues.Get(name)
}

// QueryService queries the "service" edge of the VariableReference entity.
func (vr *VariableReference) QueryService() *ServiceQuery {
	return NewVariableReferenceClient(vr.config).QueryService(vr)
}

// Update returns a builder for updating this VariableReference.
// Note that you need to call VariableReference.Unwrap() before calling this method if this VariableReference
// was returned from a transaction, and the transaction was committed or rolled back.
func (vr *VariableReference) Update() *VariableReferenceUpdateOne {
	return NewVariableReferenceClient(vr.config).UpdateOne(vr)
}

// Unwrap unwraps the VariableReference entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (vr *VariableReference) Unwrap() *VariableReference {
	_tx, ok := vr.config.driver.(*txDriver)
	if !ok {
		panic("ent: VariableReference is not a transactional entity")
	}
	vr.config.driver = _tx.drv
	return vr
}

// String implements the fmt.Stringer.
func (vr *VariableReference) String() string {
	var builder strings.Builder
	builder.WriteString("VariableReference(")
	builder.WriteString(fmt.Sprintf("id=%v, ", vr.ID))
	builder.WriteString("created_at=")
	builder.WriteString(vr.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(vr.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("target_service_id=")
	builder.WriteString(fmt.Sprintf("%v", vr.TargetServiceID))
	builder.WriteString(", ")
	builder.WriteString("target_name=")
	builder.WriteString(vr.TargetName)
	builder.WriteString(", ")
	builder.WriteString("sources=")
	builder.WriteString(fmt.Sprintf("%v", vr.Sources))
	builder.WriteString(", ")
	builder.WriteString("value_template=")
	builder.WriteString(vr.ValueTemplate)
	builder.WriteByte(')')
	return builder.String()
}

// VariableReferences is a parsable slice of VariableReference.
type VariableReferences []*VariableReference
