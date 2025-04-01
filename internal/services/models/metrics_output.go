package models

import (
	"time"

	"github.com/prometheus/common/model"
	"github.com/unbindapp/unbind-api/internal/infrastructure/prometheus"
)

type MetricsPair struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

type MetricsMapEntry struct {
	CPU     []MetricsPair `json:"cpu"`
	RAM     []MetricsPair `json:"ram"`
	Disk    []MetricsPair `json:"disk"`
	Network []MetricsPair `json:"network"`
}

type MetricsResult struct {
	Step     time.Duration              `json:"step"`
	Services map[string]MetricsMapEntry `json:"services"`
}

func TransformMetricsEntity(metrics map[string]*prometheus.ResourceMetrics, step time.Duration) *MetricsResult {
	result := &MetricsResult{
		Step:     step,
		Services: make(map[string]MetricsMapEntry),
	}
	for serviceName, metric := range metrics {
		result.Services[serviceName] = MetricsMapEntry{
			CPU:     transformMetrics(metric.CPU),
			RAM:     transformMetrics(metric.RAM),
			Disk:    transformMetrics(metric.Disk),
			Network: transformMetrics(metric.Network),
		}
	}
	return result
}

func transformMetrics(metrics []model.SamplePair) []MetricsPair {
	metricsPairs := make([]MetricsPair, len(metrics))
	for i, metric := range metrics {
		metricsPairs[i] = MetricsPair{
			Time:  metric.Timestamp.Time(),
			Value: float64(metric.Value),
		}
	}
	return metricsPairs
}
