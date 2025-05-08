package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/ent/schema"
)

func TestResolveGeneratedVariables(t *testing.T) {
	// Create a test template with both static and generated variables
	template := &schema.TemplateDefinition{
		Name:        "test-template",
		Description: "Test template",
		Version:     1,
		Services: []schema.TemplateService{
			{
				ID:      1,
				Name:    "TestService",
				Type:    schema.ServiceTypeDockerimage,
				Builder: schema.ServiceBuilderDocker,
				Variables: []schema.TemplateVariable{
					{
						Name:  "STATIC_VAR",
						Value: "static-value",
					},
					{
						Name: "GENERATED_PASSWORD",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
						},
					},
					{
						Name: "GENERATED_EMAIL",
						Generator: &schema.ValueGenerator{
							Type:       schema.GeneratorTypeEmail,
							BaseDomain: "example.com",
						},
					},
				},
			},
		},
	}

	templater := NewTemplater(&config.Config{
		ExternalUIUrl: "https://example.com",
	})

	// Test variable resolution
	resolved, err := templater.ResolveGeneratedVariables(template)
	require.NoError(t, err)
	require.NotNil(t, resolved)

	// Verify template structure is preserved
	assert.Equal(t, template.Name, resolved.Name)
	assert.Equal(t, template.Description, resolved.Description)
	assert.Equal(t, template.Version, resolved.Version)
	assert.Len(t, resolved.Services, 1)

	// Get the service and its variables
	service := resolved.Services[0]
	assert.Equal(t, template.Services[0].Name, service.Name)
	assert.Equal(t, template.Services[0].Type, service.Type)
	assert.Equal(t, template.Services[0].Builder, service.Builder)
	assert.Len(t, service.Variables, 3)

	// Verify static variable is preserved
	staticVar := findVariable(service.Variables, "STATIC_VAR")
	require.NotNil(t, staticVar)
	assert.Equal(t, "static-value", staticVar.Value)
	assert.Nil(t, staticVar.Generator)

	// Verify password is generated
	passwordVar := findVariable(service.Variables, "GENERATED_PASSWORD")
	require.NotNil(t, passwordVar)
	assert.NotEmpty(t, passwordVar.Value)
	assert.Nil(t, passwordVar.Generator)
	assert.Len(t, passwordVar.Value, 16) // Default password length

	// Verify email is generated
	emailVar := findVariable(service.Variables, "GENERATED_EMAIL")
	require.NotNil(t, emailVar)
	assert.NotEmpty(t, emailVar.Value)
	assert.Nil(t, emailVar.Generator)
	assert.Equal(t, "admin@example.com", emailVar.Value)

	// Verify all variables have no generators
	for _, v := range service.Variables {
		assert.Nil(t, v.Generator, "Variable %s should not have a generator", v.Name)
	}
}

// Helper function to find a variable by name
func findVariable(vars []schema.TemplateVariable, name string) *schema.TemplateVariable {
	for _, v := range vars {
		if v.Name == name {
			return &v
		}
	}
	return nil
}
