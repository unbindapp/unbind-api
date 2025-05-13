package templates

import (
	"fmt"
	"strings"

	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
)

func (self *Templater) ResolveTemplate(template *schema.TemplateDefinition, inputs map[int]string, kubeNameMap map[int]string, namespace string) (*schema.TemplateDefinition, error) {
	resolved, err := self.resolveGeneratedVariables(template, inputs)
	if err != nil {
		return nil, err
	}

	// Build string replace map
	stringReplaceMap := make(map[string]string)
	stringReplaceMap["NAMESPACE"] = namespace
	for _, service := range resolved.Services {
		for _, variable := range service.Variables {
			if variable.Value != "" && (variable.Generator == nil || variable.Generator.Type != schema.GeneratorTypeStringReplace) {
				stringReplaceMap[fmt.Sprintf("SERVICE_%d_%s", service.ID, variable.Name)] = variable.Value
			}
		}
	}
	for k, v := range kubeNameMap {
		stringReplaceMap[fmt.Sprintf("SERVICE_%d_KUBE_NAME", k)] = v
	}
	for k, v := range inputs {
		stringReplaceMap[fmt.Sprintf("INPUT_%d_VALUE", k)] = v
	}

	// Execute string replace on all StringReplace variables
	for i := range resolved.Services {
		for j := range resolved.Services[i].Variables {
			variable := &resolved.Services[i].Variables[j]
			if variable.Generator != nil && variable.Generator.Type == schema.GeneratorTypeStringReplace {
				for k, v := range stringReplaceMap {
					variable.Value = strings.ReplaceAll(variable.Value, fmt.Sprintf("${%s}", k), v)
				}
			}
		}
		for j := range resolved.Services[i].VariableReferences {
			ref := &resolved.Services[i].VariableReferences[j]
			if ref.TemplateString != "" {
				for k, v := range stringReplaceMap {
					ref.TemplateString = strings.ReplaceAll(ref.TemplateString, fmt.Sprintf("${%s}", k), v)
				}
			}
		}
		if resolved.Services[i].DatabaseConfig != nil && resolved.Services[i].DatabaseConfig.InitDB != "" {
			for k, v := range resolved.Services[i].InitDBReplacers {
				toReplace, ok := stringReplaceMap[v]
				if !ok || toReplace == "" {
					toReplace = v
				}
				resolved.Services[i].DatabaseConfig.InitDB = strings.ReplaceAll(resolved.Services[i].DatabaseConfig.InitDB, k, toReplace)
			}
		}
	}

	return resolved, nil
}

// resolveGeneratedVariables creates a new TemplateDefinition with all generated variables resolved to their values
func (self *Templater) resolveGeneratedVariables(template *schema.TemplateDefinition, inputs map[int]string) (*schema.TemplateDefinition, error) {
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
			ID:                 svc.ID,
			DependsOn:          svc.DependsOn,
			Icon:               svc.Icon,
			Name:               svc.Name,
			Type:               svc.Type,
			Builder:            svc.Builder,
			DatabaseType:       svc.DatabaseType,
			DatabaseConfig:     svc.DatabaseConfig,
			Image:              svc.Image,
			Ports:              svc.Ports,
			IsPublic:           svc.IsPublic,
			HostInputIDs:       svc.HostInputIDs,
			RunCommand:         svc.RunCommand,
			Volumes:            svc.Volumes,
			Variables:          svc.Variables,
			VariableReferences: svc.VariableReferences,
			SecurityContext:    svc.SecurityContext,
			HealthCheck:        svc.HealthCheck,
			VariablesMounts:    svc.VariablesMounts,
			InitDBReplacers:    svc.InitDBReplacers,
		}

		// Initialize all slices if nil
		if resolvedService.DependsOn == nil {
			resolvedService.DependsOn = []int{}
		}
		if resolvedService.Ports == nil {
			resolvedService.Ports = []schema.PortSpec{}
		}
		if resolvedService.HostInputIDs == nil {
			resolvedService.HostInputIDs = []int{}
		}
		if resolvedService.VariableReferences == nil {
			resolvedService.VariableReferences = []schema.TemplateVariableReference{}
		}
		if resolvedService.Volumes == nil {
			resolvedService.Volumes = []schema.TemplateVolume{}
		}
		if resolvedService.Variables == nil {
			resolvedService.Variables = []schema.TemplateVariable{}
		}

		// Resolve variables
		var additionalVars []schema.TemplateVariable
		if len(resolvedService.Variables) > 0 {
			resolvedVars := make([]schema.TemplateVariable, len(resolvedService.Variables))
			for j, v := range resolvedService.Variables {
				resolvedVar := schema.TemplateVariable{
					Name: v.Name,
				}

				if v.Generator != nil {
					if v.Generator.Type == schema.GeneratorTypeStringReplace {
						// Preserve StringReplace generators and their values
						resolvedVar.Generator = v.Generator
						resolvedVar.Value = v.Value
					} else {
						res, err := v.Generator.Generate(inputs)
						if err != nil {
							return nil, fmt.Errorf("failed to generate value for %s: %w", v.Name, err)
						}
						resolvedVar.Value = res.GeneratedValue
						if v.Generator.Type == schema.GeneratorTypePasswordBcrypt {
							// Inject the raw password into a variable with the same name but without the _HASH suffix
							if strings.HasSuffix(v.Name, "_HASH") {
								additionalVars = append(additionalVars, schema.TemplateVariable{
									Name:  strings.TrimSuffix(v.Name, "_HASH") + "_PLAINTEXT",
									Value: res.PlainValue,
								})
							}
						}
						if v.Generator.Type == schema.GeneratorTypeJWT {
							// Parse all of them into additional variables based on the inputs
							resolvedVar.Value = fmt.Sprintf("generated-%s-%s-%s", v.Generator.JWTParams.SecretOutputKey, v.Generator.JWTParams.AnonOutputKey, v.Generator.JWTParams.ServiceOutputKey)
							additionalVars = append(additionalVars, []schema.TemplateVariable{
								{
									Name:  v.Generator.JWTParams.SecretOutputKey,
									Value: res.JWTValues[v.Generator.JWTParams.SecretOutputKey],
								},
								{
									Name:  v.Generator.JWTParams.AnonOutputKey,
									Value: res.JWTValues[v.Generator.JWTParams.AnonOutputKey],
								},
								{
									Name:  v.Generator.JWTParams.ServiceOutputKey,
									Value: res.JWTValues[v.Generator.JWTParams.ServiceOutputKey],
								},
							}...)
						}
					}
				} else {
					resolvedVar.Value = v.Value
				}
				resolvedVars[j] = resolvedVar
			}
			resolvedService.Variables = resolvedVars
		}

		// Append additional variables
		if len(additionalVars) > 0 {
			resolvedService.Variables = append(resolvedService.Variables, additionalVars...)
		}

		// Convert variables to map for template rendering
		variables := make(map[string]string)
		for _, v := range resolvedService.Variables {
			variables[v.Name] = v.Value
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
