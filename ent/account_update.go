// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/account"
	"github.com/koo-arch/adjusta-backend/ent/predicate"
	"github.com/koo-arch/adjusta-backend/ent/user"
)

// AccountUpdate is the builder for updating Account entities.
type AccountUpdate struct {
	config
	hooks    []Hook
	mutation *AccountMutation
}

// Where appends a list predicates to the AccountUpdate builder.
func (au *AccountUpdate) Where(ps ...predicate.Account) *AccountUpdate {
	au.mutation.Where(ps...)
	return au
}

// SetEmail sets the "email" field.
func (au *AccountUpdate) SetEmail(s string) *AccountUpdate {
	au.mutation.SetEmail(s)
	return au
}

// SetNillableEmail sets the "email" field if the given value is not nil.
func (au *AccountUpdate) SetNillableEmail(s *string) *AccountUpdate {
	if s != nil {
		au.SetEmail(*s)
	}
	return au
}

// SetGoogleID sets the "google_id" field.
func (au *AccountUpdate) SetGoogleID(s string) *AccountUpdate {
	au.mutation.SetGoogleID(s)
	return au
}

// SetNillableGoogleID sets the "google_id" field if the given value is not nil.
func (au *AccountUpdate) SetNillableGoogleID(s *string) *AccountUpdate {
	if s != nil {
		au.SetGoogleID(*s)
	}
	return au
}

// SetAccessToken sets the "access_token" field.
func (au *AccountUpdate) SetAccessToken(s string) *AccountUpdate {
	au.mutation.SetAccessToken(s)
	return au
}

// SetNillableAccessToken sets the "access_token" field if the given value is not nil.
func (au *AccountUpdate) SetNillableAccessToken(s *string) *AccountUpdate {
	if s != nil {
		au.SetAccessToken(*s)
	}
	return au
}

// ClearAccessToken clears the value of the "access_token" field.
func (au *AccountUpdate) ClearAccessToken() *AccountUpdate {
	au.mutation.ClearAccessToken()
	return au
}

// SetRefreshToken sets the "refresh_token" field.
func (au *AccountUpdate) SetRefreshToken(s string) *AccountUpdate {
	au.mutation.SetRefreshToken(s)
	return au
}

// SetNillableRefreshToken sets the "refresh_token" field if the given value is not nil.
func (au *AccountUpdate) SetNillableRefreshToken(s *string) *AccountUpdate {
	if s != nil {
		au.SetRefreshToken(*s)
	}
	return au
}

// ClearRefreshToken clears the value of the "refresh_token" field.
func (au *AccountUpdate) ClearRefreshToken() *AccountUpdate {
	au.mutation.ClearRefreshToken()
	return au
}

// SetAccessTokenExpiry sets the "access_token_expiry" field.
func (au *AccountUpdate) SetAccessTokenExpiry(t time.Time) *AccountUpdate {
	au.mutation.SetAccessTokenExpiry(t)
	return au
}

// SetNillableAccessTokenExpiry sets the "access_token_expiry" field if the given value is not nil.
func (au *AccountUpdate) SetNillableAccessTokenExpiry(t *time.Time) *AccountUpdate {
	if t != nil {
		au.SetAccessTokenExpiry(*t)
	}
	return au
}

// ClearAccessTokenExpiry clears the value of the "access_token_expiry" field.
func (au *AccountUpdate) ClearAccessTokenExpiry() *AccountUpdate {
	au.mutation.ClearAccessTokenExpiry()
	return au
}

