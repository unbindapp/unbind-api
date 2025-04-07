package databases

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// FetchDatabaseDefinition fetches a db from GitHub
func (self *DatabaseProvider) FetchDatabaseDefinition(ctx context.Context, tagVersion, dbType string) (*Definition, error) {
	// Base version URL
	baseURL := fmt.Sprintf(BaseDatabaseURL, tagVersion)
	// Fetch files
	metadataURL := fmt.Sprintf("%s/definitions/%s/%s/metadata.yaml", baseURL, DB_CATEGORY, dbType)

	defURL := fmt.Sprintf("%s/definitions/%s/%s/definition.yaml", baseURL, DB_CATEGORY, dbType)

	metadataBytes, err := self.fetchURL(ctx, metadataURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}

	defBytes, err := self.fetchURL(ctx, defURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch database definition: %w", err)
	}

	// Parse metadata
	var metadata DefinitionMetadata
	if err := yaml.Unmarshal(metadataBytes, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Initialize imports map
	metadata.Schema.Imports = make(map[string]interface{})

	// Process imports
	for _, imp := range metadata.Imports {
		// Handle relative paths for imports
		// Determine the base directory of the current database
		dbBasePath := fmt.Sprintf("definitions/%s/%s", DB_CATEGORY, dbType)

		// Resolve the relative path
		importPath := resolveRelativePath(dbBasePath, imp.Path)
		importURL := fmt.Sprintf("%s/%s", baseURL, importPath)

		importBytes, err := self.fetchURL(ctx, importURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch import %s: %w", imp.Path, err)
		}

		var importSchema interface{}
		if err := yaml.Unmarshal(importBytes, &importSchema); err != nil {
			return nil, fmt.Errorf("failed to parse import %s: %w", imp.Path, err)
		}

		metadata.Schema.Imports[imp.As] = importSchema
	}

	// Resolve references
	resolvedSchema, err := self.resolveReferences(metadata.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve references: %w", err)
	}

	db := &Definition{
		Name:        metadata.Name,
		Description: metadata.Description,
		Type:        metadata.Type,
		Version:     metadata.Version,
		Schema:      resolvedSchema,
		Content:     string(defBytes),
	}

	return db, nil
}

// resolveRelativePath resolves a relative path from a base path
func resolveRelativePath(basePath, relativePath string) string {
	// If the relativePath is not relative, return it as is
	if !strings.HasPrefix(relativePath, "../") && !strings.HasPrefix(relativePath, "./") {
		return relativePath
	}

	baseParts := strings.Split(basePath, "/")
	relativeParts := strings.Split(relativePath, "/")

	// Handle "./" prefix by just removing it
	if strings.HasPrefix(relativePath, "./") {
		relativeParts = relativeParts[1:] // Skip the "." part
		return strings.Join(append(baseParts, relativeParts...), "/")
	}

	// Handle "../" by moving up in the directory hierarchy
	resultParts := append([]string{}, baseParts...)

	for i, part := range relativeParts {
		if part == ".." {
			if len(resultParts) > 0 {
				resultParts = resultParts[:len(resultParts)-1] // Move up one directory
			}
		} else {
			resultParts = append(resultParts, relativeParts[i:]...)
			break
		}
	}

	return strings.Join(resultParts, "/")
}

// resolveReferences resolves $ref references in the schema
func (self *DatabaseProvider) resolveReferences(schema DefinitionParameterSchema) (DefinitionParameterSchema, error) {
	// Create a deep copy of the schema
	resolvedSchema := schema

	// Resolve references in properties
	for name, prop := range resolvedSchema.Properties {
		if prop.Ref != "" {
			// Parse the reference path
			// Format: "#/imports/s3Schema"
			parts := strings.Split(prop.Ref, "/")
			if len(parts) < 3 || parts[0] != "#" || parts[1] != "imports" {
				return resolvedSchema, fmt.Errorf("invalid reference: %s", prop.Ref)
			}

			importName := parts[2]
			importSchema, ok := schema.Imports[importName]
			if !ok {
				return resolvedSchema, fmt.Errorf("import not found: %s", importName)
			}

			// Convert to YAML and back to get the right structure
			importBytes, err := yaml.Marshal(importSchema)
			if err != nil {
				return resolvedSchema, err
			}

			var importedProp ParameterProperty
			if err := yaml.Unmarshal(importBytes, &importedProp); err != nil {
				return resolvedSchema, err
			}

			// Replace the reference with the actual schema
			resolvedSchema.Properties[name] = importedProp
		}
	}

	return resolvedSchema, nil
}
