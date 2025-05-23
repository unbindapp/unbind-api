package k8s

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/internal/common/utils"
	"github.com/unbindapp/unbind-api/internal/models"
	v1 "github.com/unbindapp/unbind-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DiscoverEndpointsByLabels returns both internal (services) and external (ingresses) endpoints
// matching the provided labels in a namespace
func (self *KubeClient) DiscoverEndpointsByLabels(ctx context.Context, namespace string, labels map[string]string, ignoreDNSChecks bool, client *kubernetes.Clientset) (*models.EndpointDiscovery, error) {
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

		// Only process ClusterIP services as internal
		if svc.Spec.Type == corev1.ServiceTypeClusterIP {
			endpoint := models.ServiceEndpoint{
				KubernetesName: svc.Name,
				DNS:            fmt.Sprintf("%s.%s", svc.Name, namespace),
				Ports:          make([]schema.PortSpec, len(svc.Spec.Ports)),
				TeamID:         teamID,
				ProjectID:      projectID,
				EnvironmentID:  environmentID,
				ServiceID:      serviceID,
			}

			// Add port information
			for i, port := range svc.Spec.Ports {
				endpoint.Ports[i] = schema.PortSpec{
					Port:     port.Port,
					Protocol: utils.ToPtr(schema.Protocol(port.Protocol)),
				}
			}

			discovery.Internal = append(discovery.Internal, endpoint)
		} else if svc.Spec.Type == corev1.ServiceTypeNodePort || svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
			// Process NodePort and LoadBalancer services as external
			endpoint := models.IngressEndpoint{
				KubernetesName: svc.Name,
				Hosts:          []models.ExtendedHostSpec{},
				TeamID:         teamID,
				ProjectID:      projectID,
				EnvironmentID:  environmentID,
				ServiceID:      serviceID,
			}

			// Get the node IPs, use internal client for this
			nodes, err := self.GetInternalClient().CoreV1().Nodes().List(ctx, metav1.ListOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to list nodes: %w", err)
			}

			var nodeIPs []string
			for _, node := range nodes.Items {
				for _, addr := range node.Status.Addresses {
					if addr.Type == corev1.NodeExternalIP {
						nodeIPs = append(nodeIPs, addr.Address)
						break
					}
				}
			}

			// Add each port as a host with the node IPs
			for _, port := range svc.Spec.Ports {
				if port.NodePort > 0 {
					for _, nodeIP := range nodeIPs {
						endpoint.Hosts = append(endpoint.Hosts, models.ExtendedHostSpec{
							HostSpec: v1.HostSpec{
								Host: nodeIP,
								Path: "/",
								Port: utils.ToPtr(port.NodePort),
							},
							Issued:        false,
							DnsConfigured: true,
							Cloudflare:    false,
						})
					}
				}
			}

			// Add LoadBalancer external IPs if available
			if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
				for _, ingress := range svc.Status.LoadBalancer.Ingress {
					if ingress.IP != "" {
						for _, port := range svc.Spec.Ports {
							// Also add the external IP with the NodePort if it exists
							if port.NodePort > 0 {
								host := fmt.Sprintf("%s:%d", ingress.IP, port.NodePort)
								endpoint.Hosts = append(endpoint.Hosts, models.ExtendedHostSpec{
									HostSpec: v1.HostSpec{
										Host: host,
										Path: "/",
										Port: utils.ToPtr(port.NodePort),
									},
									Issued:        false,
									DnsConfigured: true,
									Cloudflare:    false,
								})
							}
						}
					}
				}
			}

			if len(endpoint.Hosts) > 0 {
				discovery.External = append(discovery.External, endpoint)
			}
		}
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
			KubernetesName: ing.Name,
			IsIngress:      true,
			Hosts:          []models.ExtendedHostSpec{},
			TeamID:         teamID,
			ProjectID:      projectID,
			EnvironmentID:  environmentID,
			ServiceID:      serviceID,
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
				dnsConfigured := false
				cloudflare := false

				if !ignoreDNSChecks {
					if tls.SecretName != "" {
						secret, err := client.CoreV1().Secrets(namespace).Get(ctx, tls.SecretName, metav1.GetOptions{})
						issued = err == nil && isCertificateIssued(secret)
					}

					ips, err := self.GetIngressNginxIP(ctx)
					if err != nil {
						return nil, fmt.Errorf("failed to get ingress nginx IP: %w", err)
					}
					// Check ipv4 first
					dnsConfigured, _ = self.dnsChecker.IsPointingToIP(host, ips.IPv4)
					if !dnsConfigured {
						// Check ipv6
						dnsConfigured, _ = self.dnsChecker.IsPointingToIP(host, ips.IPv6)
					}
					if !dnsConfigured {
						// Check cloudflare
						cloudflare, _ = self.dnsChecker.IsUsingCloudflareProxy(host)

						if cloudflare {
							// Spin up an ingress to verify
							testIngress, path, err := self.CreateVerificationIngress(ctx, host, self.GetInternalClient())
							if err != nil {
								log.Warnf("Error creating ingress test for domain %s: %v", host, err)
							} else {
								defer func() {
									err := self.DeleteVerificationIngress(ctx, testIngress.Name, self.GetInternalClient())
									if err != nil {
										log.Warnf("Error deleting ingress test for domain %s: %v", host, err)
									}
								}()

								url := fmt.Sprintf("https://%s/.unbind-challenge/%s", host, path)

								// Create a new request with context
								req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
								if err != nil {
									log.Warnf("Error creating HTTP request for domain %s: %v", host, err)
								} else {
									// Retry delaying 200ms between tries
									maxRetries := 20
									for attempt := 0; attempt < maxRetries; attempt++ {
										if attempt > 0 {
											time.Sleep(200 * time.Millisecond)
										}

										// Execute the request
										resp, err := self.httpClient.Do(req)
										if err != nil {
											log.Warnf("Attempt %d: Error executing HTTP request for domain %s: %v", attempt+1, host, err)
											continue // Try again after sleep
										}

										func() {
											defer resp.Body.Close()
											// Check for the special header
											if resp.Header.Get("X-DNS-Check") == "resolved" {
												dnsConfigured = true
												return // Exit the closure
											}
										}()

										if dnsConfigured {
											break
										}
									}
								}
							}
						}
					}
				}

				endpoint.Hosts = append(endpoint.Hosts, models.ExtendedHostSpec{
					HostSpec: v1.HostSpec{
						Host: host,
						Path: path,
						Port: utils.ToPtr[int32](443),
					},
					Issued:        issued,
					DnsConfigured: dnsConfigured,
					Cloudflare:    cloudflare,
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

// CreateVerificationIngress creates an ingress with a configuration snippet to help verify
// that a domain is pointing to the Kubernetes cluster
func (self *KubeClient) CreateVerificationIngress(
	ctx context.Context,
	domain string,
	client *kubernetes.Clientset,
) (*networkingv1.Ingress, string, error) {
	pathType := networkingv1.PathTypeImplementationSpecific
	ingressClassName := "nginx"

	name, err := utils.GenerateSlug(domain)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate slug for domain %s: %w", domain, err)
	}

	path := fmt.Sprintf("/.unbind-challenge/dns-check/%s", uuid.NewString())

	// Create the ingress specification
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: self.config.GetSystemNamespace(),
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/ssl-redirect": "false",
				"nginx.ingress.kubernetes.io/configuration-snippet": `
add_header X-DNS-Check "resolved" always;
return 200 "Domain verification successful: CloudFlare connection to Kubernetes confirmed";
add_header Content-Type text/plain;
`,
			},
			Labels: map[string]string{
				"app":       "unbind-verification",
				"type":      "domain-verification",
				"domain":    domain,
				"temporary": "true",
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &ingressClassName,
			Rules: []networkingv1.IngressRule{
				{
					Host: domain,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     path,
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "dummy-service",
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Create the ingress in the cluster
	ing, err := client.NetworkingV1().Ingresses(self.config.GetSystemNamespace()).Create(ctx, ingress, metav1.CreateOptions{})
	return ing, path, err
}

// DeleteVerificationIngress deletes the verification ingress for a domain
func (self *KubeClient) DeleteVerificationIngress(
	ctx context.Context,
	ingressName string,
	client *kubernetes.Clientset,
) error {
	// Delete the ingress
	err := client.NetworkingV1().Ingresses(self.config.GetSystemNamespace()).Delete(ctx, ingressName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete verification ingress %s: %w", ingressName, err)
	}

	return nil
}

// DeleteOldVerificationIngresses deletes verification ingresses created more than 10 minutes ago
func (self *KubeClient) DeleteOldVerificationIngresses(
	ctx context.Context,
	client *kubernetes.Clientset,
) error {
	// Create a label selector for the verification ingresses
	labelSelector := "app=unbind-verification,type=domain-verification,temporary=true"

	// Get the list of ingresses matching the label selector
	ingresses, err := client.NetworkingV1().Ingresses(self.config.GetSystemNamespace()).List(
		ctx,
		metav1.ListOptions{
			LabelSelector: labelSelector,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to list verification ingresses: %w", err)
	}

	// Get the current time to compare against creation time
	currentTime := time.Now()
	cutoffTime := currentTime.Add(-10 * time.Minute)

	// Delete each matching ingress that is older than 10 minutes
	for _, ingress := range ingresses.Items {
		creationTime := ingress.GetCreationTimestamp().Time

		// Skip ingresses that are less than 10 minutes old
		if creationTime.After(cutoffTime) {
			continue
		}

		err := client.NetworkingV1().Ingresses(self.config.GetSystemNamespace()).Delete(
			ctx,
			ingress.Name,
			metav1.DeleteOptions{},
		)
		if err != nil {
			return fmt.Errorf("failed to delete old verification ingress %s: %w", ingress.Name, err)
		}
	}

	return nil
}
