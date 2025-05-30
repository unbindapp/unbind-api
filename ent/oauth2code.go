// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/oauth2code"
	"github.com/unbindapp/unbind-api/ent/user"
)

// Oauth2Code is the model entity for the Oauth2Code schema.
type Oauth2Code struct {
	config `json:"-"`
	// ID of the ent.
	// The primary key of the entity.
	ID uuid.UUID `json:"id"`
	// The time at which the entity was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// The time at which the entity was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// AuthCode holds the value of the "auth_code" field.
	AuthCode string `json:"-"`
	// ClientID holds the value of the "client_id" field.
	ClientID string `json:"client_id,omitempty"`
	// Scope holds the value of the "scope" field.
	Scope string `json:"scope,omitempty"`
	// ExpiresAt holds the value of the "expires_at" field.
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	// Revoked holds the value of the "revoked" field.
	Revoked bool `json:"revoked,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the Oauth2CodeQuery when eager-loading is set.
	Edges             Oauth2CodeEdges `json:"edges"`
	user_oauth2_codes *uuid.UUID
	selectValues      sql.SelectValues
}

// Oauth2CodeEdges holds the relations/edges for other nodes in the graph.
type Oauth2CodeEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e Oauth2CodeEdges) UserOrErr() (*User, error) {
	if e.User != nil {
		return e.User, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "user"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Oauth2Code) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case oauth2code.FieldRevoked:
			values[i] = new(sql.NullBool)
		case oauth2code.FieldAuthCode, oauth2code.FieldClientID, oauth2code.FieldScope:
			values[i] = new(sql.NullString)
		case oauth2code.FieldCreatedAt, oauth2code.FieldUpdatedAt, oauth2code.FieldExpiresAt:
			values[i] = new(sql.NullTime)
		case oauth2code.FieldID:
			values[i] = new(uuid.UUID)
		case oauth2code.ForeignKeys[0]: // user_oauth2_codes
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Oauth2Code fields.
func (o *Oauth2Code) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case oauth2code.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				o.ID = *value
			}
		case oauth2code.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				o.CreatedAt = value.Time
			}
		case oauth2code.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				o.UpdatedAt = value.Time
			}
		case oauth2code.FieldAuthCode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field auth_code", values[i])
			} else if value.Valid {
				o.AuthCode = value.String
			}
		case oauth2code.FieldClientID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field client_id", values[i])
			} else if value.Valid {
				o.ClientID = value.String
			}
		case oauth2code.FieldScope:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field scope", values[i])
			} else if value.Valid {
				o.Scope = value.String
			}
		case oauth2code.FieldExpiresAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field expires_at", values[i])
			} else if value.Valid {
				o.ExpiresAt = value.Time
			}
		case oauth2code.FieldRevoked:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field revoked", values[i])
			} else if value.Valid {
				o.Revoked = value.Bool
			}
		case oauth2code.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field user_oauth2_codes", values[i])
			} else if value.Valid {
				o.user_oauth2_codes = new(uuid.UUID)
				*o.user_oauth2_codes = *value.S.(*uuid.UUID)
			}
		default:
			o.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Oauth2Code.
// This includes values selected through modifiers, order, etc.
func (o *Oauth2Code) Value(name string) (ent.Value, error) {
	return o.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the Oauth2Code entity.
func (o *Oauth2Code) QueryUser() *UserQuery {
	return NewOauth2CodeClient(o.config).QueryUser(o)
}

// Update returns a builder for updating this Oauth2Code.
// Note that you need to call Oauth2Code.Unwrap() before calling this method if this Oauth2Code
// was returned from a transaction, and the transaction was committed or rolled back.
func (o *Oauth2Code) Update() *Oauth2CodeUpdateOne {
	return NewOauth2CodeClient(o.config).UpdateOne(o)
}

// Unwrap unwraps the Oauth2Code entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (o *Oauth2Code) Unwrap() *Oauth2Code {
	_tx, ok := o.config.driver.(*txDriver)
	if !ok {
		panic("ent: Oauth2Code is not a transactional entity")
	}
	o.config.driver = _tx.drv
	return o
}

// String implements the fmt.Stringer.
func (o *Oauth2Code) String() string {
	var builder strings.Builder
	builder.WriteString("Oauth2Code(")
	builder.WriteString(fmt.Sprintf("id=%v, ", o.ID))
	builder.WriteString("created_at=")
	builder.WriteString(o.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(o.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("auth_code=<sensitive>")
	builder.WriteString(", ")
	builder.WriteString("client_id=")
	builder.WriteString(o.ClientID)
	builder.WriteString(", ")
	builder.WriteString("scope=")
	builder.WriteString(o.Scope)
	builder.WriteString(", ")
	builder.WriteString("expires_at=")
	builder.WriteString(o.ExpiresAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("revoked=")
	builder.WriteString(fmt.Sprintf("%v", o.Revoked))
	builder.WriteByte(')')
	return builder.String()
}

// Oauth2Codes is a parsable slice of Oauth2Code.
type Oauth2Codes []*Oauth2Code
