package templates

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	postgresv1 "github.com/zalando/postgres-operator/pkg/apis/acid.zalan.do/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	yamlDecoder "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

// TemplateRenderer renders templates into Kubernetes resources
type TemplateRenderer struct {
	// Use a custom scheme that can be extended with CRDs
	scheme *runtime.Scheme
}

func NewTemplateRenderer() *TemplateRenderer {
	// Create a new scheme that includes the core Kubernetes types
	s := runtime.NewScheme()
	scheme.AddToScheme(s)

	// Add CRD types to the scheme
	apiextensionsv1.AddToScheme(s)

	r := &TemplateRenderer{
		scheme: s,
	}

	// Register known CRDs
	zalandoGVK := schema.GroupVersionKind{
		Group:   "acid.zalan.do",
		Version: "v1",
		Kind:    "postgresql",
	}
	r.RegisterCRD(zalandoGVK, &postgresv1.Postgresql{})

	return r
}

// RenderContext holds the data for template rendering
type RenderContext struct {
	// Service info
	Name          string
	Namespace     string
	TeamID        string
	ProjectID     string
	EnvironmentID string
	ServiceID     string

	// Template info
	Template struct {
		Name    string
		Version string
		Type    string
	}

	// Parameters from service
	Parameters map[string]interface{}
}

// Render renders a template to YAML string that can be applied to Kubernetes
func (r *TemplateRenderer) Render(unbindTemplate *Template, context *RenderContext) (string, error) {
	// Apply defaults to parameters
	context.Parameters = r.applyDefaults(context.Parameters, unbindTemplate.Schema)

	// Set template info
	context.Template.Name = unbindTemplate.Name
	context.Template.Version = unbindTemplate.Version
	context.Template.Type = unbindTemplate.Type

	// Process template
	var buf bytes.Buffer
	tmpl, err := template.New("template").Funcs(sprig.FuncMap()).Parse(unbindTemplate.Content)
	if err != nil {
		return "", fmt.Errorf("template parsing error: %w", err)
	}

	if err := tmpl.Execute(&buf, context); err != nil {
		return "", fmt.Errorf("template execution error: %w", err)
	}

	return buf.String(), nil
}

// RenderToObjects parses the rendered YAML into Kubernetes objects
func (r *TemplateRenderer) RenderToObjects(renderedYAML string) ([]runtime.Object, error) {
	var resources []runtime.Object
	decoder := serializer.NewCodecFactory(r.scheme).UniversalDeserializer()

	// Split YAML documents
	docs := strings.Split(renderedYAML, "---")
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		// First try to decode using scheme-aware decoder
		obj, _, err := decoder.Decode([]byte(doc), nil, nil)
		if err == nil {
			resources = append(resources, obj)
			continue
		}

		// Fall back to unstructured decoding
		// If we haven't called RegisterCRD - this is where we end up
		unstObj, err := r.decodeToUnstructured(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to decode resource: %w", err)
		}

		resources = append(resources, unstObj)
	}

	return resources, nil
}

// decodeToUnstructured decodes a YAML document to an unstructured.Unstructured object
func (r *TemplateRenderer) decodeToUnstructured(doc string) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}

	reader := strings.NewReader(doc)
	decoder := yamlDecoder.NewYAMLOrJSONDecoder(reader, 4096)

	if err := decoder.Decode(obj); err != nil {
		return nil, err
	}

	// Verify it's a valid Kubernetes resource with apiVersion and kind
	if obj.GetAPIVersion() == "" || obj.GetKind() == "" {
		return nil, fmt.Errorf("resource is missing apiVersion and/or kind")
	}

	return obj, nil
}

// RegisterCRD registers a custom resource definition with the renderer
// For custom CRDs it enabled better validation
func (r *TemplateRenderer) RegisterCRD(gvk schema.GroupVersionKind, obj runtime.Object) {
	// Add the type to the scheme
	r.scheme.AddKnownTypeWithName(gvk, obj)

	// Add the list type if appropriate
	r.scheme.AddKnownTypeWithName(
		schema.GroupVersionKind{Group: gvk.Group, Version: gvk.Version, Kind: gvk.Kind + "List"},
		&unstructured.UnstructuredList{})
}

