// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/buildjob"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/service"
)

// BuildJob is the model entity for the BuildJob schema.
type BuildJob struct {
	config `json:"-"`
	// ID of the ent.
	// The primary key of the entity.
	ID uuid.UUID `json:"id"`
	// The time at which the entity was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// The time at which the entity was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// ServiceID holds the value of the "service_id" field.
	ServiceID uuid.UUID `json:"service_id,omitempty"`
	// Status holds the value of the "status" field.
	Status schema.BuildJobStatus `json:"status,omitempty"`
	// Error holds the value of the "error" field.
	Error string `json:"error,omitempty"`
	// StartedAt holds the value of the "started_at" field.
	StartedAt *time.Time `json:"started_at,omitempty"`
	// CompletedAt holds the value of the "completed_at" field.
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	// The name of the kubernetes job
	KubernetesJobName string `json:"kubernetes_job_name,omitempty"`
	// The status of the kubernetes job
	KubernetesJobStatus string `json:"kubernetes_job_status,omitempty"`
	// Attempts holds the value of the "attempts" field.
	Attempts int `json:"attempts,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the BuildJobQuery when eager-loading is set.
	Edges        BuildJobEdges `json:"edges"`
	selectValues sql.SelectValues
}

// BuildJobEdges holds the relations/edges for other nodes in the graph.
type BuildJobEdges struct {
	// Service holds the value of the service edge.
	Service *Service `json:"service,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ServiceOrErr returns the Service value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e BuildJobEdges) ServiceOrErr() (*Service, error) {
	if e.Service != nil {
		return e.Service, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: service.Label}
	}
	return nil, &NotLoadedError{edge: "service"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*BuildJob) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case buildjob.FieldAttempts:
			values[i] = new(sql.NullInt64)
		case buildjob.FieldStatus, buildjob.FieldError, buildjob.FieldKubernetesJobName, buildjob.FieldKubernetesJobStatus:
			values[i] = new(sql.NullString)
		case buildjob.FieldCreatedAt, buildjob.FieldUpdatedAt, buildjob.FieldStartedAt, buildjob.FieldCompletedAt:
			values[i] = new(sql.NullTime)
		case buildjob.FieldID, buildjob.FieldServiceID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the BuildJob fields.
func (bj *BuildJob) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case buildjob.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				bj.ID = *value
			}
		case buildjob.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				bj.CreatedAt = value.Time
			}
		case buildjob.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				bj.UpdatedAt = value.Time
			}
		case buildjob.FieldServiceID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field service_id", values[i])
			} else if value != nil {
				bj.ServiceID = *value
			}
		case buildjob.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				bj.Status = schema.BuildJobStatus(value.String)
			}
		case buildjob.FieldError:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field error", values[i])
			} else if value.Valid {
				bj.Error = value.String
			}
		case buildjob.FieldStartedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field started_at", values[i])
			} else if value.Valid {
				bj.StartedAt = new(time.Time)
				*bj.StartedAt = value.Time
			}
		case buildjob.FieldCompletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field completed_at", values[i])
			} else if value.Valid {
				bj.CompletedAt = new(time.Time)
				*bj.CompletedAt = value.Time
			}
		case buildjob.FieldKubernetesJobName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kubernetes_job_name", values[i])
			} else if value.Valid {
				bj.KubernetesJobName = value.String
			}
		case buildjob.FieldKubernetesJobStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kubernetes_job_status", values[i])
			} else if value.Valid {
				bj.KubernetesJobStatus = value.String
			}
		case buildjob.FieldAttempts:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field attempts", values[i])
			} else if value.Valid {
				bj.Attempts = int(value.Int64)
			}
		default:
			bj.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the BuildJob.
// This includes values selected through modifiers, order, etc.
func (bj *BuildJob) Value(name string) (ent.Value, error) {
	return bj.selectValues.Get(name)
}

// QueryService queries the "service" edge of the BuildJob entity.
func (bj *BuildJob) QueryService() *ServiceQuery {
	return NewBuildJobClient(bj.config).QueryService(bj)
}

// Update returns a builder for updating this BuildJob.
// Note that you need to call BuildJob.Unwrap() before calling this method if this BuildJob
// was returned from a transaction, and the transaction was committed or rolled back.
func (bj *BuildJob) Update() *BuildJobUpdateOne {
	return NewBuildJobClient(bj.config).UpdateOne(bj)
}

// Unwrap unwraps the BuildJob entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (bj *BuildJob) Unwrap() *BuildJob {
	_tx, ok := bj.config.driver.(*txDriver)
	if !ok {
		panic("ent: BuildJob is not a transactional entity")
	}
	bj.config.driver = _tx.drv
	return bj
}

// String implements the fmt.Stringer.
func (bj *BuildJob) String() string {
	var builder strings.Builder
	builder.WriteString("BuildJob(")
	builder.WriteString(fmt.Sprintf("id=%v, ", bj.ID))
	builder.WriteString("created_at=")
	builder.WriteString(bj.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(bj.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("service_id=")
	builder.WriteString(fmt.Sprintf("%v", bj.ServiceID))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", bj.Status))
	builder.WriteString(", ")
	builder.WriteString("error=")
	builder.WriteString(bj.Error)
	builder.WriteString(", ")
	if v := bj.StartedAt; v != nil {
		builder.WriteString("started_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := bj.CompletedAt; v != nil {
		builder.WriteString("completed_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("kubernetes_job_name=")
	builder.WriteString(bj.KubernetesJobName)
	builder.WriteString(", ")
	builder.WriteString("kubernetes_job_status=")
	builder.WriteString(bj.KubernetesJobStatus)
	builder.WriteString(", ")
	builder.WriteString("attempts=")
	builder.WriteString(fmt.Sprintf("%v", bj.Attempts))
	builder.WriteByte(')')
	return builder.String()
}

// BuildJobs is a parsable slice of BuildJob.
type BuildJobs []*BuildJob
