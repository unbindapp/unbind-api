package service_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/models"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
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
	env, project, err := self.VerifyInputs(ctx, teamID, projectID, environmentID)
	if err != nil {
		return nil, err
	}

	// Get service
	service, err := self.repo.Service().GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	if service.EnvironmentID != env.ID {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
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
		true,
		client,
	)

	if err != nil {
		return nil, err
	}

	// Build a map of discovered hosts
	endpointMap := make(map[string]struct{})
	for _, host := range endpoints.External {
		if !host.IsIngress {
			continue
		}
		endpointMap[host.Host] = struct{}{}
	}

	// Append hosts missing from discovery
	for _, host := range service.Edges.ServiceConfig.Hosts {
		if _, exists := endpointMap[host.Host]; !exists {
			path := host.Path
			if path == "" {
				path = "/"
			}
			var port int32 = 443
			if host.TargetPort != nil {
				port = *host.TargetPort
			}
			newHost := models.IngressEndpoint{
				KubernetesName: service.KubernetesName,
				IsIngress:      true,
				Host:           host.Host,
				Path:           path,
				Port: schema.PortSpec{
					Port:     port,
					Protocol: utils.ToPtr(schema.ProtocolTCP),
				},
				DNSStatus:     models.DNSStatusUnknown,
				TlsStatus:     models.TlsStatusPending,
				TeamID:        project.Edges.Team.ID,
				ProjectID:     project.ID,
				EnvironmentID: env.ID,
				ServiceID:     serviceID,
			}

			if host.Path != "" {
				newHost.Path = host.Path
			}

			endpoints.External = append(endpoints.External, newHost)
		}

	}

	return endpoints, nil
}
