package environment_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/environment"
)

func (self *EnvironmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Environment, error) {
	return self.base.DB.Environment.Query().Where(environment.ID(id)).WithProject().Only(ctx)
}
