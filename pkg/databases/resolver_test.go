package databases

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchDatabase(t *testing.T) {
	// Setup mock server
	server := setupMockServer() // Ensure setupMockServer includes Postgres paths
	defer server.Close()

	// Override the base URL constant for testing
	originalBaseURL := BaseDatabaseURL
	BaseDatabaseURL = server.URL + "/%s"
	defer func() {
		BaseDatabaseURL = originalBaseURL
	}()

	// Create provider
	provider := NewDatabaseProvider()

	// Test FetchDatabase for Postgres
	ctx := context.Background()
	database, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "postgres")

	// Log the result for debugging
	t.Logf("Postgres Database result: %+v", database)
	if err != nil {
		t.Logf("Error: %v", err)
	}

	// --- Assertions for Postgres ---
	require.NoError(t, err)
	assert.NotNil(t, database)
	assert.Equal(t, "PostgreSQL Database", database.Name)
	assert.Equal(t, "Standard PostgreSQL database using zalando postgres-operator", database.Description)
	assert.Equal(t, "postgres-operator", database.Type)
	assert.Equal(t, "1.0.0", database.Version)
	assert.Equal(t, 5432, database.Port)

	assert.NotNil(t, database.Schema)
	t.Logf("Postgres Schema properties: %+v", database.Schema.Properties)

	// Check Postgres specific properties
	versionProp, hasVersion := database.Schema.Properties["version"]
	assert.True(t, hasVersion, "Postgres schema should have a 'version' property")
	if hasVersion {
		assert.Equal(t, "string", versionProp.Type)
		assert.Contains(t, versionProp.Enum, "17", "Postgres version enum should contain '17'")
	}

	s3Prop, hasS3 := database.Schema.Properties["s3"]
	assert.True(t, hasS3, "Postgres schema should have an 's3' property (from s3Schema import)")
	if hasS3 {
		assert.Equal(t, "object", s3Prop.Type)
		assert.Contains(t, s3Prop.Properties, "bucketName")
		assert.Contains(t, s3Prop.Properties, "region")
		assert.Contains(t, s3Prop.Properties, "enabled")
	}

	labelsProp, hasLabels := database.Schema.Properties["labels"]
	assert.True(t, hasLabels, "Postgres Schema should have a 'labels' property (from labels import)")
	if hasLabels {
		assert.Equal(t, "object", labelsProp.Type)
		assert.NotNil(t, labelsProp.AdditionalProperties)
		if labelsProp.AdditionalProperties != nil {
			assert.Equal(t, "string", labelsProp.AdditionalProperties.Type)
		}
	}

	envProp, hasEnv := database.Schema.Properties["environment"]
	assert.True(t, hasEnv, "Postgres schema should have an 'environment' property")
	if hasEnv {
		assert.Equal(t, "object", envProp.Type)
		assert.NotNil(t, envProp.AdditionalProperties)
		if envProp.AdditionalProperties != nil {
			assert.Equal(t, "string", envProp.AdditionalProperties.Type)
		}
		assert.NotNil(t, envProp.Default)
		assert.IsType(t, map[string]interface{}{}, envProp.Default) // Check it's a map
		assert.Empty(t, envProp.Default)                            // Check it's empty
	}

	commonProp, hasCommon := database.Schema.Properties["common"]
	assert.True(t, hasCommon, "Postgres schema should have a 'common' property (from base import)")
	if hasCommon {
		assert.Equal(t, "object", commonProp.Type)
		assert.Contains(t, commonProp.Properties, "replicas")
		assert.Contains(t, commonProp.Properties, "storage")
		assert.Contains(t, commonProp.Properties, "resources")
	}
}

