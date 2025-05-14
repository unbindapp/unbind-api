package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// WireGuardTemplate returns the predefined WireGuard template
func wireGuardTCPTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "WireGuard TCP",
		Icon:        "wireguard",
		Keywords:    []string{"wireguard", "vpn", "tcp tunnel", "udp2raw", "openvpn"},
		Description: "Fast, modern, and open source VPN with a TCP tunnel",
		Version:     1,
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
						Port:     51820,
						Protocol: utils.ToPtr(schema.ProtocolUDP),
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
						Name:  "WG_PORT",
						Value: "51820",
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
					Privileged: utils.ToPtr(true), // Required for iptables manipulation
					Capabilities: &schema.Capabilities{
						Add: []schema.Capability{
							"NET_ADMIN",
							"SYS_MODULE",
						},
					},
				},
			},
			// UDP2RAW Service - TCP tunnel for WireGuard
			{
				ID:      2,
				Name:    "WireGuard TCP Tunnel",
				Type:    schema.ServiceTypeDockerimage,
				Builder: schema.ServiceBuilderDocker,
				Image:   utils.ToPtr("ghcr.io/unbindapp/udp2raw:latest"),
				Ports: []schema.PortSpec{
					{
						IsNodePort:      true,
						InputTemplateID: utils.ToPtr(2),
						Protocol:        utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic: true, // Expose TCP port publicly
				Variables: []schema.TemplateVariable{
					{
						Name:  "UDP2RAW_MODE",
						Value: "server",
					},
					{
						Name:  "UDP2RAW_LOCAL_ADDR",
						Value: "0.0.0.0",
					},
					{
						Name: "UDP2RAW_LOCAL_PORT",
						Generator: &schema.ValueGenerator{
							Type:    schema.GeneratorTypeInput,
							InputID: 2,
						},
					},
					{
						Name:  "UDP2RAW_REMOTE_PORT",
						Value: "51820",
					},
					{
						Name: "UDP2RAW_KEY",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
						},
					},
					{
						Name:  "UDP2RAW_RAW_MODE",
						Value: "easy-faketcp",
					},
					{
						Name:  "UDP2RAW_CIPHER",
						Value: "aes128cbc",
					},
					{
						Name:  "UDP2RAW_AUTH",
						Value: "md5",
					},
					{
						Name:  "UDP2RAW_LOG_LEVEL",
						Value: "info",
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:                1,
						TargetName:              "UDP2RAW_REMOTE_ADDR",
						IsHost:                  true,
						ResolveAsNormalVariable: true,
					},
				},
				SecurityContext: &schema.SecurityContext{
					Privileged: utils.ToPtr(true), // Required for iptables manipulation
					Capabilities: &schema.Capabilities{
						Add: []schema.Capability{
							"NET_ADMIN",
							"NET_RAW",
						},
					},
				},
			},
		},
	}
}
