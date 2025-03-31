package metric_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/infrastructure/prometheus"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

func (self *MetricsService) GetMetrics(ctx context.Context, requesterUserID uuid.UUID, input *models.MetricsQueryInput) (*models.MetricsResult, error) {
	team, project, environment, service, err := self.validatePermissionsAndParseInputs(ctx, requesterUserID, input)
	if err != nil {
		return nil, err
	}

	// Build options
	metricsFilters := prometheus.MetricsFilter{}

	// Build labels to select
	switch input.Type {
	case models.MetricsTypeTeam:
		metricsFilters.TeamID = team.ID
	case models.MetricsTypeProject:
		metricsFilters.ProjectID = project.ID
	case models.MetricsTypeEnvironment:
		metricsFilters.EnvironmentID = environment.ID
	case models.MetricsTypeService:
		metricsFilters.ServiceID = service.ID
	}

	// Parse 'step' duration
	var step time.Duration
	if input.Step != "" {
		step, err = time.ParseDuration(input.Step)
		if err != nil {
			return nil, fmt.Errorf("invalid step duration: %w", err)
		}
	}

	// Get start
	var start time.Time
	if input.Start.IsZero() {
		// Default to 7 days ago
		start = time.Now().Add(-7 * 24 * time.Hour)
	} else {
		start = input.Start
	}

	// Get end
	var end time.Time
	if input.End.IsZero() {
		// Default to now
		end = time.Now()
	} else {
		end = input.End
	}

	// Get metrics
	rawMetrics, err := self.promClient.GetResourceMetrics(ctx, start, end, step, &metricsFilters)
	if err != nil {
		return nil, fmt.Errorf("error getting metrics: %w", err)
	}

	// Convert to our format
	return models.TransformMetricsEntity(rawMetrics), nil
}
