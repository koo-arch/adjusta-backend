// Code generated by ent, DO NOT EDIT.

package event

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Event {
	return predicate.Event(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Event {
	return predicate.Event(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Event {
	return predicate.Event(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Event {
	return predicate.Event(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Event {
	return predicate.Event(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Event {
	return predicate.Event(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Event {
	return predicate.Event(sql.FieldLTE(FieldID, id))
}

// EventID applies equality check predicate on the "event_id" field. It's identical to EventIDEQ.
func EventID(v string) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldEventID, v))
}

// Summary applies equality check predicate on the "summary" field. It's identical to SummaryEQ.
func Summary(v string) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldSummary, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldDescription, v))
}

// Location applies equality check predicate on the "location" field. It's identical to LocationEQ.
func Location(v string) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldLocation, v))
}

// StartTime applies equality check predicate on the "start_time" field. It's identical to StartTimeEQ.
func StartTime(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldStartTime, v))
}

// EndTime applies equality check predicate on the "end_time" field. It's identical to EndTimeEQ.
func EndTime(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldEndTime, v))
}

// EventIDEQ applies the EQ predicate on the "event_id" field.
func EventIDEQ(v string) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldEventID, v))
}

// EventIDNEQ applies the NEQ predicate on the "event_id" field.
func EventIDNEQ(v string) predicate.Event {
	return predicate.Event(sql.FieldNEQ(FieldEventID, v))
}

// EventIDIn applies the In predicate on the "event_id" field.
func EventIDIn(vs ...string) predicate.Event {
	return predicate.Event(sql.FieldIn(FieldEventID, vs...))
}

// EventIDNotIn applies the NotIn predicate on the "event_id" field.
func EventIDNotIn(vs ...string) predicate.Event {
	return predicate.Event(sql.FieldNotIn(FieldEventID, vs...))
}

// EventIDGT applies the GT predicate on the "event_id" field.
func EventIDGT(v string) predicate.Event {
	return predicate.Event(sql.FieldGT(FieldEventID, v))
}

// EventIDGTE applies the GTE predicate on the "event_id" field.
func EventIDGTE(v string) predicate.Event {
	return predicate.Event(sql.FieldGTE(FieldEventID, v))
}

// EventIDLT applies the LT predicate on the "event_id" field.
func EventIDLT(v string) predicate.Event {
	return predicate.Event(sql.FieldLT(FieldEventID, v))
}

// EventIDLTE applies the LTE predicate on the "event_id" field.
func EventIDLTE(v string) predicate.Event {
	return predicate.Event(sql.FieldLTE(FieldEventID, v))
}

// EventIDContains applies the Contains predicate on the "event_id" field.
func EventIDContains(v string) predicate.Event {
	return predicate.Event(sql.FieldContains(FieldEventID, v))
}

// EventIDHasPrefix applies the HasPrefix predicate on the "event_id" field.
func EventIDHasPrefix(v string) predicate.Event {
	return predicate.Event(sql.FieldHasPrefix(FieldEventID, v))
}

// EventIDHasSuffix applies the HasSuffix predicate on the "event_id" field.
func EventIDHasSuffix(v string) predicate.Event {
	return predicate.Event(sql.FieldHasSuffix(FieldEventID, v))
}

// EventIDEqualFold applies the EqualFold predicate on the "event_id" field.
func EventIDEqualFold(v string) predicate.Event {
	return predicate.Event(sql.FieldEqualFold(FieldEventID, v))
}

// EventIDContainsFold applies the ContainsFold predicate on the "event_id" field.
func EventIDContainsFold(v string) predicate.Event {
	return predicate.Event(sql.FieldContainsFold(FieldEventID, v))
}

// SummaryEQ applies the EQ predicate on the "summary" field.
func SummaryEQ(v string) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldSummary, v))
}

// SummaryNEQ applies the NEQ predicate on the "summary" field.
func SummaryNEQ(v string) predicate.Event {
	return predicate.Event(sql.FieldNEQ(FieldSummary, v))
}

// SummaryIn applies the In predicate on the "summary" field.
func SummaryIn(vs ...string) predicate.Event {
	return predicate.Event(sql.FieldIn(FieldSummary, vs...))
}

// SummaryNotIn applies the NotIn predicate on the "summary" field.
func SummaryNotIn(vs ...string) predicate.Event {
	return predicate.Event(sql.FieldNotIn(FieldSummary, vs...))
}

// SummaryGT applies the GT predicate on the "summary" field.
func SummaryGT(v string) predicate.Event {
	return predicate.Event(sql.FieldGT(FieldSummary, v))
}

