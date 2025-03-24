package k8s

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LoadBalancerAddresses contains the addresses for a load balancer service
type LoadBalancerAddresses struct {
	Name      string
	Namespace string
	IPv4      string
	IPv6      string
	Hostname  string
}

// GetLoadBalancerIPs returns the external IP addresses for load balancer services
// If labelSelector is provided, it will filter services based on the selector (e.g. "app.kubernetes.io/name=ingress-nginx")
func (k *KubeClient) GetLoadBalancerIPs(ctx context.Context, labelSelector string) ([]LoadBalancerAddresses, error) {
	var addresses []LoadBalancerAddresses

	// Get all namespaces
	namespaces, err := k.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	// For each namespace, get services of type LoadBalancer
	for _, ns := range namespaces.Items {
		services, err := k.clientset.CoreV1().Services(ns.Name).List(ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list services in namespace %s: %w", ns.Name, err)
		}

		for _, svc := range services.Items {
			// Skip services that are not LoadBalancers
			if svc.Spec.Type != corev1.ServiceTypeLoadBalancer {
				continue
			}

			lbAddrs := LoadBalancerAddresses{
				Name:      svc.Name,
				Namespace: svc.Namespace,
			}

			// Extract IP addresses from ingress
			if len(svc.Status.LoadBalancer.Ingress) > 0 {
				ingress := svc.Status.LoadBalancer.Ingress[0]

				// Set IPv4 if available
				if ingress.IP != "" {
					lbAddrs.IPv4 = ingress.IP
				}

				// For IPv6, we need to check the annotations or additional data
				// Some cloud providers add IPv6 as an annotation
				if v6IP, ok := svc.Annotations["ipv6.kubernetes.io/address"]; ok {
					lbAddrs.IPv6 = v6IP
				}

				// Alternatively, check for dual-stack IPs by looking at all ingress entries
				if lbAddrs.IPv6 == "" && len(svc.Status.LoadBalancer.Ingress) > 1 {
					for _, ing := range svc.Status.LoadBalancer.Ingress {
						// If we already have an IPv4 and this is a different IP, it might be IPv6
						if ing.IP != "" && ing.IP != lbAddrs.IPv4 && isIPv6(ing.IP) {
							lbAddrs.IPv6 = ing.IP
							break
						}
					}
				}

				// Set hostname if available (for providers like AWS)
				if ingress.Hostname != "" {
					lbAddrs.Hostname = ingress.Hostname
				}
			}

			addresses = append(addresses, lbAddrs)
		}
	}

	return addresses, nil
}

// isIPv6 checks if the given IP address is an IPv6 address
func isIPv6(ip string) bool {
	return strings.Count(ip, ":") >= 2
}

// GetIngressNginxIP is a convenience function to get the IP of the ingress-nginx controller
func (k *KubeClient) GetIngressNginxIP(ctx context.Context) (*LoadBalancerAddresses, error) {
	// Common label for ingress-nginx controller
	labelSelector := "app.kubernetes.io/name=ingress-nginx"

	addresses, err := k.GetLoadBalancerIPs(ctx, labelSelector)
	if err != nil {
		return nil, err
	}

	// Filter further for the controller service if needed
	for _, addr := range addresses {
		if addr.Name == "ingress-nginx-controller" {
			return &addr, nil
		}
	}

	// If we didn't find the specific controller name, return the first match if available
	if len(addresses) > 0 {
		return &addresses[0], nil
	}

	return nil, fmt.Errorf("no ingress-nginx load balancer found")
}
