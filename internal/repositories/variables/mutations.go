package variable_repo

import (
	"context"
	"strings"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *VariableRepository) CreateReference(ctx context.Context, input *models.CreateVariableReferenceInput) (*ent.VariableReference, error) {
	// Create variable reference
	return self.base.DB.VariableReference.Create().
		SetTargetServiceID(input.TargetServiceID).
		SetTargetName(strings.TrimSpace(input.TargetName)).
		SetType(input.Type).
		SetSources(input.Sources).
		SetValueTemplate(input.ValueTemplate).
		Save(ctx)
}
