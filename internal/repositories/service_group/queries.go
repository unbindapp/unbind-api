package servicegroup_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/servicegroup"
)

func (self *ServiceGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.ServiceGroup, error) {
	// Get the service group by ID
	return self.base.DB.ServiceGroup.Get(ctx, id)
}

func (self *ServiceGroupRepository) GetByEnvironmentID(ctx context.Context, environmentID uuid.UUID) ([]*ent.ServiceGroup, error) {
	// Get all service groups by environment ID
	return self.base.DB.ServiceGroup.Query().
		Where(servicegroup.EnvironmentID(environmentID)).
		Order(
			ent.Desc(servicegroup.FieldCreatedAt),
		).
		All(ctx)
}
