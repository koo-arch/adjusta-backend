// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/predicate"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
)

// EventUpdate is the builder for updating Event entities.
type EventUpdate struct {
	config
	hooks    []Hook
	mutation *EventMutation
}

// Where appends a list predicates to the EventUpdate builder.
func (eu *EventUpdate) Where(ps ...predicate.Event) *EventUpdate {
	eu.mutation.Where(ps...)
	return eu
}

// SetEventID sets the "event_id" field.
func (eu *EventUpdate) SetEventID(s string) *EventUpdate {
	eu.mutation.SetEventID(s)
	return eu
}

// SetNillableEventID sets the "event_id" field if the given value is not nil.
func (eu *EventUpdate) SetNillableEventID(s *string) *EventUpdate {
	if s != nil {
		eu.SetEventID(*s)
	}
	return eu
}

// SetSummary sets the "summary" field.
func (eu *EventUpdate) SetSummary(s string) *EventUpdate {
	eu.mutation.SetSummary(s)
	return eu
}

// SetNillableSummary sets the "summary" field if the given value is not nil.
func (eu *EventUpdate) SetNillableSummary(s *string) *EventUpdate {
	if s != nil {
		eu.SetSummary(*s)
	}
	return eu
}

// ClearSummary clears the value of the "summary" field.
func (eu *EventUpdate) ClearSummary() *EventUpdate {
	eu.mutation.ClearSummary()
	return eu
}