func TestFetchDatabase_Redis(t *testing.T) {
	// Setup mock server
	server := setupMockServer() // Ensure setupMockServer includes Redis paths
	defer server.Close()

	// Override the base URL constant for testing
	originalBaseURL := BaseDatabaseURL
	BaseDatabaseURL = server.URL + "/%s"
	defer func() {
		BaseDatabaseURL = originalBaseURL
	}()

	// Create provider
	provider := NewDatabaseProvider()

	// Test FetchDatabase for Redis
	ctx := context.Background()
	database, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "redis")

	// Log the result for debugging
	t.Logf("Redis Database result: %+v", database)
	if err != nil {
		t.Logf("Error: %v", err)
	}

	// --- Assertions for Redis ---
	require.NoError(t, err)
	assert.NotNil(t, database)
	assert.Equal(t, "Redis", database.Name)
	assert.Equal(t, "Standard Redis installation using bitnami helm chart.", database.Description)
	assert.Equal(t, "helm", database.Type)
	assert.Equal(t, "1.0.0", database.Version)
	assert.Equal(t, 6379, database.Port)

	// Verify chart info
	assert.NotNil(t, database.Chart, "Chart should not be nil for Helm type")
	if database.Chart != nil {
		assert.Equal(t, "redis", database.Chart.Name)
		assert.Equal(t, "20.13.0", database.Chart.Version)
		assert.Equal(t, "oci://registry-1.docker.io/bitnamicharts", database.Chart.Repository)
		assert.Equal(t, "bitnami", database.Chart.RepositoryName)
	}

	// Verify schema was properly resolved - with more detailed logging
	assert.NotNil(t, database.Schema)

	// Print the schema properties for debugging
	t.Logf("Redis Schema properties: %+v", database.Schema.Properties)

	// Check Redis specific schema properties (now without chartVersion)
	_, hasChartVersion := database.Schema.Properties["chartVersion"]
	assert.False(t, hasChartVersion, "Redis Schema should NOT have a 'chartVersion' property as it's moved to top-level chart")

	// Check if common exists (resolved import)
	commonProp, hasCommon := database.Schema.Properties["common"]
	assert.True(t, hasCommon, "Redis Schema should have a 'common' property (from base import)")
	if hasCommon {
		t.Logf("Common property: %+v", commonProp)
		assert.Equal(t, "object", commonProp.Type, "Common property should be type object")
		assert.Equal(t, "Common base configuration for databases", commonProp.Description)
		assert.Contains(t, commonProp.Properties, "replicas", "Common schema should have replicas property")
		assert.Contains(t, commonProp.Properties, "storage", "Common schema should have storage property")
		assert.Contains(t, commonProp.Properties, "resources", "Common schema should have resources property")
		if resourcesProp, ok := commonProp.Properties["resources"]; ok {
			assert.Contains(t, resourcesProp.Properties, "requests", "Resources schema should have requests property")
			assert.Contains(t, resourcesProp.Properties, "limits", "Resources schema should have limits property")
		}
	}

	// Check if labels exists (resolved import)
	labelsProp, hasLabels := database.Schema.Properties["labels"]
	assert.True(t, hasLabels, "Redis Schema should have a 'labels' property (from labelsSchema import)")
	if hasLabels {
		t.Logf("Labels property: %+v", labelsProp)
		assert.Equal(t, "object", labelsProp.Type)
		assert.Equal(t, "Custom labels to add to the resource", labelsProp.Description) // Description from mock common/labels.yaml
		assert.NotNil(t, labelsProp.AdditionalProperties, "Labels property should have additionalProperties")
		if labelsProp.AdditionalProperties != nil {
			assert.Equal(t, "string", labelsProp.AdditionalProperties.Type, "Labels additionalProperties should be of type string")
		}
	}

	// Check for secretName
	secretNameProp, hasSecretName := database.Schema.Properties["secretName"]
	assert.True(t, hasSecretName, "Redis Schema should have a 'secretName' property")
	if hasSecretName {
		t.Logf("secretName property: %+v", secretNameProp)
		assert.Equal(t, "string", secretNameProp.Type)
		assert.Equal(t, "Name of the existing secret containing the Redis password", secretNameProp.Description)
	}

	// Check for secretKey
	secretKeyProp, hasSecretKey := database.Schema.Properties["secretKey"]
	assert.True(t, hasSecretKey, "Redis Schema should have a 'secretKey' property")
	if hasSecretKey {
		t.Logf("secretKey property: %+v", secretKeyProp)
		assert.Equal(t, "string", secretKeyProp.Type)
		assert.Equal(t, "Key within the secret containing the password", secretKeyProp.Description)
		assert.Equal(t, "redis-password", secretKeyProp.Default)
	}
}

