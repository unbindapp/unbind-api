package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// n8nWorkersTemplate returns the predefined n8n template with Redis queue mode and an external worker
func n8nWorkersTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "n8n-workers",
		Version:     1,
		Description: "n8n - Workflow Automation Platform (queue mode with external worker & Redis)",
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
			// PostgreSQL service (workflow metadata & credentials store)
			{
				ID:           1,
				Name:         "PostgreSQL",
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("postgres"),
			},
			// Redis service (queue backend)
			{
				ID:           2,
				Name:         "Redis",
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("redis"),
			},
			// Main n8n process (API / UI)
			{
				ID:           3,
				DependsOn:    []int{1, 2},
				HostInputIDs: []int{1},
				Name:         "n8n-main",
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
					// Host related variables
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
					// Core configuration
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
					// Enable queue mode
					{
						Name:  "EXECUTIONS_MODE",
						Value: "queue",
					},
					// Route manual executions to workers as well (optional but recommended)
					{
						Name:  "OFFLOAD_MANUAL_EXECUTIONS_TO_WORKERS",
						Value: "true",
					},
					// Redis queue connection (host & port resolved via variable references below)
					{
						Name:  "QUEUE_BULL_REDIS_DB",
						Value: "0",
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
						Name:  "OFFLOAD_MANUAL_EXECUTIONS_TO_WORKERS",
						Value: "true",
					},
					{
						Name:  "DB_TYPE",
						Value: "postgresdb",
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					// Postgres references (same as the simple template)
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
					// Redis references
					{
						SourceID:   2,
						SourceName: "DATABASE_HOST",
						TargetName: "QUEUE_BULL_REDIS_HOST",
					},
					{
						SourceID:   2,
						SourceName: "DATABASE_PORT",
						TargetName: "QUEUE_BULL_REDIS_PORT",
					},
					{
						SourceID:   2,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "QUEUE_BULL_REDIS_PASSWORD",
					},
				},
			},
			// External n8n worker (executes jobs from the queue)
			{
				ID:         4,
				DependsOn:  []int{1, 2, 3},
				Name:       "n8n-worker",
				Type:       schema.ServiceTypeDockerimage,
				Builder:    schema.ServiceBuilderDocker,
				Image:      utils.ToPtr("n8nio/n8n:1.93.0"),
				RunCommand: utils.ToPtr("n8n worker"),
				IsPublic:   false,
				Ports: []schema.PortSpec{
					{
						Port:     5679,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/healthz",
					Port:                      utils.ToPtr(int32(5679)),
					PeriodSeconds:             30,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   5,
					LivenessFailureThreshold:  5,
					ReadinessFailureThreshold: 3,
				},
				Variables: []schema.TemplateVariable{
					// Queue mode for worker
					{
						Name:  "EXECUTIONS_MODE",
						Value: "queue",
					},
					{
						Name:  "N8N_WORKER_ID",
						Value: "worker-1",
					},
					{
						Name:  "QUEUE_HEALTH_CHECK_ACTIVE",
						Value: "true",
					},
					{
						Name:  "QUEUE_HEALTH_CHECK_PORT",
						Value: "5679",
					},
					{
						Name:  "N8N_RUNNERS_ENABLED",
						Value: "true",
					},
					{
						Name:  "OFFLOAD_MANUAL_EXECUTIONS_TO_WORKERS",
						Value: "true",
					},
					{
						Name:  "DB_TYPE",
						Value: "postgresdb",
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					// Postgres references (same as main)
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
					// Redis references
					{
						SourceID:   2,
						SourceName: "DATABASE_HOST",
						TargetName: "QUEUE_BULL_REDIS_HOST",
					},
					{
						SourceID:   2,
						SourceName: "DATABASE_PORT",
						TargetName: "QUEUE_BULL_REDIS_PORT",
					},
					{
						SourceID:   2,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "QUEUE_BULL_REDIS_PASSWORD",
					},
					// N8N references
					{
						SourceID:   3,
						SourceName: "N8N_ENCRYPTION_KEY",
						TargetName: "N8N_ENCRYPTION_KEY",
					},
				},
			},
		},
	}
}
