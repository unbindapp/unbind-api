package service_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	webhooks_service "github.com/unbindapp/unbind-api/internal/services/webooks"
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
		// Cancel deployments
		if err := self.deploymentController.CancelExistingJobs(ctx, service.ID); err != nil {
			log.Warnf("Error cancelling jobs for service %s: %v", service.Name, err)
		}

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

	// Trigger webhook
	go func() {
		event := schema.WebhookEventServiceDeleted
		level := webhooks_service.WebhookLevelError

		// Construct URL
		url, _ := utils.JoinURLPaths(self.cfg.ExternalUIUrl, project.TeamID.String(), "project", project.ID.String())
		// Get user
		user, err := self.repo.User().GetByID(context.Background(), requesterUserID)
		if err != nil {
			log.Errorf("Failed to get user %s: %v", requesterUserID.String(), err)
			return
		}
		data := webhooks_service.WebookData{
			Title:       "Service Deleted",
			Url:         url,
			Description: fmt.Sprintf("A service has been deleted in project %s by %s", project.DisplayName, user.Email),
			Fields: []webhooks_service.WebhookDataField{
				{
					Name:  "Service",
					Value: service.DisplayName,
				},
				{
					Name:  "Environment",
					Value: service.Edges.Environment.DisplayName,
				},
			},
		}

		if service.Description != "" {
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Description",
				Value: service.Description,
			})
		}

		data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
			Name:  "Type",
			Value: string(service.Edges.ServiceConfig.Type),
		})
		data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
			Name:  "Subtype",
			Value: service.Edges.ServiceConfig.Icon,
		})

		if err := self.webhookService.TriggerWebhooks(context.Background(), level, event, data); err != nil {
			log.Errorf("Failed to trigger webhook %s: %v", event, err)
		}
	}()

	return nil
}
