package webhook_repo

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
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

func (self *WebhookRepository) GetWebhooksForEvent(ctx context.Context, event schema.WebhookEvent) ([]*ent.Webhook, error) {
	return self.base.DB.Webhook.Query().
		Where(func(s *sql.Selector) {
			s.Where(sqljson.ValueContains(s.C(webhook.FieldEvents), event))
		}).
		All(ctx)
}