func TestFetchDatabaseErrors(t *testing.T) {
	// Setup mock server with errors
	server := setupErrorMockServer()
	defer server.Close()

	// Override the base URL constant for testing
	originalBaseURL := BaseDatabaseURL
	BaseDatabaseURL = server.URL + "/%s"
	defer func() {
		BaseDatabaseURL = originalBaseURL
	}()

	// Create provider
	provider := NewDatabaseProvider()
	ctx := context.Background()

	t.Run("Metadata not found", func(t *testing.T) {
		_, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "not-found")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch metadata")
	})

	t.Run("Invalid metadata", func(t *testing.T) {
		_, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "invalid-metadata")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse metadata")
	})

	t.Run("Database definition not found", func(t *testing.T) {
		_, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "missing-database")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch database definition")
	})

	t.Run("Import not found", func(t *testing.T) {
		_, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "missing-import")
		t.Logf("Import not found error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch import")
	})

	t.Run("Invalid import", func(t *testing.T) {
		_, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "invalid-import")
		t.Logf("Invalid import error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse import")
	})

	t.Run("Invalid reference", func(t *testing.T) {
		_, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "invalid-reference")
		t.Logf("Invalid reference error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to resolve")
	})
}

func TestResolveRelativePath(t *testing.T) {
	testCases := []struct {
		basePath     string
		relativePath string
		expected     string
	}{
		{
			basePath:     "definitions/databases/postgres",
			relativePath: "../common/s3-schema.yaml",
			expected:     "definitions/databases/common/s3-schema.yaml",
		},
		{
			basePath:     "definitions/databases/postgres",
			relativePath: "../../common/labels.yaml",
			expected:     "definitions/common/labels.yaml",
		},
		{
			basePath:     "definitions/databases/redis", // Test with redis base
			relativePath: "../../common/base.yaml",
			expected:     "definitions/common/base.yaml",
		},
		{
			basePath:     "definitions/databases/postgres",
			relativePath: "./schema.yaml",
			expected:     "definitions/databases/postgres/schema.yaml",
		},
		{
			basePath:     "definitions/databases/postgres",
			relativePath: "schema.yaml",
			expected:     "schema.yaml",
		},
		{
			basePath:     "definitions/databases/postgres/nested",
			relativePath: "../../../common/schema.yaml",
			expected:     "definitions/common/schema.yaml",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s+%s", tc.basePath, tc.relativePath), func(t *testing.T) {
			result := resolveRelativePath(tc.basePath, tc.relativePath)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// --- setupMockServer (Ensure it includes both Postgres and Redis paths as provided before) ---
func setupMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Mock server received request: %s\n", r.URL.Path)

		// --- Serve Postgres Files ---
		postgresMetadataPath := "/v0.1/definitions/databases/postgres/metadata.yaml"
		postgresDefinitionPath := "/v0.1/definitions/databases/postgres/definition.yaml"
		if r.URL.Path == postgresMetadataPath {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`name: "PostgreSQL Database"
description: "Standard PostgreSQL database using zalando postgres-operator"
type: "postgres-operator"
version: "1.0.0"
port: 5432
imports:
  - path: "../../common/s3-schema.yaml"
    as: s3Schema
  - path: "../../common/labels.yaml"
    as: labels # Using 'labels' as alias
  - path: "../../common/base.yaml"
    as: base
schema:
  properties:
    common:
      $ref: "#/imports/base"
    labels:
      $ref: "#/imports/labels" # Referencing 'labels' alias
    version:
      type: "string"
      description: "PostgreSQL version"
      default: "17"
      enum: ["14", "15", "16", "17"]
    s3:
      $ref: "#/imports/s3Schema"
    environment:
      type: "object"
      description: "Environment variables to be set in the PostgreSQL container"
      additionalProperties:
        type: "string"
      default: {}`))
			return
		}
		if r.URL.Path == postgresDefinitionPath {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`apiVersion: "acid.zalan.do/v1"
kind: postgresql
metadata:
  name: "{{ .name }}"
  namespace: "{{ .namespace }}"
# ... (rest of postgres definition.yaml content) ...
spec:
  teamId: "unbind"
  postgresql:
    version: "{{ .parameters.version }}"
  numberOfInstances: {{ .parameters.common.replicas }} # Changed from .parameters.replicas
  volume:
    size: "{{ .parameters.common.storage }}" # Changed from .parameters.storage
  env: []`)) // Simplified env for brevity
			return
		}

		// --- Serve Redis Files ---
		redisMetadataPath := "/v0.1/definitions/databases/redis/metadata.yaml"
		redisDefinitionPath := "/v0.1/definitions/databases/redis/definition.yaml"
		if r.URL.Path == redisMetadataPath {
			w.WriteHeader(http.StatusOK)
			// Updated Redis metadata with chart info at top level
			w.Write([]byte(`name: "Redis"
description: "Standard Redis installation using bitnami helm chart."
type: "helm"
version: "1.0.0"
port: 6379
chart:
  name: "redis"
  version: "20.13.0"
  repository: "oci://registry-1.docker.io/bitnamicharts"
  repositoryName: "bitnami"
imports:
  - path: "../../common/base.yaml"
    as: base
  - path: "../../common/labels.yaml"
    as: labelsSchema
schema:
  properties:
    common:
      $ref: "#/imports/base"
    labels:
      $ref: "#/imports/labelsSchema"
    secretName:
      type: "string"
      description: "Name of the existing secret containing the Redis password"
    secretKey:
      type: "string"
      description: "Key within the secret containing the password"
      default: "redis-password"
`))
			return
		}
		if r.URL.Path == redisDefinitionPath {
			w.WriteHeader(http.StatusOK)
			// Use the latest definition provided previously
			w.Write([]byte(`
# definition.yaml for Redis (using Bitnami Helm Chart)
# Conditionally configures Standalone vs Replication based on replica count.
{{- $requestedReplicas := .Parameters.common.replicas | default 1 -}}
commonLabels:
  unbind/usd-type: {{ .Definition.Type | quote }}
  unbind/usd-version: {{ .Definition.Version | quote }}
  unbind/usd-category: databases
  {{- range $key, $value := .Parameters.labels }}
  {{ $key }}: {{ $value | quote }}
  {{- end }}
auth:
  enabled: true
  existingSecret: {{ .Parameters.secretName | quote }}
  existingSecretPasswordKey: {{ .Parameters.secretKey | quote }}
{{- if gt $requestedReplicas 1 }}
  architecture: replication
  master:
    count: 1
    persistence: { enabled: true, size: "{{ $.Parameters.common.storage | default "1Gi" }}" }
    resourcesPreset: "none"
    resources:
      requests: { cpu: "{{ $.Parameters.common.resources.requests.cpu | default "100m" }}", memory: "{{ $.Parameters.common.resources.requests.memory | default "128Mi" }}" }
      limits: { cpu: "{{ $.Parameters.common.resources.limits.cpu | default "500m" }}", memory: "{{ $.Parameters.common.resources.limits.memory | default "256Mi" }}" }
  replica:
    replicaCount: {{ sub $requestedReplicas 1 }}
    persistence: { enabled: true, size: "{{ $.Parameters.common.storage | default "1Gi" }}" }
    resourcesPreset: "none"
    resources:
      requests: { cpu: "{{ $.Parameters.common.resources.requests.cpu | default "100m" }}", memory: "{{ $.Parameters.common.resources.requests.memory | default "128Mi" }}" }
      limits: { cpu: "{{ $.Parameters.common.resources.limits.cpu | default "500m" }}", memory: "{{ $.Parameters.common.resources.limits.memory | default "256Mi" }}" }
{{- else }}
  architecture: standalone
  master:
    count: 1
    persistence: { enabled: true, size: "{{ $.Parameters.common.storage | default "1Gi" }}" }
    resourcesPreset: "none"
    resources:
      requests: { cpu: "{{ $.Parameters.common.resources.requests.cpu | default "100m" }}", memory: "{{ $.Parameters.common.resources.requests.memory | default "128Mi" }}" }
      limits: { cpu: "{{ $.Parameters.common.resources.limits.cpu | default "500m" }}", memory: "{{ $.Parameters.common.resources.limits.memory | default "256Mi" }}" }
  replica:
    replicaCount: 0
    persistence: { enabled: false }
{{- end }}
`)) // Compacted YAML for brevity
			return
		}

		// --- Serve Common Files (used by both Redis and potentially others) ---
		basePath := "/v0.1/definitions/common/base.yaml"
		labelsPath := "/v0.1/definitions/common/labels.yaml"
		s3SchemaPath := "/v0.1/definitions/common/s3-schema.yaml"

		if r.URL.Path == basePath {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`type: "object"
description: "Common base configuration for databases"
properties:
  replicas: {type: "integer", description: "Number of replicas", default: 1, minimum: 1, maximum: 5}
  storage: {type: "string", description: "Storage size", default: "1Gi"}
  resources:
    type: "object"
    description: "Resource requirements"
    properties:
      requests:
        type: "object"
        properties: { cpu: {type: "string", default: "100m"}, memory: {type: "string", default: "128Mi"} }
      limits:
        type: "object"
        properties: { cpu: {type: "string", default: "500m"}, memory: {type: "string", default: "256Mi"} }`)) // Compacted
			return
		}
		if r.URL.Path == labelsPath {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`type: "object"
description: "Custom labels to add to the resource" # Generic description
additionalProperties: {type: "string"}
default: {}`)) // Compacted
			return
		}
		if r.URL.Path == s3SchemaPath { // Keep for postgres test
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`type: "object"
description: "S3 configuration"
properties: { bucketName: {type: "string"}, region: {type: "string", default: "us-east-1"}, enabled: {type: "boolean", default: false} }
required: [bucketName]`)) // Compacted
			return
		}

		// --- Fallback for any other path ---
		fmt.Printf("Mock server: Path not found: %s\n", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}))
}

func setupErrorMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Error mock server received request: %s\n", r.URL.Path)
		path := r.URL.Path

		// Explicit switch on the requested path
		switch path {
		// --- Test case 1: Metadata Not Found ---
		case "/v0.1/definitions/databases/not-found/metadata.yaml":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))

		// --- Test case 2: Invalid Metadata (Metadata Fetch/Parse Error) ---
		case "/v0.1/definitions/databases/invalid-metadata/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			// Use raw string literal with CLEARLY invalid YAML structure
			w.Write([]byte(`
name: Invalid Metadata Example
invalid_structure: { key: "value" `)) // Missing closing brace '}'
		case "/v0.1/definitions/databases/invalid-metadata/definition.yaml":
			// This might not even be called if metadata parsing fails first, but provide it anyway.
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"key": "value"}`))

		// --- Test case 3: Missing Database Definition (Definition Fetch Error) ---
		case "/v0.1/definitions/databases/missing-database/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			// VALID metadata
			w.Write([]byte(`
name: "Test Database Missing Definition"
description: "Valid metadata, but definition file is missing"
type: "test"
version: "1.0.0"
schema:
  properties:
    test: { type: "string" }
`))
		case "/v0.1/definitions/databases/missing-database/definition.yaml":
			w.WriteHeader(http.StatusNotFound) // Definition fetch fails
			w.Write([]byte("Not Found"))

		// --- Test case 4: Missing Import (Import Fetch Error) ---
		case "/v0.1/definitions/databases/missing-import/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			// VALID metadata
			w.Write([]byte(`
