package templates_service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	service_repo "github.com/unbindapp/unbind-api/internal/repositories/service"
	"github.com/unbindapp/unbind-api/internal/services/models"
	"github.com/unbindapp/unbind-api/pkg/databases"
	"github.com/unbindapp/unbind-api/pkg/templates"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

func (self *TemplatesService) DeployTemplate(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.TemplateDeployInput) ([]*models.ServiceResponse, error) {
	// Check permissions
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to manage teams
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeEnvironment,
			ResourceID:   input.EnvironmentID,
		},
	}

	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Verify inputs
	_, project, err := self.VerifyInputs(ctx, input.TeamID, input.ProjectID, input.EnvironmentID)
	if err != nil {
		return nil, err
	}

	// Get the template
	template, err := self.repo.Template().GetByID(ctx, input.TemplateID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Template not found")
		}
		return nil, err
	}

	// Validate template inputs and handle hosts
	validatedInputs := make(map[int]string)
	hostInputMap, validatedInputs, err :=
		self.resolveHostInputs(ctx, &template.Definition, input.Inputs)
	if err != nil {
		return nil, err
	}

	// * Validate all non-host inputs
	for _, defInput := range template.Definition.Inputs {
		if defInput.Type == schema.InputTypeHost || defInput.Type == schema.InputTypeNodePort {
			continue
		}

		// Get value from provided inputs or use default
		var value string
		var exists bool
		for _, input := range input.Inputs {
			if input.ID == defInput.ID {
				value = input.Value
				exists = true
				break
			}
		}

		if !exists {
			if defInput.Default != nil {
				value = *defInput.Default
			} else if defInput.Type == schema.InputTypePassword {
				// Generate a password
				pwd, err := utils.GenerateSecurePassword(32, true)
				if err != nil {
					return nil, err
				}
				value = pwd
			} else if defInput.Required {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("input %s is required", defInput.Name))
			} else {
				continue
			}
		}

		validatedInputs[defInput.ID] = value
	}

	// Generate node ports for node port inputs
	portMap := make(map[int]int32)
	portInputCount := 0
	for _, defInput := range template.Definition.Inputs {
		if defInput.Type == schema.InputTypeNodePort {
			portInputCount++
			// Ignore if input already exists
			for _, input := range input.Inputs {
				if input.ID == defInput.ID {
					// Parse as int32
					port, err := strconv.Atoi(input.Value)
					if err != nil {
						return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("invalid node port %s for input %d", input.Value, input.ID))
					}
					portMap[defInput.ID] = int32(port)
					break
				}
			}

			// Generate a port
			nodePort, err := self.k8s.GetUnusedNodePort(ctx)
			if err != nil {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("failed to generate node port: %v", err))
			}
			portMap[defInput.ID] = nodePort
			validatedInputs[defInput.ID] = strconv.Itoa(int(portMap[defInput.ID]))
		}
	}

	if len(portMap) != portInputCount {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "not all node ports were provided")
	}

	// Parse and generate variables
	kubeNameMap := make(map[int]string)

	// Generate kube name map
	for _, service := range template.Definition.Services {
		kubernetesName, err := utils.GenerateSlug(service.Name)
		if err != nil {
			return nil, err
		}
		kubeNameMap[service.ID] = kubernetesName
	}

	templater := templates.NewTemplater(self.cfg)
	generatedTemplate, err := templater.ResolveTemplate(&template.Definition, validatedInputs, kubeNameMap, project.Edges.Team.Namespace)
	if err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	// Make sure we have all our hosts
	for _, svc := range generatedTemplate.Services {
		for _, id := range svc.HostInputIDs {
			if _, ok := hostInputMap[id]; !ok {
				return nil, errdefs.NewCustomError(
					errdefs.ErrTypeInvalidInput,
					fmt.Sprintf("service %q references unresolved host input ID %d", svc.Name, id))
			}
		}
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Create services
	var secretNames []string
	var newServices []*ent.Service
	dbServiceMap := make(map[int]*ent.Service)

	// Generate a launch ID
	templateInstanceID := uuid.New()

	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		serviceGroup, err := self.repo.ServiceGroup().Create(ctx, tx, input.GroupName, input.EnvironmentID)
		if err != nil {
			return fmt.Errorf("failed to create service group: %w", err)
		}

		for _, templateService := range generatedTemplate.Services {
			// Fetch DB metadata (if a database)
			var dbVersion *string
			if templateService.Type == schema.ServiceTypeDatabase {
				// Fetch the template
				dbDefinition, err := self.dbProvider.FetchDatabaseDefinition(ctx, self.cfg.UnbindServiceDefVersion, *templateService.DatabaseType)
				if err != nil {
					if errors.Is(err, databases.ErrDatabaseNotFound) {
						return errdefs.NewCustomError(errdefs.ErrTypeNotFound,
							fmt.Sprintf("Database %s not found", *templateService.DatabaseType))
					}
					return err
				}

				// Nuke whatever they tell us for ports
				templateService.Ports = []schema.PortSpec{
					{
						Port:     int32(dbDefinition.Port),
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				}

				versionProperty, ok := dbDefinition.Schema.Properties["version"]
				if ok {
					dbVersionDefault, _ := versionProperty.Default.(string)
					if dbVersionDefault != "" {
						dbVersion = utils.ToPtr(dbVersionDefault)
					}
				}

				// ! TODO - validate validity
				if templateService.DatabaseConfig != nil && templateService.DatabaseConfig.Version != "" {
					dbVersion = utils.ToPtr(templateService.DatabaseConfig.Version)
				}

				if dbVersion == nil {
					return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
						fmt.Sprintf("Database version not found for %s", *templateService.DatabaseType))
				}

				if templateService.DatabaseConfig == nil {
					templateService.DatabaseConfig = &schema.DatabaseConfig{
						Version: *dbVersion,
					}
				}
			}

			kubernetesName, ok := kubeNameMap[templateService.ID]
			if !ok {
				return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("service %q not found", templateService.Name))
			}

			// Create kubernetes secrets
			secret, _, err := self.k8s.GetOrCreateSecret(ctx, kubernetesName, project.Edges.Team.Namespace, client)
			if err != nil {
				return fmt.Errorf("failed to create secret: %v", err)
			}
			secretNames = append(secretNames, secret.Name)

			// Add variables to the secret
			secretData := make(map[string][]byte)
			for _, variable := range templateService.Variables {
				secretData[variable.Name] = []byte(variable.Value)
			}

			// Resolve any references that should be treated as local
			for _, ref := range templateService.VariableReferences {
				if ref.IsHost && ref.ResolveAsNormalVariable {
					// See if we have a kubeNameMap for the source ID
					sourceKubeName, ok := kubeNameMap[ref.SourceID]
					if !ok {
						return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("source service not found for local variable reference %s", ref.TargetName))
					}
					// Add the reference to the secret data
					secretData[ref.TargetName] = fmt.Appendf(nil, "%s.%s", sourceKubeName, project.Edges.Team.Namespace)
				}
			}

			if _, err := self.k8s.UpsertSecretValues(ctx, secret.Name, project.Edges.Team.Namespace, secretData, client); err != nil {
				return fmt.Errorf("failed to update secret values: %v", err)
			}

			// Create the service
			createService, err := self.repo.Service().Create(ctx, tx,
				&service_repo.CreateServiceInput{
					KubernetesName:     kubernetesName,
					ServiceType:        templateService.Type,
					Name:               templateService.Name,
					EnvironmentID:      input.EnvironmentID,
					KubernetesSecret:   secret.Name,
					Database:           templateService.DatabaseType,
					DatabaseVersion:    dbVersion,
					TemplateID:         utils.ToPtr(template.ID),
					TemplateInstanceID: utils.ToPtr(templateInstanceID),
					ServiceGroupID:     utils.ToPtr(serviceGroup.ID),
				})
			if err != nil {
				return fmt.Errorf("failed to create service: %w", err)
			}
			createService.Edges.ServiceGroup = serviceGroup

			// Create volumes
			var pvcID *string
			var pvcMountPath *string
			for _, volume := range templateService.Volumes {
				// Get size from input ID
				size, exists := validatedInputs[volume.Size.FromInputID]
				if !exists {
					return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "size input not found")
				}
				// Parse size
				_, err = utils.ValidateStorageQuantity(size)
				if err != nil {
					return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
				}

				// Build labels to set
				labels := map[string]string{
					"unbind-team":        input.TeamID.String(),
					"unbind-project":     input.ProjectID.String(),
					"unbind-environment": input.EnvironmentID.String(),
				}

				//  Generate a name
				pvcName, err := utils.GenerateSlug(volume.Name)
				if err != nil {
					return err
				}

				// Get the PVCs
				pvc, err := self.k8s.CreatePersistentVolumeClaim(ctx,
					project.Edges.Team.Namespace,
					pvcName,
					volume.Name,
					labels,
					size,
					[]corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					nil,
					client,
				)
				if err != nil {
					return err
				}
				pvcID = utils.ToPtr(pvc.ID)
				pvcMountPath = utils.ToPtr(volume.MountPath)
			}

			var ports []schema.PortSpec
			for _, port := range templateService.Ports {
				// Check if we have an input value
				if port.InputTemplateID != nil {
					// Get the value from the input map
					portValue := portMap[*port.InputTemplateID]
					if port.IsNodePort {
						port.NodePort = utils.ToPtr(portValue)
					}
					port.Port = portValue
				}
				ports = append(ports, port)
			}

			var hosts []v1.HostSpec
			for _, hostInputID := range templateService.HostInputIDs {
				hostSpec, exists := hostInputMap[hostInputID]
				if !exists {
					return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "host input not found")
				}
				hosts = append(hosts, hostSpec)
			}

			// Create the service config
			createInput := &service_repo.MutateConfigInput{
				ServiceID:               createService.ID,
				Builder:                 utils.ToPtr(templateService.Builder),
				Ports:                   ports,
				Hosts:                   hosts,
				Replicas:                utils.ToPtr[int32](1),
				Public:                  &templateService.IsPublic,
				DatabaseConfig:          templateService.DatabaseConfig,
				Image:                   templateService.Image,
				CustomDefinitionVersion: utils.ToPtr(self.cfg.UnbindServiceDefVersion),
				PVCID:                   pvcID,
				PVCVolumeMountPath:      pvcMountPath,
				RunCommand:              templateService.RunCommand,
				SecurityContext:         templateService.SecurityContext,
				HealthCheck:             templateService.HealthCheck,
				VariableMounts:          templateService.VariablesMounts,
			}

			serviceConfig, err := self.repo.Service().CreateConfig(ctx, tx, createInput)
			if err != nil {
				return fmt.Errorf("failed to create service config: %w", err)
			}
			createService.Edges.ServiceConfig = serviceConfig

			// Append
			newServices = append(newServices, createService)
			dbServiceMap[templateService.ID] = createService
		}

		// Resolve variable references
		for _, templateService := range generatedTemplate.Services {
			referenceInput := []*models.VariableReferenceInputItem{}
			for _, variableReference := range templateService.VariableReferences {
				if variableReference.ResolveAsNormalVariable {
					continue
				}
				// Get source
				sourceService := dbServiceMap[variableReference.SourceID]
				if sourceService == nil {
					log.Error("failed to find service for variable reference", "serviceID", variableReference.SourceID, "template", templateService.Name)
					return fmt.Errorf("failed to find service for variable reference: %w", err)
				}

				// Deal with is_host first, not really a reference just resolved on the fly
				if variableReference.IsHost {
					// Determine the host key (it may not exist yet so we can't DNS lookup)
					key := sourceService.KubernetesName

					if sourceService.Type == schema.ServiceTypeDatabase && sourceService.Database != nil {
						// Special DB cases
						switch *sourceService.Database {
						case "mysql":
							// Moco primary instance
							key = fmt.Sprintf("moco-%s-primary", key)
						case "redis":
							key = fmt.Sprintf("%s-headless", key)
						case "clickhouse":
							key = fmt.Sprintf("clickhouse-%s", key)
						}
					}

					// Host Refs
					referenceInput = append(referenceInput, &models.VariableReferenceInputItem{
						Name: variableReference.TargetName,
						Sources: []schema.VariableReferenceSource{
							{
								Type:                 schema.VariableReferenceTypeInternalEndpoint,
								SourceName:           sourceService.Name,
								SourceIcon:           sourceService.Edges.ServiceConfig.Icon,
								SourceID:             sourceService.ID,
								SourceType:           schema.VariableReferenceSourceTypeService,
								SourceKubernetesName: sourceService.KubernetesName,
								Key:                  key,
							},
						},
						Value: fmt.Sprintf("${%s.%s}", sourceService.KubernetesName, key),
					})

					continue
				}

				// Standard variable references
				value := fmt.Sprintf("${%s.%s}", sourceService.KubernetesName, variableReference.SourceName)
				sources := []schema.VariableReferenceSource{
					{
						Type:                 schema.VariableReferenceTypeVariable,
						SourceName:           sourceService.Name,
						SourceIcon:           sourceService.Edges.ServiceConfig.Icon,
						SourceID:             sourceService.ID,
						SourceType:           schema.VariableReferenceSourceTypeService,
						SourceKubernetesName: sourceService.KubernetesName,
						Key:                  variableReference.SourceName,
					},
				}
				if variableReference.TemplateString != "" {
					// Replace the key with the right one
					value = strings.ReplaceAll(variableReference.TemplateString, fmt.Sprintf("${%s}", variableReference.SourceName), fmt.Sprintf("${%s.%s}", sourceService.KubernetesName, variableReference.SourceName))
					if len(variableReference.AdditionalTemplateSources) > 0 {
						for _, additionalSource := range variableReference.AdditionalTemplateSources {
							sources = append(sources, schema.VariableReferenceSource{
								Type:                 schema.VariableReferenceTypeVariable,
								SourceName:           sourceService.Name,
								SourceIcon:           sourceService.Edges.ServiceConfig.Icon,
								SourceID:             sourceService.ID,
								SourceType:           schema.VariableReferenceSourceTypeService,
								SourceKubernetesName: sourceService.KubernetesName,
								Key:                  additionalSource,
							})

							value = strings.ReplaceAll(value, fmt.Sprintf("${%s}", additionalSource), fmt.Sprintf("${%s.%s}", sourceService.KubernetesName, additionalSource))
						}
					}
				}
				referenceInput = append(referenceInput, &models.VariableReferenceInputItem{
					Name:    variableReference.TargetName,
					Sources: sources,
					Value:   value,
				})
			}

			if len(referenceInput) > 0 {
				targetService := dbServiceMap[templateService.ID]
				if targetService == nil {
					log.Error("failed to find service for variable reference", "targetID", templateService.ID)
					return fmt.Errorf("failed to find service for variable reference: %w", err)
				}
				_, err := self.repo.Variables().UpdateReferences(ctx, tx, models.VariableUpdateBehaviorUpsert, targetService.ID, referenceInput)
				if err != nil {
					if ent.IsConstraintError(err) {
						return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Variable reference already exists")
					}
					return err
				}
			}
		}

		return nil
	}); err != nil {
		// Attempt to delete the created secrets
		for _, secretName := range secretNames {
			if err := self.k8s.DeleteSecret(ctx, secretName, project.Edges.Team.Namespace, client); err != nil {
				log.Error("failed to delete secret", "secretName", secretName, "error", err)
			}
		}

		return nil, err
	}

	// Deploy all services without variable references
	for _, service := range generatedTemplate.Services {
		// Get the service from our map
		dbService := dbServiceMap[service.ID]

		// Populate build environment
		env, err := self.deployCtl.PopulateBuildEnvironment(ctx, dbService.ID, nil)
		if err != nil {
			return nil, err
		}

		// Create deployment request
		deployReq := deployctl.DeploymentJobRequest{
			ServiceID:   dbService.ID,
			Environment: env,
			Source:      schema.DeploymentSourceManual,
		}

		// If service has dependencies, add to dependent queue
		if len(service.VariableReferences) > 0 {
			err = self.deployCtl.EnqueueDependentDeployment(ctx, deployReq)
		} else {
			// Otherwise deploy immediately
			_, err = self.deployCtl.EnqueueDeploymentJob(ctx, deployReq)
		}

		if err != nil {
			return nil, err
		}
	}

	return models.TransformServiceEntities(newServices), nil
}

