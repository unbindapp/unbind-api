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
	"github.com/unbindapp/unbind-api/ent/project"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/team"
	"github.com/unbindapp/unbind-api/ent/webhook"
)

// Webhook is the model entity for the Webhook schema.
type Webhook struct {
	config `json:"-"`
	// ID of the ent.
	// The primary key of the entity.
	ID uuid.UUID `json:"id"`
	// The time at which the entity was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// The time at which the entity was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// URL holds the value of the "url" field.
	URL string `json:"url,omitempty"`
	// Type holds the value of the "type" field.
	Type schema.WebhookType `json:"type,omitempty"`
	// Events holds the value of the "events" field.
	Events []schema.WebhookEvent `json:"events,omitempty"`
	// TeamID holds the value of the "team_id" field.
	TeamID uuid.UUID `json:"team_id,omitempty"`
	// ProjectID holds the value of the "project_id" field.
	ProjectID *uuid.UUID `json:"project_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the WebhookQuery when eager-loading is set.
	Edges        WebhookEdges `json:"edges"`
	selectValues sql.SelectValues
}

// WebhookEdges holds the relations/edges for other nodes in the graph.
type WebhookEdges struct {
	// Team holds the value of the team edge.
	Team *Team `json:"team,omitempty"`
	// Project holds the value of the project edge.
	Project *Project `json:"project,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// TeamOrErr returns the Team value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e WebhookEdges) TeamOrErr() (*Team, error) {
	if e.Team != nil {
		return e.Team, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: team.Label}
	}
	return nil, &NotLoadedError{edge: "team"}
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e WebhookEdges) ProjectOrErr() (*Project, error) {
	if e.Project != nil {
		return e.Project, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: project.Label}
	}
	return nil, &NotLoadedError{edge: "project"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Webhook) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case webhook.FieldProjectID:
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		case webhook.FieldEvents:
			values[i] = new([]byte)
		case webhook.FieldURL, webhook.FieldType:
			values[i] = new(sql.NullString)
		case webhook.FieldCreatedAt, webhook.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case webhook.FieldID, webhook.FieldTeamID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Webhook fields.
func (w *Webhook) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case webhook.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				w.ID = *value
			}
		case webhook.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				w.CreatedAt = value.Time
			}
		case webhook.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				w.UpdatedAt = value.Time
			}
		case webhook.FieldURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field url", values[i])
			} else if value.Valid {
				w.URL = value.String
			}
		case webhook.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				w.Type = schema.WebhookType(value.String)
			}
		case webhook.FieldEvents:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field events", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &w.Events); err != nil {
					return fmt.Errorf("unmarshal field events: %w", err)
				}
			}
		case webhook.FieldTeamID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field team_id", values[i])
			} else if value != nil {
				w.TeamID = *value
			}
		case webhook.FieldProjectID:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field project_id", values[i])
			} else if value.Valid {
				w.ProjectID = new(uuid.UUID)
				*w.ProjectID = *value.S.(*uuid.UUID)
			}
		default:
			w.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Webhook.
// This includes values selected through modifiers, order, etc.
func (w *Webhook) Value(name string) (ent.Value, error) {
	return w.selectValues.Get(name)
}

// QueryTeam queries the "team" edge of the Webhook entity.
func (w *Webhook) QueryTeam() *TeamQuery {
	return NewWebhookClient(w.config).QueryTeam(w)
}

// QueryProject queries the "project" edge of the Webhook entity.
func (w *Webhook) QueryProject() *ProjectQuery {
	return NewWebhookClient(w.config).QueryProject(w)
}

// Update returns a builder for updating this Webhook.
// Note that you need to call Webhook.Unwrap() before calling this method if this Webhook
// was returned from a transaction, and the transaction was committed or rolled back.
func (w *Webhook) Update() *WebhookUpdateOne {
	return NewWebhookClient(w.config).UpdateOne(w)
}

// Unwrap unwraps the Webhook entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (w *Webhook) Unwrap() *Webhook {
	_tx, ok := w.config.driver.(*txDriver)
	if !ok {
		panic("ent: Webhook is not a transactional entity")
	}
	w.config.driver = _tx.drv
	return w
}

// String implements the fmt.Stringer.
func (w *Webhook) String() string {
	var builder strings.Builder
	builder.WriteString("Webhook(")
	builder.WriteString(fmt.Sprintf("id=%v, ", w.ID))
	builder.WriteString("created_at=")
	builder.WriteString(w.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(w.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("url=")
	builder.WriteString(w.URL)
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(fmt.Sprintf("%v", w.Type))
	builder.WriteString(", ")
	builder.WriteString("events=")
	builder.WriteString(fmt.Sprintf("%v", w.Events))
	builder.WriteString(", ")
	builder.WriteString("team_id=")
	builder.WriteString(fmt.Sprintf("%v", w.TeamID))
	builder.WriteString(", ")
	if v := w.ProjectID; v != nil {
		builder.WriteString("project_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Webhooks is a parsable slice of Webhook.
type Webhooks []*Webhook
