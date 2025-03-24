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
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
)

// ServiceConfig is the model entity for the ServiceConfig schema.
type ServiceConfig struct {
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
	// Type of service
	Type schema.ServiceType `json:"type,omitempty"`
	// Builder holds the value of the "builder" field.
	Builder schema.ServiceBuilder `json:"builder,omitempty"`
	// Provider (e.g. Go, Python, Node, Deno)
	Provider *enum.Provider `json:"provider,omitempty"`
	// Framework of service - corresponds mostly to railpack results - e.g. Django, Next, Express, Gin
	Framework *enum.Framework `json:"framework,omitempty"`
	// Branch to build from
	GitBranch *string `json:"git_branch,omitempty"`
	// External domains and paths for the service
	Hosts []schema.HostSpec `json:"hosts,omitempty"`
	// Container ports to expose
	Ports []schema.PortSpec `json:"ports,omitempty"`
	// Number of replicas for the service
	Replicas int32 `json:"replicas,omitempty"`
	// Whether to automatically deploy on git push
	AutoDeploy bool `json:"auto_deploy,omitempty"`
	// Custom run command
	RunCommand *string `json:"run_command,omitempty"`
	// Whether the service is publicly accessible, creates an ingress resource
	Public bool `json:"public,omitempty"`
	// Custom Docker image if not building from git
	Image string `json:"image,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ServiceConfigQuery when eager-loading is set.
	Edges        ServiceConfigEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ServiceConfigEdges holds the relations/edges for other nodes in the graph.
type ServiceConfigEdges struct {
	// Service holds the value of the service edge.
	Service *Service `json:"service,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ServiceOrErr returns the Service value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceConfigEdges) ServiceOrErr() (*Service, error) {
	if e.Service != nil {
		return e.Service, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: service.Label}
	}
	return nil, &NotLoadedError{edge: "service"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ServiceConfig) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case serviceconfig.FieldHosts, serviceconfig.FieldPorts:
			values[i] = new([]byte)
		case serviceconfig.FieldAutoDeploy, serviceconfig.FieldPublic:
			values[i] = new(sql.NullBool)
		case serviceconfig.FieldReplicas:
			values[i] = new(sql.NullInt64)
		case serviceconfig.FieldType, serviceconfig.FieldBuilder, serviceconfig.FieldProvider, serviceconfig.FieldFramework, serviceconfig.FieldGitBranch, serviceconfig.FieldRunCommand, serviceconfig.FieldImage:
			values[i] = new(sql.NullString)
		case serviceconfig.FieldCreatedAt, serviceconfig.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case serviceconfig.FieldID, serviceconfig.FieldServiceID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ServiceConfig fields.
func (sc *ServiceConfig) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case serviceconfig.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				sc.ID = *value
			}
		case serviceconfig.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				sc.CreatedAt = value.Time
			}
		case serviceconfig.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				sc.UpdatedAt = value.Time
			}
		case serviceconfig.FieldServiceID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field service_id", values[i])
			} else if value != nil {
				sc.ServiceID = *value
			}
		case serviceconfig.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				sc.Type = schema.ServiceType(value.String)
			}
		case serviceconfig.FieldBuilder:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field builder", values[i])
			} else if value.Valid {
				sc.Builder = schema.ServiceBuilder(value.String)
			}
		case serviceconfig.FieldProvider:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field provider", values[i])
			} else if value.Valid {
				sc.Provider = new(enum.Provider)
				*sc.Provider = enum.Provider(value.String)
			}
		case serviceconfig.FieldFramework:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field framework", values[i])
			} else if value.Valid {
				sc.Framework = new(enum.Framework)
				*sc.Framework = enum.Framework(value.String)
			}
		case serviceconfig.FieldGitBranch:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field git_branch", values[i])
			} else if value.Valid {
				sc.GitBranch = new(string)
				*sc.GitBranch = value.String
			}
		case serviceconfig.FieldHosts:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field hosts", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &sc.Hosts); err != nil {
					return fmt.Errorf("unmarshal field hosts: %w", err)
				}
			}
		case serviceconfig.FieldPorts:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field ports", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &sc.Ports); err != nil {
					return fmt.Errorf("unmarshal field ports: %w", err)
				}
			}
		case serviceconfig.FieldReplicas:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field replicas", values[i])
			} else if value.Valid {
				sc.Replicas = int32(value.Int64)
			}
		case serviceconfig.FieldAutoDeploy:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field auto_deploy", values[i])
			} else if value.Valid {
				sc.AutoDeploy = value.Bool
			}
		case serviceconfig.FieldRunCommand:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field run_command", values[i])
			} else if value.Valid {
				sc.RunCommand = new(string)
				*sc.RunCommand = value.String
			}
		case serviceconfig.FieldPublic:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field public", values[i])
			} else if value.Valid {
				sc.Public = value.Bool
			}
		case serviceconfig.FieldImage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field image", values[i])
			} else if value.Valid {
				sc.Image = value.String
			}
		default:
			sc.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ServiceConfig.