// returns: map[inputID]HostSpec, map[inputID]string (value to hand to the templater)
func (self *TemplatesService) resolveHostInputs(
	ctx context.Context,
	tmpl *schema.TemplateDefinition,
	rawInputs []models.TemplateInputValue,
) (map[int]v1.HostSpec, map[int]string, error) {
	hostSpecByID := make(map[int]v1.HostSpec)
	valueByID := make(map[int]string)

	for _, in := range tmpl.Inputs {
		if in.Type != schema.InputTypeHost {
			continue
		}

		// 1. find a value (provided → default → "")
		var hostVal string
		provided := false
		for _, u := range rawInputs {
			if u.ID == in.ID {
				hostVal, provided = u.Value, true
				break
			}
		}
		if !provided && in.Default != nil {
			hostVal = *in.Default
		}

		// 2. if we still have nothing AND it's required → try generate
		if hostVal == "" && in.Required {
			// generateWildcardHost needs *one* service / port context – first match wins
			genForSvc := (*schema.TemplateService)(nil)
			var targetPort *int32
			for _, svc := range tmpl.Services {
				for _, id := range svc.HostInputIDs {
					if id == in.ID {
						genForSvc = &svc
						break
					}
				}
				if genForSvc != nil {
					break
				}
			}
			if genForSvc == nil {
				return nil, nil, errdefs.NewCustomError(
					errdefs.ErrTypeInvalidInput,
					fmt.Sprintf("no service refers to required host input %q", in.Name))
			}

			if in.TargetPort != nil {
				p := int32(*in.TargetPort)
				targetPort = &p
			} else if len(genForSvc.Ports) > 0 {
				targetPort = &genForSvc.Ports[0].Port
			}
			kubename, _ := utils.GenerateSlug(genForSvc.Name)
			spec, err := self.generateWildcardHost(ctx, nil, kubename, genForSvc.Ports, targetPort)
			if err != nil {
				return nil, nil, err
			}
			hostVal = spec.Host
			hostSpecByID[in.ID] = *spec // we already have the full spec
			valueByID[in.ID] = hostVal
			continue
		}

		// 3. we have a hostVal from user / default – clean & store
		if hostVal != "" {
			cleaned, err := utils.CleanAndValidateHost(hostVal)
			if err != nil {
				return nil, nil, errdefs.NewCustomError(
					errdefs.ErrTypeInvalidInput,
					fmt.Sprintf("invalid host for input %q: %v", in.Name, err))
			}
			hostVal = cleaned
			// pick a port (targetPort or first port of *any* svc using this input)
			var port *int32
			if in.TargetPort != nil {
				p := int32(*in.TargetPort)
				port = &p
			} else {
				for _, svc := range tmpl.Services {
					for _, id := range svc.HostInputIDs {
						if id == in.ID && len(svc.Ports) > 0 {
							port = &svc.Ports[0].Port
							break
						}
					}
					if port != nil {
						break
					}
				}
			}
			hostSpecByID[in.ID] = v1.HostSpec{Host: hostVal, Path: "/", Port: port}
			valueByID[in.ID] = hostVal
		}
	}
	return hostSpecByID, valueByID, nil
}

func (self *TemplatesService) generateWildcardHost(ctx context.Context, tx repository.TxInterface, kubernetesName string, ports []schema.PortSpec, targetPort *int32) (*v1.HostSpec, error) {
	settings, err := self.repo.System().GetSystemSettings(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get system settings: %w", err)
	}

	if settings.WildcardBaseURL == nil || *settings.WildcardBaseURL == "" {
		return nil, nil // No wildcard base URL configured
	}

	domain, err := utils.GenerateSubdomain(kubernetesName, *settings.WildcardBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate subdomain: %w", err)
	}

	domainCount, err := self.repo.Service().CountDomainCollisons(ctx, tx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to count domain collisions: %w", err)
	}

	if domainCount > 0 {
		domain, err = utils.GenerateSubdomain(fmt.Sprintf("%s-%d", kubernetesName, domainCount), *settings.WildcardBaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to generate subdomain with suffix: %w", err)
		}
	}

	// Use targetPort if provided, otherwise use the first port from ports
	port := targetPort
	if port == nil && len(ports) > 0 {
		port = utils.ToPtr(ports[0].Port)
	}

	return &v1.HostSpec{
		Host: domain,
		Path: "/",
		Port: port,
	}, nil
}
