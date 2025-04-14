package webhook_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/webhook"
)

func (self *WebhookRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Webhook, error) {
	return self.base.DB.Webhook.
		Query().
		Where(webhook.ID(id)).
		Only(ctx)
}

func (self *WebhookRepository) GetByTeam(ctx context.Context, teamID uuid.UUID) ([]*ent.Webhook, error) {
	return self.base.DB.Webhook.Query().
		Where(webhook.TeamID(teamID)).
		All(ctx)
}

func (self *WebhookRepository) GetByProject(ctx context.Context, projectID uuid.UUID) ([]*ent.Webhook, error) {
	return self.base.DB.Webhook.Query().
		Where(webhook.ProjectID(projectID)).
		All(ctx)
}
