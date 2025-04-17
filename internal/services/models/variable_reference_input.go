package models

import (
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type CreateVariableReferenceInput struct {
	TargetServiceID uuid.UUID                          `json:"target_service_id" doc:"The ID of the service to which this variable reference belongs" required:"true"`
	TargetName      string                             `json:"target_name" doc:"The name of the target variable" required:"true"`
	Type            schema.VariableReferenceType       `json:"type" doc:"The type of variable reference" required:"true"`
	SourceType      schema.VariableReferenceSourceType `json:"source_type" doc:"The source type of the variable reference" required:"true"`
	SourceID        uuid.UUID                          `json:"source_id" doc:"The ID of the source of the variable reference" required:"true"`
	SourceName      string                             `json:"source_name" doc:"The name of the source of the variable reference" required:"true"`
	SourceKey       string                             `json:"source_key" doc:"The key of the source of the variable reference" required:"false"`
	ValueTemplate   *string                            `json:"value_template" doc:"The template for the value of the variable reference, e.g. 'https://${}'" required:"false"`
}
