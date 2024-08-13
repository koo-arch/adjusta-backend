// Code generated by ent, DO NOT EDIT.

package event

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the event type in the database.
	Label = "event"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldEventID holds the string denoting the event_id field in the database.
	FieldEventID = "event_id"
	// FieldSummary holds the string denoting the summary field in the database.
	FieldSummary = "summary"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldLocation holds the string denoting the location field in the database.
	FieldLocation = "location"
	// FieldStartTime holds the string denoting the start_time field in the database.
	FieldStartTime = "start_time"
	// FieldEndTime holds the string denoting the end_time field in the database.
	FieldEndTime = "end_time"
	// EdgeCalendar holds the string denoting the calendar edge name in mutations.
	EdgeCalendar = "calendar"
	// Table holds the table name of the event in the database.
	Table = "events"
	// CalendarTable is the table that holds the calendar relation/edge.
	CalendarTable = "events"
	// CalendarInverseTable is the table name for the Calendar entity.
	// It exists in this package in order to avoid circular dependency with the "calendar" package.
	CalendarInverseTable = "calendars"
	// CalendarColumn is the table column denoting the calendar relation/edge.
	CalendarColumn = "calendar_events"
)

// Columns holds all SQL columns for event fields.
var Columns = []string{
	FieldID,
	FieldEventID,
	FieldSummary,
	FieldDescription,
	FieldLocation,
	FieldStartTime,
	FieldEndTime,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "events"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"calendar_events",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Event queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByEventID orders the results by the event_id field.
func ByEventID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEventID, opts...).ToFunc()
}

// BySummary orders the results by the summary field.
func BySummary(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSummary, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByLocation orders the results by the location field.
func ByLocation(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLocation, opts...).ToFunc()
}

// ByStartTime orders the results by the start_time field.
func ByStartTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStartTime, opts...).ToFunc()
}

// ByEndTime orders the results by the end_time field.
func ByEndTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEndTime, opts...).ToFunc()
}

// ByCalendarField orders the results by calendar field.
func ByCalendarField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCalendarStep(), sql.OrderByField(field, opts...))
	}
}
func newCalendarStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CalendarInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, CalendarTable, CalendarColumn),
	)
}