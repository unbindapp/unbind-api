package templates

import (
	"github.com/unbindapp/unbind-api/config"
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

type Templater struct {
	cfg *config.Config
}

func NewTemplater(cfg *config.Config) *Templater {
	return &Templater{
		cfg: cfg,
	}
}
func (self *Templater) AvailableTemplates() []*schema.TemplateDefinition {
	return []*schema.TemplateDefinition{
		wordPressTemplate(),
	}
}

// WordPressTemplate returns the predefined WordPress template
func wordPressTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "wordpress",
		Version:     1,
		Description: "WordPress with MySQL",
		Services: []schema.TemplateService{
			{
				ID:           1,
				Name:         "MySQL",
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("mysql"),
			},
			{
				ID:        2,
				Icon:      "wordpress",
				DependsOn: []int{1},
				Name:      "Wordpress",
				Type:      schema.ServiceTypeDockerimage,
				Builder:   schema.ServiceBuilderDocker,
				Image:     utils.ToPtr("wordpress:6.8"),
				Ports: []schema.PortSpec{
					{
						Port:     80,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic: true,
				Variables: []schema.TemplateVariable{
					{
						Name:  "WORDPRESS_USERNAME",
						Value: "admin",
					},
					{
						Name: "WORDPRESS_PASSWORD",
						// Generate the value in rendering
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
						},
					},
					{
						Name: "WORDPRESS_EMAIL",
						// Generate the value in rendering
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypeEmail,
						},
					},
					{
						Name:  "WORDPRESS_DATABASE_NAME",
						Value: "moco",
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1, // MySQL (now ID 1)
						SourceName: "DATABASE_USERNAME",
						TargetName: "WORDPRESS_DATABASE_USER",
					},
					{
						SourceID:   1, // MySQL (now ID 1)
						SourceName: "DATABASE_PASSWORD",
						TargetName: "WORDPRESS_DATABASE_PASSWORD",
					},
					{
						SourceID:   1, // MySQL (now ID 1)
						TargetName: "WORDPRESS_DATABASE_HOST",
						IsHost:     true,
					},
				},
			},
		},
	}
}
