package prometheus

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// PVCVolumeStats holds the used and capacity stats for a PVC.
type PVCVolumeStats struct {
	PVCName    string   `json:"pvc_name"`
	UsedGB     *float64 `json:"used_gb,omitempty"`
	CapacityGB *float64 `json:"capacity_gb,omitempty"`
}

// VolumeStatsWithHistory combines current PVC stats with historical usage data.
type VolumeStatsWithHistory struct {
	Stats   PVCVolumeStats     `json:"stats"`
	History []model.SamplePair `json:"history"`
}

// GetPVCsVolumeStats queries Prometheus for volume usage and capacity for a list of PVCs.
func (self *PrometheusClient) GetPVCsVolumeStats(ctx context.Context, pvcNames []string) ([]PVCVolumeStats, error) {
	if len(pvcNames) == 0 {
		return []PVCVolumeStats{}, nil
	}

	pvcRegex := strings.Join(pvcNames, "|")

	// Note: The query uses %s four times for pvcRegex.
	query := fmt.Sprintf(`
(
  label_replace(
    max by (persistentvolumeclaim) (
      kubelet_volume_stats_used_bytes{persistentvolumeclaim=~"%s", job="kubelet"} / 1024 / 1024 / 1024
    ),
    "kind", "used", "persistentvolumeclaim", ".*"
  )
)
or
(
  label_replace(
    max by (persistentvolumeclaim) (
      (kubelet_volume_stats_capacity_bytes{persistentvolumeclaim=~"%s", job="kubelet"} / 1024 / 1024 / 1024)
      or
      (kube_persistentvolumeclaim_resource_requests_storage_bytes{persistentvolumeclaim=~"%s", job="kubelet"} / 1024 / 1024 / 1024)
    ), 
    "kind", "capacity", "persistentvolumeclaim", ".*"
  )
)
`, pvcRegex, pvcRegex, pvcRegex)

	result, warnings, err := self.api.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("Prometheus query failed for PVC stats: %w", err)
	}
	if warnings != nil && len(warnings) > 0 {
		// TODO: Decide how to handle Prometheus warnings. For now, they are ignored.
		// fmt.Printf("Prometheus query warnings for PVC stats: %v\n", warnings) // Example of logging if needed
	}

	vectorData, ok := result.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from Prometheus for PVC stats: expected model.Vector, got %T", result)
	}

	statsMap := make(map[string]*PVCVolumeStats) // pvcName -> stats

	for _, sample := range vectorData {
		metric := sample.Metric
		pvcNameFromMetric := string(metric["persistentvolumeclaim"])
		kind := string(metric["kind"])
		value := float64(sample.Value)

		// Ensure this pvcName is one we requested and initialize if not present
		isTargetPVC := false
		for _, requestedName := range pvcNames {
			if pvcNameFromMetric == requestedName {
				isTargetPVC = true
				break
			}
		}

		if !isTargetPVC {
			// This metric is for a PVC not in our original list, skip.
			// This might happen if the regex matches more broadly than expected, though unlikely with simple PVC names.
			continue
		}

		if _, exists := statsMap[pvcNameFromMetric]; !exists {
			statsMap[pvcNameFromMetric] = &PVCVolumeStats{PVCName: pvcNameFromMetric}
		}

		statEntry := statsMap[pvcNameFromMetric] // Guaranteed to exist now

		switch kind {
		case "used":
			v := value // Create a new variable to take its address
			statEntry.UsedGB = &v
		case "capacity":
			v := value // Create a new variable to take its address
			statEntry.CapacityGB = &v
		}
	}

	// Construct the final result slice, ensuring an entry for every requested PVC name,
	// in the order they were requested.
	finalStats := make([]PVCVolumeStats, len(pvcNames))
	for i, name := range pvcNames {
		if data, found := statsMap[name]; found {
			finalStats[i] = *data
		} else {
			// This PVC was in the input list, but no data was found for it in Prometheus results.
			// Its UsedGB and CapacityGB will remain nil by default.
			finalStats[i] = PVCVolumeStats{PVCName: name}
		}
	}

	return finalStats, nil
}

// GetVolumeStatsWithHistory gets both current PVC stats and historical usage data for a specific volume.
// This is equivalent to the diskQuery but targeted at a specific volume by name.
func (self *PrometheusClient) GetVolumeStatsWithHistory(
	ctx context.Context,
	pvcName string,
	start time.Time,
	end time.Time,
	step time.Duration,
) (*VolumeStatsWithHistory, error) {
	// Get current stats using existing method
	stats, err := self.GetPVCsVolumeStats(ctx, []string{pvcName})
	if err != nil {
		return nil, fmt.Errorf("failed to get PVC stats: %w", err)
	}

	var pvcStats PVCVolumeStats
	if len(stats) > 0 {
		pvcStats = stats[0]
	} else {
		// No stats found, but still create entry with the name
		pvcStats = PVCVolumeStats{PVCName: pvcName}
	}

	// Get historical data using a query similar to diskQuery but for specific PVC
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	// Historical usage query for the specific PVC - similar to diskQuery logic
	historyQuery := fmt.Sprintf(`
		max by (persistentvolumeclaim) (
			kubelet_volume_stats_used_bytes{persistentvolumeclaim="%s"}
		)
	`, pvcName)

	result, _, err := self.api.QueryRange(ctx, historyQuery, r)
	if err != nil {
		return nil, fmt.Errorf("failed to query volume history for PVC %s: %w", pvcName, err)
	}

	history := []model.SamplePair{}
	if matrix, ok := result.(model.Matrix); ok && len(matrix) > 0 {
		// Since we filter by specific PVC name, there should be at most one series
		history = matrix[0].Values
	}

	return &VolumeStatsWithHistory{
		Stats:   pvcStats,
		History: history,
	}, nil
}
