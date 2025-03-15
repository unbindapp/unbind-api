package project_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/errdefs"
	"github.com/unbindapp/unbind-api/internal/repository"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repository/permissions"
	"github.com/unbindapp/unbind-api/internal/utils"
	"github.com/unbindapp/unbind-api/internal/validate"
)

type CreateProjectInput struct {
	TeamID      uuid.UUID `validate:"required,uuid4"`
	Name        string    `validate:"required"`
	DisplayName string    `validate:"required"`
	Description string    `validate:"required"`
}

type ProjectResponse struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Description  string                 `json:"description"`
	Status       string                 `json:"status"`
	TeamID       uuid.UUID              `json:"team_id"`
	CreatedAt    time.Time              `json:"created_at"`
	Environments []*EnvironmentResponse `json:"environments"`
}

type EnvironmentResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}

func (self *ProjectService) CreateProject(ctx context.Context, requesterUserID uuid.UUID, input *CreateProjectInput) (*ProjectResponse, error) {
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
	_, err := self.repo.Team().GetByID(ctx, input.TeamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, err
	}

	// Create the project
	var project *ent.Project
	var environment *ent.Environment
	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		project, err = self.repo.Project().Create(ctx, tx, input.TeamID, input.Name, input.DisplayName, input.Description)
		if err != nil {
			return err
		}
		// Create a default environment
		name, err := utils.GenerateSlug("production")
		if err != nil {
			return err
		}
		environment, err = self.repo.Environment().Create(ctx, tx, name, "Production", "Default production environment", project.ID)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		DisplayName: project.DisplayName,
		Description: project.Description,
		Status:      project.Status,
		TeamID:      project.TeamID,
		CreatedAt:   project.CreatedAt,
		Environments: []*EnvironmentResponse{
			{
				ID:          environment.ID,
				Name:        environment.Name,
				DisplayName: environment.DisplayName,
				Description: environment.Description,
				Active:      environment.Active,
				CreatedAt:   environment.CreatedAt,
			},
		},
	}, nil
}
