package service_service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	service_repo "github.com/unbindapp/unbind-api/internal/repositories/service"
	"github.com/unbindapp/unbind-api/internal/services/models"
	webhooks_service "github.com/unbindapp/unbind-api/internal/services/webooks"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
	"github.com/unbindapp/unbind-api/pkg/databases"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// CreateServiceInput defines the input for creating a new service
type CreateServiceInput struct {
	TeamID        uuid.UUID `format:"uuid" required:"true" json:"team_id"`
	ProjectID     uuid.UUID `format:"uuid" required:"true" json:"project_id"`
	EnvironmentID uuid.UUID `format:"uuid" required:"true" json:"environment_id"`
	Name          string    `required:"true" json:"name"`
	Description   string    `json:"description,omitempty"`

	// GitHub integration
	GitHubInstallationID *int64  `json:"github_installation_id,omitempty"`
	RepositoryOwner      *string `json:"repository_owner,omitempty"`
	RepositoryName       *string `json:"repository_name,omitempty"`

	// Configuration
	Type              schema.ServiceType    `required:"true" doc:"Type of service, e.g. 'github', 'docker-image'" json:"type"`
	Builder           schema.ServiceBuilder `required:"true" doc:"Builder of the service - docker, nixpacks, railpack" json:"builder"`
	Hosts             []v1.HostSpec         `json:"hosts,omitempty"`
	Ports             []schema.PortSpec     `json:"ports,omitempty"`
	Replicas          *int32                `minimum:"0" maximum:"10" json:"replicas,omitempty"`
	AutoDeploy        *bool                 `json:"auto_deploy,omitempty"`
	InstallCommand    *string               `json:"install_command,omitempty"`
	BuildCommand      *string               `json:"build_command,omitempty"`
	RunCommand        *string               `json:"run_command,omitempty"`
	IsPublic          *bool                 `json:"is_public,omitempty"`
	Image             *string               `json:"image,omitempty"`
	DockerfilePath    *string               `json:"dockerfile_path,omitempty" required:"false" doc:"Optional path to Dockerfile, if using docker builder"`
	DockerfileContext *string               `json:"dockerfile_context,omitempty" required:"false" doc:"Optional path to Dockerfile context, if using docker builder"`

	// Databases (special case)
	DatabaseType         *string                `json:"database_type,omitempty"`
	DatabaseConfig       *schema.DatabaseConfig `json:"database_config,omitempty"`
	S3BackupSourceID     *uuid.UUID             `json:"s3_backup_source_id,omitempty" format:"uuid"`
	S3BackupBucket       *string                `json:"s3_backup_bucket,omitempty"`
	BackupSchedule       *string                `json:"backup_schedule,omitempty" required:"false" doc:"Cron expression for the backup schedule, e.g. '0 0 * * *'"`
	BackupRetentionCount *int                   `json:"backup_retention,omitempty" required:"false" doc:"Number of base backups to retain, e.g. 3"`

	// PVC
	Volumes *[]schema.ServiceVolume `json:"volumes,omitempty" required:"false" doc:"Volumes to mount in the service"`

	// Health check
	HealthCheck *schema.HealthCheck `json:"health_check,omitempty" doc:"Health check configuration for the service"`

	// Variable mounts
	VariableMounts []*schema.VariableMount `json:"variable_mounts,omitempty" doc:"Mount variables as volumes"`

	// Init containers
	InitContainers []*schema.InitContainer `json:"init_containers,omitempty" doc:"Init containers to run before the main container"`
}

