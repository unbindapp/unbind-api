package storage_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *StorageService) ListS3StorageBackends(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, teamID uuid.UUID) ([]*models.S3Response, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		// Team viewer can view s3 sources
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

	// Check if the team exists
	team, err := self.repo.Team().GetByID(ctx, teamID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Team not found")
		}
		return nil, err
	}

	// Get all s3 sources for this team
	s3Sources, err := self.repo.S3().GetByTeam(ctx, team.ID)
	if err != nil {
		return nil, err
	}
	if len(s3Sources) == 0 {
		return []*models.S3Response{}, nil
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Create credential map
	accessKeyMap := make(map[uuid.UUID]string)
	secretKeyMap := make(map[uuid.UUID]string)
	for _, s3 := range s3Sources {
		// Get the secret
		secret, err := self.k8s.GetSecret(ctx, s3.KubernetesSecret, team.Namespace, client)
		if err != nil {
			log.Errorf("Failed to get secret %s for s3 %s: %v", s3.KubernetesSecret, s3.ID, err)
			continue
		}

		accessKeyMap[s3.ID] = string(secret.Data["access_key_id"])
		secretKeyMap[s3.ID] = string(secret.Data["secret_key"])
	}

	return models.TransformS3Entities(s3Sources, accessKeyMap, secretKeyMap), nil

}
