package models

import (
	"github.com/unbindapp/unbind-api/ent/schema"
)

type VariableReferenceInputItem struct {
	TargetName    string                           `json:"target_name" doc:"The name of the target variable" required:"true"`
	Sources       []schema.VariableReferenceSource `json:"sources" doc:"The sources to reference in the template interpolation" nullable:"false"`
	ValueTemplate string                           `json:"value_template" doc:"The template for the value of the variable reference, e.g. 'https://${sourcename.sourcekey}'" required:"true"`
}
