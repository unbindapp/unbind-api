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

	// Parse generates variables
	templater := templates.NewTemplater(self.cfg)
	generatedTemplate, err := templater.ResolveGeneratedVariables(&template.Definition)
	if err != nil {
		return nil, err
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Create services
	var secretNames []string
	var newServices []*ent.Service
	kubeNameMap := make(map[int]string)
	dbServiceMap := make(map[int]*ent.Service)

	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		for _, templateService := range generatedTemplate.Services {
			// Fetch DB metadata (if a databsae)
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

			var hosts []v1.HostSpec

			if len(templateService.Ports) > 0 && templateService.IsPublic && templateService.Type != schema.ServiceTypeDatabase {
				generatedHost, err := self.generateWildcardHost(ctx, tx, kubernetesName, templateService.Ports)
				if err != nil {
					return fmt.Errorf("failed to generate wildcard host: %w", err)
				}
				if generatedHost == nil {
					// No wildcard host generated, set to false
					templateService.IsPublic = false
				} else {
					hosts = append(hosts, *generatedHost)
				}
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
			if _, err := self.k8s.UpsertSecretValues(ctx, secret.Name, project.Edges.Team.Namespace, secretData, client); err != nil {
				return fmt.Errorf("failed to update secret values: %v", err)
			}

			// Create the service
			createService, err := self.repo.Service().Create(ctx, tx,
				&service_repo.CreateServiceInput{
					KubernetesName:   kubernetesName,
					ServiceType:      templateService.Type,
					Name:             templateService.Name,
					EnvironmentID:    input.EnvironmentID,
					KubernetesSecret: secret.Name,
					Database:         templateService.DatabaseType,
					DatabaseVersion:  dbVersion,
					TemplateID:       utils.ToPtr(template.ID),
				})
			if err != nil {
				return fmt.Errorf("failed to create service: %w", err)
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
			// Map template to kubernetes name
			kubeNameMap[templateService.ID] = kubernetesName
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
							// Add moco in front
							key = fmt.Sprintf("moco-%s", key)
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
						Value: fmt.Sprintf("${%s.%s}", sourceService.KubernetesName, variableReference.SourceName),
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

func (self *TemplatesService) generateWildcardHost(ctx context.Context, tx repository.TxInterface, kubernetesName string, ports []schema.PortSpec) (*v1.HostSpec, error) {
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

	return &v1.HostSpec{
		Host: domain,
		Path: "/",
		Port: utils.ToPtr(ports[0].Port),
	}, nil
}
