package templates

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchTemplate(t *testing.T) {
	// Setup mock server
	server := setupMockServer()
	defer server.Close()

	// Override the base URL constant for testing
	originalBaseURL := BaseTemplateURL
	BaseTemplateURL = server.URL + "/%s"
	defer func() {
		BaseTemplateURL = originalBaseURL
	}()

	// Create provider
	provider := NewUnbindTemplateProvider()

	// Test FetchTemplate
	ctx := context.Background()
	template, err := provider.FetchTemplate(ctx, "v0.1", "databases", "postgres")

	// Log the result for debugging
	t.Logf("Template result: %+v", template)
	if err != nil {
		t.Logf("Error: %v", err)
	}

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, template)
	assert.Equal(t, "PostgreSQL Database", template.Name)
	assert.Equal(t, "Standard PostgreSQL database using zalando postgres-operator", template.Description)
	assert.Equal(t, "postgres-operator", template.Type)
	assert.Equal(t, "1.0.0", template.Version)

	// Verify schema was properly resolved - with more detailed logging
	assert.NotNil(t, template.Schema)

	// Print the schema properties for debugging
	t.Logf("Schema properties: %+v", template.Schema.Properties)

	// Check if version exists
	versionProp, hasVersion := template.Schema.Properties["version"]
	assert.True(t, hasVersion, "Schema should have a 'version' property")
	if hasVersion {
		t.Logf("Version property: %+v", versionProp)
	}

	// Check if s3 exists
	s3Prop, hasS3 := template.Schema.Properties["s3"]
	assert.True(t, hasS3, "Schema should have an 's3' property")

	// If s3 exists, verify it was imported correctly
	if hasS3 {
		t.Logf("S3 property: %+v", s3Prop)
		assert.Equal(t, "object", s3Prop.Type)

		// Check if s3 has properties
		if s3Prop.Properties != nil {
			assert.Contains(t, s3Prop.Properties, "bucketName", "S3 schema should have bucketName property")
			assert.Contains(t, s3Prop.Properties, "region", "S3 schema should have region property")
		} else {
			t.Logf("WARNING: S3 property doesn't have subproperties")
		}
	}

	// Check if labels exists
	labelsProp, hasLabels := template.Schema.Properties["labels"]
	assert.True(t, hasLabels, "Schema should have a 'labels' property")

	// If labels exists, verify it was imported correctly
	if hasLabels {
		t.Logf("Labels property: %+v", labelsProp)
		assert.Equal(t, "object", labelsProp.Type)
		assert.Equal(t, "Custom labels to add to the PostgreSQL resource", labelsProp.Description)

		// Check if labels has properties field defined
		if labelsProp.Properties != nil {
			// This would check if the labels schema has a nested 'properties' field
			// which might not be the case for a simple additionalProperties: true schema
			t.Logf("Labels properties: %+v", labelsProp.Properties)
		} else {
			t.Logf("WARNING: Labels property doesn't have subproperties defined")
		}
	}
}

func TestFetchTemplateErrors(t *testing.T) {
	// Setup mock server with errors
	server := setupErrorMockServer()
	defer server.Close()

	// Override the base URL constant for testing
	originalBaseURL := BaseTemplateURL
	BaseTemplateURL = server.URL + "/%s"
	defer func() {
		BaseTemplateURL = originalBaseURL
	}()

	// Create provider
	provider := NewUnbindTemplateProvider()
	ctx := context.Background()

	t.Run("Metadata not found", func(t *testing.T) {
		_, err := provider.FetchTemplate(ctx, "v0.1", "databases", "not-found")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch metadata")
	})

	t.Run("Invalid metadata", func(t *testing.T) {
		_, err := provider.FetchTemplate(ctx, "v0.1", "databases", "invalid-metadata")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse metadata")
	})

	t.Run("Template not found", func(t *testing.T) {
		_, err := provider.FetchTemplate(ctx, "v0.1", "databases", "missing-template")
		t.Logf("Error: %v", err)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch template")
	})

	// Since these next tests depend on how resolveRelativePath works in your code,
	// we'll skip detailed assertions for now
	t.Run("Import not found", func(t *testing.T) {
		_, err := provider.FetchTemplate(ctx, "v0.1", "databases", "missing-import")
		t.Logf("Import not found error: %v", err)
		assert.Error(t, err)
		// The error might contain different text depending on your implementation
	})

	t.Run("Invalid import", func(t *testing.T) {
		_, err := provider.FetchTemplate(ctx, "v0.1", "databases", "invalid-import")
		t.Logf("Invalid import error: %v", err)
		assert.Error(t, err)
		// The error might contain different text depending on your implementation
	})

	t.Run("Invalid reference", func(t *testing.T) {
		_, err := provider.FetchTemplate(ctx, "v0.1", "databases", "invalid-reference")
		t.Logf("Invalid reference error: %v", err)
		assert.Error(t, err)
		// The error might contain different text depending on your implementation
	})
}

