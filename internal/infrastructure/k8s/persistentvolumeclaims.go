package k8s

import (
	"context"
	"fmt"

	"github.com/unbindapp/unbind-api/internal/common/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PVCInfo holds prettier information about a PVC.
type PVCInfo struct {
	Name              string   `json:"name"`
	Size              string   `json:"size"` // e.g., "10Gi"
	TeamID            string   `json:"teamId"`
	ProjectID         *string  `json:"projectId,omitempty"`
	EnvironmentID     *string  `json:"environmentId,omitempty"`
	BoundToServiceIDs []string `json:"boundToServiceIds,omitempty"`
	Status            string   `json:"status"` // e.g., "Bound", "Pending"
}

// CreatePersistentVolumeClaim creates a new PersistentVolumeClaim in the specified namespace.
func (k *KubeClient) CreatePersistentVolumeClaim(
	ctx context.Context,
	namespace string,
	pvcName string,
	labels map[string]string,
	storageRequest string,
	accessModes []corev1.PersistentVolumeAccessMode,
	storageClassName *string,
) (*corev1.PersistentVolumeClaim, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}
	if pvcName == "" {
		return nil, fmt.Errorf("pvcName cannot be empty")
	}
	if storageRequest == "" {
		return nil, fmt.Errorf("storageRequest cannot be empty")
	}
	if len(accessModes) == 0 {
		return nil, fmt.Errorf("at least one accessMode must be provided")
	}

	storageQuantity, err := resource.ParseQuantity(storageRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to parse storageRequest '%s': %w", storageRequest, err)
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: accessModes,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: storageQuantity,
				},
			},
		},
	}

	if storageClassName != nil && *storageClassName != "" {
		pvc.Spec.StorageClassName = storageClassName
	}

	createdPvc, err := k.clientset.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, pvc, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create PersistentVolumeClaim '%s' in namespace '%s': %w", pvcName, namespace, err)
	}

	return createdPvc, nil
}

// GetPersistentVolumeClaim retrieves a specific PersistentVolumeClaim by its name and namespace.
func (k *KubeClient) GetPersistentVolumeClaim(ctx context.Context, namespace string, pvcName string) (*corev1.PersistentVolumeClaim, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}
	if pvcName == "" {
		return nil, fmt.Errorf("pvcName cannot be empty")
	}

	pvc, err := k.clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get PersistentVolumeClaim '%s' in namespace '%s': %w", pvcName, namespace, err)
	}
	return pvc, nil
}

// ListPersistentVolumeClaims lists all PersistentVolumeClaims in a given namespace, optionally filtered by a label selector,
func (k *KubeClient) ListPersistentVolumeClaims(ctx context.Context, namespace string, labelSelector string) ([]PVCInfo, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}

	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}

	pvcList, err := k.clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list PersistentVolumeClaims in namespace '%s' with selector '%s': %w", namespace, labelSelector, err)
	}

	var result []PVCInfo
	const ( // Define label keys for consistency
		teamLabel        = "unbind-team"
		projectLabel     = "unbind-project"
		environmentLabel = "unbind-environment"
		serviceLabel     = "unbind-service"
	)

	for _, pvc := range pvcList.Items {
		pvcLabels := pvc.GetLabels()
		teamID := pvcLabels[teamLabel]

		// Rule 1: Omit if no unbind-team label
		if teamID == "" {
			continue
		}

		projectID := pvcLabels[projectLabel]
		environmentID := pvcLabels[environmentLabel]
		size := ""
		if storageRequest, ok := pvc.Spec.Resources.Requests[corev1.ResourceStorage]; ok {
			size = storageRequest.String()
		}

		var boundToServiceIDs []string
		includePVC := true

		// Rule 2: If bound to pods, they must have unbind-service label
		if pvc.Status.Phase == corev1.ClaimBound {
			pods, err := k.GetPodsUsingPVC(ctx, pvc.Namespace, pvc.Name)
			if err != nil {
				// Log or handle error for this specific PVC, but potentially continue processing others
				// For now, let's skip this PVC if we can't get its pod info
				// Consider logging: log.Printf("Warning: could not get pods for PVC %s: %v", pvc.Name, err)
				continue
			}

			if len(pods) > 0 {
				foundServiceIDs := make(map[string]struct{}) // To store unique service IDs
				anyPodMissingServiceLabel := false

				for _, pod := range pods {
					podServiceLabel := pod.GetLabels()[serviceLabel]
					if podServiceLabel == "" {
						anyPodMissingServiceLabel = true
						break
					}
					foundServiceIDs[podServiceLabel] = struct{}{}
				}

				if anyPodMissingServiceLabel {
					includePVC = false // Omit if any bound pod is missing the service label
				} else {
					for id := range foundServiceIDs {
						boundToServiceIDs = append(boundToServiceIDs, id)
					}
				}
			}
			// If len(pods) == 0 but PVC is Bound, it remains includable, BoundToServiceIDs will be empty.
		}

		if !includePVC {
			continue
		}

		result = append(result, PVCInfo{
			Name:              pvc.Name,
			Size:              size,
			TeamID:            teamID,
			ProjectID:         utils.ToPtr(projectID),
			EnvironmentID:     utils.ToPtr(environmentID),
			BoundToServiceIDs: boundToServiceIDs,
			Status:            string(pvc.Status.Phase),
		})
	}

	return result, nil
}

