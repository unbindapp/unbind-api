package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// MeilisearchTemplate returns the predefined Meilisearch template
func meiliSearchTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "Meilisearch",
		DisplayRank: uint(80000),
		Icon:        "meilisearch",
		Keywords:    []string{"full text search", "elasticsearch", "search engine", "ram"},
		Description: "Fast & open source search engine",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "API Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the Meilisearch API server.",
				Required:    true,
				TargetPort:  utils.ToPtr(7700),
			},
			{
				ID:   2,
				Name: "Storage Size",
				Type: schema.InputTypeVolumeSize,
				Volume: &schema.TemplateVolume{
					Name:      "meilisearch-data",
					MountPath: "/meili_data",
				},
				Description: "Size of the persistent storage for Meilisearch data.",
				Required:    true,
				Default:     utils.ToPtr("1Gi"),
			},
		},
		Services: []schema.TemplateService{
			{
				ID:       1,
				Name:     "Meilisearch",
				Type:     schema.ServiceTypeDockerimage,
				Builder:  schema.ServiceBuilderDocker,
				InputIDs: []int{1, 2},
				Image:    utils.ToPtr("getmeili/meilisearch:v1.14"),
				Ports: []schema.PortSpec{
					{
						Port:     7700,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/health",
					Port:                      utils.ToPtr(int32(7700)),
					PeriodSeconds:             2,
					TimeoutSeconds:            10,
					StartupFailureThreshold:   15,
					LivenessFailureThreshold:  15,
					ReadinessFailureThreshold: 3,
				},
				Variables: []schema.TemplateVariable{
					{
						Name: "MEILI_MASTER_KEY",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
						},
					},
					{
						Name:  "MEILI_HTTP_ADDR",
						Value: "0.0.0.0:7700",
					},
					{
						Name:  "MEILI_ENV",
						Value: "production",
					},
				},
			},
		},
	}
}
