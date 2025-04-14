package webhooks_service

import (
	"context"

	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

type WebookData struct {
	Title       string             `json:"title"`
	Url         string             `json:"url"`
	Description string             `json:"description"`
	Username    string             `json:"username"`
	Fields      []WebhookDataField `json:"fields"`
}

type WebhookDataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type WebhookLevel string

const (
	WebhookLevelError   WebhookLevel = "error"
	WebhookLevelWarning WebhookLevel = "warning"
	WebhookLevelInfo    WebhookLevel = "info"
)

func (self *WebhookLevel) DiscordColor() *string {
	switch *self {
	case WebhookLevelError:
		return utils.ToPtr("9110797")
	case WebhookLevelWarning:
		return utils.ToPtr("8396800")
	default:
		return utils.ToPtr("802316") // Success
	}
}

func (self *WebhooksService) TriggerWebhooks(ctx context.Context, level WebhookLevel, event schema.WebhookEvent, message WebookData) error {
	// Get all webhooks matching event
	webhooks, err := self.repo.Webhooks().GetWebhooksForEvent(ctx, event)
	if err != nil {
		return err
	}

	for _, webhook := range webhooks {
		target, err := self.DetectTargetFromURL(webhook.URL)
		if err != nil {
			log.Errorf("Failed to detect target from webhook URL %s: %v", webhook.URL, err)
		}

		switch target {
		case schema.WebhookTargetDiscord:
			return self.sendDiscordWebhook(level, event, message)
		}
	}

	return nil
}
