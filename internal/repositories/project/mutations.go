package project_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

func (self *ProjectRepository) Create(ctx context.Context, tx repository.TxInterface, teamID uuid.UUID, name, displayName, description string) (*ent.Project, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}
	// Create the project in the database
	return db.Project.Create().
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

func (self *ProjectRepository) Delete(ctx context.Context, projectID uuid.UUID) error {
	return self.base.DB.Project.DeleteOneID(projectID).Exec(ctx)
}
