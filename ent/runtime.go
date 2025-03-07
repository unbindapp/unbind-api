// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/githubapp"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/ent/oauth2code"
	"github.com/unbindapp/unbind-api/ent/oauth2token"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	githubappMixin := schema.GithubApp{}.Mixin()
	githubappMixinFields0 := githubappMixin[0].Fields()
	_ = githubappMixinFields0
	githubappFields := schema.GithubApp{}.Fields()
	_ = githubappFields
	// githubappDescCreatedAt is the schema descriptor for created_at field.
	githubappDescCreatedAt := githubappMixinFields0[0].Descriptor()
	// githubapp.DefaultCreatedAt holds the default value on creation for the created_at field.
	githubapp.DefaultCreatedAt = githubappDescCreatedAt.Default.(func() time.Time)
	// githubappDescUpdatedAt is the schema descriptor for updated_at field.
	githubappDescUpdatedAt := githubappMixinFields0[1].Descriptor()
	// githubapp.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	githubapp.DefaultUpdatedAt = githubappDescUpdatedAt.Default.(func() time.Time)
	// githubapp.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	githubapp.UpdateDefaultUpdatedAt = githubappDescUpdatedAt.UpdateDefault.(func() time.Time)
	// githubappDescName is the schema descriptor for name field.
	githubappDescName := githubappFields[1].Descriptor()
	// githubapp.NameValidator is a validator for the "name" field. It is called by the builders before save.
	githubapp.NameValidator = githubappDescName.Validators[0].(func(string) error)
	// githubappDescID is the schema descriptor for id field.
	githubappDescID := githubappFields[0].Descriptor()
	// githubapp.IDValidator is a validator for the "id" field. It is called by the builders before save.
	githubapp.IDValidator = githubappDescID.Validators[0].(func(int64) error)
	githubinstallationMixin := schema.GithubInstallation{}.Mixin()
	githubinstallationMixinFields0 := githubinstallationMixin[0].Fields()
	_ = githubinstallationMixinFields0
	githubinstallationFields := schema.GithubInstallation{}.Fields()
	_ = githubinstallationFields
	// githubinstallationDescCreatedAt is the schema descriptor for created_at field.
	githubinstallationDescCreatedAt := githubinstallationMixinFields0[0].Descriptor()
	// githubinstallation.DefaultCreatedAt holds the default value on creation for the created_at field.
	githubinstallation.DefaultCreatedAt = githubinstallationDescCreatedAt.Default.(func() time.Time)
	// githubinstallationDescUpdatedAt is the schema descriptor for updated_at field.
	githubinstallationDescUpdatedAt := githubinstallationMixinFields0[1].Descriptor()
	// githubinstallation.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	githubinstallation.DefaultUpdatedAt = githubinstallationDescUpdatedAt.Default.(func() time.Time)
	// githubinstallation.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	githubinstallation.UpdateDefaultUpdatedAt = githubinstallationDescUpdatedAt.UpdateDefault.(func() time.Time)
	// githubinstallationDescAccountLogin is the schema descriptor for account_login field.
	githubinstallationDescAccountLogin := githubinstallationFields[3].Descriptor()
	// githubinstallation.AccountLoginValidator is a validator for the "account_login" field. It is called by the builders before save.
	githubinstallation.AccountLoginValidator = githubinstallationDescAccountLogin.Validators[0].(func(string) error)
	// githubinstallationDescAccountURL is the schema descriptor for account_url field.
	githubinstallationDescAccountURL := githubinstallationFields[5].Descriptor()
	// githubinstallation.AccountURLValidator is a validator for the "account_url" field. It is called by the builders before save.
	githubinstallation.AccountURLValidator = githubinstallationDescAccountURL.Validators[0].(func(string) error)
	// githubinstallationDescSuspended is the schema descriptor for suspended field.
	githubinstallationDescSuspended := githubinstallationFields[7].Descriptor()
	// githubinstallation.DefaultSuspended holds the default value on creation for the suspended field.
	githubinstallation.DefaultSuspended = githubinstallationDescSuspended.Default.(bool)
	// githubinstallationDescActive is the schema descriptor for active field.
	githubinstallationDescActive := githubinstallationFields[8].Descriptor()
	// githubinstallation.DefaultActive holds the default value on creation for the active field.
	githubinstallation.DefaultActive = githubinstallationDescActive.Default.(bool)
	// githubinstallationDescID is the schema descriptor for id field.
	githubinstallationDescID := githubinstallationFields[0].Descriptor()
	// githubinstallation.IDValidator is a validator for the "id" field. It is called by the builders before save.
	githubinstallation.IDValidator = githubinstallationDescID.Validators[0].(func(int64) error)
	oauth2codeMixin := schema.Oauth2Code{}.Mixin()
	oauth2codeMixinFields0 := oauth2codeMixin[0].Fields()
	_ = oauth2codeMixinFields0
	oauth2codeMixinFields1 := oauth2codeMixin[1].Fields()
	_ = oauth2codeMixinFields1
	oauth2codeFields := schema.Oauth2Code{}.Fields()
	_ = oauth2codeFields
	// oauth2codeDescCreatedAt is the schema descriptor for created_at field.
	oauth2codeDescCreatedAt := oauth2codeMixinFields1[0].Descriptor()
	// oauth2code.DefaultCreatedAt holds the default value on creation for the created_at field.
	oauth2code.DefaultCreatedAt = oauth2codeDescCreatedAt.Default.(func() time.Time)
	// oauth2codeDescUpdatedAt is the schema descriptor for updated_at field.
	oauth2codeDescUpdatedAt := oauth2codeMixinFields1[1].Descriptor()
	// oauth2code.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	oauth2code.DefaultUpdatedAt = oauth2codeDescUpdatedAt.Default.(func() time.Time)
	// oauth2code.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	oauth2code.UpdateDefaultUpdatedAt = oauth2codeDescUpdatedAt.UpdateDefault.(func() time.Time)
	// oauth2codeDescRevoked is the schema descriptor for revoked field.
	oauth2codeDescRevoked := oauth2codeFields[4].Descriptor()
	// oauth2code.DefaultRevoked holds the default value on creation for the revoked field.
	oauth2code.DefaultRevoked = oauth2codeDescRevoked.Default.(bool)
	// oauth2codeDescID is the schema descriptor for id field.
	oauth2codeDescID := oauth2codeMixinFields0[0].Descriptor()
	// oauth2code.DefaultID holds the default value on creation for the id field.
	oauth2code.DefaultID = oauth2codeDescID.Default.(func() uuid.UUID)
	oauth2tokenMixin := schema.Oauth2Token{}.Mixin()
	oauth2tokenMixinFields0 := oauth2tokenMixin[0].Fields()
	_ = oauth2tokenMixinFields0
	oauth2tokenMixinFields1 := oauth2tokenMixin[1].Fields()
	_ = oauth2tokenMixinFields1
	oauth2tokenFields := schema.Oauth2Token{}.Fields()
	_ = oauth2tokenFields
	// oauth2tokenDescCreatedAt is the schema descriptor for created_at field.
	oauth2tokenDescCreatedAt := oauth2tokenMixinFields1[0].Descriptor()
	// oauth2token.DefaultCreatedAt holds the default value on creation for the created_at field.
	oauth2token.DefaultCreatedAt = oauth2tokenDescCreatedAt.Default.(func() time.Time)
	// oauth2tokenDescUpdatedAt is the schema descriptor for updated_at field.
	oauth2tokenDescUpdatedAt := oauth2tokenMixinFields1[1].Descriptor()
	// oauth2token.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	oauth2token.DefaultUpdatedAt = oauth2tokenDescUpdatedAt.Default.(func() time.Time)
	// oauth2token.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	oauth2token.UpdateDefaultUpdatedAt = oauth2tokenDescUpdatedAt.UpdateDefault.(func() time.Time)
	// oauth2tokenDescRevoked is the schema descriptor for revoked field.
	oauth2tokenDescRevoked := oauth2tokenFields[4].Descriptor()
	// oauth2token.DefaultRevoked holds the default value on creation for the revoked field.
	oauth2token.DefaultRevoked = oauth2tokenDescRevoked.Default.(bool)
	// oauth2tokenDescID is the schema descriptor for id field.
	oauth2tokenDescID := oauth2tokenMixinFields0[0].Descriptor()
	// oauth2token.DefaultID holds the default value on creation for the id field.
	oauth2token.DefaultID = oauth2tokenDescID.Default.(func() uuid.UUID)
	userMixin := schema.User{}.Mixin()
	userMixinFields0 := userMixin[0].Fields()
	_ = userMixinFields0
	userMixinFields1 := userMixin[1].Fields()
	_ = userMixinFields1
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userMixinFields1[0].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userMixinFields1[1].Descriptor()
	// user.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	user.DefaultUpdatedAt = userDescUpdatedAt.Default.(func() time.Time)
	// user.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	user.UpdateDefaultUpdatedAt = userDescUpdatedAt.UpdateDefault.(func() time.Time)
	// userDescID is the schema descriptor for id field.
	userDescID := userMixinFields0[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
}
