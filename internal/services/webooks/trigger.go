package webhooks_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

type WebookData struct {
	Title       string             `json:"title"`
	Url         string             `json:"url"`
	Description string             `json:"description"`
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

func (self *WebhookLevel) DecimalColor() *string {
	switch *self {
	case WebhookLevelError:
		return utils.ToPtr("9110797")
	case WebhookLevelWarning:
		return utils.ToPtr("8396800")
	default:
		return utils.ToPtr("802316") // Success
	}
}

// Helper method for WebhookLevel to return Slack color
func (level WebhookLevel) HexColor() *string {
	switch level {
	case WebhookLevelError:
		return utils.ToPtr("#8B0000")
	case WebhookLevelWarning:
		return utils.ToPtr("#802000")
	default:
		return utils.ToPtr("#0C3B0C")
	}
}

// Emoji indicator for the ones not supporting style (telegram)
func (level WebhookLevel) Emoji() string {
	var levelBar string
	switch level {
	case WebhookLevelError:
		levelBar = "🔴"
	case WebhookLevelWarning:
		levelBar = "🟠"
	default:
		levelBar = "🟢"
	}
	return levelBar
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
			return self.sendDiscordWebhook(level, event, message, webhook.URL)
		case schema.WebhookTargetSlack:
			return self.sendSlackWebhook(level, event, message, webhook.URL)
		case schema.WebhookTargetTelegram:
			return self.sendTelegramWebhook(level, event, message, webhook.URL)
		default:
			// Just encode our payload
			msg := DefaultPayload{
				Level: level,
				Event: event,
				Data:  message,
			}

			// Encode the payload
			payload := new(bytes.Buffer)
			err := json.NewEncoder(payload).Encode(msg)
			if err != nil {
				log.Errorf("Failed to encode slack webhook payload: %v", err)
				return err
			}

			// Create the request
			req, err := http.NewRequest(http.MethodPost, webhook.URL, payload)
			if err != nil {
				log.Errorf("Failed to create slack webhook request: %v", err)
				return err
			}
			req.Header.Set("Content-Type", "application/json")

			// Send the request
			resp, err := self.httpClient.Do(req)
			if err != nil {
				log.Errorf("Failed to send slack webhook: %v", err)
				return err
			}
			defer resp.Body.Close()

			// Check response
			if resp.StatusCode != http.StatusOK {
				bodyBytes, readErr := io.ReadAll(resp.Body)
				if readErr != nil {
					log.Errorf("Failed to send slack webhook: status=%d, error reading body: %v",
						resp.StatusCode, readErr)
					return fmt.Errorf("failed to send slack webhook: %s, couldn't read response", resp.Status)
				}

				// Log both status code and response body
				bodyString := string(bodyBytes)
				log.Errorf("Failed to send slack webhook: status=%d, body=%s",
					resp.StatusCode, bodyString)
				return fmt.Errorf("failed to send slack webhook: %s, response: %s", resp.Status, bodyString)
			}
		}
	}

	return nil
}

type DefaultPayload struct {
	Level WebhookLevel        `json:"level"`
	Event schema.WebhookEvent `json:"event"`
	Data  WebookData          `json:"data"`
}
