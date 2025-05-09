package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// GhostTemplate returns the predefined Ghost template
func ghostTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "ghost",
		Version:     1,
		Description: "Ghost CMS with MySQL",
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the Ghost instance.",
				Required:    true,
			},
		},
		Services: []schema.TemplateService{
			{
				ID:           1,
				Name:         "MySQL",
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("mysql"),
			},
			{
				ID:           2,
				DependsOn:    []int{1},
				HostInputIDs: []int{1},
				Name:         "Ghost",
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				Image:        utils.ToPtr("ghost:5"),
				Ports: []schema.PortSpec{
					{
						Port:     2368,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic: true,
				Variables: []schema.TemplateVariable{
					{
						Name: "url",
						Generator: &schema.ValueGenerator{
							Type:      schema.GeneratorTypeInput,
							InputID:   1,
							AddPrefix: "https://",
						},
					},
					{
						Name:  "database__client",
						Value: "mysql",
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_USERNAME",
						TargetName: "database__connection__user",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "database__connection__password",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_DEFAULT_DB_NAME",
						TargetName: "database__connection__database",
					},
					{
						SourceID:   1,
						TargetName: "database__connection__host",
						IsHost:     true,
					},
				},
			},
		},
	}
}
