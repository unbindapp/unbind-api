// Code generated by ifacemaker; DO NOT EDIT.

package variable_repo

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// VariableRepositoryInterface ...
type VariableRepositoryInterface interface {
	CreateReference(ctx context.Context, input *models.CreateVariableReferenceInput) (*ent.VariableReference, error)
}
