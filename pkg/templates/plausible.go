package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// PlausibleTemplate returns the predefined Plausible Analytics template
func plausibleTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "Plausible",
		DisplayRank: uint(10000),
		Icon:        "plausible",
		Keywords:    []string{"analytics", "privacy-friendly", "open source", "Google Analytics", "umami"},
		Description: "Privacy-friendly Google Analytics alternative",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the Plausible Analytics instance.",
				Required:    true,
			},
			{
				ID:          2,
				Name:        "PostgreSQL Database Size",
				Type:        schema.InputTypeDatabaseSize,
				Description: "Size of the persistent storage for PostgreSQL database.",
				Required:    true,
				Default:     utils.ToPtr("1Gi"),
			},
			{
				ID:          3,
				Name:        "Clickhouse Database Size",
				Type:        schema.InputTypeDatabaseSize,
				Description: "Size of the persistent storage for Clickhouse database.",
				Required:    true,
				Default:     utils.ToPtr("1Gi"),
			},
		},
		Services: []schema.TemplateService{
			{
				ID:           1,
				Name:         "PostgreSQL",
				InputIDs:     []int{2},
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("postgres"),
			},
			{
				ID:           2,
				Name:         "ClickHouse",
				InputIDs:     []int{3},
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("clickhouse"),
			},
			{
				ID:         3,
				DependsOn:  []int{1, 2},
				InputIDs:   []int{1},
				Name:       "Plausible",
				RunCommand: utils.ToPtr("sh -c \"/entrypoint.sh db createdb && /entrypoint.sh db migrate && /entrypoint.sh run\""),
				Type:       schema.ServiceTypeDockerimage,
				Builder:    schema.ServiceBuilderDocker,
				Image:      utils.ToPtr("ghcr.io/plausible/community-edition:v3.0.1"),
				Ports: []schema.PortSpec{
					{
						Port:     8000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/api/health",
					Port:                      utils.ToPtr(int32(8000)),
					PeriodSeconds:             10,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   5,
					LivenessFailureThreshold:  5,
					ReadinessFailureThreshold: 3,
				},
				Variables: []schema.TemplateVariable{
					{
						Name:  "HTTP_PORT",
						Value: "8000",
					},
					{
						Name: "BASE_URL",
						Generator: &schema.ValueGenerator{
							Type:      schema.GeneratorTypeInput,
							InputID:   1,
							AddPrefix: "https://",
						},
					},
					{
						Name: "SECRET_KEY_BASE",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
							// Plausible wants a 64-byte key
							HashType: utils.ToPtr(schema.ValueHashTypeSHA512),
						},
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					// PostgreSQL references
					{
						SourceID:   1,
						SourceName: "DATABASE_URL",
						TargetName: "DATABASE_URL",
					},
					// ClickHouse references
					{
						SourceID:   2,
						SourceName: "DATABASE_HTTP_URL",
						TargetName: "CLICKHOUSE_DATABASE_URL",
					},
				},
			},
		},
	}
}
