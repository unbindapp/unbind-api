package models

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type VariableType string

const (
	TeamVariable        VariableType = "team"
	ProjectVariable     VariableType = "project"
	EnvironmentVariable VariableType = "environment"
	ServiceVariable     VariableType = "service"
)

var VariableTypeValues = []VariableType{
	TeamVariable,
	ProjectVariable,
	EnvironmentVariable,
	ServiceVariable,
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u VariableType) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["VariableType"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "VariableType")
		schemaRef.Title = "VariableType"
		for _, v := range VariableTypeValues {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["VariableType"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/VariableType"}
}

type VariableResponse struct {
	Type  VariableType `json:"type"`
	Name  string       `json:"name"`
	Value string       `json:"value"`
}

type VariableDeleteInput struct {
	Name string `json:"name" required:"true"`
}
