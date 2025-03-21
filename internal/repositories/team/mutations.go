package team_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
)

func (self *TeamRepository) Update(ctx context.Context, teamID uuid.UUID, displayName string, description *string) (*ent.Team, error) {
	m := self.base.DB.Team.UpdateOneID(teamID)
	if displayName != "" {
		m.SetDisplayName(displayName)
	}
	if description != nil {
		if *description == "" {
			m.ClearDescription()
		} else {
			m.SetDescription(*description)
		}
	}
	return m.Save(ctx)
}
