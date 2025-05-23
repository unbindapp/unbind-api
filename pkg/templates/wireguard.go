package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// WireGuardTemplate returns the predefined WireGuard template
func wireGuardTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "WireGuard",
		DisplayRank: uint(110000),
		Icon:        "wireguard",
		Keywords:    []string{"wireguard", "vpn", "tcp tunnel", "udp2raw", "openvpn"},
		Description: "Fast, modern, and open source VPN.",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          "input_domain",
				Name:        "Domain",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the WireGuard instance.",
				Required:    true,
				TargetPort:  utils.ToPtr(51821), // Target TCP port
			},
			{
				ID:           "input_nodeport",
				Name:         "NodePort",
				Type:         schema.InputTypeGeneratedNodePort,
				PortProtocol: utils.ToPtr(schema.ProtocolUDP),
				Description:  "The NodePort to use for the WireGuard TCP tunnel.",
				Hidden:       true,
			},
			{
				ID:   "input_storage_size",
				Name: "Storage Size",
				Type: schema.InputTypeVolumeSize,
				Volume: &schema.TemplateVolume{
					Name:      "wireguard-config",
					MountPath: "/etc/wireguard",
				},
				Description: "Size of the storage for the WireGuard config data.",
				Required:    true,
				Default:     utils.ToPtr("1"),
			},
		},
		Services: []schema.TemplateService{
			// WireGuard Service
			{
				ID:       "service_wireguard",
				Name:     "WireGuard",
				Type:     schema.ServiceTypeDockerimage,
				Builder:  schema.ServiceBuilderDocker,
				InputIDs: []string{"input_domain", "input_nodeport", "input_storage_size"},
				Image:    utils.ToPtr("ghcr.io/wg-easy/wg-easy:14"),
				Ports: []schema.PortSpec{
					{
						Port:     51821,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				Variables: []schema.TemplateVariable{
					{
						Name: "WG_HOST",
						Generator: &schema.ValueGenerator{
							Type:    schema.GeneratorTypeInput,
							InputID: "input_domain",
						},
					},
					{
						Name: "WG_PORT",
						Generator: &schema.ValueGenerator{
							Type:    schema.GeneratorTypeInput,
							InputID: "input_nodeport",
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
