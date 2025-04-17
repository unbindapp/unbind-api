package variable_repo

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *VariableRepository) CreateReference(ctx context.Context, input *models.CreateVariableReferenceInput) (*ent.VariableReference, error) {
	// Create variable reference
	return self.base.DB.VariableReference.Create().
		SetTargetServiceID(input.TargetServiceID).
		SetTargetName(strings.TrimSpace(input.TargetName)).
		SetSources(input.Sources).
		SetValueTemplate(input.ValueTemplate).
		Save(ctx)
}

func (self *VariableRepository) AttachError(ctx context.Context, id uuid.UUID, err error) (*ent.VariableReference, error) {
	// Attach error to variable reference
	return self.base.DB.VariableReference.UpdateOneID(id).
		SetError(err.Error()).
		Save(ctx)
}
