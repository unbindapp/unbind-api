package databases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

func TestDefinitionRendering(t *testing.T) {
	// Create a sample template similar to the postgres one
	template := &Definition{
		Name:        "PostgreSQL Database",
		Category:    DB_CATEGORY,
		Port:        5432,
		Description: "Standard PostgreSQL database using zalando postgres-operator",
		Type:        "postgres-operator",
		Version:     "1.0.0",
		Schema: DefinitionParameterSchema{
			Properties: map[string]ParameterProperty{
				"enableMasterLoadBalancer": {
					Type:        "boolean",
					Description: "Enable master load balancer",
					Default:     false,
				},
				"common": {
					Type: "object",
					Properties: map[string]ParameterProperty{
						"replicas": {
							Type:        "integer",
							Description: "Number of replicas",
							Default:     1,
							Minimum:     utils.ToPtr[float64](1),
							Maximum:     utils.ToPtr[float64](5),
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
				"environment": {
					Type:        "object",
					Description: "Environment variables to be set in the PostgreSQL container",
					AdditionalProperties: &ParameterProperty{
						Type: "string",
					},
					Default: map[string]interface{}{},
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
  enableMasterLoadBalancer: {{ .Parameters.enableMasterLoadBalancer | default false }}
  dockerImage: {{ if .Parameters.dockerImage }}{{ .Parameters.dockerImage }}{{ else }}{{ printf "unbindapp/spilo:%s-latest" (.Parameters.version | default "17") }}{{ end }}
  postgresql:
    version: {{ .Parameters.version | default "17" | quote }}
  numberOfInstances: {{ .Parameters.common.replicas | default 1 }}
  allowedSourceRanges:
    - 0.0.0.0/0
  patroni:
    pg_hba:
    # Keep for pam authentication
    - "hostssl all +pamrole all pam"
    # Force SSL for external
    - "hostssl all all 0.0.0.0/0 md5"
    # Allow non‑SSL md5 inside the pod network (k3s)
    - "host    all all 10.42.0.0/16 md5"
    # Allow non-ssl md5 inside the pod network (others like microk8s)
    - "host    all all 10.1.0.0/16 md5"
    # Allow non-ssl md5 inside the pod network (others common k8s distributions)
    - "host    all all 10.0.0.0/8 md5"
    # Local loopback
    - "host    all all 127.0.0.1/32 trust"
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
    {{- end }}
    {{- if .Parameters.environment }}
    {{- range $key, $value := .Parameters.environment }}
    - name: {{ $key }}
      value: {{ $value | quote }}
    {{- end }}
    {{- end }}`,
	}

	// Create a renderer
	renderer := NewDatabaseRenderer()

	t.Run("Basic Rendering", func(t *testing.T) {
		// Create render context with minimal parameters
		ctx := &RenderContext{
			Name:      "test-postgres",
			Namespace: "default",
			TeamID:    "team1",
			Parameters: map[string]interface{}{
				"common": map[string]interface{}{
					"replicas": 3,
				},
			},
			Definition: Definition{
				Type:    "postgres-operator",
				Version: "1.0.0",
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify the output contains the expected values for PostgreSQL object
		assert.Contains(t, result, "name: test-postgres")
		assert.Contains(t, result, "namespace: default")
		assert.Contains(t, result, "team: team1")
		assert.Contains(t, result, "version: \"17\"")                 // Default value
		assert.Contains(t, result, "numberOfInstances: 3")            // Provided value
		assert.Contains(t, result, "size: 1Gi")                       // Default value
		assert.Contains(t, result, "enableMasterLoadBalancer: false") // Default value

		// Test the new dockerImage format with version-latest
		assert.Contains(t, result, "dockerImage: unbindapp/spilo:17-latest") // Should use default version with -latest

		// Check for Patroni configuration
		assert.Contains(t, result, "patroni:")
		assert.Contains(t, result, "\"hostssl all +pamrole all pam\"")
		assert.Contains(t, result, "\"host    all all 10.42.0.0/16 md5\"")

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
		assert.Len(t, objects, 1)
	})

	t.Run("With Environment Variables", func(t *testing.T) {
		// Create render context with environment variables
		ctx := &RenderContext{
			Name:      "test-postgres-env",
			Namespace: "default",
			TeamID:    "team1",
			Parameters: map[string]interface{}{
				"common": map[string]interface{}{
					"replicas": 2,
				},
				"environment": map[string]interface{}{
					"POSTGRES_LOG_STATEMENT": "all",
					"PGAUDIT_LOG":            "DDL",
					"MAX_CONNECTIONS":        "100",
					"LOGGING_COLLECTOR":      "on",
					"LOG_STATEMENT":          "ddl",
				},
			},
			Definition: Definition{
				Type:    "postgres-operator",
				Version: "1.0.0",
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify environment variables are included in PostgreSQL object
		assert.Contains(t, result, `name: POSTGRES_LOG_STATEMENT`)
		assert.Contains(t, result, `value: "all"`)
		assert.Contains(t, result, `name: PGAUDIT_LOG`)
		assert.Contains(t, result, `value: "DDL"`)
		assert.Contains(t, result, `name: MAX_CONNECTIONS`)
		assert.Contains(t, result, `value: "100"`)
		assert.Contains(t, result, `name: LOGGING_COLLECTOR`)
		assert.Contains(t, result, `value: "on"`)
		assert.Contains(t, result, `name: LOG_STATEMENT`)
		assert.Contains(t, result, `value: "ddl"`)

		// Ensure values are properly quoted
		assert.NotContains(t, result, "value: all") // Should be quoted

		// Standard envs should still be there
		assert.Contains(t, result, `name: ALLOW_NOSSL`)
		assert.Contains(t, result, `value: "true"`)

		// Check docker image format
		assert.Contains(t, result, "dockerImage: unbindapp/spilo:17-latest") // Should use default version with -latest

		// Parse to objects to ensure we have correct resources
		objects, err := renderer.RenderToObjects(result)
		require.NoError(t, err)
		assert.Len(t, objects, 1)
	})

	t.Run("S3 with Environment Variables", func(t *testing.T) {
		// Create render context with S3 enabled and environment variables
		ctx := &RenderContext{
			Name:      "test-postgres-s3-env",
			Namespace: "default",
			TeamID:    "team1",
			Parameters: map[string]interface{}{
				"common": map[string]interface{}{
					"replicas": 2,
				},
				"s3": map[string]interface{}{
					"enabled": true,
					"bucket":  "test-bucket",
				},
				"environment": map[string]interface{}{
					"POSTGRES_LOG_STATEMENT": "all",
					"PGUSER_SUPERUSER":       "true",
				},
			},
			Definition: Definition{
				Type:    "postgres-operator",
				Version: "1.0.0",
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify S3 section is included
		assert.Contains(t, result, "AWS_ACCESS_KEY_ID")
		assert.Contains(t, result, "value: test-bucket")

		// Verify environment variables are included
		assert.Contains(t, result, `name: POSTGRES_LOG_STATEMENT`)
		assert.Contains(t, result, `value: "all"`)
		assert.Contains(t, result, `name: PGUSER_SUPERUSER`)
		assert.Contains(t, result, `value: "true"`)

		// Check docker image format
		assert.Contains(t, result, "dockerImage: unbindapp/spilo:17-latest") // Should use default version with -latest

		// Parse to objects to ensure we have correct resources
		objects, err := renderer.RenderToObjects(result)
		require.NoError(t, err)
		assert.Len(t, objects, 1)
	})

	t.Run("Empty Environment Variables", func(t *testing.T) {
		// Create render context with empty environment variables
		ctx := &RenderContext{
			Name:      "test-postgres-empty-env",
			Namespace: "default",
			TeamID:    "team1",
			Parameters: map[string]interface{}{
				"common": map[string]interface{}{
					"replicas": 2,
				},
				"environment": map[string]interface{}{},
			},
			Definition: Definition{
				Type:    "postgres-operator",
				Version: "1.0.0",
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// No additional environment variables should be rendered
		// but the template should still work without errors
		assert.Contains(t, result, `name: ALLOW_NOSSL`)
		assert.Contains(t, result, `value: "true"`)

		// Parse to objects to ensure we have correct resources
		objects, err := renderer.RenderToObjects(result)
		require.NoError(t, err)
		assert.Len(t, objects, 1)
	})

	t.Run("Custom Version", func(t *testing.T) {
		// Create render context with custom version
		ctx := &RenderContext{
			Name:      "test-postgres-version",
			Namespace: "default",
			TeamID:    "team1",
			Parameters: map[string]interface{}{
				"common": map[string]interface{}{
					"replicas": 2,
				},
				"version": "16", // Using a different PostgreSQL version
			},
			Definition: Definition{
				Type:    "postgres-operator",
				Version: "1.0.0",
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify the output contains the expected values
		assert.Contains(t, result, "version: \"16\"") // Custom version

		// Test the new dockerImage format with custom version-latest
		assert.Contains(t, result, "dockerImage: unbindapp/spilo:16-latest") // Should use provided version with -latest

		// Parse to objects
		objects, err := renderer.RenderToObjects(result)
		require.NoError(t, err)
		assert.Len(t, objects, 1)
	})

	t.Run("Custom Docker Image", func(t *testing.T) {
		// Create render context with custom docker image
		ctx := &RenderContext{
			Name:      "test-postgres-docker",
			Namespace: "default",
			TeamID:    "team1",
			Parameters: map[string]interface{}{
				"common": map[string]interface{}{
					"replicas": 2,
				},
				"version":     "15",                     // Using a different PostgreSQL version
				"dockerImage": "custom/postgres:latest", // Custom docker image
			},
			Definition: Definition{
				Type:    "postgres-operator",
				Version: "1.0.0",
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify the output contains the expected values
		assert.Contains(t, result, "version: \"15\"") // Custom version

		// Test that the custom docker image overrides the formatted version pattern
		assert.Contains(t, result, "dockerImage: custom/postgres:latest") // Should use the exact provided image
		assert.NotContains(t, result, "dockerImage: unbindapp/spilo")     // Should not use the default image format

		// Parse to objects
		objects, err := renderer.RenderToObjects(result)
		require.NoError(t, err)
		assert.Len(t, objects, 1)
	})

	t.Run("With Custom Labels", func(t *testing.T) {
		// Create render context with custom labels
		ctx := &RenderContext{
			Name:      "test-postgres-labels",
			Namespace: "default",
			TeamID:    "team1",
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
				"storage": "10Gi",
			},
			Definition: Definition{
				Type:    "postgres-operator",
				Version: "1.0.0",
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify custom labels are included in PostgreSQL object
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

		// Check docker image format
		assert.Contains(t, result, "dockerImage: unbindapp/spilo:17-latest") // Should use default version with -latest
	})

	t.Run("S3 Enabled", func(t *testing.T) {
		// Create render context with S3 enabled
		ctx := &RenderContext{
			Name:      "test-postgres-s3",
			Namespace: "default",
			TeamID:    "team1",
			Parameters: map[string]interface{}{
				"common": map[string]interface{}{
					"replicas": 2,
				},
				"s3": map[string]interface{}{
					"enabled": true,
					"bucket":  "test-bucket",
				},
			},
			Definition: Definition{
				Type:    "postgres-operator",
				Version: "1.0.0",
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

		// Check docker image format
		assert.Contains(t, result, "dockerImage: unbindapp/spilo:17-latest") // Should use default version with -latest

		// Parse to objects to ensure we have correct resources
		objects, err := renderer.RenderToObjects(result)
		require.NoError(t, err)
		assert.Len(t, objects, 1) // PostgreSQL object
	})

	t.Run("Complex Environment Variables", func(t *testing.T) {
		// Create render context with complex environment variables including special characters
		ctx := &RenderContext{
			Name:      "test-postgres-complex-env",
			Namespace: "default",
			TeamID:    "team1",
			Parameters: map[string]interface{}{
				"common": map[string]interface{}{
					"replicas": 2,
				},
				"environment": map[string]interface{}{
					"POSTGRES_SHARED_BUFFERS":       "256MB",
					"POSTGRES_EFFECTIVE_CACHE_SIZE": "1GB",
					"POSTGRES_WORK_MEM":             "16MB",
					"POSTGRES_CONF_ADDITIONAL":      "log_min_duration_statement=200\nrandom_page_cost=1.1",
					"SPECIAL_CHARS_TEST":            "!@#$%^&*()_+",
					"QUOTED_VALUE":                  "\"quoted string\"",
				},
			},
			Definition: Definition{
				Type:    "postgres-operator",
				Version: "1.0.0",
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify complex environment variables are included in PostgreSQL object
		assert.Contains(t, result, `name: POSTGRES_SHARED_BUFFERS`)
		assert.Contains(t, result, `value: "256MB"`)
		assert.Contains(t, result, `name: POSTGRES_EFFECTIVE_CACHE_SIZE`)
		assert.Contains(t, result, `value: "1GB"`)
		assert.Contains(t, result, `name: SPECIAL_CHARS_TEST`)
		assert.Contains(t, result, `value: "!@#$%^&*()_+"`)
		assert.Contains(t, result, `name: QUOTED_VALUE`)
		assert.Contains(t, result, `value: "\"quoted string\""`)

		// Ensure multiline values are handled correctly
		assert.Contains(t, result, `name: POSTGRES_CONF_ADDITIONAL`)
		assert.Contains(t, result, `value: "log_min_duration_statement=200\nrandom_page_cost=1.1"`)

		// Parse to objects to ensure we have correct resources
		objects, err := renderer.RenderToObjects(result)
		require.NoError(t, err)
		assert.Len(t, objects, 1)
	})
}

func TestHelmRendering(t *testing.T) {
	// Create a sample Helm template for Redis with chart info at top level
	template := &Definition{
		Name:        "Redis",
		Category:    DB_CATEGORY,
		Port:        6379,
		Description: "Standard Redis installation using bitnami helm chart.",
		Type:        "helm",
		Version:     "1.0.0",
		Chart: &HelmChartInfo{
			Name:           "redis",
			Version:        "20.13.0",
			Repository:     "oci://registry-1.docker.io/bitnamicharts",
			RepositoryName: "bitnami",
		},
		Schema: DefinitionParameterSchema{
			Properties: map[string]ParameterProperty{
				"secretName": {
					Type:        "string",
					Description: "Name of the secret to store Redis password",
				},
				"secretKey": {
					Type:        "string",
					Description: "Key in the secret that contains the redis password",
				},
				"common": {
					Type: "object",
					Properties: map[string]ParameterProperty{
						"replicas": {
							Type:        "integer",
							Description: "Number of replicas",
							Default:     1,
							Minimum:     utils.ToPtr[float64](1),
							Maximum:     utils.ToPtr[float64](5),
						},
						"storage": {
							Type:        "string",
							Description: "Storage size",
							Default:     "1Gi",
						},
						"resources": {
							Type: "object",
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
											Default:     "500m",
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
				"labels": {
					Type:        "object",
					Description: "Custom labels to add to the Redis resources",
					AdditionalProperties: &ParameterProperty{
						Type: "string",
					},
				},
			},
			Required: []string{"secretName", "secretKey"},
		},
		Content: `{{- $requestedReplicas := .Parameters.common.replicas | default 1 -}}
# This file defines the values passed to the Bitnami Redis Helm chart.
# Apply common labels to all resources created by the chart
commonLabels:
  # Standard USD labels
  unbind/usd-type: {{ .Definition.Type }}
  unbind/usd-version: {{ .Definition.Version }}
  unbind/usd-category: databases
  {{- range $key, $value := .Parameters.labels }}
  {{ $key }}: {{ $value }}
  {{- end }}
auth:
  enabled: true
  existingSecret: {{ .Parameters.secretName }}
  existingSecretPasswordKey: {{ .Parameters.secretKey }}
{{- if gt $requestedReplicas 1 }}
# --- Replication Mode ---
# Triggered if .Parameters.common.replicas >= 2
# Deploys 1 Master + (N-1) Replicas
architecture: replication
# Master configuration (always 1 in this replication mode)
master:
  count: 1
  persistence:
    enabled: true
    size: {{ $.Parameters.common.storage | default "1Gi" }}
  # Disable resource presets to use specific values below
  resourcesPreset: "none"
  # Map resource requests/limits from common parameters
  resources:
    requests:
      cpu: {{ $.Parameters.common.resources.requests.cpu | default "100m" }}
      memory: {{ $.Parameters.common.resources.requests.memory | default "128Mi" }}
    limits:
      cpu: {{ $.Parameters.common.resources.limits.cpu | default "500m" }}
      memory: {{ $.Parameters.common.resources.limits.memory | default "256Mi" }}
# Replica configuration (N-1 replicas)
replica:
  replicaCount: {{ sub $requestedReplicas 1 }}
  persistence:
    enabled: true
    size: {{ $.Parameters.common.storage | default "1Gi" }}
  # Disable resource presets to use specific values below
  resourcesPreset: "none"
  resources:
    requests:
      cpu: {{ $.Parameters.common.resources.requests.cpu | default "100m" }}
      memory: {{ $.Parameters.common.resources.requests.memory | default "128Mi" }}
    limits:
      cpu: {{ $.Parameters.common.resources.limits.cpu | default "500m" }}
      memory: {{ $.Parameters.common.resources.limits.memory | default "256Mi" }}
{{- else }}
# --- Standalone Mode ---
# Triggered if .Parameters.common.replicas = 1 (or omitted, as default is 1)
# Deploys 1 Master, 0 Replicas
architecture: standalone
# Master configuration (the only node)
master:
  count: 1
  persistence:
    # Enable persistence for the single master node
    enabled: true
    # Map storage size from common parameters
    size: {{ $.Parameters.common.storage | default "1Gi" }}
  # Disable resource presets to use specific values below
  resourcesPreset: "none"
  # Map resource requests/limits from common parameters
  resources:
    requests:
      cpu: {{ $.Parameters.common.resources.requests.cpu | default "100m" }}
      memory: {{ $.Parameters.common.resources.requests.memory | default "128Mi" }}
    limits:
      cpu: {{ $.Parameters.common.resources.limits.cpu | default "500m" }}
      memory: {{ $.Parameters.common.resources.limits.memory | default "256Mi" }}
# Replica configuration (disabled for standalone)
replica:
  replicaCount: 0
  persistence:
    enabled: false
{{- end }}`,
	}

	// Create a renderer
	renderer := NewDatabaseRenderer()

	t.Run("Basic Helm Redis Rendering With Top-Level Chart", func(t *testing.T) {
		// Create render context with required parameters
		ctx := &RenderContext{
			Name:      "test-redis",
			Namespace: "default",
			TeamID:    "team1",
			Parameters: map[string]interface{}{
				"secretName": "redis-secret",
				"secretKey":  "redis-password",
				"common": map[string]interface{}{
					"replicas": 1, // Standalone mode
				},
			},
			Definition: Definition{
				Type:    "helm",
				Version: "1.0.0",
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify the output contains both HelmRepository and HelmRelease resources
		assert.Contains(t, result, "kind: HelmRepository")
		assert.Contains(t, result, "kind: HelmRelease")

		// Verify repository details
		assert.Contains(t, result, "name: bitnami")
		assert.Contains(t, result, "oci://registry-1.docker.io/bitnamicharts")

		// Verify release details
		assert.Contains(t, result, "name: test-redis")
		assert.Contains(t, result, "namespace: default")

		// Verify chart details from top-level config
		assert.Contains(t, result, "chart: redis")
		assert.Contains(t, result, "version: 20.13.0")

		// Verify Helm values - fixed to match actual output format
		assert.Contains(t, result, "architecture: standalone")
		assert.Contains(t, result, "unbind/usd-type: helm")                     // No quotes in actual output
		assert.Contains(t, result, "existingSecret: redis-secret")              // No quotes in actual output
		assert.Contains(t, result, "existingSecretPasswordKey: redis-password") // No quotes in actual output

		// Parse to objects to verify structure
		objects, err := renderer.RenderToObjects(result)
		require.NoError(t, err)
		assert.Len(t, objects, 2) // HelmRepository and HelmRelease
	})

	t.Run("Redis Replication Mode With Top-Level Chart", func(t *testing.T) {
		// Create render context with replication mode
		ctx := &RenderContext{
			Name:      "test-redis-replication",
			Namespace: "default",
			TeamID:    "team1",
			Parameters: map[string]interface{}{
				"secretName": "redis-secret",
				"secretKey":  "redis-password",
				"common": map[string]interface{}{
					"replicas": 3, // Replication mode with 3 nodes (1 master + 2 replicas)
					"storage":  "5Gi",
				},
			},
			Definition: Definition{
				Type:    "helm",
				Version: "1.0.0",
			},
		}

		// Render the template
		result, err := renderer.Render(template, ctx)
		require.NoError(t, err)

		// Verify top-level chart info is used
		assert.Contains(t, result, "chart: redis")
		assert.Contains(t, result, "version: 20.13.0")
		assert.Contains(t, result, "oci://registry-1.docker.io/bitnamicharts")

		// Verify replication mode configuration
		assert.Contains(t, result, "architecture: replication")
		assert.Contains(t, result, "replicaCount: 2") // 3 total - 1 master = 2 replicas
		assert.Contains(t, result, "size: 5Gi")       // Custom storage size

		// Parse to objects
		objects, err := renderer.RenderToObjects(result)
		require.NoError(t, err)
		assert.Len(t, objects, 2) // HelmRepository and HelmRelease
	})
}
