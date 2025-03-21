package project_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type CreateProjectInput struct {
	TeamID      uuid.UUID `validate:"required,uuid4"`
	Name        string    `validate:"required"`
	DisplayName string    `validate:"required"`
	Description string
}

func (self *ProjectService) CreateProject(ctx context.Context, requesterUserID uuid.UUID, input *CreateProjectInput, bearerToken string) (*models.ProjectResponse, error) {
	// Validate input
	if err := validate.Validator().Struct(input); err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to manage system resources
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		// Has permission to manage teams
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to manage the specific team
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   input.TeamID.String(),
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Check if the team exists
	team, err := self.repo.Team().GetByID(ctx, input.TeamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, err
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Create the project
	var project *ent.Project
	var environment *ent.Environment
	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		// Create secret for this project
		secret, _, err := self.k8s.GetOrCreateSecret(ctx, input.Name, team.Namespace, client)
		if err != nil {
			return err
		}

		project, err = self.repo.Project().Create(ctx, tx, input.TeamID, input.Name, input.DisplayName, &input.Description, secret.Name)
		if err != nil {
			return err
		}

		defaultEnvName := "production"
		defaultEnvDescription := "Default production environment"
		// Create a default environment
		name, err := utils.GenerateSlug(defaultEnvName)
		if err != nil {
			return err
		}
		// Create secret for this environment
		secret, _, err = self.k8s.GetOrCreateSecret(ctx, name, team.Namespace, client)
		if err != nil {
			return err
		}

		environment, err = self.repo.Environment().Create(ctx, tx, name, defaultEnvName, defaultEnvDescription, secret.Name, project.ID)
		if err != nil {
			return err
		}

		project.Edges.Environments = append(project.Edges.Environments, environment)
		return nil
	}); err != nil {
		return nil, err
	}

	return models.TransformProjectEntity(project), nil
}
