package system_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

type SystemMeta struct {
	ExternalIPV6 string `json:"external_ipv6" nullable:"false"`
	ExternalIPV4 string `json:"external_ipv4" nullable:"false"`
}

type SystemMetaResponse struct {
	Body struct {
		Data *SystemMeta `json:"data" nullable:"false"`
	}
}

func (self *HandlerGroup) GetSystemInformation(ctx context.Context, input *server.BaseAuthInput) (*SystemMetaResponse, error) {
	// Get k8s IPs for load balancer server
	ips, err := self.srv.KubeClient.GetIngressNginxIP(ctx)
	if err != nil {
		log.Error("Error getting ingress nginx IP", "err", err)
		return nil, huma.Error500InternalServerError("Error getting ingress nginx IP")
	}

	// Get system meta
	meta := &SystemMeta{
		ExternalIPV6: ips.IPv6,
		ExternalIPV4: ips.IPv4,
	}

	return &SystemMetaResponse{
		Body: struct {
			Data *SystemMeta `json:"data" nullable:"false"`
		}{Data: meta},
	}, nil
}
