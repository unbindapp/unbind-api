package schema

import (
	"github.com/danielgtaylor/huma/v2"
	"reflect"
)

// Enums
type WebhookTarget string

const (
	WebhookTargetDiscord WebhookTarget = "discord"
	WebhookTargetSlack   WebhookTarget = "slack"
	WebhookTargetOther   WebhookTarget = "other"
)

var allWebhookTargets = []WebhookTarget{
	WebhookTargetDiscord,
	WebhookTargetSlack,
	WebhookTargetOther,
}

// Values provides list valid values for Enum.
func (s WebhookTarget) Values() (kinds []string) {
	for _, s := range allWebhookTargets {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
func (u WebhookTarget) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["WebhookTarget"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "WebhookTarget")
		schemaRef.Title = "WebhookTarget"
		for _, v := range allWebhookTargets {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["WebhookTarget"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/WebhookTarget"}
}

// Type of webhook
type WebhookType string

const (
	WebhookTypeTeam    WebhookType = "team"
	WebhookTypeProject WebhookType = "project"
)

var allWebhookTypes = []WebhookType{
	WebhookTypeTeam,
	WebhookTypeProject,
}

// Values provides list valid values for Enum.
func (s WebhookType) Values() (kinds []string) {
	for _, s := range allWebhookTypes {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
func (u WebhookType) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["WebhookType"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "WebhookType")
		schemaRef.Title = "WebhookType"
		for _, v := range allWebhookTypes {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["WebhookType"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/WebhookType"}
}

// Project webhook events
type WebhookProjectEvent string

const (
	WebhookProjectEventServiceCreated      WebhookProjectEvent = "service.created"
	WebhookProjectEventServiceUpdated      WebhookProjectEvent = "service.updated"
	WebhookProjectEventServiceDeleted      WebhookProjectEvent = "service.deleted"
	WebhookProjectEventDeploymentQueued    WebhookProjectEvent = "deployment.queued"
	WebhookProjectEventDeploymentBuilding  WebhookProjectEvent = "deployment.building"
	WebhookProjectEventDeploymentSucceeded WebhookProjectEvent = "deployment.succeeded"
	WebhookProjectEventDeploymentFailed    WebhookProjectEvent = "deployment.failed"
	WebhookProjectEventDeploymentCancelled WebhookProjectEvent = "deployment.cancelled"
)

var allWebhookProjectEvents = []WebhookProjectEvent{
	WebhookProjectEventServiceCreated,
	WebhookProjectEventServiceUpdated,
	WebhookProjectEventServiceDeleted,
	WebhookProjectEventDeploymentQueued,
	WebhookProjectEventDeploymentBuilding,
	WebhookProjectEventDeploymentSucceeded,
	WebhookProjectEventDeploymentFailed,
	WebhookProjectEventDeploymentCancelled,
}

// Values provides list valid values for Enum.
func (s WebhookProjectEvent) Values() (kinds []string) {
	for _, s := range allWebhookProjectEvents {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
func (u WebhookProjectEvent) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["WebhookProjectEvent"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "WebhookProjectEvent")
		schemaRef.Title = "WebhookProjectEvent"
		for _, v := range allWebhookProjectEvents {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["WebhookProjectEvent"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/WebhookProjectEvent"}
}

// Team webhook events
type WebhookTeamEvent string

const (
	WebhookTeamEventProjectCreated WebhookTeamEvent = "project.created"
	WebhookTeamEventProjectUpdated WebhookTeamEvent = "project.updated"
	WebhookTeamEventProjectDeleted WebhookTeamEvent = "project.deleted"
)

var allWebhookTeamEvents = []WebhookTeamEvent{
	WebhookTeamEventProjectCreated,
	WebhookTeamEventProjectUpdated,
	WebhookTeamEventProjectDeleted,
}

// Values provides list valid values for Enum.
func (s WebhookTeamEvent) Values() (kinds []string) {
	for _, s := range allWebhookTeamEvents {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
func (u WebhookTeamEvent) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["WebhookTeamEvent"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "WebhookTeamEvent")
		schemaRef.Title = "WebhookTeamEvent"
		for _, v := range allWebhookTeamEvents {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["WebhookTeamEvent"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/WebhookTeamEvent"}
}

// WebhookEventsOneOf represents either project events or team events
type WebhookEventsOneOf struct {
	// This field is just a placeholder for Go's type system
	// The actual schema will be defined in the Schema method
	Events []string `json:"events"`
}

// Schema creates a oneOf schema for webhook events
func (w WebhookEventsOneOf) Schema(r huma.Registry) *huma.Schema {
	projectEventArraySchema := &huma.Schema{
		Type: "array",
		Items: &huma.Schema{
			Ref: "#/components/schemas/WebhookProjectEvent",
		},
	}

	teamEventArraySchema := &huma.Schema{
		Type: "array",
		Items: &huma.Schema{
			Ref: "#/components/schemas/WebhookTeamEvent",
		},
	}

	return &huma.Schema{
		Type: "object",
		Properties: map[string]*huma.Schema{
			"events": {
				OneOf: []*huma.Schema{
					projectEventArraySchema,
					teamEventArraySchema,
				},
			},
		},
		Required: []string{"events"},
	}
}
