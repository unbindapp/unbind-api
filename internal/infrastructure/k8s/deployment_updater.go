package k8s

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// UpdateDeploymentImages updates container images in deployments based on the new version
func (k *KubeClient) UpdateDeploymentImages(ctx context.Context, newVersion string) error {
	// Get all deployments in the system namespace
	deployments, err := k.clientset.AppsV1().Deployments(k.config.GetSystemNamespace()).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	// Map of image prefixes to their new versions
	imageUpdates := map[string]string{
		"ghcr.io/unbindapp/unbind-ui":       fmt.Sprintf("ghcr.io/unbindapp/unbind-ui:%s", newVersion),
		"ghcr.io/unbindapp/unbind-operator": fmt.Sprintf("ghcr.io/unbindapp/unbind-operator:%s", newVersion),
		"ghcr.io/unbindapp/unbind-api":      fmt.Sprintf("ghcr.io/unbindapp/unbind-api:%s", newVersion),
	}

	// First update all non-api deployments
	for _, deployment := range deployments.Items {
		// Skip if this is the API deployment
		if strings.Contains(deployment.Name, "api") {
			continue
		}

		updated := false
		for i, container := range deployment.Spec.Template.Spec.Containers {
			for prefix, newImage := range imageUpdates {
				if strings.HasPrefix(container.Image, prefix+":") {
					deployment.Spec.Template.Spec.Containers[i].Image = newImage
					updated = true
				}
			}
		}

		if updated {
			// Update the deployment
			_, err := k.clientset.AppsV1().Deployments(k.config.GetSystemNamespace()).Update(ctx, &deployment, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("failed to update deployment %s: %w", deployment.Name, err)
			}
		}
	}

	// Then update the API deployment
	for _, deployment := range deployments.Items {
		// Only process API deployment
		if !strings.Contains(deployment.Name, "api") {
			continue
		}

		updated := false
		for i, container := range deployment.Spec.Template.Spec.Containers {
			for prefix, newImage := range imageUpdates {
				if strings.HasPrefix(container.Image, prefix+":") {
					deployment.Spec.Template.Spec.Containers[i].Image = newImage
					updated = true
				}
			}
		}

		if updated {
			// Update the deployment
			_, err := k.clientset.AppsV1().Deployments(k.config.GetSystemNamespace()).Update(ctx, &deployment, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("failed to update deployment %s: %w", deployment.Name, err)
			}
		}
	}

	return nil
}

// CheckDeploymentsReady checks if all deployments are running with the specified version
func (k *KubeClient) CheckDeploymentsReady(ctx context.Context, version string) (bool, error) {
	// Get all deployments in the system namespace
	deployments, err := k.clientset.AppsV1().Deployments(k.config.GetSystemNamespace()).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to list deployments: %w", err)
	}

	// Expected images for each component
	expectedImages := map[string]string{
		"ghcr.io/unbindapp/unbind-ui":       fmt.Sprintf("ghcr.io/unbindapp/unbind-ui:%s", version),
		"ghcr.io/unbindapp/unbind-operator": fmt.Sprintf("ghcr.io/unbindapp/unbind-operator:%s", version),
		"ghcr.io/unbindapp/unbind-api":      fmt.Sprintf("ghcr.io/unbindapp/unbind-api:%s", version),
	}

	// Check each deployment
	for _, deployment := range deployments.Items {
		// Get the deployment's pods
		pods, err := k.clientset.CoreV1().Pods(k.config.GetSystemNamespace()).List(ctx, metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(metav1.SetAsLabelSelector(deployment.Spec.Selector.MatchLabels)),
		})
		if err != nil {
			return false, fmt.Errorf("failed to list pods for deployment %s: %w", deployment.Name, err)
		}

		// Check if any pods are running with the new version
		hasNewVersion := false
		for _, pod := range pods.Items {
			// Check if pod is running
			if pod.Status.Phase != corev1.PodRunning {
				continue
			}

			// Check if pod has the new version
			for _, container := range pod.Spec.Containers {
				for prefix, expectedImage := range expectedImages {
					if strings.HasPrefix(container.Image, prefix+":") {
						if container.Image == expectedImage {
							hasNewVersion = true
							break
						}
					}
				}
			}
		}

		if !hasNewVersion {
			return false, nil
		}
	}

	return true, nil
}
