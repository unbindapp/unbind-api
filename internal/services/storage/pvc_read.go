package storage_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *StorageService) ListPVCs(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.ListPVCInput) ([]k8s.PVCInfo, error) {
	team, _, _, err := self.validatePermissionsAndParseInputs(ctx, schema.ActionViewer, requesterUserID, input.Type, input.TeamID, input.ProjectID, input.EnvironmentID)
	if err != nil {
		return nil, err
	}

	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Build labels to query
	labels := map[string]string{
		"unbind-team": input.TeamID.String(),
	}
	switch input.Type {
	case models.PvcScopeProject:
		labels["unbind-project"] = input.ProjectID.String()
	case models.PvcScopeEnvironment:
		labels["unbind-environment"] = input.EnvironmentID.String()
	}

	// Get the PVCs
	return self.k8s.ListPersistentVolumeClaims(ctx, team.Namespace, labels, client)
}

func (self *StorageService) GetPVC(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.GetPVCInput) (*k8s.PVCInfo, error) {
	team, _, _, err := self.validatePermissionsAndParseInputs(ctx, schema.ActionViewer, requesterUserID, input.Type, input.TeamID, input.ProjectID, input.EnvironmentID)
	if err != nil {
		return nil, err
	}

	client, err := self.k8s.CreateClientWithToken(bearerToken)
	if err != nil {
		return nil, err
	}

	// Get the PVC
	pvc, err := self.k8s.GetPersistentVolumeClaim(ctx, team.Namespace, input.PVCName, client)
	if err != nil {
		return nil, err
	}

	// Ensure it belongs to the team, project, or environment depending on the scope
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
		if pvc.TeamID != input.TeamID || (pvc.EnvironmentID == nil || *pvc.EnvironmentID != input.EnvironmentID) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, "PVC not found")
		}
	}

	return pvc, nil
}
