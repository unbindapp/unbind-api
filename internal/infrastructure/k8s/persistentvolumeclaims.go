package k8s

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/internal/common/errdefs"
	"github.com/unbindapp/unbind-api/internal/models"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreatePersistentVolumeClaim creates a new PersistentVolumeClaim in the specified namespace.
func (self *KubeClient) CreatePersistentVolumeClaim(
	ctx context.Context,
	namespace string,
	pvcName string,
	displayName string,
	labels map[string]string,
	storageRequest string,
	accessModes []corev1.PersistentVolumeAccessMode,
	storageClassName *string,
	client *kubernetes.Clientset,
) (*models.PVCInfo, error) {
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

	labels["pvc-display-name"] = displayName

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

	_, err = client.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, pvc, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create PersistentVolumeClaim '%s' in namespace '%s': %w", pvcName, namespace, err)
	}

	// Return the created PVC info using GetPersistentVolumeClaim
	return self.GetPersistentVolumeClaim(ctx, namespace, pvcName, client)
}

// UpdatePersistentVolumeClaim updates an existing PersistentVolumeClaim with new parameters (size, name)
func (self *KubeClient) UpdatePersistentVolumeClaim(
	ctx context.Context,
	namespace string,
	pvcName string,
	newSize *string,
	client *kubernetes.Clientset,
) (*models.PVCInfo, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}
	if pvcName == "" {
		return nil, fmt.Errorf("pvcName cannot be empty")
	}

	// Get the existing PVC
	pvc, err := client.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, fmt.Sprintf("PersistentVolumeClaim '%s' not found", pvcName))
		}
		return nil, fmt.Errorf("failed to get PersistentVolumeClaim '%s' in namespace '%s': %w", pvcName, namespace, err)
	}

	// Update the PVC size if provided
	if newSize != nil {
		newStorageQuantity, err := resource.ParseQuantity(*newSize)
		if err != nil {
			return nil, fmt.Errorf("failed to parse newSize '%s': %w", *newSize, err)
		}
		pvc.Spec.Resources.Requests[corev1.ResourceStorage] = newStorageQuantity
	}

	// Update the PVC in Kubernetes
	_, err = client.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, pvc, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update PersistentVolumeClaim '%s' in namespace '%s': %w", pvcName, namespace, err)
	}
	// Return the updated PVC info using GetPersistentVolumeClaim
	return self.GetPersistentVolumeClaim(ctx, namespace, pvcName, client)
}

