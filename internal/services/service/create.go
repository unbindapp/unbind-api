package service_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// CreateServiceInput defines the input for creating a new service
type CreateServiceInput struct {
	EnvironmentID uuid.UUID `validate:"required,uuid4" required:"true" json:"environment_id"`
	DisplayName   string    `validate:"required" required:"true" json:"display_name"`
	Description   string    `validate:"optional" json:"description"`
	// ! TODO - infer this? make it optional? Make validation dynamic somehow?
	Type    service.Type    `validate:"required,oneof=database api web custom" required:"true" json:"type" doc:"One of database, api, web, or custom"`
	Subtype service.Subtype `validate:"required,oneof=react go node next other" required:"true" json:"subtype" doc:"One of react, go, node, next, other"`

	// GitHub integration
	GitHubInstallationID *int64  `json:"github_installation_id,omitempty"`
	RepositoryOwner      *string `json:"repository_owner,omitempty"`
	RepositoryName       *string `json:"repository_name,omitempty"`

	// Configuration
	GitBranch  *string `json:"git_branch,omitempty"`
	Host       *string `json:"host,omitempty"`
	Port       *int    `validate:"min=1,max=65535" json:"port,omitempty"`
	Replicas   *int32  `validate:"min=1,max=10" json:"replicas,omitempty"`
	AutoDeploy *bool   `json:"auto_deploy,omitempty"`
	RunCommand *string `json:"run_command,omitempty"`
	Public     *bool   `json:"public,omitempty"`
	Image      *string `json:"image,omitempty"`
}

// CreateService creates a new service and its configuration
func (self *ServiceService) CreateService(ctx context.Context, requesterUserID uuid.UUID, input *CreateServiceInput) (*models.ServiceResponse, error) {
	// Validate input
	if err := validate.Validator().Struct(input); err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	// Validate that if GitHub info is provided, all fields are set
	if input.GitHubInstallationID != nil {
		if input.RepositoryOwner == nil || input.RepositoryName == nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
				"Both GitHub repository owner and name must be provided together")
		}
	}
	// Verify that the environment exists
	environment, err := self.repo.Environment().GetByID(ctx, input.EnvironmentID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Environment not found")
		}
		return nil, err
	}

	// Get project ID from environment for permission checking
	project, err := self.repo.Project().GetByID(ctx, environment.ProjectID)
	if err != nil {
		return nil, err
	}

	// Check permissions
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to manage teams
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to manage this team
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   project.TeamID.String(),
		},
		// Has permission to manage projects
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   "*",
		},
		// Has permission to manage this project
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   project.ID.String(),
		},
		// Has permission to manage this specific environment
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   input.EnvironmentID.String(),
		},
	}

	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// If GitHub integration is provided, verify repository access
	if input.GitHubInstallationID != nil {
		// Get GitHub installation
		installation, err := self.repo.Github().GetInstallationByID(ctx, *input.GitHubInstallationID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "GitHub installation not found")
			}
			return nil, err
		}

		// Verify repository access
		canAccess, err := self.githubClient.VerifyRepositoryAccess(ctx, installation, *input.RepositoryOwner, *input.RepositoryName)
		if err != nil {
			log.Error("Error verifying repository access", "err", err)
			return nil, err
		}

		if !canAccess {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
				"Repository not accessible with the specified GitHub installation")
		}
	}

	// Create service and config in a transaction
	var service *ent.Service
	var serviceConfig *ent.ServiceConfig

	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		// Create the service
		createService, err := self.repo.Service().Create(ctx, tx,
			input.DisplayName,
			input.Description,
			input.Type,
			input.Subtype,
			input.EnvironmentID,
			input.GitHubInstallationID,
			input.RepositoryName)
		if err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}
		service = createService

		// ctx context.Context, tx repository.TxInterface, serviceID uuid.UUID, gitBranch *string, port int, host *string, replicas *int32, autoDeploy *bool, runCommand *string, public *bool, image *string) (*ent.ServiceConfig, error)

		// Create the service config
		serviceConfig, err = self.repo.Service().CreateConfig(ctx, tx,
			service.ID,
			input.GitBranch,
			input.Port,
			input.Host,
			input.Replicas,
			input.AutoDeploy,
			input.RunCommand,
			input.Public,
			input.Image)
		if err != nil {
			return fmt.Errorf("failed to create service config: %w", err)
		}
		service.Edges.ServiceConfig = serviceConfig
		return nil

	}); err != nil {
		return nil, err
	}

	return models.TransformServiceEntity(service), nil
}
