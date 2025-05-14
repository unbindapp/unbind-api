package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// PocketBaseTemplate returns the predefined PocketBase template
func pocketBaseTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "PocketBase",
		DisplayRank: uint(80000),
		Icon:        "pocketbase",
		Keywords:    []string{"pocketbase", "database", "backend", "supabase", "firebase"},
		Description: "Open source backend in 1 file",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the PocketBase instance.",
				Required:    true,
			},
			{
				ID:          2,
				Name:        "Storage Size",
				Type:        schema.InputTypeVolumeSize,
				Description: "Size of the persistent storage for PocketBase data.",
				Required:    true,
				Default:     utils.ToPtr("1Gi"),
			},
		},
		Services: []schema.TemplateService{
			{
				ID:           1,
				Name:         "PocketBase",
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				HostInputIDs: []int{1},
				Image:        utils.ToPtr("ghcr.io/unbindapp/pocketbase:v0.28.1"),
				Ports: []schema.PortSpec{
					{
						Port:     8090,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic: true,
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
				Volumes: []schema.TemplateVolume{
					{
						Name: "pb-data",
						Size: schema.TemplateVolumeSize{
							FromInputID: 2,
						},
						MountPath: "/pb_data",
					},
				},
			},
		},
	}
}