// GetPersistentVolumeClaim retrieves a specific PersistentVolumeClaim by its name and namespace.
func (self *KubeClient) GetPersistentVolumeClaim(ctx context.Context, namespace string, pvcName string, client *kubernetes.Clientset) (*models.PVCInfo, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}
	if pvcName == "" {
		return nil, fmt.Errorf("pvcName cannot be empty")
	}

	pvc, err := client.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errdefs.NewCustomError(errdefs.ErrTypeNotFound, fmt.Sprintf("PersistentVolumeClaim '%s' not found", pvcName))
		}
		return nil, fmt.Errorf("failed to get PersistentVolumeClaim '%s' in namespace '%s': %w", pvcName, namespace, err)
	}

	const ( // Define label keys for consistency
		teamLabel        = "unbind-team"
		projectLabel     = "unbind-project"
		environmentLabel = "unbind-environment"
		serviceLabel     = "unbind-service"
	)

	pvcLabels := pvc.GetLabels()
	teamIDStr := pvcLabels[teamLabel]
	// Skip if the PVC doesn't have the unbind-team label
	if teamIDStr == "" {
		return nil, fmt.Errorf("PVC '%s' does not have required team label", pvcName)
	}

	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid team ID in PVC '%s': %w", pvcName, err)
	}

	projectIDStr := pvcLabels[projectLabel]
	environmentIDStr := pvcLabels[environmentLabel]
	sizeGBValueStr := ""
	var sizeGBValue float64
	if storageRequest, ok := pvc.Spec.Resources.Requests[corev1.ResourceStorage]; ok {
		bytesValue := storageRequest.Value()
		gbValue := float64(bytesValue) / (1024 * 1024 * 1024)
		sizeGBValueStr = fmt.Sprintf("%.2f", gbValue) // Format to 2 decimal places
		// if gbValue is a whole number, remove .00
		if strings.HasSuffix(sizeGBValueStr, ".00") {
			sizeGBValueStr = strings.TrimSuffix(sizeGBValueStr, ".00")
		}
		sizeGBValue, err = strconv.ParseFloat(sizeGBValueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sizeGBValue '%s': %w", sizeGBValueStr, err)
		}
	}

	var boundToServiceID *uuid.UUID

	// Check if bound to pods with unbind-service label
	pods, err := self.GetPodsUsingPVC(ctx, pvc.Namespace, pvc.Name, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods using PVC '%s': %w", pvcName, err)
	}

	isBound := len(pods) > 0
	isDatabase := false
	for _, pod := range pods {
		// Check if pod has database label
		// unbind/usd-category : databases
		if podLabel, ok := pod.GetLabels()["unbind/usd-category"]; ok && podLabel == "databases" {
			isDatabase = true
		}
		podServiceLabel := pod.GetLabels()[serviceLabel]
		if podServiceLabel != "" {
			// Parse the service ID from the label
			serviceID, err := uuid.Parse(podServiceLabel)
			if err != nil {
				continue
			}
			boundToServiceID = &serviceID
			break
		}
	}

	if isBound && boundToServiceID == nil {
		return nil, fmt.Errorf("PVC '%s' is bound but no valid service ID found", pvcName)
	}

	var projectID *uuid.UUID
	if projectIDStr != "" {
		projectIDParsed, err := uuid.Parse(projectIDStr)
		if err == nil {
			projectID = &projectIDParsed
		}
	}

	var environmentID *uuid.UUID
	if environmentIDStr != "" {
		environmentIDParsed, err := uuid.Parse(environmentIDStr)
		if err == nil {
			environmentID = &environmentIDParsed
		}
	}

	// Check if PVC can be deleted (no owners and not in use)
	canDelete := len(pvc.OwnerReferences) == 0 && !isBound

	// Get type
	pvcType := models.PvcScopeTeam
	if projectID != nil && environmentID != nil {
		pvcType = models.PvcScopeEnvironment
	} else if projectID != nil {
		pvcType = models.PvcScopeProject
	}

	return &models.PVCInfo{
		ID:                 pvc.Name,
		Type:               pvcType,
		CapacityGB:         sizeGBValue,
		TeamID:             teamID,
		ProjectID:          projectID,
		EnvironmentID:      environmentID,
		MountedOnServiceID: boundToServiceID,
		Status:             models.PersistentVolumeClaimPhase(pvc.Status.Phase),
		IsDatabase:         isDatabase,
		IsAvailable:        canDelete,
		CanDelete:          canDelete,
		CreatedAt:          pvc.CreationTimestamp.Time,
	}, nil
}

