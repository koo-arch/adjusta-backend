// Code generated by ent, DO NOT EDIT.

package runtime

import (
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/jwtkey"
	"github.com/koo-arch/adjusta-backend/ent/oauthtoken"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
	"github.com/koo-arch/adjusta-backend/ent/schema"
	"github.com/koo-arch/adjusta-backend/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	calendarFields := schema.Calendar{}.Fields()
	_ = calendarFields
	// calendarDescIsPrimary is the schema descriptor for is_primary field.
	calendarDescIsPrimary := calendarFields[3].Descriptor()
	// calendar.DefaultIsPrimary holds the default value on creation for the is_primary field.
	calendar.DefaultIsPrimary = calendarDescIsPrimary.Default.(bool)
	// calendarDescID is the schema descriptor for id field.
	calendarDescID := calendarFields[0].Descriptor()
	// calendar.DefaultID holds the default value on creation for the id field.
	calendar.DefaultID = calendarDescID.Default.(func() uuid.UUID)
	eventFields := schema.Event{}.Fields()
	_ = eventFields
	// eventDescID is the schema descriptor for id field.
	eventDescID := eventFields[0].Descriptor()
	// event.DefaultID holds the default value on creation for the id field.
	event.DefaultID = eventDescID.Default.(func() uuid.UUID)
	jwtkeyFields := schema.JWTKey{}.Fields()
	_ = jwtkeyFields
	// jwtkeyDescKey is the schema descriptor for key field.
	jwtkeyDescKey := jwtkeyFields[0].Descriptor()
	// jwtkey.KeyValidator is a validator for the "key" field. It is called by the builders before save.
	jwtkey.KeyValidator = jwtkeyDescKey.Validators[0].(func(string) error)
	// jwtkeyDescType is the schema descriptor for type field.
	jwtkeyDescType := jwtkeyFields[1].Descriptor()
	// jwtkey.DefaultType holds the default value on creation for the type field.
	jwtkey.DefaultType = jwtkeyDescType.Default.(string)
	// jwtkey.TypeValidator is a validator for the "type" field. It is called by the builders before save.
	jwtkey.TypeValidator = jwtkeyDescType.Validators[0].(func(string) error)
	oauthtokenFields := schema.OAuthToken{}.Fields()
	_ = oauthtokenFields
	// oauthtokenDescID is the schema descriptor for id field.
	oauthtokenDescID := oauthtokenFields[0].Descriptor()
	// oauthtoken.DefaultID holds the default value on creation for the id field.
	oauthtoken.DefaultID = oauthtokenDescID.Default.(func() uuid.UUID)
	proposeddateHooks := schema.ProposedDate{}.Hooks()
	proposeddate.Hooks[0] = proposeddateHooks[0]
	proposeddateFields := schema.ProposedDate{}.Fields()
	_ = proposeddateFields
	// proposeddateDescPriority is the schema descriptor for priority field.
	proposeddateDescPriority := proposeddateFields[4].Descriptor()
	// proposeddate.DefaultPriority holds the default value on creation for the priority field.
	proposeddate.DefaultPriority = proposeddateDescPriority.Default.(int)
	// proposeddateDescID is the schema descriptor for id field.
	proposeddateDescID := proposeddateFields[0].Descriptor()
	// proposeddate.DefaultID holds the default value on creation for the id field.
	proposeddate.DefaultID = proposeddateDescID.Default.(func() uuid.UUID)
	userHooks := schema.User{}.Hooks()
	user.Hooks[0] = userHooks[0]
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[1].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = userDescEmail.Validators[0].(func(string) error)
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
}

const (
	Version = "v0.13.1"                                         // Version of ent codegen.
	Sum     = "h1:uD8QwN1h6SNphdCCzmkMN3feSUzNnVvV/WIkHKMbzOE=" // Sum of ent codegen.
)
