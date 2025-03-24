package project_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

func (self *ProjectRepository) Create(ctx context.Context, tx repository.TxInterface, teamID uuid.UUID, name, displayName string, description *string, kubernetesSecret string) (*ent.Project, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}
	// Create the project in the database
	return db.Project.Create().
		SetTeamID(teamID).
		SetName(name).
		SetDisplayName(displayName).
		SetNillableDescription(description).
		SetKubernetesSecret(kubernetesSecret).
		Save(ctx)
}

func (self *ProjectRepository) Update(ctx context.Context, projectID uuid.UUID, displayName string, description *string) (*ent.Project, error) {
	m := self.base.DB.Project.UpdateOneID(projectID)
	if displayName != "" {
		m.SetDisplayName(displayName)
	}
	if description != nil {
		// Reset on empty string
		if *description == "" {
			m.ClearDescription()
		} else {
			m.SetDescription(*description)
		}
	}
	return m.Save(ctx)
}

func (self *ProjectRepository) Delete(ctx context.Context, tx repository.TxInterface, projectID uuid.UUID) error {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}
	return db.Project.DeleteOneID(projectID).Exec(ctx)
}