// ListPersistentVolumeClaims lists all PersistentVolumeClaims in a given namespace, optionally filtered by a label selector,
func (self *KubeClient) ListPersistentVolumeClaims(ctx context.Context, namespace string, labels map[string]string, client *kubernetes.Clientset) ([]*models.PVCInfo, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}

	listOptions := metav1.ListOptions{}
	var selectors []string
	for key, value := range labels {
		selectors = append(selectors, fmt.Sprintf("%s=%s", key, value))
	}
	listOptions.LabelSelector = strings.Join(selectors, ",")

	pvcList, err := client.CoreV1().PersistentVolumeClaims(namespace).List(ctx, listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list PersistentVolumeClaims in namespace '%s' with selector '%s': %w", namespace, listOptions.LabelSelector, err)
	}

	// List all pods ONCE and build a map of PVC -> Pods using it
	podList, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace '%s': %w", namespace, err)
	}

	// Build map: PVC name -> list of pods using it
	pvcToPods := make(map[string][]corev1.Pod)
	for _, pod := range podList.Items {
		for _, volume := range pod.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil {
				pvcName := volume.PersistentVolumeClaim.ClaimName
				pvcToPods[pvcName] = append(pvcToPods[pvcName], pod)
			}
		}
	}

	var result []*models.PVCInfo
	const (
		teamLabel        = "unbind-team"
		projectLabel     = "unbind-project"
		environmentLabel = "unbind-environment"
		serviceLabel     = "unbind-service"
	)

	for _, pvc := range pvcList.Items {
		pvcLabels := pvc.GetLabels()
		teamIDStr := pvcLabels[teamLabel]
		// Skip if the PVC doesn't have the unbind-team label
		if teamIDStr == "" {
			continue
		}

		teamID, err := uuid.Parse(teamIDStr)

		// Skip if the team ID is not valid
		if err != nil {
			continue
		}

		projectIDStr := pvcLabels[projectLabel]
		environmentIDStr := pvcLabels[environmentLabel]
		displayName := pvcLabels["pvc-display-name"]
		if displayName == "" {
			displayName = pvc.Name
		}
		sizeGBValueStr := ""
		var sizeGBValue float64
		var bytesValue int64
		if pvc.Status.Capacity != nil {
			if capacityQuantity, ok := pvc.Status.Capacity[corev1.ResourceStorage]; ok {
				bytesValue = capacityQuantity.Value()
			}
		} else {
			// Fallback to requests if capacity is not set
			if storageRequest, ok := pvc.Spec.Resources.Requests[corev1.ResourceStorage]; ok {
				bytesValue = storageRequest.Value()
			}
		}
		gbValue := float64(bytesValue) / (1024 * 1024 * 1024)
		sizeGBValueStr = fmt.Sprintf("%.2f", gbValue) // Format to 2 decimal places
		// if gbValue is a whole number, remove .00
		if strings.HasSuffix(sizeGBValueStr, ".00") {
			sizeGBValueStr = strings.TrimSuffix(sizeGBValueStr, ".00")
		}
		sizeGBValue, err = strconv.ParseFloat(sizeGBValueStr, 64)
		if err != nil {
			continue
		}

		var boundToServiceID *uuid.UUID
		isDatabase := false

		// Rule 2: If bound to pods, they must have unbind-service label
		pods := pvcToPods[pvc.Name]
		isBound := len(pods) > 0

		for _, pod := range pods {
			// Check if pod has database label
			// unbind/usd-category : databases
			if podLabel, ok := pod.GetLabels()["unbind/usd-category"]; ok && podLabel == "databases" {
				isDatabase = true
			}

			podServiceLabel := pod.GetLabels()[serviceLabel]
			if podServiceLabel != "" {
				// Parse the service ID from the label
				serviceID, err := uuid.Parse(podServiceLabel)
				if err != nil {
					continue
				}
				boundToServiceID = &serviceID
				break
			}
		}

		if isBound && boundToServiceID == nil {
			// If the PVC is bound but no service ID is found, skip this PVC
			continue
		}

		var projectID *uuid.UUID
		if projectIDStr != "" {
			projectIDParsed, err := uuid.Parse(projectIDStr)
			if err == nil {
				projectID = &projectIDParsed
			}
		}

		var environmentID *uuid.UUID
		if environmentIDStr != "" {
			environmentIDParsed, err := uuid.Parse(environmentIDStr)
			if err == nil {
				environmentID = &environmentIDParsed
			}
		}
		// Check if PVC can be deleted (no owners and not in use)
		canDelete := len(pvc.OwnerReferences) == 0 && !isBound

		// Figure out type
		pvcType := models.PvcScopeTeam
		if projectID != nil && environmentID != nil {
			pvcType = models.PvcScopeEnvironment
		} else if projectID != nil {
			pvcType = models.PvcScopeProject
		}

		result = append(result, &models.PVCInfo{
			ID:                 pvc.Name,
			Type:               pvcType,
			CapacityGB:         sizeGBValue,
			TeamID:             teamID,
			ProjectID:          projectID,
			EnvironmentID:      environmentID,
			MountedOnServiceID: boundToServiceID,
			Status:             models.PersistentVolumeClaimPhase(pvc.Status.Phase),
			IsDatabase:         isDatabase,
			IsAvailable:        canDelete,
			CanDelete:          canDelete,
			CreatedAt:          pvc.CreationTimestamp.Time,
		})
	}

	// Sort the result by CreatedAt in descending order
	slices.SortFunc(result, func(a, b *models.PVCInfo) int {
		if a.CreatedAt.After(b.CreatedAt) {
			return -1
		} else if a.CreatedAt.Before(b.CreatedAt) {
			return 1
		}
		return 0
	})

	return result, nil
}

