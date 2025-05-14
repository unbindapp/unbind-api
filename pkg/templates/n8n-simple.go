package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// n8nSimpleTemplate returns the predefined n8n template
func n8nSimpleTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "n8n-simple",
		Version:     1,
		Description: "n8n - Workflow Automation Platform, with Internal Worker",
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the n8n instance.",
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
				DependsOn:    []int{1},
				HostInputIDs: []int{1},
				Name:         "n8n",
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				Image:        utils.ToPtr("n8nio/n8n:1.93.0"),
				Ports: []schema.PortSpec{
					{
						Port:     5678,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic: true,
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/healthz",
					Port:                      utils.ToPtr(int32(5678)),
					PeriodSeconds:             5,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   3,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 3,
				},
				Variables: []schema.TemplateVariable{
					{
						Name: "N8N_HOST",
						Generator: &schema.ValueGenerator{
							Type:    schema.GeneratorTypeInput,
							InputID: 1,
						},
					},
					{
						Name: "N8N_EDITOR_BASE_URL",
						Generator: &schema.ValueGenerator{
							Type:      schema.GeneratorTypeInput,
							InputID:   1,
							AddPrefix: "https://",
						},
					},
					{
						Name: "VUE_APP_URL_BASE_API",
						Generator: &schema.ValueGenerator{
							Type:      schema.GeneratorTypeInput,
							InputID:   1,
							AddPrefix: "https://",
						},
					},
					{
						Name: "WEBHOOK_TUNNEL_URL",
						Generator: &schema.ValueGenerator{
							Type:      schema.GeneratorTypeInput,
							InputID:   1,
							AddPrefix: "https://",
						},
					},
					{
						Name:  "N8N_PORT",
						Value: "5678",
					},
					{
						Name: "N8N_ENCRYPTION_KEY",
						Generator: &schema.ValueGenerator{
							Type:     schema.GeneratorTypePassword,
							HashType: utils.ToPtr(schema.ValueHashTypeSHA256),
						},
					},
					{
						Name:  "N8N_USER_MANAGEMENT_DISABLED",
						Value: "false",
					},
					{
						Name:  "NODE_FUNCTION_ALLOW_EXTERNAL",
						Value: "axios,qs",
					},
					{
						Name:  "N8N_RUNNERS_ENABLED",
						Value: "true",
					},
					{
						Name:  "DB_TYPE",
						Value: "postgresdb",
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_DEFAULT_DB_NAME",
						TargetName: "DB_POSTGRESDB_DATABASE",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_HOST",
						TargetName: "DB_POSTGRESDB_HOST",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_PORT",
						TargetName: "DB_POSTGRESDB_PORT",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_USERNAME",
						TargetName: "DB_POSTGRESDB_USER",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "DB_POSTGRESDB_PASSWORD",
					},
				},
			},
		},
	}
}
