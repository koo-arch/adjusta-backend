// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/koo-arch/adjusta-backend/ent/jwtkey"
	"github.com/koo-arch/adjusta-backend/ent/predicate"
)

// JWTKeyUpdate is the builder for updating JWTKey entities.
type JWTKeyUpdate struct {
	config
	hooks    []Hook
	mutation *JWTKeyMutation
}

// Where appends a list predicates to the JWTKeyUpdate builder.
func (jku *JWTKeyUpdate) Where(ps ...predicate.JWTKey) *JWTKeyUpdate {
	jku.mutation.Where(ps...)
	return jku
}

// SetKey sets the "key" field.
func (jku *JWTKeyUpdate) SetKey(s string) *JWTKeyUpdate {
	jku.mutation.SetKey(s)
	return jku
}

// SetNillableKey sets the "key" field if the given value is not nil.
func (jku *JWTKeyUpdate) SetNillableKey(s *string) *JWTKeyUpdate {
	if s != nil {
		jku.SetKey(*s)
	}
	return jku
}

// SetType sets the "type" field.
func (jku *JWTKeyUpdate) SetType(s string) *JWTKeyUpdate {
	jku.mutation.SetType(s)
	return jku
}

// SetNillableType sets the "type" field if the given value is not nil.
func (jku *JWTKeyUpdate) SetNillableType(s *string) *JWTKeyUpdate {
	if s != nil {
		jku.SetType(*s)
	}
	return jku
}

// Mutation returns the JWTKeyMutation object of the builder.
func (jku *JWTKeyUpdate) Mutation() *JWTKeyMutation {
	return jku.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (jku *JWTKeyUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, jku.sqlSave, jku.mutation, jku.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (jku *JWTKeyUpdate) SaveX(ctx context.Context) int {
	affected, err := jku.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (jku *JWTKeyUpdate) Exec(ctx context.Context) error {
	_, err := jku.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (jku *JWTKeyUpdate) ExecX(ctx context.Context) {
	if err := jku.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (jku *JWTKeyUpdate) check() error {
	if v, ok := jku.mutation.Key(); ok {
		if err := jwtkey.KeyValidator(v); err != nil {
			return &ValidationError{Name: "key", err: fmt.Errorf(`ent: validator failed for field "JWTKey.key": %w`, err)}
		}
	}
	if v, ok := jku.mutation.GetType(); ok {
		if err := jwtkey.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "JWTKey.type": %w`, err)}
		}
	}
	return nil
}

func (jku *JWTKeyUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := jku.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(jwtkey.Table, jwtkey.Columns, sqlgraph.NewFieldSpec(jwtkey.FieldID, field.TypeInt))
	if ps := jku.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := jku.mutation.Key(); ok {
		_spec.SetField(jwtkey.FieldKey, field.TypeString, value)
	}
	if value, ok := jku.mutation.GetType(); ok {
		_spec.SetField(jwtkey.FieldType, field.TypeString, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, jku.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{jwtkey.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	jku.mutation.done = true
	return n, nil
}

// JWTKeyUpdateOne is the builder for updating a single JWTKey entity.
type JWTKeyUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *JWTKeyMutation
}

// SetKey sets the "key" field.
func (jkuo *JWTKeyUpdateOne) SetKey(s string) *JWTKeyUpdateOne {
	jkuo.mutation.SetKey(s)
	return jkuo
}

// SetNillableKey sets the "key" field if the given value is not nil.
func (jkuo *JWTKeyUpdateOne) SetNillableKey(s *string) *JWTKeyUpdateOne {
	if s != nil {
		jkuo.SetKey(*s)
	}
	return jkuo
}

// SetType sets the "type" field.
func (jkuo *JWTKeyUpdateOne) SetType(s string) *JWTKeyUpdateOne {
	jkuo.mutation.SetType(s)
	return jkuo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (jkuo *JWTKeyUpdateOne) SetNillableType(s *string) *JWTKeyUpdateOne {
	if s != nil {
		jkuo.SetType(*s)
	}
	return jkuo
}

// Mutation returns the JWTKeyMutation object of the builder.
func (jkuo *JWTKeyUpdateOne) Mutation() *JWTKeyMutation {
	return jkuo.mutation
}

// Where appends a list predicates to the JWTKeyUpdate builder.
func (jkuo *JWTKeyUpdateOne) Where(ps ...predicate.JWTKey) *JWTKeyUpdateOne {
	jkuo.mutation.Where(ps...)
	return jkuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (jkuo *JWTKeyUpdateOne) Select(field string, fields ...string) *JWTKeyUpdateOne {
	jkuo.fields = append([]string{field}, fields...)
	return jkuo
}

// Save executes the query and returns the updated JWTKey entity.
func (jkuo *JWTKeyUpdateOne) Save(ctx context.Context) (*JWTKey, error) {
	return withHooks(ctx, jkuo.sqlSave, jkuo.mutation, jkuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (jkuo *JWTKeyUpdateOne) SaveX(ctx context.Context) *JWTKey {
	node, err := jkuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (jkuo *JWTKeyUpdateOne) Exec(ctx context.Context) error {
	_, err := jkuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (jkuo *JWTKeyUpdateOne) ExecX(ctx context.Context) {
	if err := jkuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (jkuo *JWTKeyUpdateOne) check() error {
	if v, ok := jkuo.mutation.Key(); ok {
		if err := jwtkey.KeyValidator(v); err != nil {
			return &ValidationError{Name: "key", err: fmt.Errorf(`ent: validator failed for field "JWTKey.key": %w`, err)}
		}
	}
	if v, ok := jkuo.mutation.GetType(); ok {
		if err := jwtkey.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "JWTKey.type": %w`, err)}
		}
	}
	return nil
}

func (jkuo *JWTKeyUpdateOne) sqlSave(ctx context.Context) (_node *JWTKey, err error) {
	if err := jkuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(jwtkey.Table, jwtkey.Columns, sqlgraph.NewFieldSpec(jwtkey.FieldID, field.TypeInt))
	id, ok := jkuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "JWTKey.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := jkuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, jwtkey.FieldID)
		for _, f := range fields {
			if !jwtkey.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != jwtkey.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := jkuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := jkuo.mutation.Key(); ok {
		_spec.SetField(jwtkey.FieldKey, field.TypeString, value)
	}
	if value, ok := jkuo.mutation.GetType(); ok {
		_spec.SetField(jwtkey.FieldType, field.TypeString, value)
	}
	_node = &JWTKey{config: jkuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, jkuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{jwtkey.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	jkuo.mutation.done = true
	return _node, nil
}
