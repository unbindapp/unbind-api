package project_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	entProject "github.com/unbindapp/unbind-api/ent/project"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
)

type DeleteProjectInput struct {
	TeamID    uuid.UUID `validate:"required,uuid4"`
	ProjectID uuid.UUID `validate:"required,uuid4"`
}

func (self *ProjectService) DeleteProject(ctx context.Context, requesterUserID uuid.UUID, input *DeleteProjectInput, bearerToken string) error {
	// Validate input
	if err := validate.Validator().Struct(input); err != nil {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to delete system resources
		{
			Action:       permission.ActionDelete,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		// Has permission to delete teams
		{
			Action:       permission.ActionDelete,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to delete the specific team
		{
			Action:       permission.ActionDelete,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   input.TeamID.String(),
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return err
	}

	// Check if the team exists
	_, err := self.repo.Team().GetByID(ctx, input.TeamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return err
	}

	// Make sure project exists and is in the team
	project, err := self.repo.Project().GetByID(ctx, input.ProjectID)
	if err != nil {
		if ent.IsNotFound(err) {
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Project not found")
		}
		return err
	}
	if project.TeamID != input.TeamID {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Project not in team")
	}

	// Create kubernetes client
	k8sClient, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return err
	}

	// Delete the project
	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		client := tx.Client()
		if _, err := client.Project.Delete().Where(entProject.ID(input.ProjectID)).Exec(ctx); err != nil {
			return err
		}

		// Delete kubernetes secrets
		if err := self.k8s.DeleteSecret(ctx, project.Edges.Environments[0].KubernetesSecret, project.Edges.Team.Namespace, k8sClient); err != nil {
			return err
		}

		// Delete project secret
		if err := self.k8s.DeleteSecret(ctx, project.KubernetesSecret, project.Edges.Team.Namespace, k8sClient); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
