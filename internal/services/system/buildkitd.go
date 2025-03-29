package system_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/log"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
)

type BuildkitSettingsResponse struct {
	Replicas       int `json:"replicas" doc:"The number of buildkitd replicas, higher will allow faster concurrent builds"`
	MaxParallelism int `json:"max_parallelism" doc:"buildkitd max_parallelism setting"`
}

func (self *SystemService) GetBuildkitSettings(ctx context.Context, requesterUserID uuid.UUID) (*BuildkitSettingsResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeSystem,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		requesterUserID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	settings, err := self.repo.System().GetBuildkitSettings(ctx)
	if err != nil {
		log.Errorf("Failed to get buildkit settings from DB: %v", err)
		return nil, err
	}
	return &BuildkitSettingsResponse{
		Replicas:       settings.Replicas,
		MaxParallelism: settings.MaxParallelism,
	}, nil
}

// UpdateBuildkitSettings updates the buildkit settings in the database and kubernetes
func (self *SystemService) UpdateBuildkitSettings(ctx context.Context, requesterUserID uuid.UUID, replicas int, maxParallelism int) (*BuildkitSettingsResponse, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeSystem,
		},
	}

	if err := self.repo.Permissions().Check(
		ctx,
		requesterUserID,
		permissionChecks,
	); err != nil {
		return nil, err
	}

	// Get buildkit settings from DB
	settings, err := self.repo.System().GetBuildkitSettings(ctx)
	if err != nil {
		log.Errorf("Failed to get buildkit settings from DB: %v", err)
		return nil, err
	}

	updatedSettings, err := self.repo.System().UpdateBuildkitSettings(ctx, settings.ID, replicas, maxParallelism)
	if err != nil {
		log.Errorf("Failed to update buildkit settings in DB: %v", err)
		return nil, err
	}
	return &BuildkitSettingsResponse{
		Replicas:       updatedSettings.Replicas,
		MaxParallelism: updatedSettings.MaxParallelism,
	}, nil
}