name: "Test Database Missing Import"
description: "Metadata refers to an import file that doesn't exist"
type: "test"
version: "1.0.0"
imports:
  - path: "../../common/missing.yaml"
    as: miss
schema:
  properties:
    p1:
      $ref: "#/imports/miss"
`))
		case "/v0.1/definitions/databases/missing-import/definition.yaml":
			// Provide a valid definition file for this case
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"key": "value"}`))
		case "/v0.1/definitions/common/missing.yaml":
			w.WriteHeader(http.StatusNotFound) // Import fetch fails
			w.Write([]byte("Not Found"))

		// --- Test case 5: Invalid Import (Import Parse Error) ---
		case "/v0.1/definitions/databases/invalid-import/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			// VALID metadata
			w.Write([]byte(`
name: "Test Database Invalid Import"
description: "Metadata refers to an import file with invalid content"
type: "test"
version: "1.0.0"
imports:
  - path: "../../common/invalid.yaml"
    as: inv
schema:
  properties:
    p1:
      $ref: "#/imports/inv"
`))
		case "/v0.1/definitions/databases/invalid-import/definition.yaml":
			// Provide a valid definition file for this case
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"key": "value"}`))
		case "/v0.1/definitions/common/invalid.yaml":
			w.WriteHeader(http.StatusOK)
			// Import content is invalid YAML, causing parse error during import processing
			w.Write([]byte(`invalid: yaml: [`)) // Keep this invalid

		// --- Test case 6: Invalid Reference (Reference Resolution Error) ---
		case "/v0.1/definitions/databases/invalid-reference/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			// VALID metadata
			w.Write([]byte(`
name: "Test Database Invalid Reference"
description: "Metadata uses a $ref to an import alias that doesn't exist"
type: "test"
version: "1.0.0"
imports:
  - path: "../../common/valid-schema.yaml" # This import exists and is valid
    as: validSchema
schema:
  properties:
    p1:
      $ref: "#/imports/nonExistentImport" # This import alias is not defined above
    valid:
      $ref: "#/imports/validSchema" # This one is fine
`))
		case "/v0.1/definitions/databases/invalid-reference/definition.yaml":
			// Provide a valid definition file for this case
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"key": "value"}`))
		case "/v0.1/definitions/common/valid-schema.yaml": // Serve the valid import needed by this test case
			w.WriteHeader(http.StatusOK)
			// Use raw string literal for valid YAML
			w.Write([]byte(`
type: "object"
properties:
  propA: { type: "string" }
`))

		// Default Fallback for any unhandled paths
		default:
			fmt.Printf("Error mock server: Path not handled: %s\n", path)
			// Return 404 for any unhandled path in the error server context
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found by Error Mock"))
		}
	}))
}

