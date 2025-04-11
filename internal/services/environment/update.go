package environment_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/validate"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

type UpdateEnvironmentInput struct {
	TeamID        uuid.UUID `json:"team_id" validate:"required,uuid4" required:"true"`
	ProjectID     uuid.UUID `json:"project_id" validate:"required,uuid4" required:"true"`
	EnvironmentID uuid.UUID `json:"environment_id" validate:"required,uuid4" required:"true"`
	Name          *string   `json:"name"`
	Description   *string   `json:"description"`
}

func (self *EnvironmentService) UpdateEnvironment(ctx context.Context, requesterUserID uuid.UUID, input *UpdateEnvironmentInput) (*models.EnvironmentResponse, error) {
	// Validate input
	if err := validate.Validator().Struct(input); err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	permissionChecks := []permissions_repo.PermissionCheck{
		// Project editor can create environments
		{
			Action:       schema.ActionEditor,
			ResourceType: schema.ResourceTypeEnvironment,
			ResourceID:   input.EnvironmentID,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Verify inputs
	_, environment, err := self.VerifyInputs(ctx, input.TeamID, input.ProjectID, input.EnvironmentID)
	if err != nil {
		return nil, err
	}

	// Update the environment
	updated, err := self.repo.Environment().Update(ctx, environment.ID, input.Name, input.Description)
	if err != nil {
		return nil, err
	}

	// Convert to response
	resp := models.TransformEnvironmentEntity(updated)

	// Summarizes services
	counts, providerSummaries, err := self.repo.Service().SummarizeServices(ctx, []uuid.UUID{environment.ID})
	if err != nil {
		return nil, err
	}
	resp.ServiceCount, _ = counts[environment.ID]
	resp.ServiceIcons, _ = providerSummaries[environment.ID]

	return resp, nil
}
