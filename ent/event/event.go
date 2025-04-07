// Code generated by ent, DO NOT EDIT.

package event

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the event type in the database.
	Label = "event"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldSummary holds the string denoting the summary field in the database.
	FieldSummary = "summary"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldLocation holds the string denoting the location field in the database.
	FieldLocation = "location"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldConfirmedDateID holds the string denoting the confirmed_date_id field in the database.
	FieldConfirmedDateID = "confirmed_date_id"
	// FieldGoogleEventID holds the string denoting the google_event_id field in the database.
	FieldGoogleEventID = "google_event_id"
	// FieldSlug holds the string denoting the slug field in the database.
	FieldSlug = "slug"
	// EdgeCalendar holds the string denoting the calendar edge name in mutations.
	EdgeCalendar = "calendar"
	// EdgeProposedDates holds the string denoting the proposed_dates edge name in mutations.
	EdgeProposedDates = "proposed_dates"
	// Table holds the table name of the event in the database.
	Table = "events"
	// CalendarTable is the table that holds the calendar relation/edge.
	CalendarTable = "events"
	// CalendarInverseTable is the table name for the Calendar entity.
	// It exists in this package in order to avoid circular dependency with the "calendar" package.
	CalendarInverseTable = "calendars"
	// CalendarColumn is the table column denoting the calendar relation/edge.
	CalendarColumn = "calendar_events"
	// ProposedDatesTable is the table that holds the proposed_dates relation/edge.
	ProposedDatesTable = "proposed_dates"
	// ProposedDatesInverseTable is the table name for the ProposedDate entity.
	// It exists in this package in order to avoid circular dependency with the "proposeddate" package.
	ProposedDatesInverseTable = "proposed_dates"
	// ProposedDatesColumn is the table column denoting the proposed_dates relation/edge.
	ProposedDatesColumn = "event_proposed_dates"
)

// Columns holds all SQL columns for event fields.
var Columns = []string{
	FieldID,
	FieldSummary,
	FieldDescription,
	FieldLocation,
	FieldStatus,
	FieldConfirmedDateID,
	FieldGoogleEventID,
	FieldSlug,
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

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/koo-arch/adjusta-backend/ent/runtime"
var (
	Hooks [1]ent.Hook
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// Status defines the type for the "status" enum field.
type Status string

// StatusPending is the default value of the Status enum.
const DefaultStatus = StatusPending

// Status values.
const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCancelled Status = "cancelled"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusPending, StatusConfirmed, StatusCancelled:
		return nil
	default:
		return fmt.Errorf("event: invalid enum value for status field: %q", s)
	}
}

// OrderOption defines the ordering options for the Event queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
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

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByConfirmedDateID orders the results by the confirmed_date_id field.
func ByConfirmedDateID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldConfirmedDateID, opts...).ToFunc()
}

// ByGoogleEventID orders the results by the google_event_id field.
func ByGoogleEventID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGoogleEventID, opts...).ToFunc()
}

// BySlug orders the results by the slug field.
func BySlug(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSlug, opts...).ToFunc()
}

// ByCalendarField orders the results by calendar field.
func ByCalendarField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCalendarStep(), sql.OrderByField(field, opts...))
	}
}

// ByProposedDatesCount orders the results by proposed_dates count.
func ByProposedDatesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newProposedDatesStep(), opts...)
	}
}

// ByProposedDates orders the results by proposed_dates terms.
func ByProposedDates(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProposedDatesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newCalendarStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CalendarInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, CalendarTable, CalendarColumn),
	)
}
func newProposedDatesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProposedDatesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ProposedDatesTable, ProposedDatesColumn),
	)
}
