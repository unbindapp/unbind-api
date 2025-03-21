package builds_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/permission"
	"github.com/unbindapp/unbind-api/internal/buildctl"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *BuildsService) CreateManualBuildJob(ctx context.Context, requesterUserId uuid.UUID, input *models.CreateBuildJobInput) (*models.BuildJobResponse, error) {
	if err := validate.Validator().Struct(input); err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserId, []permissions_repo.PermissionCheck{
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeSystem,
			ResourceID:   "*",
		},
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   "*",
		},
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeTeam,
			ResourceID:   input.TeamID.String(),
		},
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   "*",
		},
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeProject,
			ResourceID:   input.ProjectID.String(),
		},
		{
			Action:       permission.ActionManage,
			ResourceType: permission.ResourceTypeService,
			ResourceID:   input.ServiceID.String(),
		},
	}); err != nil {
		return nil, err
	}

	if err := self.validateInputs(ctx, input); err != nil {
		return nil, err
	}

	// Enqueue build job
	env, err := self.buildController.PopulateBuildEnvironment(ctx, input.ServiceID)
	if err != nil {
		return nil, err
	}

	job, err := self.buildController.EnqueueBuildJob(ctx, buildctl.BuildJobRequest{
		ServiceID:   input.ServiceID,
		Environment: env,
	})
	if err != nil {
		return nil, err
	}

	return models.TransformBuildJobEntity(job), nil

}