// SetUserID sets the "user" edge to the User entity by ID.
func (au *AccountUpdate) SetUserID(id uuid.UUID) *AccountUpdate {
	au.mutation.SetUserID(id)
	return au
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (au *AccountUpdate) SetNillableUserID(id *uuid.UUID) *AccountUpdate {
	if id != nil {
		au = au.SetUserID(*id)
	}
	return au
}

// SetUser sets the "user" edge to the User entity.
func (au *AccountUpdate) SetUser(u *User) *AccountUpdate {
	return au.SetUserID(u.ID)
}

// Mutation returns the AccountMutation object of the builder.
func (au *AccountUpdate) Mutation() *AccountMutation {
	return au.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (au *AccountUpdate) ClearUser() *AccountUpdate {
	au.mutation.ClearUser()
	return au
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (au *AccountUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, au.sqlSave, au.mutation, au.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (au *AccountUpdate) SaveX(ctx context.Context) int {
	affected, err := au.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (au *AccountUpdate) Exec(ctx context.Context) error {
	_, err := au.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (au *AccountUpdate) ExecX(ctx context.Context) {
	if err := au.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (au *AccountUpdate) check() error {
	if v, ok := au.mutation.Email(); ok {
		if err := account.EmailValidator(v); err != nil {
			return &ValidationError{Name: "email", err: fmt.Errorf(`ent: validator failed for field "Account.email": %w`, err)}
		}
	}
	if v, ok := au.mutation.GoogleID(); ok {
		if err := account.GoogleIDValidator(v); err != nil {
			return &ValidationError{Name: "google_id", err: fmt.Errorf(`ent: validator failed for field "Account.google_id": %w`, err)}
		}
	}
	return nil
}

func (au *AccountUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := au.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(account.Table, account.Columns, sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID))
	if ps := au.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := au.mutation.Email(); ok {
		_spec.SetField(account.FieldEmail, field.TypeString, value)
	}
	if value, ok := au.mutation.GoogleID(); ok {
		_spec.SetField(account.FieldGoogleID, field.TypeString, value)
	}
	if value, ok := au.mutation.AccessToken(); ok {
		_spec.SetField(account.FieldAccessToken, field.TypeString, value)
	}
	if au.mutation.AccessTokenCleared() {
		_spec.ClearField(account.FieldAccessToken, field.TypeString)
	}
	if value, ok := au.mutation.RefreshToken(); ok {
		_spec.SetField(account.FieldRefreshToken, field.TypeString, value)
	}
	if au.mutation.RefreshTokenCleared() {
		_spec.ClearField(account.FieldRefreshToken, field.TypeString)
	}
	if value, ok := au.mutation.AccessTokenExpiry(); ok {
		_spec.SetField(account.FieldAccessTokenExpiry, field.TypeTime, value)
	}
	if au.mutation.AccessTokenExpiryCleared() {
		_spec.ClearField(account.FieldAccessTokenExpiry, field.TypeTime)
	}
	if au.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   account.UserTable,
			Columns: []string{account.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := au.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   account.UserTable,
			Columns: []string{account.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, au.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{account.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	au.mutation.done = true
	return n, nil
}

// AccountUpdateOne is the builder for updating a single Account entity.
type AccountUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *AccountMutation
}

// SetEmail sets the "email" field.
func (auo *AccountUpdateOne) SetEmail(s string) *AccountUpdateOne {
	auo.mutation.SetEmail(s)
	return auo
}

// SetNillableEmail sets the "email" field if the given value is not nil.
func (auo *AccountUpdateOne) SetNillableEmail(s *string) *AccountUpdateOne {
	if s != nil {
		auo.SetEmail(*s)
	}
	return auo
}

// SetGoogleID sets the "google_id" field.
func (auo *AccountUpdateOne) SetGoogleID(s string) *AccountUpdateOne {
	auo.mutation.SetGoogleID(s)
	return auo
}

// SetNillableGoogleID sets the "google_id" field if the given value is not nil.
func (auo *AccountUpdateOne) SetNillableGoogleID(s *string) *AccountUpdateOne {
	if s != nil {
		auo.SetGoogleID(*s)
	}
	return auo
}

// SetAccessToken sets the "access_token" field.
func (auo *AccountUpdateOne) SetAccessToken(s string) *AccountUpdateOne {
	auo.mutation.SetAccessToken(s)
	return auo
}

// SetNillableAccessToken sets the "access_token" field if the given value is not nil.
func (auo *AccountUpdateOne) SetNillableAccessToken(s *string) *AccountUpdateOne {
	if s != nil {
		auo.SetAccessToken(*s)
	}
	return auo
}

// ClearAccessToken clears the value of the "access_token" field.
func (auo *AccountUpdateOne) ClearAccessToken() *AccountUpdateOne {
	auo.mutation.ClearAccessToken()
	return auo
}

// SetRefreshToken sets the "refresh_token" field.
func (auo *AccountUpdateOne) SetRefreshToken(s string) *AccountUpdateOne {
	auo.mutation.SetRefreshToken(s)
	return auo
}

// SetNillableRefreshToken sets the "refresh_token" field if the given value is not nil.
func (auo *AccountUpdateOne) SetNillableRefreshToken(s *string) *AccountUpdateOne {
	if s != nil {
		auo.SetRefreshToken(*s)
	}
	return auo
}

// ClearRefreshToken clears the value of the "refresh_token" field.
func (auo *AccountUpdateOne) ClearRefreshToken() *AccountUpdateOne {
	auo.mutation.ClearRefreshToken()
	return auo
}

// SetAccessTokenExpiry sets the "access_token_expiry" field.
func (auo *AccountUpdateOne) SetAccessTokenExpiry(t time.Time) *AccountUpdateOne {
	auo.mutation.SetAccessTokenExpiry(t)
	return auo
}

// SetNillableAccessTokenExpiry sets the "access_token_expiry" field if the given value is not nil.
func (auo *AccountUpdateOne) SetNillableAccessTokenExpiry(t *time.Time) *AccountUpdateOne {
	if t != nil {
		auo.SetAccessTokenExpiry(*t)
	}
	return auo
}

// ClearAccessTokenExpiry clears the value of the "access_token_expiry" field.
func (auo *AccountUpdateOne) ClearAccessTokenExpiry() *AccountUpdateOne {
	auo.mutation.ClearAccessTokenExpiry()
	return auo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (auo *AccountUpdateOne) SetUserID(id uuid.UUID) *AccountUpdateOne {
	auo.mutation.SetUserID(id)
	return auo
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (auo *AccountUpdateOne) SetNillableUserID(id *uuid.UUID) *AccountUpdateOne {
	if id != nil {
		auo = auo.SetUserID(*id)
	}
	return auo
}

// SetUser sets the "user" edge to the User entity.
func (auo *AccountUpdateOne) SetUser(u *User) *AccountUpdateOne {
	return auo.SetUserID(u.ID)
}

// Mutation returns the AccountMutation object of the builder.
func (auo *AccountUpdateOne) Mutation() *AccountMutation {
	return auo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (auo *AccountUpdateOne) ClearUser() *AccountUpdateOne {
	auo.mutation.ClearUser()
	return auo
}

// Where appends a list predicates to the AccountUpdate builder.
func (auo *AccountUpdateOne) Where(ps ...predicate.Account) *AccountUpdateOne {
	auo.mutation.Where(ps...)
	return auo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (auo *AccountUpdateOne) Select(field string, fields ...string) *AccountUpdateOne {
	auo.fields = append([]string{field}, fields...)
	return auo
}

// Save executes the query and returns the updated Account entity.
func (auo *AccountUpdateOne) Save(ctx context.Context) (*Account, error) {
	return withHooks(ctx, auo.sqlSave, auo.mutation, auo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (auo *AccountUpdateOne) SaveX(ctx context.Context) *Account {
	node, err := auo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (auo *AccountUpdateOne) Exec(ctx context.Context) error {
	_, err := auo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (auo *AccountUpdateOne) ExecX(ctx context.Context) {
	if err := auo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (auo *AccountUpdateOne) check() error {
	if v, ok := auo.mutation.Email(); ok {
		if err := account.EmailValidator(v); err != nil {
			return &ValidationError{Name: "email", err: fmt.Errorf(`ent: validator failed for field "Account.email": %w`, err)}
		}
	}
	if v, ok := auo.mutation.GoogleID(); ok {
		if err := account.GoogleIDValidator(v); err != nil {
			return &ValidationError{Name: "google_id", err: fmt.Errorf(`ent: validator failed for field "Account.google_id": %w`, err)}
		}
	}
	return nil
}

func (auo *AccountUpdateOne) sqlSave(ctx context.Context) (_node *Account, err error) {
	if err := auo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(account.Table, account.Columns, sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID))
	id, ok := auo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Account.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := auo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, account.FieldID)
		for _, f := range fields {
			if !account.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != account.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := auo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := auo.mutation.Email(); ok {
		_spec.SetField(account.FieldEmail, field.TypeString, value)
	}
	if value, ok := auo.mutation.GoogleID(); ok {
		_spec.SetField(account.FieldGoogleID, field.TypeString, value)
	}
	if value, ok := auo.mutation.AccessToken(); ok {
		_spec.SetField(account.FieldAccessToken, field.TypeString, value)
	}
	if auo.mutation.AccessTokenCleared() {
		_spec.ClearField(account.FieldAccessToken, field.TypeString)
	}
	if value, ok := auo.mutation.RefreshToken(); ok {
		_spec.SetField(account.FieldRefreshToken, field.TypeString, value)
	}
	if auo.mutation.RefreshTokenCleared() {
		_spec.ClearField(account.FieldRefreshToken, field.TypeString)
	}
	if value, ok := auo.mutation.AccessTokenExpiry(); ok {
		_spec.SetField(account.FieldAccessTokenExpiry, field.TypeTime, value)
	}
	if auo.mutation.AccessTokenExpiryCleared() {
		_spec.ClearField(account.FieldAccessTokenExpiry, field.TypeTime)
	}
	if auo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   account.UserTable,
			Columns: []string{account.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := auo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   account.UserTable,
			Columns: []string{account.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Account{config: auo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, auo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{account.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	auo.mutation.done = true
	return _node, nil
}