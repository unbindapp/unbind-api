package project_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/project"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Project, error) {
	return self.base.DB.Project.Query().Where(project.ID(id)).WithTeam().WithEnvironments().Only(ctx)
}

func (self *ProjectRepository) GetTeamID(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	team, err := self.base.DB.Project.Query().Where(project.ID(id)).QueryTeam().Only(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return team.ID, nil
}

func (self *ProjectRepository) GetByTeam(ctx context.Context, teamID uuid.UUID, sortField models.SortByField, sortOrder models.SortOrder) ([]*ent.Project, error) {
	q := self.base.DB.Project.Query().
		Where(project.TeamID(teamID)).
		WithEnvironments()

	switch sortField {
	case models.SortByCreatedAt:
		ent.Asc(project.FieldCreatedAt)
		q.Order(sortOrder.SortFunction()(project.FieldCreatedAt))
	case models.SortByUpdatedAt:
		q.Order(sortOrder.SortFunction()(project.FieldUpdatedAt))
	}

	return q.All(ctx)
}
