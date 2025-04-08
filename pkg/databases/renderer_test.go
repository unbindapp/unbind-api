package databases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefinitionRendering(t *testing.T) {
	// Create a sample template similar to the postgres one
	template := &Definition{
		Name:        "PostgreSQL Database",
		Category:    DB_CATEGORY,
		Description: "Standard PostgreSQL database using zalando postgres-operator",
		Type:        "postgres-operator",
		Version:     "1.0.0",
		Schema: DefinitionParameterSchema{
			Properties: map[string]ParameterProperty{
				"common": {
					Type: "object",
					Properties: map[string]ParameterProperty{
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
					},
				},
				"dockerImage": {
					Type:        "string",
					Description: "Spilo image version",
					Default:     "unbindapp/spilo:17",
				},
				"version": {
					Type:        "string",
					Description: "PostgreSQL version",
					Default:     "17",
					Enum:        []string{"14", "15", "16", "17"},
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
						"region": {
							Type:        "string",
							Description: "AWS region",
							Default:     "us-east-1",
						},
						"endpoint": {
							Type:        "string",
							Description: "S3 endpoint URL",
							Default:     "https://s3.amazonaws.com",
						},
						"forcePathStyle": {
							Type:        "boolean",
							Description: "Force path style URLs for S3 objects",
							Default:     true,
						},
						"backupRetention": {
							Type:        "integer",
							Description: "Number of backups to retain",
							Default:     5,
						},
						"backupSchedule": {
							Type:        "string",
							Description: "Cron schedule for backups",
							Default:     "5 5 * * *",
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
			Required: []string{"common"},
		},
		Content: `apiVersion: acid.zalan.do/v1
kind: postgresql
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  labels:
    # Zalando labels
    team: {{ .TeamID }}
    # usd-specific labels
    unbind/usd-type: {{ .Definition.Type }}
    unbind/usd-version: {{ .Definition.Version }}
    unbind/usd-category: databases
    {{- range $key, $value := .Parameters.labels }}
    {{ $key }}: {{ $value | quote }}
    {{- end }}
spec:
  teamId: {{ .TeamID }}
  dockerImage: {{ .Parameters.dockerImage | default "unbindapp/spilo:17" }}
  postgresql:
    version: {{ .Parameters.version | default "17" | quote }}
  numberOfInstances: {{ .Parameters.common.replicas | default 1 }}
  allowedSourceRanges:
    - 0.0.0.0/0
  volume:
    size: {{ .Parameters.common.storage | default "1Gi" }}
  resources:
    requests:
      cpu: {{ .Parameters.common.resources.requests.cpu | default "100m" }}
      memory: {{ .Parameters.common.resources.requests.memory | default "128Mi" }}
    limits:
      cpu: {{ .Parameters.common.resources.limits.cpu | default "200m" }}
      memory: {{ .Parameters.common.resources.limits.memory | default "256Mi" }}
  env:
    - name: ALLOW_NOSSL
      value: "true"
    {{- if .Parameters.s3.enabled }}
    - name: AWS_ACCESS_KEY_ID
      valueFrom:
        secretKeyRef:
          name: {{ .Name }}-s3-credentials
          key: accessKey
    - name: AWS_ENDPOINT
      value: {{ .Parameters.s3.endpoint }}
    - name: AWS_REGION
      value: {{ .Parameters.s3.region }}
    - name: AWS_S3_FORCE_PATH_STYLE
      value: {{ .Parameters.s3.forcePathStyle | quote }}
    - name: AWS_SECRET_ACCESS_KEY
      valueFrom:
        secretKeyRef:
          name: {{ .Name }}-s3-credentials
          key: secretKey
    - name: BACKUP_NUM_TO_RETAIN
      value: {{ .Parameters.s3.backupRetention | default 5 | quote }}
    - name: BACKUP_SCHEDULE
      value: {{ .Parameters.s3.backupSchedule | default "5 5 * * *" }}
    - name: CLONE_AWS_ACCESS_KEY_ID
      valueFrom:
        secretKeyRef:
          name: {{ .Name }}-s3-credentials
          key: accessKey
    - name: CLONE_AWS_ENDPOINT
      value: {{ .Parameters.s3.endpoint }}
    - name: CLONE_AWS_REGION
      value: {{ .Parameters.s3.region }}
    - name: CLONE_AWS_S3_FORCE_PATH_STYLE
      value: {{ .Parameters.s3.forcePathStyle | quote }}
    - name: CLONE_AWS_SECRET_ACCESS_KEY
      valueFrom:
        secretKeyRef:
          name: {{ .Name }}-s3-credentials
          key: secretKey
    - name: CLONE_METHOD
      value: "CLONE_WITH_WALE"
    - name: CLONE_USE_WALG_RESTORE
      value: "true"
    - name: CLONE_WAL_BUCKET_SCOPE_PREFIX
      value: ""
    - name: CLONE_WAL_S3_BUCKET
      value: {{ .Parameters.s3.bucket }}
    - name: USE_WALG_BACKUP
      value: "true"
    - name: USE_WALG_RESTORE
      value: "true"
    - name: WAL_BUCKET_SCOPE_PREFIX
      value: ""
    - name: WAL_BUCKET_SCOPE_SUFFIX
      value: ""
    - name: WAL_S3_BUCKET
      value: {{ .Parameters.s3.bucket }}
    - name: WALG_DISABLE_S3_SSE
      value: "true"
    {{- end }}`,
	}

	// Create a renderer
	renderer := NewDatabaseRenderer()

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
				"common": map[string]interface{}{
					"replicas": 3,
				},
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify the output contains the expected values
		assert.Contains(t, result, "name: test-postgres")
		assert.Contains(t, result, "namespace: default")
		assert.Contains(t, result, "team: team1")
		assert.Contains(t, result, "version: \"17\"")      // Default value
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
				"common": map[string]interface{}{
					"replicas": 2,
				},
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
		assert.Contains(t, result, "team: team1")
		assert.Contains(t, result, "unbind/usd-type: postgres-operator")
		assert.Contains(t, result, "unbind/usd-version: 1.0.0")

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
				"common": map[string]interface{}{
					"replicas": 2,
				},
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
				"common": map[string]interface{}{
					"replicas": 2,
				},
				"labels": map[string]interface{}{},
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify the template renders correctly with an empty labels map
		assert.Contains(t, result, "team: team1")
		assert.NotContains(t, result, "{{ $key }}: {{ $value | quote }}") // Definition should not contain raw template syntax
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
				"common": map[string]interface{}{
					"replicas": 2,
				},
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
		assert.Contains(t, result, "value: us-east-1")                // Default region
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
				"version":     "15",
				"dockerImage": "unbindapp/spilo:15",
				"common": map[string]interface{}{
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
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify the output contains the custom values
		assert.Contains(t, result, "version: \"15\"")
		assert.Contains(t, result, "dockerImage: unbindapp/spilo:15")
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
				"version": "15",
				"common": map[string]interface{}{
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
				"labels": map[string]interface{}{
					"environment": "production",
					"app":         "custom-app",
					"tier":        "database",
				},
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify both custom values and labels are included
		assert.Contains(t, result, "version: \"15\"")
		assert.Contains(t, result, "numberOfInstances: 5")
		assert.Contains(t, result, `environment: "production"`)
		assert.Contains(t, result, `app: "custom-app"`)
		assert.Contains(t, result, `tier: "database"`)
	})

	t.Run("Parameter Validation", func(t *testing.T) {
		// Test validation of required fields
		err := renderer.Validate(map[string]interface{}{}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing required field: common")

		// Test validation of enum values
		err = renderer.Validate(map[string]interface{}{
			"common": map[string]interface{}{
				"replicas": 2,
			},
			"version": "18", // Not in enum
		}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be one of: [14 15 16 17]")

		// Test validation of range
		err = renderer.Validate(map[string]interface{}{
			"common": map[string]interface{}{
				"replicas": 10, // Above maximum
			},
		}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be <= 5")

		// Test validation of types
		err = renderer.Validate(map[string]interface{}{
			"common": map[string]interface{}{
				"replicas": "not-a-number",
			},
		}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a number")

		// Test validation with valid labels
		err = renderer.Validate(map[string]interface{}{
			"common": map[string]interface{}{
				"replicas": 3,
			},
			"labels": map[string]interface{}{
				"app": "test",
			},
		}, template.Schema)
		assert.NoError(t, err)

		// Test validation with invalid labels type
		err = renderer.Validate(map[string]interface{}{
			"common": map[string]interface{}{
				"replicas": 3,
			},
			"labels": "not-an-object", // Should be map/object
		}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be an object")

		// Test validation with invalid label value type
		err = renderer.Validate(map[string]interface{}{
			"common": map[string]interface{}{
				"replicas": 3,
			},
			"labels": map[string]interface{}{
				"app": 123, // Should be string
			},
		}, template.Schema)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a string")

		// Valid parameters should pass validation
		err = renderer.Validate(map[string]interface{}{
			"common": map[string]interface{}{
				"replicas": 3,
				"storage":  "5Gi",
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu": "200m",
					},
				},
			},
			"version": "16",
			"labels": map[string]interface{}{
				"app":         "postgres",
				"environment": "staging",
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
				"common": map[string]interface{}{
					"replicas": 2,
				},
			},
		}

		// Apply defaults
		params := renderer.applyDefaults(ctx.Parameters, template.Schema)

		// Check defaults were applied
		assert.Equal(t, "17", params["version"])
		assert.Equal(t, "unbindapp/spilo:17", params["dockerImage"])

		// Check common structure
		common, ok := params["common"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, 2, common["replicas"]) // Original value preserved
		assert.Equal(t, "1Gi", common["storage"])

		// Check nested defaults in common
		resources, ok := common["resources"].(map[string]interface{})
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
		assert.Equal(t, "us-east-1", s3["region"])
		assert.Equal(t, true, s3["forcePathStyle"])
		assert.Equal(t, 5, s3["backupRetention"])

		// Labels should not have defaults since it's additionalProperties
		_, hasLabels := params["labels"]
		assert.False(t, hasLabels, "Labels should not have defaults")
	})
}

// Helper function to create a float pointer
func floatPtr(f float64) *float64 {
	return &f
}