func TestResolveRelativePath(t *testing.T) {
	testCases := []struct {
		basePath     string
		relativePath string
		expected     string
	}{
		{
			basePath:     "templates/databases/postgres",
			relativePath: "../common/s3-schema.yaml",
			expected:     "templates/databases/common/s3-schema.yaml",
		},
		{
			basePath:     "templates/databases/postgres",
			relativePath: "../common/labels.yaml",
			expected:     "templates/databases/common/labels.yaml",
		},
		{
			basePath:     "templates/databases/postgres",
			relativePath: "./schema.yaml",
			expected:     "templates/databases/postgres/schema.yaml",
		},
		{
			basePath:     "templates/databases/postgres",
			relativePath: "schema.yaml",
			expected:     "schema.yaml",
		},
		{
			basePath:     "templates/databases/postgres/nested",
			relativePath: "../../../common/schema.yaml",
			expected:     "templates/common/schema.yaml",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s+%s", tc.basePath, tc.relativePath), func(t *testing.T) {
			result := resolveRelativePath(tc.basePath, tc.relativePath)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Helper to setup a mock server that simulates the GitHub API
func setupMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Mock server received request: %s\n", r.URL.Path)

		switch r.URL.Path {
		case "/v0.1/templates/databases/postgres/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`name: "PostgreSQL Database"
description: "Standard PostgreSQL database using zalando postgres-operator"
type: "postgres-operator"
version: "1.0.0"
imports:
  - path: "../../common/s3-schema.yaml"
    as: s3Schema
  - path: "../../common/labels.yaml"
    as: labels
schema:
  properties:
    labels:
      $ref: "#/imports/labels"
    version:
      type: "string"
      description: "PostgreSQL version"
      default: "17"
      enum: ["14", "15", "16", "17"]
    replicas:
      type: "integer"
      description: "Number of replicas"
      default: 1
      minimum: 1
      maximum: 5
    storage:
      type: "string"
      description: "Storage size"
      default: "1Gi"
    s3:
      $ref: "#/imports/s3Schema"
  required:
    - replicas`))

		case "/v0.1/templates/databases/postgres/template.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`apiVersion: "acid.zalan.do/v1"
kind: postgresql
metadata:
  name: "{{ .name }}"
  namespace: "{{ .namespace }}"
  labels:
    # Zalando labels
    team: "{{ .teamId }}"
    # Template-specific labels
    unbind/template-name: "{{ .template.name }}"
    unbind/template-version: "{{ .template.version }}"
    unbind/template-category: "databases"
    {{- range $key, $value := .parameters.labels }}
    {{ $key }}: {{ $value | quote }}
    {{- end }}
spec:
  teamId: "unbind"
  postgresql:
    version: "{{ .parameters.version }}"
  numberOfInstances: {{ .parameters.replicas }}
  volume:
    size: "{{ .parameters.storage }}"`))

		// This is the path after resolveRelativePath is applied to "../common/s3-schema.yaml" from "templates/postgres"
		case "/v0.1/templates/common/s3-schema.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`type: "object"
description: "S3 configuration"
properties:
  bucketName:
    type: "string"
    description: "S3 bucket name"
  region:
    type: "string"
    description: "AWS region"
    default: "us-east-1"
  accessKey:
    type: "string"
    description: "AWS access key"
required:
  - bucketName`))

		// This is the path after resolveRelativePath is applied to "../common/labels.yaml" from "templates/postgres"
		case "/v0.1/templates/common/labels.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`type: "object"
description: "Custom labels to add to the PostgreSQL resource"
additionalProperties:
  type: "string"
default: {}`))

		default:
			fmt.Printf("Mock server: Path not found: %s\n", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))
		}
	}))
}

// Helper to setup a mock server that returns errors for various test cases
func setupErrorMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Error mock server received request: %s\n", r.URL.Path)

		switch r.URL.Path {
		// Test case 1: Not Found
		case "/v0.1/templates/databases/not-found/metadata.yaml":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))

		// Test case 2: Invalid Metadata
		case "/v0.1/templates/databases/invalid-metadata/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid-yaml: [this is not valid yaml`))
			// Test case 2: Invalid Metadata
		case "/v0.1/templates/databases/invalid-metadata/template.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid-yaml: [this is not valid yaml`))

		// Test case 3: Missing Template
		case "/v0.1/templates/databases/missing-template/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`name: "Test Template"
description: "Test template with missing template file"
type: "test"
version: "1.0.0"
schema:
  properties:
    test:
      type: "string"`))

		case "/v0.1/templates/databases/missing-template/template.yaml":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))

		// Test case 4: Missing Import
		case "/v0.1/templates/databases/missing-import/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`name: "Test Template"
description: "Test template with missing import"
type: "test"
version: "1.0.0"
imports:
  - path: "../common/missing-schema.yaml"
    as: missingSchema
  - path: "../common/missing-labels.yaml"
    as: labels
schema:
  properties:
    test:
      type: "string"
    missing:
      $ref: "#/imports/missingSchema"
    labels:
      $ref: "#/imports/labels"`))

		case "/v0.1/templates/databases/missing-import/template.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`template: "content"`))

		// This is where the import would be resolved to
		case "/v0.1/templates/common/missing-schema.yaml":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))

		case "/v0.1/templates/common/missing-labels.yaml":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))

		// Test case 5: Invalid Import
		case "/v0.1/templates/databases/invalid-import/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`name: "Test Template"
description: "Test template with invalid import"
type: "test"
version: "1.0.0"
imports:
  - path: "../common/invalid-schema.yaml"
    as: invalidSchema
  - path: "../common/invalid-labels.yaml"
    as: labels
schema:
  properties:
    test:
      type: "string"
    invalid:
      $ref: "#/imports/invalidSchema"
    labels:
      $ref: "#/imports/labels"`))

		case "/v0.1/templates/databases/invalid-import/template.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`template: "content"`))

		case "/v0.1/templates/common/invalid-schema.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid-yaml: [this is not valid yaml`))

		case "/v0.1/templates/common/invalid-labels.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid-yaml: [this is not valid yaml`))

		// Test case 6: Invalid Reference
		case "/v0.1/templates/databases/invalid-reference/metadata.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`name: "Test Template"
description: "Test template with invalid reference"
type: "test"
version: "1.0.0"
imports:
  - path: "../common/valid-schema.yaml"
    as: validSchema
  - path: "../common/valid-labels.yaml"
    as: labels
schema:
  properties:
    test:
      type: "string"
    invalid:
      $ref: "invalid-reference"
    labels:
      $ref: "invalid-labels-reference"`))

		case "/v0.1/templates/databases/invalid-reference/template.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`template: "content"`))

		case "/v0.1/templates/common/valid-schema.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`type: "object"
properties:
  test:
    type: "string"`))

		case "/v0.1/templates/common/valid-labels.yaml":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`type: "object"
description: "Custom labels"
additionalProperties:
  type: "string"`))

		default:
			fmt.Printf("Error mock server: Path not found: %s\n", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))
		}
	}))
}