// SummaryGTE applies the GTE predicate on the "summary" field.
func SummaryGTE(v string) predicate.Event {
	return predicate.Event(sql.FieldGTE(FieldSummary, v))
}

// SummaryLT applies the LT predicate on the "summary" field.
func SummaryLT(v string) predicate.Event {
	return predicate.Event(sql.FieldLT(FieldSummary, v))
}

// SummaryLTE applies the LTE predicate on the "summary" field.
func SummaryLTE(v string) predicate.Event {
	return predicate.Event(sql.FieldLTE(FieldSummary, v))
}

// SummaryContains applies the Contains predicate on the "summary" field.
func SummaryContains(v string) predicate.Event {
	return predicate.Event(sql.FieldContains(FieldSummary, v))
}

// SummaryHasPrefix applies the HasPrefix predicate on the "summary" field.
func SummaryHasPrefix(v string) predicate.Event {
	return predicate.Event(sql.FieldHasPrefix(FieldSummary, v))
}

// SummaryHasSuffix applies the HasSuffix predicate on the "summary" field.
func SummaryHasSuffix(v string) predicate.Event {
	return predicate.Event(sql.FieldHasSuffix(FieldSummary, v))
}

// SummaryIsNil applies the IsNil predicate on the "summary" field.
func SummaryIsNil() predicate.Event {
	return predicate.Event(sql.FieldIsNull(FieldSummary))
}

// SummaryNotNil applies the NotNil predicate on the "summary" field.
func SummaryNotNil() predicate.Event {
	return predicate.Event(sql.FieldNotNull(FieldSummary))
}

// SummaryEqualFold applies the EqualFold predicate on the "summary" field.
func SummaryEqualFold(v string) predicate.Event {
	return predicate.Event(sql.FieldEqualFold(FieldSummary, v))
}

