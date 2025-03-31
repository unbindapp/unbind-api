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

	// Calculate step
	duration := end.Sub(start)
	var step time.Duration
	switch {
	case duration <= 24*time.Hour:
		// 5 minute step
		step = 5 * time.Minute
	case duration <= 3*24*time.Hour:
		// 30 minute step
		step = 30 * time.Minute
	case duration <= 7*24*time.Hour:
		// 1 hour step
		step = 1 * time.Hour
	case duration <= 30*24*time.Hour:
		// 6 hour step
		step = 6 * time.Hour
	default:
		// 1 day step
		step = 24 * time.Hour
	}

	// Get metrics
	rawMetrics, err := self.promClient.GetResourceMetrics(ctx, start, end, step, &metricsFilters)
	if err != nil {
		return nil, fmt.Errorf("error getting metrics: %w", err)
	}

	// Convert to our format
	return models.TransformMetricsEntity(rawMetrics), nil
}
