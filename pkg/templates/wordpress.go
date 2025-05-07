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
						Name:  "WORDPRESS_DB_NAME",
						Value: "moco",
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_USERNAME",
						TargetName: "WORDPRESS_DB_USER",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "WORDPRESS_DB_PASSWORD",
					},
					{
						SourceID:   1,
						TargetName: "WORDPRESS_DB_HOST",
						IsHost:     true,
					},
				},
			},
		},
	}
}
