package templates

import (
	"fmt"

	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
)

// ResolveGeneratedVariables creates a new TemplateDefinition with all generated variables resolved to their values
func (self *Templater) ResolveGeneratedVariables(template *schema.TemplateDefinition, inputs map[int]string) (*schema.TemplateDefinition, error) {
	// Create a deep copy of the template
	resolved := &schema.TemplateDefinition{
		Name:        template.Name,
		Description: template.Description,
		Version:     template.Version,
		Services:    make([]schema.TemplateService, len(template.Services)),
		Inputs:      template.Inputs,
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
			Volumes:            svc.Volumes,
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
				value, err := v.Generator.Generate(inputs)
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

// ProcessTemplateInputs processes the template inputs and returns a map of resolved values
func (self *Templater) ProcessTemplateInputs(template *schema.TemplateDefinition, values map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	for _, input := range template.Inputs {
		// Skip host type as it's handled separately
		if input.Type == schema.InputTypeHost {
			continue
		}

		// Get value from provided values or use default
		value, exists := values[input.Name]
		if !exists {
			if input.Default != nil {
				value = *input.Default
			} else if input.Required {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("input %s is required", input.Name))
			} else {
				continue
			}
		}

		result[input.Name] = value
	}

	return result, nil
}
