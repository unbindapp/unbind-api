package deployment_repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/deployment"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

func (self *DeploymentRepository) GetJobsByStatus(ctx context.Context, status schema.DeploymentStatus) ([]*ent.Deployment, error) {
	return self.base.DB.Deployment.Query().
		Where(deployment.StatusEQ(status)).
		All(ctx)
}

func (self *DeploymentRepository) GetByServiceIDPaginated(ctx context.Context, serviceID uuid.UUID, cursor *time.Time, statusFilter *schema.DeploymentStatus) (jobs []*ent.Deployment, nextCursor *time.Time, err error) {
	perPage := 10

	query := self.base.DB.Deployment.Query().
		Where(deployment.ServiceIDEQ(serviceID))

	if cursor != nil {
		query = query.Where(deployment.CreatedAtLT(*cursor))
	}

	if statusFilter != nil {
		query = query.Where(deployment.StatusEQ(*statusFilter))
	}

	all, err := query.
		Order(ent.Desc(deployment.FieldCreatedAt)).
		Limit(perPage + 1).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	// If we have more than the perPage limit, we have a next page, get its cursot and truncate the results
	if len(all) > perPage {
		nextCursor = utils.ToPtr(all[perPage].CreatedAt)
		all = all[:perPage]
	}

	return all, nextCursor, nil
}
