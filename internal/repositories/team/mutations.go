package team_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
)

func (self *TeamRepository) Update(ctx context.Context, teamID uuid.UUID, displayName string) (*ent.Team, error) {
	return self.base.DB.Team.UpdateOneID(teamID).
		SetDisplayName(displayName).
		Save(ctx)
}
