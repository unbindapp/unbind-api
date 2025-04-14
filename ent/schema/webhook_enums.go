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
	if r.Map()["WebhookEvent"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "WebhookEvent")
		schemaRef.Title = "WebhookEvent"
		for _, v := range allWebhookEvents {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["WebhookEvent"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/WebhookEvent"}
}
