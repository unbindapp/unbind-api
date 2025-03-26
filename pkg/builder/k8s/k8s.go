package k8s

import (
	"github.com/unbindapp/unbind-api/internal/common/log"
	"github.com/unbindapp/unbind-api/pkg/builder/config"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8SClient struct {
	config    *config.Config
	namespace string
	client    *dynamic.DynamicClient
}

func NewK8SClient(cfg *config.Config) *K8SClient {
	// Get config
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		if cfg.KubeConfig != "" {
			// Use the configured kubeconfig file instead
			kubeConfig, err = clientcmd.BuildConfigFromFlags("", cfg.KubeConfig)
			if err != nil {
				log.Fatalf("Error building kubeconfig from %s: %v", cfg.KubeConfig, err)
			}
			log.Infof("Using kubeconfig from: %s", cfg.KubeConfig)
		} else {
			log.Fatalf("Error getting in-cluster config: %v", err)
		}
	}

	clientset, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	return &K8SClient{
		config:    cfg,
		client:    clientset,
		namespace: cfg.DeploymentNamespace,
	}
}
