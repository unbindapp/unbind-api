package webhooks_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

func (self *WebhooksService) sendDiscordWebhook(level WebhookLevel, event schema.WebhookEvent, data WebookData) error {
	// Convert to discord format
	fields := make([]Field, len(data.Fields))
	embed := Embed{
		Title:       &data.Title,
		Url:         &data.Url,
		Description: &data.Description,
		Color:       level.DiscordColor(),
		Footer: &Footer{
			Text: utils.ToPtr(fmt.Sprintf("%s: %s", event, time.Now().Format(time.RFC1123))),
		},
	}
	for i, entry := range data.Fields {
		fields[i] = Field{
			Name:  &entry.Name,
			Value: &entry.Value,
		}
	}
	embed.Fields = &fields

	// Execute
	msg := Message{
		Username: &data.Username,
		Embeds: &[]Embed{
			embed,
		},
	}

	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(msg)
	if err != nil {
		log.Errorf("Failed to encode discord webhook payload: %v", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, data.Url, payload)
	if err != nil {
		log.Errorf("Failed to create discord webhook request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Do
	resp, err := self.httpClient.Do(req)
	if err != nil {
		log.Errorf("Failed to send discord webhook: %v", err)
		return err
	}

	if !slices.Contains([]int{http.StatusOK, http.StatusNoContent}, resp.StatusCode) {
		log.Errorf("Failed to send discord webhook: %v", resp.Status)
		return fmt.Errorf("failed to send discord webhook: %s", resp.Status)
	}

	return nil
}

type Message struct {
	Username        *string          `json:"username,omitempty"`
	AvatarUrl       *string          `json:"avatar_url,omitempty"`
	Content         *string          `json:"content,omitempty"`
	Embeds          *[]Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
}

type Embed struct {
	Title       *string    `json:"title,omitempty"`
	Url         *string    `json:"url,omitempty"`
	Description *string    `json:"description,omitempty"`
	Color       *string    `json:"color,omitempty"`
	Author      *Author    `json:"author,omitempty"`
	Fields      *[]Field   `json:"fields,omitempty"`
	Thumbnail   *Thumbnail `json:"thumbnail,omitempty"`
	Image       *Image     `json:"image,omitempty"`
	Footer      *Footer    `json:"footer,omitempty"`
}

type Author struct {
	Name    *string `json:"name,omitempty"`
	Url     *string `json:"url,omitempty"`
	IconUrl *string `json:"icon_url,omitempty"`
}

type Field struct {
	Name   *string `json:"name,omitempty"`
	Value  *string `json:"value,omitempty"`
	Inline *bool   `json:"inline,omitempty"`
}

type Thumbnail struct {
	Url *string `json:"url,omitempty"`
}

type Image struct {
	Url *string `json:"url,omitempty"`
}

type Footer struct {
	Text    *string `json:"text,omitempty"`
	IconUrl *string `json:"icon_url,omitempty"`
}

type AllowedMentions struct {
	Parse *[]string `json:"parse,omitempty"`
	Users *[]string `json:"users,omitempty"`
	Roles *[]string `json:"roles,omitempty"`
}
