package service_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/internal/repository"
)

// Create the service config
func (self *ServiceRepository) CreateConfig(
	ctx context.Context,
	tx repository.TxInterface,
	serviceID uuid.UUID,
	gitBranch *string,
	port *int,
	host *string,
	replicas *int32,
	autoDeploy *bool,
	runCommand *string,
	public *bool,
	image *string,
) (*ent.ServiceConfig, error) {
	db := self.base.DB
	if tx != nil {
		db = tx.Client()
	}

	return db.ServiceConfig.Create().
		SetServiceID(serviceID).
		SetNillableGitBranch(gitBranch).
		SetNillablePort(port).
		SetNillableHost(host).
		SetNillableReplicas(replicas).
		SetNillableAutoDeploy(autoDeploy).
		SetNillableRunCommand(runCommand).
		SetNillablePublic(public).
		SetNillableImage(image).
		Save(ctx)
}
