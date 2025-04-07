package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// Enum for category name
type TemplateCategoryName string

const (
	TemplateCategoryDatabases TemplateCategoryName = "databases"
)

var allTemplateCategories = []TemplateCategoryName{
	TemplateCategoryDatabases,
}

func (self TemplateCategoryName) String() string {
	return string(self)
}

// Values provides list valid values for Enum.
func (TemplateCategoryName) Values() (kinds []string) {
	for _, s := range allTemplateCategories {
		kinds = append(kinds, string(s))
	}
	return
}

// Register enum in OpenAPI specification
// https://github.com/danielgtaylor/huma/issues/621
func (u TemplateCategoryName) Schema(r huma.Registry) *huma.Schema {
	if r.Map()["TemplateCategoryName"] == nil {
		schemaRef := r.Schema(reflect.TypeOf(""), true, "TemplateCategoryName")
		schemaRef.Title = "TemplateCategoryName"
		for _, v := range allTemplateCategories {
			schemaRef.Enum = append(schemaRef.Enum, string(v))
		}
		r.Map()["TemplateCategoryName"] = schemaRef
	}
	return &huma.Schema{Ref: "#/components/schemas/TemplateCategoryName"}
}

// TemplateCategory represents a category of templates
type TemplateCategory struct {
	Name      TemplateCategoryName `json:"name"`
	Templates []string             `json:"templates"`
}

// TemplateList represents a list of template categories
type TemplateList struct {
	Categories []TemplateCategory `json:"categories"`
}

// ListTemplates lists all available templates, optionally filtered by category
func (self *UnbindTemplateProvider) ListTemplates(ctx context.Context, tagVersion string, categoryFilter TemplateCategoryName) (*TemplateList, error) {
	// Base version URL
	baseURL := fmt.Sprintf(BaseTemplateURL, tagVersion)

	// Fetch the index file that contains all categories and templates
	indexURL := fmt.Sprintf("%s/index.json", baseURL)

	indexBytes, err := self.fetchURL(ctx, indexURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch template index: %w", err)
	}

	// Parse the index
	var templateList TemplateList
	if err := json.Unmarshal(indexBytes, &templateList); err != nil {
		return nil, fmt.Errorf("failed to parse template index: %w", err)
	}

	// Apply category filter if provided
	if categoryFilter != "" {
		filteredCategories := []TemplateCategory{}
		for _, category := range templateList.Categories {
			if strings.EqualFold(category.Name.String(), categoryFilter.String()) {
				filteredCategories = append(filteredCategories, category)
				break
			}
		}
		templateList.Categories = filteredCategories
	}

	return &templateList, nil
}