// SummaryContainsFold applies the ContainsFold predicate on the "summary" field.
func SummaryContainsFold(v string) predicate.Event {
	return predicate.Event(sql.FieldContainsFold(FieldSummary, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Event {
	return predicate.Event(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Event {
	return predicate.Event(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Event {
	return predicate.Event(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Event {
	return predicate.Event(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Event {
	return predicate.Event(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Event {
	return predicate.Event(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Event {
	return predicate.Event(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Event {
	return predicate.Event(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Event {
	return predicate.Event(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Event {
	return predicate.Event(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.Event {
	return predicate.Event(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.Event {
	return predicate.Event(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Event {
	return predicate.Event(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Event {
	return predicate.Event(sql.FieldContainsFold(FieldDescription, v))
}

// LocationEQ applies the EQ predicate on the "location" field.
func LocationEQ(v string) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldLocation, v))
}

// LocationNEQ applies the NEQ predicate on the "location" field.
func LocationNEQ(v string) predicate.Event {
	return predicate.Event(sql.FieldNEQ(FieldLocation, v))
}

// LocationIn applies the In predicate on the "location" field.
func LocationIn(vs ...string) predicate.Event {
	return predicate.Event(sql.FieldIn(FieldLocation, vs...))
}

// LocationNotIn applies the NotIn predicate on the "location" field.
func LocationNotIn(vs ...string) predicate.Event {
	return predicate.Event(sql.FieldNotIn(FieldLocation, vs...))
}

// LocationGT applies the GT predicate on the "location" field.
func LocationGT(v string) predicate.Event {
	return predicate.Event(sql.FieldGT(FieldLocation, v))
}

// LocationGTE applies the GTE predicate on the "location" field.
func LocationGTE(v string) predicate.Event {
	return predicate.Event(sql.FieldGTE(FieldLocation, v))
}

// LocationLT applies the LT predicate on the "location" field.
func LocationLT(v string) predicate.Event {
	return predicate.Event(sql.FieldLT(FieldLocation, v))
}

// LocationLTE applies the LTE predicate on the "location" field.
func LocationLTE(v string) predicate.Event {
	return predicate.Event(sql.FieldLTE(FieldLocation, v))
}

// LocationContains applies the Contains predicate on the "location" field.
func LocationContains(v string) predicate.Event {
	return predicate.Event(sql.FieldContains(FieldLocation, v))
}

// LocationHasPrefix applies the HasPrefix predicate on the "location" field.
func LocationHasPrefix(v string) predicate.Event {
	return predicate.Event(sql.FieldHasPrefix(FieldLocation, v))
}

// LocationHasSuffix applies the HasSuffix predicate on the "location" field.
func LocationHasSuffix(v string) predicate.Event {
	return predicate.Event(sql.FieldHasSuffix(FieldLocation, v))
}

// LocationIsNil applies the IsNil predicate on the "location" field.
func LocationIsNil() predicate.Event {
	return predicate.Event(sql.FieldIsNull(FieldLocation))
}

// LocationNotNil applies the NotNil predicate on the "location" field.
func LocationNotNil() predicate.Event {
	return predicate.Event(sql.FieldNotNull(FieldLocation))
}

// LocationEqualFold applies the EqualFold predicate on the "location" field.
func LocationEqualFold(v string) predicate.Event {
	return predicate.Event(sql.FieldEqualFold(FieldLocation, v))
}

// LocationContainsFold applies the ContainsFold predicate on the "location" field.
func LocationContainsFold(v string) predicate.Event {
	return predicate.Event(sql.FieldContainsFold(FieldLocation, v))
}

// StartTimeEQ applies the EQ predicate on the "start_time" field.
func StartTimeEQ(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldStartTime, v))
}

// StartTimeNEQ applies the NEQ predicate on the "start_time" field.
func StartTimeNEQ(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldNEQ(FieldStartTime, v))
}

// StartTimeIn applies the In predicate on the "start_time" field.
func StartTimeIn(vs ...time.Time) predicate.Event {
	return predicate.Event(sql.FieldIn(FieldStartTime, vs...))
}

// StartTimeNotIn applies the NotIn predicate on the "start_time" field.
func StartTimeNotIn(vs ...time.Time) predicate.Event {
	return predicate.Event(sql.FieldNotIn(FieldStartTime, vs...))
}

// StartTimeGT applies the GT predicate on the "start_time" field.
func StartTimeGT(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldGT(FieldStartTime, v))
}

// StartTimeGTE applies the GTE predicate on the "start_time" field.
func StartTimeGTE(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldGTE(FieldStartTime, v))
}

// StartTimeLT applies the LT predicate on the "start_time" field.
func StartTimeLT(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldLT(FieldStartTime, v))
}

// StartTimeLTE applies the LTE predicate on the "start_time" field.
func StartTimeLTE(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldLTE(FieldStartTime, v))
}

// StartTimeIsNil applies the IsNil predicate on the "start_time" field.
func StartTimeIsNil() predicate.Event {
	return predicate.Event(sql.FieldIsNull(FieldStartTime))
}

// StartTimeNotNil applies the NotNil predicate on the "start_time" field.
func StartTimeNotNil() predicate.Event {
	return predicate.Event(sql.FieldNotNull(FieldStartTime))
}

// EndTimeEQ applies the EQ predicate on the "end_time" field.
func EndTimeEQ(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldEQ(FieldEndTime, v))
}

// EndTimeNEQ applies the NEQ predicate on the "end_time" field.
func EndTimeNEQ(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldNEQ(FieldEndTime, v))
}

// EndTimeIn applies the In predicate on the "end_time" field.
func EndTimeIn(vs ...time.Time) predicate.Event {
	return predicate.Event(sql.FieldIn(FieldEndTime, vs...))
}

// EndTimeNotIn applies the NotIn predicate on the "end_time" field.
func EndTimeNotIn(vs ...time.Time) predicate.Event {
	return predicate.Event(sql.FieldNotIn(FieldEndTime, vs...))
}

// EndTimeGT applies the GT predicate on the "end_time" field.
func EndTimeGT(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldGT(FieldEndTime, v))
}

// EndTimeGTE applies the GTE predicate on the "end_time" field.
func EndTimeGTE(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldGTE(FieldEndTime, v))
}

// EndTimeLT applies the LT predicate on the "end_time" field.
func EndTimeLT(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldLT(FieldEndTime, v))
}

// EndTimeLTE applies the LTE predicate on the "end_time" field.
func EndTimeLTE(v time.Time) predicate.Event {
	return predicate.Event(sql.FieldLTE(FieldEndTime, v))
}

// EndTimeIsNil applies the IsNil predicate on the "end_time" field.
func EndTimeIsNil() predicate.Event {
	return predicate.Event(sql.FieldIsNull(FieldEndTime))
}

// EndTimeNotNil applies the NotNil predicate on the "end_time" field.
func EndTimeNotNil() predicate.Event {
	return predicate.Event(sql.FieldNotNull(FieldEndTime))
}

// HasCalendar applies the HasEdge predicate on the "calendar" edge.
func HasCalendar() predicate.Event {
	return predicate.Event(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, CalendarTable, CalendarColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCalendarWith applies the HasEdge predicate on the "calendar" edge with a given conditions (other predicates).
func HasCalendarWith(preds ...predicate.Calendar) predicate.Event {
	return predicate.Event(func(s *sql.Selector) {
		step := newCalendarStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Event) predicate.Event {
	return predicate.Event(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Event) predicate.Event {
	return predicate.Event(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Event) predicate.Event {
	return predicate.Event(sql.NotPredicates(p))
}
