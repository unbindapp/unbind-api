package k8s

import (
	"context"
	"fmt"
	"strings"

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
		"unbindapp/unbind-api":      fmt.Sprintf("unbindapp/unbind-api:%s", newVersion),
		"unbindapp/unbind-ui":       fmt.Sprintf("unbindapp/unbind-ui:%s", newVersion),
		"unbindapp/unbind-operator": fmt.Sprintf("unbindapp/unbind-operator:%s", newVersion),
	}

	// Update each deployment
	for _, deployment := range deployments.Items {
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

// WaitForDeploymentsReady waits for all deployments to be ready after an update
func (k *KubeClient) WaitForDeploymentsReady(ctx context.Context) error {
	deployments, err := k.clientset.AppsV1().Deployments(k.config.GetSystemNamespace()).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	for _, deployment := range deployments.Items {
		// Wait for deployment to be ready
		for {
			updated, err := k.clientset.AppsV1().Deployments(k.config.GetSystemNamespace()).Get(ctx, deployment.Name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get deployment %s: %w", deployment.Name, err)
			}

			if updated.Status.ReadyReplicas == *updated.Spec.Replicas {
				break
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Continue waiting
			}
		}
	}

	return nil
}
