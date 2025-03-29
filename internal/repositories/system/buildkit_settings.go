package system_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
)

func (self *SystemRepository) GetBuildkitSettings(ctx context.Context) (*ent.BuildkitSettings, error) {
	return self.base.DB.BuildkitSettings.Query().First(ctx)
}

func (self *SystemRepository) CreateBuildkitSettings(ctx context.Context, replicas int, parallelism int) (*ent.BuildkitSettings, error) {
	return self.base.DB.BuildkitSettings.Create().SetReplicas(replicas).SetMaxParallelism(parallelism).Save(ctx)
}

func (self *SystemRepository) UpdateBuildkitSettings(ctx context.Context, id uuid.UUID, replicas int, parallelism int) (*ent.BuildkitSettings, error) {
	return self.base.DB.BuildkitSettings.UpdateOneID(id).SetReplicas(replicas).SetMaxParallelism(parallelism).Save(ctx)
}
