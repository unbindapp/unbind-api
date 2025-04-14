package webhooks_service

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
)

func (self *WebhooksService) DetectTargetFromURL(urlStr string) (schema.WebhookTarget, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return schema.WebhookTargetOther, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("Invalid webhook URL: %s", urlStr))
	}

	if parsedURL.Host == "discord.com" || parsedURL.Host == "discordapp.com" {
		// Check if the path matches the webhook pattern
		if strings.HasPrefix(parsedURL.Path, "/api/webhooks/") {
			return schema.WebhookTargetDiscord, nil
		}
	}

	// Regular Slack webhooks
	if strings.HasSuffix(parsedURL.Host, ".slack.com") &&
		strings.Contains(parsedURL.Path, "/services/") {
		return schema.WebhookTargetSlack, nil
	}

	// Slack API webhooks
	if parsedURL.Host == "hooks.slack.com" &&
		strings.HasPrefix(parsedURL.Path, "/services/") {
		return schema.WebhookTargetSlack, nil
	}

	// Match Slack webhook URL patterns
	// Example: https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX
	slackRegex := regexp.MustCompile(`^https://[^/]+/services/[A-Z0-9]+/[A-Z0-9]+/[A-Za-z0-9]+$`)
	if slackRegex.MatchString(parsedURL.String()) {
		return schema.WebhookTargetSlack, nil
	}

	return schema.WebhookTargetOther, nil
}
