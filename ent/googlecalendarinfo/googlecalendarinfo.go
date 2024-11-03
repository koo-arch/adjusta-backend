// Code generated by ent, DO NOT EDIT.

package googlecalendarinfo

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the googlecalendarinfo type in the database.
	Label = "google_calendar_info"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldGoogleCalendarID holds the string denoting the google_calendar_id field in the database.
	FieldGoogleCalendarID = "google_calendar_id"
	// FieldSummary holds the string denoting the summary field in the database.
	FieldSummary = "summary"
	// FieldIsPrimary holds the string denoting the is_primary field in the database.
	FieldIsPrimary = "is_primary"
	// EdgeCalendars holds the string denoting the calendars edge name in mutations.
	EdgeCalendars = "calendars"
	// Table holds the table name of the googlecalendarinfo in the database.
	Table = "google_calendar_infos"
	// CalendarsTable is the table that holds the calendars relation/edge. The primary key declared below.
	CalendarsTable = "calendar_google_calendar_infos"
	// CalendarsInverseTable is the table name for the Calendar entity.
	// It exists in this package in order to avoid circular dependency with the "calendar" package.
	CalendarsInverseTable = "calendars"
)

// Columns holds all SQL columns for googlecalendarinfo fields.
var Columns = []string{
	FieldID,
	FieldGoogleCalendarID,
	FieldSummary,
	FieldIsPrimary,
}

var (
	// CalendarsPrimaryKey and CalendarsColumn2 are the table columns denoting the
	// primary key for the calendars relation (M2M).
	CalendarsPrimaryKey = []string{"calendar_id", "google_calendar_info_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// GoogleCalendarIDValidator is a validator for the "google_calendar_id" field. It is called by the builders before save.
	GoogleCalendarIDValidator func(string) error
	// DefaultIsPrimary holds the default value on creation for the "is_primary" field.
	DefaultIsPrimary bool
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the GoogleCalendarInfo queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByGoogleCalendarID orders the results by the google_calendar_id field.
func ByGoogleCalendarID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGoogleCalendarID, opts...).ToFunc()
}

// BySummary orders the results by the summary field.
func BySummary(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSummary, opts...).ToFunc()
}

// ByIsPrimary orders the results by the is_primary field.
func ByIsPrimary(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsPrimary, opts...).ToFunc()
}

// ByCalendarsCount orders the results by calendars count.
func ByCalendarsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCalendarsStep(), opts...)
	}
}

// ByCalendars orders the results by calendars terms.
func ByCalendars(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCalendarsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newCalendarsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CalendarsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, CalendarsTable, CalendarsPrimaryKey...),
	)
}