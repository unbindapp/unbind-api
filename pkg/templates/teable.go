package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// teableTemplate returns the predefined Teable template
func teableTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "Teable",
		DisplayRank: uint(90500),
		Icon:        "teable",
		Keywords:    []string{"airtable", "teable", "database", "visual", "interface", "relational", "postgresql"},
		Description: "Teable is a powerful visual interface built on relational databases (PostgreSQL).",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          "input_domain",
				Name:        "Domain",
				Type:        schema.InputTypeHost,
				Description: "The domain to use for the Teable instance.",
				Required:    true,
			},
			{
				ID:          "input_database_size",
				Name:        "Database Size",
				Type:        schema.InputTypeDatabaseSize,
				Description: "Size of the storage for the PostgreSQL database.",
				Required:    true,
				Default:     utils.ToPtr("1"),
			},
			{
				ID:          "input_teable_size",
				Name:        "Storage Size",
				Type:        schema.InputTypeVolumeSize,
				Description: "Size of the storage for Teable assets.",
				Required:    true,
				Default:     utils.ToPtr("1"),
				Volume: &schema.TemplateVolume{
					Name:      "teable-data",
					MountPath: "/app/.assets",
				},
			},
		},
		Services: []schema.TemplateService{
			{
				ID:           "service_postgres",
				Name:         "PostgreSQL",
				InputIDs:     []string{"input_database_size"},
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("postgres"),
			},
			{
				ID:           "service_redis",
				Name:         "Redis",
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				DatabaseType: utils.ToPtr("redis"),
				DatabaseConfig: &schema.DatabaseConfig{
					StorageSize: "0.25",
				},
			},
			{
				ID:        "service_teable",
				Name:      "Teable",
				InputIDs:  []string{"input_domain", "input_teable_size", "input_postgres_password", "input_redis_password", "input_secret_key"},
				Type:      schema.ServiceTypeDockerimage,
				Builder:   schema.ServiceBuilderDocker,
				Image:     utils.ToPtr("ghcr.io/teableio/teable:latest"),
				DependsOn: []string{"service_postgres", "service_redis"},
				Volumes: []schema.TemplateVolume{
					{
						Name:      "teable_data",
						MountPath: "/app/.assets",
					},
				},
				Ports: []schema.PortSpec{
					{
						Port:     3000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				Variables: []schema.TemplateVariable{
					{
						Name:  "BACKEND_CACHE_PROVIDER",
						Value: "redis",
					},
					{
						Name: "PUBLIC_ORIGIN",
						Generator: &schema.ValueGenerator{
							Type:      schema.GeneratorTypeInput,
							InputID:   "input_domain",
							AddPrefix: "https://",
						},
					},
					{
						Name: "SECRET_KEY",
						Generator: &schema.ValueGenerator{
							Type:     schema.GeneratorTypePassword,
							HashType: utils.ToPtr(schema.ValueHashTypeSHA256),
						},
					},
					{
						Name:  "PORT",
						Value: "3000",
					},
					{
						Name:  "NEXT_ENV_IMAGES_ALL_REMOTE",
						Value: "true",
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   "service_redis",
						SourceName: "DATABASE_URL",
						TargetName: "BACKEND_CACHE_REDIS_URI",
					},
					{
						SourceID:   "service_postgres",
						SourceName: "DATABASE_URL",
						TargetName: "PRISMA_DATABASE_URL",
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/health",
					Port:                      utils.ToPtr(int32(3000)),
					PeriodSeconds:             5,
					TimeoutSeconds:            20,
					StartupFailureThreshold:   10,
					LivenessFailureThreshold:  10,
					ReadinessFailureThreshold: 10,
				},
			},
		},
	}
}
