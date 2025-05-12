package schema

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"reflect"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/ent/schema/mixin"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// Template holds the schema definition for the Template entity.
type Template struct {
	ent.Schema
}

// Mixin of the Template.
func (Template) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.PKMixin{},
		mixin.TimeMixin{},
	}
}

// Fields of the Template.
func (Template) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Int("version"),
		field.Bool("immutable").Default(false).Comment("If true, the template cannot be modified or deleted (system bundle)"),
		field.JSON("definition", TemplateDefinition{}),
	}
}

// Edges of the Template.
func (Template) Edges() []ent.Edge {
	return []ent.Edge{
		// O2M with service
		edge.To("services", Service.Type),
	}
}

// Indexes of the Template.
func (Template) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "version").Unique(),
	}
}

// Annotations of the Template
func (Template) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "templates",
		},
	}
}

// TemplateDefinition represents a complete template configuration
type TemplateDefinition struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     int               `json:"version"`
	Services    []TemplateService `json:"services" nullable:"false"`
	Inputs      []TemplateInput   `json:"inputs" nullable:"false"`
}

// TemplateService represents a service within a template
type TemplateService struct {
	ID                 int                         `json:"id"`
	DependsOn          []int                       `json:"depends_on" nullable:"false"` // IDs of services that must be started before this one
	Icon               string                      `json:"icon,omitempty"`              // Icon name
	Name               string                      `json:"name"`
	Type               ServiceType                 `json:"type"`
	Builder            ServiceBuilder              `json:"builder"`
	DatabaseType       *string                     `json:"database_type,omitempty"`
	DatabaseConfig     *DatabaseConfig             `json:"database_config,omitempty"` // Database configuration
	Image              *string                     `json:"image,omitempty"`
	Ports              []PortSpec                  `json:"ports" nullable:"false"` // Ports to expose
	IsPublic           bool                        `json:"is_public"`
	HostInputIDs       []int                       `json:"host_input_ids,omitempty"` // IDs of inputs that are hostnames
	RunCommand         *string                     `json:"run_command,omitempty"`
	Volumes            []TemplateVolume            `json:"volumes" nullable:"false"`             // Volumes to mount
	Variables          []TemplateVariable          `json:"variables" nullable:"false"`           // Variables this service needs
	VariableReferences []TemplateVariableReference `json:"variable_references" nullable:"false"` // Variables this service needs
	SecurityContext    *SecurityContext            `json:"security_context,omitempty"`           // Security context for the service
}

// TemplateVariable represents a configurable variable in a template
type TemplateVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	// If set, the value will be generated using the generator
	Generator *ValueGenerator `json:"generator,omitempty"`
}

// TenokateVariableReference represents a reference to a variable in a template
type TemplateVariableReference struct {
	SourceID   int    `json:"source_id"`
	TargetName string `json:"target_name"` // Name of the variable
	SourceName string `json:"source_name"` // Name of the variable
	IsHost     bool   `json:"is_host"`     // If true, variable will be <kubernetesName>.<serviceName>, sort of customized by type (e.g. mysql adds moco- prefix)
}

// Types of generators
type GeneratorType string

const (
	GeneratorTypePassword GeneratorType = "password"
	GeneratorTypeInput    GeneratorType = "input"
)

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u GeneratorType) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["GeneratorType"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "GeneratorType")
		schemaRef.Title = "GeneratorType"
		schemaRef.Enum = append(schemaRef.Enum,
			[]any{
				string(GeneratorTypePassword),
				string(GeneratorTypeInput),
			}...)
		r.Map()["GeneratorType"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/GeneratorType"}
}

// ValueGenerator represents how to generate a value
type ValueGenerator struct {
	Type       GeneratorType  `json:"type"`
	InputID    int            `json:"input_id,omitempty"`    // For input
	BaseDomain string         `json:"base_domain,omitempty"` // For email
	AddPrefix  string         `json:"add_prefix,omitempty"`  // Add a prefix to the generated value
	HashType   *ValueHashType `json:"hash_type,omitempty"`   // Hash the generated value
}

type ValueHashType string

const (
	ValueHashTypeSHA256 ValueHashType = "sha256"
	ValueHashTypeSHA512 ValueHashType = "sha512"
)

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u ValueHashType) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["ValueHashType"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "ValueHashType")
		schemaRef.Title = "ValueHashType"
		schemaRef.Enum = append(schemaRef.Enum,
			[]any{
				string(ValueHashTypeSHA256),
				string(ValueHashTypeSHA512),
			}...)
		r.Map()["ValueHashType"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/ValueHashType"}
}

func (self *ValueGenerator) Generate(inputs map[int]string) (string, error) {
	switch self.Type {
	case GeneratorTypePassword:
		pwd, err := utils.GenerateSecurePassword(32)
		if err != nil {
			return "", err
		}
		if self.HashType != nil {
			switch *self.HashType {
			case ValueHashTypeSHA256:
				hash := sha256.Sum256([]byte(pwd))
				pwd = hex.EncodeToString(hash[:])
			case ValueHashTypeSHA512:
				hash := sha512.Sum512([]byte(pwd))
				pwd = hex.EncodeToString(hash[:])
			}
		}
		return self.AddPrefix + pwd, nil
	case GeneratorTypeInput:
		// Find the input by ID
		inputValue, ok := inputs[self.InputID]
		if !ok {
			return "", fmt.Errorf("input ID %d not found in inputs map", self.InputID)
		}
		return self.AddPrefix + inputValue, nil
	default:
		return "", fmt.Errorf("unknown generator type: %s", self.Type)
	}
}

// TemplateInputType represents the type of user input
type TemplateInputType string

const (
	InputTypeVariable   TemplateInputType = "variable"
	InputTypeHost       TemplateInputType = "host"
	InputTypeVolumeSize TemplateInputType = "volume_size"
	InputTypeNodePort   TemplateInputType = "node_port"
)

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u TemplateInputType) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["TemplateInputType"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "TemplateInputType")
		schemaRef.Title = "TemplateInputType"
		schemaRef.Enum = append(schemaRef.Enum,
			[]any{
				string(InputTypeVariable),
				string(InputTypeHost),
				string(InputTypeVolumeSize),
				string(InputTypeNodePort),
			}...)
		r.Map()["TemplateInputType"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/TemplateInputType"}
}

// TemplateInput represents a user input field in the template
type TemplateInput struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Type        TemplateInputType `json:"type"`
	Description string            `json:"description"`
	Default     *string           `json:"default,omitempty"`
	Required    bool              `json:"required"`
	TargetPort  *int              `json:"target_port,omitempty"`
}

// TemplateVolume represents a volume configuration in the template
type TemplateVolume struct {
	Name      string             `json:"name"`
	Size      TemplateVolumeSize `json:"size"`
	MountPath string             `json:"mountPath"`
}

type TemplateVolumeSize struct {
	FromInputID int `json:"from_input_id"`
}
