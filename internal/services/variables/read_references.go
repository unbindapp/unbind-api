package variables_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *VariablesService) GetAvailableVariableReferences(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, teamID uuid.UUID) (*models.AvailableVariableReferenceResponse, error) {
	// ! TODO - we're going to need to change all of our permission checks to filter not reject
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeTeam,
			ResourceID:   teamID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Get available variable references
	team, err := self.repo.Team().GetByID(ctx, teamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, err
	}

	// Get all kubernetes secrets in the team
	teamSecret := team.KubernetesSecret
	projectSecrets := make(map[uuid.UUID]string)
	environmentSecrets := make(map[uuid.UUID]string)
	serviceSecrets := make(map[uuid.UUID]string)
	for _, project := range team.Edges.Projects {
		projectSecrets[project.ID] = project.KubernetesSecret
		for _, environment := range project.Edges.Environments {
			environmentSecrets[environment.ID] = environment.KubernetesSecret
			for _, service := range environment.Edges.Services {
				serviceSecrets[service.ID] = service.KubernetesSecret
			}
		}
	}

	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Get all kubernetes secrets in the team
	k8sSecrets, err := self.k8s.GetAllSecrets(
		ctx,
		team.ID,
		teamSecret,
		projectSecrets,
		environmentSecrets,
		serviceSecrets,
		client,
		team.Namespace,
	)

	if err != nil {
		return nil, err
	}

	// Get all the DNS/endpoints
	endpoints, err := self.k8s.DiscoverEndpointsByLabels(
		ctx,
		team.Namespace,
		map[string]string{
			"unbind-team": team.ID.String(),
		},
		client,
	)

	if err != nil {
		return nil, err
	}

	return models.TransformAvailableVariableResponse(k8sSecrets, endpoints), nil
}
