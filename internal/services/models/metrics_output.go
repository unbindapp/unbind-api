package models

import (
	"sort"
	"time"

	"github.com/prometheus/common/model"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/infrastructure/prometheus"
)

// MetricDetail represents a single metric with aggregate value and breakdown by service
type MetricDetail struct {
	Timestamp time.Time `json:"timestamp"`
	// Aggregated value for the metric
	Value float64 `json:"value" doc:"Aggregated value for the timestamp"`
	// Map of IDs to their respective values
	Breakdown map[string]*float64 `json:"breakdown" doc:"Map of IDs to their respective values"`
}

// MetricsMapEntry contains arrays of metric details for each resource type
type MetricsMapEntry struct {
	CPU     []MetricDetail `json:"cpu" nullable:"false"`
	RAM     []MetricDetail `json:"ram" nullable:"false"`
	Disk    []MetricDetail `json:"disk" nullable:"false"`
	Network []MetricDetail `json:"network" nullable:"false"`
}

// MetricsResult is the top-level structure containing the sampling interval and metrics
type MetricsResult struct {
	Step         time.Duration   `json:"step"`
	BrokenDownBy MetricsType     `json:"broken_down_by" doc:"The type of metric that is broken down, e.g. team, project"`
	Metrics      MetricsMapEntry `json:"metrics" nullable:"false"`
}

func TransformMetricsEntity(metrics map[string]*prometheus.ResourceMetrics, step time.Duration, sumBy prometheus.MetricsFilterSumBy) *MetricsResult {
	brokenDownBy := MetricsTypeService
	switch sumBy {
	case prometheus.MetricsFilterSumByProject:
		brokenDownBy = MetricsTypeProject
	case prometheus.MetricsFilterSumByEnvironment:
		brokenDownBy = MetricsTypeEnvironment
	case prometheus.MetricsFilterSumByService:
		brokenDownBy = MetricsTypeService
	}
	result := &MetricsResult{
		Step:         step,
		BrokenDownBy: brokenDownBy,
		Metrics: MetricsMapEntry{
			CPU:     aggregateMetricsByTime(metrics, MetricTypeCPU),
			RAM:     aggregateMetricsByTime(metrics, MetricTypeRAM),
			Disk:    aggregateMetricsByTime(metrics, MetricTypeDisk),
			Network: aggregateMetricsByTime(metrics, MetricTypeNetwork),
		},
	}

	return result
}

type MetricType int

const (
	MetricTypeCPU MetricType = iota
	MetricTypeRAM
	MetricTypeDisk
	MetricTypeNetwork
)

func aggregateMetricsByTime(metrics map[string]*prometheus.ResourceMetrics, metricType MetricType) []MetricDetail {
	// Gather all possible timestamps
	timestampMetricMap := make(map[time.Time]MetricDetail)

	// This will make our array allocations more efficient later
	for _, metric := range metrics {
		switch metricType {
		case MetricTypeCPU:
			for _, sample := range metric.CPU {
				timestampMetricMap[sample.Timestamp.Time()] = MetricDetail{}
			}
		case MetricTypeRAM:
			for _, sample := range metric.RAM {
				timestampMetricMap[sample.Timestamp.Time()] = MetricDetail{}
			}
		case MetricTypeDisk:
			for _, sample := range metric.Disk {
				timestampMetricMap[sample.Timestamp.Time()] = MetricDetail{}
			}
		case MetricTypeNetwork:
			for _, sample := range metric.Network {
				timestampMetricMap[sample.Timestamp.Time()] = MetricDetail{}
			}
		}
	}

	// Figure out all the unique IDs available
	// This is used to build the breakdown map with nil values if necessary
	keys := make([]string, len(metrics))
	keyIndex := 0
	for k := range metrics {
		keys[keyIndex] = k
		keyIndex++
	}

	for id, samples := range metrics {
		var samplePair []model.SamplePair
		switch metricType {
		case MetricTypeCPU:
			samplePair = samples.CPU
		case MetricTypeRAM:
			samplePair = samples.RAM
		case MetricTypeDisk:
			samplePair = samples.Disk
		case MetricTypeNetwork:
			samplePair = samples.Network
		}

		for _, sample := range samplePair {
			metricDetail := timestampMetricMap[sample.Timestamp.Time()]
			if metricDetail.Breakdown == nil {
				metricDetail.Breakdown = make(map[string]*float64)
				// Initialize breakdown map with nil values for all keys
				for _, key := range keys {
					metricDetail.Breakdown[key] = nil
				}
			}
			metricDetail.Breakdown[id] = utils.ToPtr(float64(sample.Value))
			metricDetail.Timestamp = sample.Timestamp.Time()
			metricDetail.Value += float64(sample.Value)
			timestampMetricMap[sample.Timestamp.Time()] = metricDetail
		}
	}

	// Convert map to slice
	metricDetails := make([]MetricDetail, len(timestampMetricMap))
	i := 0
	for _, metricDetail := range timestampMetricMap {
		metricDetails[i] = metricDetail
		i++
	}
	// Sort by timestamp
	sort.Slice(metricDetails, func(i, j int) bool {
		return metricDetails[i].Timestamp.Before(metricDetails[j].Timestamp)
	})
	return metricDetails
}
