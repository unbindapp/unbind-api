// Code generated by ifacemaker; DO NOT EDIT.

package project_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/repository"
)

// ProjectRepositoryInterface ...
type ProjectRepositoryInterface interface {
	Create(ctx context.Context, tx repository.TxInterface, teamID uuid.UUID, name, displayName, description string) (*ent.Project, error)
	Update(ctx context.Context, projectID uuid.UUID, displayName, description string) (*ent.Project, error)
	Delete(ctx context.Context, projectID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*ent.Project, error)
	GetTeamID(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	GetByTeam(ctx context.Context, teamID uuid.UUID) ([]*ent.Project, error)
}
