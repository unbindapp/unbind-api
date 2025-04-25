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
	sumBy MetricsFilterSumBy,
	start time.Time,
	end time.Time,
	step time.Duration,
	filter *MetricsFilter,
) (map[string]*ResourceMetrics, error) {
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	// Build the label selector for kube_pod_labels
	kubeLabelsSelector := buildLabelSelector(filter)

	// Queries with label filtering
	cpuQuery := fmt.Sprintf(`sum by (%s) (
			rate(container_cpu_usage_seconds_total{container!="POD", container!=""}[%ds])
			* on(namespace, pod) group_left(label_unbind_team,label_unbind_project,label_unbind_environment,label_unbind_service)
			kube_pod_labels%s
	)`, sumBy.Label(), int(step.Seconds()), kubeLabelsSelector)

	ramQuery := fmt.Sprintf(`sum by (%s) (
		container_memory_working_set_bytes{container!="POD", container!=""}
		* on(namespace, pod) group_left(label_unbind_team,label_unbind_project,label_unbind_environment,label_unbind_service)
		kube_pod_labels%s
	)`, sumBy.Label(), kubeLabelsSelector)

	networkQuery := fmt.Sprintf(`sum by (%s) (
		(increase(container_network_receive_bytes_total{pod!=""}[%ds]) +
		increase(container_network_transmit_bytes_total{pod!=""}[%ds]))
		* on(namespace, pod) group_left(label_unbind_team,label_unbind_project,label_unbind_environment,label_unbind_service)
		kube_pod_labels%s
	)`, sumBy.Label(), int(step.Seconds()), int(step.Seconds()), kubeLabelsSelector)

	diskQuery := fmt.Sprintf(`sum by (%s) (
		kubelet_volume_stats_used_bytes
			* on(namespace,persistentvolumeclaim) group_left(pod)
				kube_pod_spec_volumes_persistentvolumeclaims_info
			* on(namespace,pod) group_left(label_unbind_service)
				kube_pod_labels%s
	)`, sumBy.Label(), kubeLabelsSelector)

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
	metricsResult := make(map[string]*ResourceMetrics)

	extractGroupedMetrics(cpuResult, sumBy.Label(), metricsResult, func(metrics *ResourceMetrics, samples []model.SamplePair) {
		metrics.CPU = samples
	})

	extractGroupedMetrics(ramResult, sumBy.Label(), metricsResult, func(metrics *ResourceMetrics, samples []model.SamplePair) {
		metrics.RAM = samples
	})

	extractGroupedMetrics(networkResult, sumBy.Label(), metricsResult, func(metrics *ResourceMetrics, samples []model.SamplePair) {
		metrics.Network = samples
	})

	extractGroupedMetrics(diskResult, sumBy.Label(), metricsResult, func(metrics *ResourceMetrics, samples []model.SamplePair) {
		metrics.Disk = samples
	})

	return metricsResult, nil
}

func extractGroupedMetrics(
	result model.Value,
	sumByLabel string,
	groupedMetrics map[string]*ResourceMetrics,
	assignFunc func(*ResourceMetrics, []model.SamplePair),
) {
	if matrix, ok := result.(model.Matrix); ok {
		for _, series := range matrix {
			// Get service name from the metric labels
			serviceName := string(series.Metric[model.LabelName(sumByLabel)])
			if serviceName == "" {
				serviceName = "unknown" // Default for metrics without service label
			}

			// Create service entry if it doesn't exist
			if _, exists := groupedMetrics[serviceName]; !exists {
				groupedMetrics[serviceName] = &ResourceMetrics{}
			}

			// Assign metrics using the provided function
			assignFunc(groupedMetrics[serviceName], series.Values)
		}
	}
}
