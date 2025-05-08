package storage_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/services/models"
	v1 "k8s.io/api/core/v1"
)

func (self *StorageService) CreatePVC(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.CreatePVCInput) (*k8s.PVCInfo, error) {
	team, _, _, err := self.validatePermissionsAndParseInputs(ctx, schema.ActionEditor, requesterUserID, input.Type, input.TeamID, input.ProjectID, input.EnvironmentID)
	if err != nil {
		return nil, err
	}

	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Parse size
	_, err = validateStorageQuantity(input.Size)
	if err != nil {
		return nil, errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, err.Error())
	}

	// Build labels to set
	labels := map[string]string{
		"unbind-team": input.TeamID.String(),
	}
	switch input.Type {
	case models.PvcScopeProject:
		labels["unbind-project"] = input.ProjectID.String()
	case models.PvcScopeEnvironment:
		labels["unbind-project"] = input.ProjectID.String()
		labels["unbind-environment"] = input.EnvironmentID.String()
	}

	//  Generate a name
	kubernetesName, err := utils.GenerateSlug(input.Name)
	if err != nil {
		return nil, err
	}

	// Get the PVCs
	return self.k8s.CreatePersistentVolumeClaim(ctx,
		team.Namespace,
		kubernetesName,
		input.Name,
		labels,
		input.Size,
		[]v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
		nil,
		client,
	)
}
