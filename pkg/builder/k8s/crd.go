package k8s

import (
	"context"
	"fmt"
	"strings"

	// Import the operator API package
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
		Namespace:        self.builderConfig.DeploymentNamespace,
		Provider:         self.builderConfig.ServiceProvider,
		Framework:        self.builderConfig.ServiceFramework,
		TeamRef:          self.builderConfig.ServiceTeamRef,
		ProjectRef:       self.builderConfig.ServiceProjectRef,
		EnvironmentID:    self.builderConfig.ServiceEnvironmentRef,
		KubernetesSecret: self.builderConfig.ServiceSecretName,
		GitRepoURL:       self.builderConfig.GitRepoURL,
		GitRef:           self.builderConfig.GitRef,
		Image:            image,
		Hosts:            self.builderConfig.Hosts,
		Ports:            self.builderConfig.Ports,
		Public:           self.builderConfig.ServicePublic,
		Replicas:         self.builderConfig.ServiceReplicas,
	}

	// Set GitHub installation ID if provided
	if self.builderConfig.GithubInstallationID != 0 {
		installationID := self.builderConfig.GithubInstallationID
		params.GithubInstallationID = &installationID
	}

	// Create the Service object
	service, err := CreateServiceObject(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create service object: %v", err)
	}

	return self.k8s.DeployUnbindService(ctx, service)
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
