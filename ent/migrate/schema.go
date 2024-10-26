// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AccountsColumns holds the columns for the "accounts" table.
	AccountsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "email", Type: field.TypeString},
		{Name: "google_id", Type: field.TypeString},
		{Name: "access_token", Type: field.TypeString, Nullable: true},
		{Name: "refresh_token", Type: field.TypeString, Nullable: true},
		{Name: "access_token_expiry", Type: field.TypeTime, Nullable: true},
		{Name: "user_accounts", Type: field.TypeUUID, Nullable: true},
	}
	// AccountsTable holds the schema information for the "accounts" table.
	AccountsTable = &schema.Table{
		Name:       "accounts",
		Columns:    AccountsColumns,
		PrimaryKey: []*schema.Column{AccountsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "accounts_users_accounts",
				Columns:    []*schema.Column{AccountsColumns[6]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// CalendarsColumns holds the columns for the "calendars" table.
	CalendarsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "calendar_id", Type: field.TypeString},
		{Name: "summary", Type: field.TypeString},
		{Name: "is_primary", Type: field.TypeBool, Default: false},
		{Name: "account_calendars", Type: field.TypeUUID, Nullable: true},
	}
	// CalendarsTable holds the schema information for the "calendars" table.
	CalendarsTable = &schema.Table{
		Name:       "calendars",
		Columns:    CalendarsColumns,
		PrimaryKey: []*schema.Column{CalendarsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "calendars_accounts_calendars",
				Columns:    []*schema.Column{CalendarsColumns[4]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// EventsColumns holds the columns for the "events" table.
	EventsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "summary", Type: field.TypeString, Nullable: true},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "location", Type: field.TypeString, Nullable: true},
		{Name: "calendar_events", Type: field.TypeUUID, Nullable: true},
	}
	// EventsTable holds the schema information for the "events" table.
	EventsTable = &schema.Table{
		Name:       "events",
		Columns:    EventsColumns,
		PrimaryKey: []*schema.Column{EventsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "events_calendars_events",
				Columns:    []*schema.Column{EventsColumns[4]},
				RefColumns: []*schema.Column{CalendarsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// JwtKeysColumns holds the columns for the "jwt_keys" table.
	JwtKeysColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "key", Type: field.TypeString},
		{Name: "type", Type: field.TypeString, Default: "access"},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "expires_at", Type: field.TypeTime},
	}
	// JwtKeysTable holds the schema information for the "jwt_keys" table.
	JwtKeysTable = &schema.Table{
		Name:       "jwt_keys",
		Columns:    JwtKeysColumns,
		PrimaryKey: []*schema.Column{JwtKeysColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "idx_type_expires",
				Unique:  false,
				Columns: []*schema.Column{JwtKeysColumns[1], JwtKeysColumns[4]},
			},
		},
	}
	// ProposedDatesColumns holds the columns for the "proposed_dates" table.
	ProposedDatesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "google_event_id", Type: field.TypeString, Nullable: true},
		{Name: "start_time", Type: field.TypeTime},
		{Name: "end_time", Type: field.TypeTime},
		{Name: "is_finalized", Type: field.TypeBool, Default: false},
		{Name: "priority", Type: field.TypeInt, Default: 0},
		{Name: "event_proposed_dates", Type: field.TypeUUID, Nullable: true},
	}
	// ProposedDatesTable holds the schema information for the "proposed_dates" table.
	ProposedDatesTable = &schema.Table{
		Name:       "proposed_dates",
		Columns:    ProposedDatesColumns,
		PrimaryKey: []*schema.Column{ProposedDatesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "proposed_dates_events_proposed_dates",
				Columns:    []*schema.Column{ProposedDatesColumns[6]},
				RefColumns: []*schema.Column{EventsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "email", Type: field.TypeString, Unique: true},
		{Name: "refresh_token", Type: field.TypeString, Nullable: true},
		{Name: "refresh_token_expiry", Type: field.TypeTime, Nullable: true},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AccountsTable,
		CalendarsTable,
		EventsTable,
		JwtKeysTable,
		ProposedDatesTable,
		UsersTable,
	}
)

func init() {
	AccountsTable.ForeignKeys[0].RefTable = UsersTable
	CalendarsTable.ForeignKeys[0].RefTable = AccountsTable
	EventsTable.ForeignKeys[0].RefTable = CalendarsTable
	ProposedDatesTable.ForeignKeys[0].RefTable = EventsTable
}
