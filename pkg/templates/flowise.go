package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// FlowiseTemplate returns the predefined Flowise template
func flowiseTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "Flowise",
		DisplayRank: uint(90000),
		Icon:        "flowise",
		Keywords:    []string{"llm", "ai", "chatbot", "langchain", "flow", "workflow", "automation", "low code", "low-code", "no code", "no-code", "chatbot", "ai"},
		Description: "Low code tool for building LLM flows",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the Flowise instance.",
				Required:    true,
			},
			{
				ID:   2,
				Name: "Storage Size",
				Type: schema.InputTypeVolumeSize,
				Volume: &schema.TemplateVolume{
					Name:      "flowise-data",
					MountPath: "/root/.flowise",
				},
				Description: "Size of the persistent storage for Flowise data.",
				Required:    true,
				Default:     utils.ToPtr("1Gi"),
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
				InputIDs:     []int{3},
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("postgres"),
			},
			{
				ID:         2,
				DependsOn:  []int{1},
				InputIDs:   []int{1, 2},
				Name:       "Flowise",
				Type:       schema.ServiceTypeDockerimage,
				Builder:    schema.ServiceBuilderDocker,
				Image:      utils.ToPtr("flowiseai/flowise:2.2.8"),
				RunCommand: utils.ToPtr("flowise start"),
				Ports: []schema.PortSpec{
					{
						Port:     3000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/api/v1/ping",
					Port:                      utils.ToPtr(int32(3000)),
					PeriodSeconds:             5,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   3,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 3,
				},
				Variables: []schema.TemplateVariable{
					{
						Name:  "DEBUG",
						Value: "false",
					},
					{
						Name:  "PORT",
						Value: "3000",
					},
					{
						Name:  "FLOWISE_USERNAME",
						Value: "admin",
					},
					{
						Name: "FLOWISE_PASSWORD",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
						},
					},
					{
						Name:  "APIKEY_PATH",
						Value: "/root/.flowise",
					},
					{
						Name:  "SECRETKEY_PATH",
						Value: "/root/.flowise",
					},
					{
						Name:  "LOG_LEVEL",
						Value: "info",
					},
					{
						Name:  "LOG_PATH",
						Value: "/root/.flowise/logs",
					},
					{
						Name:  "DATABASE_TYPE",
						Value: "postgres",
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_HOST",
						TargetName: "DATABASE_HOST",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_PORT",
						TargetName: "DATABASE_PORT",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_DEFAULT_DB_NAME",
						TargetName: "DATABASE_NAME",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_USERNAME",
						TargetName: "DATABASE_USER",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "DATABASE_PASSWORD",
					},
				},
			},
		},
	}
}
