// Code generated by ent, DO NOT EDIT.

package account

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the account type in the database.
	Label = "account"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldGoogleID holds the string denoting the google_id field in the database.
	FieldGoogleID = "google_id"
	// FieldAccessToken holds the string denoting the access_token field in the database.
	FieldAccessToken = "access_token"
	// FieldRefreshToken holds the string denoting the refresh_token field in the database.
	FieldRefreshToken = "refresh_token"
	// FieldAccessTokenExpiry holds the string denoting the access_token_expiry field in the database.
	FieldAccessTokenExpiry = "access_token_expiry"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// Table holds the table name of the account in the database.
	Table = "accounts"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "accounts"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_accounts"
)

// Columns holds all SQL columns for account fields.
var Columns = []string{
	FieldID,
	FieldEmail,
	FieldGoogleID,
	FieldAccessToken,
	FieldRefreshToken,
	FieldAccessTokenExpiry,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "accounts"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_accounts",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/koo-arch/adjusta-backend/ent/runtime"
var (
	Hooks [2]ent.Hook
	// EmailValidator is a validator for the "email" field. It is called by the builders before save.
	EmailValidator func(string) error
	// GoogleIDValidator is a validator for the "google_id" field. It is called by the builders before save.
	GoogleIDValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Account queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByEmail orders the results by the email field.
func ByEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmail, opts...).ToFunc()
}

// ByGoogleID orders the results by the google_id field.
func ByGoogleID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGoogleID, opts...).ToFunc()
}

// ByAccessToken orders the results by the access_token field.
func ByAccessToken(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAccessToken, opts...).ToFunc()
}

// ByRefreshToken orders the results by the refresh_token field.
func ByRefreshToken(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRefreshToken, opts...).ToFunc()
}

// ByAccessTokenExpiry orders the results by the access_token_expiry field.
func ByAccessTokenExpiry(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAccessTokenExpiry, opts...).ToFunc()
}

// ByUserField orders the results by user field.
func ByUserField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUserStep(), sql.OrderByField(field, opts...))
	}
}
func newUserStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UserInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
	)
}
