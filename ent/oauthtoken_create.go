// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/oauthtoken"
	"github.com/koo-arch/adjusta-backend/ent/user"
)

// OAuthTokenCreate is the builder for creating a OAuthToken entity.
type OAuthTokenCreate struct {
	config
	mutation *OAuthTokenMutation
	hooks    []Hook
}

// SetAccessToken sets the "access_token" field.
func (otc *OAuthTokenCreate) SetAccessToken(s string) *OAuthTokenCreate {
	otc.mutation.SetAccessToken(s)
	return otc
}

// SetNillableAccessToken sets the "access_token" field if the given value is not nil.
func (otc *OAuthTokenCreate) SetNillableAccessToken(s *string) *OAuthTokenCreate {
	if s != nil {
		otc.SetAccessToken(*s)
	}
	return otc
}

// SetRefreshToken sets the "refresh_token" field.
func (otc *OAuthTokenCreate) SetRefreshToken(s string) *OAuthTokenCreate {
	otc.mutation.SetRefreshToken(s)
	return otc
}

// SetNillableRefreshToken sets the "refresh_token" field if the given value is not nil.
func (otc *OAuthTokenCreate) SetNillableRefreshToken(s *string) *OAuthTokenCreate {
	if s != nil {
		otc.SetRefreshToken(*s)
	}
	return otc
}

// SetExpiry sets the "expiry" field.
func (otc *OAuthTokenCreate) SetExpiry(t time.Time) *OAuthTokenCreate {
	otc.mutation.SetExpiry(t)
	return otc
}

// SetNillableExpiry sets the "expiry" field if the given value is not nil.
func (otc *OAuthTokenCreate) SetNillableExpiry(t *time.Time) *OAuthTokenCreate {
	if t != nil {
		otc.SetExpiry(*t)
	}
	return otc
}

// SetID sets the "id" field.
func (otc *OAuthTokenCreate) SetID(u uuid.UUID) *OAuthTokenCreate {
	otc.mutation.SetID(u)
	return otc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (otc *OAuthTokenCreate) SetNillableID(u *uuid.UUID) *OAuthTokenCreate {
	if u != nil {
		otc.SetID(*u)
	}
	return otc
}

// SetUserID sets the "user" edge to the User entity by ID.
func (otc *OAuthTokenCreate) SetUserID(id uuid.UUID) *OAuthTokenCreate {
	otc.mutation.SetUserID(id)
	return otc
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (otc *OAuthTokenCreate) SetNillableUserID(id *uuid.UUID) *OAuthTokenCreate {
	if id != nil {
		otc = otc.SetUserID(*id)
	}
	return otc
}

// SetUser sets the "user" edge to the User entity.
func (otc *OAuthTokenCreate) SetUser(u *User) *OAuthTokenCreate {
	return otc.SetUserID(u.ID)
}

// Mutation returns the OAuthTokenMutation object of the builder.
func (otc *OAuthTokenCreate) Mutation() *OAuthTokenMutation {
	return otc.mutation
}

// Save creates the OAuthToken in the database.
func (otc *OAuthTokenCreate) Save(ctx context.Context) (*OAuthToken, error) {
	otc.defaults()
	return withHooks(ctx, otc.sqlSave, otc.mutation, otc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (otc *OAuthTokenCreate) SaveX(ctx context.Context) *OAuthToken {
	v, err := otc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (otc *OAuthTokenCreate) Exec(ctx context.Context) error {
	_, err := otc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (otc *OAuthTokenCreate) ExecX(ctx context.Context) {
	if err := otc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (otc *OAuthTokenCreate) defaults() {
	if _, ok := otc.mutation.ID(); !ok {
		v := oauthtoken.DefaultID()
		otc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (otc *OAuthTokenCreate) check() error {
	return nil
}

func (otc *OAuthTokenCreate) sqlSave(ctx context.Context) (*OAuthToken, error) {
	if err := otc.check(); err != nil {
		return nil, err
	}
	_node, _spec := otc.createSpec()
	if err := sqlgraph.CreateNode(ctx, otc.driver, _spec); err != nil {
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
	otc.mutation.id = &_node.ID
	otc.mutation.done = true
	return _node, nil
}

func (otc *OAuthTokenCreate) createSpec() (*OAuthToken, *sqlgraph.CreateSpec) {
	var (
		_node = &OAuthToken{config: otc.config}
		_spec = sqlgraph.NewCreateSpec(oauthtoken.Table, sqlgraph.NewFieldSpec(oauthtoken.FieldID, field.TypeUUID))
	)
	if id, ok := otc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := otc.mutation.AccessToken(); ok {
		_spec.SetField(oauthtoken.FieldAccessToken, field.TypeString, value)
		_node.AccessToken = value
	}
	if value, ok := otc.mutation.RefreshToken(); ok {
		_spec.SetField(oauthtoken.FieldRefreshToken, field.TypeString, value)
		_node.RefreshToken = value
	}
	if value, ok := otc.mutation.Expiry(); ok {
		_spec.SetField(oauthtoken.FieldExpiry, field.TypeTime, value)
		_node.Expiry = value
	}
	if nodes := otc.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   oauthtoken.UserTable,
			Columns: []string{oauthtoken.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.user_oauth_token = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OAuthTokenCreateBulk is the builder for creating many OAuthToken entities in bulk.
type OAuthTokenCreateBulk struct {
	config
	err      error
	builders []*OAuthTokenCreate
}

// Save creates the OAuthToken entities in the database.
func (otcb *OAuthTokenCreateBulk) Save(ctx context.Context) ([]*OAuthToken, error) {
	if otcb.err != nil {
		return nil, otcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(otcb.builders))
	nodes := make([]*OAuthToken, len(otcb.builders))
	mutators := make([]Mutator, len(otcb.builders))
	for i := range otcb.builders {
		func(i int, root context.Context) {
			builder := otcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*OAuthTokenMutation)
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
					_, err = mutators[i+1].Mutate(root, otcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, otcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, otcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (otcb *OAuthTokenCreateBulk) SaveX(ctx context.Context) []*OAuthToken {
	v, err := otcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (otcb *OAuthTokenCreateBulk) Exec(ctx context.Context) error {
	_, err := otcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (otcb *OAuthTokenCreateBulk) ExecX(ctx context.Context) {
	if err := otcb.Exec(ctx); err != nil {
		panic(err)
	}
}
