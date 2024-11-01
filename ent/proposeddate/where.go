// Code generated by ent, DO NOT EDIT.

package proposeddate

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldLTE(FieldID, id))
}

// GoogleEventID applies equality check predicate on the "google_event_id" field. It's identical to GoogleEventIDEQ.
func GoogleEventID(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldEQ(FieldGoogleEventID, v))
}

// StartTime applies equality check predicate on the "start_time" field. It's identical to StartTimeEQ.
func StartTime(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldEQ(FieldStartTime, v))
}

// EndTime applies equality check predicate on the "end_time" field. It's identical to EndTimeEQ.
func EndTime(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldEQ(FieldEndTime, v))
}

// Priority applies equality check predicate on the "priority" field. It's identical to PriorityEQ.
func Priority(v int) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldEQ(FieldPriority, v))
}

// GoogleEventIDEQ applies the EQ predicate on the "google_event_id" field.
func GoogleEventIDEQ(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldEQ(FieldGoogleEventID, v))
}

// GoogleEventIDNEQ applies the NEQ predicate on the "google_event_id" field.
func GoogleEventIDNEQ(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldNEQ(FieldGoogleEventID, v))
}

// GoogleEventIDIn applies the In predicate on the "google_event_id" field.
func GoogleEventIDIn(vs ...string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldIn(FieldGoogleEventID, vs...))
}

// GoogleEventIDNotIn applies the NotIn predicate on the "google_event_id" field.
func GoogleEventIDNotIn(vs ...string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldNotIn(FieldGoogleEventID, vs...))
}

// GoogleEventIDGT applies the GT predicate on the "google_event_id" field.
func GoogleEventIDGT(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldGT(FieldGoogleEventID, v))
}

// GoogleEventIDGTE applies the GTE predicate on the "google_event_id" field.
func GoogleEventIDGTE(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldGTE(FieldGoogleEventID, v))
}

// GoogleEventIDLT applies the LT predicate on the "google_event_id" field.
func GoogleEventIDLT(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldLT(FieldGoogleEventID, v))
}

// GoogleEventIDLTE applies the LTE predicate on the "google_event_id" field.
func GoogleEventIDLTE(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldLTE(FieldGoogleEventID, v))
}

// GoogleEventIDContains applies the Contains predicate on the "google_event_id" field.
func GoogleEventIDContains(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldContains(FieldGoogleEventID, v))
}

// GoogleEventIDHasPrefix applies the HasPrefix predicate on the "google_event_id" field.
func GoogleEventIDHasPrefix(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldHasPrefix(FieldGoogleEventID, v))
}

// GoogleEventIDHasSuffix applies the HasSuffix predicate on the "google_event_id" field.
func GoogleEventIDHasSuffix(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldHasSuffix(FieldGoogleEventID, v))
}

// GoogleEventIDIsNil applies the IsNil predicate on the "google_event_id" field.
func GoogleEventIDIsNil() predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldIsNull(FieldGoogleEventID))
}

// GoogleEventIDNotNil applies the NotNil predicate on the "google_event_id" field.
func GoogleEventIDNotNil() predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldNotNull(FieldGoogleEventID))
}

// GoogleEventIDEqualFold applies the EqualFold predicate on the "google_event_id" field.
func GoogleEventIDEqualFold(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldEqualFold(FieldGoogleEventID, v))
}

// GoogleEventIDContainsFold applies the ContainsFold predicate on the "google_event_id" field.
func GoogleEventIDContainsFold(v string) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldContainsFold(FieldGoogleEventID, v))
}

// StartTimeEQ applies the EQ predicate on the "start_time" field.
func StartTimeEQ(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldEQ(FieldStartTime, v))
}

// StartTimeNEQ applies the NEQ predicate on the "start_time" field.
func StartTimeNEQ(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldNEQ(FieldStartTime, v))
}

