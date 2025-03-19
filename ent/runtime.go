// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/buildjob"
	"github.com/unbindapp/unbind-api/ent/deployment"
	"github.com/unbindapp/unbind-api/ent/environment"
	"github.com/unbindapp/unbind-api/ent/githubapp"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/ent/group"
	"github.com/unbindapp/unbind-api/ent/oauth2code"
	"github.com/unbindapp/unbind-api/ent/oauth2token"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/ent/project"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
	"github.com/unbindapp/unbind-api/ent/team"
	"github.com/unbindapp/unbind-api/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	buildjobMixin := schema.BuildJob{}.Mixin()
	buildjobMixinFields0 := buildjobMixin[0].Fields()
	_ = buildjobMixinFields0
	buildjobMixinFields1 := buildjobMixin[1].Fields()
	_ = buildjobMixinFields1
	buildjobFields := schema.BuildJob{}.Fields()
	_ = buildjobFields
	// buildjobDescCreatedAt is the schema descriptor for created_at field.
	buildjobDescCreatedAt := buildjobMixinFields1[0].Descriptor()
	// buildjob.DefaultCreatedAt holds the default value on creation for the created_at field.
	buildjob.DefaultCreatedAt = buildjobDescCreatedAt.Default.(func() time.Time)
	// buildjobDescUpdatedAt is the schema descriptor for updated_at field.
	buildjobDescUpdatedAt := buildjobMixinFields1[1].Descriptor()
	// buildjob.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	buildjob.DefaultUpdatedAt = buildjobDescUpdatedAt.Default.(func() time.Time)
	// buildjob.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	buildjob.UpdateDefaultUpdatedAt = buildjobDescUpdatedAt.UpdateDefault.(func() time.Time)
	// buildjobDescAttempts is the schema descriptor for attempts field.
	buildjobDescAttempts := buildjobFields[7].Descriptor()
	// buildjob.DefaultAttempts holds the default value on creation for the attempts field.
	buildjob.DefaultAttempts = buildjobDescAttempts.Default.(int)
	// buildjobDescID is the schema descriptor for id field.
	buildjobDescID := buildjobMixinFields0[0].Descriptor()
	// buildjob.DefaultID holds the default value on creation for the id field.
	buildjob.DefaultID = buildjobDescID.Default.(func() uuid.UUID)
	deploymentMixin := schema.Deployment{}.Mixin()
	deploymentMixinFields0 := deploymentMixin[0].Fields()
	_ = deploymentMixinFields0
	deploymentMixinFields1 := deploymentMixin[1].Fields()
	_ = deploymentMixinFields1
	deploymentFields := schema.Deployment{}.Fields()
	_ = deploymentFields
	// deploymentDescCreatedAt is the schema descriptor for created_at field.
	deploymentDescCreatedAt := deploymentMixinFields1[0].Descriptor()
	// deployment.DefaultCreatedAt holds the default value on creation for the created_at field.
	deployment.DefaultCreatedAt = deploymentDescCreatedAt.Default.(func() time.Time)
	// deploymentDescUpdatedAt is the schema descriptor for updated_at field.
	deploymentDescUpdatedAt := deploymentMixinFields1[1].Descriptor()
	// deployment.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	deployment.DefaultUpdatedAt = deploymentDescUpdatedAt.Default.(func() time.Time)
	// deployment.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	deployment.UpdateDefaultUpdatedAt = deploymentDescUpdatedAt.UpdateDefault.(func() time.Time)
	// deploymentDescID is the schema descriptor for id field.
	deploymentDescID := deploymentMixinFields0[0].Descriptor()
	// deployment.DefaultID holds the default value on creation for the id field.
	deployment.DefaultID = deploymentDescID.Default.(func() uuid.UUID)
	environmentMixin := schema.Environment{}.Mixin()
	environmentMixinFields0 := environmentMixin[0].Fields()
	_ = environmentMixinFields0
	environmentMixinFields1 := environmentMixin[1].Fields()
	_ = environmentMixinFields1
	environmentFields := schema.Environment{}.Fields()
	_ = environmentFields
	// environmentDescCreatedAt is the schema descriptor for created_at field.
	environmentDescCreatedAt := environmentMixinFields1[0].Descriptor()
	// environment.DefaultCreatedAt holds the default value on creation for the created_at field.
	environment.DefaultCreatedAt = environmentDescCreatedAt.Default.(func() time.Time)
	// environmentDescUpdatedAt is the schema descriptor for updated_at field.
	environmentDescUpdatedAt := environmentMixinFields1[1].Descriptor()
	// environment.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	environment.DefaultUpdatedAt = environmentDescUpdatedAt.Default.(func() time.Time)
	// environment.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	environment.UpdateDefaultUpdatedAt = environmentDescUpdatedAt.UpdateDefault.(func() time.Time)
	// environmentDescName is the schema descriptor for name field.
	environmentDescName := environmentFields[0].Descriptor()
	// environment.NameValidator is a validator for the "name" field. It is called by the builders before save.
	environment.NameValidator = environmentDescName.Validators[0].(func(string) error)
	// environmentDescActive is the schema descriptor for active field.
	environmentDescActive := environmentFields[3].Descriptor()
	// environment.DefaultActive holds the default value on creation for the active field.
	environment.DefaultActive = environmentDescActive.Default.(bool)
	// environmentDescID is the schema descriptor for id field.
	environmentDescID := environmentMixinFields0[0].Descriptor()
	// environment.DefaultID holds the default value on creation for the id field.
	environment.DefaultID = environmentDescID.Default.(func() uuid.UUID)
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
	githubappDescName := githubappFields[2].Descriptor()
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
	groupMixin := schema.Group{}.Mixin()
	groupMixinFields0 := groupMixin[0].Fields()
	_ = groupMixinFields0
	groupMixinFields1 := groupMixin[1].Fields()
	_ = groupMixinFields1
	groupFields := schema.Group{}.Fields()
	_ = groupFields
	// groupDescCreatedAt is the schema descriptor for created_at field.
	groupDescCreatedAt := groupMixinFields1[0].Descriptor()
	// group.DefaultCreatedAt holds the default value on creation for the created_at field.
	group.DefaultCreatedAt = groupDescCreatedAt.Default.(func() time.Time)
	// groupDescUpdatedAt is the schema descriptor for updated_at field.
	groupDescUpdatedAt := groupMixinFields1[1].Descriptor()
	// group.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	group.DefaultUpdatedAt = groupDescUpdatedAt.Default.(func() time.Time)
	// group.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	group.UpdateDefaultUpdatedAt = groupDescUpdatedAt.UpdateDefault.(func() time.Time)
	// groupDescName is the schema descriptor for name field.
	groupDescName := groupFields[0].Descriptor()
	// group.NameValidator is a validator for the "name" field. It is called by the builders before save.
	group.NameValidator = groupDescName.Validators[0].(func(string) error)
	// groupDescSuperuser is the schema descriptor for superuser field.
	groupDescSuperuser := groupFields[2].Descriptor()
	// group.DefaultSuperuser holds the default value on creation for the superuser field.
	group.DefaultSuperuser = groupDescSuperuser.Default.(bool)
	// groupDescID is the schema descriptor for id field.
	groupDescID := groupMixinFields0[0].Descriptor()
	// group.DefaultID holds the default value on creation for the id field.
	group.DefaultID = groupDescID.Default.(func() uuid.UUID)
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
	permissionMixin := schema.Permission{}.Mixin()
	permissionMixinFields0 := permissionMixin[0].Fields()
	_ = permissionMixinFields0
	permissionMixinFields1 := permissionMixin[1].Fields()
	_ = permissionMixinFields1
	permissionFields := schema.Permission{}.Fields()
	_ = permissionFields
	// permissionDescCreatedAt is the schema descriptor for created_at field.
	permissionDescCreatedAt := permissionMixinFields1[0].Descriptor()
	// permission.DefaultCreatedAt holds the default value on creation for the created_at field.
	permission.DefaultCreatedAt = permissionDescCreatedAt.Default.(func() time.Time)
	// permissionDescUpdatedAt is the schema descriptor for updated_at field.
	permissionDescUpdatedAt := permissionMixinFields1[1].Descriptor()
	// permission.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	permission.DefaultUpdatedAt = permissionDescUpdatedAt.Default.(func() time.Time)
	// permission.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	permission.UpdateDefaultUpdatedAt = permissionDescUpdatedAt.UpdateDefault.(func() time.Time)
	// permissionDescResourceID is the schema descriptor for resource_id field.
	permissionDescResourceID := permissionFields[2].Descriptor()
	// permission.ResourceIDValidator is a validator for the "resource_id" field. It is called by the builders before save.
	permission.ResourceIDValidator = permissionDescResourceID.Validators[0].(func(string) error)
	// permissionDescID is the schema descriptor for id field.
	permissionDescID := permissionMixinFields0[0].Descriptor()
	// permission.DefaultID holds the default value on creation for the id field.
	permission.DefaultID = permissionDescID.Default.(func() uuid.UUID)
	projectMixin := schema.Project{}.Mixin()
	projectMixinFields0 := projectMixin[0].Fields()
	_ = projectMixinFields0
	projectMixinFields1 := projectMixin[1].Fields()
	_ = projectMixinFields1
	projectFields := schema.Project{}.Fields()
	_ = projectFields
	// projectDescCreatedAt is the schema descriptor for created_at field.
	projectDescCreatedAt := projectMixinFields1[0].Descriptor()
	// project.DefaultCreatedAt holds the default value on creation for the created_at field.
	project.DefaultCreatedAt = projectDescCreatedAt.Default.(func() time.Time)
	// projectDescUpdatedAt is the schema descriptor for updated_at field.
	projectDescUpdatedAt := projectMixinFields1[1].Descriptor()
	// project.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	project.DefaultUpdatedAt = projectDescUpdatedAt.Default.(func() time.Time)
	// project.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	project.UpdateDefaultUpdatedAt = projectDescUpdatedAt.UpdateDefault.(func() time.Time)
	// projectDescName is the schema descriptor for name field.
	projectDescName := projectFields[0].Descriptor()
	// project.NameValidator is a validator for the "name" field. It is called by the builders before save.
	project.NameValidator = projectDescName.Validators[0].(func(string) error)
	// projectDescStatus is the schema descriptor for status field.
	projectDescStatus := projectFields[3].Descriptor()
	// project.DefaultStatus holds the default value on creation for the status field.
	project.DefaultStatus = projectDescStatus.Default.(string)
	// projectDescID is the schema descriptor for id field.
	projectDescID := projectMixinFields0[0].Descriptor()
	// project.DefaultID holds the default value on creation for the id field.
	project.DefaultID = projectDescID.Default.(func() uuid.UUID)
	serviceMixin := schema.Service{}.Mixin()
	serviceMixinFields0 := serviceMixin[0].Fields()
	_ = serviceMixinFields0
	serviceMixinFields1 := serviceMixin[1].Fields()
	_ = serviceMixinFields1
	serviceFields := schema.Service{}.Fields()
	_ = serviceFields
	// serviceDescCreatedAt is the schema descriptor for created_at field.
	serviceDescCreatedAt := serviceMixinFields1[0].Descriptor()
	// service.DefaultCreatedAt holds the default value on creation for the created_at field.
	service.DefaultCreatedAt = serviceDescCreatedAt.Default.(func() time.Time)
	// serviceDescUpdatedAt is the schema descriptor for updated_at field.
	serviceDescUpdatedAt := serviceMixinFields1[1].Descriptor()
	// service.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	service.DefaultUpdatedAt = serviceDescUpdatedAt.Default.(func() time.Time)
	// service.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	service.UpdateDefaultUpdatedAt = serviceDescUpdatedAt.UpdateDefault.(func() time.Time)
	// serviceDescName is the schema descriptor for name field.
	serviceDescName := serviceFields[0].Descriptor()
	// service.NameValidator is a validator for the "name" field. It is called by the builders before save.
	service.NameValidator = serviceDescName.Validators[0].(func(string) error)
	// serviceDescID is the schema descriptor for id field.
	serviceDescID := serviceMixinFields0[0].Descriptor()
	// service.DefaultID holds the default value on creation for the id field.
	service.DefaultID = serviceDescID.Default.(func() uuid.UUID)
	serviceconfigMixin := schema.ServiceConfig{}.Mixin()
	serviceconfigMixinFields0 := serviceconfigMixin[0].Fields()
	_ = serviceconfigMixinFields0
	serviceconfigMixinFields1 := serviceconfigMixin[1].Fields()
	_ = serviceconfigMixinFields1
	serviceconfigFields := schema.ServiceConfig{}.Fields()
	_ = serviceconfigFields
	// serviceconfigDescCreatedAt is the schema descriptor for created_at field.
	serviceconfigDescCreatedAt := serviceconfigMixinFields1[0].Descriptor()
	// serviceconfig.DefaultCreatedAt holds the default value on creation for the created_at field.
	serviceconfig.DefaultCreatedAt = serviceconfigDescCreatedAt.Default.(func() time.Time)
	// serviceconfigDescUpdatedAt is the schema descriptor for updated_at field.
	serviceconfigDescUpdatedAt := serviceconfigMixinFields1[1].Descriptor()
	// serviceconfig.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	serviceconfig.DefaultUpdatedAt = serviceconfigDescUpdatedAt.Default.(func() time.Time)
	// serviceconfig.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	serviceconfig.UpdateDefaultUpdatedAt = serviceconfigDescUpdatedAt.UpdateDefault.(func() time.Time)
	// serviceconfigDescPort is the schema descriptor for port field.
	serviceconfigDescPort := serviceconfigFields[3].Descriptor()
	// serviceconfig.DefaultPort holds the default value on creation for the port field.
	serviceconfig.DefaultPort = serviceconfigDescPort.Default.(int)
	// serviceconfigDescReplicas is the schema descriptor for replicas field.
	serviceconfigDescReplicas := serviceconfigFields[4].Descriptor()
	// serviceconfig.DefaultReplicas holds the default value on creation for the replicas field.
	serviceconfig.DefaultReplicas = serviceconfigDescReplicas.Default.(int32)
	// serviceconfigDescAutoDeploy is the schema descriptor for auto_deploy field.
	serviceconfigDescAutoDeploy := serviceconfigFields[5].Descriptor()
	// serviceconfig.DefaultAutoDeploy holds the default value on creation for the auto_deploy field.
	serviceconfig.DefaultAutoDeploy = serviceconfigDescAutoDeploy.Default.(bool)
	// serviceconfigDescPublic is the schema descriptor for public field.
	serviceconfigDescPublic := serviceconfigFields[7].Descriptor()
	// serviceconfig.DefaultPublic holds the default value on creation for the public field.
	serviceconfig.DefaultPublic = serviceconfigDescPublic.Default.(bool)
	// serviceconfigDescID is the schema descriptor for id field.
	serviceconfigDescID := serviceconfigMixinFields0[0].Descriptor()
	// serviceconfig.DefaultID holds the default value on creation for the id field.
	serviceconfig.DefaultID = serviceconfigDescID.Default.(func() uuid.UUID)
	teamMixin := schema.Team{}.Mixin()
	teamMixinFields0 := teamMixin[0].Fields()
	_ = teamMixinFields0
	teamMixinFields1 := teamMixin[1].Fields()
	_ = teamMixinFields1
	teamFields := schema.Team{}.Fields()
	_ = teamFields
	// teamDescCreatedAt is the schema descriptor for created_at field.
	teamDescCreatedAt := teamMixinFields1[0].Descriptor()
	// team.DefaultCreatedAt holds the default value on creation for the created_at field.
	team.DefaultCreatedAt = teamDescCreatedAt.Default.(func() time.Time)
	// teamDescUpdatedAt is the schema descriptor for updated_at field.
	teamDescUpdatedAt := teamMixinFields1[1].Descriptor()
	// team.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	team.DefaultUpdatedAt = teamDescUpdatedAt.Default.(func() time.Time)
	// team.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	team.UpdateDefaultUpdatedAt = teamDescUpdatedAt.UpdateDefault.(func() time.Time)
	// teamDescName is the schema descriptor for name field.
	teamDescName := teamFields[0].Descriptor()
	// team.NameValidator is a validator for the "name" field. It is called by the builders before save.
	team.NameValidator = teamDescName.Validators[0].(func(string) error)
	// teamDescID is the schema descriptor for id field.
	teamDescID := teamMixinFields0[0].Descriptor()
	// team.DefaultID holds the default value on creation for the id field.
	team.DefaultID = teamDescID.Default.(func() uuid.UUID)
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
