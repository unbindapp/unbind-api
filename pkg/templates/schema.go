package templates

// TemplateMetadata represents the metadata of a template
type TemplateMetadata struct {
	Name        string                  `yaml:"name" json:"name"`
	Description string                  `yaml:"description" json:"description"`
	Type        string                  `yaml:"type" json:"type"`
	Version     string                  `yaml:"version" json:"version"`
	Imports     []TemplateImport        `yaml:"imports,omitempty" json:"imports,omitempty"`
	Schema      TemplateParameterSchema `yaml:"schema" json:"schema"`
}

// TemplateImport represents an import of an external schema
type TemplateImport struct {
	Path string `yaml:"path"`
	As   string `yaml:"as"`
}

// TemplateParameterSchema defines the structure for allowed parameters
type TemplateParameterSchema struct {
	Properties map[string]ParameterProperty `yaml:"properties" json:"properties"`
	Required   []string                     `yaml:"required,omitempty" json:"required,omitempty"`
	Imports    map[string]interface{}       `json:"-" yaml:"-"` // For resolved imports
}

// ParameterProperty defines a single parameter's schema
type ParameterProperty struct {
	Type        string                       `yaml:"type" json:"type"`
	Description string                       `yaml:"description,omitempty" json:"description,omitempty"`
	Default     interface{}                  `yaml:"default,omitempty" json:"default,omitempty"`
	Enum        []string                     `yaml:"enum,omitempty" json:"enum,omitempty"`
	Properties  map[string]ParameterProperty `yaml:"properties,omitempty" json:"properties,omitempty"`
	Ref         string                       `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Minimum     *float64                     `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	Maximum     *float64                     `yaml:"maximum,omitempty" json:"maximum,omitempty"`
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
