package k8s

// ? "Team" is synonymous with a Kubernetes namespace

import (
	"context"
	"fmt"

	"github.com/unbindapp/unbind-api/internal/common/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Gets specified namespaces
func (k *KubeClient) GetNamespaces(ctx context.Context, namespaceNames []string, bearerToken string) ([]*v1.Namespace, error) {
	client, err := k.CreateClientWithToken(bearerToken)
	if err != nil {
		log.Errorf("Error creating client with token: %v", err)
		return nil, fmt.Errorf("error creating client with token: %v", err)
	}

	namespaces := make([]*v1.Namespace, 0, len(namespaceNames))

	// Get each namespace by name instead of listing all
	for i, namespaceName := range namespaceNames {
		// Skip empty namespace names
		if namespaceName == "" {
			continue
		}

		// Get the specific namespace
		ns, err := client.CoreV1().Namespaces().Get(ctx, namespaceName, metav1.GetOptions{})
		if err != nil {
			// Log the error but continue with other namespaces
			log.Warnf("Error getting namespace %s: %v", namespaceName, err)
			continue
		}

		namespaces[i] = ns
	}

	return namespaces, nil
}
