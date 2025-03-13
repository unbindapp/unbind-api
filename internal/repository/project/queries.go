package project_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/project"
)

func (self *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Project, error) {
	return self.base.DB.Project.Query().Where(project.ID(id)).WithTeam().Only(ctx)
}

func (self *ProjectRepository) GetTeamID(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	team, err := self.base.DB.Project.Query().Where(project.ID(id)).QueryTeam().Only(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return team.ID, nil
}
