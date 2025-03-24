package service_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// UpdateServiceConfigInput defines the input for updating a service configuration
type UpdateServiceInput struct {
	TeamID        uuid.UUID `validate:"required,uuid4" required:"true" json:"team_id"`
	ProjectID     uuid.UUID `validate:"required,uuid4" required:"true" json:"project_id"`
	EnvironmentID uuid.UUID `validate:"required,uuid4" required:"true" json:"environment_id"`
	ServiceID     uuid.UUID `validate:"required,uuid4" required:"true" json:"service_id"`
	DisplayName   *string   `validate:"optional" required:"false" json:"display_name"`
	Description   *string   `validate:"optional" required:"false" json:"description"`

	// Configuration
	GitBranch  *string                `json:"git_branch,omitempty" required:"false"`
	Type       *schema.ServiceType    `json:"type,omitempty" required:"false"`
	Builder    *schema.ServiceBuilder `json:"builder,omitempty" required:"false"`
	Host       *string                `json:"host,omitempty" required:"false"`
	Port       *int                   `validate:"min=1,max=65535" json:"port,omitempty" required:"false"`
	Replicas   *int32                 `validate:"min=1,max=10" json:"replicas,omitempty" required:"false"`
	AutoDeploy *bool                  `json:"auto_deploy,omitempty" required:"false"`
	RunCommand *string                `json:"run_command,omitempty" required:"false"`
	Public     *bool                  `json:"public,omitempty" required:"false"`
	Image      *string                `json:"image,omitempty" required:"false"`
}

// UpdateService updates a service and its configuration
func (self *ServiceService) UpdateService(ctx context.Context, requesterUserID uuid.UUID, input *UpdateServiceInput) (*models.ServiceResponse, error) {
	// Validate input
	if err := validate.Validator().Struct(input); err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
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
			ResourceID:   input.TeamID.String(),
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
			ResourceID:   input.ProjectID.String(),
		},
		// Has permission to manage this specific environment
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   input.EnvironmentID.String(),
		},
		// Has permission to update this service
		{
			Action:       permission.ActionUpdate,
			ResourceType: permission.ResourceTypeService,
			ResourceID:   input.ServiceID.String(),
		},
	}

	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Verify inputs
	_, _, err := self.VerifyInputs(ctx, input.TeamID, input.ProjectID, input.EnvironmentID)
	if err != nil {
		return nil, err
	}

	// Perform update
	var service *ent.Service
	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		// Update the service
		if err := self.repo.Service().Update(ctx, tx, input.ServiceID, input.DisplayName, input.Description); err != nil {
			return fmt.Errorf("failed to update service: %w", err)
		}

		host := input.Host
		if service.Edges.ServiceConfig.Host == nil && input.Public != nil && *input.Public && input.Host == nil {
			host, err := utils.GenerateSubdomain(service.DisplayName, service.Edges.Environment.DisplayName, self.cfg.ExternalURL)
			if err != nil {
				log.Warn("failed to generate subdomain", "error", err)
			} else {
				input.Host = &host
			}
		}

		// Update the service config
		if err := self.repo.Service().UpdateConfig(ctx,
			tx,
			input.ServiceID,
			input.Type,
			input.Builder,
			input.GitBranch,
			input.Port,
			host,
			input.Replicas,
			input.AutoDeploy,
			input.RunCommand,
			input.Public,
			input.Image); err != nil {
			return fmt.Errorf("failed to update service config: %w", err)
		}

		// Re-fetch the service
		service, err = self.repo.Service().GetByID(ctx, input.ServiceID)
		if err != nil {
			return fmt.Errorf("failed to re-fetch service: %w", err)
		}

		return nil

	}); err != nil {
		return nil, err
	}

	return models.TransformServiceEntity(service), nil
}
