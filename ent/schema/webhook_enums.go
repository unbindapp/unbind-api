package schema

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
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
// https://github.com/danielgtaylor/huma/issues/621
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
// https://github.com/danielgtaylor/huma/issues/621
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

// Webhook events
type WebhookEvent string

const (
	WebhookEventProjectCreated      WebhookEvent = "project.created"
	WebhookEventProjectUpdated      WebhookEvent = "project.updated"
	WebhookEventProjectDeleted      WebhookEvent = "project.deleted"
	WebhookEventServiceCreated      WebhookEvent = "service.created"
	WebhookEventServiceUpdated      WebhookEvent = "service.updated"
	WebhookEventServiceDeleted      WebhookEvent = "service.deleted"
	WebhookEventDeploymentQueued    WebhookEvent = "deployment.queued"
	WebhookEventDeploymentBuilding  WebhookEvent = "deployment.building"
	WebhookEventDeploymentSucceeded WebhookEvent = "deployment.succeeded"
	WebhookEventDeploymentFailed    WebhookEvent = "deployment.failed"
	WebhookEventDeploymentCancelled WebhookEvent = "deployment.cancelled"
)

var allWebhookEvents = []WebhookEvent{
	WebhookEventProjectCreated,
	WebhookEventProjectUpdated,
	WebhookEventProjectDeleted,
	WebhookEventServiceCreated,
	WebhookEventServiceUpdated,
	WebhookEventServiceDeleted,
	WebhookEventDeploymentQueued,
	WebhookEventDeploymentBuilding,
	WebhookEventDeploymentSucceeded,
	WebhookEventDeploymentFailed,
	WebhookEventDeploymentCancelled,
}

// Values provides list valid values for Enum.
func (s WebhookEvent) Values() (kinds []string) {
	for _, s := range allWebhookEvents {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u WebhookEvent) Schema(r huma.Registry) *huma.Schema {
	// Register WebhookTeamEvents schema
	if r.Map()["WebhookTeamEvents"] == nil {
		// Create a schema for team events
		teamSchema := &huma.Schema{
			Title: "WebhookTeamEvents",
			Type:  "object",
			Properties: map[string]*huma.Schema{
				"event": {
					Type: "string",
					Enum: []interface{}{
						string(WebhookEventProjectCreated),
						string(WebhookEventProjectUpdated),
						string(WebhookEventProjectDeleted),
					},
				},
				// Add other team event properties here
			},
			Required: []string{"event"},
		}
		r.Map()["WebhookTeamEvents"] = teamSchema
	}

	// Register WebhookProjectEvents schema
	if r.Map()["WebhookProjectEvents"] == nil {
		// Create a schema for project events (non-team events)
		projectEvents := []interface{}{
			string(WebhookEventServiceCreated),
			string(WebhookEventServiceUpdated),
			string(WebhookEventServiceDeleted),
			string(WebhookEventDeploymentQueued),
			string(WebhookEventDeploymentBuilding),
			string(WebhookEventDeploymentSucceeded),
			string(WebhookEventDeploymentFailed),
			string(WebhookEventDeploymentCancelled),
		}

		projectSchema := &huma.Schema{
			Title: "WebhookProjectEvents",
			Type:  "object",
			Properties: map[string]*huma.Schema{
				"event": {
					Type: "string",
					Enum: projectEvents,
				},
				// Add other project event properties here
			},
			Required: []string{"event"},
		}
		r.Map()["WebhookProjectEvents"] = projectSchema
	}

	// Register WebhookEvent as a oneOf schema
	if r.Map()["WebhookEvent"] == nil {
		webhookEventsSchema := &huma.Schema{
			Title: "WebhookEvent",
			OneOf: []*huma.Schema{
				{Ref: "#/components/schemas/WebhookTeamEvents"},
				{Ref: "#/components/schemas/WebhookProjectEvents"},
			},
		}
		r.Map()["WebhookEvent"] = webhookEventsSchema
	}

	// Return reference to WebhookEvent schema
	return &huma.Schema{Ref: "#/components/schemas/WebhookEvent"}
}
