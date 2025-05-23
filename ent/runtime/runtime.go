// Code generated by ent, DO NOT EDIT.

package runtime

import (
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/googlecalendarinfo"
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
	calendarMixin := schema.Calendar{}.Mixin()
	calendarMixinInters1 := calendarMixin[1].Interceptors()
	calendar.Interceptors[0] = calendarMixinInters1[0]
	calendarMixinFields0 := calendarMixin[0].Fields()
	_ = calendarMixinFields0
	calendarFields := schema.Calendar{}.Fields()
	_ = calendarFields
	// calendarDescCreatedAt is the schema descriptor for created_at field.
	calendarDescCreatedAt := calendarMixinFields0[0].Descriptor()
	// calendar.DefaultCreatedAt holds the default value on creation for the created_at field.
	calendar.DefaultCreatedAt = calendarDescCreatedAt.Default.(func() time.Time)
	// calendarDescUpdatedAt is the schema descriptor for updated_at field.
	calendarDescUpdatedAt := calendarMixinFields0[1].Descriptor()
	// calendar.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	calendar.DefaultUpdatedAt = calendarDescUpdatedAt.Default.(func() time.Time)
	// calendar.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	calendar.UpdateDefaultUpdatedAt = calendarDescUpdatedAt.UpdateDefault.(func() time.Time)
	// calendarDescID is the schema descriptor for id field.
	calendarDescID := calendarFields[0].Descriptor()
	// calendar.DefaultID holds the default value on creation for the id field.
	calendar.DefaultID = calendarDescID.Default.(func() uuid.UUID)
	eventMixin := schema.Event{}.Mixin()
	eventHooks := schema.Event{}.Hooks()
	event.Hooks[0] = eventHooks[0]
	eventMixinInters1 := eventMixin[1].Interceptors()
	event.Interceptors[0] = eventMixinInters1[0]
	eventMixinFields0 := eventMixin[0].Fields()
	_ = eventMixinFields0
	eventFields := schema.Event{}.Fields()
	_ = eventFields
	// eventDescCreatedAt is the schema descriptor for created_at field.
	eventDescCreatedAt := eventMixinFields0[0].Descriptor()
	// event.DefaultCreatedAt holds the default value on creation for the created_at field.
	event.DefaultCreatedAt = eventDescCreatedAt.Default.(func() time.Time)
	// eventDescUpdatedAt is the schema descriptor for updated_at field.
	eventDescUpdatedAt := eventMixinFields0[1].Descriptor()
	// event.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	event.DefaultUpdatedAt = eventDescUpdatedAt.Default.(func() time.Time)
	// event.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	event.UpdateDefaultUpdatedAt = eventDescUpdatedAt.UpdateDefault.(func() time.Time)
	// eventDescID is the schema descriptor for id field.
	eventDescID := eventFields[0].Descriptor()
	// event.DefaultID holds the default value on creation for the id field.
	event.DefaultID = eventDescID.Default.(func() uuid.UUID)
	googlecalendarinfoMixin := schema.GoogleCalendarInfo{}.Mixin()
	googlecalendarinfoMixinInters1 := googlecalendarinfoMixin[1].Interceptors()
	googlecalendarinfo.Interceptors[0] = googlecalendarinfoMixinInters1[0]
	googlecalendarinfoMixinFields0 := googlecalendarinfoMixin[0].Fields()
	_ = googlecalendarinfoMixinFields0
	googlecalendarinfoFields := schema.GoogleCalendarInfo{}.Fields()
	_ = googlecalendarinfoFields
	// googlecalendarinfoDescCreatedAt is the schema descriptor for created_at field.
	googlecalendarinfoDescCreatedAt := googlecalendarinfoMixinFields0[0].Descriptor()
	// googlecalendarinfo.DefaultCreatedAt holds the default value on creation for the created_at field.
	googlecalendarinfo.DefaultCreatedAt = googlecalendarinfoDescCreatedAt.Default.(func() time.Time)
	// googlecalendarinfoDescUpdatedAt is the schema descriptor for updated_at field.
	googlecalendarinfoDescUpdatedAt := googlecalendarinfoMixinFields0[1].Descriptor()
	// googlecalendarinfo.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	googlecalendarinfo.DefaultUpdatedAt = googlecalendarinfoDescUpdatedAt.Default.(func() time.Time)
	// googlecalendarinfo.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	googlecalendarinfo.UpdateDefaultUpdatedAt = googlecalendarinfoDescUpdatedAt.UpdateDefault.(func() time.Time)
	// googlecalendarinfoDescGoogleCalendarID is the schema descriptor for google_calendar_id field.
	googlecalendarinfoDescGoogleCalendarID := googlecalendarinfoFields[1].Descriptor()
	// googlecalendarinfo.GoogleCalendarIDValidator is a validator for the "google_calendar_id" field. It is called by the builders before save.
	googlecalendarinfo.GoogleCalendarIDValidator = googlecalendarinfoDescGoogleCalendarID.Validators[0].(func(string) error)
	// googlecalendarinfoDescIsPrimary is the schema descriptor for is_primary field.
	googlecalendarinfoDescIsPrimary := googlecalendarinfoFields[3].Descriptor()
	// googlecalendarinfo.DefaultIsPrimary holds the default value on creation for the is_primary field.
	googlecalendarinfo.DefaultIsPrimary = googlecalendarinfoDescIsPrimary.Default.(bool)
	// googlecalendarinfoDescID is the schema descriptor for id field.
	googlecalendarinfoDescID := googlecalendarinfoFields[0].Descriptor()
	// googlecalendarinfo.DefaultID holds the default value on creation for the id field.
	googlecalendarinfo.DefaultID = googlecalendarinfoDescID.Default.(func() uuid.UUID)
	jwtkeyMixin := schema.JWTKey{}.Mixin()
	jwtkeyMixinInters1 := jwtkeyMixin[1].Interceptors()
	jwtkey.Interceptors[0] = jwtkeyMixinInters1[0]
	jwtkeyMixinFields0 := jwtkeyMixin[0].Fields()
	_ = jwtkeyMixinFields0
	jwtkeyFields := schema.JWTKey{}.Fields()
	_ = jwtkeyFields
	// jwtkeyDescCreatedAt is the schema descriptor for created_at field.
	jwtkeyDescCreatedAt := jwtkeyMixinFields0[0].Descriptor()
	// jwtkey.DefaultCreatedAt holds the default value on creation for the created_at field.
	jwtkey.DefaultCreatedAt = jwtkeyDescCreatedAt.Default.(func() time.Time)
	// jwtkeyDescUpdatedAt is the schema descriptor for updated_at field.
	jwtkeyDescUpdatedAt := jwtkeyMixinFields0[1].Descriptor()
	// jwtkey.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	jwtkey.DefaultUpdatedAt = jwtkeyDescUpdatedAt.Default.(func() time.Time)
	// jwtkey.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	jwtkey.UpdateDefaultUpdatedAt = jwtkeyDescUpdatedAt.UpdateDefault.(func() time.Time)
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
	oauthtokenMixin := schema.OAuthToken{}.Mixin()
	oauthtokenMixinInters1 := oauthtokenMixin[1].Interceptors()
	oauthtoken.Interceptors[0] = oauthtokenMixinInters1[0]
	oauthtokenMixinFields0 := oauthtokenMixin[0].Fields()
	_ = oauthtokenMixinFields0
	oauthtokenFields := schema.OAuthToken{}.Fields()
	_ = oauthtokenFields
	// oauthtokenDescCreatedAt is the schema descriptor for created_at field.
	oauthtokenDescCreatedAt := oauthtokenMixinFields0[0].Descriptor()
	// oauthtoken.DefaultCreatedAt holds the default value on creation for the created_at field.
	oauthtoken.DefaultCreatedAt = oauthtokenDescCreatedAt.Default.(func() time.Time)
	// oauthtokenDescUpdatedAt is the schema descriptor for updated_at field.
	oauthtokenDescUpdatedAt := oauthtokenMixinFields0[1].Descriptor()
	// oauthtoken.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	oauthtoken.DefaultUpdatedAt = oauthtokenDescUpdatedAt.Default.(func() time.Time)
	// oauthtoken.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	oauthtoken.UpdateDefaultUpdatedAt = oauthtokenDescUpdatedAt.UpdateDefault.(func() time.Time)
	// oauthtokenDescID is the schema descriptor for id field.
	oauthtokenDescID := oauthtokenFields[0].Descriptor()
	// oauthtoken.DefaultID holds the default value on creation for the id field.
	oauthtoken.DefaultID = oauthtokenDescID.Default.(func() uuid.UUID)
	proposeddateMixin := schema.ProposedDate{}.Mixin()
	proposeddateHooks := schema.ProposedDate{}.Hooks()
	proposeddate.Hooks[0] = proposeddateHooks[0]
	proposeddateMixinInters1 := proposeddateMixin[1].Interceptors()
	proposeddate.Interceptors[0] = proposeddateMixinInters1[0]
	proposeddateMixinFields0 := proposeddateMixin[0].Fields()
	_ = proposeddateMixinFields0
	proposeddateFields := schema.ProposedDate{}.Fields()
	_ = proposeddateFields
	// proposeddateDescCreatedAt is the schema descriptor for created_at field.
	proposeddateDescCreatedAt := proposeddateMixinFields0[0].Descriptor()
	// proposeddate.DefaultCreatedAt holds the default value on creation for the created_at field.
	proposeddate.DefaultCreatedAt = proposeddateDescCreatedAt.Default.(func() time.Time)
	// proposeddateDescUpdatedAt is the schema descriptor for updated_at field.
	proposeddateDescUpdatedAt := proposeddateMixinFields0[1].Descriptor()
	// proposeddate.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	proposeddate.DefaultUpdatedAt = proposeddateDescUpdatedAt.Default.(func() time.Time)
	// proposeddate.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	proposeddate.UpdateDefaultUpdatedAt = proposeddateDescUpdatedAt.UpdateDefault.(func() time.Time)
	// proposeddateDescPriority is the schema descriptor for priority field.
	proposeddateDescPriority := proposeddateFields[3].Descriptor()
	// proposeddate.DefaultPriority holds the default value on creation for the priority field.
	proposeddate.DefaultPriority = proposeddateDescPriority.Default.(int)
	// proposeddateDescID is the schema descriptor for id field.
	proposeddateDescID := proposeddateFields[0].Descriptor()
	// proposeddate.DefaultID holds the default value on creation for the id field.
	proposeddate.DefaultID = proposeddateDescID.Default.(func() uuid.UUID)
	userMixin := schema.User{}.Mixin()
	userHooks := schema.User{}.Hooks()
	user.Hooks[0] = userHooks[0]
	userMixinInters1 := userMixin[1].Interceptors()
	user.Interceptors[0] = userMixinInters1[0]
	userMixinFields0 := userMixin[0].Fields()
	_ = userMixinFields0
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userMixinFields0[0].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userMixinFields0[1].Descriptor()
	// user.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	user.DefaultUpdatedAt = userDescUpdatedAt.Default.(func() time.Time)
	// user.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	user.UpdateDefaultUpdatedAt = userDescUpdatedAt.UpdateDefault.(func() time.Time)
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
	Version = "v0.14.4"                                         // Version of ent codegen.
	Sum     = "h1:/DhDraSLXIkBhyiVoJeSshr4ZYi7femzhj6/TckzZuI=" // Sum of ent codegen.
)
