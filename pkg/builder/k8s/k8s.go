package k8s

import (
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/internal/infrastructure/k8s"
	builderConfig "github.com/unbindapp/unbind-api/pkg/builder/config"
)

type K8SClient struct {
	config        config.ConfigInterface
	builderConfig *builderConfig.Config
	namespace     string
	k8s           *k8s.KubeClient
}

func NewK8SClient(cfg config.ConfigInterface, builderConfig *builderConfig.Config) *K8SClient {
	// Get client
	kubeClient := k8s.NewKubeClient(cfg)

	return &K8SClient{
		config:        cfg,
		builderConfig: builderConfig,
		k8s:           kubeClient,
		namespace:     builderConfig.DeploymentNamespace,
	}
}
