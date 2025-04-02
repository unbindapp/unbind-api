package enum

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

// Parses detected providers into Provider
func ParseProvider(detectedProviders []string) Provider {
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
	case "ruby":
		return Ruby
	case "rust":
		return Rust
	case "staticfile":
		return Staticfile
	default:
		return UnknownProvider
	}
}

// Provider represents the detected provider
type Provider string

const (
	Node            Provider = "node"
	Deno            Provider = "deno"
	Go              Provider = "go"
	Java            Provider = "java"
	PHP             Provider = "php"
	Python          Provider = "python"
	Ruby            Provider = "ruby"
	Rust            Provider = "rust"
	Staticfile      Provider = "staticfile"
	UnknownProvider Provider = "unknown"
)

var allProviders = []Provider{Node, Deno, Go, Java, PHP, Python, Ruby, Rust, Staticfile, UnknownProvider}

// Values provides list valid values for Enum.
func (Provider) Values() (kinds []string) {
	for _, s := range allProviders {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u Provider) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["Provider"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "Provider")
		schemaRef.Title = "Provider"
		for _, v := range allProviders {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["Provider"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/Provider"}
}
