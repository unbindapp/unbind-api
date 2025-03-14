package project_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
)

func (self *ProjectRepository) Create(ctx context.Context, teamID uuid.UUID, name, displayName, description string) (*ent.Project, error) {
	// Create the project in the database
	return self.base.DB.Project.Create().
		SetTeamID(teamID).
		SetName(name).
		SetDisplayName(displayName).
		SetDescription(description).
		Save(ctx)
}

func (self *ProjectRepository) Update(ctx context.Context, projectID uuid.UUID, displayName, description string) (*ent.Project, error) {
	m := self.base.DB.Project.UpdateOneID(projectID)
	if displayName != "" {
		m.SetDisplayName(displayName)
	}
	if description != "" {
		m.SetDescription(description)
	}
	return m.Save(ctx)
}
