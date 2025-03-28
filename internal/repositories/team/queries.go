package team_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/team"
	"github.com/unbindapp/unbind-api/ent/user"
)

func (self *TeamRepository) GetAll(ctx context.Context) ([]*ent.Team, error) {
	return self.base.DB.Team.Query().
		All(ctx)
}

func (self *TeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Team, error) {
	return self.base.DB.Team.Query().
		Where(team.ID(id)).
		WithProjects(func(pq *ent.ProjectQuery) {
			pq.WithEnvironments(
				func(eq *ent.EnvironmentQuery) {
					eq.WithServices(
						func(sq *ent.ServiceQuery) {
							sq.WithGithubInstallation(
								func(giq *ent.GithubInstallationQuery) {
									giq.WithGithubApp()
								},
							)
							sq.WithServiceConfig()
						},
					)
				},
			)
		}).
		Only(ctx)
}

func (self *TeamRepository) HasUserWithID(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (bool, error) {
	return self.base.DB.User.Query().
		Where(user.ID(userID)).
		QueryTeams().
		Where(team.ID(teamID)).
		Exist(ctx)
}
