package environment_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
)

func (self *EnvironmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Environment, error) {
	return self.base.DB.Environment.Get(ctx, id)
}
