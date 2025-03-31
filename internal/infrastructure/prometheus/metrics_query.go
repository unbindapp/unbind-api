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

	// Build the label selector for kube_pod_labels
	kubeLabelsSelector := buildLabelSelector(filter)

	// Queries with label filtering
	cpuQuery := fmt.Sprintf(`sum(
		node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{container!="POD", container!=""}
		* on(namespace, pod) group_left(label_unbind_team,label_unbind_project,label_unbind_environment,label_unbind_service)
		kube_pod_labels%s
	)`, kubeLabelsSelector)

	ramQuery := fmt.Sprintf(`sum(
		container_memory_working_set_bytes{container!="POD", container!=""}
		* on(namespace, pod) group_left(label_unbind_team,label_unbind_project,label_unbind_environment,label_unbind_service)
		kube_pod_labels%s
	)`, kubeLabelsSelector)

	networkQuery := fmt.Sprintf(`sum(
		(increase(container_network_receive_bytes_total{pod!=""}[%ds]) +
		increase(container_network_transmit_bytes_total{pod!=""}[%ds]))
		* on(namespace, pod) group_left(label_unbind_team,label_unbind_project,label_unbind_environment,label_unbind_service)
		kube_pod_labels%s
	)`, int(step.Seconds()), int(step.Seconds()), kubeLabelsSelector)

	diskQuery := fmt.Sprintf(`sum(
		container_fs_usage_bytes{container!="POD", container!=""}
		* on(namespace, pod) group_left(label_unbind_team,label_unbind_project,label_unbind_environment,label_unbind_service)
		kube_pod_labels%s
	)`, kubeLabelsSelector)

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

	// Extract metrics from all series, not just the first one
	metrics.CPU = extractMetrics(cpuResult)
	metrics.RAM = extractMetrics(ramResult)
	metrics.Network = extractMetrics(networkResult)
	metrics.Disk = extractMetrics(diskResult)

	return metrics, nil
}

func extractMetrics(result model.Value) []model.SamplePair {
	var samples []model.SamplePair

	// Handle matrix results (for range queries)
	if matrix, ok := result.(model.Matrix); ok {
		for _, series := range matrix {
			samples = append(samples, series.Values...)
		}
	}

	return samples
}
