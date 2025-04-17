package variables_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/repositories/repositories"
)

// Integrate variables management with internal permissions and kubernetes RBAC
type VariablesService struct {
	repo repositories.RepositoriesInterface
	k8s  *k8s.KubeClient
}

func NewVariablesService(repo repositories.RepositoriesInterface, k8s *k8s.KubeClient) *VariablesService {
	return &VariablesService{
		repo: repo,
		k8s:  k8s,
	}
}

func (self *VariablesService) validateInputs(ctx context.Context, teamID, projectID, environmentID, serviceID uuid.UUID) (*ent.Team, *ent.Project, *ent.Environment, *ent.Service, error) {
	// Get available variable references
	service, err := self.repo.Service().GetByID(ctx, serviceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
		}
		return nil, nil, nil, nil, err
	}

	environment := service.Edges.Environment
	project := environment.Edges.Project
	team := project.Edges.Team
	if team.ID != teamID ||
		project.ID != projectID ||
		environment.ID != environmentID {
		return nil, nil, nil, nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "Service not found")
	}

	return team, project, environment, service, nil
}
