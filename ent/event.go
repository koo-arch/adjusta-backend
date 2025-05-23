// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/event"
)

// Event is the model entity for the Event schema.
type Event struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	// Summary holds the value of the "summary" field.
	Summary string `json:"summary,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Location holds the value of the "location" field.
	Location string `json:"location,omitempty"`
	// Status holds the value of the "status" field.
	Status event.Status `json:"status,omitempty"`
	// ConfirmedDateID holds the value of the "confirmed_date_id" field.
	ConfirmedDateID uuid.UUID `json:"confirmed_date_id,omitempty"`
	// GoogleEventID holds the value of the "google_event_id" field.
	GoogleEventID string `json:"google_event_id,omitempty"`
	// Slug holds the value of the "slug" field.
	Slug string `json:"slug,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EventQuery when eager-loading is set.
	Edges           EventEdges `json:"edges"`
	calendar_events *uuid.UUID
	selectValues    sql.SelectValues
}

// EventEdges holds the relations/edges for other nodes in the graph.
type EventEdges struct {
	// Calendar holds the value of the calendar edge.
	Calendar *Calendar `json:"calendar,omitempty"`
	// ProposedDates holds the value of the proposed_dates edge.
	ProposedDates []*ProposedDate `json:"proposed_dates,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// CalendarOrErr returns the Calendar value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EventEdges) CalendarOrErr() (*Calendar, error) {
	if e.Calendar != nil {
		return e.Calendar, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: calendar.Label}
	}
	return nil, &NotLoadedError{edge: "calendar"}
}

// ProposedDatesOrErr returns the ProposedDates value or an error if the edge
// was not loaded in eager-loading.
func (e EventEdges) ProposedDatesOrErr() ([]*ProposedDate, error) {
	if e.loadedTypes[1] {
		return e.ProposedDates, nil
	}
	return nil, &NotLoadedError{edge: "proposed_dates"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Event) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case event.FieldSummary, event.FieldDescription, event.FieldLocation, event.FieldStatus, event.FieldGoogleEventID, event.FieldSlug:
			values[i] = new(sql.NullString)
		case event.FieldCreatedAt, event.FieldUpdatedAt, event.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		case event.FieldID, event.FieldConfirmedDateID:
			values[i] = new(uuid.UUID)
		case event.ForeignKeys[0]: // calendar_events
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Event fields.
func (e *Event) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case event.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				e.ID = *value
			}
		case event.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				e.CreatedAt = value.Time
			}
		case event.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				e.UpdatedAt = value.Time
			}
		case event.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				e.DeletedAt = new(time.Time)
				*e.DeletedAt = value.Time
			}
		case event.FieldSummary:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field summary", values[i])
			} else if value.Valid {
				e.Summary = value.String
			}
		case event.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				e.Description = value.String
			}
		case event.FieldLocation:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field location", values[i])
			} else if value.Valid {
				e.Location = value.String
			}
		case event.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				e.Status = event.Status(value.String)
			}
		case event.FieldConfirmedDateID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field confirmed_date_id", values[i])
			} else if value != nil {
				e.ConfirmedDateID = *value
			}
		case event.FieldGoogleEventID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field google_event_id", values[i])
			} else if value.Valid {
				e.GoogleEventID = value.String
			}
		case event.FieldSlug:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field slug", values[i])
			} else if value.Valid {
				e.Slug = value.String
			}
		case event.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field calendar_events", values[i])
			} else if value.Valid {
				e.calendar_events = new(uuid.UUID)
				*e.calendar_events = *value.S.(*uuid.UUID)
			}
		default:
			e.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Event.
// This includes values selected through modifiers, order, etc.
func (e *Event) Value(name string) (ent.Value, error) {
	return e.selectValues.Get(name)
}

// QueryCalendar queries the "calendar" edge of the Event entity.
func (e *Event) QueryCalendar() *CalendarQuery {
	return NewEventClient(e.config).QueryCalendar(e)
}

// QueryProposedDates queries the "proposed_dates" edge of the Event entity.
func (e *Event) QueryProposedDates() *ProposedDateQuery {
	return NewEventClient(e.config).QueryProposedDates(e)
}

// Update returns a builder for updating this Event.
// Note that you need to call Event.Unwrap() before calling this method if this Event
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Event) Update() *EventUpdateOne {
	return NewEventClient(e.config).UpdateOne(e)
}

// Unwrap unwraps the Event entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Event) Unwrap() *Event {
	_tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("ent: Event is not a transactional entity")
	}
	e.config.driver = _tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Event) String() string {
	var builder strings.Builder
	builder.WriteString("Event(")
	builder.WriteString(fmt.Sprintf("id=%v, ", e.ID))
	builder.WriteString("created_at=")
	builder.WriteString(e.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(e.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := e.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("summary=")
	builder.WriteString(e.Summary)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(e.Description)
	builder.WriteString(", ")
	builder.WriteString("location=")
	builder.WriteString(e.Location)
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", e.Status))
	builder.WriteString(", ")
	builder.WriteString("confirmed_date_id=")
	builder.WriteString(fmt.Sprintf("%v", e.ConfirmedDateID))
	builder.WriteString(", ")
	builder.WriteString("google_event_id=")
	builder.WriteString(e.GoogleEventID)
	builder.WriteString(", ")
	builder.WriteString("slug=")
	builder.WriteString(e.Slug)
	builder.WriteByte(')')
	return builder.String()
}

// Events is a parsable slice of Event.
type Events []*Event
