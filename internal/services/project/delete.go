package project_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
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
			Action:       schema.ActionAdmin,
			ResourceType: schema.ResourceTypeProject,
			ResourceID:   input.ProjectID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return err
	}

	// Check if the team exists
	team, err := self.repo.Team().GetByID(ctx, input.TeamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return err
	}

	// Make sure project exists and is in the team
	var project *ent.Project
	for _, p := range team.Edges.Projects {
		if p.ID == input.ProjectID {
			project = p
			break
		}
	}
	if project == nil {
		return errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Project not found")
	}

	// Create kubernetes client
	k8sClient, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return err
	}

	environmnets, err := self.repo.Environment().GetForProject(ctx, input.ProjectID)
	if err != nil {
		return err
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return err
	}

	// Delete the project in cascading fashion
	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		// Delete environments
		for _, environment := range environmnets {
			// Delete services
			for _, service := range environment.Edges.Services {
				if err := self.k8s.DeleteUnbindService(ctx, team.Namespace, service.Name); err != nil {
					log.Error("Error deleting service from k8s", "svc", service.Name, "err", err)

					return err
				}

				// Delete secret
				if err := self.k8s.DeleteSecret(ctx, service.KubernetesSecret, team.Namespace, client); err != nil {
					log.Error("Error deleting secret from k8s", "secret", service.KubernetesSecret, "err", err)
					return err
				}

				if err := self.repo.Service().Delete(ctx, tx, service.ID); err != nil {
					return err
				}
			}

			// Delete environment
			if err := self.k8s.DeleteSecret(ctx, environment.KubernetesSecret, team.Namespace, client); err != nil {
				log.Error("Error deleting secret", "secret", environment.KubernetesSecret, "err", err)
			}

			if err := self.repo.Environment().Delete(ctx, tx, environment.ID); err != nil {
				return err
			}
		}

		// Delete project secret
		if err := self.k8s.DeleteSecret(ctx, project.KubernetesSecret, team.Namespace, k8sClient); err != nil {
			return err
		}

		// Delete project by ID
		if err := self.repo.Project().Delete(ctx, tx, input.ProjectID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
