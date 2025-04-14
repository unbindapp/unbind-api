// Code generated by ifacemaker; DO NOT EDIT.

package repositories

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	deployment_repo "github.com/unbindapp/unbind-api/internal/repositories/deployment"
	environment_repo "github.com/unbindapp/unbind-api/internal/repositories/environment"
	github_repo "github.com/unbindapp/unbind-api/internal/repositories/github"
	group_repo "github.com/unbindapp/unbind-api/internal/repositories/group"
	oauth_repo "github.com/unbindapp/unbind-api/internal/repositories/oauth"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	project_repo "github.com/unbindapp/unbind-api/internal/repositories/project"
	service_repo "github.com/unbindapp/unbind-api/internal/repositories/service"
	system_repo "github.com/unbindapp/unbind-api/internal/repositories/system"
	team_repo "github.com/unbindapp/unbind-api/internal/repositories/team"
	user_repo "github.com/unbindapp/unbind-api/internal/repositories/user"
	webhook_repo "github.com/unbindapp/unbind-api/internal/repositories/webhook"
)

// RepositoriesInterface ...
type RepositoriesInterface interface {
	// Ent() returns the ent client
	Ent() *ent.Client
	// Github returns the GitHub repository
	Github() github_repo.GithubRepositoryInterface
	// User returns the User repository
	User() user_repo.UserRepositoryInterface
	// Oauth returns the OAuth repository
	Oauth() oauth_repo.OauthRepositoryInterface
	// Group returns the Group repository
	Group() group_repo.GroupRepositoryInterface
	// Project returns the Project repository
	Project() project_repo.ProjectRepositoryInterface
	// Team returns the Team repository
	Team() team_repo.TeamRepositoryInterface
	// Permissions returns the Permissions repository
	Permissions() permissions_repo.PermissionsRepositoryInterface
	// Environment returns the Environment repository
	Environment() environment_repo.EnvironmentRepositoryInterface
	// Service returns the Service repository
	Service() service_repo.ServiceRepositoryInterface
	// Deployment returns the Deployment repository
	Deployment() deployment_repo.DeploymentRepositoryInterface
	// System returns the System repository
	System() system_repo.SystemRepositoryInterface
	// Webhooks returns the Webhook repository
	Webhooks() webhook_repo.WebhookRepositoryInterface
	WithTx(ctx context.Context, fn func(tx repository.TxInterface) error) error
}
