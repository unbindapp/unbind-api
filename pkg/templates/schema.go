package templates

// TemplateMetadata represents the metadata of a template
type TemplateMetadata struct {
	Name        string                  `yaml:"name"`
	Description string                  `yaml:"description"`
	Type        string                  `yaml:"type"`
	Version     string                  `yaml:"version"`
	Maintainer  string                  `yaml:"maintainer"`
	Imports     []TemplateImport        `yaml:"imports,omitempty"`
	Schema      TemplateParameterSchema `yaml:"schema"`
}

// TemplateImport represents an import of an external schema
type TemplateImport struct {
	Path string `yaml:"path"`
	As   string `yaml:"as"`
}

// TemplateParameterSchema defines the structure for allowed parameters
type TemplateParameterSchema struct {
	Properties map[string]ParameterProperty `yaml:"properties"`
	Required   []string                     `yaml:"required,omitempty"`
	Imports    map[string]interface{}       `json:"-" yaml:"-"` // For resolved imports
}

// ParameterProperty defines a single parameter's schema
type ParameterProperty struct {
	Type        string                       `yaml:"type"`
	Description string                       `yaml:"description,omitempty"`
	Default     interface{}                  `yaml:"default,omitempty"`
	Enum        []string                     `yaml:"enum,omitempty"`
	Properties  map[string]ParameterProperty `yaml:"properties,omitempty"`
	Ref         string                       `yaml:"$ref,omitempty"`
	Minimum     *float64                     `yaml:"minimum,omitempty"`
	Maximum     *float64                     `yaml:"maximum,omitempty"`
}

// Template represents a fully resolved template
type Template struct {
	Name        string                  `json:"name"`
	Category    TemplateCategoryName    `json:"category"`
	Description string                  `json:"description"`
	Type        string                  `json:"type"`
	Version     string                  `json:"version"`
	Schema      TemplateParameterSchema `json:"schema"`
	Content     string                  `json:"-"`
}
