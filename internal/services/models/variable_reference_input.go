package models

import (
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
)

type CreateVariableReferenceInput struct {
	TargetServiceID uuid.UUID                        `json:"target_service_id" doc:"The ID of the service to which this variable reference belongs" required:"true"`
	TargetName      string                           `json:"target_name" doc:"The name of the target variable" required:"true"`
	Sources         []schema.VariableReferenceSource `json:"sources" doc:"The sources to reference in the template interpolation" nullable:"false"`
	ValueTemplate   string                           `json:"value_template" doc:"The template for the value of the variable reference, e.g. 'https://${sourcename.sourcekey}'" required:"true"`
}

type ResolveVariableReferenceInput struct {
	TeamID     uuid.UUID                          `query:"team_id"`
	Type       schema.VariableReferenceType       `query:"type"`
	Name       string                             `query:"name"`
	SourceType schema.VariableReferenceSourceType `query:"source_type"`
	SourceID   uuid.UUID                          `query:"source_id"`
	Key        string                             `query:"key"`
}
