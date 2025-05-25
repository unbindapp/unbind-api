package storage_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/models"
	repository "github.com/unbindapp/unbind-api/internal/repositories"
)

func (self *StorageService) UpdatePVC(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.UpdatePVCInput) (*models.PVCInfo, error) {
	if input.CapacityGB == nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "Size is required")
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
	var newCapacity *string
	if input.CapacityGB != nil {
		// Parse size
		newCapacity = utils.ToPtr(fmt.Sprintf("%fGi", *input.CapacityGB))
		newSize, err := utils.ValidateStorageQuantity(*newCapacity)
		if err != nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
		}

		existingSize, err := utils.ValidateStorageQuantityGB(pvc.CapacityGB)
		if err != nil {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
		}

		// Ensure new size is greater than existing size
		if newSize.Cmp(existingSize) < 0 {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, "New size must be greater than existing size")
		}
	}

	var updatedPvc *models.PVCInfo
	if err := self.repo.WithTx(ctx, func(tx repository.TxInterface) error {
		err := self.repo.System().UpsertPVCMetadata(ctx, tx, pvc.ID, input.Name, input.Description)
		if err != nil {
			return err
		}

		updatedPvc, err = self.k8s.UpdatePersistentVolumeClaim(ctx,
			team.Namespace,
			pvc.ID,
			newCapacity,
			client,
		)

		if err != nil {
			return err
		}

		// Get latest metadata
		pvcMetadata, err := self.repo.System().GetPVCMetadata(ctx, tx, []string{pvc.ID})
		if err != nil {
			return err
		}

		if metadata, ok := pvcMetadata[pvc.ID]; ok {
			if metadata.Name != nil {
				updatedPvc.Name = *metadata.Name
			} else {
				updatedPvc.Name = pvc.ID
			}
			updatedPvc.Description = metadata.Description
		} else {
			updatedPvc.Name = pvc.ID
		}
		return nil

	}); err != nil {
		return nil, err
	}

	return updatedPvc, nil
}
