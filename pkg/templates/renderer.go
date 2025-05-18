package templates

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

func (self *Templater) ResolveTemplate(template *schema.TemplateDefinition, inputs map[string]string, kubeNameMap map[string]string, namespace string) (*schema.TemplateDefinition, error) {
	resolved, err := self.resolveGeneratedVariables(template, inputs)
	if err != nil {
		return nil, err
	}

	// Resolve volumes
	resolved, err = self.resolveVolumes(resolved, inputs)
	if err != nil {
		return nil, err
	}

	// Resolve node ports
	resolved, err = self.resolveNodePorts(resolved, inputs)
	if err != nil {
		return nil, err
	}

	// Resolve database sizes
	resolved, err = self.resolveDatabaseSizes(resolved, inputs)
	if err != nil {
		return nil, err
	}

	// Build string replace map
	stringReplaceMap := make(map[string]string)
	stringReplaceMap["NAMESPACE"] = namespace
	for _, service := range resolved.Services {
		for _, variable := range service.Variables {
			if variable.Value != "" && (variable.Generator == nil || variable.Generator.Type != schema.GeneratorTypeStringReplace) {
				stringReplaceMap[fmt.Sprintf("%s_%s", strings.ToUpper(service.ID), variable.Name)] = variable.Value
			}
		}
	}
	for k, v := range kubeNameMap {
		stringReplaceMap[fmt.Sprintf("%s_KUBE_NAME", strings.ToUpper(k))] = v
	}
	for k, v := range inputs {
		stringReplaceMap[fmt.Sprintf("%s_VALUE", strings.ToUpper(k))] = v
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

// resolveDatabaseSizes resolves DatabaseSize inputs and attaches them to the relevant services
func (self *Templater) resolveDatabaseSizes(template *schema.TemplateDefinition, inputs map[string]string) (*schema.TemplateDefinition, error) {
	// We need to see if we have inputs of type DatabaseSize
	databaseSizeInputsMap := make(map[string]schema.TemplateInput)
	for _, input := range template.Inputs {
		if input.Type == schema.InputTypeDatabaseSize {
			databaseSizeInputsMap[input.ID] = input
		}
	}

	if len(databaseSizeInputsMap) == 0 {
		return template, nil
	}

	// Attach to relevant services
	for i := range template.Services {
		for _, inputID := range template.Services[i].InputIDs {
			if input, ok := databaseSizeInputsMap[inputID]; ok {
				// Check if the service has a database defined
				if template.Services[i].DatabaseConfig == nil {
					template.Services[i].DatabaseConfig = &schema.DatabaseConfig{}
				}

				size := inputs[inputID]
				// Verify size
				if _, err := utils.ValidateStorageQuantity(size); err != nil {
					return template, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("input %s is of type DatabaseSize but has an invalid size: %s", input.Name, err.Error()))
				}

				// Add the size to the database config
				template.Services[i].DatabaseConfig.StorageSize = size
			}
		}
	}

	return template, nil
}

// resolveNodePorts resolves NodePort inputs and attaches them to the relevant services
func (self *Templater) resolveNodePorts(template *schema.TemplateDefinition, inputs map[string]string) (*schema.TemplateDefinition, error) {
	// We need to see if we have inputs of type NodePort
	nodePortInputsMap := make(map[string]schema.TemplateInput)
	for _, input := range template.Inputs {
		if input.Type == schema.InputTypeGeneratedNodePort {
			nodePortInputsMap[input.ID] = input
		}
	}

	if len(nodePortInputsMap) == 0 {
		return template, nil
	}
	// {
	// 	IsNodePort:      true,
	// 	InputTemplateID: utils.ToPtr(2),
	// 	Protocol:        utils.ToPtr(schema.ProtocolUDP),
	// },
	// Attach to relevant services
	for i := range template.Services {
		for _, inputID := range template.Services[i].InputIDs {
			if input, ok := nodePortInputsMap[inputID]; ok {
				// Parse as int32
				asInt, err := strconv.Atoi(inputs[inputID])
				if err != nil {
					return template, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("input %s is of type NodePort but has an invalid value: %s", input.Name, err.Error()))
				}

				protocol := utils.ToPtr(schema.ProtocolTCP)
				if input.PortProtocol != nil {
					protocol = input.PortProtocol
				}

				template.Services[i].Ports = append(template.Services[i].Ports, schema.PortSpec{
					IsNodePort: true,
					NodePort:   utils.ToPtr(int32(asInt)),
					Port:       int32(asInt),
					Protocol:   protocol,
				})
			}
		}
	}

	return template, nil
}

// resolveVolumes resolves VolumeSize inputs and attaches them to the relevant services
func (self *Templater) resolveVolumes(template *schema.TemplateDefinition, inputs map[string]string) (*schema.TemplateDefinition, error) {
	// We need to see if we have inputs of type VolumeSize
	volumeSizeInputsMap := make(map[string]schema.TemplateInput)
	for _, input := range template.Inputs {
		if input.Type == schema.InputTypeVolumeSize {
			if input.Volume == nil {
				return template, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("input %s is of type VolumeSize but has no volume defined", input.Name))
			}
			volumeSizeInputsMap[input.ID] = input
		}
	}

	if len(volumeSizeInputsMap) == 0 {
		return template, nil
	}

	// Attach to relevant services
	for i := range template.Services {
		for _, inputID := range template.Services[i].InputIDs {
			if input, ok := volumeSizeInputsMap[inputID]; ok {
				// Check if the service has a volume defined
				if template.Services[i].Volumes == nil {
					template.Services[i].Volumes = []schema.TemplateVolume{}
				}

				size := inputs[inputID] + "Gi"
				// Verify size
				if _, err := utils.ValidateStorageQuantity(size); err != nil {
					return template, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("input %s is of type VolumeSize but has an invalid size: %s", input.Name, err.Error()))
				}

				// Add the volume to the service
				template.Services[i].Volumes = append(template.Services[i].Volumes, schema.TemplateVolume{
					Name:      input.Volume.Name,
					SizeGB:    size,
					MountPath: input.Volume.MountPath,
				})
			}
		}
	}

	return template, nil
}

// resolveGeneratedVariables creates a new TemplateDefinition with all generated variables resolved to their values
func (self *Templater) resolveGeneratedVariables(template *schema.TemplateDefinition, inputs map[string]string) (*schema.TemplateDefinition, error) {
	// Create a deep copy of the template
	resolved := &schema.TemplateDefinition{
		Name:        template.Name,
		Description: template.Description,
		Icon:        template.Icon,
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
			Icon:               svc.Icon,
			DependsOn:          svc.DependsOn,
			Name:               svc.Name,
			Type:               svc.Type,
			Builder:            svc.Builder,
			DatabaseType:       svc.DatabaseType,
			DatabaseConfig:     svc.DatabaseConfig,
			Image:              svc.Image,
			Ports:              svc.Ports,
			InputIDs:           svc.InputIDs,
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
			resolvedService.DependsOn = []string{}
		}
		if resolvedService.Ports == nil {
			resolvedService.Ports = []schema.PortSpec{}
		}
		if resolvedService.InputIDs == nil {
			resolvedService.InputIDs = []string{}
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
						// Set base domain for email generator
						if v.Generator.Type == schema.GeneratorTypeEmail {
							v.Generator.BaseDomain = self.cfg.ExternalUIUrl
						}
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
