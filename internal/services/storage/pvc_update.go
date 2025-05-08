package storage_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *StorageService) UpdatePVC(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.UpdatePVCInput) (*k8s.PVCInfo, error) {
	if input.Size == nil || input.Name == nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Size or Name are required")
	}

	// Validate permissions and parse inputs
	team, _, _, err := self.validatePermissionsAndParseInputs(ctx, schema.ActionEditor, requesterUserID, input.Type, input.TeamID, input.ProjectID, input.EnvironmentID)
	if err != nil {
		return nil, err
	}

	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Get the PVC
	pvc, err := self.k8s.GetPersistentVolumeClaim(ctx, team.Namespace, input.ID, client)
	if err != nil {
		return nil, err
	}

	switch input.Type {
	case models.PvcScopeTeam:
		if pvc.TeamID != input.TeamID {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "PVC not found")
		}
	case models.PvcScopeProject:
		if pvc.TeamID != input.TeamID || (pvc.ProjectID == nil || *pvc.ProjectID != input.ProjectID) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "PVC not found")
		}
	case models.PvcScopeEnvironment:
		if pvc.TeamID != input.TeamID || (pvc.ProjectID == nil || *pvc.ProjectID != input.ProjectID) || (pvc.EnvironmentID == nil || *pvc.EnvironmentID != input.EnvironmentID) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "PVC not found")
		}
	}

	// Size validation
	if input.Size != nil {
		// Parse size
		newSize, err := validateStorageQuantity(*input.Size)
		if err != nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
		}

		existingSize, err := validateStorageQuantity(pvc.Size)
		if err != nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
		}

		// Ensure new size is greater than existing size
		if newSize.Cmp(existingSize) < 0 {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "New size must be greater than existing size")
		}
	}

	return self.k8s.UpdatePersistentVolumeClaim(ctx,
		team.Namespace,
		input.ID,
		input.Name,
		input.Size,
		client,
	)
}
