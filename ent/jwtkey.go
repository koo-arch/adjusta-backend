// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/koo-arch/adjusta-backend/ent/jwtkey"
)

// JWTKey is the model entity for the JWTKey schema.
type JWTKey struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	// Key holds the value of the "key" field.
	Key string `json:"-"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// ExpiresAt holds the value of the "expires_at" field.
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*JWTKey) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case jwtkey.FieldID:
			values[i] = new(sql.NullInt64)
		case jwtkey.FieldKey, jwtkey.FieldType:
			values[i] = new(sql.NullString)
		case jwtkey.FieldCreatedAt, jwtkey.FieldUpdatedAt, jwtkey.FieldDeletedAt, jwtkey.FieldExpiresAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the JWTKey fields.
func (jk *JWTKey) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case jwtkey.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			jk.ID = int(value.Int64)
		case jwtkey.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				jk.CreatedAt = value.Time
			}
		case jwtkey.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				jk.UpdatedAt = value.Time
			}
		case jwtkey.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				jk.DeletedAt = new(time.Time)
				*jk.DeletedAt = value.Time
			}
		case jwtkey.FieldKey:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field key", values[i])
			} else if value.Valid {
				jk.Key = value.String
			}
		case jwtkey.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				jk.Type = value.String
			}
		case jwtkey.FieldExpiresAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field expires_at", values[i])
			} else if value.Valid {
				jk.ExpiresAt = value.Time
			}
		default:
			jk.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the JWTKey.
// This includes values selected through modifiers, order, etc.
func (jk *JWTKey) Value(name string) (ent.Value, error) {
	return jk.selectValues.Get(name)
}

// Update returns a builder for updating this JWTKey.
// Note that you need to call JWTKey.Unwrap() before calling this method if this JWTKey
// was returned from a transaction, and the transaction was committed or rolled back.
func (jk *JWTKey) Update() *JWTKeyUpdateOne {
	return NewJWTKeyClient(jk.config).UpdateOne(jk)
}

// Unwrap unwraps the JWTKey entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (jk *JWTKey) Unwrap() *JWTKey {
	_tx, ok := jk.config.driver.(*txDriver)
	if !ok {
		panic("ent: JWTKey is not a transactional entity")
	}
	jk.config.driver = _tx.drv
	return jk
}

// String implements the fmt.Stringer.
func (jk *JWTKey) String() string {
	var builder strings.Builder
	builder.WriteString("JWTKey(")
	builder.WriteString(fmt.Sprintf("id=%v, ", jk.ID))
	builder.WriteString("created_at=")
	builder.WriteString(jk.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(jk.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := jk.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("key=<sensitive>")
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(jk.Type)
	builder.WriteString(", ")
	builder.WriteString("expires_at=")
	builder.WriteString(jk.ExpiresAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// JWTKeys is a parsable slice of JWTKey.
type JWTKeys []*JWTKey
