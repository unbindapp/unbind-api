// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/s3"
	"github.com/unbindapp/unbind-api/ent/team"
)

// S3 is the model entity for the S3 schema.
type S3 struct {
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
	// Endpoint holds the value of the "endpoint" field.
	Endpoint string `json:"endpoint,omitempty"`
	// Region holds the value of the "region" field.
	Region string `json:"region,omitempty"`
	// ForcePathStyle holds the value of the "force_path_style" field.
	ForcePathStyle bool `json:"force_path_style,omitempty"`
	// KubernetesSecret holds the value of the "kubernetes_secret" field.
	KubernetesSecret string `json:"kubernetes_secret,omitempty"`
	// TeamID holds the value of the "team_id" field.
	TeamID uuid.UUID `json:"team_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the S3Query when eager-loading is set.
	Edges        S3Edges `json:"edges"`
	selectValues sql.SelectValues
}

// S3Edges holds the relations/edges for other nodes in the graph.
type S3Edges struct {
	// Team holds the value of the team edge.
	Team *Team `json:"team,omitempty"`
	// ServiceBackupSource holds the value of the service_backup_source edge.
	ServiceBackupSource []*ServiceConfig `json:"service_backup_source,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// TeamOrErr returns the Team value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e S3Edges) TeamOrErr() (*Team, error) {
	if e.Team != nil {
		return e.Team, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: team.Label}
	}
	return nil, &NotLoadedError{edge: "team"}
}

// ServiceBackupSourceOrErr returns the ServiceBackupSource value or an error if the edge
// was not loaded in eager-loading.
func (e S3Edges) ServiceBackupSourceOrErr() ([]*ServiceConfig, error) {
	if e.loadedTypes[1] {
		return e.ServiceBackupSource, nil
	}
	return nil, &NotLoadedError{edge: "service_backup_source"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*S3) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case s3.FieldForcePathStyle:
			values[i] = new(sql.NullBool)
		case s3.FieldName, s3.FieldEndpoint, s3.FieldRegion, s3.FieldKubernetesSecret:
			values[i] = new(sql.NullString)
		case s3.FieldCreatedAt, s3.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case s3.FieldID, s3.FieldTeamID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the S3 fields.
func (s *S3) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case s3.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				s.ID = *value
			}
		case s3.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				s.CreatedAt = value.Time
			}
		case s3.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				s.UpdatedAt = value.Time
			}
		case s3.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				s.Name = value.String
			}
		case s3.FieldEndpoint:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field endpoint", values[i])
			} else if value.Valid {
				s.Endpoint = value.String
			}
		case s3.FieldRegion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field region", values[i])
			} else if value.Valid {
				s.Region = value.String
			}
		case s3.FieldForcePathStyle:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field force_path_style", values[i])
			} else if value.Valid {
				s.ForcePathStyle = value.Bool
			}
		case s3.FieldKubernetesSecret:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kubernetes_secret", values[i])
			} else if value.Valid {
				s.KubernetesSecret = value.String
			}
		case s3.FieldTeamID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field team_id", values[i])
			} else if value != nil {
				s.TeamID = *value
			}
		default:
			s.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the S3.
// This includes values selected through modifiers, order, etc.
func (s *S3) Value(name string) (ent.Value, error) {
	return s.selectValues.Get(name)
}

// QueryTeam queries the "team" edge of the S3 entity.
func (s *S3) QueryTeam() *TeamQuery {
	return NewS3Client(s.config).QueryTeam(s)
}

// QueryServiceBackupSource queries the "service_backup_source" edge of the S3 entity.
func (s *S3) QueryServiceBackupSource() *ServiceConfigQuery {
	return NewS3Client(s.config).QueryServiceBackupSource(s)
}

// Update returns a builder for updating this S3.
// Note that you need to call S3.Unwrap() before calling this method if this S3
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *S3) Update() *S3UpdateOne {
	return NewS3Client(s.config).UpdateOne(s)
}

// Unwrap unwraps the S3 entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *S3) Unwrap() *S3 {
	_tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: S3 is not a transactional entity")
	}
	s.config.driver = _tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *S3) String() string {
	var builder strings.Builder
	builder.WriteString("S3(")
	builder.WriteString(fmt.Sprintf("id=%v, ", s.ID))
	builder.WriteString("created_at=")
	builder.WriteString(s.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(s.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(s.Name)
	builder.WriteString(", ")
	builder.WriteString("endpoint=")
	builder.WriteString(s.Endpoint)
	builder.WriteString(", ")
	builder.WriteString("region=")
	builder.WriteString(s.Region)
	builder.WriteString(", ")
	builder.WriteString("force_path_style=")
	builder.WriteString(fmt.Sprintf("%v", s.ForcePathStyle))
	builder.WriteString(", ")
	builder.WriteString("kubernetes_secret=")
	builder.WriteString(s.KubernetesSecret)
	builder.WriteString(", ")
	builder.WriteString("team_id=")
	builder.WriteString(fmt.Sprintf("%v", s.TeamID))
	builder.WriteByte(')')
	return builder.String()
}

// S3s is a parsable slice of S3.
type S3s []*S3