// StartTimeIn applies the In predicate on the "start_time" field.
func StartTimeIn(vs ...time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldIn(FieldStartTime, vs...))
}

// StartTimeNotIn applies the NotIn predicate on the "start_time" field.
func StartTimeNotIn(vs ...time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldNotIn(FieldStartTime, vs...))
}

// StartTimeGT applies the GT predicate on the "start_time" field.
func StartTimeGT(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldGT(FieldStartTime, v))
}

// StartTimeGTE applies the GTE predicate on the "start_time" field.
func StartTimeGTE(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldGTE(FieldStartTime, v))
}

// StartTimeLT applies the LT predicate on the "start_time" field.
func StartTimeLT(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldLT(FieldStartTime, v))
}

// StartTimeLTE applies the LTE predicate on the "start_time" field.
func StartTimeLTE(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldLTE(FieldStartTime, v))
}

// EndTimeEQ applies the EQ predicate on the "end_time" field.
func EndTimeEQ(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldEQ(FieldEndTime, v))
}

// EndTimeNEQ applies the NEQ predicate on the "end_time" field.
func EndTimeNEQ(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldNEQ(FieldEndTime, v))
}

// EndTimeIn applies the In predicate on the "end_time" field.
func EndTimeIn(vs ...time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldIn(FieldEndTime, vs...))
}

// EndTimeNotIn applies the NotIn predicate on the "end_time" field.
func EndTimeNotIn(vs ...time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldNotIn(FieldEndTime, vs...))
}

// EndTimeGT applies the GT predicate on the "end_time" field.
func EndTimeGT(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldGT(FieldEndTime, v))
}

// EndTimeGTE applies the GTE predicate on the "end_time" field.
func EndTimeGTE(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldGTE(FieldEndTime, v))
}

// EndTimeLT applies the LT predicate on the "end_time" field.
func EndTimeLT(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldLT(FieldEndTime, v))
}

// EndTimeLTE applies the LTE predicate on the "end_time" field.
func EndTimeLTE(v time.Time) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldLTE(FieldEndTime, v))
}

// PriorityEQ applies the EQ predicate on the "priority" field.
func PriorityEQ(v int) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldEQ(FieldPriority, v))
}

// PriorityNEQ applies the NEQ predicate on the "priority" field.
func PriorityNEQ(v int) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldNEQ(FieldPriority, v))
}

// PriorityIn applies the In predicate on the "priority" field.
func PriorityIn(vs ...int) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldIn(FieldPriority, vs...))
}

// PriorityNotIn applies the NotIn predicate on the "priority" field.
func PriorityNotIn(vs ...int) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldNotIn(FieldPriority, vs...))
}

// PriorityGT applies the GT predicate on the "priority" field.
func PriorityGT(v int) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldGT(FieldPriority, v))
}

// PriorityGTE applies the GTE predicate on the "priority" field.
func PriorityGTE(v int) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldGTE(FieldPriority, v))
}

// PriorityLT applies the LT predicate on the "priority" field.
func PriorityLT(v int) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldLT(FieldPriority, v))
}

// PriorityLTE applies the LTE predicate on the "priority" field.
func PriorityLTE(v int) predicate.ProposedDate {
	return predicate.ProposedDate(sql.FieldLTE(FieldPriority, v))
}

// HasEvent applies the HasEdge predicate on the "event" edge.
func HasEvent() predicate.ProposedDate {
	return predicate.ProposedDate(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EventTable, EventColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEventWith applies the HasEdge predicate on the "event" edge with a given conditions (other predicates).
func HasEventWith(preds ...predicate.Event) predicate.ProposedDate {
	return predicate.ProposedDate(func(s *sql.Selector) {
		step := newEventStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ProposedDate) predicate.ProposedDate {
	return predicate.ProposedDate(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ProposedDate) predicate.ProposedDate {
	return predicate.ProposedDate(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ProposedDate) predicate.ProposedDate {
	return predicate.ProposedDate(sql.NotPredicates(p))
}
