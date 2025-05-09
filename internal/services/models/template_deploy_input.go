package models

import "github.com/google/uuid"

// For template inputs specifically
type TemplateInput struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

type TemplateDeployInput struct {
	TemplateID    uuid.UUID       `json:"template_id" format:"uuid" required:"true"`
	TeamID        uuid.UUID       `json:"team_id" format:"uuid" required:"true"`
	ProjectID     uuid.UUID       `json:"project_id" format:"uuid" required:"true"`
	EnvironmentID uuid.UUID       `json:"environment_id" format:"uuid" required:"true"`
	Inputs        []TemplateInput `json:"inputs,omitempty" required:"false"`
}
