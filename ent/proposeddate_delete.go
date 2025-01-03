// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/koo-arch/adjusta-backend/ent/predicate"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
)

// ProposedDateDelete is the builder for deleting a ProposedDate entity.
type ProposedDateDelete struct {
	config
	hooks    []Hook
	mutation *ProposedDateMutation
}

// Where appends a list predicates to the ProposedDateDelete builder.
func (pdd *ProposedDateDelete) Where(ps ...predicate.ProposedDate) *ProposedDateDelete {
	pdd.mutation.Where(ps...)
	return pdd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pdd *ProposedDateDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, pdd.sqlExec, pdd.mutation, pdd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (pdd *ProposedDateDelete) ExecX(ctx context.Context) int {
	n, err := pdd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pdd *ProposedDateDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(proposeddate.Table, sqlgraph.NewFieldSpec(proposeddate.FieldID, field.TypeUUID))
	if ps := pdd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pdd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	pdd.mutation.done = true
	return affected, err
}

// ProposedDateDeleteOne is the builder for deleting a single ProposedDate entity.
type ProposedDateDeleteOne struct {
	pdd *ProposedDateDelete
}

// Where appends a list predicates to the ProposedDateDelete builder.
func (pddo *ProposedDateDeleteOne) Where(ps ...predicate.ProposedDate) *ProposedDateDeleteOne {
	pddo.pdd.mutation.Where(ps...)
	return pddo
}

// Exec executes the deletion query.
func (pddo *ProposedDateDeleteOne) Exec(ctx context.Context) error {
	n, err := pddo.pdd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{proposeddate.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pddo *ProposedDateDeleteOne) ExecX(ctx context.Context) {
	if err := pddo.Exec(ctx); err != nil {
		panic(err)
	}
}
