package storage_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *StorageService) ListPVCs(ctx context.Context, requesterUserID uuid.UUID, bearerToken string, input *models.ListPVCInput) ([]k8s.PVCInfo, error) {
	team, _, _, err := self.validatePermissionsAndParseInputs(ctx, requesterUserID, input.Type, input.TeamID, input.ProjectID, input.EnvironmentID)
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
