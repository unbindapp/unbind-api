package prometheus

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// Get metrics for specific nodes
func (self *PrometheusClient) GetNodeMetrics(
	ctx context.Context,
	sumBy NodeMetricsFilterSumBy,
	start time.Time,
	end time.Time,
	step time.Duration,
	filter *NodeMetricsFilter,
) (map[string]*NodeMetrics, error) {
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	// Build the label selector for node metrics
	nodeSelector := buildNodeLabelSelector(filter)

	// Queries for node-level metrics
	cpuQuery := fmt.Sprintf(`sum by (%s) (
		rate(node_cpu_seconds_total{mode!="idle"}[%ds])%s
	)`, sumBy.Label(), int(step.Seconds()), nodeSelector)

	ramQuery := fmt.Sprintf(`sum by (%s) (
		(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes)%s
	)`, sumBy.Label(), nodeSelector)

	networkQuery := fmt.Sprintf(`sum by (%s) (
		rate(node_network_receive_bytes_total[%ds]) + 
		rate(node_network_transmit_bytes_total[%ds])%s
	)`, sumBy.Label(), int(step.Seconds()), int(step.Seconds()), nodeSelector)

	diskQuery := fmt.Sprintf(`sum by (%s) (
		rate(node_disk_read_bytes_total[%ds]) + 
		rate(node_disk_written_bytes_total[%ds])%s
	)`, sumBy.Label(), int(step.Seconds()), int(step.Seconds()), nodeSelector)

	fsQuery := fmt.Sprintf(`sum by (%s) (
		node_filesystem_size_bytes%s - node_filesystem_free_bytes%s
	)`, sumBy.Label(), nodeSelector, nodeSelector)

	// Load average
	loadQuery := fmt.Sprintf(`sum by (%s) (
		node_load1%s
	)`, sumBy.Label(), nodeSelector)

	// Execute queries
	cpuResult, _, err := self.api.QueryRange(ctx, cpuQuery, r)
	if err != nil {
		return nil, fmt.Errorf("error querying node CPU metrics: %w", err)
	}

	ramResult, _, err := self.api.QueryRange(ctx, ramQuery, r)
	if err != nil {
		return nil, fmt.Errorf("error querying node RAM metrics: %w", err)
	}

	networkResult, _, err := self.api.QueryRange(ctx, networkQuery, r)
	if err != nil {
		return nil, fmt.Errorf("error querying node network metrics: %w", err)
	}

	diskResult, _, err := self.api.QueryRange(ctx, diskQuery, r)
	if err != nil {
		return nil, fmt.Errorf("error querying node disk I/O metrics: %w", err)
	}

	fsResult, _, err := self.api.QueryRange(ctx, fsQuery, r)
	if err != nil {
		return nil, fmt.Errorf("error querying node filesystem metrics: %w", err)
	}

	loadResult, _, err := self.api.QueryRange(ctx, loadQuery, r)
	if err != nil {
		return nil, fmt.Errorf("error querying node load metrics: %w", err)
	}

	// Process results
	metricsResult := make(map[string]*NodeMetrics)

	extractNodeMetrics(cpuResult, sumBy.Label(), metricsResult, func(metrics *NodeMetrics, samples []model.SamplePair) {
		metrics.CPU = samples
	})

	extractNodeMetrics(ramResult, sumBy.Label(), metricsResult, func(metrics *NodeMetrics, samples []model.SamplePair) {
		metrics.RAM = samples
	})

	extractNodeMetrics(networkResult, sumBy.Label(), metricsResult, func(metrics *NodeMetrics, samples []model.SamplePair) {
		metrics.Network = samples
	})

	extractNodeMetrics(diskResult, sumBy.Label(), metricsResult, func(metrics *NodeMetrics, samples []model.SamplePair) {
		metrics.Disk = samples
	})

	extractNodeMetrics(fsResult, sumBy.Label(), metricsResult, func(metrics *NodeMetrics, samples []model.SamplePair) {
		metrics.FileSystem = samples
	})

	extractNodeMetrics(loadResult, sumBy.Label(), metricsResult, func(metrics *NodeMetrics, samples []model.SamplePair) {
		metrics.Load = samples
	})

	return metricsResult, nil
}

func extractNodeMetrics(
	result model.Value,
	sumByLabel string,
	groupedMetrics map[string]*NodeMetrics,
	assignFunc func(*NodeMetrics, []model.SamplePair),
) {
	if matrix, ok := result.(model.Matrix); ok {
		for _, series := range matrix {
			// Get node identifier from the metric labels
			nodeID := string(series.Metric[model.LabelName(sumByLabel)])
			if nodeID == "" {
				nodeID = "unknown" // Default for metrics without the specified label
			}

			// Create node entry if it doesn't exist
			if _, exists := groupedMetrics[nodeID]; !exists {
				groupedMetrics[nodeID] = &NodeMetrics{}
			}

			// Assign metrics using the provided function
			assignFunc(groupedMetrics[nodeID], series.Values)
		}
	}
}

// buildNodeLabelSelector constructs a Prometheus label selector string for node metrics filtering
func buildNodeLabelSelector(filter *NodeMetricsFilter) string {
	if filter == nil {
		return ""
	}

	var selector string

	// Add node name filter
	if len(filter.Name) > 0 {
		selector += buildLabelValueFilter("node", filter.Name)
	}

	// Add zone filter
	if len(filter.Zone) > 0 {
		selector += buildLabelValueFilter("zone", filter.Zone)
	}

	// Add region filter
	if len(filter.Region) > 0 {
		selector += buildLabelValueFilter("region", filter.Region)
	}

	// Add cluster filter
	if len(filter.Cluster) > 0 {
		selector += buildLabelValueFilter("cluster", filter.Cluster)
	}

	return selector
}

// buildLabelValueFilter creates a label filter expression for Prometheus
func buildLabelValueFilter(label string, values []string) string {
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		return fmt.Sprintf(`, %s="%s"`, label, values[0])
	}

	filter := fmt.Sprintf(`, %s=~"%s"`, label, values[0])
	for i := 1; i < len(values); i++ {
		filter += fmt.Sprintf("|%s", values[i])
	}
	filter += `"`

	return filter
}