// Validate validates parameters against the schema
func (r *TemplateRenderer) Validate(params map[string]interface{}, schema TemplateParameterSchema) error {
	// Check required fields
	for _, required := range schema.Required {
		if _, ok := params[required]; !ok {
			return fmt.Errorf("missing required field: %s", required)
		}
	}

	// Validate each field
	for name, value := range params {
		prop, ok := schema.Properties[name]
		if !ok {
			return fmt.Errorf("unknown field: %s", name)
		}

		if err := r.validateProperty(name, value, prop); err != nil {
			return err
		}
	}

	return nil
}

// validateProperty validates a property against its schema
func (r *TemplateRenderer) validateProperty(name string, value interface{}, prop ParameterProperty) error {
	// Skip validation if value is nil
	if value == nil {
		return nil
	}

	// Type checking
	switch prop.Type {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("field %s must be a string", name)
		}

		// Enum validation
		if len(prop.Enum) > 0 {
			strValue := value.(string)
			valid := false
			for _, enum := range prop.Enum {
				if strValue == enum {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("field %s must be one of: %v", name, prop.Enum)
			}
		}
	case "integer", "number":
		// Convert to float for comparison
		var floatVal float64
		switch v := value.(type) {
		case int:
			floatVal = float64(v)
		case int32:
			floatVal = float64(v)
		case int64:
			floatVal = float64(v)
		case float32:
			floatVal = float64(v)
		case float64:
			floatVal = v
		default:
			return fmt.Errorf("field %s must be a number", name)
		}

		// Range validation
		if prop.Minimum != nil && floatVal < *prop.Minimum {
			return fmt.Errorf("field %s must be >= %v", name, *prop.Minimum)
		}
		if prop.Maximum != nil && floatVal > *prop.Maximum {
			return fmt.Errorf("field %s must be <= %v", name, *prop.Maximum)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("field %s must be a boolean", name)
		}
	case "object":
		objValue, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("field %s must be an object", name)
		}

		// Validate nested properties
		if prop.Properties != nil {
			for propName, propValue := range objValue {
				nestedProp, ok := prop.Properties[propName]
				if !ok {
					return fmt.Errorf("unknown field: %s.%s", name, propName)
				}

				if err := r.validateProperty(fmt.Sprintf("%s.%s", name, propName), propValue, nestedProp); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// applyDefaults applies default values from the schema
func (r *TemplateRenderer) applyDefaults(params map[string]interface{}, schema TemplateParameterSchema) map[string]interface{} {
	result := make(map[string]interface{})

	// Copy all provided parameters
	for k, v := range params {
		result[k] = v
	}

	// Apply defaults recursively for all properties in the schema
	r.applyDefaultsRecursive(result, schema.Properties, "")

	return result
}

// applyDefaultsRecursive applies defaults recursively
func (r *TemplateRenderer) applyDefaultsRecursive(params map[string]interface{}, properties map[string]ParameterProperty, prefix string) {
	for name, prop := range properties {
		fullName := name
		if prefix != "" {
			fullName = prefix + "." + name
		}

		parts := strings.Split(fullName, ".")
		paramMap := params

		// Navigate to the correct nested map for deeply nested properties
		for i := 0; i < len(parts)-1; i++ {
			part := parts[i]

			// Create path if it doesn't exist
			if _, exists := paramMap[part]; !exists {
				paramMap[part] = make(map[string]interface{})
			}

			var ok bool
			paramMap, ok = paramMap[part].(map[string]interface{})
			if !ok {
				// If it's not a map, we can't proceed further with this path
				break
			}
		}

		lastPart := parts[len(parts)-1]

		// Apply default for this property if not set
		if _, exists := paramMap[lastPart]; !exists && prop.Default != nil {
			paramMap[lastPart] = prop.Default
		}

		// Recursively apply defaults to nested objects
		if prop.Type == "object" && prop.Properties != nil {
			// Ensure the object exists
			if _, exists := paramMap[lastPart]; !exists {
				paramMap[lastPart] = make(map[string]interface{})
			}

			// Get the nested object
			if nestedMap, ok := paramMap[lastPart].(map[string]interface{}); ok {
				r.applyDefaultsRecursive(nestedMap, prop.Properties, "")
			}
		}
	}
}
