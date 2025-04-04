package deployments_service

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	"github.com/unbindapp/unbind-api/internal/infrastructure/loki"
	"github.com/unbindapp/unbind-api/internal/integrations/github"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Integrate builds management with internal permissions and kubernetes RBAC
type DeploymentService struct {
	repo                 repositories.RepositoriesInterface
	deploymentController *deployctl.DeploymentController
	githubClient         *github.GithubClient
	lokiQuerier          *loki.LokiLogQuerier
}

func NewDeploymentService(repo repositories.RepositoriesInterface, deploymentController *deployctl.DeploymentController, githubClient *github.GithubClient, lokiQuerier *loki.LokiLogQuerier) *DeploymentService {
	return &DeploymentService{
		repo:                 repo,
		deploymentController: deploymentController,
		githubClient:         githubClient,
		lokiQuerier:          lokiQuerier,
	}
}

func (self *DeploymentService) validateInputs(ctx context.Context, input models.DeploymentInputRequirements) (*ent.Service, error) {
	// Get team
	team, err := self.repo.Team().GetByID(ctx, input.GetTeamID())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, input.GetTeamID().String())
		}
		return nil, err
	}

	// Validate project
	var project *ent.Project
	for _, proj := range team.Edges.Projects {
		if proj.ID == input.GetProjectID() {
			project = proj
			break
		}
	}

	if project == nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "project does not belong to team")
	}

	// Validate environment
	var environment *ent.Environment
	for _, env := range project.Edges.Environments {
		if env.ID == input.GetEnvironmentID() {
			environment = env
			break
		}
	}

	if environment == nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "environment does not belong to project")
	}

	// Validate service
	var service *ent.Service
	for _, svc := range environment.Edges.Services {
		if svc.ID == input.GetServiceID() {
			service = svc
			break
		}
	}

	if service == nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "service does not belong to environment")
	}

	return service, nil
}
