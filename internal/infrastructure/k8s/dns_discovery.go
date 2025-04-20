package k8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/services/models"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DiscoverEndpointsByLabels returns both internal (services) and external (ingresses) endpoints
// matching the provided labels in a namespace
func (k *KubeClient) DiscoverEndpointsByLabels(ctx context.Context, namespace string, labels map[string]string, client *kubernetes.Clientset) (*models.EndpointDiscovery, error) {
	// Convert the labels map to a selector string
	var labelSelectors []string
	for key, value := range labels {
		labelSelectors = append(labelSelectors, fmt.Sprintf("%s=%s", key, value))
	}
	labelSelector := strings.Join(labelSelectors, ",")

	discovery := &models.EndpointDiscovery{
		Internal: []models.ServiceEndpoint{},
		External: []models.IngressEndpoint{},
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

		endpoint := models.ServiceEndpoint{
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

		endpoint := models.IngressEndpoint{
			Name:          ing.Name,
			Hosts:         []models.ExtendedHostSpec{},
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

				// Check if the secret is issued
				issued := false
				if tls.SecretName != "" {
					secret, err := client.CoreV1().Secrets(namespace).Get(ctx, tls.SecretName, metav1.GetOptions{})
					issued = err == nil && isCertificateIssued(secret)
				}

				endpoint.Hosts = append(endpoint.Hosts, models.ExtendedHostSpec{
					HostSpec: v1.HostSpec{
						Host: host,
						Path: path,
						Port: utils.ToPtr[int32](443),
					},
					Issued: issued,
				})
			}
		}

		discovery.External = append(discovery.External, endpoint)
	}

	return discovery, nil
}

// isCertificateIssued checks if a TLS secret contains valid certificate data
func isCertificateIssued(secret *corev1.Secret) bool {
	if secret == nil {
		return false
	}

	// Check if the secret contains the required TLS data
	_, hasCert := secret.Data["tls.crt"]
	_, hasKey := secret.Data["tls.key"]

	// Check if the secret has any cert-manager annotations
	hasCertManagerAnnotation := false
	if secret.Annotations != nil {
		for key := range secret.Annotations {
			if strings.Contains(key, "cert-manager") {
				hasCertManagerAnnotation = true
				break
			}
		}
	}

	// Both fields must exist and contain data, and it should have cert-manager annotations
	return hasCert && hasKey && len(secret.Data["tls.crt"]) > 0 && len(secret.Data["tls.key"]) > 0 && hasCertManagerAnnotation
}
