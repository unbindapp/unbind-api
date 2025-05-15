package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// StrapiTemplate returns the predefined Strapi template
func strapiTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "Strapi",
		DisplayRank: uint(65000),
		Icon:        "strapi",
		Keywords:    []string{"cms", "headless cms", "content management system"},
		Description: "Open source headless CMS",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the Strapi instance.",
				Required:    true,
			},
			{
				ID:   2,
				Name: "Upload Storage Size",
				Type: schema.InputTypeVolumeSize,
				Volume: &schema.TemplateVolume{
					Name:      "strapi-upload-data",
					MountPath: "/opt/app/public/uploads",
				},
				Description: "Size of the persistent storage for Strapi uploads.",
				Required:    true,
				Default:     utils.ToPtr("512Mi"),
			},
			{
				ID:   3,
				Name: "Source Code Storage Size",
				Type: schema.InputTypeVolumeSize,
				Volume: &schema.TemplateVolume{
					Name:      "strapi-src",
					MountPath: "/opt/app/src",
				},
				Description: "Persistent volume for Strapi source code (e.g. APIs, components).",
				Required:    true,
				Default:     utils.ToPtr("512Mi"),
			},
			{
				ID:   4,
				Name: "Config Storage Size",
				Type: schema.InputTypeVolumeSize,
				Volume: &schema.TemplateVolume{
					Name:      "strapi-config",
					MountPath: "/opt/app/config",
				},
				Description: "Persistent volume for Strapi configuration files.",
				Required:    true,
				Default:     utils.ToPtr("256Mi"),
			},
			{
				ID:          3,
				Name:        "Database Size",
				Type:        schema.InputTypeDatabaseSize,
				Description: "Size of the persistent storage for PostgreSQL database.",
				Required:    true,
				Default:     utils.ToPtr("1Gi"),
			},
		},
		Services: []schema.TemplateService{
			{
				ID:           1,
				Name:         "PostgreSQL",
				InputIDs:     []int{5},
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("postgres"),
			},
			{
				ID:        2,
				DependsOn: []int{1},
				InputIDs:  []int{1, 2, 3, 4},
				Name:      "Strapi",
				Type:      schema.ServiceTypeDockerimage,
				Builder:   schema.ServiceBuilderDocker,
				Image:     utils.ToPtr("ghcr.io/unbindapp/strapi-development:v5.12.6"),
				Ports: []schema.PortSpec{
					{
						Port:     1337,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type: schema.HealthCheckTypeHTTP,
					Path: "/",
					Port: utils.ToPtr(int32(1337)),
					// High trhesholds because it does some bootstrap stuff
					PeriodSeconds:             10,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   10,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 10,
				},
				Variables: []schema.TemplateVariable{
					{
						Name: "BASE_URL",
						Generator: &schema.ValueGenerator{
							Type:      schema.GeneratorTypeInput,
							InputID:   1,
							AddPrefix: "https://",
						},
					},
					{
						Name:  "DATABASE_CLIENT",
						Value: "postgres",
					},
					{
						Name:  "DATABASE_PORT",
						Value: "5432",
					},
					{
						Name:  "STRAPI_TELEMETRY_DISABLED",
						Value: "true",
					},
					{
						Name:  "NODE_ENV",
						Value: "production",
					},
					{
						Name:  "STRAPI_PLUGIN_I18N_INIT_LOCALE_CODE",
						Value: "en",
					},
					{
						Name:  "STRAPI_ENFORCE_SOURCEMAPS",
						Value: "false",
					},
					{
						Name: "JWT_SECRET",
						Generator: &schema.ValueGenerator{
							Type:     schema.GeneratorTypePassword,
							HashType: utils.ToPtr(schema.ValueHashTypeSHA512),
						},
					},
					{
						Name: "ADMIN_JWT_SECRET",
						Generator: &schema.ValueGenerator{
							Type:     schema.GeneratorTypePassword,
							HashType: utils.ToPtr(schema.ValueHashTypeSHA512),
						},
					},
					{
						Name: "APP_KEYS",
						Generator: &schema.ValueGenerator{
							Type:     schema.GeneratorTypePassword,
							HashType: utils.ToPtr(schema.ValueHashTypeSHA512),
						},
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_USERNAME",
						TargetName: "DATABASE_USERNAME",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "DATABASE_PASSWORD",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_DEFAULT_DB_NAME",
						TargetName: "DATABASE_NAME",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_HOST",
						TargetName: "DATABASE_HOST",
					},
				},
			},
		},
	}
}
