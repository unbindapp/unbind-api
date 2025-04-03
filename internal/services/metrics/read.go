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
	sumBy := prometheus.MetricsFilterSumByService

	// Build labels to select
	switch input.Type {
	case models.MetricsTypeTeam:
		sumBy = prometheus.MetricsFilterSumByProject
		metricsFilters.TeamID = team.ID
	case models.MetricsTypeProject:
		sumBy = prometheus.MetricsFilterSumByEnvironment
		metricsFilters.ProjectID = project.ID
	case models.MetricsTypeEnvironment:
		sumBy = prometheus.MetricsFilterSumByService
		metricsFilters.EnvironmentID = environment.ID
		// Get services in this environment
		services, err := self.repo.Service().GetByEnvironmentID(ctx, environment.ID)
		if err != nil {
			return nil, err
		}
		metricsFilters.ServiceIDs = make([]uuid.UUID, len(services))
		for i, s := range services {
			metricsFilters.ServiceIDs[i] = s.ID
		}
	case models.MetricsTypeService:
		sumBy = prometheus.MetricsFilterSumByService
		metricsFilters.ServiceIDs = []uuid.UUID{service.ID}
	}

	// Get start
	var start time.Time
	if input.Start.IsZero() {
		// Default to 24 hours ago
		start = time.Now().Add(-1 * 24 * time.Hour)
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

	// Calculate step size
	step := 5 * time.Minute // Default
	duration := end.Sub(start)
	// Default step calculations
	if !input.Start.IsZero() || !input.End.IsZero() {
		switch {
		case duration <= 5*time.Minute:
			step = 10 * time.Second
		case duration <= 15*time.Minute:
			step = 30 * time.Second
		case duration <= 30*time.Minute:
			step = 1 * time.Minute
		case duration <= 1*time.Hour:
			step = 2 * time.Minute
		case duration <= 4*time.Hour:
			step = 5 * time.Minute
		case duration <= 12*time.Hour:
			step = 15 * time.Minute
		case duration <= 24*time.Hour:
			step = 30 * time.Minute
		case duration <= 3*24*time.Hour:
			step = 2 * time.Hour
		case duration <= 7*24*time.Hour:
			step = 4 * time.Hour
		case duration <= 14*24*time.Hour:
			step = 8 * time.Hour
		case duration <= 30*24*time.Hour:
			step = 24 * time.Hour
		default:
			step = 24 * time.Hour
		}
	}

	// Get metrics
	rawMetrics, err := self.promClient.GetResourceMetrics(ctx, sumBy, start, end, step, &metricsFilters)
	if err != nil {
		return nil, fmt.Errorf("error getting metrics: %w", err)
	}

	// Convert to our format
	return models.TransformMetricsEntity(rawMetrics, step, sumBy), nil
}
