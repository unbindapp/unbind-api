package prometheus

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type ResourceMetrics struct {
	CPU     []model.SamplePair
	RAM     []model.SamplePair
	Network []model.SamplePair
	Disk    []model.SamplePair
}

// Get metrics for specific resources (team, project, environment, service)
func (self *PrometheusClient) GetResourceMetrics(
	ctx context.Context,
	start time.Time,
	end time.Time,
	step time.Duration,
	filter *MetricsFilter,
) (*ResourceMetrics, error) {
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	labelSelector := buildLabelSelector(filter)

	// Queries for different resource metrics
	cpuQuery := "sum(rate(container_cpu_usage_seconds_total" + labelSelector + "[5m])) by (pod)"
	ramQuery := "sum(container_memory_usage_bytes" + labelSelector + ") by (pod)"
	networkQuery := "sum(rate(container_network_receive_bytes_total" + labelSelector + "[5m]) + rate(container_network_transmit_bytes_total" + labelSelector + "[5m])) by (pod)"
	diskQuery := "sum(container_fs_usage_bytes" + labelSelector + ") by (pod)"

	// Execute queries
	cpuResult, _, err := self.api.QueryRange(ctx, cpuQuery, r)
	if err != nil {
		return nil, fmt.Errorf("error querying CPU metrics: %w", err)
	}

	ramResult, _, err := self.api.QueryRange(ctx, ramQuery, r)
	if err != nil {
		return nil, fmt.Errorf("error querying RAM metrics: %w", err)
	}

	networkResult, _, err := self.api.QueryRange(ctx, networkQuery, r)
	if err != nil {
		return nil, fmt.Errorf("error querying network metrics: %w", err)
	}

	diskResult, _, err := self.api.QueryRange(ctx, diskQuery, r)
	if err != nil {
		return nil, fmt.Errorf("error querying disk metrics: %w", err)
	}

	// Process results
	metrics := &ResourceMetrics{}

	// Extract CPU metrics
	if cpuMatrix, ok := cpuResult.(model.Matrix); ok && len(cpuMatrix) > 0 {
		for _, s := range cpuMatrix[0].Values {
			metrics.CPU = append(metrics.CPU, s)
		}
	}

	// Extract RAM metrics
	if ramMatrix, ok := ramResult.(model.Matrix); ok && len(ramMatrix) > 0 {
		for _, s := range ramMatrix[0].Values {
			metrics.RAM = append(metrics.RAM, s)
		}
	}

	// Extract Network metrics
	if networkMatrix, ok := networkResult.(model.Matrix); ok && len(networkMatrix) > 0 {
		for _, s := range networkMatrix[0].Values {
			metrics.Network = append(metrics.Network, s)
		}
	}

	// Extract Disk metrics
	if diskMatrix, ok := diskResult.(model.Matrix); ok && len(diskMatrix) > 0 {
		for _, s := range diskMatrix[0].Values {
			metrics.Disk = append(metrics.Disk, s)
		}
	}

	return metrics, nil
}
