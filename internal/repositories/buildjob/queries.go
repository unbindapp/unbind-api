package buildjob_repo

import (
	"context"

	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/buildjob"
)

func (self *BuildJobRepository) GetJobsByStatus(ctx context.Context, status buildjob.Status) ([]*ent.BuildJob, error) {
	return self.base.DB.BuildJob.Query().
		Where(buildjob.StatusEQ(status)).
		All(ctx)
}
