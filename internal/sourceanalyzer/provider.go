package sourceanalyzer

import (
	"encoding/json"
)

// Parses detected providers into Provider
func parseProvider(detectedProviders []string) Provider {
	provider := ""
	// Always use the first one?
	if len(detectedProviders) > 0 {
		provider = detectedProviders[0]
	}

	switch provider {
	case "node":
		return Node
	case "deno":
		return Deno
	case "golang":
		return Go
	case "java":
		return Java
	case "php":
		return PHP
	case "python":
		return Python
	case "staticfile":
		return Staticfile
	default:
		return UnknownProvider
	}
}

// Provider represents the detected provider
type Provider int

const (
	Node Provider = iota
	Deno
	Go
	Java
	PHP
	Python
	Staticfile
	UnknownProvider
)

func (p Provider) String() string {
	names := map[Provider]string{
		Node:            "node",
		Deno:            "deno",
		Go:              "go",
		Java:            "java",
		PHP:             "php",
		Python:          "python",
		Staticfile:      "staticfile",
		UnknownProvider: "unknown",
	}

	if name, ok := names[p]; ok {
		return name
	}
	return "unknown"
}

func (p Provider) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}