// DeletePersistentVolumeClaim deletes a specific PersistentVolumeClaim by its name and namespace.
func (self *KubeClient) DeletePersistentVolumeClaim(ctx context.Context, namespace string, pvcName string, client *kubernetes.Clientset) error {
	if namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}
	if pvcName == "" {
		return fmt.Errorf("pvcName cannot be empty")
	}

	// Get the PVC to check its owner references
	pvc, err := client.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return errdefs.NewCustomError(errdefs.ErrTypeNotFound, fmt.Sprintf("PersistentVolumeClaim '%s' not found", pvcName))
		}
		return fmt.Errorf("failed to get PersistentVolumeClaim '%s': %w", pvcName, err)
	}

	// Check if PVC has any owner references
	if len(pvc.OwnerReferences) > 0 {
		ownerNames := make([]string, 0, len(pvc.OwnerReferences))
		for _, owner := range pvc.OwnerReferences {
			ownerNames = append(ownerNames, fmt.Sprintf("%s/%s", owner.Kind, owner.Name))
		}
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput,
			fmt.Sprintf("Cannot delete PVC '%s' as it is owned by: %s", pvcName, strings.Join(ownerNames, ", ")))
	}

	// Check if PVC is in use by any pods
	pods, err := self.GetPodsUsingPVC(ctx, namespace, pvcName, client)
	if err != nil {
		return fmt.Errorf("failed to check if PVC is in use: %w", err)
	}
	if len(pods) > 0 {
		return errdefs.NewCustomError(errdefs.ErrTypeInvalidInput, fmt.Sprintf("Cannot delete PVC '%s' as it is currently in use by %d pod(s)", pvcName, len(pods)))
	}

	err = client.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, pvcName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete PersistentVolumeClaim '%s' in namespace '%s': %w", pvcName, namespace, err)
	}
	return nil
}

// GetPodsUsingPVC finds all pods in a given namespace that are mounting the specified PVC.
func (self *KubeClient) GetPodsUsingPVC(ctx context.Context, namespace string, pvcName string, client *kubernetes.Clientset) ([]corev1.Pod, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}
	if pvcName == "" {
		return nil, fmt.Errorf("pvcName cannot be empty")
	}

	podList, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
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
func (self *KubeClient) ResizePersistentVolumeClaim(ctx context.Context, namespace string, pvcName string, newSize string, client *kubernetes.Clientset) (*models.PVCInfo, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}
	if pvcName == "" {
		return nil, fmt.Errorf("pvcName cannot be empty")
	}
	if newSize == "" {
		return nil, fmt.Errorf("newSize cannot be empty")
	}

	// Get the raw PVC first since we need to modify it
	pvc, err := client.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get PVC '%s' for resizing: %w", pvcName, err)
	}

	newStorageQuantity, err := resource.ParseQuantity(newSize)
	if err != nil {
		return nil, fmt.Errorf("failed to parse newSize '%s': %w", newSize, err)
	}

	// Update the PVC size
	pvc.Spec.Resources.Requests[corev1.ResourceStorage] = newStorageQuantity

	// Update the PVC in Kubernetes
	_, err = client.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, pvc, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update (resize) PersistentVolumeClaim '%s': %w", pvcName, err)
	}

	// Return the updated PVC info
	return self.GetPersistentVolumeClaim(ctx, namespace, pvcName, client)
}
