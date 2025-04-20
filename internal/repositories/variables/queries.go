package variable_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/variablereference"
)

func (self *VariableRepository) GetReferencesForService(
	ctx context.Context,
	serviceID uuid.UUID,
) ([]*ent.VariableReference, error) {
	// Get all references for a service
	return self.base.DB.VariableReference.Query().
		Where(variablereference.TargetServiceIDEQ(serviceID)).
		Order(
			ent.Desc(variablereference.FieldCreatedAt),
		).
		All(ctx)
}
