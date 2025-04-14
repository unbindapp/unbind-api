package models

import (
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type WebhookCreateInput struct {
	Type      schema.WebhookType    `json:"type" validate:"required" required:"true"`
	TeamID    uuid.UUID             `json:"team_id" validate:"required,uuid4" required:"true"`
	ProjectID *uuid.UUID            `json:"project_id,omitempty" validate:"omitempty,uuid4" required:"false" doc:"required if type is project"`
	URL       string                `json:"url" validate:"required,url" required:"true"`
	Events    []schema.WebhookEvent `json:"events" validate:"required" required:"true"`
}

type WebhookUpdateInput struct {
	ID        uuid.UUID              `json:"id" validate:"required,uuid4" required:"true"`
	TeamID    uuid.UUID              `json:"team_id" validate:"required,uuid4" required:"true"`
	ProjectID *uuid.UUID             `json:"project_id,omitempty" validate:"omitempty,uuid4" required:"false" doc:"required if type is project"`
	URL       *string                `json:"url" validate:"omitempty,url" required:"false"`
	Events    *[]schema.WebhookEvent `json:"events" required:"false"`
}

type WebhookListInput struct {
	Type      schema.WebhookType `query:"type" validate:"required" required:"true"`
	TeamID    uuid.UUID          `query:"team_id" validate:"required,uuid4" required:"true"`
	ProjectID uuid.UUID          `query:"project_id" validate:"omitempty,uuid4" required:"false"`
}

type WebhookGetInput struct {
	ID        uuid.UUID `query:"id" validate:"required,uuid4" required:"true"`
	TeamID    uuid.UUID `query:"team_id" validate:"required,uuid4" required:"true"`
	ProjectID uuid.UUID `query:"project_id" validate:"omitempty,uuid4" required:"false"`
}
