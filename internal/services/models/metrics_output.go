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

type MetricsResult struct {
	CPU     []MetricsPair `json:"cpu"`
	RAM     []MetricsPair `json:"ram"`
	Disk    []MetricsPair `json:"disk"`
	Network []MetricsPair `json:"network"`
}

func TransformMetricsEntity(metrics *prometheus.ResourceMetrics) *MetricsResult {
	return &MetricsResult{
		CPU:     transformMetrics(metrics.CPU),
		RAM:     transformMetrics(metrics.RAM),
		Disk:    transformMetrics(metrics.Disk),
		Network: transformMetrics(metrics.Network),
	}
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
