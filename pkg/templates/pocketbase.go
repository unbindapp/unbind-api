package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// PocketBaseTemplate returns the predefined PocketBase template
func pocketBaseTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "PocketBase",
		DisplayRank: uint(20000),
		Icon:        "pocketbase",
		Keywords:    []string{"pocketbase", "database", "backend", "supabase", "firebase"},
		Description: "Open source backend in 1 file",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          "input_domain",
				Name:        "Domain",
				Type:        schema.InputTypeHost,
				Description: "The domain to use for the PocketBase instance.",
				Required:    true,
			},
			{
				ID:   "input_storage_size",
				Name: "Storage Size",
				Type: schema.InputTypeVolumeSize,
				Volume: &schema.TemplateVolume{
					Name:      "pb-data",
					MountPath: "/pb_data",
				},
				Description: "Size of the storage for the PocketBase data.",
				Required:    true,
				Default:     utils.ToPtr("1"),
			},
		},
		Services: []schema.TemplateService{
			{
				ID:       "service_pocketbase",
				Name:     "PocketBase",
				Type:     schema.ServiceTypeDockerimage,
				Builder:  schema.ServiceBuilderDocker,
				InputIDs: []string{"input_domain"},
				Image:    utils.ToPtr("ghcr.io/unbindapp/pocketbase:v0.28.1"),
				Ports: []schema.PortSpec{
					{
						Port:     8090,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/api/health",
					Port:                      utils.ToPtr(int32(8090)),
					PeriodSeconds:             5,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   3,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 3,
				},
				Variables: []schema.TemplateVariable{
					{
						Name: "PB_ADMIN_EMAIL",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypeEmail,
						},
					},
					{
						Name: "PB_ADMIN_PASSWORD",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
						},
					},
				},
			},
		},
	}
}
