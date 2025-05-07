package template_repo

import (
	"context"
	"sort"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/template"
)

func (self *TemplateRepository) GetAll(ctx context.Context) ([]*ent.Template, error) {
	templates, err := self.base.DB.Template.Query().
		Select(
			template.FieldID,
			template.FieldName,
			template.FieldVersion,
			template.FieldImmutable,
			template.FieldCreatedAt,
			template.FieldUpdatedAt,
		).All(ctx)
	if err != nil {
		return nil, err
	}

	// Create a map to store the newest version of each template name
	newestVersions := make(map[string]*ent.Template)

	// Iterate through all templates
	for _, tmpl := range templates {
		// If this name hasn't been seen yet, or if this version is newer than the stored one
		existingTmpl, exists := newestVersions[tmpl.Name]
		if !exists || tmpl.Version > existingTmpl.Version {
			newestVersions[tmpl.Name] = tmpl
		}
	}

	// Convert the map back to a slice
	result := make([]*ent.Template, 0, len(newestVersions))
	for _, tmpl := range newestVersions {
		result = append(result, tmpl)
	}

	// Sort by name
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}
