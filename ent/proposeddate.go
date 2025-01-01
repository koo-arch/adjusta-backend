// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
)

// ProposedDate is the model entity for the ProposedDate schema.
type ProposedDate struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// StartTime holds the value of the "start_time" field.
	StartTime time.Time `json:"start_time,omitempty"`
	// EndTime holds the value of the "end_time" field.
	EndTime time.Time `json:"end_time,omitempty"`
	// Priority holds the value of the "priority" field.
	Priority int `json:"priority,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ProposedDateQuery when eager-loading is set.
	Edges                ProposedDateEdges `json:"edges"`
	event_proposed_dates *uuid.UUID
	selectValues         sql.SelectValues
}

// ProposedDateEdges holds the relations/edges for other nodes in the graph.
type ProposedDateEdges struct {
	// Event holds the value of the event edge.
	Event *Event `json:"event,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// EventOrErr returns the Event value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProposedDateEdges) EventOrErr() (*Event, error) {
	if e.Event != nil {
		return e.Event, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: event.Label}
	}
	return nil, &NotLoadedError{edge: "event"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ProposedDate) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case proposeddate.FieldPriority:
			values[i] = new(sql.NullInt64)
		case proposeddate.FieldStartTime, proposeddate.FieldEndTime:
			values[i] = new(sql.NullTime)
		case proposeddate.FieldID:
			values[i] = new(uuid.UUID)
		case proposeddate.ForeignKeys[0]: // event_proposed_dates
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ProposedDate fields.
func (pd *ProposedDate) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case proposeddate.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				pd.ID = *value
			}
		case proposeddate.FieldStartTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field start_time", values[i])
			} else if value.Valid {
				pd.StartTime = value.Time
			}
		case proposeddate.FieldEndTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field end_time", values[i])
			} else if value.Valid {
				pd.EndTime = value.Time
			}
		case proposeddate.FieldPriority:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field priority", values[i])
			} else if value.Valid {
				pd.Priority = int(value.Int64)
			}
		case proposeddate.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field event_proposed_dates", values[i])
			} else if value.Valid {
				pd.event_proposed_dates = new(uuid.UUID)
				*pd.event_proposed_dates = *value.S.(*uuid.UUID)
			}
		default:
			pd.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ProposedDate.
// This includes values selected through modifiers, order, etc.
func (pd *ProposedDate) Value(name string) (ent.Value, error) {
	return pd.selectValues.Get(name)
}

// QueryEvent queries the "event" edge of the ProposedDate entity.
func (pd *ProposedDate) QueryEvent() *EventQuery {
	return NewProposedDateClient(pd.config).QueryEvent(pd)
}

// Update returns a builder for updating this ProposedDate.
// Note that you need to call ProposedDate.Unwrap() before calling this method if this ProposedDate
// was returned from a transaction, and the transaction was committed or rolled back.
func (pd *ProposedDate) Update() *ProposedDateUpdateOne {
	return NewProposedDateClient(pd.config).UpdateOne(pd)
}

// Unwrap unwraps the ProposedDate entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pd *ProposedDate) Unwrap() *ProposedDate {
	_tx, ok := pd.config.driver.(*txDriver)
	if !ok {
		panic("ent: ProposedDate is not a transactional entity")
	}
	pd.config.driver = _tx.drv
	return pd
}

// String implements the fmt.Stringer.
func (pd *ProposedDate) String() string {
	var builder strings.Builder
	builder.WriteString("ProposedDate(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pd.ID))
	builder.WriteString("start_time=")
	builder.WriteString(pd.StartTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("end_time=")
	builder.WriteString(pd.EndTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("priority=")
	builder.WriteString(fmt.Sprintf("%v", pd.Priority))
	builder.WriteByte(')')
	return builder.String()
}

// ProposedDates is a parsable slice of ProposedDate.
type ProposedDates []*ProposedDate