func TestHelmChartValidation(t *testing.T) {
	// Setup mock server with errors
	server := setupHelmErrorMockServer()
	defer server.Close()

	// Override the base URL constant for testing
	originalBaseURL := BaseDatabaseURL
	BaseDatabaseURL = server.URL + "/%s"
	defer func() {
		BaseDatabaseURL = originalBaseURL
	}()

	// Create provider
	provider := NewDatabaseProvider()
	ctx := context.Background()

	t.Run("Missing Chart Info", func(t *testing.T) {
		_, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "missing-chart")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "chart information is required for Helm database type")
	})

	t.Run("Incomplete Chart Info - Missing Name", func(t *testing.T) {
		_, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "incomplete-chart-name")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "chart name is required for Helm database type")
	})

	t.Run("Incomplete Chart Info - Missing Version", func(t *testing.T) {
		_, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "incomplete-chart-version")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "chart version is required for Helm database type")
	})

	t.Run("Incomplete Chart Info - Missing Repository", func(t *testing.T) {
		_, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "incomplete-chart-repo")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "chart repository is required for Helm database type")
	})

	t.Run("Incomplete Chart Info - Missing RepositoryName", func(t *testing.T) {
		_, err := provider.FetchDatabaseDefinition(ctx, "v0.1", "incomplete-chart-reponame")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "chart repositoryName is required for Helm database type")
	})
}

func setupHelmErrorMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Helm error mock server received request: %s\n", r.URL.Path)
		path := r.URL.Path

		// Provide a valid definition file for all cases
		if strings.HasSuffix(path, "definition.yaml") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`commonLabels:
  unbind/usd-type: {{ .Definition.Type }}
  unbind/usd-version: {{ .Definition.Version }}
  unbind/usd-category: databases
auth:
  enabled: true
  existingSecret: {{ .Parameters.secretName }}
  existingSecretPasswordKey: {{ .Parameters.secretKey }}
architecture: standalone`))
			return
		}

		// Test cases for missing or incomplete chart information
		switch path {
		case "/v0.1/definitions/databases/missing-chart/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`name: "Redis Missing Chart"
description: "Redis definition without chart information"
type: "helm"
version: "1.0.0"
port: 6379
schema:
  properties:
    secretName:
      type: "string"
      description: "Name of the secret to store Redis password"
    secretKey:
      type: "string"
      description: "Key in the secret that contains the redis password"`))

		case "/v0.1/definitions/databases/incomplete-chart-name/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`name: "Redis Incomplete Chart Name"
description: "Redis definition with missing chart name"
type: "helm"
version: "1.0.0"
port: 6379
chart:
  version: "20.13.0"
  repository: "oci://registry-1.docker.io/bitnamicharts"
  repositoryName: "bitnami"
schema:
  properties:
    secretName:
      type: "string"
      description: "Name of the secret to store Redis password"
    secretKey:
      type: "string"
      description: "Key in the secret that contains the redis password"`))

		case "/v0.1/definitions/databases/incomplete-chart-version/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`name: "Redis Incomplete Chart Version"
description: "Redis definition with missing chart version"
type: "helm"
version: "1.0.0"
port: 6379
chart:
  name: "redis"
  repository: "oci://registry-1.docker.io/bitnamicharts"
  repositoryName: "bitnami"
schema:
  properties:
    secretName:
      type: "string"
      description: "Name of the secret to store Redis password"
    secretKey:
      type: "string"
      description: "Key in the secret that contains the redis password"`))

		case "/v0.1/definitions/databases/incomplete-chart-repo/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`name: "Redis Incomplete Chart Repository"
description: "Redis definition with missing chart repository"
type: "helm"
version: "1.0.0"
port: 6379
chart:
  name: "redis"
  version: "20.13.0"
  repositoryName: "bitnami"
schema:
  properties:
    secretName:
      type: "string"
      description: "Name of the secret to store Redis password"
    secretKey:
      type: "string"
      description: "Key in the secret that contains the redis password"`))

		case "/v0.1/definitions/databases/incomplete-chart-reponame/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`name: "Redis Incomplete Chart RepositoryName"
description: "Redis definition with missing chart repositoryName"
type: "helm"
version: "1.0.0"
port: 6379
chart:
  name: "redis"
  version: "20.13.0"
  repository: "oci://registry-1.docker.io/bitnamicharts"
schema:
  properties:
    secretName:
      type: "string"
      description: "Name of the secret to store Redis password"
    secretKey:
      type: "string"
      description: "Key in the secret that contains the redis password"`))

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
	}))
}