// CreateService creates a new service and its configuration
func (self *ServiceService) CreateService(ctx context.Context, requesterUserID uuid.UUID, input *CreateServiceInput, bearerToken string) (*models.ServiceResponse, error) {
	var err error
	var dbDefinition *databases.Definition
	var dbVersion *string
	var protectedVariables *[]string

	switch input.Type {
	case schema.ServiceTypeGithub:
		// Validate that if GitHub info is provided, all fields are set
		if input.GitHubInstallationID != nil {
			if input.RepositoryOwner == nil || input.RepositoryName == nil {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
					"GitHub repository owner, name must be provided together")
			}
		}
	case schema.ServiceTypeDockerimage:
		// Validate that if Docker image is provided, all fields are set
		if input.Image == nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
				"Docker image must be provided")
		}
		input.Builder = schema.ServiceBuilderDocker
	case schema.ServiceTypeDatabase:
		// Fixed protected variables for databases
		protectedVariables = &[]string{
			"DATABASE_USERNAME",
			"DATABASE_PASSWORD",
			"DATABASE_HOST",
			"DATABASE_PORT",
			"DATABASE_DEFAULT_DB_NAME",
			"DATABASE_URL",
			"DATABASE_HTTP_URL",
			"DATABASE_HTTP_PORT",
		}

		// Disallow pvc for database
		if input.Volumes != nil && len(*input.Volumes) > 0 {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
				"PVC is not supported for database services")
		}

		// Validate that if database is provided, name is set
		if input.DatabaseType == nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
				"Database name must be provided")
		}
		// Fetch the template
		dbDefinition, err = self.dbProvider.FetchDatabaseDefinition(ctx, self.cfg.UnbindServiceDefVersion, *input.DatabaseType)
		if err != nil {
			if errors.Is(err, databases.ErrDatabaseNotFound) {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound,
					fmt.Sprintf("Database %s not found", *input.DatabaseType))
			}
			return nil, err
		}

		// Nuke whatever they tell us for ports
		input.Ports = []schema.PortSpec{
			{
				Port:     int32(dbDefinition.Port),
				Protocol: utils.ToPtr(schema.ProtocolTCP),
			},
		}

		if input.DatabaseConfig != nil {
			// check version
			if input.DatabaseConfig.Version != "" {
				dbVersion = utils.ToPtr(input.DatabaseConfig.Version)
			}
			if input.DatabaseConfig.StorageSize == "" {
				input.DatabaseConfig.StorageSize = "1Gi" // Default to 1Gi
			} else {
				// Validate
				input.DatabaseConfig.StorageSize = utils.EnsureSuffix(input.DatabaseConfig.StorageSize, "Gi")
				_, err := resource.ParseQuantity(input.DatabaseConfig.StorageSize)
				if err != nil {
					return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
						fmt.Sprintf("Invalid storage size: %s", err))
				}
			}
		} else {
			input.DatabaseConfig = &schema.DatabaseConfig{
				StorageSize: "1Gi", // Default to 1Gi
			}
		}

		if dbVersion == nil {
			versionProperty, ok := dbDefinition.Schema.Properties["version"]
			if ok {
				dbVersionDefault, _ := versionProperty.Default.(string)
				if dbVersionDefault != "" {
					dbVersion = utils.ToPtr(dbVersionDefault)
				}
			}
		}

		imageProperty, ok := dbDefinition.Schema.Properties["dockerImage"]
		if ok {
			if image, ok := imageProperty.Default.(string); ok {
				input.Image = utils.ToPtr(image)
			}
		}

		// Check backup schedule
		if input.BackupSchedule != nil {
			if err := utils.ValidateCronExpression(*input.BackupSchedule); err != nil {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("invalid backup schedule: %s", err))
			}
		}

		input.Builder = schema.ServiceBuilderDatabase
	default:
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("received unsupported service type %s", input.Type))
	}

	// PVC validation, requires a path
	if input.Volumes != nil {
		for _, volume := range *input.Volumes {
			if !utils.IsValidUnixPath(volume.MountPath) {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Invalid PVC mount path")
			}
		}
	}

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

	// Git integrations
	var gitOwnerName *string

	// If GitHub integration is provided, verify repository access
	var analysisResult *sourceanalyzer.AnalysisResult
	var gitBranch *string
	if input.Type == schema.ServiceTypeGithub {
		// Get GitHub installation
		installation, err := self.repo.Github().GetInstallationByID(ctx, *input.GitHubInstallationID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "GitHub installation not found")
			}
			return nil, err
		}
		// Set owner
		gitOwnerName = utils.ToPtr(installation.AccountLogin)

		// Verify repository access
		canAccess, cloneUrl, defaultBranch, err := self.githubClient.VerifyRepositoryAccess(ctx, installation, *input.RepositoryOwner, *input.RepositoryName)
		if err != nil {
			log.Error("Error verifying repository access", "err", err)
			return nil, err
		}
		gitBranch = utils.ToPtr(defaultBranch)

		if !canAccess {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
				"Repository not accessible with the specified GitHub installation")
		}

		// Clone repository to infer information
		tmpDir, err := self.githubClient.CloneRepository(ctx, installation.GithubAppID, installation.ID, installation.Edges.GithubApp.PrivateKey, cloneUrl, fmt.Sprintf("refs/heads/%s", defaultBranch), "")
		if err != nil {
			log.Error("Error cloning repository", "err", err)
			return nil, err
		}
		defer os.RemoveAll(tmpDir)

		// Perform analysis
		analysisResult, err = sourceanalyzer.AnalyzeSourceCode(tmpDir)
		if err != nil {
			log.Error("Error analyzing source code", "err", err)
			return nil, err
		}
	} else if input.Type == schema.ServiceTypeDockerimage && len(input.Ports) == 0 {
		// Detect ports from image
		ports, _ := utils.GetExposedPortsFromRegistry(*input.Image)
		for _, port := range ports {
			// Split
			portSplit := strings.Split(port, "/")
			proto := schema.ProtocolTCP
			if len(portSplit) > 1 {
				// Check if the protocol is UDP
				if strings.EqualFold(portSplit[1], "udp") {
					proto = schema.ProtocolUDP
				}
			}
			portInt, err := strconv.Atoi(portSplit[0])
			if err != nil {
				log.Errorf("Failed to parse port %s: %v", port, err)
				continue
			}
			input.Ports = append(input.Ports, schema.PortSpec{
				Port:     int32(portInt),
				Protocol: utils.ToPtr(proto),
			})
		}
	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Check if PVC is in use by a service
	if input.Volumes != nil {
		for _, volume := range *input.Volumes {
			err = self.validatePVC(ctx, input.TeamID, input.ProjectID, input.EnvironmentID, volume.ID, project.Edges.Team.Namespace, client)
			if err != nil {
				return nil, err
			}
		}
	}

	// Verify backup sources (for databases)
	// Make sure we can read and write to the S3 bucket provided
	if input.Type == schema.ServiceTypeDatabase && input.S3BackupSourceID != nil && input.S3BackupBucket != nil {
		// Check if the S3 source exists
		s3Source, err := self.repo.S3().GetByID(ctx, *input.S3BackupSourceID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "S3 source not found")
			}
			return nil, err
		}

		if err := self.verifyS3Access(ctx, s3Source, *input.S3BackupBucket, project.Edges.Team.Namespace, client); err != nil {
			return nil, err
		}
	}

	// Create service and config in a transaction
	var service *ent.Service
	var serviceConfig *ent.ServiceConfig

	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		var provider *enum.Provider
		var framework *enum.Framework
		hosts := input.Hosts
		ports := input.Ports
		isPublic := input.IsPublic
		if analysisResult != nil {
			// Service core information
			if analysisResult.Provider != enum.UnknownProvider {
				provider = utils.ToPtr(analysisResult.Provider)
			}
			if analysisResult.Framework != enum.UnknownFramework {
				framework = utils.ToPtr(analysisResult.Framework)
			}

			// Default configuration information
			if len(ports) == 0 && analysisResult.Port != nil {
				ports = append(ports, schema.PortSpec{
					Port: int32(*analysisResult.Port),
				})
			}
		}

		// Generate unique name
		kubernetesName, err := utils.GenerateSlug(input.Name)
		if err != nil {
			return err
		}

		if len(ports) > 0 && input.IsPublic == nil {
			isPublic = utils.ToPtr(true)
		}

		if len(hosts) == 0 && input.IsPublic != nil && *isPublic && input.Type != schema.ServiceTypeDatabase && len(ports) > 0 {
			generatedHost, err := self.generateWildcardHost(ctx, tx, kubernetesName, ports)
			if err != nil {
				return fmt.Errorf("failed to generate wildcard host: %w", err)
			}
			if generatedHost == nil {
				isPublic = utils.ToPtr(false)
			} else {
				hosts = append(hosts, *generatedHost)
			}
		}

		// Validate hosts
		for _, host := range hosts {
			// Count domain collisions
			domainCount, err := self.repo.Service().CountDomainCollisons(ctx, tx, host.Host)
			if err != nil {
				return fmt.Errorf("failed to count domain collisions: %w", err)
			}
			if domainCount > 0 {
				return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("domain %s already in use", host.Host))
			}
		}

		if project == nil {
			log.Errorf("Project not found")
			return fmt.Errorf("Project not found")
		}
		if project.Edges.Team == nil {
			log.Errorf("Team not found")
			return fmt.Errorf("Team not found")
		}
		// Create kubernetes secrets
		secret, _, err := self.k8s.GetOrCreateSecret(ctx, kubernetesName, project.Edges.Team.Namespace, client)
		if err != nil {
			return fmt.Errorf("failed to create secret: %v", err)
		}

		// Create the service
		createService, err := self.repo.Service().Create(ctx, tx,
			&service_repo.CreateServiceInput{
				KubernetesName:       kubernetesName,
				ServiceType:          input.Type,
				Name:                 input.Name,
				Description:          input.Description,
				EnvironmentID:        input.EnvironmentID,
				GitHubInstallationID: input.GitHubInstallationID,
				GitRepository:        input.RepositoryName,
				GitRepositoryOwner:   gitOwnerName,
				KubernetesSecret:     secret.Name,
				Database:             input.DatabaseType,
				DatabaseVersion:      dbVersion,
			})
		if err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}
		service = createService

		// Create the service config
		createInput := &service_repo.MutateConfigInput{
			ServiceID:               service.ID,
			Builder:                 utils.ToPtr(input.Builder),
			Provider:                provider,
			Framework:               framework,
			GitBranch:               gitBranch,
			Ports:                   ports,
			Hosts:                   hosts,
			Replicas:                input.Replicas,
			AutoDeploy:              input.AutoDeploy,
			InstallCommand:          input.InstallCommand,
			BuildCommand:            input.BuildCommand,
			RunCommand:              input.RunCommand,
			Public:                  isPublic,
			Image:                   input.Image,
			DockerfilePath:          input.DockerfilePath,
			DockerfileContext:       input.DockerfileContext,
			CustomDefinitionVersion: utils.ToPtr(self.cfg.UnbindServiceDefVersion),
			DatabaseConfig:          input.DatabaseConfig,
			S3BackupSourceID:        input.S3BackupSourceID,
			S3BackupBucket:          input.S3BackupBucket,
			BackupSchedule:          input.BackupSchedule,
			BackupRetentionCount:    input.BackupRetentionCount,
			Volumes:                 input.Volumes,
			HealthCheck:             input.HealthCheck,
			VariableMounts:          input.VariableMounts,
			ProtectedVariables:      protectedVariables,
			InitContainers:          input.InitContainers,
		}

		serviceConfig, err = self.repo.Service().CreateConfig(ctx, tx, createInput)
		if err != nil {
			return fmt.Errorf("failed to create service config: %w", err)
		}
		service.Edges.ServiceConfig = serviceConfig
		return nil

	}); err != nil {
		return nil, err
	}

	// Trigger webhook
	go func() {
		event := schema.WebhookEventServiceCreated
		level := webhooks_service.WebhookLevelInfo

		// Get service with edges
		service, err := self.repo.Service().GetByID(context.Background(), service.ID)
		if err != nil {
			log.Errorf("Failed to get service %s: %v", service.ID.String(), err)
			return
		}

		// Construct URL
		basePath, _ := utils.JoinURLPaths(
			self.cfg.ExternalUIUrl,
			project.TeamID.String(),
			"project",
			project.ID.String(),
		)
		url := basePath + "?environment=" + input.EnvironmentID.String() +
			"&service=" + service.ID.String()
		// Get user
		user, err := self.repo.User().GetByID(context.Background(), requesterUserID)
		if err != nil {
			log.Errorf("Failed to get user %s: %v", requesterUserID.String(), err)
			return
		}
		data := webhooks_service.WebhookData{
			Title: "Service Created",
			Url:   url,
			Fields: []webhooks_service.WebhookDataField{
				{
					Name:  "Service",
					Value: service.Name,
				},
				{
					Name:  "Project & Environment",
					Value: fmt.Sprintf("%s > %s", service.Edges.Environment.Edges.Project.Name, service.Edges.Environment.Name),
				},
				{
					Name:  "Created By",
					Value: user.Email,
				},
			},
		}

		if len(service.Edges.ServiceConfig.Hosts) > 0 {
			data.Fields = append(data.Fields, webhooks_service.WebhookDataField{
				Name:  "Service URL",
				Value: fmt.Sprintf("https://%s", service.Edges.ServiceConfig.Hosts[0].Host),
			})
		}

		if err := self.webhookService.TriggerWebhooks(context.Background(), level, event, data); err != nil {
			log.Errorf("Failed to trigger webhook %s: %v", event, err)
		}
	}()

	resp := models.TransformServiceEntity(service)

	respArr := []*models.ServiceResponse{resp}

	if err := self.addPromMetricsToServiceVolumes(ctx, project.Edges.Team.Namespace, respArr); err != nil {
		log.Errorf("Failed to get PVC stats from prometheus: %v", err)
	}

	return respArr[0], nil
}
