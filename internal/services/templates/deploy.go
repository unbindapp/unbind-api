package templates_service

import (
	"context"
	"errors"
	"fmt"

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

	// Find the first public service that's not a database for wildcard host generation
	var firstPublicService *schema.TemplateService
	for _, svc := range template.Definition.Services {
		if svc.IsPublic && svc.Type != schema.ServiceTypeDatabase && len(svc.Ports) > 0 {
			firstPublicService = &svc
			break
		}
	}

	// Validate template inputs and handle hosts
	validatedInputs := make(map[int]string)
	hostInputMap := make(map[int]v1.HostSpec) // Maps input ID to host spec

	// * Validate all non-host inputs
	for _, defInput := range template.Definition.Inputs {
		if defInput.Type == schema.InputTypeHost {
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
			} else if defInput.Required {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("input %s is required", defInput.Name))
			} else {
				continue
			}
		}

		validatedInputs[defInput.ID] = value
	}

	// * Validate or generate host inputs
	for _, defInput := range template.Definition.Inputs {
		if defInput.Type != schema.InputTypeHost {
			continue
		}

		// Get value from provided inputs or use default
		var hostValue string
		var exists bool

		for _, inputItem := range input.Inputs {
			if inputItem.ID == defInput.ID {
				hostValue = inputItem.Value
				exists = true
				break
			}
		}

		// If host value exists, validate it
		if exists {
			// Clean and validate the host
			cleanedHost, err := utils.CleanAndValidateHost(hostValue)
			if err != nil {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("invalid host for input %s: %v", defInput.Name, err))
			}
			hostValue = cleanedHost
		} else if defInput.Default != nil {
			// Use default if available
			hostValue = *defInput.Default
		} else if defInput.Required {
			// Required but missing, generate a wildcard host
			if firstPublicService == nil {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "no public service found for required host input")
			}

			kubernetesName, err := utils.GenerateSlug(firstPublicService.Name)
			if err != nil {
				return nil, err
			}

			// Generate wildcard host - NOT inside transaction
			var targetPort *int32
			if defInput.TargetPort != nil {
				port := int32(*defInput.TargetPort)
				targetPort = &port
			} else if len(firstPublicService.Ports) > 0 {
				targetPort = &firstPublicService.Ports[0].Port
			}

			generatedHost, err := self.generateWildcardHost(ctx, nil, kubernetesName, firstPublicService.Ports, targetPort)
			if err != nil {
				return nil, fmt.Errorf("failed to generate wildcard host: %w", err)
			}

			if generatedHost == nil {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "failed to generate wildcard host for required host input")
			}

			hostValue = generatedHost.Host
			// Store the host spec for later use
			hostInputMap[defInput.ID] = *generatedHost
		} else {
			// Not required and no value, skip
			continue
		}

		validatedInputs[defInput.ID] = hostValue

		// Create host spec if not already created
		if _, exists := hostInputMap[defInput.ID]; !exists {
			var port *int32
			if defInput.TargetPort != nil {
				p := int32(*defInput.TargetPort)
				port = &p
			} else if firstPublicService != nil && len(firstPublicService.Ports) > 0 {
				port = &firstPublicService.Ports[0].Port
			}

			hostInputMap[defInput.ID] = v1.HostSpec{
				Host: hostValue,
				Path: "/",
				Port: port,
			}
		}
	}

	// Parse and generate variables
	templater := templates.NewTemplater(self.cfg)
	generatedTemplate, err := templater.ResolveGeneratedVariables(&template.Definition, validatedInputs)
	if err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
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
	kubeNameMap := make(map[int]string)

	// Generate a launch ID
	templateInstanceID := uuid.New()

	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
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

				// Check metadata
				if dbVersion == nil && dbDefinition.DBVersion != "" {
					dbVersion = utils.ToPtr(dbDefinition.DBVersion)
				}
			}

			// Generate unique name
			kubernetesName, err := utils.GenerateSlug(templateService.Name)
			if err != nil {
				return err
			}
			kubeNameMap[templateService.ID] = kubernetesName

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
				})
			if err != nil {
				return fmt.Errorf("failed to create service: %w", err)
			}

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
				Ports:                   templateService.Ports,
				Hosts:                   hosts,
				Replicas:                utils.ToPtr[int32](1),
				Public:                  &templateService.IsPublic,
				Image:                   templateService.Image,
				CustomDefinitionVersion: utils.ToPtr(self.cfg.UnbindServiceDefVersion),
				PVCID:                   pvcID,
				PVCVolumeMountPath:      pvcMountPath,
				RunCommand:              templateService.RunCommand,
			}
			if templateService.Icon != "" {
				createInput.Icon = utils.ToPtr(templateService.Icon)
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
						}
					}

					// Standard variable references
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
				referenceInput = append(referenceInput, &models.VariableReferenceInputItem{
					Name: variableReference.TargetName,
					Sources: []schema.VariableReferenceSource{
						{
							Type:                 schema.VariableReferenceTypeVariable,
							SourceName:           sourceService.Name,
							SourceIcon:           sourceService.Edges.ServiceConfig.Icon,
							SourceID:             sourceService.ID,
							SourceType:           schema.VariableReferenceSourceTypeService,
							SourceKubernetesName: sourceService.KubernetesName,
							Key:                  variableReference.SourceName,
						},
					},
					Value: fmt.Sprintf("${%s.%s}", sourceService.KubernetesName, variableReference.SourceName),
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

func (self *TemplatesService) generateWildcardHost(ctx context.Context, tx repository.TxInterface, kubernetesName string, ports []schema.PortSpec, targetPort *int32) (*v1.HostSpec, error) {
	settings, err := self.repo.System().GetSystemSettings(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get system settings: %w", err)
	}

	if settings.WildcardBaseURL == nil {
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
