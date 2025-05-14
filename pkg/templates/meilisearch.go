package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// MeiliSearchTemplate returns the predefined MeiliSearch template
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
				Description: "Hostname to use for the MeiliSearch API server.",
				Required:    true,
				TargetPort:  utils.ToPtr(7700),
			},
			{
				ID:          2,
				Name:        "Storage Size",
				Type:        schema.InputTypeVolumeSize,
				Description: "Size of the persistent storage for MeiliSearch data.",
				Required:    true,
				Default:     utils.ToPtr("1Gi"),
			},
		},
		Services: []schema.TemplateService{
			{
				ID:           1,
				Name:         "MeiliSearch",
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				HostInputIDs: []int{1},
				Image:        utils.ToPtr("getmeili/meilisearch:v1.14"),
				Ports: []schema.PortSpec{
					{
						Port:     7700,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic: true,
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
				Volumes: []schema.TemplateVolume{
					{
						Name: "meilisearch-data",
						Size: schema.TemplateVolumeSize{
							FromInputID: 2,
						},
						MountPath: "/meili_data",
					},
				},
			},
		},
	}
}
