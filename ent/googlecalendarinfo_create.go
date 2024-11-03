// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/googlecalendarinfo"
)

// GoogleCalendarInfoCreate is the builder for creating a GoogleCalendarInfo entity.
type GoogleCalendarInfoCreate struct {
	config
	mutation *GoogleCalendarInfoMutation
	hooks    []Hook
}

// SetGoogleCalendarID sets the "google_calendar_id" field.
func (gcic *GoogleCalendarInfoCreate) SetGoogleCalendarID(s string) *GoogleCalendarInfoCreate {
	gcic.mutation.SetGoogleCalendarID(s)
	return gcic
}

// SetSummary sets the "summary" field.
func (gcic *GoogleCalendarInfoCreate) SetSummary(s string) *GoogleCalendarInfoCreate {
	gcic.mutation.SetSummary(s)
	return gcic
}

// SetNillableSummary sets the "summary" field if the given value is not nil.
func (gcic *GoogleCalendarInfoCreate) SetNillableSummary(s *string) *GoogleCalendarInfoCreate {
	if s != nil {
		gcic.SetSummary(*s)
	}
	return gcic
}

// SetIsPrimary sets the "is_primary" field.
func (gcic *GoogleCalendarInfoCreate) SetIsPrimary(b bool) *GoogleCalendarInfoCreate {
	gcic.mutation.SetIsPrimary(b)
	return gcic
}

// SetNillableIsPrimary sets the "is_primary" field if the given value is not nil.
func (gcic *GoogleCalendarInfoCreate) SetNillableIsPrimary(b *bool) *GoogleCalendarInfoCreate {
	if b != nil {
		gcic.SetIsPrimary(*b)
	}
	return gcic
}

// SetID sets the "id" field.
func (gcic *GoogleCalendarInfoCreate) SetID(u uuid.UUID) *GoogleCalendarInfoCreate {
	gcic.mutation.SetID(u)
	return gcic
}

// SetNillableID sets the "id" field if the given value is not nil.
func (gcic *GoogleCalendarInfoCreate) SetNillableID(u *uuid.UUID) *GoogleCalendarInfoCreate {
	if u != nil {
		gcic.SetID(*u)
	}
	return gcic
}

// AddCalendarIDs adds the "calendars" edge to the Calendar entity by IDs.
func (gcic *GoogleCalendarInfoCreate) AddCalendarIDs(ids ...uuid.UUID) *GoogleCalendarInfoCreate {
	gcic.mutation.AddCalendarIDs(ids...)
	return gcic
}

