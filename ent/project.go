// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/project"
	"github.com/unbindapp/unbind-api/ent/team"
)

// Project is the model entity for the Project schema.
type Project struct {
	config `json:"-"`
	// ID of the ent.
	// The primary key of the entity.
	ID uuid.UUID `json:"id"`
	// The time at which the entity was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// The time at which the entity was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// DisplayName holds the value of the "display_name" field.
	DisplayName string `json:"display_name,omitempty"`
	// Description holds the value of the "description" field.
	Description *string `json:"description,omitempty"`
	// Status holds the value of the "status" field.
	Status string `json:"status,omitempty"`
	// TeamID holds the value of the "team_id" field.
	TeamID uuid.UUID `json:"team_id,omitempty"`
	// Kubernetes secret for this project
	KubernetesSecret string `json:"kubernetes_secret,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ProjectQuery when eager-loading is set.
	Edges        ProjectEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ProjectEdges holds the relations/edges for other nodes in the graph.
type ProjectEdges struct {
	// Team holds the value of the team edge.
	Team *Team `json:"team,omitempty"`
	// Environments holds the value of the environments edge.
	Environments []*Environment `json:"environments,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// TeamOrErr returns the Team value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProjectEdges) TeamOrErr() (*Team, error) {
	if e.Team != nil {
		return e.Team, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: team.Label}
	}
	return nil, &NotLoadedError{edge: "team"}
}

// EnvironmentsOrErr returns the Environments value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectEdges) EnvironmentsOrErr() ([]*Environment, error) {
	if e.loadedTypes[1] {
		return e.Environments, nil
	}
	return nil, &NotLoadedError{edge: "environments"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Project) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case project.FieldName, project.FieldDisplayName, project.FieldDescription, project.FieldStatus, project.FieldKubernetesSecret:
			values[i] = new(sql.NullString)
		case project.FieldCreatedAt, project.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case project.FieldID, project.FieldTeamID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Project fields.
func (pr *Project) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case project.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				pr.ID = *value
			}
		case project.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pr.CreatedAt = value.Time
			}
		case project.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				pr.UpdatedAt = value.Time
			}
		case project.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				pr.Name = value.String
			}
		case project.FieldDisplayName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field display_name", values[i])
			} else if value.Valid {
				pr.DisplayName = value.String
			}
		case project.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				pr.Description = new(string)
				*pr.Description = value.String
			}
		case project.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				pr.Status = value.String
			}
		case project.FieldTeamID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field team_id", values[i])
			} else if value != nil {
				pr.TeamID = *value
			}
		case project.FieldKubernetesSecret:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kubernetes_secret", values[i])
			} else if value.Valid {
				pr.KubernetesSecret = value.String
			}
		default:
			pr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Project.
// This includes values selected through modifiers, order, etc.
func (pr *Project) Value(name string) (ent.Value, error) {
	return pr.selectValues.Get(name)
}

// QueryTeam queries the "team" edge of the Project entity.
func (pr *Project) QueryTeam() *TeamQuery {
	return NewProjectClient(pr.config).QueryTeam(pr)
}

// QueryEnvironments queries the "environments" edge of the Project entity.
func (pr *Project) QueryEnvironments() *EnvironmentQuery {
	return NewProjectClient(pr.config).QueryEnvironments(pr)
}

// Update returns a builder for updating this Project.
// Note that you need to call Project.Unwrap() before calling this method if this Project
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *Project) Update() *ProjectUpdateOne {
	return NewProjectClient(pr.config).UpdateOne(pr)
}

// Unwrap unwraps the Project entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pr *Project) Unwrap() *Project {
	_tx, ok := pr.config.driver.(*txDriver)
	if !ok {
		panic("ent: Project is not a transactional entity")
	}
	pr.config.driver = _tx.drv
	return pr
}

// String implements the fmt.Stringer.
func (pr *Project) String() string {
	var builder strings.Builder
	builder.WriteString("Project(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pr.ID))
	builder.WriteString("created_at=")
	builder.WriteString(pr.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(pr.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(pr.Name)
	builder.WriteString(", ")
	builder.WriteString("display_name=")
	builder.WriteString(pr.DisplayName)
	builder.WriteString(", ")
	if v := pr.Description; v != nil {
		builder.WriteString("description=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(pr.Status)
	builder.WriteString(", ")
	builder.WriteString("team_id=")
	builder.WriteString(fmt.Sprintf("%v", pr.TeamID))
	builder.WriteString(", ")
	builder.WriteString("kubernetes_secret=")
	builder.WriteString(pr.KubernetesSecret)
	builder.WriteByte(')')
	return builder.String()
}

// Projects is a parsable slice of Project.
type Projects []*Project