// DeletePersistentVolumeClaim deletes a specific PersistentVolumeClaim by its name and namespace.
func (k *KubeClient) DeletePersistentVolumeClaim(ctx context.Context, namespace string, pvcName string) error {
	if namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}
	if pvcName == "" {
		return fmt.Errorf("pvcName cannot be empty")
	}

	err := k.clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, pvcName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete PersistentVolumeClaim '%s' in namespace '%s': %w", pvcName, namespace, err)
	}
	return nil
}

// GetPodsUsingPVC finds all pods in a given namespace that are mounting the specified PVC.
func (k *KubeClient) GetPodsUsingPVC(ctx context.Context, namespace string, pvcName string) ([]corev1.Pod, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}
	if pvcName == "" {
		return nil, fmt.Errorf("pvcName cannot be empty")
	}

	podList, err := k.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace '%s': %w", namespace, err)
	}

	var podsUsingPVC []corev1.Pod
	for _, pod := range podList.Items {
		for _, volume := range pod.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil && volume.PersistentVolumeClaim.ClaimName == pvcName {
				podsUsingPVC = append(podsUsingPVC, pod)
				break // Move to the next pod once a match is found for this pod
			}
		}
	}
	return podsUsingPVC, nil
}

// ResizePersistentVolumeClaim updates the requested storage size of an existing PVC.
// Note: The StorageClass must support volume expansion (allowVolumeExpansion: true).
// The actual resize is handled by the storage provisioner.
func (k *KubeClient) ResizePersistentVolumeClaim(ctx context.Context, namespace string, pvcName string, newSize string) (*corev1.PersistentVolumeClaim, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}
	if pvcName == "" {
		return nil, fmt.Errorf("pvcName cannot be empty")
	}
	if newSize == "" {
		return nil, fmt.Errorf("newSize cannot be empty")
	}

	pvc, err := k.GetPersistentVolumeClaim(ctx, namespace, pvcName)
	if err != nil {
		return nil, fmt.Errorf("failed to get PVC '%s' for resizing: %w", pvcName, err)
	}

	newStorageQuantity, err := resource.ParseQuantity(newSize)
	if err != nil {
		return nil, fmt.Errorf("failed to parse newSize '%s': %w", newSize, err)
	}

	// Ensure the new size is greater than the current size if already bound, some providers might reject same or smaller size updates.
	// However, the API itself might allow it, and it's up to the controller to validate.
	// For simplicity, we'll allow setting any valid quantity here.
	pvc.Spec.Resources.Requests[corev1.ResourceStorage] = newStorageQuantity

	updatedPvc, err := k.clientset.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, pvc, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update (resize) PersistentVolumeClaim '%s': %w", pvcName, err)
	}

	return updatedPvc, nil
}
