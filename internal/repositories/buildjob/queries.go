package buildjob_repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/buildjob"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

func (self *BuildJobRepository) GetJobsByStatus(ctx context.Context, status schema.BuildJobStatus) ([]*ent.BuildJob, error) {
	return self.base.DB.BuildJob.Query().
		Where(buildjob.StatusEQ(status)).
		All(ctx)
}

func (self *BuildJobRepository) GetByServiceIDPaginated(ctx context.Context, serviceID uuid.UUID, cursor *time.Time, statusFilter *schema.BuildJobStatus) (jobs []*ent.BuildJob, nextCursor *time.Time, err error) {
	perPage := 10

	query := self.base.DB.BuildJob.Query().
		Where(buildjob.ServiceIDEQ(serviceID))

	if cursor != nil {
		query = query.Where(buildjob.CreatedAtLT(*cursor))
	}

	if statusFilter != nil {
		query = query.Where(buildjob.StatusEQ(*statusFilter))
	}

	all, err := query.
		Order(ent.Desc(buildjob.FieldCreatedAt)).
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
