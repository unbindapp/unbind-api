package deployments_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/deployctl"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// Integrate builds management with internal permissions and kubernetes RBAC
type DeploymentService struct {
	repo                 repositories.RepositoriesInterface
	deploymentController *deployctl.DeploymentController
}

func NewDeploymentService(repo repositories.RepositoriesInterface, deploymentController *deployctl.DeploymentController) *DeploymentService {
	return &DeploymentService{
		repo:                 repo,
		deploymentController: deploymentController,
	}
}

func (self *DeploymentService) validateInputs(ctx context.Context, input models.DeploymentInputRequirements) error {
	// Get team
	team, err := self.repo.Team().GetByID(ctx, input.GetTeamID())
	if err != nil {
		if ent.IsNotFound(err) {
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound, input.GetTeamID().String())
		}
		return err
	}

	// Get project
	project, err := self.repo.Project().GetByID(ctx, input.GetProjectID())
	if err != nil {
		if ent.IsNotFound(err) {
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound, input.GetProjectID().String())
		}
		return err
	}

	if project.TeamID != team.ID {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "project does not belong to team")
	}

	var environmentID *uuid.UUID
	for _, env := range project.Edges.Environments {
		if env.ID == input.GetEnvironmentID() {
			environmentID = utils.ToPtr(env.ID)
			break
		}
	}

	if environmentID == nil {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "environment does not belong to project")
	}

	// Get service
	service, err := self.repo.Service().GetByID(ctx, input.GetServiceID())
	if err != nil {
		if ent.IsNotFound(err) {
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound, "service not found")
		}
		return err
	}

	if service.EnvironmentID != *environmentID {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "service does not belong to environment")
	}
	return nil
}
