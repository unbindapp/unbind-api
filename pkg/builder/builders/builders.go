package builders

import "github.com/unbindapp/unbind-api/pkg/builder/config"

type Builder struct {
	config *config.Config
}

func NewBuilder(config *config.Config) *Builder {
	return &Builder{
		config: config,
	}
}
