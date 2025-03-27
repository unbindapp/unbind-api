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

// ServiceParams contains all parameters needed to create a v1.Service object
type ServiceParams struct {
	// Basic service information
	Name        string
	DisplayName string
	Description string
	Namespace   string

	// Service type configuration
	Provider         string
	Framework        string
	TeamRef          string
	ProjectRef       string
	EnvironmentID    string
	KubernetesSecret string

	// Git configuration
	GitRepoURL           string
	GitRef               string
	GithubInstallationID *int64

	// Deployment configuration
	Image    string
	Hosts    []v1.HostSpec
	Ports    []v1.PortSpec
	Public   *bool
	Replicas *int32
}

// CreateServiceObject creates a new v1.Service object with the provided parameters
func CreateServiceObject(params ServiceParams) (*v1.Service, error) {
	// Extract GitHub repository name from the Git URL
	gitRepository := extractGitRepository(params.GitRepoURL)

	// Generate a sanitized service name if not provided
	serviceName := params.Name
	if serviceName == "" {
		return nil, fmt.Errorf("service name cannot be empty")
	}

	// Create a new Service CR using the official structs
	service := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "unbind.unbind.app/v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: params.Namespace,
		},
		Spec: v1.ServiceSpec{
			Name:             serviceName,
			DisplayName:      params.DisplayName,
			Description:      params.Description,
			Type:             "git",
			Builder:          "railpack",
			Provider:         params.Provider,
			Framework:        params.Framework,
			TeamRef:          params.TeamRef,
			ProjectRef:       params.ProjectRef,
			EnvironmentID:    params.EnvironmentID,
			KubernetesSecret: params.KubernetesSecret,
			GitRepository:    gitRepository,
		},
	}

	// Set GitHub installation ID if provided
	if params.GithubInstallationID != nil {
		service.Spec.GitHubInstallationID = params.GithubInstallationID
	}

	// Build service configuration
	service.Spec.Config = v1.ServiceConfigSpec{
		GitBranch: params.GitRef,
		Image:     params.Image,
	}

	// Add host configuration if provided
	if len(params.Hosts) > 0 {
		service.Spec.Config.Hosts = params.Hosts
	}

	// Add port configuration if provided
	if len(params.Ports) > 0 {
		service.Spec.Config.Ports = params.Ports
	}

	// Set public flag if provided
	if params.Public != nil {
		service.Spec.Config.Public = *params.Public
	}

	// Set replicas if provided
	if params.Replicas != nil {
		replicas := int32(*params.Replicas)
		service.Spec.Config.Replicas = &replicas
	}

	return service, nil
}

// DeployImage creates (or replaces) the service resource in the target namespace
// for deployment after a successful build job.
func (self *K8SClient) DeployImage(ctx context.Context, crdName, image string) (*unstructured.Unstructured, *v1.Service, error) {
	// Generate a sanitized service name from the repo name
	serviceName := strings.ToLower(strings.ReplaceAll(crdName, "_", "-"))

	params := ServiceParams{
		Name:             serviceName,
		DisplayName:      serviceName,
		Description:      fmt.Sprintf("Auto-deployed service for %s", crdName),
		Namespace:        self.config.DeploymentNamespace,
		Provider:         self.config.ServiceProvider,
		Framework:        self.config.ServiceFramework,
		TeamRef:          self.config.ServiceTeamRef,
		ProjectRef:       self.config.ServiceProjectRef,
		EnvironmentID:    self.config.ServiceProjectRef,
		KubernetesSecret: self.config.ServiceSecretName,
		GitRepoURL:       self.config.GitRepoURL,
		GitRef:           self.config.GitRef,
		Image:            image,
		Hosts:            self.config.Hosts,
		Ports:            self.config.Ports,
		Public:           self.config.ServicePublic,
		Replicas:         self.config.ServiceReplicas,
	}

	// Set GitHub installation ID if provided
	if self.config.GithubInstallationID != 0 {
		installationID := self.config.GithubInstallationID
		params.GithubInstallationID = &installationID
	}

	// Create the Service object
	service, err := CreateServiceObject(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create service object: %v", err)
	}

	// Build service configuration
	service.Spec.Config = v1.ServiceConfigSpec{
		GitBranch:  self.config.GitRef,
		AutoDeploy: true,
		Image:      image,
	}

	// Convert to unstructured for the dynamic client
	unstructuredObj, err := convertToUnstructured(service)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert service to unstructured: %v", err)
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
			res, err := updateExistingServiceCR(ctx, self, serviceGVR, unstructuredObj)
			return res, service, err
		}
		return nil, nil, fmt.Errorf("failed to create service custom resource: %v", err)
	}

	return createdCR, service, nil
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