// AddCalendars adds the "calendars" edges to the Calendar entity.
func (gcic *GoogleCalendarInfoCreate) AddCalendars(c ...*Calendar) *GoogleCalendarInfoCreate {
	ids := make([]uuid.UUID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return gcic.AddCalendarIDs(ids...)
}

// Mutation returns the GoogleCalendarInfoMutation object of the builder.
func (gcic *GoogleCalendarInfoCreate) Mutation() *GoogleCalendarInfoMutation {
	return gcic.mutation
}

// Save creates the GoogleCalendarInfo in the database.
func (gcic *GoogleCalendarInfoCreate) Save(ctx context.Context) (*GoogleCalendarInfo, error) {
	gcic.defaults()
	return withHooks(ctx, gcic.sqlSave, gcic.mutation, gcic.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (gcic *GoogleCalendarInfoCreate) SaveX(ctx context.Context) *GoogleCalendarInfo {
	v, err := gcic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gcic *GoogleCalendarInfoCreate) Exec(ctx context.Context) error {
	_, err := gcic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gcic *GoogleCalendarInfoCreate) ExecX(ctx context.Context) {
	if err := gcic.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gcic *GoogleCalendarInfoCreate) defaults() {
	if _, ok := gcic.mutation.IsPrimary(); !ok {
		v := googlecalendarinfo.DefaultIsPrimary
		gcic.mutation.SetIsPrimary(v)
	}
	if _, ok := gcic.mutation.ID(); !ok {
		v := googlecalendarinfo.DefaultID()
		gcic.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (gcic *GoogleCalendarInfoCreate) check() error {
	if _, ok := gcic.mutation.GoogleCalendarID(); !ok {
		return &ValidationError{Name: "google_calendar_id", err: errors.New(`ent: missing required field "GoogleCalendarInfo.google_calendar_id"`)}
	}
	if v, ok := gcic.mutation.GoogleCalendarID(); ok {
		if err := googlecalendarinfo.GoogleCalendarIDValidator(v); err != nil {
			return &ValidationError{Name: "google_calendar_id", err: fmt.Errorf(`ent: validator failed for field "GoogleCalendarInfo.google_calendar_id": %w`, err)}
		}
	}
	if _, ok := gcic.mutation.IsPrimary(); !ok {
		return &ValidationError{Name: "is_primary", err: errors.New(`ent: missing required field "GoogleCalendarInfo.is_primary"`)}
	}
	return nil
}

func (gcic *GoogleCalendarInfoCreate) sqlSave(ctx context.Context) (*GoogleCalendarInfo, error) {
	if err := gcic.check(); err != nil {
		return nil, err
	}
	_node, _spec := gcic.createSpec()
	if err := sqlgraph.CreateNode(ctx, gcic.driver, _spec); err != nil {
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
	gcic.mutation.id = &_node.ID
	gcic.mutation.done = true
	return _node, nil
}

func (gcic *GoogleCalendarInfoCreate) createSpec() (*GoogleCalendarInfo, *sqlgraph.CreateSpec) {
	var (
		_node = &GoogleCalendarInfo{config: gcic.config}
		_spec = sqlgraph.NewCreateSpec(googlecalendarinfo.Table, sqlgraph.NewFieldSpec(googlecalendarinfo.FieldID, field.TypeUUID))
	)
	if id, ok := gcic.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := gcic.mutation.GoogleCalendarID(); ok {
		_spec.SetField(googlecalendarinfo.FieldGoogleCalendarID, field.TypeString, value)
		_node.GoogleCalendarID = value
	}
	if value, ok := gcic.mutation.Summary(); ok {
		_spec.SetField(googlecalendarinfo.FieldSummary, field.TypeString, value)
		_node.Summary = value
	}
	if value, ok := gcic.mutation.IsPrimary(); ok {
		_spec.SetField(googlecalendarinfo.FieldIsPrimary, field.TypeBool, value)
		_node.IsPrimary = value
	}
	if nodes := gcic.mutation.CalendarsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   googlecalendarinfo.CalendarsTable,
			Columns: googlecalendarinfo.CalendarsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(calendar.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// GoogleCalendarInfoCreateBulk is the builder for creating many GoogleCalendarInfo entities in bulk.
type GoogleCalendarInfoCreateBulk struct {
	config
	err      error
	builders []*GoogleCalendarInfoCreate
}

// Save creates the GoogleCalendarInfo entities in the database.
func (gcicb *GoogleCalendarInfoCreateBulk) Save(ctx context.Context) ([]*GoogleCalendarInfo, error) {
	if gcicb.err != nil {
		return nil, gcicb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(gcicb.builders))
	nodes := make([]*GoogleCalendarInfo, len(gcicb.builders))
	mutators := make([]Mutator, len(gcicb.builders))
	for i := range gcicb.builders {
		func(i int, root context.Context) {
			builder := gcicb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*GoogleCalendarInfoMutation)
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
					_, err = mutators[i+1].Mutate(root, gcicb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, gcicb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, gcicb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (gcicb *GoogleCalendarInfoCreateBulk) SaveX(ctx context.Context) []*GoogleCalendarInfo {
	v, err := gcicb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gcicb *GoogleCalendarInfoCreateBulk) Exec(ctx context.Context) error {
	_, err := gcicb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gcicb *GoogleCalendarInfoCreateBulk) ExecX(ctx context.Context) {
	if err := gcicb.Exec(ctx); err != nil {
		panic(err)
	}
}