package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// WireGuardTemplate returns the predefined WireGuard template
func wireGuardTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "wireguard",
		Version:     1,
		Description: "WireGuard VPN with web-based management interface.",
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "Wireguard Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the WireGuard instance.",
				Required:    true,
				TargetPort:  utils.ToPtr(51821), // Target TCP port
			},
			{
				ID:          2,
				Name:        "Wireguard TCP NodePort",
				Type:        schema.InputTypeNodePort,
				Description: "NodePort to use for the WireGuard TCP tunnel.",
				Required:    true,
			},
			{
				ID:          3,
				Name:        "Storage Size",
				Type:        schema.InputTypeVolumeSize,
				Description: "Size of the persistent storage for Wireguard config data.",
				Required:    true,
				Default:     utils.ToPtr("1Gi"),
			},
		},
		Services: []schema.TemplateService{
			// WireGuard Service
			{
				ID:           1,
				Name:         "WireGuard",
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				HostInputIDs: []int{1},
				Image:        utils.ToPtr("ghcr.io/wg-easy/wg-easy:14"),
				Ports: []schema.PortSpec{
					{
						IsNodePort:      true,
						InputTemplateID: utils.ToPtr(2),
						Protocol:        utils.ToPtr(schema.ProtocolUDP),
					},
					{
						Port:     51821,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic: true, // Not directly exposed, use TCP proxy instead
				Variables: []schema.TemplateVariable{
					{
						Name: "WG_HOST",
						Generator: &schema.ValueGenerator{
							Type:    schema.GeneratorTypeInput,
							InputID: 1,
						},
					},
					{
						Name: "WG_PORT",
						Generator: &schema.ValueGenerator{
							Type:    schema.GeneratorTypeInput,
							InputID: 2,
						},
					},
					{
						Name: "PASSWORD_HASH",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePasswordBcrypt,
						},
					},
					{
						Name:  "LANG",
						Value: "en",
					},
					{
						Name:  "WG_DEFAULT_DNS",
						Value: "1.1.1.1",
					},
					{
						Name:  "WG_MTU",
						Value: "1420",
					},
					{
						Name:  "UI_TRAFFIC_STATS",
						Value: "true",
					},
					{
						Name:  "UI_CHART_TYPE",
						Value: "1",
					},
				},
				Volumes: []schema.TemplateVolume{
					{
						Name: "wireguard-config",
						Size: schema.TemplateVolumeSize{
							FromInputID: 3,
						},
						MountPath: "/etc/wireguard",
					},
				},
				SecurityContext: &schema.SecurityContext{
					Capabilities: &schema.Capabilities{
						Add: []schema.Capability{
							"NET_ADMIN",
							"SYS_MODULE",
						},
					},
				},
			},
		},
	}
}
