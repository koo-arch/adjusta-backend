// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/account"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/event"
)

// CalendarCreate is the builder for creating a Calendar entity.
type CalendarCreate struct {
	config
	mutation *CalendarMutation
	hooks    []Hook
}

// SetCalendarID sets the "calendar_id" field.
func (cc *CalendarCreate) SetCalendarID(s string) *CalendarCreate {
	cc.mutation.SetCalendarID(s)
	return cc
}

// SetSummary sets the "summary" field.
func (cc *CalendarCreate) SetSummary(s string) *CalendarCreate {
	cc.mutation.SetSummary(s)
	return cc
}

// SetIsPrimary sets the "is_primary" field.
func (cc *CalendarCreate) SetIsPrimary(b bool) *CalendarCreate {
	cc.mutation.SetIsPrimary(b)
	return cc
}

// SetNillableIsPrimary sets the "is_primary" field if the given value is not nil.
func (cc *CalendarCreate) SetNillableIsPrimary(b *bool) *CalendarCreate {
	if b != nil {
		cc.SetIsPrimary(*b)
	}
	return cc
}

// SetID sets the "id" field.
func (cc *CalendarCreate) SetID(u uuid.UUID) *CalendarCreate {
	cc.mutation.SetID(u)
	return cc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (cc *CalendarCreate) SetNillableID(u *uuid.UUID) *CalendarCreate {
	if u != nil {
		cc.SetID(*u)
	}
	return cc
}

// SetAccountID sets the "account" edge to the Account entity by ID.
func (cc *CalendarCreate) SetAccountID(id uuid.UUID) *CalendarCreate {
	cc.mutation.SetAccountID(id)
	return cc
}

// SetNillableAccountID sets the "account" edge to the Account entity by ID if the given value is not nil.
func (cc *CalendarCreate) SetNillableAccountID(id *uuid.UUID) *CalendarCreate {
	if id != nil {
		cc = cc.SetAccountID(*id)
	}
	return cc
}

// SetAccount sets the "account" edge to the Account entity.
func (cc *CalendarCreate) SetAccount(a *Account) *CalendarCreate {
	return cc.SetAccountID(a.ID)
}

// AddEventIDs adds the "events" edge to the Event entity by IDs.
func (cc *CalendarCreate) AddEventIDs(ids ...uuid.UUID) *CalendarCreate {
	cc.mutation.AddEventIDs(ids...)
	return cc
}

// AddEvents adds the "events" edges to the Event entity.
func (cc *CalendarCreate) AddEvents(e ...*Event) *CalendarCreate {
	ids := make([]uuid.UUID, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return cc.AddEventIDs(ids...)
}

// Mutation returns the CalendarMutation object of the builder.
func (cc *CalendarCreate) Mutation() *CalendarMutation {
	return cc.mutation
}

// Save creates the Calendar in the database.
func (cc *CalendarCreate) Save(ctx context.Context) (*Calendar, error) {
	if err := cc.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, cc.sqlSave, cc.mutation, cc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *CalendarCreate) SaveX(ctx context.Context) *Calendar {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *CalendarCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *CalendarCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cc *CalendarCreate) defaults() error {
	if _, ok := cc.mutation.IsPrimary(); !ok {
		v := calendar.DefaultIsPrimary
		cc.mutation.SetIsPrimary(v)
	}
	if _, ok := cc.mutation.ID(); !ok {
		if calendar.DefaultID == nil {
			return fmt.Errorf("ent: uninitialized calendar.DefaultID (forgotten import ent/runtime?)")
		}
		v := calendar.DefaultID()
		cc.mutation.SetID(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (cc *CalendarCreate) check() error {
	if _, ok := cc.mutation.CalendarID(); !ok {
		return &ValidationError{Name: "calendar_id", err: errors.New(`ent: missing required field "Calendar.calendar_id"`)}
	}
	if _, ok := cc.mutation.Summary(); !ok {
		return &ValidationError{Name: "summary", err: errors.New(`ent: missing required field "Calendar.summary"`)}
	}
	if _, ok := cc.mutation.IsPrimary(); !ok {
		return &ValidationError{Name: "is_primary", err: errors.New(`ent: missing required field "Calendar.is_primary"`)}
	}
	return nil
}

func (cc *CalendarCreate) sqlSave(ctx context.Context) (*Calendar, error) {
	if err := cc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	cc.mutation.id = &_node.ID
	cc.mutation.done = true
	return _node, nil
}

func (cc *CalendarCreate) createSpec() (*Calendar, *sqlgraph.CreateSpec) {
	var (
		_node = &Calendar{config: cc.config}
		_spec = sqlgraph.NewCreateSpec(calendar.Table, sqlgraph.NewFieldSpec(calendar.FieldID, field.TypeUUID))
	)
	if id, ok := cc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := cc.mutation.CalendarID(); ok {
		_spec.SetField(calendar.FieldCalendarID, field.TypeString, value)
		_node.CalendarID = value
	}
	if value, ok := cc.mutation.Summary(); ok {
		_spec.SetField(calendar.FieldSummary, field.TypeString, value)
		_node.Summary = value
	}
	if value, ok := cc.mutation.IsPrimary(); ok {
		_spec.SetField(calendar.FieldIsPrimary, field.TypeBool, value)
		_node.IsPrimary = value
	}
	if nodes := cc.mutation.AccountIDs(); len(nodes) > 0 {
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
		_node.account_calendars = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := cc.mutation.EventsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// CalendarCreateBulk is the builder for creating many Calendar entities in bulk.
type CalendarCreateBulk struct {
	config
	err      error
	builders []*CalendarCreate
}

// Save creates the Calendar entities in the database.
func (ccb *CalendarCreateBulk) Save(ctx context.Context) ([]*Calendar, error) {
	if ccb.err != nil {
		return nil, ccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Calendar, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CalendarMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *CalendarCreateBulk) SaveX(ctx context.Context) []*Calendar {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *CalendarCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *CalendarCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}
