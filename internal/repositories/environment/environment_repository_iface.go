// Code generated by ifacemaker; DO NOT EDIT.

package environment_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

// EnvironmentRepositoryInterface ...
type EnvironmentRepositoryInterface interface {
	Create(ctx context.Context, tx repository.TxInterface, name, displayName, kuberneteSecret string, description *string, projectID uuid.UUID) (*ent.Environment, error)
	Delete(ctx context.Context, tx repository.TxInterface, environmentID uuid.UUID) error
	Update(ctx context.Context, environmentID uuid.UUID, displayName *string, description *string) (*ent.Environment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ent.Environment, error)
	// Return all environments for a project with service edge populated
	GetForProject(ctx context.Context, tx repository.TxInterface, projectID uuid.UUID) ([]*ent.Environment, error)
}
