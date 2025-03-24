package service_service

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer"
	"github.com/unbindapp/unbind-api/internal/sourceanalyzer/enum"
)

// CreateServiceInput defines the input for creating a new service
type CreateServiceInput struct {
	TeamID        uuid.UUID `validate:"required,uuid4" required:"true" json:"team_id"`
	ProjectID     uuid.UUID `validate:"required,uuid4" required:"true" json:"project_id"`
	EnvironmentID uuid.UUID `validate:"required,uuid4" required:"true" json:"environment_id"`
	DisplayName   string    `validate:"required" required:"true" json:"display_name"`
	Description   string    `json:"description,omitempty"`

	// GitHub integration
	GitHubInstallationID *int64  `json:"github_installation_id,omitempty"`
	RepositoryOwner      *string `json:"repository_owner,omitempty"`
	RepositoryName       *string `json:"repository_name,omitempty"`

	// Configuration
	Type       schema.ServiceType    `validate:"required" required:"true" doc:"Type of service, e.g. 'git', 'docker'" json:"type"`
	Builder    schema.ServiceBuilder `validate:"required" required:"true" doc:"Builder of the service - docker, nixpacks, railpack" json:"builder"`
	GitBranch  *string               `json:"git_branch,omitempty"`
	Hosts      []schema.HostSpec     `json:"hosts,omitempty"`
	Ports      []schema.PortSpec     `validate:"min=1,max=65535" json:"port,omitempty"`
	Replicas   *int32                `validate:"min=1,max=10" json:"replicas,omitempty"`
	AutoDeploy *bool                 `json:"auto_deploy,omitempty"`
	RunCommand *string               `json:"run_command,omitempty"`
	Public     *bool                 `json:"public,omitempty"`
	Image      *string               `json:"image,omitempty"`
}

// CreateService creates a new service and its configuration
func (self *ServiceService) CreateService(ctx context.Context, requesterUserID uuid.UUID, input *CreateServiceInput, bearerToken string) (*models.ServiceResponse, error) {
	// Validate input
	if err := validate.Validator().Struct(input); err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	// ! TODO - support docka
	if input.Type != schema.ServiceTypeGit {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "only git services supported")
	}

	// ! TODO - support docka
	if input.Builder != schema.ServiceBuilderNixpacks && input.Builder != schema.ServiceBuilderRailpack {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "only railpack and nixpacks builder supported")
	}

	// Validate that if GitHub info is provided, all fields are set
	if input.GitHubInstallationID != nil {
		if input.RepositoryOwner == nil || input.RepositoryName == nil || input.GitBranch == nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
				"GitHub repository owner, name, and branch must be provided together")
		}
	}

	// Check permissions
	permissionChecks := []permissions_repo.PermissionCheck{
		// Has permission to manage teams
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		// Has permission to manage this team
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   input.TeamID.String(),
		},
		// Has permission to manage projects
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   "*",
		},
		// Has permission to manage this project
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   input.ProjectID.String(),
		},
		// Has permission to manage this specific environment
		{
			Action:       permission.ActionCreate,
			ResourceType: permission.ResourceTypeEnvironment,
			ResourceID:   input.EnvironmentID.String(),
		},
	}

	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Verify inputs
	environment, project, err := self.VerifyInputs(ctx, input.TeamID, input.ProjectID, input.EnvironmentID)
	if err != nil {
		return nil, err
	}

	// If GitHub integration is provided, verify repository access
	var analysisResult *sourceanalyzer.AnalysisResult
	if input.GitHubInstallationID != nil {
		// Get GitHub installation
		installation, err := self.repo.Github().GetInstallationByID(ctx, *input.GitHubInstallationID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "GitHub installation not found")
			}
			return nil, err
		}

		// Verify repository access
		canAccess, cloneUrl, err := self.githubClient.VerifyRepositoryAccess(ctx, installation, *input.RepositoryOwner, *input.RepositoryName)
		if err != nil {
			log.Error("Error verifying repository access", "err", err)
			return nil, err
		}

		if !canAccess {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
				"Repository not accessible with the specified GitHub installation")
		}

		// Clone repository to infer information
		tmpDir, err := self.githubClient.CloneRepository(ctx, installation.GithubAppID, installation.ID, installation.Edges.GithubApp.PrivateKey, cloneUrl, fmt.Sprintf("refs/heads/%s", *input.GitBranch))
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

	}

	// Create kubernetes client
	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Create service and config in a transaction
	var service *ent.Service
	var serviceConfig *ent.ServiceConfig

	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		var provider *enum.Provider
		var framework *enum.Framework
		hosts := input.Hosts
		ports := input.Ports
		public := false
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
				public = true
			}
		}

		if len(hosts) == 0 && public {
			// Generate a subdomain
			domain, err := utils.GenerateSubdomain(input.DisplayName, environment.Name, self.cfg.ExternalURL)
			if err != nil {
				log.Warnf("Failed to generate subdomain: %v", err)
				public = false
			} else {
				// Check for collisons of the domain
				domainCount, err := self.repo.Service().CountDomainCollisons(ctx, tx, domain)
				if err != nil {
					log.Errorf("Failed to count domain collisions: %v", err)
					return err
				}
				if domainCount > 0 {
					// Re-generate with numerical suffix
					domain, err = utils.GenerateSubdomain(fmt.Sprintf("%s%d", input.DisplayName, domainCount), environment.Name, self.cfg.ExternalURL)
					if err != nil {
						log.Warnf("Failed to generate subdomain: %v", err)
						public = false
						domain = ""
					}
				}

				if domain != "" {
					hosts = append(hosts, schema.HostSpec{
						Host: domain,
						Path: "/",
						Port: utils.ToPtr(ports[0].Port),
					})
				}
			}
		}

		// Generate unique name
		name, err := utils.GenerateSlug(input.DisplayName)
		if err != nil {
			return err
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
		secret, _, err := self.k8s.GetOrCreateSecret(ctx, name, project.Edges.Team.Namespace, client)
		if err != nil {
			return fmt.Errorf("failed to create secret: %v", err)
		}

		// Create the service
		createService, err := self.repo.Service().Create(ctx, tx,
			name,
			input.DisplayName,
			input.Description,
			input.EnvironmentID,
			input.GitHubInstallationID,
			input.RepositoryName,
			secret.Name)
		if err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}
		service = createService

		// Create the service config
		serviceConfig, err = self.repo.Service().CreateConfig(ctx, tx,
			service.ID,
			input.Type,
			input.Builder,
			provider,
			framework,
			input.GitBranch,
			ports,
			hosts,
			input.Replicas,
			input.AutoDeploy,
			input.RunCommand,
			input.Public,
			input.Image)
		if err != nil {
			return fmt.Errorf("failed to create service config: %w", err)
		}
		service.Edges.ServiceConfig = serviceConfig
		return nil

	}); err != nil {
		return nil, err
	}

	return models.TransformServiceEntity(service), nil
}