// This includes values selected through modifiers, order, etc.
func (sc *ServiceConfig) Value(name string) (ent.Value, error) {
	return sc.selectValues.Get(name)
}

// QueryService queries the "service" edge of the ServiceConfig entity.
func (sc *ServiceConfig) QueryService() *ServiceQuery {
	return NewServiceConfigClient(sc.config).QueryService(sc)
}

// Update returns a builder for updating this ServiceConfig.
// Note that you need to call ServiceConfig.Unwrap() before calling this method if this ServiceConfig
// was returned from a transaction, and the transaction was committed or rolled back.
func (sc *ServiceConfig) Update() *ServiceConfigUpdateOne {
	return NewServiceConfigClient(sc.config).UpdateOne(sc)
}

// Unwrap unwraps the ServiceConfig entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (sc *ServiceConfig) Unwrap() *ServiceConfig {
	_tx, ok := sc.config.driver.(*txDriver)
	if !ok {
		panic("ent: ServiceConfig is not a transactional entity")
	}
	sc.config.driver = _tx.drv
	return sc
}

// String implements the fmt.Stringer.
func (sc *ServiceConfig) String() string {
	var builder strings.Builder
	builder.WriteString("ServiceConfig(")
	builder.WriteString(fmt.Sprintf("id=%v, ", sc.ID))
	builder.WriteString("created_at=")
	builder.WriteString(sc.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(sc.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("service_id=")
	builder.WriteString(fmt.Sprintf("%v", sc.ServiceID))
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(fmt.Sprintf("%v", sc.Type))
	builder.WriteString(", ")
	builder.WriteString("builder=")
	builder.WriteString(fmt.Sprintf("%v", sc.Builder))
	builder.WriteString(", ")
	if v := sc.Provider; v != nil {
		builder.WriteString("provider=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := sc.Framework; v != nil {
		builder.WriteString("framework=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := sc.GitBranch; v != nil {
		builder.WriteString("git_branch=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("hosts=")
	builder.WriteString(fmt.Sprintf("%v", sc.Hosts))
	builder.WriteString(", ")
	builder.WriteString("ports=")
	builder.WriteString(fmt.Sprintf("%v", sc.Ports))
	builder.WriteString(", ")
	builder.WriteString("replicas=")
	builder.WriteString(fmt.Sprintf("%v", sc.Replicas))
	builder.WriteString(", ")
	builder.WriteString("auto_deploy=")
	builder.WriteString(fmt.Sprintf("%v", sc.AutoDeploy))
	builder.WriteString(", ")
	if v := sc.RunCommand; v != nil {
		builder.WriteString("run_command=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("public=")
	builder.WriteString(fmt.Sprintf("%v", sc.Public))
	builder.WriteString(", ")
	builder.WriteString("image=")
	builder.WriteString(sc.Image)
	builder.WriteByte(')')
	return builder.String()
}

// ServiceConfigs is a parsable slice of ServiceConfig.
type ServiceConfigs []*ServiceConfig
