package service_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Get a service by ID
func (self *ServiceService) GetDNSForService(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, teamID, projectID, environmentID, serviceID uuid.UUID) (*models.EndpointDiscovery, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to admin service
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeService,
			ResourceID:   serviceID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Verify inputs
	_, project, err := self.VerifyInputs(ctx, teamID, projectID, environmentID)
	if err != nil {
		return nil, err
	}

	// Get service
	service, err := self.repo.Service().GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Get discovery
	endpoints, err := self.k8s.DiscoverEndpointsByLabels(
		ctx,
		project.Edges.Team.Namespace,
		map[string]string{
			"unbind-service": serviceID.String(),
		},
		false,
		client,
	)

	if err != nil {
		return nil, err
	}

	if service.Edges.ServiceConfig.Database != nil {
		for _, endpoint := range endpoints.Internal {
			// ! TODO parse database commonly since this is repeated logic
			switch *service.Edges.ServiceConfig.Database {
			case "redis", "postgres":
				// Get database password from secret
				secret, err := self.k8s.GetSecret(ctx, service.KubernetesName, service.Edges.Environment.Edges.Project.Edges.Team.Namespace, client)
				if err != nil {
					return nil, err
				}
				username := string(secret.Data["DATABASE_USERNAME"])
				password := string(secret.Data["DATABASE_PASSWORD"])

				if *service.Edges.ServiceConfig.Database == "redis" {
					endpoint.DNS = fmt.Sprintf("redis://%s:%s@%s:%d", username, password, endpoint.DNS, 6379)
				}
				if *service.Edges.ServiceConfig.Database == "postgres" {
					endpoint.DNS = fmt.Sprintf("postgresql://%s:%s@%s:%d/postgres?sslmode=disable", username, password, endpoint.DNS, 5432)
				}
			}
		}
	}

	return endpoints, nil
}
