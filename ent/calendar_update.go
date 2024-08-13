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
	"github.com/koo-arch/adjusta-backend/ent/account"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/predicate"
)

// CalendarUpdate is the builder for updating Calendar entities.
type CalendarUpdate struct {
	config
	hooks    []Hook
	mutation *CalendarMutation
}

// Where appends a list predicates to the CalendarUpdate builder.
func (cu *CalendarUpdate) Where(ps ...predicate.Calendar) *CalendarUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetCalendarID sets the "calendar_id" field.
func (cu *CalendarUpdate) SetCalendarID(s string) *CalendarUpdate {
	cu.mutation.SetCalendarID(s)
	return cu
}

// SetNillableCalendarID sets the "calendar_id" field if the given value is not nil.
func (cu *CalendarUpdate) SetNillableCalendarID(s *string) *CalendarUpdate {
	if s != nil {
		cu.SetCalendarID(*s)
	}
	return cu
}

// SetSummary sets the "summary" field.
func (cu *CalendarUpdate) SetSummary(s string) *CalendarUpdate {
	cu.mutation.SetSummary(s)
	return cu
}

// SetNillableSummary sets the "summary" field if the given value is not nil.
func (cu *CalendarUpdate) SetNillableSummary(s *string) *CalendarUpdate {
	if s != nil {
		cu.SetSummary(*s)
	}
	return cu
}

// SetAccountID sets the "account" edge to the Account entity by ID.
func (cu *CalendarUpdate) SetAccountID(id uuid.UUID) *CalendarUpdate {
	cu.mutation.SetAccountID(id)
	return cu
}

// SetNillableAccountID sets the "account" edge to the Account entity by ID if the given value is not nil.
func (cu *CalendarUpdate) SetNillableAccountID(id *uuid.UUID) *CalendarUpdate {
	if id != nil {
		cu = cu.SetAccountID(*id)
	}
	return cu
}

// SetAccount sets the "account" edge to the Account entity.
func (cu *CalendarUpdate) SetAccount(a *Account) *CalendarUpdate {
	return cu.SetAccountID(a.ID)
}

// AddEventIDs adds the "events" edge to the Event entity by IDs.
func (cu *CalendarUpdate) AddEventIDs(ids ...uuid.UUID) *CalendarUpdate {
	cu.mutation.AddEventIDs(ids...)
	return cu
}

