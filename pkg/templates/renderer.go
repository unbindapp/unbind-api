package templates

import (
	"fmt"

	"github.com/unbindapp/unbind-api/ent/schema"
)

// ResolveGeneratedVariables creates a new TemplateDefinition with all generated variables resolved to their values
func (self *Templater) ResolveGeneratedVariables(template *schema.TemplateDefinition) (*schema.TemplateDefinition, error) {
	// Create a deep copy of the template
	resolved := &schema.TemplateDefinition{
		Name:        template.Name,
		Description: template.Description,
		Version:     template.Version,
		Services:    make([]schema.TemplateService, len(template.Services)),
	}

	// Copy and resolve each service
	for i, svc := range template.Services {
		resolvedService := schema.TemplateService{
			ID:                 svc.ID,
			Icon:               svc.Icon,
			DependsOn:          svc.DependsOn,
			Name:               svc.Name,
			Type:               svc.Type,
			Builder:            svc.Builder,
			DatabaseType:       svc.DatabaseType,
			Image:              svc.Image,
			Ports:              svc.Ports,
			IsPublic:           svc.IsPublic,
			VariableReferences: svc.VariableReferences,
		}

		// Resolve variables
		resolvedService.Variables = make([]schema.TemplateVariable, len(svc.Variables))
		for j, v := range svc.Variables {
			resolvedVar := schema.TemplateVariable{
				Name: v.Name,
			}

			if v.Generator != nil {
				// Set base domain for email generator
				if v.Generator.Type == schema.GeneratorTypeEmail {
					v.Generator.BaseDomain = self.cfg.ExternalUIUrl
				}
				value, err := v.Generator.Generate()
				if err != nil {
					return nil, fmt.Errorf("failed to generate value for %s: %w", v.Name, err)
				}
				resolvedVar.Value = value
			} else {
				resolvedVar.Value = v.Value
			}

			resolvedService.Variables[j] = resolvedVar
		}

		resolved.Services[i] = resolvedService
	}

	return resolved, nil
}
