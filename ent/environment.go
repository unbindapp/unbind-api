// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/environment"
	"github.com/unbindapp/unbind-api/ent/project"
)

// Environment is the model entity for the Environment schema.
type Environment struct {
	config `json:"-"`
	// ID of the ent.
	// The primary key of the entity.
	ID uuid.UUID `json:"id"`
	// The time at which the entity was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// The time at which the entity was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// KubernetesName holds the value of the "kubernetes_name" field.
	KubernetesName string `json:"kubernetes_name,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description *string `json:"description,omitempty"`
	// Active holds the value of the "active" field.
	Active bool `json:"active,omitempty"`
	// ProjectID holds the value of the "project_id" field.
	ProjectID uuid.UUID `json:"project_id,omitempty"`
	// Kubernetes secret for this environment
	KubernetesSecret string `json:"kubernetes_secret,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EnvironmentQuery when eager-loading is set.
	Edges        EnvironmentEdges `json:"edges"`
	selectValues sql.SelectValues
}

// EnvironmentEdges holds the relations/edges for other nodes in the graph.
type EnvironmentEdges struct {
	// Project holds the value of the project edge.
	Project *Project `json:"project,omitempty"`
	// Services holds the value of the services edge.
	Services []*Service `json:"services,omitempty"`
	// ProjectDefault holds the value of the project_default edge.
	ProjectDefault []*Project `json:"project_default,omitempty"`
	// ServiceGroups holds the value of the service_groups edge.
	ServiceGroups []*ServiceGroup `json:"service_groups,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EnvironmentEdges) ProjectOrErr() (*Project, error) {
	if e.Project != nil {
		return e.Project, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: project.Label}
	}
	return nil, &NotLoadedError{edge: "project"}
}

// ServicesOrErr returns the Services value or an error if the edge
// was not loaded in eager-loading.
func (e EnvironmentEdges) ServicesOrErr() ([]*Service, error) {
	if e.loadedTypes[1] {
		return e.Services, nil
	}
	return nil, &NotLoadedError{edge: "services"}
}

// ProjectDefaultOrErr returns the ProjectDefault value or an error if the edge
// was not loaded in eager-loading.
func (e EnvironmentEdges) ProjectDefaultOrErr() ([]*Project, error) {
	if e.loadedTypes[2] {
		return e.ProjectDefault, nil
	}
	return nil, &NotLoadedError{edge: "project_default"}
}

// ServiceGroupsOrErr returns the ServiceGroups value or an error if the edge
// was not loaded in eager-loading.
func (e EnvironmentEdges) ServiceGroupsOrErr() ([]*ServiceGroup, error) {
	if e.loadedTypes[3] {
		return e.ServiceGroups, nil
	}
	return nil, &NotLoadedError{edge: "service_groups"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Environment) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case environment.FieldActive:
			values[i] = new(sql.NullBool)
		case environment.FieldKubernetesName, environment.FieldName, environment.FieldDescription, environment.FieldKubernetesSecret:
			values[i] = new(sql.NullString)
		case environment.FieldCreatedAt, environment.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case environment.FieldID, environment.FieldProjectID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Environment fields.
func (e *Environment) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case environment.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				e.ID = *value
			}
		case environment.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				e.CreatedAt = value.Time
			}
		case environment.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				e.UpdatedAt = value.Time
			}
		case environment.FieldKubernetesName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kubernetes_name", values[i])
			} else if value.Valid {
				e.KubernetesName = value.String
			}
		case environment.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				e.Name = value.String
			}
		case environment.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				e.Description = new(string)
				*e.Description = value.String
			}
		case environment.FieldActive:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field active", values[i])
			} else if value.Valid {
				e.Active = value.Bool
			}
		case environment.FieldProjectID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field project_id", values[i])
			} else if value != nil {
				e.ProjectID = *value
			}
		case environment.FieldKubernetesSecret:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kubernetes_secret", values[i])
			} else if value.Valid {
				e.KubernetesSecret = value.String
			}
		default:
			e.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Environment.
// This includes values selected through modifiers, order, etc.
func (e *Environment) Value(name string) (ent.Value, error) {
	return e.selectValues.Get(name)
}

// QueryProject queries the "project" edge of the Environment entity.
func (e *Environment) QueryProject() *ProjectQuery {
	return NewEnvironmentClient(e.config).QueryProject(e)
}

// QueryServices queries the "services" edge of the Environment entity.
func (e *Environment) QueryServices() *ServiceQuery {
	return NewEnvironmentClient(e.config).QueryServices(e)
}

// QueryProjectDefault queries the "project_default" edge of the Environment entity.
func (e *Environment) QueryProjectDefault() *ProjectQuery {
	return NewEnvironmentClient(e.config).QueryProjectDefault(e)
}

// QueryServiceGroups queries the "service_groups" edge of the Environment entity.
func (e *Environment) QueryServiceGroups() *ServiceGroupQuery {
	return NewEnvironmentClient(e.config).QueryServiceGroups(e)
}

// Update returns a builder for updating this Environment.
// Note that you need to call Environment.Unwrap() before calling this method if this Environment
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Environment) Update() *EnvironmentUpdateOne {
	return NewEnvironmentClient(e.config).UpdateOne(e)
}

// Unwrap unwraps the Environment entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Environment) Unwrap() *Environment {
	_tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("ent: Environment is not a transactional entity")
	}
	e.config.driver = _tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Environment) String() string {
	var builder strings.Builder
	builder.WriteString("Environment(")
	builder.WriteString(fmt.Sprintf("id=%v, ", e.ID))
	builder.WriteString("created_at=")
	builder.WriteString(e.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(e.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("kubernetes_name=")
	builder.WriteString(e.KubernetesName)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(e.Name)
	builder.WriteString(", ")
	if v := e.Description; v != nil {
		builder.WriteString("description=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("active=")
	builder.WriteString(fmt.Sprintf("%v", e.Active))
	builder.WriteString(", ")
	builder.WriteString("project_id=")
	builder.WriteString(fmt.Sprintf("%v", e.ProjectID))
	builder.WriteString(", ")
	builder.WriteString("kubernetes_secret=")
	builder.WriteString(e.KubernetesSecret)
	builder.WriteByte(')')
	return builder.String()
}

// Environments is a parsable slice of Environment.
type Environments []*Environment
