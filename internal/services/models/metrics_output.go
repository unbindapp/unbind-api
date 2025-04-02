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

type MetricsEntry struct {
	ID   string          `json:"id"`
	Data MetricsMapEntry `json:"data"`
}

type MetricsResult struct {
	Step    time.Duration  `json:"step"`
	Metrics []MetricsEntry `json:"metrics"`
}

func TransformMetricsEntity(metrics map[string]*prometheus.ResourceMetrics, step time.Duration) *MetricsResult {
	result := &MetricsResult{
		Step:    step,
		Metrics: make([]MetricsEntry, len(metrics)),
	}
	for serviceName, metric := range metrics {
		result.Metrics = append(result.Metrics, MetricsEntry{
			ID: serviceName,
			Data: MetricsMapEntry{
				CPU:     transformMetrics(metric.CPU),
				RAM:     transformMetrics(metric.RAM),
				Disk:    transformMetrics(metric.Disk),
				Network: transformMetrics(metric.Network),
			},
		})
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
