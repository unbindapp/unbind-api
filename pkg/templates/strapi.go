package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// StrapiTemplate returns the predefined Strapi template
func strapiTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "Strapi",
		Icon:        "strapi",
		Keywords:    []string{"cms", "headless cms", "content management system"},
		Description: "Open source headless CMS",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the Strapi instance.",
				Required:    true,
			},
			{
				ID:          2,
				Name:        "Storage Size",
				Type:        schema.InputTypeVolumeSize,
				Description: "Size of the persistent storage for Strapi uploads.",
				Required:    true,
				Default:     utils.ToPtr("512Mi"),
			},
		},
		Services: []schema.TemplateService{
			{
				ID:           1,
				Name:         "PostgreSQL",
				Icon:         "postgres",
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("postgres"),
			},
			{
				ID:           2,
				DependsOn:    []int{1},
				HostInputIDs: []int{1},
				Name:         "Strapi",
				Icon:         string(schema.ServiceTypeDockerimage),
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				Image:        utils.ToPtr("elestio/strapi-development:v5.12.6"),
				Ports: []schema.PortSpec{
					{
						Port:     1337,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic: true,
				Variables: []schema.TemplateVariable{
					{
						Name:  "DATABASE_CLIENT",
						Value: "postgres",
					},
					{
						Name:  "DATABASE_PORT",
						Value: "5432",
					},
					{
						Name:  "STRAPI_TELEMETRY_DISABLED",
						Value: "true",
					},
					{
						Name:  "NODE_ENV",
						Value: "production",
					},
					{
						Name:  "STRAPI_PLUGIN_I18N_INIT_LOCALE_CODE",
						Value: "en",
					},
					{
						Name: "JWT_SECRET",
						Generator: &schema.ValueGenerator{
							Type:     schema.GeneratorTypePassword,
							HashType: utils.ToPtr(schema.ValueHashTypeSHA512),
						},
					},
					{
						Name: "ADMIN_JWT_SECRET",
						Generator: &schema.ValueGenerator{
							Type:     schema.GeneratorTypePassword,
							HashType: utils.ToPtr(schema.ValueHashTypeSHA512),
						},
					},
					{
						Name: "APP_KEYS",
						Generator: &schema.ValueGenerator{
							Type:     schema.GeneratorTypePassword,
							HashType: utils.ToPtr(schema.ValueHashTypeSHA512),
						},
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_USERNAME",
						TargetName: "DATABASE_USERNAME",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "DATABASE_PASSWORD",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_DEFAULT_DB_NAME",
						TargetName: "DATABASE_NAME",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_HOST",
						TargetName: "DATABASE_HOST",
					},
				},
				Volumes: []schema.TemplateVolume{
					{
						Name: "strapi-upload-data",
						Size: schema.TemplateVolumeSize{
							FromInputID: 2,
						},
						MountPath: "/opt/app/public/uploads",
					},
				},
			},
		},
	}
}
