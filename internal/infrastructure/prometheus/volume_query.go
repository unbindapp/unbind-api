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
	CapacityGB *float64 `json:"capacity_gb,omitempty"`
}

// VolumeStatsWithHistory combines current PVC stats with historical usage data.
type VolumeStatsWithHistory struct {
	Stats   PVCVolumeStats     `json:"stats"`
	History []model.SamplePair `json:"history"`
}

// GetPVCsVolumeStats queries Prometheus for volume usage and capacity for a list of PVCs.
// If no data is found in Prometheus (unattached volumes), it falls back to the Kubernetes API.
func (self *PrometheusClient) GetPVCsVolumeStats(ctx context.Context, pvcNames []string, namespace string, client *kubernetes.Clientset) ([]PVCVolumeStats, error) {
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
		case "capacity":
			statEntry.CapacityGB = utils.ToPtr(value)
		}
	}

	finalStats := make([]PVCVolumeStats, len(pvcNames))
	for i, name := range pvcNames {
		if data, found := statsMap[name]; found {
			finalStats[i] = *data
		} else {
			// No data, may not be attached
			finalStats[i] = PVCVolumeStats{PVCName: name}

			// Try to get PVC info from K8s API if client is provided
			if client != nil && namespace != "" {
				if pvcInfo, err := self.getPVCFromK8s(ctx, namespace, name, client); err == nil {
					finalStats[i].CapacityGB = pvcInfo
				}
				// ! TODO - PVC may not exist?
			}
		}
	}

	return finalStats, nil
}

// Helper function to get PVC capacity from Kubernetes API
func (self *PrometheusClient) getPVCFromK8s(ctx context.Context, namespace, pvcName string, client *kubernetes.Clientset) (*float64, error) {
	pvc, err := client.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// First try to get capacity from PVC stats
	if pvc.Status.Capacity != nil {
		if capacity, ok := pvc.Status.Capacity[corev1.ResourceStorage]; ok {
			bytesValue := capacity.Value()
			gbValue := float64(bytesValue) / (1024 * 1024 * 1024)
			return &gbValue, nil
		}
	}

	// Fall back to requests if capacity is not set
	if storageRequest, ok := pvc.Spec.Resources.Requests[corev1.ResourceStorage]; ok {
		bytesValue := storageRequest.Value()
		gbValue := float64(bytesValue) / (1024 * 1024 * 1024)
		return &gbValue, nil
	}

	return nil, fmt.Errorf("no storage request found in PVC")
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

	var pvcStats PVCVolumeStats
	if len(stats) > 0 {
		pvcStats = stats[0]
	}

	return &VolumeStatsWithHistory{
		Stats:   pvcStats,
		History: history,
	}, nil
}
