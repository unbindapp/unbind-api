package models

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type SecretType string

const (
	TeamSecret        SecretType = "team"
	ProjectSecret     SecretType = "project"
	EnvironmentSecret SecretType = "environment"
	ServiceSecret     SecretType = "service"
)

var SecretTypeValues = []SecretType{
	TeamSecret,
	ProjectSecret,
	EnvironmentSecret,
	ServiceSecret,
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u SecretType) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["SecretType"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "SecretType")
		schemaRef.Title = "SecretType"
		for _, v := range SecretTypeValues {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["SecretType"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/SecretType"}
}

type SecretResponse struct {
	Type  SecretType `json:"type"`
	Key   string     `json:"key"`
	Value string     `json:"value"`
}

type SecretDeleteInput struct {
	Name string `json:"name" required:"true"`
}
