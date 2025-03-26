package k8s

import (
	"context"
	"fmt"
	"strings"

	// Import the operator API package
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DeployImage creates (or replaces) the service resource in the target namespace
// for deployment after a successful build job.
func (self *K8SClient) DeployImage(ctx context.Context, crdName, image string) (*unstructured.Unstructured, error) {
	// Extract GitHub repository name from the Git URL
	gitRepository := extractGitRepository(self.config.GitRepoURL)

	// Generate a sanitized service name from the repo name
	serviceName := strings.ToLower(strings.ReplaceAll(crdName, "_", "-"))

	// Create a new Service CR using the official structs
	service := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "unbind.unbind.app/v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: self.config.DeploymentNamespace,
		},
		Spec: v1.ServiceSpec{
			Name:             serviceName,
			DisplayName:      serviceName,
			Description:      fmt.Sprintf("Auto-deployed service for %s", crdName),
			Type:             "git",
			Builder:          "railpack",
			Provider:         self.config.ServiceProvider,
			Framework:        self.config.ServiceFramework,
			TeamRef:          "default",
			ProjectRef:       "default",
			EnvironmentID:    "prod",
			KubernetesSecret: self.config.ServiceSecretName,
			GitRepository:    gitRepository,
		},
	}

	// Set GitHub installation ID if provided
	if self.config.GithubInstallationID != 0 {
		installationID := self.config.GithubInstallationID
		service.Spec.GitHubInstallationID = &installationID
	}

	// Build service configuration
	service.Spec.Config = v1.ServiceConfigSpec{
		GitBranch:  self.config.GitRef,
		AutoDeploy: true,
		Image:      image,
	}

	// Add host configuration if provided
	if len(self.config.Hosts) > 0 {
		service.Spec.Config.Hosts = self.config.Hosts
	}

	// Add port configuration if provided
	if len(self.config.Ports) > 0 {
		service.Spec.Config.Ports = self.config.Ports
	}

	// Set public flag if provided
	if self.config.ServicePublic != nil {
		service.Spec.Config.Public = *self.config.ServicePublic
	}

	// Set replicas if provided
	if self.config.ServiceReplicas != nil {
		replicas := int32(*self.config.ServiceReplicas)
		service.Spec.Config.Replicas = &replicas
	}

	// Convert to unstructured for the dynamic client
	unstructuredObj, err := convertToUnstructured(service)
	if err != nil {
		return nil, fmt.Errorf("failed to convert service to unstructured: %v", err)
	}

	// Define the GroupVersionResource for the Service custom resource
	serviceGVR := schema.GroupVersionResource{
		Group:    "unbind.unbind.app",
		Version:  "v1",
		Resource: "services", // plural name of the custom resource
	}

	// Create the custom resource in the target namespace
	createdCR, err := self.client.Resource(serviceGVR).Namespace(self.config.DeploymentNamespace).Create(ctx, unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		// If the resource already exists, update it
		if apierrors.IsAlreadyExists(err) {
			return updateExistingServiceCR(ctx, self, serviceGVR, unstructuredObj)
		}
		return nil, fmt.Errorf("failed to create service custom resource: %v", err)
	}

	return createdCR, nil
}

// updateExistingServiceCR handles updating an existing Service custom resource
func updateExistingServiceCR(ctx context.Context, client *K8SClient, gvr schema.GroupVersionResource, newCR *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	// Retrieve the existing resource
	existingCR, err := client.client.Resource(gvr).Namespace(client.config.DeploymentNamespace).Get(ctx, newCR.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve existing service: %v", err)
	}

	// Set the resourceVersion on the object to be updated
	newCR.SetResourceVersion(existingCR.GetResourceVersion())

	// Update the CR
	updatedCR, err := client.client.Resource(gvr).Namespace(client.config.DeploymentNamespace).Update(ctx, newCR, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update service custom resource: %v", err)
	}

	return updatedCR, nil
}

// extractGitRepository parses a Git URL and extracts the repository information in the format "owner/repo"
func extractGitRepository(gitRepoURL string) string {
	parts := strings.Split(gitRepoURL, "/")
	if len(parts) < 2 {
		return ""
	}

	repoWithGit := parts[len(parts)-1]
	repo := strings.TrimSuffix(repoWithGit, ".git")

	if len(parts) >= 3 {
		owner := parts[len(parts)-2]
		return owner + "/" + repo
	}

	return repo
}

// convertToUnstructured converts a runtime.Object to an Unstructured object
func convertToUnstructured(obj runtime.Object) (*unstructured.Unstructured, error) {
	// Create a new unstructured object
	unstructuredObj := &unstructured.Unstructured{}

	// Convert the typed object to map[string]interface{}
	data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}

	// Set the unstructured data
	unstructuredObj.SetUnstructuredContent(data)

	return unstructuredObj, nil
}
