package metric_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/infrastructure/prometheus"
	permissions_repo "github.com/unbindapp/unbind-api/internal/repositories/permissions"
	"github.com/unbindapp/unbind-api/internal/services/models"
)

// New method for getting node metrics
func (self *MetricsService) GetNodeMetrics(ctx context.Context, requesterUserID uuid.UUID, input *models.NodeMetricsQueryInput) (*models.NodeMetricsResult, error) {
	permissionChecks := []permissions_repo.PermissionCheck{
		{
			Action:       schema.ActionViewer,
			ResourceType: schema.ResourceTypeSystem,
		},
	}

	// Check permissions
	if err := self.repo.Permissions().Check(ctx, requesterUserID, permissionChecks); err != nil {
		return nil, err
	}

	// Build options
	nodeMetricsFilters := prometheus.NodeMetricsFilter{}
	var sumBy prometheus.NodeMetricsFilterSumBy

	// Set the appropriate grouping based on the query type
	switch input.Type {
	case models.NodeMetricsTypeNode:
		sumBy = prometheus.NodeSumByName
		if input.NodeName != "" {
			nodeMetricsFilters.Name = []string{input.NodeName}
		}
	case models.NodeMetricsTypeZone:
		sumBy = prometheus.NodeSumByZone
		if input.Zone != "" {
			nodeMetricsFilters.Zone = []string{input.Zone}
		}
	case models.NodeMetricsTypeRegion:
		sumBy = prometheus.NodeSumByRegion
		if input.Region != "" {
			nodeMetricsFilters.Region = []string{input.Region}
		}
	case models.NodeMetricsTypeCluster:
		sumBy = prometheus.NodeSumByCluster
		if input.ClusterName != "" {
			nodeMetricsFilters.Cluster = []string{input.ClusterName}
		}
	default:
		return nil, fmt.Errorf("invalid node metrics type: %s", input.Type)
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
	duration := end.Sub(start)
	step := chooseStep(duration, 30, []time.Duration{
		1 * time.Minute,
		5 * time.Minute,
		15 * time.Minute,
		30 * time.Minute,
		1 * time.Hour,
		2 * time.Hour,
		4 * time.Hour,
		8 * time.Hour,
		12 * time.Hour,
		1 * 24 * time.Hour,
	})

	// Get metrics
	var filter *prometheus.NodeMetricsFilter
	if input.NodeName != "" || input.Zone != "" || input.Region != "" || input.ClusterName != "" {
		filter = &nodeMetricsFilters
	} else {
		filter = nil // Get all nodes if no specific filters
	}

	rawMetrics, err := self.promClient.GetNodeMetrics(ctx, sumBy, start, end, step, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting node metrics: %w", err)
	}

	// Convert to our format
	return models.TransformNodeMetricsEntity(rawMetrics, step, sumBy), nil
}
