package service_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
)

func (self *ServiceService) DeleteServiceByID(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, teamID, projectID, environmentID, serviceID uuid.UUID) error {
	// Check permissions
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to admin service
		{
			Action:       schema.ActionAdmin,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   serviceID,
		},
	}

	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return err
	}

	_, project, err := self.VerifyInputs(ctx, teamID, projectID, environmentID)
	if err != nil {
		return err
	}
	team := project.Edges.Team

	// Get the service
	service, err := self.repo.Service().GetByID(ctx, serviceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
		}
		return err
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return err
	}

	// Delete kubernetes resources, db resource
	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		if err := self.k8s.DeleteUnbindService(ctx, team.Namespace, service.Name); err != nil {
			log.Error("Error deleting service from k8s", "svc", service.Name, "err", err)

			return err
		}

		// Delete secret
		if err := self.k8s.DeleteSecret(ctx, service.KubernetesSecret, team.Namespace, client); err != nil {
			log.Error("Error deleting secret from k8s", "secret", service.KubernetesSecret, "err", err)
			return err
		}

		if err := self.repo.Service().Delete(ctx, tx, serviceID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
