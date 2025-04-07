package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateRendering(t *testing.T) {
	// Create a sample template similar to the postgres one
	template := &Template{
		Name:        "PostgreSQL Database",
		Category:    TemplateCategoryDatabases,
		Description: "Standard PostgreSQL database using zalando postgres-operator",
		Type:        "postgres-operator",
		Version:     "1.0.0",
		Schema: TemplateParameterSchema{
			Properties: map[string]ParameterProperty{
				"version": {
					Type:        "string",
					Description: "PostgreSQL version",
					Default:     "17",
					Enum:        []string{"14", "15", "16", "17"},
				},
				"replicas": {
					Type:        "integer",
					Description: "Number of replicas",
					Default:     1,
					Minimum:     floatPtr(1),
					Maximum:     floatPtr(5),
				},
				"storage": {
					Type:        "string",
					Description: "Storage size",
					Default:     "1Gi",
				},
				"resources": {
					Type:        "object",
					Description: "Resource requirements",
					Properties: map[string]ParameterProperty{
						"requests": {
							Type: "object",
							Properties: map[string]ParameterProperty{
								"cpu": {
									Type:        "string",
									Description: "CPU request",
									Default:     "100m",
								},
								"memory": {
									Type:        "string",
									Description: "Memory request",
									Default:     "128Mi",
								},
							},
						},
						"limits": {
							Type: "object",
							Properties: map[string]ParameterProperty{
								"cpu": {
									Type:        "string",
									Description: "CPU limit",
									Default:     "200m",
								},
								"memory": {
									Type:        "string",
									Description: "Memory limit",
									Default:     "256Mi",
								},
							},
						},
					},
				},
				"s3": {
					Type: "object",
					Properties: map[string]ParameterProperty{
						"enabled": {
							Type:        "boolean",
							Description: "Enable S3 backups",
							Default:     false,
						},
						"bucket": {
							Type:        "string",
							Description: "S3 bucket name",
						},
						"endpoint": {
							Type:        "string",
							Description: "S3 endpoint URL",
							Default:     "https://s3.amazonaws.com",
						},
					},
				},
				"labels": {
					Type:        "object",
					Description: "Custom labels to add to the PostgreSQL resource",
					AdditionalProperties: &ParameterProperty{
						Type: "string",
					},
				},
			},
			Required: []string{"replicas"},
		},
		Content: `apiVersion: acid.zalan.do/v1
kind: postgresql
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  labels:
    # Standard labels from service
    unbind-team: {{ .TeamID }}
    unbind-project: {{ .ProjectID }}
    unbind-environment: {{ .EnvironmentID }}
    unbind-service: {{ .ServiceID }}
    # Template-specific labels
    unbind/template-name: {{ .Template.Name }}
    unbind/template-version: {{ .Template.Version }}
    {{- range $key, $value := .Parameters.labels }}
    {{ $key }}: {{ $value | quote }}
    {{- end }}
spec:
  teamId: {{ .TeamID | default "default" }}
  postgresql:
    version: {{ .Parameters.version | default "17" }}
  numberOfInstances: {{ .Parameters.replicas | default 1 }}
  volume:
    size: {{ .Parameters.storage | default "1Gi" }}
  resources:
    requests:
      cpu: {{ .Parameters.resources.requests.cpu | default "100m" }}
      memory: {{ .Parameters.resources.requests.memory | default "128Mi" }}
    limits:
      cpu: {{ .Parameters.resources.limits.cpu | default "200m" }}
      memory: {{ .Parameters.resources.limits.memory | default "256Mi" }}
  {{- if .Parameters.s3.enabled }}
  env:
    - name: AWS_ACCESS_KEY_ID
      valueFrom:
        secretKeyRef:
          name: {{ .Name }}-s3-credentials
          key: accessKey
    - name: AWS_ENDPOINT
      value: {{ .Parameters.s3.endpoint }}
    - name: AWS_S3_FORCE_PATH_STYLE
      value: "true"
    - name: BACKUP_NUM_TO_RETAIN
      value: "5"
    - name: BACKUP_SCHEDULE
      value: "5 5 * * *"
    - name: WAL_S3_BUCKET
      value: {{ .Parameters.s3.bucket }}
  {{- end }}`,
	}

	// Create a renderer
	renderer := NewTemplateRenderer()

	t.Run("Basic Rendering", func(t *testing.T) {
		// Create render context with minimal parameters
		ctx := &RenderContext{
			Name:          "test-postgres",
			Namespace:     "default",
			TeamID:        "team1",
			ProjectID:     "project1",
			EnvironmentID: "env1",
			ServiceID:     "svc1",
			Parameters: map[string]interface{}{
				"replicas": 3,
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify the output contains the expected values
		assert.Contains(t, result, "name: test-postgres")
		assert.Contains(t, result, "namespace: default")
		assert.Contains(t, result, "unbind-team: team1")
		assert.Contains(t, result, "version: 17")          // Default value
		assert.Contains(t, result, "numberOfInstances: 3") // Provided value
		assert.Contains(t, result, "size: 1Gi")            // Default value

		// S3 section should not be included since enabled=false by default
		assert.NotContains(t, result, "AWS_ACCESS_KEY_ID")

		// Check for resource defaults
		assert.Contains(t, result, "cpu: 100m")
		assert.Contains(t, result, "memory: 128Mi")
		assert.Contains(t, result, "cpu: 200m")
		assert.Contains(t, result, "memory: 256Mi")

		// Parse to objects
		objects, err := renderer.RenderToObjects(result)
		require.NoError(t, err)
		assert.Len(t, objects, 1) // Should have 1 Kubernetes object
	})

	t.Run("With Custom Labels", func(t *testing.T) {
		// Create render context with custom labels
		ctx := &RenderContext{
			Name:          "test-postgres-labels",
			Namespace:     "default",
			TeamID:        "team1",
			ProjectID:     "project1",
			EnvironmentID: "env1",
			ServiceID:     "svc1",
			Parameters: map[string]interface{}{
				"replicas": 2,
				"labels": map[string]interface{}{
					"environment": "production",
					"app":         "my-app",
					"tier":        "database",
					"cost-center": "123456",
				},
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify custom labels are included
		assert.Contains(t, result, `environment: "production"`)
		assert.Contains(t, result, `app: "my-app"`)
		assert.Contains(t, result, `tier: "database"`)
		assert.Contains(t, result, `cost-center: "123456"`)

		// Make sure standard labels are still there
		assert.Contains(t, result, "unbind-team: team1")
		assert.Contains(t, result, "unbind-project: project1")
		assert.Contains(t, result, "unbind-environment: env1")

		// Ensure labels are properly quoted
		assert.NotContains(t, result, "environment: production") // Should be quoted
	})

	t.Run("With Special Characters in Labels", func(t *testing.T) {
		// Create render context with labels containing special characters
		ctx := &RenderContext{
			Name:          "test-postgres-special",
			Namespace:     "default",
			TeamID:        "team1",
			ProjectID:     "project1",
			EnvironmentID: "env1",
			ServiceID:     "svc1",
			Parameters: map[string]interface{}{
				"replicas": 2,
				"labels": map[string]interface{}{
					"app.kubernetes.io/name":      "postgres",
					"app.kubernetes.io/component": "database",
					"special/label":               "value-with-hyphens",
					"security-level":              "high",
				},
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify custom labels with special characters are included and properly quoted
		assert.Contains(t, result, `app.kubernetes.io/name: "postgres"`)
		assert.Contains(t, result, `app.kubernetes.io/component: "database"`)
		assert.Contains(t, result, `special/label: "value-with-hyphens"`)
		assert.Contains(t, result, `security-level: "high"`)
	})

	t.Run("Empty Labels Map", func(t *testing.T) {
		// Create render context with an empty labels map
		ctx := &RenderContext{
			Name:          "test-postgres-empty-labels",
			Namespace:     "default",
			TeamID:        "team1",
			ProjectID:     "project1",
			EnvironmentID: "env1",
			ServiceID:     "svc1",
			Parameters: map[string]interface{}{
				"replicas": 2,
				"labels":   map[string]interface{}{},
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify the template renders correctly with an empty labels map
		assert.Contains(t, result, "unbind-team: team1")
		assert.NotContains(t, result, "{{ $key }}: {{ $value | quote }}") // Template should not contain raw template syntax
	})

	t.Run("S3 Enabled", func(t *testing.T) {
		// Create render context with S3 enabled
		ctx := &RenderContext{
			Name:          "test-postgres-s3",
			Namespace:     "default",
			TeamID:        "team1",
			ProjectID:     "project1",
			EnvironmentID: "env1",
			ServiceID:     "svc1",
			Parameters: map[string]interface{}{
				"replicas": 2,
				"s3": map[string]interface{}{
					"enabled": true,
					"bucket":  "test-bucket",
				},
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify S3 section is included
		assert.Contains(t, result, "AWS_ACCESS_KEY_ID")
		assert.Contains(t, result, "value: test-bucket")
		assert.Contains(t, result, "value: https://s3.amazonaws.com") // Default endpoint
	})

	t.Run("Custom Values", func(t *testing.T) {
		// Create render context with all custom values
		ctx := &RenderContext{
			Name:          "custom-postgres",
			Namespace:     "custom-ns",
			TeamID:        "custom-team",
			ProjectID:     "custom-project",
			EnvironmentID: "custom-env",
			ServiceID:     "custom-svc",
			Parameters: map[string]interface{}{
				"version":  "15",
				"replicas": 5,
				"storage":  "10Gi",
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "500m",
						"memory": "1Gi",
					},
					"limits": map[string]interface{}{
						"cpu":    "1",
						"memory": "2Gi",
					},
				},
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify the output contains the custom values
		assert.Contains(t, result, "version: 15")
		assert.Contains(t, result, "numberOfInstances: 5")
		assert.Contains(t, result, "size: 10Gi")
		assert.Contains(t, result, "cpu: 500m")
		assert.Contains(t, result, "memory: 1Gi")
		assert.Contains(t, result, "cpu: 1")
		assert.Contains(t, result, "memory: 2Gi")
	})

	t.Run("Custom Values With Labels", func(t *testing.T) {
		// Create render context with custom values and labels
		ctx := &RenderContext{
			Name:          "custom-postgres-labels",
			Namespace:     "custom-ns",
			TeamID:        "custom-team",
			ProjectID:     "custom-project",
			EnvironmentID: "custom-env",
			ServiceID:     "custom-svc",
			Parameters: map[string]interface{}{
				"version":  "15",
				"replicas": 5,
				"storage":  "10Gi",
				"labels": map[string]interface{}{
					"environment": "production",
					"app":         "custom-app",
					"tier":        "database",
				},
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "500m",
						"memory": "1Gi",
					},
					"limits": map[string]interface{}{
						"cpu":    "1",
						"memory": "2Gi",
					},
				},
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify both custom values and labels are included
		assert.Contains(t, result, "version: 15")
		assert.Contains(t, result, "numberOfInstances: 5")
		assert.Contains(t, result, `environment: "production"`)
		assert.Contains(t, result, `app: "custom-app"`)
		assert.Contains(t, result, `tier: "database"`)
	})

	t.Run("Parameter Validation", func(t *testing.T) {
		// Test validation of required fields
		err := renderer.Validate(map[string]interface{}{}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing required field: replicas")

		// Test validation of enum values
		err = renderer.Validate(map[string]interface{}{
			"replicas": 2,
			"version":  "18", // Not in enum
		}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be one of: [14 15 16 17]")

		// Test validation of range
		err = renderer.Validate(map[string]interface{}{
			"replicas": 10, // Above maximum
		}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be <= 5")

		// Test validation of types
		err = renderer.Validate(map[string]interface{}{
			"replicas": "not-a-number",
		}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a number")

		// Test validation with valid labels
		err = renderer.Validate(map[string]interface{}{
			"replicas": 3,
			"labels": map[string]interface{}{
				"app": "test",
			},
		}, template.Schema)
		assert.NoError(t, err)

		// Test validation with invalid labels type
		err = renderer.Validate(map[string]interface{}{
			"replicas": 3,
			"labels":   "not-an-object", // Should be map/object
		}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be an object")

		// Test validation with invalid label value type
		err = renderer.Validate(map[string]interface{}{
			"replicas": 3,
			"labels": map[string]interface{}{
				"app": 123, // Should be string
			},
		}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a string")

		// Valid parameters should pass validation
		err = renderer.Validate(map[string]interface{}{
			"replicas": 3,
			"version":  "16",
			"storage":  "5Gi",
			"labels": map[string]interface{}{
				"app":         "postgres",
				"environment": "staging",
			},
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu": "200m",
				},
			},
		}, template.Schema)
		assert.NoError(t, err)
	})

	t.Run("Default Value Application", func(t *testing.T) {
		// Create render context with minimal parameters
		ctx := &RenderContext{
			Name:      "defaults-test",
			Namespace: "default",
			Parameters: map[string]interface{}{
				"replicas": 2,
			},
		}

		// Apply defaults
		params := renderer.applyDefaults(ctx.Parameters, template.Schema)

		// Check defaults were applied
		assert.Equal(t, "17", params["version"])
		assert.Equal(t, 2, params["replicas"]) // Original value preserved
		assert.Equal(t, "1Gi", params["storage"])

		// Check nested defaults
		resources, ok := params["resources"].(map[string]interface{})
		require.True(t, ok)

		requests, ok := resources["requests"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "100m", requests["cpu"])
		assert.Equal(t, "128Mi", requests["memory"])

		limits, ok := resources["limits"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "200m", limits["cpu"])
		assert.Equal(t, "256Mi", limits["memory"])

		// Check S3 defaults
		s3, ok := params["s3"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, false, s3["enabled"])
		assert.Equal(t, "https://s3.amazonaws.com", s3["endpoint"])

		// Labels should not have defaults since it's additionalProperties
		_, hasLabels := params["labels"]
		assert.False(t, hasLabels, "Labels should not have defaults")
	})
}

// Helper function to create a float pointer
func floatPtr(f float64) *float64 {
	return &f
}