// SetDescription sets the "description" field.
func (eu *EventUpdate) SetDescription(s string) *EventUpdate {
	eu.mutation.SetDescription(s)
	return eu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (eu *EventUpdate) SetNillableDescription(s *string) *EventUpdate {
	if s != nil {
		eu.SetDescription(*s)
	}
	return eu
}

// ClearDescription clears the value of the "description" field.
func (eu *EventUpdate) ClearDescription() *EventUpdate {
	eu.mutation.ClearDescription()
	return eu
}

// SetLocation sets the "location" field.
func (eu *EventUpdate) SetLocation(s string) *EventUpdate {
	eu.mutation.SetLocation(s)
	return eu
}

// SetNillableLocation sets the "location" field if the given value is not nil.
func (eu *EventUpdate) SetNillableLocation(s *string) *EventUpdate {
	if s != nil {
		eu.SetLocation(*s)
	}
	return eu
}

// ClearLocation clears the value of the "location" field.
func (eu *EventUpdate) ClearLocation() *EventUpdate {
	eu.mutation.ClearLocation()
	return eu
}

// SetCalendarID sets the "calendar" edge to the Calendar entity by ID.
func (eu *EventUpdate) SetCalendarID(id uuid.UUID) *EventUpdate {
	eu.mutation.SetCalendarID(id)
	return eu
}

// SetNillableCalendarID sets the "calendar" edge to the Calendar entity by ID if the given value is not nil.
func (eu *EventUpdate) SetNillableCalendarID(id *uuid.UUID) *EventUpdate {
	if id != nil {
		eu = eu.SetCalendarID(*id)
	}
	return eu
}

// SetCalendar sets the "calendar" edge to the Calendar entity.
func (eu *EventUpdate) SetCalendar(c *Calendar) *EventUpdate {
	return eu.SetCalendarID(c.ID)
}

// AddProposedDateIDs adds the "proposed_dates" edge to the ProposedDate entity by IDs.
func (eu *EventUpdate) AddProposedDateIDs(ids ...uuid.UUID) *EventUpdate {
	eu.mutation.AddProposedDateIDs(ids...)
	return eu
}

// AddProposedDates adds the "proposed_dates" edges to the ProposedDate entity.
func (eu *EventUpdate) AddProposedDates(p ...*ProposedDate) *EventUpdate {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eu.AddProposedDateIDs(ids...)
}

// Mutation returns the EventMutation object of the builder.
func (eu *EventUpdate) Mutation() *EventMutation {
	return eu.mutation
}

// ClearCalendar clears the "calendar" edge to the Calendar entity.
func (eu *EventUpdate) ClearCalendar() *EventUpdate {
	eu.mutation.ClearCalendar()
	return eu
}

// ClearProposedDates clears all "proposed_dates" edges to the ProposedDate entity.
func (eu *EventUpdate) ClearProposedDates() *EventUpdate {
	eu.mutation.ClearProposedDates()
	return eu
}

// RemoveProposedDateIDs removes the "proposed_dates" edge to ProposedDate entities by IDs.
func (eu *EventUpdate) RemoveProposedDateIDs(ids ...uuid.UUID) *EventUpdate {
	eu.mutation.RemoveProposedDateIDs(ids...)
	return eu
}

// RemoveProposedDates removes "proposed_dates" edges to ProposedDate entities.
func (eu *EventUpdate) RemoveProposedDates(p ...*ProposedDate) *EventUpdate {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eu.RemoveProposedDateIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (eu *EventUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, eu.sqlSave, eu.mutation, eu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (eu *EventUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *EventUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *EventUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eu *EventUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(event.Table, event.Columns, sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID))
	if ps := eu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eu.mutation.EventID(); ok {
		_spec.SetField(event.FieldEventID, field.TypeString, value)
	}
	if value, ok := eu.mutation.Summary(); ok {
		_spec.SetField(event.FieldSummary, field.TypeString, value)
	}
	if eu.mutation.SummaryCleared() {
		_spec.ClearField(event.FieldSummary, field.TypeString)
	}
	if value, ok := eu.mutation.Description(); ok {
		_spec.SetField(event.FieldDescription, field.TypeString, value)
	}
	if eu.mutation.DescriptionCleared() {
		_spec.ClearField(event.FieldDescription, field.TypeString)
	}
	if value, ok := eu.mutation.Location(); ok {
		_spec.SetField(event.FieldLocation, field.TypeString, value)
	}
	if eu.mutation.LocationCleared() {
		_spec.ClearField(event.FieldLocation, field.TypeString)
	}
	if eu.mutation.CalendarCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   event.CalendarTable,
			Columns: []string{event.CalendarColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(calendar.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.CalendarIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   event.CalendarTable,
			Columns: []string{event.CalendarColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(calendar.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.mutation.ProposedDatesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   event.ProposedDatesTable,
			Columns: []string{event.ProposedDatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(proposeddate.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.RemovedProposedDatesIDs(); len(nodes) > 0 && !eu.mutation.ProposedDatesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   event.ProposedDatesTable,
			Columns: []string{event.ProposedDatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(proposeddate.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.ProposedDatesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   event.ProposedDatesTable,
			Columns: []string{event.ProposedDatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(proposeddate.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, eu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{event.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	eu.mutation.done = true
	return n, nil
}

// EventUpdateOne is the builder for updating a single Event entity.
type EventUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *EventMutation
}

// SetEventID sets the "event_id" field.
func (euo *EventUpdateOne) SetEventID(s string) *EventUpdateOne {
	euo.mutation.SetEventID(s)
	return euo
}

// SetNillableEventID sets the "event_id" field if the given value is not nil.
func (euo *EventUpdateOne) SetNillableEventID(s *string) *EventUpdateOne {
	if s != nil {
		euo.SetEventID(*s)
	}
	return euo
}

// SetSummary sets the "summary" field.
func (euo *EventUpdateOne) SetSummary(s string) *EventUpdateOne {
	euo.mutation.SetSummary(s)
	return euo
}

// SetNillableSummary sets the "summary" field if the given value is not nil.
func (euo *EventUpdateOne) SetNillableSummary(s *string) *EventUpdateOne {
	if s != nil {
		euo.SetSummary(*s)
	}
	return euo
}

// ClearSummary clears the value of the "summary" field.
func (euo *EventUpdateOne) ClearSummary() *EventUpdateOne {
	euo.mutation.ClearSummary()
	return euo
}

// SetDescription sets the "description" field.
func (euo *EventUpdateOne) SetDescription(s string) *EventUpdateOne {
	euo.mutation.SetDescription(s)
	return euo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (euo *EventUpdateOne) SetNillableDescription(s *string) *EventUpdateOne {
	if s != nil {
		euo.SetDescription(*s)
	}
	return euo
}

// ClearDescription clears the value of the "description" field.
func (euo *EventUpdateOne) ClearDescription() *EventUpdateOne {
	euo.mutation.ClearDescription()
	return euo
}

// SetLocation sets the "location" field.
func (euo *EventUpdateOne) SetLocation(s string) *EventUpdateOne {
	euo.mutation.SetLocation(s)
	return euo
}

// SetNillableLocation sets the "location" field if the given value is not nil.
func (euo *EventUpdateOne) SetNillableLocation(s *string) *EventUpdateOne {
	if s != nil {
		euo.SetLocation(*s)
	}
	return euo
}

// ClearLocation clears the value of the "location" field.
func (euo *EventUpdateOne) ClearLocation() *EventUpdateOne {
	euo.mutation.ClearLocation()
	return euo
}

// SetCalendarID sets the "calendar" edge to the Calendar entity by ID.
func (euo *EventUpdateOne) SetCalendarID(id uuid.UUID) *EventUpdateOne {
	euo.mutation.SetCalendarID(id)
	return euo
}

// SetNillableCalendarID sets the "calendar" edge to the Calendar entity by ID if the given value is not nil.
func (euo *EventUpdateOne) SetNillableCalendarID(id *uuid.UUID) *EventUpdateOne {
	if id != nil {
		euo = euo.SetCalendarID(*id)
	}
	return euo
}

// SetCalendar sets the "calendar" edge to the Calendar entity.
func (euo *EventUpdateOne) SetCalendar(c *Calendar) *EventUpdateOne {
	return euo.SetCalendarID(c.ID)
}

// AddProposedDateIDs adds the "proposed_dates" edge to the ProposedDate entity by IDs.
func (euo *EventUpdateOne) AddProposedDateIDs(ids ...uuid.UUID) *EventUpdateOne {
	euo.mutation.AddProposedDateIDs(ids...)
	return euo
}

// AddProposedDates adds the "proposed_dates" edges to the ProposedDate entity.
func (euo *EventUpdateOne) AddProposedDates(p ...*ProposedDate) *EventUpdateOne {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return euo.AddProposedDateIDs(ids...)
}

// Mutation returns the EventMutation object of the builder.
func (euo *EventUpdateOne) Mutation() *EventMutation {
	return euo.mutation
}

// ClearCalendar clears the "calendar" edge to the Calendar entity.
func (euo *EventUpdateOne) ClearCalendar() *EventUpdateOne {
	euo.mutation.ClearCalendar()
	return euo
}

// ClearProposedDates clears all "proposed_dates" edges to the ProposedDate entity.
func (euo *EventUpdateOne) ClearProposedDates() *EventUpdateOne {
	euo.mutation.ClearProposedDates()
	return euo
}

// RemoveProposedDateIDs removes the "proposed_dates" edge to ProposedDate entities by IDs.
func (euo *EventUpdateOne) RemoveProposedDateIDs(ids ...uuid.UUID) *EventUpdateOne {
	euo.mutation.RemoveProposedDateIDs(ids...)
	return euo
}

// RemoveProposedDates removes "proposed_dates" edges to ProposedDate entities.
func (euo *EventUpdateOne) RemoveProposedDates(p ...*ProposedDate) *EventUpdateOne {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return euo.RemoveProposedDateIDs(ids...)
}

// Where appends a list predicates to the EventUpdate builder.
func (euo *EventUpdateOne) Where(ps ...predicate.Event) *EventUpdateOne {
	euo.mutation.Where(ps...)
	return euo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (euo *EventUpdateOne) Select(field string, fields ...string) *EventUpdateOne {
	euo.fields = append([]string{field}, fields...)
	return euo
}

// Save executes the query and returns the updated Event entity.
func (euo *EventUpdateOne) Save(ctx context.Context) (*Event, error) {
	return withHooks(ctx, euo.sqlSave, euo.mutation, euo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (euo *EventUpdateOne) SaveX(ctx context.Context) *Event {
	node, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (euo *EventUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *EventUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (euo *EventUpdateOne) sqlSave(ctx context.Context) (_node *Event, err error) {
	_spec := sqlgraph.NewUpdateSpec(event.Table, event.Columns, sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID))
	id, ok := euo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Event.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := euo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, event.FieldID)
		for _, f := range fields {
			if !event.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != event.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := euo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := euo.mutation.EventID(); ok {
		_spec.SetField(event.FieldEventID, field.TypeString, value)
	}
	if value, ok := euo.mutation.Summary(); ok {
		_spec.SetField(event.FieldSummary, field.TypeString, value)
	}
	if euo.mutation.SummaryCleared() {
		_spec.ClearField(event.FieldSummary, field.TypeString)
	}
	if value, ok := euo.mutation.Description(); ok {
		_spec.SetField(event.FieldDescription, field.TypeString, value)
	}
	if euo.mutation.DescriptionCleared() {
		_spec.ClearField(event.FieldDescription, field.TypeString)
	}
	if value, ok := euo.mutation.Location(); ok {
		_spec.SetField(event.FieldLocation, field.TypeString, value)
	}
	if euo.mutation.LocationCleared() {
		_spec.ClearField(event.FieldLocation, field.TypeString)
	}
	if euo.mutation.CalendarCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   event.CalendarTable,
			Columns: []string{event.CalendarColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(calendar.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.CalendarIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   event.CalendarTable,
			Columns: []string{event.CalendarColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(calendar.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.mutation.ProposedDatesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   event.ProposedDatesTable,
			Columns: []string{event.ProposedDatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(proposeddate.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.RemovedProposedDatesIDs(); len(nodes) > 0 && !euo.mutation.ProposedDatesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   event.ProposedDatesTable,
			Columns: []string{event.ProposedDatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(proposeddate.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.ProposedDatesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   event.ProposedDatesTable,
			Columns: []string{event.ProposedDatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(proposeddate.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Event{config: euo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, euo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{event.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	euo.mutation.done = true
	return _node, nil
}
