package models

import "encoding/json"

type SecretType int

const (
	TeamSecret SecretType = iota
	ProjectSecret
	EnvironmentSecret
	ServiceSecret
)

func (f SecretType) String() string {
	switch f {
	case TeamSecret:
		return "team"
	case ProjectSecret:
		return "project"
	case EnvironmentSecret:
		return "environment"
	case ServiceSecret:
		return "service"
	}
	return "unknown"
}

func (f SecretType) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

type SecretResponse struct {
	Type  SecretType `json:"type"`
	Key   string     `json:"key"`
	Value string     `json:"value"`
}
