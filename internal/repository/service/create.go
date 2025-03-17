package service_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/internal/repository"
	"github.com/unbindapp/unbind-api/internal/utils"
)

// Create the service
func (self *ServiceRepository) Create(
	ctx context.Context,
	tx repository.TxInterface,
	displayName string,
	description string,
	serviceType service.Type,
	subtype service.Subtype,
	environmentID uuid.UUID,
	gitHubInstallationID *int64,
	gitRepository *string,
) (*ent.Service, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	// Generate unique name
	name, err := utils.GenerateSlug(displayName)
	if err != nil {
		return nil, err
	}

	return db.Service.Create().
		SetName(name).
		SetDisplayName(displayName).
		SetDescription(description).
		SetType(serviceType).
		SetSubtype(subtype).
		SetEnvironmentID(environmentID).
		SetNillableGithubInstallationID(gitHubInstallationID).
		SetNillableGitRepository(gitRepository).
		Save(ctx)
}
