package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// WordPressTemplate returns the predefined WordPress template
func wordPressTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "WordPress",
		DisplayRank: uint(70000),
		Icon:        "wordpress",
		Keywords:    []string{"bloggin", "cms", "content management system", "WooCommerce", "ecommerce", "website", "publishing platform", "php", "mysql"},
		Description: "The open source publishing platform & CMS",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "Domain",
				Type:        schema.InputTypeHost,
				Description: "The domain to use for the WordPress instance.",
				Required:    true,
			},
			{
				ID:          2,
				Name:        "Database Size",
				Type:        schema.InputTypeDatabaseSize,
				Description: "Size of the storage for the MySQL database.",
				Required:    true,
				Default:     utils.ToPtr("1Gi"),
			},
		},
		Services: []schema.TemplateService{
			{
				ID:           1,
				Name:         "MySQL",
				InputIDs:     []int{2},
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("mysql"),
			},
			{
				ID:        2,
				DependsOn: []int{1},
				InputIDs:  []int{1},
				Name:      "WordPress",
				Type:      schema.ServiceTypeDockerimage,
				Builder:   schema.ServiceBuilderDocker,
				Image:     utils.ToPtr("wordpress:6.8"),
				Ports: []schema.PortSpec{
					{
						Port:     80,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/",
					Port:                      utils.ToPtr(int32(80)),
					PeriodSeconds:             2,
					TimeoutSeconds:            10,
					StartupFailureThreshold:   10,
					LivenessFailureThreshold:  10,
					ReadinessFailureThreshold: 3,
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
						SourceName: "DATABASE_DEFAULT_DB_NAME",
						TargetName: "WORDPRESS_DB_NAME",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_HOST",
						TargetName: "WORDPRESS_DB_HOST",
					},
				},
			},
		},
	}
}
