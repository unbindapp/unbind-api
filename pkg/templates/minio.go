package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// MinioTemplate returns the predefined MinIO template
func minioTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "MinIO",
		DisplayRank: uint(60000),
		Icon:        "minio",
		Keywords:    []string{"object storage", "file storage", "s3", "s3 compatible"},
		Description: "S3-compatible object storage",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "API Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the MinIO API server.",
				Required:    true,
				TargetPort:  utils.ToPtr(9000),
			},
			{
				ID:          2,
				Name:        "Console Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the MinIO Console.",
				Required:    true,
				TargetPort:  utils.ToPtr(9001),
			},
			{
				ID:          3,
				Name:        "Storage Size",
				Type:        schema.InputTypeVolumeSize,
				Description: "Size of the persistent storage for MinIO data.",
				Required:    true,
				Default:     utils.ToPtr("1Gi"),
			},
		},
		Services: []schema.TemplateService{
			{
				ID:           1,
				Name:         "MinIO",
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				HostInputIDs: []int{1, 2},
				Image:        utils.ToPtr("minio/minio:latest"),
				Ports: []schema.PortSpec{
					{
						Port:     9000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
					{
						Port:     9001,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic:   true,
				RunCommand: utils.ToPtr("minio server /data --console-address ':9001'"),
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeExec,
					Command:                   "mc ready local",
					PeriodSeconds:             5,
					TimeoutSeconds:            20,
					StartupFailureThreshold:   10,
					LivenessFailureThreshold:  10,
					ReadinessFailureThreshold: 3,
				},
				Variables: []schema.TemplateVariable{
					{
						Name:  "MINIO_ROOT_USER",
						Value: "minioadmin",
					},
					{
						Name: "MINIO_ROOT_PASSWORD",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
						},
					},
					{
						Name: "MINIO_SERVER_URL",
						Generator: &schema.ValueGenerator{
							Type:      schema.GeneratorTypeInput,
							InputID:   1,
							AddPrefix: "https://",
						},
					},
					{
						Name: "MINIO_BROWSER_REDIRECT_URL",
						Generator: &schema.ValueGenerator{
							Type:      schema.GeneratorTypeInput,
							InputID:   2,
							AddPrefix: "https://",
						},
					},
				},
				Volumes: []schema.TemplateVolume{
					{
						Name: "minio-data",
						Size: schema.TemplateVolumeSize{
							FromInputID: 3,
						},
						MountPath: "/data",
					},
				},
			},
		},
	}
}
