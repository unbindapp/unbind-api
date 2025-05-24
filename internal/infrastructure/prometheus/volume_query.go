package prometheus

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PVCVolumeStats holds the used and capacity stats for a PVC.
type PVCVolumeStats struct {
	PVCName    string   `json:"pvc_name"`
	UsedGB     *float64 `json:"used_gb,omitempty"`
	CapacityGB float64  `json:"capacity_gb,omitempty"`
}

// VolumeStatsWithHistory combines current PVC stats with historical usage data.
type VolumeStatsWithHistory struct {
	Stats   *PVCVolumeStats    `json:"stats"`
	History []model.SamplePair `json:"history"`
}

// GetPVCsVolumeStats queries Prometheus for volume usage and capacity for a list of PVCs.
// If no data is found in Prometheus (unattached volumes), it falls back to the Kubernetes API.
func (self *PrometheusClient) GetPVCsVolumeStats(ctx context.Context, pvcNames []string, namespace string, client *kubernetes.Clientset) ([]*PVCVolumeStats, error) {
	if len(pvcNames) == 0 {
		return []*PVCVolumeStats{}, nil
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
`, pvcRegex)

	result, _, err := self.api.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("Prometheus query failed for PVC stats: %w", err)
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
			continue
		}

		if _, exists := statsMap[pvcNameFromMetric]; !exists {
			statsMap[pvcNameFromMetric] = &PVCVolumeStats{PVCName: pvcNameFromMetric}
		}

		statEntry := statsMap[pvcNameFromMetric] // Guaranteed to exist now

		switch kind {
		case "used":
			statEntry.UsedGB = utils.ToPtr(value)
		}
	}

	finalStats := make([]*PVCVolumeStats, len(pvcNames))
	// Batch process missing PVCs if client is provided
	if pvcCapacities, err := self.getPVCsFromK8s(ctx, namespace, pvcNames, client); err == nil {
		for i, name := range pvcNames {
			if capacity, found := pvcCapacities[name]; found && capacity != nil {
				finalStats[i], ok = statsMap[name]
				if ok {
					finalStats[i].CapacityGB = *capacity
				}
			}
		}
	}

	return finalStats, nil
}

// Helper function to get PVC capacities from Kubernetes API for multiple PVCs
func (self *PrometheusClient) getPVCsFromK8s(ctx context.Context, namespace string, pvcNames []string, client *kubernetes.Clientset) (map[string]*float64, error) {
	results := make(map[string]*float64)

	// Create a set of requested PVC names for efficient lookup
	requestedPVCs := make(map[string]bool)
	for _, name := range pvcNames {
		requestedPVCs[name] = true
		results[name] = nil // Initialize with nil
	}

	// Single API call to list all PVCs in the namespace
	pvcList, err := client.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return results, err
	}

	// Filter and process only the requested PVCs
	for _, pvc := range pvcList.Items {
		if !requestedPVCs[pvc.Name] {
			continue // Skip PVCs we didn't request
		}

		var capacity *float64

		// First try to get capacity from PVC status
		if pvc.Status.Capacity != nil {
			if capacityQuantity, ok := pvc.Status.Capacity[corev1.ResourceStorage]; ok {
				bytesValue := capacityQuantity.Value()
				gbValue := float64(bytesValue) / (1024 * 1024 * 1024)
				capacity = &gbValue
			}
		}

		// Fall back to requests if capacity is not set
		if capacity == nil {
			if storageRequest, ok := pvc.Spec.Resources.Requests[corev1.ResourceStorage]; ok {
				bytesValue := storageRequest.Value()
				gbValue := float64(bytesValue) / (1024 * 1024 * 1024)
				capacity = &gbValue
			}
		}

		results[pvc.Name] = capacity
	}

	return results, nil
}

// GetVolumeStatsWithHistory gets both current PVC stats and historical usage data for a specific volume.
// This is equivalent to the diskQuery but targeted at a specific volume by name.
func (self *PrometheusClient) GetVolumeStatsWithHistory(
	ctx context.Context,
	pvcName string,
	start time.Time,
	end time.Time,
	step time.Duration,
	namespace string,
	client *kubernetes.Clientset,
) (*VolumeStatsWithHistory, error) {
	// First try to get current stats using existing method (works for attached volumes)
	stats, err := self.GetPVCsVolumeStats(ctx, []string{pvcName}, namespace, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get PVC stats: %w", err)
	}

	// Get historical data using a query similar to diskQuery but for specific PVC
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	// Historical usage query for the specific PVC - only works for attached volumes
	historyQuery := fmt.Sprintf(`
		max by (persistentvolumeclaim) (
			kubelet_volume_stats_used_bytes{persistentvolumeclaim="%s"}
		)
	`, pvcName)

	result, _, err := self.api.QueryRange(ctx, historyQuery, r)
	if err != nil {
		return nil, fmt.Errorf("failed to query volume history for PVC %s: %w", pvcName, err)
	}

	var history []model.SamplePair
	if matrix, ok := result.(model.Matrix); ok && len(matrix) > 0 {
		// Since we filter by specific PVC name, there should be at most one series
		history = matrix[0].Values
	}

	var pvcStats *PVCVolumeStats
	if len(stats) > 0 {
		pvcStats = stats[0]
	}

	return &VolumeStatsWithHistory{
		Stats:   pvcStats,
		History: history,
	}, nil
}
