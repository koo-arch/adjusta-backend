// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/user"
)

// Calendar is the model entity for the Calendar schema.
type Calendar struct {
	config
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the CalendarQuery when eager-loading is set.
	Edges          CalendarEdges `json:"edges"`
	user_calendars *uuid.UUID
	selectValues   sql.SelectValues
}

// CalendarEdges holds the relations/edges for other nodes in the graph.
type CalendarEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// GoogleCalendarInfos holds the value of the google_calendar_infos edge.
	GoogleCalendarInfos []*GoogleCalendarInfo `json:"google_calendar_infos,omitempty"`
	// Events holds the value of the events edge.
	Events []*Event `json:"events,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CalendarEdges) UserOrErr() (*User, error) {
	if e.User != nil {
		return e.User, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "user"}
}

// GoogleCalendarInfosOrErr returns the GoogleCalendarInfos value or an error if the edge
// was not loaded in eager-loading.
func (e CalendarEdges) GoogleCalendarInfosOrErr() ([]*GoogleCalendarInfo, error) {
	if e.loadedTypes[1] {
		return e.GoogleCalendarInfos, nil
	}
	return nil, &NotLoadedError{edge: "google_calendar_infos"}
}

// EventsOrErr returns the Events value or an error if the edge
// was not loaded in eager-loading.
func (e CalendarEdges) EventsOrErr() ([]*Event, error) {
	if e.loadedTypes[2] {
		return e.Events, nil
	}
	return nil, &NotLoadedError{edge: "events"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Calendar) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case calendar.FieldID:
			values[i] = new(uuid.UUID)
		case calendar.ForeignKeys[0]: // user_calendars
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Calendar fields.
func (c *Calendar) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case calendar.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				c.ID = *value
			}
		case calendar.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field user_calendars", values[i])
			} else if value.Valid {
				c.user_calendars = new(uuid.UUID)
				*c.user_calendars = *value.S.(*uuid.UUID)
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Calendar.
// This includes values selected through modifiers, order, etc.
func (c *Calendar) Value(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the Calendar entity.
func (c *Calendar) QueryUser() *UserQuery {
	return NewCalendarClient(c.config).QueryUser(c)
}

// QueryGoogleCalendarInfos queries the "google_calendar_infos" edge of the Calendar entity.
func (c *Calendar) QueryGoogleCalendarInfos() *GoogleCalendarInfoQuery {
	return NewCalendarClient(c.config).QueryGoogleCalendarInfos(c)
}

// QueryEvents queries the "events" edge of the Calendar entity.
func (c *Calendar) QueryEvents() *EventQuery {
	return NewCalendarClient(c.config).QueryEvents(c)
}

// Update returns a builder for updating this Calendar.
// Note that you need to call Calendar.Unwrap() before calling this method if this Calendar
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Calendar) Update() *CalendarUpdateOne {
	return NewCalendarClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Calendar entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Calendar) Unwrap() *Calendar {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Calendar is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Calendar) String() string {
	var builder strings.Builder
	builder.WriteString("Calendar(")
	builder.WriteString(fmt.Sprintf("id=%v", c.ID))
	builder.WriteByte(')')
	return builder.String()
}

// Calendars is a parsable slice of Calendar.
type Calendars []*Calendar
