package k8s

import (
	"log"

	"github.com/unbindapp/unbind-api/config"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	config    config.ConfigInterface
	client    *dynamic.DynamicClient
	clientset *kubernetes.Clientset
}

func NewKubeClient(cfg config.ConfigInterface) *KubeClient {
	var kubeConfig *rest.Config
	var err error

	if cfg.GetKubeConfig() != "" {
		// Use provided kubeconfig if present
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", cfg.GetKubeConfig())
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

	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	return &KubeClient{
		config:    cfg,
		client:    dynamicClient,
		clientset: clientSet,
	}
}

// This function is used to manage unbind-system resources
func (self *KubeClient) GetInternalClient() *kubernetes.Clientset {
	return self.clientset
}

func (self *KubeClient) CreateClientWithToken(token string) (*kubernetes.Clientset, error) {
	config := &rest.Config{
		Host:        self.config.GetKubeProxyURL(),
		BearerToken: token,
		// Skip TLS verification for internal cluster communication
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}

	return kubernetes.NewForConfig(config)
}
