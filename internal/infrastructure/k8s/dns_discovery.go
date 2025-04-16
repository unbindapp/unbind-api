package k8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// EndpointDiscovery contains both internal (services) and external (ingresses) endpoints
type EndpointDiscovery struct {
	Internal []ServiceEndpoint `json:"internal" nullable:"false"`
	External []IngressEndpoint `json:"external" nullable:"false"`
}

// ServiceEndpoint represents internal DNS information for a Kubernetes service
type ServiceEndpoint struct {
	Name          string            `json:"name"`
	DNS           string            `json:"dns"`
	Ports         []schema.PortSpec `json:"ports" nullable:"false"`
	TeamID        uuid.UUID         `json:"team_id"`
	ProjectID     uuid.UUID         `json:"project_id"`
	EnvironmentID uuid.UUID         `json:"environment_id"`
	ServiceID     uuid.UUID         `json:"service_id"`
}

// IngressEndpoint represents external DNS information for a Kubernetes ingress
type IngressEndpoint struct {
	Name          string        `json:"name"`
	Hosts         []v1.HostSpec `json:"hosts" nullable:"false"`
	TeamID        uuid.UUID     `json:"team_id"`
	ProjectID     uuid.UUID     `json:"project_id"`
	EnvironmentID uuid.UUID     `json:"environment_id"`
	ServiceID     uuid.UUID     `json:"service_id"`
}

// DiscoverEndpointsByLabels returns both internal (services) and external (ingresses) endpoints
// matching the provided labels in a namespace
func (k *KubeClient) DiscoverEndpointsByLabels(ctx context.Context, namespace string, labels map[string]string, client *kubernetes.Clientset) (*EndpointDiscovery, error) {
	// Convert the labels map to a selector string
	var labelSelectors []string
	for key, value := range labels {
		labelSelectors = append(labelSelectors, fmt.Sprintf("%s=%s", key, value))
	}
	labelSelector := strings.Join(labelSelectors, ",")

	discovery := &EndpointDiscovery{
		Internal: []ServiceEndpoint{},
		External: []IngressEndpoint{},
	}

	// Get services matching the label selector
	services, err := client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list services with labels %s: %w", labelSelector, err)
	}

	// Process services (internal endpoints)
	for _, svc := range services.Items {
		teamID, _ := uuid.Parse(svc.Labels["unbind-team"])
		projectID, _ := uuid.Parse(svc.Labels["unbind-project"])
		environmentID, _ := uuid.Parse(svc.Labels["unbind-environment"])
		serviceID, _ := uuid.Parse(svc.Labels["unbind-service"])

		endpoint := ServiceEndpoint{
			Name:          svc.Name,
			DNS:           fmt.Sprintf("%s.%s", svc.Name, namespace),
			Ports:         make([]schema.PortSpec, len(svc.Spec.Ports)),
			TeamID:        teamID,
			ProjectID:     projectID,
			EnvironmentID: environmentID,
			ServiceID:     serviceID,
		}

		// Add port information
		for i, port := range svc.Spec.Ports {
			endpoint.Ports[i] = schema.PortSpec{
				Port:     port.Port,
				Protocol: utils.ToPtr(schema.Protocol(port.Protocol)),
			}
		}

		discovery.Internal = append(discovery.Internal, endpoint)
	}

	// Get ingresses matching the label selector
	ingresses, err := client.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses with labels %s: %w", labelSelector, err)
	}

	// Process ingresses (external endpoints)
	for _, ing := range ingresses.Items {
		teamID, _ := uuid.Parse(ing.Labels["unbind-team"])
		projectID, _ := uuid.Parse(ing.Labels["unbind-project"])
		environmentID, _ := uuid.Parse(ing.Labels["unbind-environment"])
		serviceID, _ := uuid.Parse(ing.Labels["unbind-service"])

		endpoint := IngressEndpoint{
			Name:          ing.Name,
			Hosts:         []v1.HostSpec{},
			TeamID:        teamID,
			ProjectID:     projectID,
			EnvironmentID: environmentID,
			ServiceID:     serviceID,
		}

		// Make a map of paths to iterate TLS
		pathMap := make(map[string]string)

		for _, rule := range ing.Spec.Rules {
			host := rule.Host

			if rule.HTTP != nil {
				for _, path := range rule.HTTP.Paths {
					pathMap[host] = path.Path
				}
			}
		}

		// Only consider TLS for ingresses, get path from map above
		for _, tls := range ing.Spec.TLS {
			for _, host := range tls.Hosts {
				path := pathMap[host]

				endpoint.Hosts = append(endpoint.Hosts, v1.HostSpec{
					Host: host,
					Path: path,
				})
			}
		}

		discovery.External = append(discovery.External, endpoint)
	}

	return discovery, nil
}
