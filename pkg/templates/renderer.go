package templates

import (
	"fmt"
	"strings"

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

	// Ensure Inputs is initialized if nil
	if resolved.Inputs == nil {
		resolved.Inputs = []schema.TemplateInput{}
	}

	// Copy and resolve each service
	for i, svc := range template.Services {
		resolvedService := schema.TemplateService{
			ID:              svc.ID,
			Icon:            svc.Icon,
			Name:            svc.Name,
			Type:            svc.Type,
			Builder:         svc.Builder,
			DatabaseType:    svc.DatabaseType,
			DatabaseConfig:  svc.DatabaseConfig,
			Image:           svc.Image,
			IsPublic:        svc.IsPublic,
			RunCommand:      svc.RunCommand,
			SecurityContext: svc.SecurityContext,
		}

		// Initialize all slices if nil
		if svc.DependsOn == nil {
			resolvedService.DependsOn = []int{}
		} else {
			resolvedService.DependsOn = svc.DependsOn
		}

		if svc.Ports == nil {
			resolvedService.Ports = []schema.PortSpec{}
		} else {
			resolvedService.Ports = svc.Ports
		}

		if svc.HostInputIDs == nil {
			resolvedService.HostInputIDs = []int{}
		} else {
			resolvedService.HostInputIDs = svc.HostInputIDs
		}

		if svc.VariableReferences == nil {
			resolvedService.VariableReferences = []schema.TemplateVariableReference{}
		} else {
			resolvedService.VariableReferences = svc.VariableReferences
		}

		if svc.Volumes == nil {
			resolvedService.Volumes = []schema.TemplateVolume{}
		} else {
			resolvedService.Volumes = svc.Volumes
		}

		// Resolve variables
		var additionalVars []schema.TemplateVariable
		if svc.Variables == nil {
			resolvedService.Variables = []schema.TemplateVariable{}
		} else {
			resolvedService.Variables = make([]schema.TemplateVariable, len(svc.Variables))
			for j, v := range svc.Variables {
				resolvedVar := schema.TemplateVariable{
					Name: v.Name,
				}

				if v.Generator != nil {
					value, rawPassword, err := v.Generator.Generate(inputs)
					if err != nil {
						return nil, fmt.Errorf("failed to generate value for %s: %w", v.Name, err)
					}
					resolvedVar.Value = value
					if v.Generator.Type == schema.GeneratorTypePasswordBcrypt {
						// Inject the raw password into a variable with the same name but without the _HASH suffix
						if strings.HasSuffix(v.Name, "_HASH") {
							additionalVars = append(additionalVars, schema.TemplateVariable{
								Name:  strings.TrimSuffix(v.Name, "_HASH"),
								Value: rawPassword,
							})
						}
					}
				} else {
					resolvedVar.Value = v.Value
				}
				resolvedService.Variables[j] = resolvedVar
			}
		}

		// Append additional variables
		if len(additionalVars) > 0 {
			resolvedService.Variables = append(resolvedService.Variables, additionalVars...)
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