// AddEvents adds the "events" edges to the Event entity.
func (cu *CalendarUpdate) AddEvents(e ...*Event) *CalendarUpdate {
	ids := make([]uuid.UUID, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return cu.AddEventIDs(ids...)
}

// Mutation returns the CalendarMutation object of the builder.
func (cu *CalendarUpdate) Mutation() *CalendarMutation {
	return cu.mutation
}

// ClearAccount clears the "account" edge to the Account entity.
func (cu *CalendarUpdate) ClearAccount() *CalendarUpdate {
	cu.mutation.ClearAccount()
	return cu
}

// ClearEvents clears all "events" edges to the Event entity.
func (cu *CalendarUpdate) ClearEvents() *CalendarUpdate {
	cu.mutation.ClearEvents()
	return cu
}

// RemoveEventIDs removes the "events" edge to Event entities by IDs.
func (cu *CalendarUpdate) RemoveEventIDs(ids ...uuid.UUID) *CalendarUpdate {
	cu.mutation.RemoveEventIDs(ids...)
	return cu
}

// RemoveEvents removes "events" edges to Event entities.
func (cu *CalendarUpdate) RemoveEvents(e ...*Event) *CalendarUpdate {
	ids := make([]uuid.UUID, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return cu.RemoveEventIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *CalendarUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, cu.sqlSave, cu.mutation, cu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CalendarUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CalendarUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CalendarUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cu *CalendarUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(calendar.Table, calendar.Columns, sqlgraph.NewFieldSpec(calendar.FieldID, field.TypeUUID))
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.CalendarID(); ok {
		_spec.SetField(calendar.FieldCalendarID, field.TypeString, value)
	}
	if value, ok := cu.mutation.Summary(); ok {
		_spec.SetField(calendar.FieldSummary, field.TypeString, value)
	}
	if cu.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   calendar.AccountTable,
			Columns: []string{calendar.AccountColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.AccountIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   calendar.AccountTable,
			Columns: []string{calendar.AccountColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cu.mutation.EventsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   calendar.EventsTable,
			Columns: []string{calendar.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedEventsIDs(); len(nodes) > 0 && !cu.mutation.EventsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   calendar.EventsTable,
			Columns: []string{calendar.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.EventsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   calendar.EventsTable,
			Columns: []string{calendar.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{calendar.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cu.mutation.done = true
	return n, nil
}

// CalendarUpdateOne is the builder for updating a single Calendar entity.
type CalendarUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *CalendarMutation
}

// SetCalendarID sets the "calendar_id" field.
func (cuo *CalendarUpdateOne) SetCalendarID(s string) *CalendarUpdateOne {
	cuo.mutation.SetCalendarID(s)
	return cuo
}

// SetNillableCalendarID sets the "calendar_id" field if the given value is not nil.
func (cuo *CalendarUpdateOne) SetNillableCalendarID(s *string) *CalendarUpdateOne {
	if s != nil {
		cuo.SetCalendarID(*s)
	}
	return cuo
}

// SetSummary sets the "summary" field.
func (cuo *CalendarUpdateOne) SetSummary(s string) *CalendarUpdateOne {
	cuo.mutation.SetSummary(s)
	return cuo
}

// SetNillableSummary sets the "summary" field if the given value is not nil.
func (cuo *CalendarUpdateOne) SetNillableSummary(s *string) *CalendarUpdateOne {
	if s != nil {
		cuo.SetSummary(*s)
	}
	return cuo
}

// SetAccountID sets the "account" edge to the Account entity by ID.
func (cuo *CalendarUpdateOne) SetAccountID(id uuid.UUID) *CalendarUpdateOne {
	cuo.mutation.SetAccountID(id)
	return cuo
}

// SetNillableAccountID sets the "account" edge to the Account entity by ID if the given value is not nil.
func (cuo *CalendarUpdateOne) SetNillableAccountID(id *uuid.UUID) *CalendarUpdateOne {
	if id != nil {
		cuo = cuo.SetAccountID(*id)
	}
	return cuo
}

// SetAccount sets the "account" edge to the Account entity.
func (cuo *CalendarUpdateOne) SetAccount(a *Account) *CalendarUpdateOne {
	return cuo.SetAccountID(a.ID)
}

// AddEventIDs adds the "events" edge to the Event entity by IDs.
func (cuo *CalendarUpdateOne) AddEventIDs(ids ...uuid.UUID) *CalendarUpdateOne {
	cuo.mutation.AddEventIDs(ids...)
	return cuo
}

// AddEvents adds the "events" edges to the Event entity.
func (cuo *CalendarUpdateOne) AddEvents(e ...*Event) *CalendarUpdateOne {
	ids := make([]uuid.UUID, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return cuo.AddEventIDs(ids...)
}

// Mutation returns the CalendarMutation object of the builder.
func (cuo *CalendarUpdateOne) Mutation() *CalendarMutation {
	return cuo.mutation
}

// ClearAccount clears the "account" edge to the Account entity.
func (cuo *CalendarUpdateOne) ClearAccount() *CalendarUpdateOne {
	cuo.mutation.ClearAccount()
	return cuo
}

// ClearEvents clears all "events" edges to the Event entity.
func (cuo *CalendarUpdateOne) ClearEvents() *CalendarUpdateOne {
	cuo.mutation.ClearEvents()
	return cuo
}

// RemoveEventIDs removes the "events" edge to Event entities by IDs.
func (cuo *CalendarUpdateOne) RemoveEventIDs(ids ...uuid.UUID) *CalendarUpdateOne {
	cuo.mutation.RemoveEventIDs(ids...)
	return cuo
}

// RemoveEvents removes "events" edges to Event entities.
func (cuo *CalendarUpdateOne) RemoveEvents(e ...*Event) *CalendarUpdateOne {
	ids := make([]uuid.UUID, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return cuo.RemoveEventIDs(ids...)
}

// Where appends a list predicates to the CalendarUpdate builder.
func (cuo *CalendarUpdateOne) Where(ps ...predicate.Calendar) *CalendarUpdateOne {
	cuo.mutation.Where(ps...)
	return cuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *CalendarUpdateOne) Select(field string, fields ...string) *CalendarUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Calendar entity.
func (cuo *CalendarUpdateOne) Save(ctx context.Context) (*Calendar, error) {
	return withHooks(ctx, cuo.sqlSave, cuo.mutation, cuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CalendarUpdateOne) SaveX(ctx context.Context) *Calendar {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *CalendarUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CalendarUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cuo *CalendarUpdateOne) sqlSave(ctx context.Context) (_node *Calendar, err error) {
	_spec := sqlgraph.NewUpdateSpec(calendar.Table, calendar.Columns, sqlgraph.NewFieldSpec(calendar.FieldID, field.TypeUUID))
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Calendar.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, calendar.FieldID)
		for _, f := range fields {
			if !calendar.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != calendar.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.CalendarID(); ok {
		_spec.SetField(calendar.FieldCalendarID, field.TypeString, value)
	}
	if value, ok := cuo.mutation.Summary(); ok {
		_spec.SetField(calendar.FieldSummary, field.TypeString, value)
	}
	if cuo.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   calendar.AccountTable,
			Columns: []string{calendar.AccountColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.AccountIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   calendar.AccountTable,
			Columns: []string{calendar.AccountColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cuo.mutation.EventsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   calendar.EventsTable,
			Columns: []string{calendar.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedEventsIDs(); len(nodes) > 0 && !cuo.mutation.EventsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   calendar.EventsTable,
			Columns: []string{calendar.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.EventsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   calendar.EventsTable,
			Columns: []string{calendar.EventsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(event.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Calendar{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{calendar.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cuo.mutation.done = true
	return _node, nil
}