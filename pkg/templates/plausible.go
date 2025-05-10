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
				ID:              1,
				Name:            "PostgreSQL",
				Type:            schema.ServiceTypeDatabase,
				Builder:         schema.ServiceBuilderDatabase,
				DatabaseType:    utils.ToPtr("postgres"),
				DatabaseVersion: utils.ToPtr("16"),
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
				RunCommand:   utils.ToPtr("sh -c \"/entrypoint.sh db createdb && /entrypoint.sh db migrate && /entrypoint.sh run\""),
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				Image:        utils.ToPtr("ghcr.io/plausible/community-edition:v3.0.1"),
				Ports: []schema.PortSpec{
					{
						Port:     8000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic: true,
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
