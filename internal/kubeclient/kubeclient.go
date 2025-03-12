package kubeclient

import (
	"log"

	"github.com/unbindapp/unbind-api/config"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	config *config.Config
	client *dynamic.DynamicClient
}

func NewKubeClient(cfg *config.Config) *KubeClient {
	var kubeConfig *rest.Config
	var err error

	if cfg.KubeConfig != "" {
		// Use provided kubeconfig if present
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", cfg.KubeConfig)
		if err != nil {
			log.Fatalf("Error building kubeconfig: %v", err)
		}
	} else {
		// Fall back to in-cluster config
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Error getting in-cluster config: %v", err)
		}
	}

	clientset, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	return &KubeClient{
		config: cfg,
		client: clientset,
	}
}

func (k *KubeClient) createClientWithToken(token string) (*kubernetes.Clientset, error) {
	config := &rest.Config{
		Host:        k.config.KubeProxyURL,
		BearerToken: token,
		// Skip TLS verification for internal cluster communication
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
		Impersonate: rest.ImpersonationConfig{
			// ! TODO - dynamic groups
			UserName: "admin@unbind.app",
			Groups:   []string{"oidc:users"},
		},
	}

	return kubernetes.NewForConfig(config)
}
