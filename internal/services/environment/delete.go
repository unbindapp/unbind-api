package environment_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/common/log"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
)

func (self *EnvironmentService) DeleteEnvironmentByID(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, teamID, projectID, environmentID uuid.UUID) error {
	// Check permissions
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to manage teams
		{
			Action:       permission.ActionDelete,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to manage this team
		{
			Action:       permission.ActionDelete,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   teamID.String(),
		},
		// Has permission to manage this project
		{
			Action:       permission.ActionDelete,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   projectID.String(),
		},
		// Has permission to manage this specific environment
		{
			Action:       permission.ActionDelete,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   environmentID.String(),
		},
	}

	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return err
	}

	// Verify inputs
	team, environment, err := self.VerifyInputs(ctx, teamID, projectID, environmentID)
	if err != nil {
		return err
	}

	// Get services in this environment
	services, err := self.repo.Service().GetByEnvironmentID(ctx, environmentID)
	if err != nil {
		return err
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return err
	}

	// Delete kubernetes resources, db resource
	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		// Delete services
		for _, service := range services {
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

		if err := self.repo.Environment().Delete(ctx, tx, environmentID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
