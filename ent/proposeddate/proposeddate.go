// Code generated by ent, DO NOT EDIT.

package proposeddate

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the proposeddate type in the database.
	Label = "proposed_date"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldGoogleEventID holds the string denoting the google_event_id field in the database.
	FieldGoogleEventID = "google_event_id"
	// FieldStartTime holds the string denoting the start_time field in the database.
	FieldStartTime = "start_time"
	// FieldEndTime holds the string denoting the end_time field in the database.
	FieldEndTime = "end_time"
	// FieldPriority holds the string denoting the priority field in the database.
	FieldPriority = "priority"
	// EdgeEvent holds the string denoting the event edge name in mutations.
	EdgeEvent = "event"
	// Table holds the table name of the proposeddate in the database.
	Table = "proposed_dates"
	// EventTable is the table that holds the event relation/edge.
	EventTable = "proposed_dates"
	// EventInverseTable is the table name for the Event entity.
	// It exists in this package in order to avoid circular dependency with the "event" package.
	EventInverseTable = "events"
	// EventColumn is the table column denoting the event relation/edge.
	EventColumn = "event_proposed_dates"
)

// Columns holds all SQL columns for proposeddate fields.
var Columns = []string{
	FieldID,
	FieldGoogleEventID,
	FieldStartTime,
	FieldEndTime,
	FieldPriority,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "proposed_dates"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"event_proposed_dates",
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
	// DefaultPriority holds the default value on creation for the "priority" field.
	DefaultPriority int
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the ProposedDate queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByGoogleEventID orders the results by the google_event_id field.
func ByGoogleEventID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGoogleEventID, opts...).ToFunc()
}

// ByStartTime orders the results by the start_time field.
func ByStartTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStartTime, opts...).ToFunc()
}

// ByEndTime orders the results by the end_time field.
func ByEndTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEndTime, opts...).ToFunc()
}

// ByPriority orders the results by the priority field.
func ByPriority(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPriority, opts...).ToFunc()
}

// ByEventField orders the results by event field.
func ByEventField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEventStep(), sql.OrderByField(field, opts...))
	}
}
func newEventStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EventInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, EventTable, EventColumn),
	)
}
