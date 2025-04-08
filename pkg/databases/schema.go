package databases

// DefinitionMetadata represents the metadata of a definition
type DefinitionMetadata struct {
	Name        string                    `yaml:"name" json:"name"`
	Port        int                       `yaml:"port" json:"port"`
	Description string                    `yaml:"description" json:"description"`
	Type        string                    `yaml:"type" json:"type"`
	Version     string                    `yaml:"version" json:"version"`
	Imports     []DefinitionImport        `yaml:"imports,omitempty" json:"imports,omitempty"`
	Schema      DefinitionParameterSchema `yaml:"schema" json:"schema"`
}

// DefinitionImport represents an import of an external schema
type DefinitionImport struct {
	Path string `yaml:"path"`
	As   string `yaml:"as"`
}

// DefinitionParameterSchema defines the structure for allowed parameters
type DefinitionParameterSchema struct {
	Properties map[string]ParameterProperty `yaml:"properties" json:"properties"`
	Required   []string                     `yaml:"required,omitempty" json:"required,omitempty"`
	Imports    map[string]interface{}       `json:"-" yaml:"-"` // For resolved imports
}

// ParameterProperty defines a single parameter's schema
type ParameterProperty struct {
	Type                 string                       `yaml:"type" json:"type"`
	Description          string                       `yaml:"description,omitempty" json:"description,omitempty"`
	Default              interface{}                  `yaml:"default,omitempty" json:"default,omitempty"`
	Enum                 []string                     `yaml:"enum,omitempty" json:"enum,omitempty"`
	Properties           map[string]ParameterProperty `yaml:"properties,omitempty" json:"properties,omitempty"`
	AdditionalProperties *ParameterProperty           `yaml:"additionalProperties,omitempty" json:"additionalProperties,omitempty"`
	Ref                  string                       `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Minimum              *float64                     `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	Maximum              *float64                     `yaml:"maximum,omitempty" json:"maximum,omitempty"`
}

// Definition represents a fully resolved definition
type Definition struct {
	Name        string                    `json:"name"`
	Port        int                       `json:"port"`
	Category    string                    `json:"category"`
	Description string                    `json:"description"`
	Type        string                    `json:"type"`
	Version     string                    `json:"version"`
	Schema      DefinitionParameterSchema `json:"schema"`
	Content     string                    `json:"-"`
}
