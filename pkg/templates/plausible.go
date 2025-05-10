package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// PlausibleTemplate returns the predefined Plausible Analytics template
func plausibleTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "plausible",
		Version:     1,
		Description: "Plausible Analytics - Simple, privacy-friendly Google Analytics alternative",
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the Plausible Analytics instance.",
				Required:    true,
			},
		},
		Services: []schema.TemplateService{
			{
				ID:           1,
				Name:         "PostgreSQL",
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("postgres"),
			},
			{
				ID:           2,
				Name:         "ClickHouse",
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("clickhouse"),
			},
			{
				ID:           3,
				DependsOn:    []int{1, 2},
				HostInputIDs: []int{1},
				Name:         "Plausible",
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				Image:        utils.ToPtr("plausible/community-edition:latest"),
				Ports: []schema.PortSpec{
					{
						Port:     8000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic: true,
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
						Name: "SECRET_KEY_BASE",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
						},
					},
					{
						Name:  "ADMIN_USER_EMAIL",
						Value: "admin@example.com",
					},
					{
						Name: "ADMIN_USER_NAME",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypeEmail,
						},
					},
					{
						Name: "ADMIN_USER_PWD",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
						},
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					// PostgreSQL references
					{
						SourceID:   1,
						SourceName: "DATABASE_USERNAME",
						TargetName: "POSTGRES_USER",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "POSTGRES_PASSWORD",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_DEFAULT_DB_NAME",
						TargetName: "POSTGRES_DB",
					},
					{
						SourceID:   1,
						TargetName: "POSTGRES_HOST",
						IsHost:     true,
					},
					// ClickHouse references
					{
						SourceID:   2,
						SourceName: "DATABASE_USERNAME",
						TargetName: "CLICKHOUSE_USER",
					},
					{
						SourceID:   2,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "CLICKHOUSE_PASSWORD",
					},
					{
						SourceID:   2,
						SourceName: "DATABASE_DEFAULT_DB_NAME",
						TargetName: "CLICKHOUSE_DB",
					},
					{
						SourceID:   2,
						TargetName: "CLICKHOUSE_HOST",
						IsHost:     true,
					},
				},
			},
		},
	}
}
